package queries

import "github.com/google/uuid"

// GetTokenByID returns access & refresh tokens by verification ID.
type GetTokenByID struct {
	// ID is a verification ID.
	ID uuid.UUID
}

// GetTokenByRefreshTokenID returns access & refresh tokens by refresh token.
type GetTokenByRefreshTokenID struct {
	// ID is a refresh token ID.
	ID uuid.UUID
}
