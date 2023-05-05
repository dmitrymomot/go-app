package eventstore

import (
	"github.com/google/uuid"
)

type Event struct {
	EventID      uuid.UUID `json:"event_id"`
	AggregateID  uuid.UUID `json:"aggregate_id"`
	EventType    string    `json:"event_type"`
	EventVersion int32     `json:"event_version"`
	EventData    []byte    `json:"event_data"`
	EventTime    int64     `json:"event_time"`
}

type Snapshot struct {
	SnapshotID         uuid.UUID `json:"snapshot_id"`
	AggregateID        uuid.UUID `json:"aggregate_id"`
	AggregateType      string    `json:"aggregate_type"`
	SnapshotVersion    int32     `json:"snapshot_version"`
	SnapshotData       []byte    `json:"snapshot_data"`
	SnapshotTime       int64     `json:"snapshot_time"`
	LatestEventVersion int32     `json:"latest_event_version"`
}
