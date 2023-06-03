package command_handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrymomot/go-app/internal/auth/dto"
	auth_repository "github.com/dmitrymomot/go-app/internal/auth/repository"
	"github.com/dmitrymomot/random"
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

// generateVerificationParams struct contains parameters for generateVerification function.
type generateVerificationParams struct {
	UserID uuid.UUID
	Email  string
	Type   dto.VerificationType
}

// generateAndSendVerification is a helper function that generates OTP for verification and stores it in the database.
// It also sends an email with OTP to the user.
// It returns otp or an error if OTP generation or storage fails.
func generateAndSendVerification(ctx context.Context, repo auth_repository.TxQuerier, sender userEmailVerificationSender, arg generateVerificationParams) error {
	// Generate OTP hash.
	otp := random.String(6, random.Numeric)
	otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return ErrFailedToGenerateOTP
	}

	// Store or update user verification.
	verificationID, err := repo.StoreOrUpdateVerification(ctx, auth_repository.StoreOrUpdateVerificationParams{
		UserID:           arg.UserID,
		VerificationType: string(arg.Type),
		Email:            arg.Email,
		OtpHash:          otpHash,
	})
	if err != nil {
		return ErrFailedToStoreVerification
	}

	// Send auth email.
	if err := sender.SendEmail(ctx, arg.Email, verificationID, otp); err != nil {
		return ErrFailedToSendVerificationEmail
	}

	return nil
}
