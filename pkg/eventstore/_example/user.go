package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dmitrymomot/go-app/pkg/eventstore"
	"github.com/dmitrymomot/go-utils"
	"github.com/google/uuid"
)

type (
	User struct {
		ID               uuid.UUID `json:"id,omitempty"`
		Name             string    `json:"name,omitempty"`
		Email            string    `json:"email,omitempty"`
		Status           string    `json:"status,omitempty"`
		CreatedAt        int64     `json:"created_at,omitempty"`
		UpdatedAt        int64     `json:"updated_at,omitempty"`
		lastEventVersion int32
	}
)

// NewUser creates a new user
func NewUser() *User {
	return &User{
		ID:        uuid.New(),
		CreatedAt: time.Now().UnixNano(),
	}
}

// OnEvent applies the event to the aggregate
func (u *User) ApplyEvent(event eventstore.Event) error {
	fmt.Println("OnEvent", event.EventType, string(event.EventData))
	switch event.EventType {
	case utils.FullyQualifiedStructName(UserCreatedEvent{}):
		var e UserCreatedEvent
		if err := json.Unmarshal(event.EventData, &e); err != nil {
			return err
		}
		u.ID = event.AggregateID
		u.Name = e.Name
		u.Email = e.Email
		u.Status = e.Status
		u.CreatedAt = e.CreatedAt
	case utils.FullyQualifiedStructName(UserNameUpdatedEvent{}):
		var e UserNameUpdatedEvent
		if err := json.Unmarshal(event.EventData, &e); err != nil {
			return err
		}
		u.Name = e.Name
		u.UpdatedAt = e.UpdatedAt
	case utils.FullyQualifiedStructName(UserEmailUpdatedEvent{}):
		var e UserEmailUpdatedEvent
		if err := json.Unmarshal(event.EventData, &e); err != nil {
			return err
		}
		u.Email = e.Email
		u.UpdatedAt = e.UpdatedAt
	case utils.FullyQualifiedStructName(UserStatusUpdatedEvent{}):
		var e UserStatusUpdatedEvent
		if err := json.Unmarshal(event.EventData, &e); err != nil {
			return err
		}
		u.Status = e.Status
		u.UpdatedAt = e.UpdatedAt
	}

	u.lastEventVersion = event.EventVersion

	return nil
}

// LastEventVersion returns the last applied event version
func (u *User) LastEventVersion() int32 {
	return u.lastEventVersion
}
