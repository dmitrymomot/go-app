
-- +migrate Up
-- +migrate StatementBegin
BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS event_store (
    event_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    aggregate_id uuid NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    event_version INT NOT NULL,
    event_data JSONB NOT NULL,
    event_time BIGINT NOT NULL
);
CREATE INDEX IF NOT EXISTS event_store_aggregate_id_idx ON event_store
    USING BTREE (aggregate_id, event_version ASC);
CLUSTER event_store USING event_store_aggregate_id_idx;
COMMIT;
-- +migrate StatementEnd

-- +migrate Down
DROP TABLE IF EXISTS event_store;
