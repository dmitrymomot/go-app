package queries

import "context"

// Query is an interface for queries.
type Query[Q any, R any] interface {
	// Handle is a function that returns query result or an error.
	Handle(ctx context.Context, arg Q) (R, error)
}

// query is an implementation of Query interface.
type query[Q any, R any] struct {
	fn func(ctx context.Context, arg Q) (R, error)
}

// NewQuery is a constructor of query.
// It returns Query[Q] interface.
func NewQuery[Q any, R any](fn func(ctx context.Context, arg Q) (R, error)) Query[Q, R] {
	return &query[Q, R]{
		fn: fn,
	}
}

// Handle is a function that returns query result or an error.
func (q *query[Q, R]) Handle(ctx context.Context, arg Q) (R, error) {
	return q.fn(ctx, arg)
}
