package main

import (
	"database/sql"

	"github.com/dmitrymomot/go-env"
	_ "github.com/joho/godotenv/autoload" // Load .env file automatically
	_ "github.com/lib/pq"                 // init pg driver
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
)

var (
	// DB
	dbConnString    = env.MustString("DATABASE_URL")
	migrationsDir   = env.GetString("DATABASE_MIGRATIONS_DIR", "./internal/repository/sql/migrations")
	migrationsTable = env.GetString("DATABASE_MIGRATIONS_TABLE", "migrations")

	// Build tag is set up while deployment
	buildTag        = "undefined"
	buildTagRuntime = env.GetString("COMMIT_HASH", buildTag)
)

func main() {
	// Init logger
	logrus.SetReportCaller(false)
	logger := logrus.WithFields(logrus.Fields{
		"app":       "db-migrate",
		"build_tag": buildTagRuntime,
	})
	logger.Logger.SetLevel(logrus.InfoLevel)

	// Init db connection
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		logger.WithError(err).Fatal("Failed to init db connection")
	}
	defer db.Close()

	// check db connection
	if err := db.Ping(); err != nil {
		logger.WithError(err).Fatal("Failed to ping db")
	}

	m := migrate.MigrationSet{
		TableName: migrationsTable,
	}
	migrations := &migrate.FileMigrationSource{
		Dir: migrationsDir,
	}
	n, err := m.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		logger.WithError(err).Fatal("Failed to apply migrations")
	}

	logger.Infof("Applied %d migrations!", n)
}
