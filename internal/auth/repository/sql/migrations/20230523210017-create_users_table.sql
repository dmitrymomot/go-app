
-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL CHECK (email <> ''),
    verified BOOLEAN NOT NULL DEFAULT false,
    updated_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX IF NOT EXISTS users_email_idx ON users (email);
CREATE TRIGGER update_users_modtime BEFORE 
UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
-- +migrate StatementEnd


-- +migrate Down
-- +migrate StatementBegin
DROP TRIGGER IF EXISTS update_users_modtime ON users;
DROP TABLE IF EXISTS users;
-- +migrate StatementEnd