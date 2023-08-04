package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrymomot/go-app/internal/auth/dto"
	auth_repository "github.com/dmitrymomot/go-app/internal/auth/repository"
	"github.com/dmitrymomot/go-utils"
	"github.com/dmitrymomot/random"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	// Service is an interface that describes the authentication service.
	Service interface {
		// Auth method takes user email and sends a verification code to it.
		Auth(ctx context.Context, email string) error
		// Verify method takes verification id and code and returns access & refresh tokens.
		Token(ctx context.Context, id uuid.UUID, code string) (dto.Token, error)
		// Refresh method takes refresh token and returns new access & refresh tokens.
		Refresh(ctx context.Context, token string) (dto.Token, error)
		// Logout method takes access token and removes session from the database.
		Logout(ctx context.Context, token string) error
		// LogoutAll method takes access token and removes all user sessions from the database.
		LogoutAll(ctx context.Context, token string) error
		// LogoutOther method takes access token and removes all user sessions except the current one from the database.
		LogoutOther(ctx context.Context, token string) error

		// UpdateEmail method takes user id and new email and sends a verification code to it.
		UpdateEmail(ctx context.Context, uid uuid.UUID, email string) error
		// VerifyUpdateEmail method takes verification id and code and updates user email.
		VerifyUpdateEmail(ctx context.Context, id uuid.UUID, code string) error

		// DeleteProfile method takes user id and sends a verification code to it.
		DeleteProfile(ctx context.Context, uid uuid.UUID) error
		// VerifyDeleteProfile method takes verification id and code and deletes user profile.
		VerifyDeleteProfile(ctx context.Context, id uuid.UUID, code string) error

		// CleanUp method removes all expired sessions & verification codes from the database.
		CleanUp(ctx context.Context) error
	}

	// service is an implementation of the Service interface.
	service struct {
		repo     auth_repository.TxQuerier
		verifier Verifier
		token    TokenService
	}
)

// NewService is a function that creates a new instance of the Service interface.
func NewService(repo auth_repository.TxQuerier, verifier Verifier, token TokenService) Service {
	return &service{
		repo:     repo,
		verifier: verifier,
		token:    token,
	}
}

// Auth method takes user email and sends a verification code to it.
func (s *service) Auth(ctx context.Context, email string) error {
	// Sanitize email (trim spaces, convert to lowercase, etc).
	email, err := utils.SanitizeEmail(email)
	if err != nil {
		return fmt.Errorf("failed to request auth: %w", err)
	}

	// Begin DB transaction.
	txRepo, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer txRepo.Rollback() // nolint: errcheck

	// Find user by email or create new one.
	var userID uuid.UUID
	{
		user, err := txRepo.FindUserByEmail(ctx, email)
		if err != nil {
			if !auth_repository.IsNotFoundError(err) {
				return fmt.Errorf("failed to find user by email: %w", err)
			}

			userID, err = txRepo.CreateUser(ctx, email)
			if err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			userID = user.ID
		}
	}

	// Generate OTP.
	otp := random.String(6, random.Numeric)
	// Hash OTP.
	otpHash, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash OTP: %w", err)
	}

	// Store verification.
	verificationID, err := txRepo.StoreOrUpdateVerification(ctx, auth_repository.StoreOrUpdateVerificationParams{
		UserID:           userID,
		VerificationType: string(dto.VerificationTypeAuth),
		Email:            email,
		OtpHash:          otpHash,
	})
	if err != nil {
		return fmt.Errorf("failed to save verification: %w", err)
	}

	// Send verification email.
	if err := s.sender.Auth(ctx, email, verificationID, otp); err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	// Commit transaction.
	if err := txRepo.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Verify method takes verification id and code and returns access & refresh tokens.
func (s *service) Token(ctx context.Context, id uuid.UUID, code string) (dto.Token, error) {
	// Begin DB transaction.
	txRepo, err := s.repo.BeginTx(ctx)
	if err != nil {
		return dto.Token{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer txRepo.Rollback() // nolint: errcheck

	// Find verification by id.
	verification, err := txRepo.FindVerificationByID(ctx, id)
	if err != nil {
		if auth_repository.IsNotFoundError(err) {
			return dto.Token{}, fmt.Errorf("verification not found: %w", err)
		}

		return dto.Token{}, fmt.Errorf("failed to find verification by id: %w", err)
	}

	// Check verification type.
	if verification.VerificationType != string(dto.VerificationTypeAuth) {
		return dto.Token{}, fmt.Errorf("invalid verification type: %w", err)
	}

	// Check verification expiration.
	if verification.ExpiresAt.Before(time.Now()) {
		return dto.Token{}, fmt.Errorf("verification expired")
	}

	// Compare OTP.
	if err := bcrypt.CompareHashAndPassword(verification.OtpHash, []byte(code)); err != nil {
		return dto.Token{}, fmt.Errorf("invalid OTP: %w", err)
	}

	// Generate access & refresh tokens.
	var (
		accessTokenID  = uuid.New()
		refreshTokenID = uuid.New()
		issuer         = "auth"
		subject        = verification.UserID
		audience       = "auth"
	)

	if err := txRepo.StoreToken(ctx, auth_repository.StoreTokenParams{
		ID:             verification.ID,
		UserID:         verification.UserID,
		AccessTokenID:  accessTokenID,
		RefreshTokenID: refreshTokenID,
	}); err != nil {
		return dto.Token{}, fmt.Errorf("failed to create tokens: %w", err)
	}

	// Commit transaction.
	if err := txRepo.Commit(); err != nil {
		return dto.Token{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return tokens, nil
}
