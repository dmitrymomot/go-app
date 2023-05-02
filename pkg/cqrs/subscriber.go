package cqrs

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

// NewSubscriber creates a new subscriber for the given consumer group.
// If consumer group empty, fan-out mode will be used.
func NewSubscriber(redisClient redis.UniversalClient, consumerGroup string, logger Logger) (message.Subscriber, error) {
	if redisClient == nil {
		return nil, fmt.Errorf("redis client is nil")
	}
	if logger == nil {
		logger = watermill.NewStdLogger(false, false)
	}

	sbscr, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client:        redisClient,
			ConsumerGroup: consumerGroup,
			Unmarshaller:  redisstream.DefaultMarshallerUnmarshaller{},
		},
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create subscriber: %w", err)
	}

	return sbscr, nil
}
