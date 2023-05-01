package httpserver

import "net/http"

// HealthCheckHandler is a simple health check handler
func HealthCheckHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}
