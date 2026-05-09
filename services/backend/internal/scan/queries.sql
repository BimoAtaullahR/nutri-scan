-- name: GetScan :one
SELECT id, anonymous_user_id, status, estimated_energy_kcal, created_at, updated_at
FROM scans
WHERE id = $1;
