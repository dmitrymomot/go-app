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

func TestRequestAuthUser_NewUser(t *testing.T) {
	uid := uuid.New()
	vid := uuid.New()
	email := "test@mail.dev"
	repo := &mocks_repository.TxQuerier{}
	repo.On("BeginTx", mock.Anything, mock.Anything).Return(repo, nil).Once()
	repo.On("Commit", mock.Anything).Return(nil).Once()
	repo.On("Rollback", mock.Anything).Return(nil).Once()
	repo.On("FindUserByEmail", mock.Anything, email).Return(auth_repository.User{}, sql.ErrNoRows).Once()
	repo.On("CreateUser", mock.Anything, email).Return(uid, nil).Once()
	repo.On("StoreOrUpdateVerification", mock.Anything, mock.Anything).Return(vid, nil).Once()

	mailSender := &mocks_mail.AuthUserEmailSender{}
	mailSender.On("SendEmail", mock.Anything, email, vid, mock.Anything).Return(nil).Once()

	fn := command_handlers.RequestAuthUser(repo, mailSender)
	err := fn(context.Background(), commands.RequestAuthUser{Email: email})
	require.NoError(t, err)
}

func TestRequestAuthUser_ExistedUser(t *testing.T) {
	uid := uuid.New()
	vid := uuid.New()
	email := "test@mail.dev"
	user := auth_repository.User{
		ID:       uid,
		Email:    email,
		Verified: false,
	}
	repo := &mocks_repository.TxQuerier{}
	repo.On("BeginTx", mock.Anything, mock.Anything).Return(repo, nil).Once()
	repo.On("Commit", mock.Anything).Return(nil).Once()
	repo.On("Rollback", mock.Anything).Return(nil).Once()
	repo.On("FindUserByEmail", mock.Anything, email).Return(user, sql.ErrNoRows).Once()
	repo.On("StoreOrUpdateVerification", mock.Anything, mock.Anything).Return(vid, nil).Once()

	mailSender := &mocks_mail.AuthUserEmailSender{}
	mailSender.On("SendEmail", mock.Anything, email, vid, mock.Anything).Return(nil).Once()

	fn := command_handlers.RequestAuthUser(repo, mailSender)
	err := fn(context.Background(), commands.RequestAuthUser{Email: email})
	require.NoError(t, err)
}

func TestRequestAuthUser_FailedToSendEmail(t *testing.T) {
	uid := uuid.New()
	vid := uuid.New()
	email := "test@mail.dev"
	user := auth_repository.User{
		ID:       uid,
		Email:    email,
		Verified: false,
	}
	repo := &mocks_repository.TxQuerier{}
	repo.On("BeginTx", mock.Anything, mock.Anything).Return(repo, nil).Once()
	repo.On("Commit", mock.Anything).Return(nil).Once()
	repo.On("Rollback", mock.Anything).Return(nil).Once()
	repo.On("FindUserByEmail", mock.Anything, email).Return(user, sql.ErrNoRows).Once()
	repo.On("StoreOrUpdateVerification", mock.Anything, mock.Anything).Return(vid, nil).Once()

	mailSender := &mocks_mail.AuthUserEmailSender{}
	sendEmailErr := fmt.Errorf("something went wrong")
	mailSender.On("SendEmail", mock.Anything, email, vid, mock.Anything).Return(sendEmailErr).Once()

	fn := command_handlers.RequestAuthUser(repo, mailSender)
	err := fn(context.Background(), commands.RequestAuthUser{Email: email})
	require.Error(t, err)
	require.EqualError(t, err, fmt.Sprintf("failed to send auth email: %s", sendEmailErr.Error()))
}
