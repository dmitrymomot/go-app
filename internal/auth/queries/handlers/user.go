package query_handlers

import (
	"context"
	"fmt"

	"github.com/dmitrymomot/go-app/internal/auth/dto"
	"github.com/dmitrymomot/go-app/internal/auth/queries"
	auth_repository "github.com/dmitrymomot/go-app/internal/auth/repository"
)

// GetUserByID returns user by id query handler.
func GetUserByID(repo auth_repository.TxQuerier) func(ctx context.Context, query queries.GetUserByID) (dto.User, error) {
	return func(ctx context.Context, query queries.GetUserByID) (dto.User, error) {
		user, err := repo.FindUserByID(ctx, query.UserID)
		if err != nil {
			if auth_repository.IsNotFoundError(err) {
				return dto.User{}, queries.ErrUserNotFound
			}
			return dto.User{}, fmt.Errorf("failed to get user by id: %w", err)
		}

		return castUser(user), nil
	}
}

// GetUserByEmail returns user by email query handler.
func GetUserByEmail(repo auth_repository.TxQuerier) func(ctx context.Context, query queries.GetUserByEmail) (dto.User, error) {
	return func(ctx context.Context, query queries.GetUserByEmail) (dto.User, error) {
		user, err := repo.FindUserByEmail(ctx, query.Email)
		if err != nil {
			if auth_repository.IsNotFoundError(err) {
				return dto.User{}, queries.ErrUserNotFound
			}
			return dto.User{}, fmt.Errorf("failed to get user by email: %w", err)
		}

		return castUser(user), nil
	}
}
