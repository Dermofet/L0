package app

import (
	"L0/cmd/L0/config"
	"L0/internal/api/http"
	"L0/internal/cache"
	"L0/internal/db"
	"L0/internal/nats"
	"L0/internal/repository"
	"context"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

// App represents the main application.
type App struct {
	config     *config.Config
	dbConn     *sqlx.DB
	logger     *zap.Logger
	httpServer http.Server
	cache      cache.Cache
}

// NewApp creates a new instance of the application.
func NewApp(cfg *config.Config, logger *zap.Logger) *App {
	return &App{
		config: cfg,
		logger: logger,
		cache:  cache.NewCache(),
	}
}

// Start starts the application.
func (a *App) Start(ctx context.Context) {
	appCtx, cancelApp := context.WithCancel(ctx)
	logger := a.logger
	defer func() {
		if e := recover(); e != nil {
			logger.Panic("application shutdown", zap.Error(fmt.Errorf("%s", e)))
			cancelApp()
		}
	}()

	// Initialize the database
	dbConn, err := a.initDb(appCtx,
		a.config.DB.Host,
		a.config.DB.Port,
		a.config.DB.Name,
		a.config.DB.Username,
		a.config.DB.Password,
		a.config.DB.SSLMode,
	)
	if err != nil {
		logger.Fatal("init db error", zap.Error(err))
	}
	a.dbConn = dbConn

	// Start database migrations
	err = a.startMigrate(appCtx, migrationsPath, a.config.DB.Name, a.dbConn)
	if err != nil {
		logger.Error("db migration error", zap.Error(err))
	}

	wg := &sync.WaitGroup{}

	// Start HTTP server
	wg.Add(1)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				logger.Panic("http start panic", zap.Error(fmt.Errorf("%s", e)))
			}
			wg.Done()
		}()

		addr := fmt.Sprintf("%s:%d", a.config.HttpServer.Host, a.config.HttpServer.Port)
		conn, err := stan.Connect(
			a.config.Nats.ClusterID,
			a.config.Nats.Client2ID,
			stan.NatsURL(fmt.Sprintf("nats://%s:%d", a.config.Nats.Host, a.config.Nats.Port)),
		)
		if err != nil {
			logger.Error("NATS connection error", zap.Error(err))
			return
		}
		a.httpServer = http.NewServer(addr, a.dbConn, logger, a.cache, conn, a.config.Nats.Subject)
		if a.httpServer == nil {
			cancelApp()
			logger.Fatal("can't create http server")
			return
		}

		err = a.httpServer.Run(appCtx)
		cancelApp()
		if err != nil {
			logger.Error("can't start http server", zap.Error(err))
			return
		}
	}()

	// Initialize order repository
	orderRepository := repository.NewOrderRepository(db.NewSource(a.dbConn))

	// Load cache
	if err := a.cache.Load(appCtx, orderRepository); err != nil {
		logger.Error("can't load cache", zap.Error(err))
	}
	logger.Info("load cache")

	// Start NATS subscription
	wg.Add(1)
	go func() {
		defer wg.Done()

		conn, err := stan.Connect(
			a.config.Nats.ClusterID,
			a.config.Nats.Client1ID,
			stan.NatsURL(fmt.Sprintf("nats://%s:%d", a.config.Nats.Host, a.config.Nats.Port)),
		)
		if err != nil {
			logger.Error("NATS connection error", zap.Error(err))
			return
		}

		natsService := nats.NewNatsService(
			orderRepository,
			a.cache,
			conn,
			a.config.Nats.Subject,
		)
		err = natsService.Subscribe(
			context.Background(),
		)
		if err != nil {
			a.logger.Error("NATS subscription error", zap.Error(err))
		}
	}()
}

// GracefulShutdown performs a graceful shutdown of the application.
func (a *App) GracefulShutdown(ctx context.Context) error {
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("can't shutdown http-server: %w", err)
	}
	err = a.dbConn.Close()
	if err != nil {
		return fmt.Errorf("can't shutdown db: %w", err)
	}
	return nil
}

// initDb initializes the database.
func (a *App) initDb(
	ctx context.Context,
	host string,
	port int,
	dbName string,
	user string,
	password string,
	sslmode string,
) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(
		ctx,
		"postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbName, sslmode),
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
