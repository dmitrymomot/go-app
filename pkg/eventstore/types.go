package eventstore

import (
	"encoding/json"

	"github.com/google/uuid"
)

// Aggregator is the interface for an event aggregator
type Aggregator interface {
	// Init loads the state of the aggregate from the snapshot
	Init(snapshot Snapshot) error
	// OnEvent applies the event to the aggregate
	OnEvent(event Event) error
	// AggregateID returns the id of the aggregate
	AggregateID() uuid.UUID
	// AggregateType returns the type of the aggregate
	AggregateType() string
	// AggregateVersion returns the version of the aggregate
	AggregateVersion() int32
	// AggregateState returns the state of the aggregate
	AggregateState() json.RawMessage
	// LatestEventVersion returns the version of the latest event
	LatestEventVersion() int32
	// LatestEventTime returns the time of the latest event
	LatestEventTime() int64
}
