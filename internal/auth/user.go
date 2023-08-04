package auth

import (
	"context"

	"github.com/dmitrymomot/go-app/internal/auth/dto"
	"github.com/google/uuid"
)

type (
	// UserService is an interface that describes the user service.
	UserService interface {
		// Create method creates a new user.
		Create(ctx context.Context, email string) (dto.User, error)
		// Get method returns a user by id.
		Get(ctx context.Context, uid uuid.UUID) (dto.User, error)
		// GetByEmail method returns a user by email.
		GetByEmail(ctx context.Context, email string) (dto.User, error)
		// UpdateEmail method updates user email.
		UpdateEmail(ctx context.Context, uid uuid.UUID, newEmail string) error
		// Delete method deletes a user.
		Delete(ctx context.Context, uid uuid.UUID) error
	}

	// UserRepository is an interface that describes the user repository.
	UserRepository interface{}

	// userService is an implementation of the UserService interface.
	userService struct {
		repo UserRepository
	}
)

// NewUserService returns a new instance of the UserService interface.
func NewUserService(repo UserRepository, verifier Verifier) UserService {}
