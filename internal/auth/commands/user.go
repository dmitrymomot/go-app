package commands

import "github.com/google/uuid"

// RequestAuthUser is a command for requesting user authentication link.
type RequestAuthUser struct {
	Email string `json:"email"`
}

// UpdateUserEmail is a command for updating user email.
type UpdateUserEmail struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}

// UpdateUserVerified is a command for updating user verified status.
type UpdateUserVerified struct {
	UserID   uuid.UUID `json:"user_id"`
	Verified bool      `json:"verified"`
}

// DeleteUser is a command for deleting user.
type DeleteUser struct {
	UserID uuid.UUID `json:"user_id"`
}
