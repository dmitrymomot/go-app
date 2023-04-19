package httpserver

import "time"

// WithShutdownTimeout sets the timeout for graceful shutdown.
// If the timeout is exceeded, the server will be shutdown immediately.
// Default is 5 seconds.
func WithShutdownTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = d
	}
}

// WithLogger sets the logger for the server.
// If no logger is set, no log messages will be printed.
func WithLogger(logger Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}
