package eventstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

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
		repo:        newQueries(db),
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
	defer tx.Rollback()

	repo := es.repo.WithTx(tx)
	if err := repo.CreateEventsTable(ctx, es.eventStream); err != nil {
		return fmt.Errorf("failed to create events table: %w", err)
	}
	if err := repo.CreateSnapshotTable(ctx, es.eventStream); err != nil {
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
func (es *EventStore) AppendEvent(ctx context.Context, event Event) (Event, error) {
	tx, err := es.db.Begin()
	if err != nil {
		return Event{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	repo := es.repo.WithTx(tx)

	e, err := repo.StoreEvent(ctx, StoreEventParams{
		AggregateID: event.AggregateID,
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

// LoadLatestSnapshot loads the latest snapshot for the given aggregate id
func (es *EventStore) LoadLatestSnapshot(ctx context.Context, aggregateID uuid.UUID) (Snapshot, error) {
	snapshot, err := es.repo.LoadLatestSnapshot(ctx, aggregateID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snapshot{}, ErrNoSnapshotFound
		}
		return Snapshot{}, fmt.Errorf("failed to load latest snapshot: %w", err)
	}

	if snapshot.SnapshotID == uuid.Nil {
		return Snapshot{}, ErrNoSnapshotFound
	}

	return snapshot, nil
}

// LoadSnapshot loads the snapshot for the given aggregate id and snapshot version
func (es *EventStore) LoadSnapshot(ctx context.Context, aggregateID uuid.UUID, snapshotVersion int32) (Snapshot, error) {
	snapshot, err := es.repo.LoadSnapshot(ctx, LoadSnapshotParams{
		AggregateID:     aggregateID,
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

	return snapshot, nil
}

// LoadCurrentState loads the current state for the given aggregate id.
// It loads the latest snapshot and all events since the snapshot version.
// Then it applies all events to the snapshot and returns the state.
func (es *EventStore) LoadCurrentState(ctx context.Context, aggregateID uuid.UUID) (Snapshot, error) {
	// load latest snapshot for the given aggregate id
	snapshot, err := es.repo.LoadLatestSnapshot(ctx, aggregateID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return Snapshot{}, fmt.Errorf("failed to load latest snapshot: %w", err)
		}
	}

	// if there is no snapshot, we can just load all events
	events, err := es.repo.LoadNewestEvents(ctx, LoadNewestEventsParams{
		AggregateID:        aggregateID,
		LatestEventVersion: snapshot.SnapshotVersion,
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return Snapshot{}, fmt.Errorf("failed to load newest events: %w", err)
		}
	}

	// apply events to the snapshot
	snapshot, err = applyEventsToSnapshot(snapshot, events)
	if err != nil {
		return Snapshot{}, fmt.Errorf("failed to apply events to snapshot: %w", err)
	}

	return snapshot, nil
}

// StoreSnapshot stores a snapshot for the given aggregate id
func (es *EventStore) StoreSnapshot(ctx context.Context, snapshot Snapshot) (Snapshot, error) {
	tx, err := es.db.Begin()
	if err != nil {
		return Snapshot{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	repo := es.repo.WithTx(tx)

	snapshot, err = repo.StoreSnapshot(ctx, StoreSnapshotParams{
		AggregateID:     snapshot.AggregateID,
		SnapshotVersion: snapshot.SnapshotVersion,
		SnapshotData:    snapshot.SnapshotData,
		SnapshotTime:    snapshot.SnapshotTime,
	})
	if err != nil {
		return Snapshot{}, fmt.Errorf("failed to store snapshot: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return Snapshot{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return snapshot, nil
}

// MakeSnapshot makes a snapshot for the given aggregate id.
// It loads the latest snapshot and all events since the snapshot version.
// Then it applies all events to the snapshot and stores the new snapshot.
func (es *EventStore) MakeSnapshot(ctx context.Context, aggregateID uuid.UUID) (Snapshot, error) {
	tx, err := es.db.Begin()
	if err != nil {
		return Snapshot{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	repo := es.repo.WithTx(tx)

	// load latest snapshot for the given aggregate id
	snapshot, err := repo.LoadLatestSnapshot(ctx, aggregateID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return Snapshot{}, fmt.Errorf("failed to load latest snapshot: %w", err)
		}
	}

	// if there is no snapshot, we can just load all events
	events, err := repo.LoadNewestEvents(ctx, LoadNewestEventsParams{
		AggregateID:        aggregateID,
		LatestEventVersion: snapshot.SnapshotVersion,
	})
	if err != nil {
		return Snapshot{}, fmt.Errorf("failed to load newest events: %w", err)
	}

	// apply events to the snapshot
	snapshot, err = applyEventsToSnapshot(snapshot, events)
	if err != nil {
		return Snapshot{}, fmt.Errorf("failed to apply events to snapshot: %w", err)
	}

	snapshot, err = repo.StoreSnapshot(ctx, StoreSnapshotParams{
		AggregateID:     snapshot.AggregateID,
		SnapshotVersion: snapshot.SnapshotVersion,
		SnapshotData:    snapshot.SnapshotData,
		SnapshotTime:    snapshot.SnapshotTime,
	})

	if err := tx.Commit(); err != nil {
		return Snapshot{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return snapshot, nil
}

// apply events to snapshot
func applyEventsToSnapshot(snapshot Snapshot, events []Event) (Snapshot, error) {
	// if there are no events since the last snapshot, we can return the snapshot as is
	if len(events) == 0 {
		if snapshot.SnapshotID == uuid.Nil {
			return Snapshot{}, ErrFailedToMakeSnapshot
		}
		return snapshot, nil
	}

	// merge snapshot data with event data
	snapshotData := make(map[string]interface{})
	if snapshot.SnapshotData != nil {
		if err := json.Unmarshal(snapshot.SnapshotData, &snapshotData); err != nil {
			return Snapshot{}, fmt.Errorf("failed to unmarshal snapshot data: %w", err)
		}
	}

	for _, event := range events {
		if event.EventData == nil {
			continue
		}
		eventData := make(map[string]interface{})
		if err := json.Unmarshal(event.EventData, &eventData); err != nil {
			return Snapshot{}, fmt.Errorf("failed to unmarshal event data: %w", err)
		}
		snapshotData = utils.MergeIntoMapRecursively(snapshotData, eventData)
	}

	// store new snapshot
	snapshotDataJSON, err := json.Marshal(snapshotData)
	if err != nil {
		return Snapshot{}, fmt.Errorf("failed to marshal snapshot data: %w", err)
	}

	return Snapshot{
		AggregateID:     snapshot.AggregateID,
		SnapshotVersion: events[len(events)-1].EventVersion,
		SnapshotData:    snapshotDataJSON,
		SnapshotTime:    events[len(events)-1].EventTime,
	}, nil
}
