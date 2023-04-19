package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dmitrymomot/go-server/pkg/httpserver"
	"github.com/dmitrymomot/go-utils"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.WithFields(logrus.Fields{
		"app": "go-server",
	})
	defer func() { logger.Info("Server successfully shutdown") }()

	// Create a context with a timeout and set the Server's context
	ctx, cancel := utils.NewContextWithCancel(logger)
	defer cancel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "Hello, World!") })
	server := httpserver.NewServer(":8080", handler,
		httpserver.WithShutdownTimeout(5*time.Second),
		httpserver.WithLogger(logger),
	)

	// Run the server
	if err := server.Run(ctx); err != nil {
		logger.Errorf("Server returned an error: %v", err)
	}

	// Wait for the server to shutdown
	logger.Info("Waiting for server shutdown")
	<-ctx.Done()
	logger.Info("Server shutdown")
}
