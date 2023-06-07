package query_handlers_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	mocks_repository "github.com/dmitrymomot/go-app/internal/auth/mocks/repository"
	"github.com/dmitrymomot/go-app/internal/auth/queries"
	query_handlers "github.com/dmitrymomot/go-app/internal/auth/queries/handlers"
	auth_repository "github.com/dmitrymomot/go-app/internal/auth/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetUserByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		uid := uuid.New()
		repo := &mocks_repository.TxQuerier{}
		repo.On("FindUserByID", mock.Anything, uid).Return(auth_repository.User{
			ID:        uid,
			UpdatedAt: sql.NullTime{Valid: true, Time: time.Now()},
		}, nil).Once()

		fn := query_handlers.GetUserByID(repo)
		resp, err := fn(context.Background(), queries.GetUserByID{UserID: uid})
		require.NoError(t, err)
		require.Equal(t, uid, resp.ID)
	})

	t.Run("error", func(t *testing.T) {
		uid := uuid.New()
		errTest := errors.New("test error")
		repo := &mocks_repository.TxQuerier{}
		repo.On("FindUserByID", mock.Anything, uid).Return(auth_repository.User{}, errTest).Once()

		fn := query_handlers.GetUserByID(repo)
		_, err := fn(context.Background(), queries.GetUserByID{UserID: uid})
		require.Error(t, err)
		require.ErrorIs(t, err, errTest)
	})

	t.Run("not found error", func(t *testing.T) {
		uid := uuid.New()
		repo := &mocks_repository.TxQuerier{}
		repo.On("FindUserByID", mock.Anything, uid).Return(auth_repository.User{}, sql.ErrNoRows).Once()

		fn := query_handlers.GetUserByID(repo)
		_, err := fn(context.Background(), queries.GetUserByID{UserID: uid})
		require.Error(t, err)
		require.ErrorIs(t, err, queries.ErrUserNotFound)
	})
}

func TestGetUserByEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		email := "test@mail.dev"
		repo := &mocks_repository.TxQuerier{}
		repo.On("FindUserByEmail", mock.Anything, email).Return(auth_repository.User{Email: email}, nil).Once()

		fn := query_handlers.GetUserByEmail(repo)
		resp, err := fn(context.Background(), queries.GetUserByEmail{Email: email})
		require.NoError(t, err)
		require.Equal(t, email, resp.Email)
	})

	t.Run("error", func(t *testing.T) {
		email := "test@mail.dev"
		errTest := errors.New("test error")
		repo := &mocks_repository.TxQuerier{}
		repo.On("FindUserByEmail", mock.Anything, email).Return(auth_repository.User{}, errTest).Once()

		fn := query_handlers.GetUserByEmail(repo)
		_, err := fn(context.Background(), queries.GetUserByEmail{Email: email})
		require.Error(t, err)
		require.ErrorIs(t, err, errTest)
	})

	t.Run("not found error", func(t *testing.T) {
		email := "test@mail.dev"
		repo := &mocks_repository.TxQuerier{}
		repo.On("FindUserByEmail", mock.Anything, email).Return(auth_repository.User{}, sql.ErrNoRows).Once()

		fn := query_handlers.GetUserByEmail(repo)
		_, err := fn(context.Background(), queries.GetUserByEmail{Email: email})
		require.Error(t, err)
		require.ErrorIs(t, err, queries.ErrUserNotFound)
	})
}
