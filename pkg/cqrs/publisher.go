package cqrs

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/redis/go-redis/v9"
)

// NewPublisher initializes the events publisher based on the Redis stream client.
func NewPublisher(redisClient redis.UniversalClient, logger Logger) (*redisstream.Publisher, error) {
	if redisClient == nil {
		return nil, fmt.Errorf("redis client is nil")
	}
	if logger == nil {
		logger = watermill.NewStdLogger(false, false)
	}

	// Initialize the events publisher based on the Redis stream client.
	publisher, err := redisstream.NewPublisher(
		redisstream.PublisherConfig{
			Client:     redisClient,
			Marshaller: redisstream.DefaultMarshallerUnmarshaller{},
		},
		logger,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create events publisher: %w", err)
	}

	return publisher, nil
}
