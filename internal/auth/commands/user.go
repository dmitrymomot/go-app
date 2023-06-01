package commands

import "github.com/google/uuid"

// RequestToUpdateUserEmail is a command for requesting user email update.
type RequestToUpdateUserEmail struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}

// UpdateUserEmail is a command for updating user email.
type UpdateUserEmail struct {
	VerificationID uuid.UUID `json:"verification_id"`
	OTP            string    `json:"otp"`
}

// RequestToDeleteUser is a command for requesting user deletion.
type RequestToDeleteUser struct {
	UserID uuid.UUID `json:"user_id"`
}

// DeleteUser is a command for deleting user.
type DeleteUser struct {
	VerificationID uuid.UUID `json:"verification_id"`
	OTP            string    `json:"otp"`
}
