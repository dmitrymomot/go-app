
-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS tokens (
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    access_token_id UUID NOT NULL,
    access_expires_at TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '2 hours',
    refresh_token_id UUID NOT NULL,
    refresh_expires_at TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '30 days',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NULL,
    metadata JSON DEFAULT NULL
);
CREATE INDEX IF NOT EXISTS tokens_user_id_idx ON tokens (user_id);
CREATE INDEX IF NOT EXISTS tokens_access_token_id_idx ON tokens (access_token_id,access_expires_at);
CREATE INDEX IF NOT EXISTS tokens_refresh_token_id_idx ON tokens (refresh_token_id,refresh_expires_at);
CREATE TRIGGER update_tokens_modtime BEFORE
UPDATE ON tokens FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
-- +migrate StatementEnd


-- +migrate Down
-- +migrate StatementBegin
DROP TRIGGER IF EXISTS update_tokens_modtime ON tokens;
DROP TABLE IF EXISTS tokens;
-- +migrate StatementEnd
