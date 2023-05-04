package eventstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

// LoadLatestSnapshot loads the latest snapshot for the given aggregate id and aggregate type
func (es *EventStore) LoadLatestSnapshot(ctx context.Context, aggregateID uuid.UUID, aggregateType string) (Snapshot, error) {
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

	return snapshot, nil
}

// LoadSnapshot loads the snapshot for the given aggregate id and snapshot version
func (es *EventStore) LoadSnapshot(ctx context.Context, aggregateID uuid.UUID, aggregateType string, snapshotVersion int32) (Snapshot, error) {
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

	return snapshot, nil
}

// LoadCurrentState loads the current state for the given aggregate id.
// It loads the latest snapshot and all events since the snapshot version.
// Then it applies all events to the snapshot and returns the state.
func (es *EventStore) LoadCurrentState(ctx context.Context, aggregate Aggregator) (Aggregator, error) {
	// load latest snapshot for the given aggregate id
	snapshot, err := es.repo.LoadLatestSnapshot(ctx, LoadLatestSnapshotParams{
		AggregateID:   aggregate.AggregateID(),
		AggregateType: aggregate.AggregateType(),
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return aggregate, fmt.Errorf("failed to load latest snapshot: %w", err)
		}
	}

	// init aggregate with snapshot
	if err := aggregate.Init(snapshot); err != nil {
		return aggregate, fmt.Errorf("failed to init aggregate: %w", err)
	}

	// if there is no snapshot, we can just load all events,
	// otherwise we load all events since the snapshot version
	events, err := es.repo.LoadNewestEvents(ctx, LoadNewestEventsParams{
		AggregateID:        aggregate.AggregateID(),
		LatestEventVersion: aggregate.LatestEventVersion(),
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return aggregate, fmt.Errorf("failed to load newest events: %w", err)
		}
	}

	// aggregate state is up to date
	for _, event := range events {
		if err := aggregate.OnEvent(event); err != nil {
			return aggregate, fmt.Errorf("failed to apply event: %w", err)
		}
	}

	return aggregate, nil
}

// StoreSnapshot stores a snapshot for the given aggregator.
// It stores the snapshot and then initializes the aggregator with the snapshot.
func (es *EventStore) StoreSnapshot(ctx context.Context, agg Aggregator) (Aggregator, error) {
	tx, err := es.db.Begin()
	if err != nil {
		return agg, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	repo := es.repo.WithTx(tx)

	snapshot, err := repo.StoreSnapshot(ctx, StoreSnapshotParams{
		AggregateID:        agg.AggregateID(),
		AggregateType:      agg.AggregateType(),
		SnapshotData:       agg.AggregateState(),
		SnapshotTime:       agg.LatestEventTime(),
		LatestEventVersion: agg.LatestEventVersion(),
	})
	if err != nil {
		return agg, fmt.Errorf("failed to store snapshot: %w", err)
	}

	// init aggregate with snapshot
	if err := agg.Init(snapshot); err != nil {
		return agg, fmt.Errorf("failed to init aggregate: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return agg, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return agg, nil
}

// MakeSnapshot makes a snapshot for the given aggregate id.
// Helpful when you want to load all events since the snapshot and apply them to the aggregator.
// It loads the latest snapshot and all events since the snapshot version.
// Then it applies all events to the snapshot and stores the new snapshot.
func (es *EventStore) MakeSnapshot(ctx context.Context, agg Aggregator) (Aggregator, error) {
	tx, err := es.db.Begin()
	if err != nil {
		return agg, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	repo := es.repo.WithTx(tx)

	// load latest snapshot for the given aggregate id
	snapshot, err := repo.LoadLatestSnapshot(ctx, LoadLatestSnapshotParams{
		AggregateID:   agg.AggregateID(),
		AggregateType: agg.AggregateType(),
	})
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return agg, fmt.Errorf("failed to load latest snapshot: %w", err)
		}
	}

	// init aggregate with snapshot
	if err := agg.Init(snapshot); err != nil {
		return agg, fmt.Errorf("failed to init aggregate: %w", err)
	}

	// if there is no snapshot, we can just load all events
	events, err := repo.LoadNewestEvents(ctx, LoadNewestEventsParams{
		AggregateID:        agg.AggregateID(),
		LatestEventVersion: agg.LatestEventVersion(),
	})
	if err != nil {
		return agg, fmt.Errorf("failed to load newest events: %w", err)
	}

	// apply events to the aggregate
	for _, event := range events {
		if err := agg.OnEvent(event); err != nil {
			return agg, fmt.Errorf("failed to apply event %s: %w", event.EventType, err)
		}
	}

	snapshot, err = repo.StoreSnapshot(ctx, StoreSnapshotParams{
		AggregateID:        agg.AggregateID(),
		AggregateType:      agg.AggregateType(),
		SnapshotData:       agg.AggregateState(),
		SnapshotTime:       agg.LatestEventTime(),
		LatestEventVersion: agg.LatestEventVersion(),
	})
	if err != nil {
		return agg, fmt.Errorf("failed to store snapshot: %w", err)
	}

	// init aggregate with snapshot
	if err := agg.Init(snapshot); err != nil {
		return agg, fmt.Errorf("failed to init aggregate: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return agg, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return agg, nil
}
