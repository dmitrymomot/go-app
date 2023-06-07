package auth

import (
	"context"

	"github.com/dmitrymomot/go-app/internal/auth/dto"
	"github.com/google/uuid"
)

type (
	// TokenService is an interface that describes the token service.
	TokenService interface {
		// Token method takes user id, audience and returns access & refresh tokens.
		// It also store session in the database.
		Token(ctx context.Context, uid uuid.UUID, audience ...string) (dto.Token, error)
		// Refresh method takes refresh token and returns new access & refresh tokens.
		Refresh(ctx context.Context, token string) (dto.Token, error)
		// Logout method takes access token and removes session from the database.
		Logout(ctx context.Context, token string) error
		// LogoutAll method takes access token and removes all user sessions from the database.
		LogoutAll(ctx context.Context, token string) error
		// LogoutOther method takes access token and removes all user sessions
		// except the current one from the database.
		LogoutOther(ctx context.Context, token string) error
		// CleanUp method removes all expired sessions from the database.
		CleanUp(ctx context.Context) error
	}
)
