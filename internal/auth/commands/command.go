package commands

import (
	"context"
)

// Command is an interface for commands.
type Command[Cmd any] interface {
	// Handle is a function that returns command result or an error.
	Handle(ctx context.Context, arg Cmd) error
}

// command is an implementation of Command interface.
type command[Cmd any] struct {
	fn func(ctx context.Context, arg Cmd) error
}

// NewCommand is a constructor of command.
// It returns Command[Cmd] interface.
func NewCommand[Cmd any](fn func(ctx context.Context, arg Cmd) error) Command[Cmd] {
	return &command[Cmd]{
		fn: fn,
	}
}

// Handle is a function that returns command result or an error.
func (c *command[Cmd]) Handle(ctx context.Context, arg Cmd) error {
	return c.fn(ctx, arg)
}
