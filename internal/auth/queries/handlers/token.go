package query_handlers

import (
	"context"

	"github.com/dmitrymomot/go-app/internal/auth/dto"
	"github.com/dmitrymomot/go-app/internal/auth/queries"
	auth_repository "github.com/dmitrymomot/go-app/internal/auth/repository"
)

// GetTokenByID query handler returns access & refresh tokens by verification ID.
func GetTokenByID(repo auth_repository.TxQuerier) func(ctx context.Context, arg queries.GetTokenByID) (dto.Token, error) {
	return func(ctx context.Context, arg queries.GetTokenByID) (dto.Token, error) {
		token, err := repo.FindTokenByID(ctx, arg.ID)
		if err != nil {
			if auth_repository.IsNotFoundError(err) {
				return dto.Token{}, queries.ErrInvalidToken
			}
			return dto.Token{}, err
		}

		return castToken(token), nil
	}
}

// GetTokenByRefreshTokenID query handler returns access & refresh tokens by refresh token.
func GetTokenByRefreshTokenID(repo auth_repository.TxQuerier) func(ctx context.Context, arg queries.GetTokenByRefreshTokenID) (dto.Token, error) {
	return func(ctx context.Context, arg queries.GetTokenByRefreshTokenID) (dto.Token, error) {
		token, err := repo.FindTokenByRefreshTokenID(ctx, arg.ID)
		if err != nil {
			if auth_repository.IsNotFoundError(err) {
				return dto.Token{}, queries.ErrInvalidToken
			}
			return dto.Token{}, err
		}

		return castToken(token), nil
	}
}
