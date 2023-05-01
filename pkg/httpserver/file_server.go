package httpserver

import (
	"fmt"
	"net/http"
	"strings"
)

// File server handler
func FileServer(dir, prefix string) http.HandlerFunc {
	prefix = fmt.Sprintf("/%s/", strings.Trim(prefix, "/"))
	fs := http.StripPrefix(prefix, http.FileServer(http.Dir(dir)))
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, prefix) {
			fs.ServeHTTP(w, r)
			return
		}
		http.NotFound(w, r)
	}
}

// File server middleware
func FileServerMiddleware(dir, prefix string) func(http.Handler) http.Handler {
	prefix = fmt.Sprintf("/%s/", strings.Trim(prefix, "/"))
	fs := http.StripPrefix(prefix, http.FileServer(http.Dir(dir)))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, prefix) {
				fs.ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
