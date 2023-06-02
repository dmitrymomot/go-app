package dto

import "github.com/google/uuid"

// VerificationType is a type of verification.
type VerificationType string

// Predefined verification types.
const (
	VerificationTypeAuth       VerificationType = "auth"        // verification to issue auth token
	VerificationTypeNewEmail   VerificationType = "new_email"   // verification to update email
	VerificationTypeDeleteUser VerificationType = "delete_user" // verification to delete user account
)

// VerificationID is a verification id type.
type VerificationID struct {
	ID uuid.UUID `json:"id"`
}
