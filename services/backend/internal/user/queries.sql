-- name: CreateAnonymousUser :exec
INSERT INTO anonymous_users (id, token_hash)
VALUES ($1, $2);

-- name: GetAnonymousUserByTokenHash :one
SELECT id, token_hash, created_at
FROM anonymous_users
WHERE token_hash = $1;
