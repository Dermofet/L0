// Package http provides functionality for creating and running an HTTP server.
package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"

	"L0/internal/cache"
)

// RequestTimeOut defines the timeout duration for HTTP requests.
const RequestTimeOut = 30 * time.Second

// Server represents an HTTP server.
type Server interface {
	// Run starts the HTTP server and listens for incoming requests.
	// It takes a context as an input parameter.
	// Returns an error if the server fails to start.
	Run(ctx context.Context) error

	// Shutdown gracefully shuts down the HTTP server.
	// It takes a context as an input parameter.
	// Returns an error if the server fails to shut down gracefully.
	Shutdown(ctx context.Context) error
}

// server represents an HTTP server instance.
type server struct {
	server *http.Server
	db     *sqlx.DB
	logger *zap.Logger
}

// NewServer creates a new instance of the HTTP server.
// It takes the server address, database connection, logger, and cache as input parameters.
// Returns the HTTP server instance.
func NewServer(
	addr string,
	db *sqlx.DB,
	logger *zap.Logger,
	cache cache.Cache,
	connect stan.Conn,
	subject string,
) *server {
	s := &server{
		db:     db,
		logger: logger,
	}

	r := NewRouter(db, logger, cache, connect, subject)
	err := r.Init()
	if err != nil {
		s.logger.Error("can't init router:", zap.Error(err))
		return nil
	}

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           r.router,
		ReadHeaderTimeout: RequestTimeOut,
	}
	s.server = httpServer

	return s
}

// Run starts the HTTP server and listens for incoming requests.
// It takes a context as an input parameter.
// Returns an error if the server fails to start.
func (s *server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		err := s.server.Shutdown(ctx)
		if err != nil {
			s.logger.Error("can't shutdown http-server", zap.Error(err))
			return
		}
	}()

	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server.
// It takes a context as an input parameter.
// Returns an error if the server fails to shut down gracefully.
func (s *server) Shutdown(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("http server shutdown error: %w", err)
	}
	return nil
}
