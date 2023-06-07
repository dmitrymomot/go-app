package commands

import "github.com/google/uuid"

// RequestAuthUser is a command for requesting user authentication link.
type RequestAuthUser struct {
	Email string `json:"email"`
}

// AuthUser is a command for authenticating user.
type AuthUser struct {
	VerificationID uuid.UUID `json:"verification_id"`
	OTP            string    `json:"otp"`
}
