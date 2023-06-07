package commands_test

import (
	"context"
	"errors"
	"testing"

	"github.com/dmitrymomot/go-app/internal/auth/commands"
	"github.com/stretchr/testify/require"
)

func TestCommand_Handle(t *testing.T) {
	// Test Case 1: Ensure Handle function returns no error with correct input.
	t.Run("no error", func(t *testing.T) {
		c := commands.NewCommand(func(ctx context.Context, arg int) error { return nil })

		err := c.Handle(context.Background(), 5)
		require.NoError(t, err)
	})

	// Test Case 2: Ensure Handle function returns error with incorrect input.
	t.Run("error", func(t *testing.T) {
		errorMsg := "invalid argument"
		c := commands.NewCommand(func(ctx context.Context, arg int) error {
			if arg <= 0 {
				return errors.New(errorMsg)
			}
			return nil
		})

		err := c.Handle(context.Background(), -5)
		require.Error(t, err)
		require.EqualError(t, err, errorMsg)
	})
}
