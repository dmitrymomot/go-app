package main

import (
	"github.com/dmitrymomot/go-app/pkg/httpserver"
	"github.com/dmitrymomot/go-app/pkg/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// init router with default middlewares and routes
func initRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(
		// middleware.Logger,
		middleware.Recoverer,
		middleware.AllowContentType(allowContentTypes...),
		middleware.CleanPath,
		middleware.StripSlashes,
		middleware.GetHead,
		middleware.NoCache,
		middleware.RealIP,
		middleware.RequestID,
		middleware.Timeout(httpRequestTimeout),

		// Basic CORS
		// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
		cors.Handler(cors.Options{
			AllowedOrigins:   corsAllowedOrigins,
			AllowedMethods:   corsAllowedMethods,
			AllowedHeaders:   corsAllowedHeaders,
			AllowCredentials: corsAllowedCredentials,
			MaxAge:           corsMaxAge, // Maximum value not ignored by any of major browsers
		}),

		// Uses for testing error response with needed status code
		middlewares.Testing(),
	)

	// Default error handlers
	r.NotFound(httpserver.NotFoundHandler())
	r.MethodNotAllowed(httpserver.MethodNotAllowedHandler())

	// Default routes
	r.HandleFunc("/health", httpserver.HealthCheckHandler())
	r.Handle("/static/*", httpserver.FileServer("./public", "/static/"))

	if appDebug {
		r.Mount("/debug", middleware.Profiler())
	}

	return r
}
