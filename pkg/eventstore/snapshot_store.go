package eventstore

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

const loadLatestSnapshot = `-- name: LoadLatestSnapshot :one
SELECT snapshot_id, aggregate_id, snapshot_version, snapshot_data, snapshot_time FROM snapshot_store 
WHERE aggregate_id=$1
ORDER BY snapshot_version DESC
LIMIT 1
`

func (q *queries) LoadLatestSnapshot(ctx context.Context, aggregateID uuid.UUID) (Snapshot, error) {
	row := q.db.QueryRowContext(ctx, loadLatestSnapshot, aggregateID)
	var i Snapshot
	err := row.Scan(
		&i.SnapshotID,
		&i.AggregateID,
		&i.SnapshotVersion,
		&i.SnapshotData,
		&i.SnapshotTime,
	)
	return i, err
}

const loadSnapshot = `-- name: LoadSnapshot :one
SELECT snapshot_id, aggregate_id, snapshot_version, snapshot_data, snapshot_time FROM snapshot_store 
WHERE aggregate_id=$1 AND snapshot_version=$2
LIMIT 1
`

type LoadSnapshotParams struct {
	AggregateID     uuid.UUID `json:"aggregate_id"`
	SnapshotVersion int32     `json:"snapshot_version"`
}

func (q *queries) LoadSnapshot(ctx context.Context, arg LoadSnapshotParams) (Snapshot, error) {
	row := q.db.QueryRowContext(ctx, loadSnapshot, arg.AggregateID, arg.SnapshotVersion)
	var i Snapshot
	err := row.Scan(
		&i.SnapshotID,
		&i.AggregateID,
		&i.SnapshotVersion,
		&i.SnapshotData,
		&i.SnapshotTime,
	)
	return i, err
}

const storeSnapshot = `-- name: StoreSnapshot :one
INSERT INTO snapshot_store (aggregate_id, snapshot_version, snapshot_data, snapshot_time)
VALUES ($1, $2, $3, $4) RETURNING snapshot_id, aggregate_id, snapshot_version, snapshot_data, snapshot_time
`

type StoreSnapshotParams struct {
	AggregateID     uuid.UUID       `json:"aggregate_id"`
	SnapshotVersion int32           `json:"snapshot_version"`
	SnapshotData    json.RawMessage `json:"snapshot_data"`
	SnapshotTime    int64           `json:"snapshot_time"`
}

func (q *queries) StoreSnapshot(ctx context.Context, arg StoreSnapshotParams) (Snapshot, error) {
	row := q.db.QueryRowContext(ctx, storeSnapshot,
		arg.AggregateID,
		arg.SnapshotVersion,
		arg.SnapshotData,
		arg.SnapshotTime,
	)
	var i Snapshot
	err := row.Scan(
		&i.SnapshotID,
		&i.AggregateID,
		&i.SnapshotVersion,
		&i.SnapshotData,
		&i.SnapshotTime,
	)
	return i, err
}
