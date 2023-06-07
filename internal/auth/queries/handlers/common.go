package query_handlers

import (
	"github.com/dmitrymomot/go-app/internal/auth/dto"
	auth_repository "github.com/dmitrymomot/go-app/internal/auth/repository"
)

// cast auth_repository.User to dto.User
func castUser(user auth_repository.User) dto.User {
	result := dto.User{
		ID:        user.ID,
		Email:     user.Email,
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt.Unix(),
	}
	if user.UpdatedAt.Valid {
		result.UpdatedAt = user.UpdatedAt.Time.Unix()
	}

	return result
}

// generate access & refresh tokens by auth_repository.Token
func castToken(token auth_repository.Token) dto.Token {
	return dto.Token{}
}
