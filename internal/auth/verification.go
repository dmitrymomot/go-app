package auth

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
	// Verifier is an interface that describes the verification service.
	Verifier interface {
		// Auth method takes user email and sends a verification code to it.
		Auth(ctx context.Context, uid uuid.UUID, email string) error
		// UpdateEmail method takes user id and new email and sends a verification code to it.
		UpdateEmail(ctx context.Context, uid uuid.UUID, email string) error
		// DeleteProfile method takes user id and sends a verification code to it.
		DeleteProfile(ctx context.Context, uid uuid.UUID, email string) error
		// Verify method takes verification id and code and checks it.
		Verify(ctx context.Context, id uuid.UUID, code string, verificationType dto.VerificationType) error
		// CleanUp method removes all expired verification codes from the database.
		CleanUp(ctx context.Context) error
	}

	// verifier is an implementation of the Verifier interface.
	verifier struct {
		repo       auth_repository.TxQuerier
		sender     userEmailVerificationSender
		defaultTTL time.Duration
	}

	// User email verification sender interface.
	userEmailVerificationSender interface {
		Auth(ctx context.Context, email string, verificationID uuid.UUID, otp string) error
		UpdateEmail(ctx context.Context, email string, verificationID uuid.UUID, otp string) error
		DeleteProfile(ctx context.Context, email string, verificationID uuid.UUID, otp string) error
	}
)

// NewVerifier returns a new instance of the Verifier interface.
func NewVerifier(repo auth_repository.TxQuerier, sender userEmailVerificationSender, defaultTTL time.Duration) Verifier {
	if defaultTTL == 0 {
		defaultTTL = time.Minute * 15
	}
	return &verifier{
		repo:       repo,
		sender:     sender,
		defaultTTL: defaultTTL,
	}
}

// Auth method takes user email and sends a verification code to it.
func (v *verifier) Auth(ctx context.Context, uid uuid.UUID, email string) error {
	// Begin DB transaction.
	txRepo, err := v.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer txRepo.Rollback() // nolint: errcheck

	// Store and send verification.
	if err := v.storeAndSendVerification(ctx, txRepo, dto.VerificationTypeAuth, uid, email, v.defaultTTL); err != nil {
		return fmt.Errorf("failed to store and send verification: %w", err)
	}

	// Commit transaction.
	if err := txRepo.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateEmail method takes user id and new email and sends a verification code to it.
func (v *verifier) UpdateEmail(ctx context.Context, uid uuid.UUID, email string) error {
	// Begin DB transaction.
	txRepo, err := v.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer txRepo.Rollback() // nolint: errcheck

	// Store and send verification.
	if err := v.storeAndSendVerification(ctx, txRepo, dto.VerificationTypeNewEmail, uid, email, v.defaultTTL); err != nil {
		return fmt.Errorf("failed to store and send verification: %w", err)
	}

	// Commit transaction.
	if err := txRepo.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DeleteProfile method takes user id and sends a verification code to it.
func (v *verifier) DeleteProfile(ctx context.Context, uid uuid.UUID, email string) error {
	// Begin DB transaction.
	txRepo, err := v.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer txRepo.Rollback() // nolint: errcheck

	// Store and send verification.
	if err := v.storeAndSendVerification(
		ctx, txRepo,
		dto.VerificationTypeDeleteUser,
		uid, email,
		v.defaultTTL,
	); err != nil {
		return fmt.Errorf("failed to store and send verification: %w", err)
	}

	// Commit transaction.
	if err := txRepo.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Verify method takes verification id and code and checks it.
func (v *verifier) Verify(ctx context.Context, id uuid.UUID, code string, verificationType dto.VerificationType) error {
	// Begin DB transaction.
	txRepo, err := v.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer txRepo.Rollback() // nolint: errcheck

	// Get verification.
	verification, err := txRepo.FindVerificationByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get verification: %w", err)
	}

	// Check verification.
	if verification.VerificationType == string(verificationType) {
		return fmt.Errorf("verification type mismatch")
	}
	if verification.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("verification expired")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(verification.OtpHash), []byte(code)); err != nil {
		return fmt.Errorf("verification code mismatch")
	}

	// Delete verification from the database.
	if err := txRepo.DeleteVerificationByID(ctx, id); err != nil {
		return fmt.Errorf("failed to delete verification: %w", err)
	}

	// Commit transaction.
	if err := txRepo.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// CleanUp method removes all expired verification codes from the database.
func (v *verifier) CleanUp(ctx context.Context) error {
	if err := v.repo.CleanUpVerifications(ctx); err != nil {
		if !auth_repository.IsNotFoundError(err) {
			return fmt.Errorf("failed to clean up verifications: %w", err)
		}
	}
	return nil
}

// storeAndSendVerification method stores verification and sends it to the user.
func (v *verifier) storeAndSendVerification(ctx context.Context, repo auth_repository.TxQuerier, verificationType dto.VerificationType, uid uuid.UUID, email string, ttl time.Duration) error {
	// Generate OTP.
	otp := random.String(6, random.Numeric)
	// Hash OTP.
	otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash OTP: %w", err)
	}

	// Store verification.
	verificationID, err := repo.StoreOrUpdateVerification(ctx, auth_repository.StoreOrUpdateVerificationParams{
		UserID:           uid,
		VerificationType: string(verificationType),
		Email:            email,
		OtpHash:          otpHash,
	})
	if err != nil {
		return fmt.Errorf("failed to save verification: %w", err)
	}

	switch verificationType {
	case dto.VerificationTypeAuth:
		// Send verification email.
		if err := v.sender.Auth(ctx, email, verificationID, otp); err != nil {
			return fmt.Errorf("failed to send verification email: %w", err)
		}
	case dto.VerificationTypeNewEmail:
		// Send verification email.
		if err := v.sender.UpdateEmail(ctx, email, verificationID, otp); err != nil {
			return fmt.Errorf("failed to send verification email: %w", err)
		}
	case dto.VerificationTypeDeleteUser:
		// Send verification email.
		if err := v.sender.DeleteProfile(ctx, email, verificationID, otp); err != nil {
			return fmt.Errorf("failed to send verification email: %w", err)
		}
	default:
		return fmt.Errorf("unknown verification type: %s", verificationType)
	}

	return nil
}
