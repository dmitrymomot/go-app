package command_handlers_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/dmitrymomot/go-app/internal/auth/commands"
	command_handlers "github.com/dmitrymomot/go-app/internal/auth/commands/handlers"
	mocks_repository "github.com/dmitrymomot/go-app/internal/auth/mocks/repository"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCleanUpExpiredTokens(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mocks_repository.TxQuerier{}
		repo.On("CleanUpTokens", mock.Anything).Return(nil).Once()

		fn := command_handlers.CleanUpExpiredTokens(repo)
		err := fn(context.Background(), commands.CleanUpExpiredTokens{})
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		errTest := errors.New("test error")
		repo := &mocks_repository.TxQuerier{}
		repo.On("CleanUpTokens", mock.Anything).Return(errTest).Once()

		fn := command_handlers.CleanUpExpiredTokens(repo)
		err := fn(context.Background(), commands.CleanUpExpiredTokens{})
		require.Error(t, err)
		require.ErrorIs(t, err, errTest)
	})

	t.Run("not found error", func(t *testing.T) {
		repo := &mocks_repository.TxQuerier{}
		repo.On("CleanUpTokens", mock.Anything).Return(sql.ErrNoRows).Once()

		fn := command_handlers.CleanUpExpiredTokens(repo)
		err := fn(context.Background(), commands.CleanUpExpiredTokens{})
		require.NoError(t, err)
	})
}
