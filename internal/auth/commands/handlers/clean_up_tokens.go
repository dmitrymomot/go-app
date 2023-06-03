package command_handlers

import (
	"context"
	"fmt"

	"github.com/dmitrymomot/go-app/internal/auth/commands"
	auth_repository "github.com/dmitrymomot/go-app/internal/auth/repository"
)

// CleanUpExpiredTokens command handler.
func CleanUpExpiredTokens(repo auth_repository.TxQuerier) func(ctx context.Context, cmd commands.CleanUpExpiredTokens) error {
	return func(ctx context.Context, cmd commands.CleanUpExpiredTokens) error {
		if err := repo.CleanUpTokens(ctx); err != nil {
			if !auth_repository.IsNotFoundError(err) {
				return fmt.Errorf("failed to clean up expired tokens: %w", err)
			}
		}
		return nil
	}
}
