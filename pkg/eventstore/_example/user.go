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
		ID          uuid.UUID `json:"id,omitempty"`
		Name        string    `json:"name,omitempty"`
		Email       string    `json:"email,omitempty"`
		Status      string    `json:"status,omitempty"`
		CreatedAt   int64     `json:"created_at,omitempty"`
		UpdatedAt   int64     `json:"updated_at,omitempty"`
		snapshot    eventstore.Snapshot
		latestEvent *eventstore.Event
	}
)

// NewUser creates a new user
func NewUser() *User {
	return &User{
		ID:        uuid.New(),
		CreatedAt: time.Now().UnixNano(),
	}
}

// Init loads the state of the aggregate from the snapshot
func (u *User) Init(snapshot eventstore.Snapshot) error {
	u.snapshot = snapshot
	if snapshot.SnapshotData != nil {
		if err := json.Unmarshal(snapshot.SnapshotData, u); err != nil {
			return err
		}
	}
	return nil
}

// OnEvent applies the event to the aggregate
func (u *User) OnEvent(event eventstore.Event) error {
	fmt.Println("OnEvent", event.EventType, string(event.EventData))
	switch event.EventType {
	case "UserCreatedEvent":
		var e UserCreatedEvent
		if err := json.Unmarshal(event.EventData, &e); err != nil {
			return err
		}
		utils.PrettyPrint("UserCreatedEvent", e)
		u.ID = event.AggregateID
		u.Name = e.Name
		u.Email = e.Email
		u.Status = e.Status
		u.CreatedAt = e.CreatedAt
	case "UserNameUpdatedEvent":
		var e UserNameUpdatedEvent
		if err := json.Unmarshal(event.EventData, &e); err != nil {
			return err
		}
		u.Name = e.Name
		u.UpdatedAt = e.UpdatedAt
	case "UserEmailUpdatedEvent":
		var e UserEmailUpdatedEvent
		if err := json.Unmarshal(event.EventData, &e); err != nil {
			return err
		}
		u.Email = e.Email
		u.UpdatedAt = e.UpdatedAt
	case "UserStatusUpdatedEvent":
		var e UserStatusUpdatedEvent
		if err := json.Unmarshal(event.EventData, &e); err != nil {
			return err
		}
		u.Status = e.Status
		u.UpdatedAt = e.UpdatedAt
	}
	u.latestEvent = &event
	return nil
}

// AggregateID returns the id of the aggregate
func (u *User) AggregateID() uuid.UUID {
	return u.ID
}

// AggregateType returns the type of the aggregate
func (u *User) AggregateType() string {
	return "UserModel"
}

// AggregateVersion returns the version of the aggregate
func (u *User) AggregateVersion() int32 {
	if u.latestEvent != nil {
		return u.snapshot.SnapshotVersion + 1
	}
	return u.snapshot.SnapshotVersion
}

// AggregateState returns the state of the aggregate
func (u *User) AggregateState() ([]byte, error) {
	data, err := json.Marshal(u)
	if err != nil {
		return nil, fmt.Errorf("could not marshal user: %w", err)
	}
	return data, nil
}

// LatestEventVersion returns the version of the latest event
func (u *User) LatestEventVersion() int32 {
	if u.latestEvent != nil {
		return u.latestEvent.EventVersion
	}
	return u.snapshot.LatestEventVersion
}

// LatestEventTime returns the time of the latest event
func (u *User) LatestEventTime() int64 {
	if u.latestEvent != nil {
		return u.latestEvent.EventTime
	}
	return u.snapshot.SnapshotTime
}
