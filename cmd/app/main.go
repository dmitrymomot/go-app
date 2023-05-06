package main

import (
	"database/sql"
	"fmt"

	"github.com/dmitrymomot/go-pkg/httpserver"
	"github.com/dmitrymomot/go-utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func main() {
	logger := logrus.WithFields(logrus.Fields{
		"app":       appName,
		"build_tag": buildTag,
		"component": "main",
	})
	defer func() { logger.Info("Server successfully shutdown") }()

	// Init DB connection
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		logger.WithError(err).Fatal("Failed to init db connection")
	}
	defer db.Close()

	db.SetMaxOpenConns(dbMaxOpenConns)
	db.SetMaxIdleConns(dbMaxIdleConns)

	if err := db.Ping(); err != nil {
		logger.WithError(err).Fatal("Failed to ping db")
	}

	// Create a context with a timeout and set the Server's context
	ctx, cancel := utils.NewContextWithCancel(logger.WithField("component", "context"))
	defer cancel()

	// Create a new errgroup
	eg, _ := errgroup.WithContext(ctx)

	// Init router with default middlewares and routes
	r := initRouter()

	// TODO: Add your routes here

	// Create a new server
	server := httpserver.NewServer(
		fmt.Sprintf(":%d", httpPort), r,
		httpserver.WithShutdownTimeout(httpShutdownTimeout),
		httpserver.WithLogger(logger.WithField("component", "http-server")),
	)
	defer server.Shutdown()

	// Run the server
	eg.Go(func() error { return server.Run(ctx) })

	// Wait for the server to finish
	if err := eg.Wait(); err != nil {
		logger.WithError(err).Error("Server stopped with error")
	}
}
