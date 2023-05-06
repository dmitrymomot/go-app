//go:build mage
// +build mage

package main

import (
	"fmt"
	"time"

	"github.com/dmitrymomot/go-env"
	_ "github.com/joho/godotenv/autoload" // Load .env file automatically

	"github.com/fatih/color"
	"github.com/magefile/mage/sh"
)

// mg contains helpful utility functions, like Deps

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// Run runs the application with environment variables
func Run() error {
	// clean up go build cache
	color.Yellow("Cleaning up go build cache...")
	sh.RunV("go", "clean", "-cache")

	// run the application
	color.Cyan(fmt.Sprintf("Running the application on http://localhost:%d", env.GetInt("HTTP_PORT", 8080)))
	return sh.RunV("go", "run", "./cmd/app/")
}

// PrepareMigration prepares the database migration
func PrepareMigration() error {
	color.Cyan("Preparing database migration...")
	return sh.RunV("./scripts/prepare-migrations.sh")
}

// MigrateUp runs the database migrations up
func MigrateUp() error {
	if err := PrepareMigration(); err != nil {
		return fmt.Errorf("failed to prepare migration: %w", err)
	}

	color.Cyan("Running database migrations...")
	return sh.RunV("go", "run", "./cmd/migrate/main.go")
}

// Up runs the application and the database migrations up
func Up() error {
	color.Cyan("Starting the database...")
	if err := sh.RunV("docker-compose", "-f deployments/docker-compose.yml", "up", "-d"); err != nil {
		return err
	}

	color.Yellow("Waiting for the database to start...")
	time.Sleep(5 * time.Second)

	if err := MigrateUp(); err != nil {
		return err
	}

	return Run()
}

// Down stops the application and the database
func Down() error {
	color.Yellow("Stopping the database and removing all data...")
	return sh.RunV("docker-compose", "-f deployments/docker-compose.yml", "down", "--volumes", "--rmi=local")
}
