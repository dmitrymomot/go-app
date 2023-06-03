package command_handlers

import "errors"

// Predefined errors.
var (
	ErrVerificationNotFound    = errors.New("verification not found")
	ErrVerificationExpired     = errors.New("OTP is expired")
	ErrVerificationInvalidOTP  = errors.New("invalid OTP")
	ErrVerificationInvalidType = errors.New("invalid type of verification")
)
