package eventstore

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

const loadLatestSnapshot = `-- name: LoadLatestSnapshot :one
SELECT snapshot_id, aggregate_id, aggregate_type, snapshot_version, snapshot_data, snapshot_time, latest_event_version FROM %s 
WHERE aggregate_id=$1 AND aggregate_type=$2
ORDER BY snapshot_version DESC
LIMIT 1
`

type LoadLatestSnapshotParams struct {
	AggregateID   uuid.UUID `json:"aggregate_id"`
	AggregateType string    `json:"aggregate_type"`
}

func (q *queries) LoadLatestSnapshot(ctx context.Context, arg LoadLatestSnapshotParams) (Snapshot, error) {
	query := fmt.Sprintf(loadLatestSnapshot, q.snapshotTableName)
	row := q.db.QueryRowContext(ctx, query, arg.AggregateID, arg.AggregateType)
	var i Snapshot
	err := row.Scan(
		&i.SnapshotID,
		&i.AggregateID,
		&i.AggregateType,
		&i.SnapshotVersion,
		&i.SnapshotData,
		&i.SnapshotTime,
		&i.LatestEventVersion,
	)
	return i, err
}

const loadSnapshot = `-- name: LoadSnapshot :one
SELECT snapshot_id, aggregate_id, aggregate_type, snapshot_version, snapshot_data, snapshot_time, latest_event_version FROM %s 
WHERE aggregate_id=$1 AND aggregate_type=$2 AND snapshot_version=$3
LIMIT 1
`

type LoadSnapshotParams struct {
	AggregateID     uuid.UUID `json:"aggregate_id"`
	AggregateType   string    `json:"aggregate_type"`
	SnapshotVersion int32     `json:"snapshot_version"`
}

func (q *queries) LoadSnapshot(ctx context.Context, arg LoadSnapshotParams) (Snapshot, error) {
	query := fmt.Sprintf(loadSnapshot, q.snapshotTableName)
	row := q.db.QueryRowContext(ctx, query, arg.AggregateID, arg.AggregateType, arg.SnapshotVersion)
	var i Snapshot
	err := row.Scan(
		&i.SnapshotID,
		&i.AggregateID,
		&i.AggregateType,
		&i.SnapshotVersion,
		&i.SnapshotData,
		&i.SnapshotTime,
		&i.LatestEventVersion,
	)
	return i, err
}

const storeSnapshot = `-- name: StoreSnapshot :one
INSERT INTO %[1]s (aggregate_id, aggregate_type, snapshot_version, snapshot_data, snapshot_time, latest_event_version)
VALUES ($1::uuid, $2::varchar, COALESCE((SELECT MAX(snapshot_version)+1 FROM %[1]s WHERE aggregate_id = $1::uuid AND aggregate_type = $2::varchar),1), $3, $4, $5) RETURNING snapshot_id, aggregate_id, aggregate_type, snapshot_version, snapshot_data, snapshot_time, latest_event_version
`

type StoreSnapshotParams struct {
	AggregateID        uuid.UUID       `json:"aggregate_id"`
	AggregateType      string          `json:"aggregate_type"`
	SnapshotData       json.RawMessage `json:"snapshot_data"`
	SnapshotTime       int64           `json:"snapshot_time"`
	LatestEventVersion int32           `json:"latest_event_version"`
}

func (q *queries) StoreSnapshot(ctx context.Context, arg StoreSnapshotParams) (Snapshot, error) {
	query := fmt.Sprintf(storeSnapshot, q.snapshotTableName)
	row := q.db.QueryRowContext(ctx, query,
		arg.AggregateID,
		arg.AggregateType,
		arg.SnapshotData,
		arg.SnapshotTime,
		arg.LatestEventVersion,
	)
	var i Snapshot
	err := row.Scan(
		&i.SnapshotID,
		&i.AggregateID,
		&i.AggregateType,
		&i.SnapshotVersion,
		&i.SnapshotData,
		&i.SnapshotTime,
		&i.LatestEventVersion,
	)
	return i, err
}
