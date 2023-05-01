package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dmitrymomot/go-app/pkg/worker"
	"github.com/hibiken/asynq"
)

type (
	// Enqueuer is a helper struct for enqueuing email tasks.
	// It encapsulates the worker.Enqueuer struct and adds queue methods.
	// See pkg/worker/enqueuer.go.
	Enqueuer struct {
		*worker.Enqueuer
	}
)

// NewEnqueuer creates a new email enqueuer.
func NewEnqueuer(e *worker.Enqueuer) *Enqueuer {
	return &Enqueuer{Enqueuer: e}
}

// SendConfirmationEmail enqueues a task to send an example email.
// This function returns an error if the task could not be enqueued.
func (e *Enqueuer) SendConfirmationEmail(ctx context.Context, email, otp string) error {
	payload, err := json.Marshal(SendConfirmationEmailPayload{
		Email: email,
		OTP:   otp,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	return e.EnqueueTask(ctx, asynq.NewTask(SendConfirmationEmailTask, payload))
}
