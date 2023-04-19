package middlewares

import (
	"net/http"
	"strconv"
)

// Testing middleware. Helps to test any HTTP error.
// Pass must_err query parameter with code you want get
// E.g.: /login?must_err=403
func Testing() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if errCodeStr := r.URL.Query().Get("must_err"); len(errCodeStr) == 3 {
				if errCode, err := strconv.Atoi(errCodeStr); err == nil && errCode >= 400 && errCode < 600 {
					// Return error with needed status code

					// TODO: Add error message
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
