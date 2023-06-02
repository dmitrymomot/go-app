package command_handlers

import (
	"context"
	"fmt"
	"time"

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
) func(context.Context, commands.RequestToDeleteUser) (dto.VerificationID, error) {
	return func(ctx context.Context, arg commands.RequestToDeleteUser) (dto.VerificationID, error) {
		txRepo, err := repo.BeginTx(ctx)
		if err != nil {
			return dto.VerificationID{}, fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer txRepo.Rollback() // nolint: errcheck

		// Find user by id.
		user, err := txRepo.FindUserByID(ctx, arg.UserID)
		if err != nil {
			return dto.VerificationID{}, fmt.Errorf("failed to find user: %w", err)
		}

		// Generate OTP hash.
		otp := random.String(6, random.Numeric)
		otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
		if err != nil {
			return dto.VerificationID{}, fmt.Errorf("failed to generate OTP hash: %w", err)
		}

		// Store or update user verification.
		verificationID, err := txRepo.StoreOrUpdateVerification(ctx, auth_repository.StoreOrUpdateVerificationParams{
			UserID:           user.ID,
			VerificationType: string(dto.VerificationTypeDeleteUser),
			Email:            user.Email,
			OtpHash:          otpHash,
		})
		if err != nil {
			return dto.VerificationID{}, fmt.Errorf("failed to store or update verification: %w", err)
		}

		// Send email with OTP.
		if err := sender.SendEmail(ctx, user.Email, verificationID, otp); err != nil {
			return dto.VerificationID{}, fmt.Errorf("failed to send email: %w", err)
		}

		// Commit transaction.
		if err := txRepo.Commit(); err != nil {
			return dto.VerificationID{}, fmt.Errorf("failed to commit transaction: %w", err)
		}

		return dto.VerificationID{
			ID: verificationID,
		}, nil
	}
}

// DeleteUser is a command handler for deleting user.
// It returns an error if the user is not found.
// Otherwise, it returns nil.
func DeleteUser(
	repo auth_repository.TxQuerier,
) func(context.Context, commands.DeleteUser) (dto.UserID, error) {
	return func(ctx context.Context, arg commands.DeleteUser) (dto.UserID, error) {
		txRepo, err := repo.BeginTx(ctx)
		if err != nil {
			return dto.UserID{}, fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer txRepo.Rollback() // nolint: errcheck

		// Find verification by id.
		verification, err := txRepo.FindVerificationByID(ctx, arg.VerificationID)
		if err != nil {
			return dto.UserID{}, fmt.Errorf("failed to find verification: %w", err)
		}

		// Check verification type.
		if verification.VerificationType != string(dto.VerificationTypeDeleteUser) {
			return dto.UserID{}, fmt.Errorf("invalid verification type")
		}
		if verification.ExpiresAt.Before(time.Now()) {
			return dto.UserID{}, fmt.Errorf("otp is expired")
		}

		// Check OTP hash.
		if err := bcrypt.CompareHashAndPassword(verification.OtpHash, []byte(arg.OTP)); err != nil {
			return dto.UserID{}, fmt.Errorf("invalid OTP")
		}

		// Find user by id.
		user, err := txRepo.FindUserByID(ctx, verification.UserID)
		if err != nil {
			return dto.UserID{}, fmt.Errorf("failed to find user: %w", err)
		}
		if user.Email != verification.Email {
			return dto.UserID{}, fmt.Errorf("invalid user")
		}

		// Delete user.
		if err := txRepo.DeleteUserByID(ctx, user.ID); err != nil {
			return dto.UserID{}, fmt.Errorf("failed to delete user: %w", err)
		}

		// Delete user verification.
		if err := txRepo.DeleteVerificationByID(ctx, verification.ID); err != nil {
			return dto.UserID{}, fmt.Errorf("failed to delete verification: %w", err)
		}

		// Commit transaction.
		if err := txRepo.Commit(); err != nil {
			return dto.UserID{}, fmt.Errorf("failed to commit transaction: %w", err)
		}

		return dto.UserID{
			ID: user.ID,
		}, nil
	}
}
