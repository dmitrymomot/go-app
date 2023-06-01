package command_handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrymomot/go-app/internal/auth/commands"
	"github.com/dmitrymomot/go-app/internal/auth/dto"
	auth_repository "github.com/dmitrymomot/go-app/internal/auth/repository"
	"github.com/dmitrymomot/go-utils"
	"github.com/dmitrymomot/random"
	"golang.org/x/crypto/bcrypt"
)

// RequestToUpdateUserEmail is a handler for RequestToUpdateUserEmail command.
func RequestToUpdateUserEmail(
	repo auth_repository.TxQuerier,
	sender userEmailVerificationSender,
) func(context.Context, commands.RequestToUpdateUserEmail) error {
	return func(ctx context.Context, arg commands.RequestToUpdateUserEmail) error {
		email, err := utils.SanitizeEmail(arg.Email)
		if err != nil {
			return fmt.Errorf("failed to update email address: %w", err)
		}

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

		if user.Email == email {
			return fmt.Errorf("email is already used")
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
			VerificationType: string(dto.VerificationTypeNewEmail),
			Email:            email,
			OtpHash:          otpHash,
		})
		if err != nil {
			return fmt.Errorf("failed to store or update verification: %w", err)
		}

		// Send auth email.
		if err := sender.SendEmail(ctx, email, verificationID, otp); err != nil {
			return fmt.Errorf("failed to send update email verification: %w", err)
		}

		if err := txRepo.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		return nil
	}
}

// UpdateUserEmail is a handler for UpdateUserEmail command.
func UpdateUserEmail(
	repo auth_repository.TxQuerier,
) func(context.Context, commands.UpdateUserEmail) error {
	return func(ctx context.Context, arg commands.UpdateUserEmail) error {
		txRepo, err := repo.BeginTx(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer txRepo.Rollback() // nolint: errcheck

		// Find verification by id.
		verification, err := txRepo.FindVerificationByID(ctx, arg.VerificationID)
		if err != nil {
			return fmt.Errorf("failed to find verification: %w", err)
		}

		// Check verification type.
		if verification.VerificationType != string(dto.VerificationTypeNewEmail) {
			return fmt.Errorf("invalid verification type")
		}
		if verification.ExpiresAt.Before(time.Now()) {
			return fmt.Errorf("otp is expired")
		}

		// Find user by id.
		user, err := txRepo.FindUserByID(ctx, verification.UserID)
		if err != nil {
			return fmt.Errorf("failed to find user: %w", err)
		}
		if user.Email == verification.Email {
			return fmt.Errorf("email is already used")
		}

		// Check OTP hash.
		if err := bcrypt.CompareHashAndPassword(verification.OtpHash, []byte(arg.OTP)); err != nil {
			return fmt.Errorf("invalid OTP")
		}

		// Update user email.
		if err := txRepo.UpdateUserEmailByID(ctx, auth_repository.UpdateUserEmailByIDParams{
			ID:       user.ID,
			Email:    verification.Email,
			Verified: true,
		}); err != nil {
			return fmt.Errorf("failed to update user email: %w", err)
		}

		// Delete user verification.
		if err := txRepo.DeleteVerificationByID(ctx, verification.ID); err != nil {
			return fmt.Errorf("failed to delete verification: %w", err)
		}

		if err := txRepo.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		return nil
	}
}
