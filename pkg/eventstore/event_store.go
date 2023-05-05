package eventstore

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

const loadAllEvents = `-- name: LoadAllEvents :many
SELECT event_id, aggregate_id, event_type, event_version, event_data, event_time FROM %s 
WHERE aggregate_id = $1
ORDER BY event_time ASC
`

func (q *queries) LoadAllEvents(ctx context.Context, aggregateID uuid.UUID) ([]Event, error) {
	query := fmt.Sprintf(loadAllEvents, q.eventTableName)
	rows, err := q.db.QueryContext(ctx, query, aggregateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Event{}
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.EventID,
			&i.AggregateID,
			&i.EventType,
			&i.EventVersion,
			&i.EventData,
			&i.EventTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const loadEventsRange = `-- name: LoadEventsRange :many
SELECT event_id, aggregate_id, event_type, event_version, event_data, event_time FROM %s 
WHERE aggregate_id = $1 AND event_version >= $2 AND event_version <= $3
ORDER BY event_time ASC
`

type LoadEventsRangeParams struct {
	AggregateID      uuid.UUID `json:"aggregate_id"`
	FromEventVersion int32     `json:"from_event_version"`
	ToEventVersion   int32     `json:"to_event_version"`
}

func (q *queries) LoadEventsRange(ctx context.Context, arg LoadEventsRangeParams) ([]Event, error) {
	query := fmt.Sprintf(loadEventsRange, q.eventTableName)
	rows, err := q.db.QueryContext(ctx, query, arg.AggregateID, arg.FromEventVersion, arg.ToEventVersion)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Event{}
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.EventID,
			&i.AggregateID,
			&i.EventType,
			&i.EventVersion,
			&i.EventData,
			&i.EventTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const loadNewestEvents = `-- name: LoadNewestEvents :many
SELECT event_id, aggregate_id, event_type, event_version, event_data, event_time FROM %s 
WHERE aggregate_id = $1 AND event_version > $2
ORDER BY event_time ASC
`

type LoadNewestEventsParams struct {
	AggregateID        uuid.UUID `json:"aggregate_id"`
	LatestEventVersion int32     `json:"latest_event_version"`
}

func (q *queries) LoadNewestEvents(ctx context.Context, arg LoadNewestEventsParams) ([]Event, error) {
	query := fmt.Sprintf(loadNewestEvents, q.eventTableName)
	rows, err := q.db.QueryContext(ctx, query, arg.AggregateID, arg.LatestEventVersion)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Event{}
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.EventID,
			&i.AggregateID,
			&i.EventType,
			&i.EventVersion,
			&i.EventData,
			&i.EventTime,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const storeEvent = `-- name: StoreEvent :one
INSERT INTO %[1]s (aggregate_id, event_type, event_version, event_data, event_time)
VALUES ($1, $2, COALESCE((SELECT MAX(event_version)+1 FROM %[1]s WHERE aggregate_id = $1),1), $3, $4) RETURNING event_id, aggregate_id, event_type, event_version, event_data, event_time
`

type StoreEventParams struct {
	AggregateID uuid.UUID       `json:"aggregate_id"`
	EventType   string          `json:"event_type"`
	EventData   json.RawMessage `json:"event_data"`
	EventTime   int64           `json:"event_time"`
}

func (q *queries) StoreEvent(ctx context.Context, arg StoreEventParams) (Event, error) {
	query := fmt.Sprintf(storeEvent, q.eventTableName)
	row := q.db.QueryRowContext(ctx, query,
		arg.AggregateID,
		arg.EventType,
		arg.EventData,
		arg.EventTime,
	)
	var i Event
	err := row.Scan(
		&i.EventID,
		&i.AggregateID,
		&i.EventType,
		&i.EventVersion,
		&i.EventData,
		&i.EventTime,
	)
	return i, err
}
