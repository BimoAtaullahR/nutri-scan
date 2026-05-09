-- name: ListCompletedScansForWeeklyTrend :many
SELECT id, anonymous_user_id, estimated_energy_kcal, created_at
FROM scans
WHERE anonymous_user_id = $1
  AND status = 'completed'
  AND created_at >= $2
ORDER BY created_at ASC;
