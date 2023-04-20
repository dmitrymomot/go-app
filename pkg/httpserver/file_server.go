package httpserver

import (
	"fmt"
	"net/http"
	"strings"
)

// File server handler
func FileServer(dir, prefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prefix = fmt.Sprintf("/%s/", strings.Trim(prefix, "/"))
		if strings.HasPrefix(r.URL.Path, prefix) {
			http.StripPrefix(prefix, http.FileServer(http.Dir(dir))).ServeHTTP(w, r)
			return
		}
		http.NotFound(w, r)
	}
}
