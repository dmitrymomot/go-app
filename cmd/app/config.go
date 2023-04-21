package main

import (
	"time"

	"github.com/dmitrymomot/go-env"
	_ "github.com/joho/godotenv/autoload" // Load .env file automatically
)

// Application glabal variables
var (
	appName     = env.GetString("APP_NAME", "go-server")
	appDebug    = env.GetBool("APP_DEBUG", false)
	appLogLevel = env.GetString("APP_LOG_LEVEL", "info")
	buildTag    = env.GetString("COMMIT_HASH", "undefined")

	// HTTP server
	httpPort            = env.GetInt("HTTP_PORT", 8080)
	httpRequestTimeout  = env.GetDuration("HTTP_REQUEST_TIMEOUT", 10*time.Second)
	httpShutdownTimeout = env.GetDuration("HTTP_SERVER_SHUTDOWN_TIMEOUT", 15*time.Second)
	allowContentTypes   = env.GetStrings("ALLOW_CONTENT_TYPES", ",", []string{"application/json"})

	// CORS
	corsAllowedOrigins     = env.GetStrings("CORS_ALLOWED_ORIGINS", ",", []string{"*"})
	corsAllowedMethods     = env.GetStrings("CORS_ALLOWED_METHODS", ",", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"})
	corsAllowedHeaders     = env.GetStrings("CORS_ALLOWED_HEADERS", ",", []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-ID", "X-Request-Id", "Origin", "User-Agent", "Accept-Encoding", "Accept-Language", "Cache-Control", "Connection", "DNT", "Host", "Pragma", "Referer"})
	corsAllowedCredentials = env.GetBool("CORS_ALLOWED_CREDENTIALS", true)
	corsMaxAge             = env.GetInt("CORS_MAX_AGE", 300)

	// DB
	dbConnString   = env.MustString("DATABASE_URL")
	dbMaxOpenConns = env.GetInt("DATABASE_MAX_OPEN_CONNS", 20)
	dbMaxIdleConns = env.GetInt("DATABASE_IDLE_CONNS", 2)

	// Redis
	redisConnString = env.GetString("REDIS_URL", "redis://localhost:6379/0")

	// Queue
	workerConcurrency = env.GetInt("WORKER_CONCURRENCY", 10)
	queueName         = env.GetString("QUEUE_NAME", "default")
	queueTaskDeadline = env.GetDuration("QUEUE_TASK_DEADLINE", time.Minute)
	queueMaxRetry     = env.GetInt("QUEUE_TASK_RETRY_LIMIT", 3)
)
