
-- +migrate Up
-- +migrate StatementBegin
BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS snapshots (
    snapshot_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    aggregate_id uuid NOT NULL,
    aggregate_type VARCHAR(100) NOT NULL,
    snapshot_version INT NOT NULL,
    snapshot_data JSONB NOT NULL,
    snapshot_time BIGINT NOT NULL,
    latest_event_version INT NOT NULL
);
CREATE INDEX IF NOT EXISTS snapshots_aggregate_id_idx ON snapshots
    USING BTREE (aggregate_id, aggregate_type, snapshot_version DESC);
CLUSTER snapshots USING snapshots_aggregate_id_idx;
COMMIT;
-- +migrate StatementEnd

-- +migrate Down
DROP TABLE IF EXISTS snapshots;
