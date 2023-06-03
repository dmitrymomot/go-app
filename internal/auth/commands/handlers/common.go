package command_handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrymomot/go-app/internal/auth/dto"
	auth_repository "github.com/dmitrymomot/go-app/internal/auth/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	// User email verification sender interface.
	userEmailVerificationSender interface {
		SendEmail(ctx context.Context, email string, verificationID uuid.UUID, otp string) error
	}
)

// getVerificationParams struct contains parameters for getVerification function.
type getVerificationParams struct {
	ID   uuid.UUID
	Type dto.VerificationType
	OTP  string
}

// getVerificationByID is a helper function that returns verification by ID.
// It returns error if verification is not found or if verification is expired.
func getVerificationByID(ctx context.Context, repo auth_repository.TxQuerier, arg getVerificationParams) (auth_repository.Verification, error) {
	// Find verification by ID.
	verification, err := repo.FindVerificationByID(ctx, arg.ID)
	if err != nil {
		if auth_repository.IsNotFoundError(err) {
			return auth_repository.Verification{}, ErrVerificationNotFound
		}
		return auth_repository.Verification{}, fmt.Errorf("failed to find verification by ID: %w", err)
	}

	// Check verification type.
	if verification.VerificationType != string(arg.Type) {
		return auth_repository.Verification{}, ErrVerificationInvalidType
	}
	// Check verification expiration.
	if verification.ExpiresAt.Before(time.Now()) {
		return auth_repository.Verification{}, ErrVerificationExpired
	}
	// Check OTP.
	if err := bcrypt.CompareHashAndPassword(verification.OtpHash, []byte(arg.OTP)); err != nil {
		return auth_repository.Verification{}, ErrVerificationInvalidOTP
	}

	return verification, nil
}
