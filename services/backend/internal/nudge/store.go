package nudge

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNudgeNotFound = errors.New("nudge not found")

type ResponseRecord struct {
	ID              string
	ScanID          string
	AnonymousUserID string
	Response        string
	CreatedAt       time.Time
}

type Store interface {
	VerifyNudgeOwnership(ctx context.Context, anonymousUserID string, scanID string, nudgeID string) error
	RecordResponse(ctx context.Context, record ResponseRecord) (ResponseRecord, error)
}

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (s *PostgresStore) VerifyNudgeOwnership(ctx context.Context, anonymousUserID string, scanID string, nudgeID string) error {
	parsedAnonymousUserID, err := uuid.Parse(anonymousUserID)
	if err != nil {
		return err
	}
	parsedScanID, err := uuid.Parse(scanID)
	if err != nil {
		return err
	}

	var foundNudgeID *string
	err = s.pool.QueryRow(ctx, `
		SELECT nudge_decision->>'nudgeId'
		FROM scans
		WHERE id = $1 AND anonymous_user_id = $2
	`, parsedScanID, parsedAnonymousUserID).Scan(&foundNudgeID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNudgeNotFound
	}
	if err != nil {
		return err
	}
	if foundNudgeID == nil || *foundNudgeID != nudgeID {
		return ErrNudgeNotFound
	}

	return nil
}

func (s *PostgresStore) RecordResponse(ctx context.Context, record ResponseRecord) (ResponseRecord, error) {
	recordID, err := uuid.Parse(record.ID)
	if err != nil {
		return ResponseRecord{}, err
	}
	scanID, err := uuid.Parse(record.ScanID)
	if err != nil {
		return ResponseRecord{}, err
	}
	anonymousUserID, err := uuid.Parse(record.AnonymousUserID)
	if err != nil {
		return ResponseRecord{}, err
	}

	err = s.pool.QueryRow(ctx, `
		INSERT INTO nudge_responses (id, scan_id, anonymous_user_id, response)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`, recordID, scanID, anonymousUserID, record.Response).Scan(&record.CreatedAt)
	
	return record, err
}
