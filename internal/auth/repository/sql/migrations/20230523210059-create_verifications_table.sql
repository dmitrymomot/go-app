
-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS verifications (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    verification_type VARCHAR(50) NOT NULL CHECK (verification_type <> ''),
    email VARCHAR(255) NOT NULL CHECK (email <> ''),
    otp_hash bytea NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    expires_at TIMESTAMP NOT NULL DEFAULT now() + INTERVAL '15 minutes'
);
CREATE INDEX IF NOT EXISTS verifications_user_id_idx ON verifications (user_id);
CREATE INDEX IF NOT EXISTS verifications_expiry_idx ON verifications (expires_at);
-- +migrate StatementEnd


-- +migrate Down
DROP TABLE IF EXISTS verifications;