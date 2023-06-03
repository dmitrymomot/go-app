package queries_test

import (
	"context"
	"errors"
	"testing"

	"github.com/dmitrymomot/go-app/internal/auth/queries"
	"github.com/stretchr/testify/require"
)

func TestQuery_Handler(t *testing.T) {
	// Test Case 1: Ensure Handle function returns no error with correct input.
	t.Run("no error", func(t *testing.T) {
		q := queries.NewQuery(func(ctx context.Context, arg int) (int, error) { return arg, nil })

		result, err := q.Handle(context.Background(), 5)
		require.NoError(t, err)
		require.Equal(t, 5, result)
	})

	// Test Case 2: Ensure Handle function returns error with incorrect input.
	t.Run("error", func(t *testing.T) {
		errorMsg := "invalid argument"
		q := queries.NewQuery(func(ctx context.Context, arg int) (int, error) {
			if arg <= 0 {
				return 0, errors.New(errorMsg)
			}
			return arg, nil
		})

		result, err := q.Handle(context.Background(), -5)
		require.Error(t, err)
		require.EqualError(t, err, errorMsg)
		require.Equal(t, 0, result)
	})
}
