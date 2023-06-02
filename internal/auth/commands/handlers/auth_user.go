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
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// RequestAuthUser is a handler for RequestAuthUser command.
func RequestAuthUser(repo auth_repository.TxQuerier, sender userEmailVerificationSender) func(ctx context.Context, arg commands.RequestAuthUser) (dto.VerificationID, error) {
	return func(ctx context.Context, arg commands.RequestAuthUser) (dto.VerificationID, error) {
		email, err := utils.SanitizeEmail(arg.Email)
		if err != nil {
			return dto.VerificationID{}, fmt.Errorf("invalid email address: %w", err)
		}

		txRepo, err := repo.BeginTx(ctx)
		if err != nil {
			return dto.VerificationID{}, fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer txRepo.Rollback() // nolint: errcheck

		// Find user by email or create new one.
		user, err := txRepo.FindUserByEmail(ctx, email)
		if err != nil {
			if !auth_repository.IsNotFoundError(err) {
				return dto.VerificationID{}, fmt.Errorf("failed to find user by email: %w", err)
			}
		}
		id := user.ID
		if user.ID == uuid.Nil {
			id, err = txRepo.CreateUser(ctx, email)
			if err != nil {
				return dto.VerificationID{}, fmt.Errorf("failed to create user: %w", err)
			}
		}

		// Generate OTP hash.
		otp := random.String(6, random.Numeric)
		otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
		if err != nil {
			return dto.VerificationID{}, fmt.Errorf("failed to generate OTP hash: %w", err)
		}

		// Store or update user verification.
		verificationID, err := txRepo.StoreOrUpdateVerification(ctx, auth_repository.StoreOrUpdateVerificationParams{
			UserID:           id,
			VerificationType: string(dto.VerificationTypeAuth),
			Email:            email,
			OtpHash:          otpHash,
		})
		if err != nil {
			return dto.VerificationID{}, fmt.Errorf("failed to store verification: %w", err)
		}

		// Send auth email.
		if err := sender.SendEmail(ctx, email, verificationID, otp); err != nil {
			return dto.VerificationID{}, fmt.Errorf("failed to send auth email: %w", err)
		}

		if err := txRepo.Commit(); err != nil {
			return dto.VerificationID{}, fmt.Errorf("failed to commit transaction: %w", err)
		}

		return dto.VerificationID{
			ID: verificationID,
		}, nil
	}
}

// AuthUser is a handler for AuthUser command.
func AuthUser(repo auth_repository.TxQuerier) func(ctx context.Context, arg commands.AuthUser) (dto.Token, error) {
	return func(ctx context.Context, arg commands.AuthUser) (dto.Token, error) {
		txRepo, err := repo.BeginTx(ctx)
		if err != nil {
			return dto.Token{}, fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer txRepo.Rollback() // nolint: errcheck

		// Find user verification by id.
		verification, err := txRepo.FindVerificationByID(ctx, arg.VerificationID)
		if err != nil {
			return dto.Token{}, fmt.Errorf("failed to find verification by id: %w", err)
		}

		// Check verification type.
		if verification.VerificationType != string(dto.VerificationTypeAuth) {
			return dto.Token{}, fmt.Errorf("invalid verification type")
		}
		// Check verification expiration.
		if verification.ExpiresAt.Before(time.Now()) {
			return dto.Token{}, fmt.Errorf("otp is expired")
		}
		// Check OTP.
		if err := bcrypt.CompareHashAndPassword(verification.OtpHash, []byte(arg.OTP)); err != nil {
			return dto.Token{}, fmt.Errorf("invalid OTP")
		}

		// Update user verification.
		if err := txRepo.UpdateUserVerificationStatusByID(ctx, auth_repository.UpdateUserVerificationStatusByIDParams{
			ID:       verification.UserID,
			Verified: true,
		}); err != nil {
			return dto.Token{}, fmt.Errorf("failed to update user verification status: %w", err)
		}

		// Delete user verification.
		if err := txRepo.DeleteVerificationByID(ctx, verification.ID); err != nil {
			return dto.Token{}, fmt.Errorf("failed to delete verification: %w", err)
		}

		// TODO: Create a new token.
		if err := txRepo.StoreToken(ctx, auth_repository.StoreTokenParams{}); err != nil {
			return dto.Token{}, fmt.Errorf("failed to store token: %w", err)
		}

		// Commit transaction.
		if err := txRepo.Commit(); err != nil {
			return dto.Token{}, fmt.Errorf("failed to commit transaction: %w", err)
		}

		return dto.Token{}, nil
	}
}
