package cqrs

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

// NewEventHandlerSubscriber creates a new subscriber for the given handler name.
// It's used by CQRS to create a subscriber for each event handler.
func NewEventHandlerSubscriber(redisClient redis.UniversalClient, handlerName string, logger Logger) (message.Subscriber, error) {
	if redisClient == nil {
		return nil, fmt.Errorf("redis client is nil")
	}
	if logger == nil {
		logger = watermill.NewStdLogger(false, false)
	}

	sbscr, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client:       redisClient,
			Consumer:     handlerName,
			Unmarshaller: redisstream.DefaultMarshallerUnmarshaller{},
		},
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create subscriber for handler %s: %w", handlerName, err)
	}

	return sbscr, nil
}

// NewSubscriber creates a new subscriber.
func NewSubscriber(redisClient redis.UniversalClient, logger Logger) (message.Subscriber, error) {
	if redisClient == nil {
		return nil, fmt.Errorf("redis client is nil")
	}
	if logger == nil {
		logger = watermill.NewStdLogger(false, false)
	}

	sbscr, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client:       redisClient,
			Unmarshaller: redisstream.DefaultMarshallerUnmarshaller{},
		},
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create subscriber: %w", err)
	}

	return sbscr, nil
}
