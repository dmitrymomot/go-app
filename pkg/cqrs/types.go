package cqrs

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
)

// Logger is an interface, that you need to implement to support Watermill logging.
// You can use watermill.StdLoggerAdapter as a reference implementation.
type Logger interface {
	Error(msg string, err error, fields watermill.LogFields)
	Info(msg string, fields watermill.LogFields)
	Debug(msg string, fields watermill.LogFields)
	Trace(msg string, fields watermill.LogFields)
	With(fields watermill.LogFields) watermill.LoggerAdapter
}

type (
	// CommanfHandlerFactory is a function that creates a new CommandHandler.
	CommanfHandlerFactory func(cb CommandBus, eb EventBus) CommandHandler

	// EventHandlerFactory is a function that creates a new EventHandler.
	EventHandlerFactory func(cb CommandBus, eb EventBus) EventHandler
)

// CommandBus interface
type CommandBus interface {
	// Send sends command to the command bus.
	Send(ctx context.Context, cmd interface{}) error
}

// EventBus interface
type EventBus interface {
	// Publish sends event to the event bus.
	Publish(ctx context.Context, event interface{}) error
}

// CommandHandler receives a command defined by NewCommand and handles it with the Handle method.
// If using DDD, CommandHandler may modify and persist the aggregate.
//
// In contrast to EventHandler, every Command must have only one CommandHandler.
//
// One instance of CommandHandler is used during handling messages.
// When multiple commands are delivered at the same time, Handle method can be executed multiple times at the same time.
// Because of that, Handle method needs to be thread safe!
type CommandHandler interface {
	// HandlerName is the name used in message.Router while creating handler.
	//
	// It will be also passed to CommandsSubscriberConstructor.
	// May be useful, for example, to create a consumer group per each handler.
	//
	// WARNING: If HandlerName was changed and is used for generating consumer groups,
	// it may result with **reconsuming all messages**!
	HandlerName() string

	NewCommand() interface{}

	Handle(ctx context.Context, cmd interface{}) error
}

// EventHandler receives events defined by NewEvent and handles them with its Handle method.
// If using DDD, CommandHandler may modify and persist the aggregate.
// It can also invoke a process manager, a saga or just build a read model.
//
// In contrast to CommandHandler, every Event can have multiple EventHandlers.
//
// One instance of EventHandler is used during handling messages.
// When multiple events are delivered at the same time, Handle method can be executed multiple times at the same time.
// Because of that, Handle method needs to be thread safe!
type EventHandler interface {
	// HandlerName is the name used in message.Router while creating handler.
	//
	// It will be also passed to EventsSubscriberConstructor.
	// May be useful, for example, to create a consumer group per each handler.
	//
	// WARNING: If HandlerName was changed and is used for generating consumer groups,
	// it may result with **reconsuming all messages** !!!
	HandlerName() string

	NewEvent() interface{}

	Handle(ctx context.Context, event interface{}) error
}
