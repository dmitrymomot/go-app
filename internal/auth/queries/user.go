package queries

import "github.com/google/uuid"

// GetUserByID query struct.
type GetUserByID struct {
	UserID uuid.UUID `json:"user_id"`
}

// GetUserByEmail returns user by email query handler.
type GetUserByEmail struct {
	Email string `json:"email"`
}
