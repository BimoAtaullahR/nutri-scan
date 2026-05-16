package scan

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/BimoAtaullahR/nutri-scan/services/backend/internal/nudge"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ScanStatusProcessing = "processing"
	ScanStatusCompleted  = "completed"
	ScanStatusFailed     = "failed"

	MealTypeBreakfast = "breakfast"
	MealTypeLunch     = "lunch"
	MealTypeDinner    = "dinner"
	MealTypeSnack     = "snack"
)

var ErrScanNotFound = errors.New("scan not found")

type Scan struct {
	ID                   string
	AnonymousUserID      string
	Status               string
	MealType             string
	EstimatedEnergyRange *EnergyRange
	Inference            *InferenceResult
	NudgeDecision        *nudge.Decision
	FailureReason        string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type EnergyRange struct {
	MinKcal int `json:"minKcal"`
	MaxKcal int `json:"maxKcal"`
}

type Store interface {
	CreateProcessingScan(ctx context.Context, scan Scan) error
	CompleteScan(ctx context.Context, anonymousUserID string, scanID string, inference InferenceResult, nudgeDecision nudge.Decision) (Scan, error)
	FailScan(ctx context.Context, anonymousUserID string, scanID string, reason string) (Scan, error)
	GetScan(ctx context.Context, anonymousUserID string, scanID string) (Scan, error)
}

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (s *PostgresStore) CreateProcessingScan(ctx context.Context, scan Scan) error {
	scanID, err := uuid.Parse(scan.ID)
	if err != nil {
		return err
	}
	anonymousUserID, err := uuid.Parse(scan.AnonymousUserID)
	if err != nil {
		return err
	}

	_, err = s.pool.Exec(ctx, `
		INSERT INTO scans (id, anonymous_user_id, status, meal_type)
		VALUES ($1, $2, $3, $4)
	`, scanID, anonymousUserID, ScanStatusProcessing, scan.MealType)
	return err
}

func (s *PostgresStore) CompleteScan(ctx context.Context, anonymousUserID string, scanID string, inference InferenceResult, nudgeDecision nudge.Decision) (Scan, error) {
	inferencePayload, err := json.Marshal(inference)
	if err != nil {
		return Scan{}, err
	}
	nudgeDecisionPayload, err := json.Marshal(nudgeDecision)
	if err != nil {
		return Scan{}, err
	}

	var minKcal any
	var maxKcal any
	var estimatedEnergyKcal any
	if inference.EstimatedEnergyRange != nil {
		minKcal = inference.EstimatedEnergyRange.MinKcal
		maxKcal = inference.EstimatedEnergyRange.MaxKcal
		estimatedEnergyKcal = int(math.Round(float64(inference.EstimatedEnergyRange.MinKcal+inference.EstimatedEnergyRange.MaxKcal) / 2))
	}

	return s.scanRow(ctx, `
		UPDATE scans
		SET status = $3,
			estimated_energy_min_kcal = $4,
			estimated_energy_max_kcal = $5,
			inference_payload = $6,
			estimated_energy_kcal = $7,
			nudge_decision = $8,
			failure_reason = NULL,
			updated_at = now()
		WHERE id = $1 AND anonymous_user_id = $2
		RETURNING id, anonymous_user_id, status, meal_type, estimated_energy_min_kcal, estimated_energy_max_kcal, inference_payload, nudge_decision, failure_reason, created_at, updated_at
	`, scanID, anonymousUserID, ScanStatusCompleted, minKcal, maxKcal, inferencePayload, estimatedEnergyKcal, nudgeDecisionPayload)
}

func (s *PostgresStore) FailScan(ctx context.Context, anonymousUserID string, scanID string, reason string) (Scan, error) {
	return s.scanRow(ctx, `
		UPDATE scans
		SET status = $3,
			failure_reason = $4,
			updated_at = now()
		WHERE id = $1 AND anonymous_user_id = $2
		RETURNING id, anonymous_user_id, status, meal_type, estimated_energy_min_kcal, estimated_energy_max_kcal, inference_payload, nudge_decision, failure_reason, created_at, updated_at
	`, scanID, anonymousUserID, ScanStatusFailed, reason)
}

func (s *PostgresStore) GetScan(ctx context.Context, anonymousUserID string, scanID string) (Scan, error) {
	return s.scanRow(ctx, `
		SELECT id, anonymous_user_id, status, meal_type, estimated_energy_min_kcal, estimated_energy_max_kcal, inference_payload, nudge_decision, failure_reason, created_at, updated_at
		FROM scans
		WHERE id = $1 AND anonymous_user_id = $2
	`, scanID, anonymousUserID)
}

func (s *PostgresStore) scanRow(ctx context.Context, sql string, scanID string, anonymousUserID string, args ...any) (Scan, error) {
	parsedScanID, err := uuid.Parse(scanID)
	if err != nil {
		return Scan{}, err
	}
	parsedAnonymousUserID, err := uuid.Parse(anonymousUserID)
	if err != nil {
		return Scan{}, err
	}

	allArgs := append([]any{parsedScanID, parsedAnonymousUserID}, args...)
	var id uuid.UUID
	var ownerID uuid.UUID
	var scan Scan
	var minKcal *int
	var maxKcal *int
	var inferencePayload []byte
	var nudgeDecisionPayload []byte
	var failureReason *string

	err = s.pool.QueryRow(ctx, sql, allArgs...).Scan(
		&id,
		&ownerID,
		&scan.Status,
		&scan.MealType,
		&minKcal,
		&maxKcal,
		&inferencePayload,
		&nudgeDecisionPayload,
		&failureReason,
		&scan.CreatedAt,
		&scan.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Scan{}, ErrScanNotFound
	}
	if err != nil {
		return Scan{}, err
	}

	scan.ID = id.String()
	scan.AnonymousUserID = ownerID.String()
	if minKcal != nil && maxKcal != nil {
		scan.EstimatedEnergyRange = &EnergyRange{
			MinKcal: *minKcal,
			MaxKcal: *maxKcal,
		}
	}
	if len(inferencePayload) > 0 {
		var inference InferenceResult
		if err := json.Unmarshal(inferencePayload, &inference); err != nil {
			return Scan{}, err
		}
		scan.Inference = &inference
	}
	if len(nudgeDecisionPayload) > 0 {
		var nudgeDecision nudge.Decision
		if err := json.Unmarshal(nudgeDecisionPayload, &nudgeDecision); err != nil {
			return Scan{}, err
		}
		scan.NudgeDecision = &nudgeDecision
	}
	if failureReason != nil {
		scan.FailureReason = *failureReason
	}

	return scan, nil
}
