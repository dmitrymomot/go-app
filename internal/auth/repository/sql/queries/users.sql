-- name: CreateUser :one
-- Store or update a user
INSERT INTO users (email) VALUES ($1) ON CONFLICT (email) DO NOTHING RETURNING id;

-- name: FindUserByEmail :one
-- Find a user by email
SELECT * FROM users WHERE email = $1;

-- name: FindUserByID :one
-- Find a user by ID
SELECT * FROM users WHERE id = $1;

-- name: UpdateUserEmailByID :exec
-- Update a user's email by ID
UPDATE users SET email = $1, verified = $3 WHERE id = $2;

-- name: UpdateUserVerificationStatusByID :exec
-- Update a user's verification status by ID
UPDATE users SET verified = $1 WHERE id = $2;

-- name: DeleteUserByID :exec
-- Delete a user by ID
DELETE FROM users WHERE id = $1;