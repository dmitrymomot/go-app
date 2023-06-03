package command_handlers

import "errors"

// Predefined errors.
var (
	ErrVerificationNotFound          = errors.New("verification not found")
	ErrVerificationExpired           = errors.New("OTP is expired")
	ErrVerificationInvalidOTP        = errors.New("invalid OTP")
	ErrVerificationInvalidType       = errors.New("invalid type of verification")
	ErrFailedToGenerateOTP           = errors.New("failed to generate OTP")
	ErrFailedToStoreVerification     = errors.New("failed to store verification OTP")
	ErrFailedToSendVerificationEmail = errors.New("failed to send verification email")
)
