-- name: StoreOrUpdateVerification :one
-- Store or update a user's verification code.
INSERT INTO verifications (user_id, verification_type, email, otp_hash) VALUES ($1, $2, $3, $4) RETURNING id;

-- name: FindVerificationByID :one
-- Find a verification by ID
SELECT * FROM verifications WHERE id = $1 AND expires_at > now();

-- name: DeleteVerificationByID :exec
-- Delete a verification by ID
DELETE FROM verifications WHERE id = $1;

-- name: CleanUpVerifications :exec
-- Clean up expired verifications
DELETE FROM verifications WHERE expires_at < now();