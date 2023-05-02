package cqrs

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

// NewFacade creates a new CQRS Facade.
// Read more about cqrs component: https://watermill.io/docs/cqrs/
func NewFacade(
	redisClient redis.UniversalClient, logger Logger, router *message.Router,
	commandsPublisher, eventsPublisher message.Publisher,
	commandsSubscriber message.Subscriber,
	commandHandlers []CommanfHandlerFactory, eventHandlers []EventHandlerFactory,
) (*cqrs.Facade, error) {
	if logger == nil {
		logger = watermill.NewStdLogger(false, false)
	}

	// We are using JSON marshaler, but you can use any marshaler you want.
	cqrsMarshaler := cqrs.JSONMarshaler{}

	// cqrs.Facade is facade for Command and Event buses and processors.
	// You can use facade, or create buses and processors manually (you can inspire with cqrs.NewFacade)
	cqrsFacade, err := cqrs.NewFacade(cqrs.FacadeConfig{
		GenerateCommandsTopic: func(commandName string) string {
			// we are using queue RabbitMQ config, so we need to have topic per command type
			return commandName
		},
		CommandHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.CommandHandler {
			ch := make([]cqrs.CommandHandler, len(commandHandlers))
			for i, factory := range commandHandlers {
				ch[i] = factory(cb, eb)
			}
			return ch
		},
		CommandsPublisher: commandsPublisher,
		CommandsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			// we can reuse subscriber, because all commands have separated topics
			return commandsSubscriber, nil
		},
		GenerateEventsTopic: func(eventName string) string {
			return eventName
		},
		EventHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.EventHandler {
			eh := make([]cqrs.EventHandler, len(eventHandlers))
			for i, factory := range eventHandlers {
				eh[i] = factory(cb, eb)
			}
			return eh
		},
		EventsPublisher: eventsPublisher,
		EventsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			return NewSubscriber(redisClient, handlerName, logger)
		},
		Router:                router,
		CommandEventMarshaler: cqrsMarshaler,
		Logger:                logger,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create cqrs facade: %w", err)
	}

	cqrsFacade.CommandBus()

	return cqrsFacade, nil
}
