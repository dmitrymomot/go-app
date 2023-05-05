package main

import "github.com/google/uuid"

type (
	// BaseEvent is the base event
	BaseEvent struct {
		UserID uuid.UUID `json:"user_id,omitempty"`
	}

	// UserCreatedEvent is the event that will be stored in the event store
	UserCreatedEvent struct {
		Name      string `json:"name,omitempty"`
		Email     string `json:"email,omitempty"`
		Status    string `json:"status,omitempty"`
		CreatedAt int64  `json:"created_at,omitempty"`
	}

	// UserNameUpdatedEvent is the event that will be stored in the event store
	UserNameUpdatedEvent struct {
		Name      string `json:"name,omitempty"`
		UpdatedAt int64  `json:"updated_at,omitempty"`
	}

	// UserEmailUpdatedEvent is the event that will be stored in the event store
	UserEmailUpdatedEvent struct {
		Email     string `json:"email,omitempty"`
		UpdatedAt int64  `json:"updated_at,omitempty"`
	}

	// UserStatusUpdatedEvent is the event that will be stored in the event store
	UserStatusUpdatedEvent struct {
		Status    string `json:"status,omitempty"`
		UpdatedAt int64  `json:"updated_at,omitempty"`
	}
)

// EventID returns the id of the event
func (be BaseEvent) EventID() uuid.UUID {
	return uuid.Nil
}

// AggregateID returns the id of the aggregate
func (be BaseEvent) AggregateID() uuid.UUID {
	return be.UserID
}

// EventVersion returns the version of the event
func (be BaseEvent) EventVersion() int32 {
	return 0
}

// EventTime returns the time of the event
func (be BaseEvent) EventTime() int64 {
	return 0
}

// EventData returns the data of the event
func (be BaseEvent) EventData() []byte {
	return nil
}
