package command_handlers_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dmitrymomot/go-app/internal/auth/commands"
	command_handlers "github.com/dmitrymomot/go-app/internal/auth/commands/handlers"
	"github.com/dmitrymomot/go-app/internal/auth/dto"
	mocks_mail "github.com/dmitrymomot/go-app/internal/auth/mocks/mail"
	mocks_repository "github.com/dmitrymomot/go-app/internal/auth/mocks/repository"
	auth_repository "github.com/dmitrymomot/go-app/internal/auth/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestRequestToUpdateUserEmail(t *testing.T) {
	uid := uuid.New()
	vid := uuid.New()
	email := "test@mail.dev"
	newEmail := "test2@mail.dev"

	mailSender := &mocks_mail.UserEmailVerificationSender{}
	mailSender.On("SendEmail", mock.Anything, newEmail, vid, mock.Anything).Return(nil).Once()

	newRepoFn := func() *mocks_repository.TxQuerier {
		repo := &mocks_repository.TxQuerier{}
		repo.On("BeginTx", mock.Anything, mock.Anything).Return(repo, nil).Once()
		repo.On("Commit", mock.Anything).Return(nil).Once()
		repo.On("Rollback", mock.Anything).Return(nil).Once()
		repo.On("StoreOrUpdateVerification", mock.Anything, mock.Anything).Return(vid, nil).Once()
		return repo
	}

	t.Run("user not found", func(t *testing.T) {
		repo := newRepoFn()
		repo.On("FindUserByID", mock.Anything, uid).Return(auth_repository.User{}, sql.ErrNoRows).Once()

		fn := command_handlers.RequestToUpdateUserEmail(repo, mailSender)
		_, err := fn(context.Background(), commands.RequestToUpdateUserEmail{
			UserID: uid,
			Email:  newEmail,
		})
		require.Error(t, err)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("email is already used", func(t *testing.T) {
		repo := newRepoFn()
		repo.On("FindUserByID", mock.Anything, uid).Return(auth_repository.User{
			ID:    uid,
			Email: email,
		}, nil).Once()

		fn := command_handlers.RequestToUpdateUserEmail(repo, mailSender)
		_, err := fn(context.Background(), commands.RequestToUpdateUserEmail{
			UserID: uid,
			Email:  email,
		})
		require.Error(t, err)
		require.EqualError(t, err, "email is already used")
	})

	t.Run("email is not valid", func(t *testing.T) {
		repo := newRepoFn()
		repo.On("FindUserByID", mock.Anything, uid).Return(auth_repository.User{
			ID:    uid,
			Email: email,
		}, nil).Once()

		fn := command_handlers.RequestToUpdateUserEmail(repo, mailSender)
		_, err := fn(context.Background(), commands.RequestToUpdateUserEmail{
			UserID: uid,
			Email:  "test",
		})
		require.Error(t, err)
		require.EqualError(t, err, "failed to update email address: invalid email address")
	})

	t.Run("success", func(t *testing.T) {
		repo := newRepoFn()
		repo.On("FindUserByID", mock.Anything, uid).Return(auth_repository.User{
			ID:    uid,
			Email: email,
		}, nil).Once()

		fn := command_handlers.RequestToUpdateUserEmail(repo, mailSender)
		resp, err := fn(context.Background(), commands.RequestToUpdateUserEmail{
			UserID: uid,
			Email:  newEmail,
		})
		require.NoError(t, err)
		require.Equal(t, vid, resp.ID)
	})
}

func TestUpdateUserEmail(t *testing.T) {
	uid := uuid.New()
	vid := uuid.New()
	email := "test@mail.dev"
	newEmail := "test2@mail.dev"

	otp := "123456"
	otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	require.NoError(t, err)

	newRepoFn := func() *mocks_repository.TxQuerier {
		repo := &mocks_repository.TxQuerier{}

		repo.On("BeginTx", mock.Anything, mock.Anything).Return(repo, nil).Once()
		repo.On("Commit", mock.Anything).Return(nil).Once()
		repo.On("Rollback", mock.Anything).Return(nil).Once()
		repo.On("DeleteVerificationByID", mock.Anything, vid).Return(nil).Once()

		repo.On("UpdateUserEmailByID", mock.Anything, auth_repository.UpdateUserEmailByIDParams{
			ID:       uid,
			Email:    newEmail,
			Verified: true,
		}).Return(nil).Once()

		return repo
	}

	t.Run("success", func(t *testing.T) {
		repo := newRepoFn()

		repo.On("FindUserByID", mock.Anything, uid).Return(auth_repository.User{
			ID:    uid,
			Email: email,
		}, nil).Once()
		repo.On("FindVerificationByID", mock.Anything, vid).Return(auth_repository.Verification{
			ID:               vid,
			UserID:           uid,
			VerificationType: string(dto.VerificationTypeNewEmail),
			Email:            newEmail,
			OtpHash:          otpHash,
			ExpiresAt:        time.Now().Add(time.Minute),
		}, nil).Once()

		fn := command_handlers.UpdateUserEmail(repo)
		resp, err := fn(context.Background(), commands.UpdateUserEmail{
			VerificationID: vid,
			OTP:            otp,
		})
		require.NoError(t, err)
		require.Equal(t, uid, resp.ID)
	})

	t.Run("user not found", func(t *testing.T) {
		repo := newRepoFn()
		repo.On("FindVerificationByID", mock.Anything, vid).Return(auth_repository.Verification{
			ID:               vid,
			UserID:           uid,
			VerificationType: string(dto.VerificationTypeNewEmail),
			Email:            newEmail,
			OtpHash:          otpHash,
			ExpiresAt:        time.Now().Add(time.Minute),
		}, nil).Once()
		repo.On("FindUserByID", mock.Anything, uid).Return(auth_repository.User{}, sql.ErrNoRows).Once()

		fn := command_handlers.UpdateUserEmail(repo)
		_, err := fn(context.Background(), commands.UpdateUserEmail{
			VerificationID: vid,
			OTP:            otp,
		})
		require.Error(t, err)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("otp is expired", func(t *testing.T) {
		repo := newRepoFn()
		repo.On("FindVerificationByID", mock.Anything, vid).Return(auth_repository.Verification{
			ID:               vid,
			UserID:           uid,
			VerificationType: string(dto.VerificationTypeNewEmail),
			Email:            newEmail,
			ExpiresAt:        time.Now().Add(-time.Minute),
		}, nil).Once()

		fn := command_handlers.UpdateUserEmail(repo)
		_, err := fn(context.Background(), commands.UpdateUserEmail{
			VerificationID: vid,
			OTP:            otp,
		})
		require.Error(t, err)
		require.EqualError(t, err, "otp is expired")
	})
}
