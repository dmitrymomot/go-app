package queries

import "errors"

// Predefined errors.
var (
	ErrUserNotFound  = errors.New("user not found")
	ErrTokenNotFound = errors.New("token not found")
	ErrInvalidToken  = errors.New("invalid or expired token")
)
