package httpserver

import "errors"

// Predefined errors
var (
	ErrNotFound         = errors.New("not_found")
	ErrMethodNotAllowed = errors.New("method_not_allowed")
)
