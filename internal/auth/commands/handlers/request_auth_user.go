package command_handlers

import (
	"context"
	"fmt"

	"github.com/dmitrymomot/go-app/internal/auth/commands"
	"github.com/dmitrymomot/go-app/internal/auth/dto"
	auth_repository "github.com/dmitrymomot/go-app/internal/auth/repository"
	"github.com/dmitrymomot/random"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	// authUserSender is an interface for user sender.
	authUserSender interface {
		SendEmail(ctx context.Context, email string, verificationID uuid.UUID, otp string) error
	}
)

// RequestAuthUser is a handler for RequestAuthUser command.
func RequestAuthUser(repo auth_repository.TxQuerier, sender authUserSender) func(ctx context.Context, arg commands.RequestAuthUser) error {
	return func(ctx context.Context, arg commands.RequestAuthUser) error {
		txRepo, err := repo.BeginTx(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer txRepo.Rollback() // nolint: errcheck

		// Find user by email or create new one.
		user, err := txRepo.FindUserByEmail(ctx, arg.Email)
		if err != nil {
			if !auth_repository.IsNotFoundError(err) {
				return fmt.Errorf("failed to find user by email: %w", err)
			}
		}
		id := user.ID
		if user.ID == uuid.Nil {
			id, err = txRepo.CreateUser(ctx, arg.Email)
			if err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}
		}

		// Generate OTP hash.
		otp := random.String(6, random.Numeric)
		otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to generate OTP hash: %w", err)
		}

		// Store or update user verification.
		verificationID, err := txRepo.StoreOrUpdateVerification(ctx, auth_repository.StoreOrUpdateVerificationParams{
			UserID:           id,
			VerificationType: string(dto.VerificationTypeEmail),
			Email:            arg.Email,
			OtpHash:          otpHash,
		})
		if err != nil {
			return fmt.Errorf("failed to store or update verification: %w", err)
		}

		// Send auth email.
		if err := sender.SendEmail(ctx, arg.Email, verificationID, otp); err != nil {
			return fmt.Errorf("failed to send auth email: %w", err)
		}

		if err := txRepo.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		return nil
	}
}
