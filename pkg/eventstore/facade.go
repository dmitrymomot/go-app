package eventstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/dmitrymomot/go-utils"
	"github.com/google/uuid"
)

// EventStore is the facade for the event store
type EventStore struct {
	db          *sql.DB
	repo        *queries
	eventStream string
}

// NewEventStore creates a new event store
func NewEventStore(ctx context.Context, db *sql.DB, eventStreamName string) (*EventStore, error) {
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	if eventStreamName == "" {
		return nil, fmt.Errorf("event stream name cannot be empty")
	}

	es := &EventStore{
		db:          db,
		repo:        newQueries(db, eventStreamName),
		eventStream: eventStreamName,
	}

	if err := es.prepare(ctx); err != nil {
		return nil, fmt.Errorf("failed to prepare event store: %w", err)
	}

	return es, nil
}

// Prepare prepares the event store.
// It creates the events and snapshot tables for the given event streams.
func (es *EventStore) prepare(ctx context.Context) error {
	tx, err := es.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // nolint:errcheck

	repo := es.repo.WithTx(tx)
	if err := repo.CreateEventsTable(ctx); err != nil {
		return fmt.Errorf("failed to create events table: %w", err)
	}
	if err := repo.CreateSnapshotTable(ctx); err != nil {
		return fmt.Errorf("failed to create snapshot table: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// EventStream returns the event stream name
func (es *EventStore) EventStream() string {
	return es.eventStream
}

// AppendEvents appends events to the event store
func (es *EventStore) AppendEvent(ctx context.Context, aggregateID uuid.UUID, event interface{}) (Event, error) {
	var (
		eventType string
		eventTime int64
		eventData []byte
	)

	// if event implements the EventType interface, we can get the event type
	// otherwise we use the type name of the event
	if e, ok := event.(EventType); ok {
		eventType = e.EventType()
	} else {
		eventType = utils.FullyQualifiedStructName(event)
	}

	// if event implements the EventTime interface, we can get the event time
	// otherwise we use the current time
	if e, ok := event.(EventTime); ok {
		eventTime = e.EventTime()
	} else {
		eventTime = time.Now().UnixNano()
	}

	// if event implements the EventData interface, we can get the event data
	// otherwise we use the event as json
	if e, ok := event.(EventData); ok {
		eventData = e.EventData()
	} else {
		var err error
		eventData, err = json.Marshal(event)
		if err != nil {
			return Event{}, fmt.Errorf("failed to marshal event: %w", err)
		}
	}

	tx, err := es.db.Begin()
	if err != nil {
		return Event{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // nolint:errcheck

	repo := es.repo.WithTx(tx)

	e, err := repo.StoreEvent(ctx, StoreEventParams{
		AggregateID: aggregateID,
		EventType:   eventType,
		EventTime:   eventTime,
		EventData:   eventData,
	})
	if err != nil {
		return Event{}, fmt.Errorf("failed to append event: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return Event{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return e, nil
}

// LoadEvents loads all events for the given aggregate id
func (es *EventStore) LoadEvents(ctx context.Context, aggregateID uuid.UUID) ([]Event, error) {
	events, err := es.repo.LoadAllEvents(ctx, aggregateID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Event{}, nil
		}
		return []Event{}, fmt.Errorf("failed to load events: %w", err)
	}

	return events, nil
}

// LoadEventsRange loads events for the given aggregate id and event version range
func (es *EventStore) LoadEventsRange(ctx context.Context, aggregateID uuid.UUID, fromEventVersion, toEventVersion int32) ([]Event, error) {
	events, err := es.repo.LoadEventsRange(ctx, LoadEventsRangeParams{
		AggregateID:      aggregateID,
		FromEventVersion: fromEventVersion,
		ToEventVersion:   toEventVersion,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Event{}, nil
		}
		return []Event{}, fmt.Errorf("failed to load events range: %w", err)
	}

	return events, nil
}

// LoadNewestEvents loads the newest events for the given aggregate id.
// Helpful when you have a snapshot and want to load all events since the snapshot.
func (es *EventStore) LoadNewestEvents(ctx context.Context, aggregateID uuid.UUID, snapshotVersion int32) ([]Event, error) {
	events, err := es.repo.LoadNewestEvents(ctx, LoadNewestEventsParams{
		AggregateID:        aggregateID,
		LatestEventVersion: snapshotVersion,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []Event{}, nil
		}
		return []Event{}, fmt.Errorf("failed to load newest events: %w", err)
	}

	return events, nil
}

// LoadLatestSnapshot loads the latest snapshot for the given aggregate id and aggregate type
func (es *EventStore) LoadLatestSnapshot(ctx context.Context, aggregateID uuid.UUID, aggregate interface{}) (Snapshot, error) {
	if reflect.ValueOf(aggregate).Kind() != reflect.Ptr {
		return Snapshot{}, ErrAggregateMustBePointer
	}

	var aggregateType string
	if a, ok := aggregate.(AggregateType); ok {
		aggregateType = a.AggregateType()
	} else {
		aggregateType = utils.FullyQualifiedStructName(aggregate)
	}

	snapshot, err := es.repo.LoadLatestSnapshot(ctx, LoadLatestSnapshotParams{
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snapshot{}, ErrNoSnapshotFound
		}
		return Snapshot{}, fmt.Errorf("failed to load latest snapshot: %w", err)
	}

	if snapshot.SnapshotID == uuid.Nil {
		return Snapshot{}, ErrNoSnapshotFound
	}

	if a, ok := aggregate.(LoadAggregateFromSnapshot); ok {
		if err := a.LoadSnapshot(snapshot); err != nil {
			return Snapshot{}, fmt.Errorf("failed to load aggregate from snapshot: %w", err)
		}
	} else {
		if err := json.Unmarshal(snapshot.SnapshotData, aggregate); err != nil {
			return Snapshot{}, fmt.Errorf("failed to unmarshal snapshot data: %w", err)
		}
	}

	return snapshot, nil
}

// LoadSnapshot loads the snapshot for the given aggregate id and snapshot version
func (es *EventStore) LoadSnapshot(ctx context.Context, aggregateID uuid.UUID, aggregate interface{}, snapshotVersion int32) (Snapshot, error) {
	if reflect.ValueOf(aggregate).Kind() != reflect.Ptr {
		return Snapshot{}, ErrAggregateMustBePointer
	}

	var aggregateType string
	if a, ok := aggregate.(AggregateType); ok {
		aggregateType = a.AggregateType()
	} else {
		aggregateType = utils.FullyQualifiedStructName(aggregate)
	}

	snapshot, err := es.repo.LoadSnapshot(ctx, LoadSnapshotParams{
		AggregateID:     aggregateID,
		AggregateType:   aggregateType,
		SnapshotVersion: snapshotVersion,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snapshot{}, ErrNoSnapshotFound
		}
		return Snapshot{}, fmt.Errorf("failed to load latest snapshot: %w", err)
	}

	if snapshot.SnapshotID == uuid.Nil {
		return Snapshot{}, ErrNoSnapshotFound
	}

	if snapshot.SnapshotData != nil {
		if a, ok := aggregate.(LoadAggregateFromSnapshot); ok {
			if err := a.LoadSnapshot(snapshot); err != nil {
				return Snapshot{}, fmt.Errorf("failed to load aggregate from snapshot: %w", err)
			}
		} else {
			if err := json.Unmarshal(snapshot.SnapshotData, aggregate); err != nil {
				return Snapshot{}, fmt.Errorf("failed to unmarshal snapshot data: %w", err)
			}
		}
	}

	return snapshot, nil
}

// LoadCurrentState loads the current state for the given aggregate id.
// It loads the latest snapshot and all events since the snapshot version.
// Then it applies all events to the snapshot and returns the state.
func (es *EventStore) LoadCurrentState(ctx context.Context, aggregateID uuid.UUID, aggregate interface{}) error {
	if reflect.ValueOf(aggregate).Kind() != reflect.Ptr {
		return ErrAggregateMustBePointer
	}

	var aggregateType string
	if a, ok := aggregate.(AggregateType); ok {
		aggregateType = a.AggregateType()
	} else {
		aggregateType = utils.FullyQualifiedStructName(aggregate)
	}

	// load latest snapshot for the given aggregate id
	snapshot, err := es.repo.LoadLatestSnapshot(ctx, LoadLatestSnapshotParams{
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to load latest snapshot: %w", err)
		}
	}

	if snapshot.SnapshotID != uuid.Nil && snapshot.SnapshotData != nil {
		// init aggregate with snapshot
		if a, ok := aggregate.(LoadAggregateFromSnapshot); ok {
			if err := a.LoadSnapshot(snapshot); err != nil {
				return fmt.Errorf("failed to load aggregate from snapshot: %w", err)
			}
		} else {
			if err := json.Unmarshal(snapshot.SnapshotData, aggregate); err != nil {
				return fmt.Errorf("failed to unmarshal snapshot data: %w", err)
			}
		}
	}

	// if there is no snapshot, we can just load all events,
	// otherwise we load all events since the snapshot version
	events, err := es.repo.LoadNewestEvents(ctx, LoadNewestEventsParams{
		AggregateID:        aggregateID,
		LatestEventVersion: snapshot.LatestEventVersion,
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to load newest events: %w", err)
		}
	}

	// nothing to do if there are no events
	if len(events) == 0 {
		return nil
	}

	// apply events to the aggregate
	if a, ok := aggregate.(ApplyEventToAggregate); ok {
		for _, event := range events {
			if err := a.ApplyEvent(event); err != nil {
				return fmt.Errorf("failed to apply event %s: %w", event.EventType, err)
			}
		}
	} else {
		return ErrApplyEventToAggregate
	}

	return nil
}

// StoreSnapshot stores a snapshot for the given aggregator.
// It stores the snapshot and then initializes the aggregator with the snapshot.
func (es *EventStore) StoreSnapshot(ctx context.Context, aggregateID uuid.UUID, aggregate interface{}) error {
	if reflect.ValueOf(aggregate).Kind() != reflect.Ptr {
		return ErrAggregateMustBePointer
	}

	var aggregateType string
	if a, ok := aggregate.(AggregateType); ok {
		aggregateType = a.AggregateType()
	} else {
		aggregateType = utils.FullyQualifiedStructName(aggregate)
	}

	var aggregateState []byte
	if a, ok := aggregate.(GetSnapshotData); ok {
		var err error
		aggregateState, err = a.GetSnapshotData()
		if err != nil {
			return fmt.Errorf("failed to get aggregate state: %w", err)
		}
	} else {
		var err error
		aggregateState, err = json.Marshal(aggregate)
		if err != nil {
			return fmt.Errorf("failed to marshal aggregate: %w", err)
		}
	}

	var latestEventVersion int32
	if a, ok := aggregate.(LastEventVersion); ok {
		latestEventVersion = a.LastEventVersion()
	} else {
		return ErrLastEventVersion
	}

	tx, err := es.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // nolint:errcheck

	repo := es.repo.WithTx(tx)

	if _, err := repo.StoreSnapshot(ctx, StoreSnapshotParams{
		AggregateID:        aggregateID,
		AggregateType:      aggregateType,
		SnapshotData:       aggregateState,
		SnapshotTime:       time.Now().UnixNano(),
		LatestEventVersion: latestEventVersion,
	}); err != nil {
		return fmt.Errorf("failed to store snapshot: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// MakeSnapshot makes a snapshot for the given aggregate id.
// Helpful when you want to load all events since the snapshot and apply them to the aggregator.
// It loads the latest snapshot and all events since the snapshot version.
// Then it applies all events to the snapshot and stores the new snapshot.
func (es *EventStore) MakeSnapshot(ctx context.Context, aggregateID uuid.UUID, aggregate interface{}) error {
	var aggregateType string
	if a, ok := aggregate.(AggregateType); ok {
		aggregateType = a.AggregateType()
	} else {
		aggregateType = utils.FullyQualifiedStructName(aggregate)
	}

	tx, err := es.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // nolint:errcheck

	repo := es.repo.WithTx(tx)

	// load latest snapshot for the given aggregate id
	snapshot, err := repo.LoadLatestSnapshot(ctx, LoadLatestSnapshotParams{
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to load latest snapshot: %w", err)
		}
	}

	if snapshot.SnapshotID != uuid.Nil && snapshot.SnapshotData != nil {
		// init aggregate with snapshot
		if a, ok := aggregate.(LoadAggregateFromSnapshot); ok {
			if err := a.LoadSnapshot(snapshot); err != nil {
				return fmt.Errorf("failed to load aggregate from snapshot: %w", err)
			}
		} else {
			if err := json.Unmarshal(snapshot.SnapshotData, aggregate); err != nil {
				return fmt.Errorf("failed to unmarshal snapshot data: %w", err)
			}
		}
	}

	// if there is no snapshot, we can just load all events
	events, err := repo.LoadNewestEvents(ctx, LoadNewestEventsParams{
		AggregateID:        aggregateID,
		LatestEventVersion: snapshot.LatestEventVersion,
	})
	if err != nil {
		return fmt.Errorf("failed to load newest events: %w", err)
	}

	// nothing to do if there are no events
	if len(events) == 0 {
		return nil
	}

	// apply events to the aggregate
	if a, ok := aggregate.(ApplyEventToAggregate); ok {
		for _, event := range events {
			if err := a.ApplyEvent(event); err != nil {
				return fmt.Errorf("failed to apply event %s: %w", event.EventType, err)
			}
		}
	} else {
		return ErrApplyEventToAggregate
	}

	var aggregateState []byte
	if a, ok := aggregate.(GetSnapshotData); ok {
		var err error
		aggregateState, err = a.GetSnapshotData()
		if err != nil {
			return fmt.Errorf("failed to get aggregate state: %w", err)
		}
	} else {
		var err error
		aggregateState, err = json.Marshal(aggregate)
		if err != nil {
			return fmt.Errorf("failed to marshal aggregate: %w", err)
		}
	}

	snapshot, err = repo.StoreSnapshot(ctx, StoreSnapshotParams{
		AggregateID:        aggregateID,
		AggregateType:      aggregateType,
		SnapshotData:       aggregateState,
		SnapshotTime:       time.Now().UnixNano(),
		LatestEventVersion: events[len(events)-1].EventVersion,
	})
	if err != nil {
		return fmt.Errorf("failed to store snapshot: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
