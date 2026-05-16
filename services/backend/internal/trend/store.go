package trend

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DailyTrendSummary struct {
	Date            time.Time
	EatenEnergyKcal int
	ScanCount       int
}

type Store interface {
	GetWeeklyTrendSummaries(ctx context.Context, anonymousUserID string, weekStart, weekEnd time.Time) ([]DailyTrendSummary, error)
}

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (s *PostgresStore) GetWeeklyTrendSummaries(ctx context.Context, anonymousUserID string, weekStart, weekEnd time.Time) ([]DailyTrendSummary, error) {
	parsedUserID, err := uuid.Parse(anonymousUserID)
	if err != nil {
		return nil, err
	}

	rows, err := s.pool.Query(ctx, `
		SELECT date_trunc('day', created_at) AS day_date, COALESCE(SUM(estimated_energy_kcal), 0), COUNT(*)
		FROM scans
		WHERE anonymous_user_id = $1
		  AND status = 'completed'
		  AND created_at >= $2
		  AND created_at < $3
		GROUP BY date_trunc('day', created_at)
	`, parsedUserID, weekStart, weekEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []DailyTrendSummary
	for rows.Next() {
		var ts DailyTrendSummary
		if err := rows.Scan(&ts.Date, &ts.EatenEnergyKcal, &ts.ScanCount); err != nil {
			return nil, err
		}
		// ensure we convert to local time if date_trunc returns UTC
		ts.Date = ts.Date.In(time.Local)
		summaries = append(summaries, ts)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return summaries, nil
}
