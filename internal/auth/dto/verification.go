package dto

// VerificationType is a type of verification.
type VerificationType string

// Predefined verification types.
const (
	VerificationTypeAuth       VerificationType = "auth"        // verification to issue auth token
	VerificationTypeNewEmail   VerificationType = "new_email"   // verification to update email
	VerificationTypeDeleteUser VerificationType = "delete_user" // verification to delete user account
)
