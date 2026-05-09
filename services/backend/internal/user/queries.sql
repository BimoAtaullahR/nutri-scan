-- name: GetAnonymousUserByTokenHash :one
SELECT id, token_hash, created_at
FROM anonymous_users
WHERE token_hash = $1;
