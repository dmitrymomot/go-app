package eventstore

import (
	"context"
	"fmt"
)

const createEventsTable = `
-- name: CreateEventStoreTable :exec
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS %[1]s_events (
    event_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    aggregate_id uuid NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    event_version INT NOT NULL,
    event_data JSONB NOT NULL,
    event_time BIGINT NOT NULL
);
CREATE INDEX IF NOT EXISTS %[1]s_events_aggregate_id_idx ON %[1]s_events
    USING BTREE (aggregate_id, event_version ASC);
CLUSTER %[1]s_events USING %[1]s_events_aggregate_id_idx;
`

func (q *queries) CreateEventsTable(ctx context.Context, tableName string) error {
	_, err := q.db.ExecContext(ctx, fmt.Sprintf(createEventsTable, tableName))
	if err != nil {
		return fmt.Errorf("failed to create events table: %w", err)
	}
	return nil
}

const createSnapshotTable = `
-- name: CreateSnapshotStoreTable :exec
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS %[1]s_snapshots (
    snapshot_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    aggregate_id uuid NOT NULL,
    snapshot_version INT NOT NULL,
    snapshot_data JSONB NOT NULL,
    snapshot_time BIGINT NOT NULL
);
CREATE INDEX IF NOT EXISTS %[1]s_snapshots_aggregate_id_idx ON %[1]s_snapshots
    USING BTREE (aggregate_id, snapshot_version DESC);
CLUSTER %[1]s_snapshots USING %[1]s_snapshots_aggregate_id_idx;
`

func (q *queries) CreateSnapshotTable(ctx context.Context, tableName string) error {
	_, err := q.db.ExecContext(ctx, fmt.Sprintf(createSnapshotTable, tableName))
	if err != nil {
		return fmt.Errorf("failed to create snapshot table: %w", err)
	}
	return nil
}
