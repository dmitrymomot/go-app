package command_handlers_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/dmitrymomot/go-app/internal/auth/commands"
	command_handlers "github.com/dmitrymomot/go-app/internal/auth/commands/handlers"
	mocks_mail "github.com/dmitrymomot/go-app/internal/auth/mocks/mail"
	mocks_repository "github.com/dmitrymomot/go-app/internal/auth/mocks/repository"
	auth_repository "github.com/dmitrymomot/go-app/internal/auth/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRequestAuthUser(t *testing.T) {
	uid := uuid.New()
	vid := uuid.New()
	email := "test@mail.dev"

	mailSender := &mocks_mail.UserEmailVerificationSender{}
	mailSender.On("SendEmail", mock.Anything, email, vid, mock.Anything).Return(nil)

	newRepoFn := func() *mocks_repository.TxQuerier {
		repo := &mocks_repository.TxQuerier{}
		repo.On("BeginTx", mock.Anything, mock.Anything).Return(repo, nil).Once()
		repo.On("Commit", mock.Anything).Return(nil).Once()
		repo.On("Rollback", mock.Anything).Return(nil).Once()
		return repo
	}

	t.Run("new user", func(t *testing.T) {
		repo := newRepoFn()
		repo.On("FindUserByEmail", mock.Anything, email).Return(auth_repository.User{}, sql.ErrNoRows).Once()
		repo.On("CreateUser", mock.Anything, email).Return(uid, nil).Once()
		repo.On("StoreOrUpdateVerification", mock.Anything, mock.Anything).Return(vid, nil).Once()

		fn := command_handlers.RequestAuthUser(repo, mailSender)
		err := fn(context.Background(), commands.RequestAuthUser{Email: email})
		require.NoError(t, err)
	})

	t.Run("existing user", func(t *testing.T) {
		repo := newRepoFn()
		repo.On("FindUserByEmail", mock.Anything, email).Return(auth_repository.User{
			ID:       uid,
			Email:    email,
			Verified: false,
		}, nil).Once()
		repo.On("StoreOrUpdateVerification", mock.Anything, mock.Anything).Return(vid, nil).Once()

		fn := command_handlers.RequestAuthUser(repo, mailSender)
		err := fn(context.Background(), commands.RequestAuthUser{Email: email})
		require.NoError(t, err)
	})

	t.Run("failed to create user", func(t *testing.T) {
		repo := newRepoFn()
		repo.On("FindUserByEmail", mock.Anything, email).Return(auth_repository.User{}, sql.ErrNoRows).Once()
		createUserErr := fmt.Errorf("something went wrong")
		repo.On("CreateUser", mock.Anything, email).Return(uuid.Nil, createUserErr).Once()

		fn := command_handlers.RequestAuthUser(repo, mailSender)
		err := fn(context.Background(), commands.RequestAuthUser{Email: email})
		require.Error(t, err)
		require.EqualError(t, err, fmt.Sprintf("failed to create user: %s", createUserErr.Error()))
	})

	t.Run("failed to store verification", func(t *testing.T) {
		repo := newRepoFn()
		repo.On("FindUserByEmail", mock.Anything, email).Return(auth_repository.User{}, sql.ErrNoRows).Once()
		repo.On("CreateUser", mock.Anything, email).Return(uid, nil).Once()
		storeVerificationErr := fmt.Errorf("something went wrong")
		repo.On("StoreOrUpdateVerification", mock.Anything, mock.Anything).Return(uuid.Nil, storeVerificationErr).Once()

		fn := command_handlers.RequestAuthUser(repo, mailSender)
		err := fn(context.Background(), commands.RequestAuthUser{Email: email})
		require.Error(t, err)
		require.ErrorIs(t, err, command_handlers.ErrFailedToStoreVerification)
	})

	t.Run("failed to begin transaction", func(t *testing.T) {
		repo := &mocks_repository.TxQuerier{}
		beginTxErr := fmt.Errorf("something went wrong")
		repo.On("BeginTx", mock.Anything, mock.Anything).Return(nil, beginTxErr).Once()

		fn := command_handlers.RequestAuthUser(repo, mailSender)
		err := fn(context.Background(), commands.RequestAuthUser{Email: email})
		require.Error(t, err)
		require.EqualError(t, err, fmt.Sprintf("failed to begin transaction: %s", beginTxErr.Error()))
	})

	t.Run("failed to send email", func(t *testing.T) {
		repo := newRepoFn()
		repo.On("FindUserByEmail", mock.Anything, email).Return(auth_repository.User{}, sql.ErrNoRows).Once()
		repo.On("CreateUser", mock.Anything, email).Return(uid, nil).Once()
		repo.On("StoreOrUpdateVerification", mock.Anything, mock.Anything).Return(vid, nil).Once()

		mailSender := &mocks_mail.UserEmailVerificationSender{}
		sendEmailErr := fmt.Errorf("something went wrong")
		mailSender.On("SendEmail", mock.Anything, email, vid, mock.Anything).Return(sendEmailErr)

		fn := command_handlers.RequestAuthUser(repo, mailSender)
		err := fn(context.Background(), commands.RequestAuthUser{Email: email})
		require.Error(t, err)
		require.ErrorIs(t, err, command_handlers.ErrFailedToSendVerificationEmail)
	})
}
