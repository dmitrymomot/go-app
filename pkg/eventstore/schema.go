package eventstore

import (
	"context"
	"fmt"
)

const createEventsTable = `
-- name: CreateEventStoreTable :exec
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS %[1]s (
    event_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    aggregate_id uuid NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    event_version INT NOT NULL,
    event_data BYTEA NOT NULL,
    event_time BIGINT NOT NULL
);
CREATE INDEX IF NOT EXISTS %[1]s_aggregate_id_idx ON %[1]s
    USING BTREE (aggregate_id, event_version ASC);
CLUSTER %[1]s USING %[1]s_aggregate_id_idx;
`

func (q *queries) CreateEventsTable(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, fmt.Sprintf(createEventsTable, q.eventTableName))
	if err != nil {
		return fmt.Errorf("failed to create events table: %w", err)
	}
	return nil
}

const createSnapshotTable = `
-- name: CreateSnapshotStoreTable :exec
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS %[1]s (
    snapshot_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    aggregate_id uuid NOT NULL,
    aggregate_type VARCHAR(100) NOT NULL,
    snapshot_version INT NOT NULL,
    snapshot_data BYTEA NOT NULL,
    snapshot_time BIGINT NOT NULL,
    latest_event_version INT NOT NULL
);
CREATE INDEX IF NOT EXISTS %[1]s_aggregate_id_idx ON %[1]s
    USING BTREE (aggregate_id, aggregate_type, snapshot_version DESC);
CLUSTER %[1]s USING %[1]s_aggregate_id_idx;
`

func (q *queries) CreateSnapshotTable(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, fmt.Sprintf(createSnapshotTable, q.snapshotTableName))
	if err != nil {
		return fmt.Errorf("failed to create snapshot table: %w", err)
	}
	return nil
}
