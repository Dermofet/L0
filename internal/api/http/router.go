// Package http provides functionality for handling HTTP requests.
package http

import (
	"fmt"
	"net/http"
	"time"

	"L0/internal/api/http/handlers"
	"L0/internal/cache"
	"L0/internal/db"
	"L0/internal/nats"
	"L0/internal/repository"
	"L0/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

// routerHandlers contains handlers for router.
type routerHandlers struct {
	orderHandlers handlers.OrderHandlers
}

// router represents an HTTP router.
type router struct {
	router   *gin.Engine
	db       *sqlx.DB
	handlers routerHandlers
	logger   *zap.Logger
	cache    cache.Cache
	connect  stan.Conn
	subject  string
}

// NewRouter creates a new instance of HTTP router.
func NewRouter(db *sqlx.DB, logger *zap.Logger, cache cache.Cache, connect stan.Conn, subject string) *router {
	return &router{
		router:  gin.New(),
		db:      db,
		logger:  logger,
		cache:   cache,
		connect: connect,
		subject: subject,
	}
}

// Init initializes the HTTP router.
func (r *router) Init() error {
	r.router.Use(
		gin.Logger(),
		gin.CustomRecovery(r.recovery),
	)
	err := r.registerRoutes()
	if err != nil {
		return fmt.Errorf("can't init router: %w", err)
	}

	return nil
}

// recovery recovers from panics in HTTP handlers.
func (r *router) recovery(c *gin.Context, recovered interface{}) {
	defer func() {
		if e := recover(); e != nil {
			r.logger.Fatal("http server panic", zap.Any("panic", recovered))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()
}

// registerRoutes registers routes in the HTTP router.
func (r *router) registerRoutes() error {
	r.router.NoMethod(handlers.NotImplementedHandler)
	r.router.NoRoute(handlers.NotImplementedHandler)

	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})

	r.router.Use(corsMiddleware)
	r.router.Static("/static", "/backend/internal/static")
	r.router.LoadHTMLFiles("/backend/internal/templates/order.html")

	pgSource := db.NewSource(r.db)
	orderRepository := repository.NewOrderRepository(pgSource)
	orderInteractor := usecase.NewOrderInteractor(orderRepository, r.cache)
	natsService := nats.NewNatsService(
		orderRepository,
		r.cache,
		r.connect,
		r.subject,
	)
	r.handlers.orderHandlers = handlers.NewOrderHandlers(orderInteractor, natsService)

	orderGroup := r.router.Group("/orders")
	orderGroup.GET("/", r.handlers.orderHandlers.GetHTMLOrderHandler)
	orderGroup.GET("/id/:uid", r.handlers.orderHandlers.GetByIdHandler)
	orderGroup.GET("/all", r.handlers.orderHandlers.GetAllHandler)
	orderGroup.POST("/new", r.handlers.orderHandlers.CreateHandler)
	orderGroup.DELETE("/id/:uid", r.handlers.orderHandlers.DeleteHandler)

	return nil
}
