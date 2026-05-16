package summary

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MealSummary struct {
	MealType        string
	EatenEnergyKcal int
	ScanCount       int
}

type Store interface {
	GetDailyMealSummaries(ctx context.Context, anonymousUserID string, dateStart, dateEnd time.Time) ([]MealSummary, error)
}

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (s *PostgresStore) GetDailyMealSummaries(ctx context.Context, anonymousUserID string, dateStart, dateEnd time.Time) ([]MealSummary, error) {
	parsedUserID, err := uuid.Parse(anonymousUserID)
	if err != nil {
		return nil, err
	}

	rows, err := s.pool.Query(ctx, `
		SELECT meal_type, COALESCE(SUM(estimated_energy_kcal), 0), COUNT(*)
		FROM scans
		WHERE anonymous_user_id = $1
		  AND status = 'completed'
		  AND created_at >= $2
		  AND created_at < $3
		GROUP BY meal_type
	`, parsedUserID, dateStart, dateEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []MealSummary
	for rows.Next() {
		var ms MealSummary
		if err := rows.Scan(&ms.MealType, &ms.EatenEnergyKcal, &ms.ScanCount); err != nil {
			return nil, err
		}
		summaries = append(summaries, ms)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return summaries, nil
}
