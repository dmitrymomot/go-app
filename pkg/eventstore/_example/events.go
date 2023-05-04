package main

type (
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
