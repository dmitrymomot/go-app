package command_handlers

import (
	"context"
	"fmt"

	"github.com/dmitrymomot/go-app/internal/auth/commands"
	auth_repository "github.com/dmitrymomot/go-app/internal/auth/repository"
)

// CleanUpVerifications command handler.
func CleanUpVerifications(repo auth_repository.TxQuerier) func(ctx context.Context, cmd commands.CleanUpVerifications) error {
	return func(ctx context.Context, cmd commands.CleanUpVerifications) error {
		if err := repo.CleanUpVerifications(ctx); err != nil {
			if !auth_repository.IsNotFoundError(err) {
				return fmt.Errorf("failed to clean up verifications: %w", err)
			}
		}
		return nil
	}
}
