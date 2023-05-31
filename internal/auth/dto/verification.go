package dto

// Verification is a data transfer object for verification.
// It is used for internal communication between packages.
type Verification struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Code   string `json:"code"`
}

// VerificationType is a type of verification.
type VerificationType string

const (
	// VerificationTypeEmail is a verification type for email.
	VerificationTypeEmail VerificationType = "email"
)
