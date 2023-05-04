
-- +migrate Up
-- +migrate StatementBegin
BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS snapshot_store (
    snapshot_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    aggregate_id uuid NOT NULL,
    snapshot_version INT NOT NULL,
    snapshot_data JSONB NOT NULL,
    snapshot_time BIGINT NOT NULL
);
CREATE INDEX IF NOT EXISTS snapshot_store_aggregate_id_idx ON snapshot_store
    USING BTREE (aggregate_id, snapshot_version DESC);
CLUSTER snapshot_store USING snapshot_store_aggregate_id_idx;
COMMIT;
-- +migrate StatementEnd

-- +migrate Down
DROP TABLE IF EXISTS snapshot_store;
