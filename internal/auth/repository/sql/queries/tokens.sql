-- name: StoreToken :exec
-- Store a token
INSERT INTO tokens (user_id, access_token_id, refresh_token_id, metadata) VALUES (
    @user_id, @access_token_id, @refresh_token_id, @metadata::json
) RETURNING *;

-- name: FindTokenByAccessTokenID :one
-- Find a token by access token ID
SELECT * FROM tokens WHERE access_token_id = @access_token_id AND access_expires_at > now();

-- name: FindTokenByRefreshTokenID :one
-- Find a token by refresh token ID
SELECT * FROM tokens WHERE refresh_token_id = @refresh_token_id AND refresh_expires_at > now();

-- name: RefreshToken :exec
-- Refresh a token
UPDATE tokens SET 
    access_token_id = @access_token_id, 
    access_expires_at = @access_expires_at, 
    refresh_token_id = @refresh_token_id, 
    refresh_expires_at = @refresh_expires_at 
WHERE refresh_token_id = @old_refresh_token_id AND refresh_expires_at > now();

-- name: DeleteTokenByAccessTokenID :exec
-- Delete a token by access token ID
DELETE FROM tokens WHERE access_token_id = @access_token_id;

-- name: DeleteTokenByRefreshTokenID :exec
-- Delete a token by refresh token ID
DELETE FROM tokens WHERE refresh_token_id = @refresh_token_id;

-- name: DeleteTokensByUserID :exec
-- Delete all tokens for a user
DELETE FROM tokens WHERE user_id = @user_id;

-- name: CleanUpTokens :exec
-- Clean up expired tokens
DELETE FROM tokens WHERE refresh_expires_at < now();