-- name: GetScan :one
SELECT id, anonymous_user_id, status, meal_type, estimated_energy_kcal, estimated_energy_min_kcal, estimated_energy_max_kcal, inference_payload, nudge_decision, failure_reason, created_at, updated_at
FROM scans
WHERE id = $1 AND anonymous_user_id = $2;
