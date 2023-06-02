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
		c := commands.NewCommand(func(ctx context.Context, arg int) (string, error) { return "test", nil })

		resp, err := c.Handle(context.Background(), 5)
		require.NoError(t, err)
		require.Equal(t, "test", resp)
	})

	// Test Case 2: Ensure Handle function returns error with incorrect input.
	t.Run("error", func(t *testing.T) {
		errorMsg := "invalid argument"
		c := commands.NewCommand(func(ctx context.Context, arg int) (string, error) {
			if arg <= 0 {
				return "", errors.New(errorMsg)
			}
			return "test", nil
		})

		resp, err := c.Handle(context.Background(), -5)
		require.Error(t, err)
		require.EqualError(t, err, errorMsg)
		require.Equal(t, "", resp)
	})
}
