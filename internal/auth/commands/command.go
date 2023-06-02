package commands

import (
	"context"
)

// Command is an interface for commands.
type Command[Cmd any, Resp any] interface {
	// Handle is a function that returns command result or an error.
	Handle(ctx context.Context, arg Cmd) (Resp, error)
}

// command is an implementation of Command interface.
type command[Cmd any, Resp any] struct {
	fn func(ctx context.Context, arg Cmd) (Resp, error)
}

// NewCommand is a constructor of command.
// It returns Command[Cmd] interface.
func NewCommand[Cmd any, Resp any](fn func(ctx context.Context, arg Cmd) (Resp, error)) Command[Cmd, Resp] {
	return &command[Cmd, Resp]{
		fn: fn,
	}
}

// Handle is a function that returns command result or an error.
func (c *command[Cmd, Resp]) Handle(ctx context.Context, arg Cmd) (Resp, error) {
	return c.fn(ctx, arg)
}
