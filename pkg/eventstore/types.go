package eventstore

type EventType interface {
	// EventType returns the type of the event
	EventType() string
}

type EventVersion interface {
	// EventVersion returns the version of the event
	EventVersion() int32
}

type EventTime interface {
	// EventTime returns the time of the event
	EventTime() int64
}

type EventData interface {
	// EventData returns the data of the event
	EventData() []byte
}

type AggregateType interface {
	// AggregateType returns the type of the aggregate
	AggregateType() string
}

type LoadAggregateFromSnapshot interface {
	// LoadAggregateFromSnapshot loads the aggregate from the snapshot
	LoadSnapshot(snapshot Snapshot) error
}

type ApplyEventToAggregate interface {
	// ApplyEvent applies the event to the aggregate
	ApplyEvent(event Event) error
}

type GetSnapshotData interface {
	// GetSnapshot returns the snapshot of the aggregate
	GetSnapshotData() ([]byte, error)
}

type LastEventVersion interface {
	// LastEventVersion returns the version of the latest event
	LastEventVersion() int32
}
