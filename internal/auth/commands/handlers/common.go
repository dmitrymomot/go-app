package command_handlers

import (
	"context"

	"github.com/google/uuid"
)

type (
	// User email verification sender interface.
	userEmailVerificationSender interface {
		SendEmail(ctx context.Context, email string, verificationID uuid.UUID, otp string) error
	}
)
