package command_handlers

import (
	"context"
	"fmt"

	"github.com/dmitrymomot/go-app/internal/auth/commands"
	"github.com/dmitrymomot/go-app/internal/auth/dto"
	auth_repository "github.com/dmitrymomot/go-app/internal/auth/repository"
	"github.com/dmitrymomot/random"
	"golang.org/x/crypto/bcrypt"
)

// RequestToDeleteUser is a command handler for deleting user.
// It returns an error if the user is not found.
// Otherwise, it returns nil.
func RequestToDeleteUser(
	repo auth_repository.TxQuerier,
	sender userEmailVerificationSender,
) func(context.Context, commands.RequestToDeleteUser) error {
	return func(ctx context.Context, arg commands.RequestToDeleteUser) error {
		txRepo, err := repo.BeginTx(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer txRepo.Rollback() // nolint: errcheck

		// Find user by id.
		user, err := txRepo.FindUserByID(ctx, arg.UserID)
		if err != nil {
			return fmt.Errorf("failed to find user: %w", err)
		}

		// Generate OTP hash.
		otp := random.String(6, random.Numeric)
		otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to generate OTP hash: %w", err)
		}

		// Store or update user verification.
		verificationID, err := txRepo.StoreOrUpdateVerification(ctx, auth_repository.StoreOrUpdateVerificationParams{
			UserID:           user.ID,
			VerificationType: string(dto.VerificationTypeDeleteUser),
			Email:            user.Email,
			OtpHash:          otpHash,
		})
		if err != nil {
			return fmt.Errorf("failed to store or update verification: %w", err)
		}

		// Send email with OTP.
		if err := sender.SendEmail(ctx, user.Email, verificationID, otp); err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}

		// Commit transaction.
		if err := txRepo.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		return nil
	}
}

// DeleteUser is a command handler for deleting user.
// It returns an error if the user is not found.
// Otherwise, it returns nil.
func DeleteUser(
	repo auth_repository.TxQuerier,
) func(context.Context, commands.DeleteUser) error {
	return func(ctx context.Context, arg commands.DeleteUser) error {
		txRepo, err := repo.BeginTx(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer txRepo.Rollback() // nolint: errcheck

		// Find user verification by id.
		verification, err := getVerificationByID(ctx, txRepo, getVerificationParams{
			ID:   arg.VerificationID,
			Type: dto.VerificationTypeDeleteUser,
			OTP:  arg.OTP,
		})
		if err != nil {
			return fmt.Errorf("verification failed: %w", err)
		}

		// Find user by id.
		user, err := txRepo.FindUserByID(ctx, verification.UserID)
		if err != nil {
			return fmt.Errorf("failed to find user: %w", err)
		}
		if user.Email != verification.Email {
			return fmt.Errorf("invalid user")
		}

		// Delete user.
		if err := txRepo.DeleteUserByID(ctx, user.ID); err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}

		// Delete user verification.
		if err := txRepo.DeleteVerificationByID(ctx, verification.ID); err != nil {
			return fmt.Errorf("failed to delete verification: %w", err)
		}

		// Commit transaction.
		if err := txRepo.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		return nil
	}
}
