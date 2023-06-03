package command_handlers

import (
	"context"
	"fmt"

	"github.com/dmitrymomot/go-app/internal/auth/commands"
	"github.com/dmitrymomot/go-app/internal/auth/dto"
	auth_repository "github.com/dmitrymomot/go-app/internal/auth/repository"
	"github.com/dmitrymomot/go-utils"
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

		// Send update email verification email.
		if err := generateAndSendVerification(ctx, txRepo, sender, generateVerificationParams{
			UserID: user.ID,
			Type:   dto.VerificationTypeNewEmail,
			Email:  email,
		}); err != nil {
			return fmt.Errorf("failed to generate and send verification: %w", err)
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

		// Find user verification by id.
		verification, err := getVerificationByID(ctx, txRepo, getVerificationParams{
			ID:   arg.VerificationID,
			Type: dto.VerificationTypeNewEmail,
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
		if user.Email == verification.Email {
			return fmt.Errorf("email is already used")
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
