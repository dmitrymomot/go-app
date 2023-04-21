package httpserver

import (
	"fmt"
	"net/http"

	"github.com/dmitrymomot/go-app/pkg/response"
)

// NotFoundHandler is a handler for 404 Not Found
func NotFoundHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := response.NewError(
			http.StatusNotFound,
			ErrNotFound,
			fmt.Sprintf("The requested URL %s was not found on this server.", r.URL.Path),
			nil,
		)
		response.JSON(w, resp)
		return
	}
}

// MethodNotAllowedHandler is a handler for 405 Method Not Allowed
func MethodNotAllowedHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := response.NewError(
			http.StatusMethodNotAllowed,
			ErrMethodNotAllowed,
			fmt.Sprintf("The requested method %s is not allowed for the URL %s.", r.Method, r.URL.Path),
			nil,
		)
		response.JSON(w, resp)
		return
	}
}
