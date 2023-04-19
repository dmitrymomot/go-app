package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

type (
	// Logger is an interface for logging.
	Logger interface {
		Printf(format string, v ...interface{})
	}

	// http.Server wrapper to handle graceful shutdown
	Server struct {
		Server          *http.Server
		shutdownTimeout time.Duration
		logger          Logger
		shutdownC       chan struct{}
		once            sync.Once
	}

	// Option is a function that configures the Server.
	Option func(*Server)
)

// NewServer creates a new Server instance.
func NewServer(addr string, handler http.Handler, opts ...Option) *Server {
	s := &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		shutdownTimeout: 5 * time.Second,
		shutdownC:       make(chan struct{}),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// log messages if a logger is configured
func (s *Server) log(format string, v ...interface{}) {
	if s.logger != nil {
		s.logger.Printf(format, v...)
	}
}

// Run starts the HTTP server and listens for the operating system interrupt signal.
// If the interrupt signal is received, the HTTP server is gracefully shutdown.
func (s *Server) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	// Start the HTTP server
	g.Go(func() error {
		s.log("Starting HTTP server on %s", s.Server.Addr)
		if err := s.Server.ListenAndServe(); err != http.ErrServerClosed {
			return fmt.Errorf("failed to start HTTP server: %w", err)
		}
		return nil
	})

	// Listen for the operating system interrupt signal (e.g., Ctrl+C)
	g.Go(func() error {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

		select {
		case <-ctx.Done():
			s.log("Context cancelled, shutting down...")
		case <-sigint:
			s.log("Received interrupt signal, shutting down...")
		case <-s.shutdownC:
			s.log("Shutdown requested, shutting down...")
		}

		// Perform graceful shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer shutdownCancel()

		if err := s.Server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("failed to gracefully shutdown HTTP server: %w", err)
		}
		return nil
	})

	// Wait for all goroutines to complete or return an error
	return g.Wait()
}

// Shutdown gracefully shuts down the HTTP server without interrupt signal.
// It cancels the context and waits for shutdownTimeout.
func (s *Server) Shutdown() {
	s.once.Do(func() {
		close(s.shutdownC)
	})
}
