package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrAnonymousUserNotFound = errors.New("anonymous user not found")
var ErrUserProfileNotFound = errors.New("user profile not found")

type AnonymousUser struct {
	ID        string
	TokenHash string
	CreatedAt time.Time
}

type UserProfile struct {
	AnonymousUserID string
	HeightCm        float64
	WeightKg        float64
	AgeRange        string
	BMI             float64
	BMICategory     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Store interface {
	CreateAnonymousUser(ctx context.Context, anonymousUser AnonymousUser) error
	GetAnonymousUserByTokenHash(ctx context.Context, tokenHash string) (AnonymousUser, error)
	GetUserProfile(ctx context.Context, anonymousUserID string) (UserProfile, error)
	UpsertUserProfile(ctx context.Context, profile UserProfile) (UserProfile, error)
}

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (s *PostgresStore) CreateAnonymousUser(ctx context.Context, anonymousUser AnonymousUser) error {
	id, err := uuid.Parse(anonymousUser.ID)
	if err != nil {
		return err
	}

	_, err = s.pool.Exec(ctx, `
		INSERT INTO anonymous_users (id, token_hash)
		VALUES ($1, $2)
	`, id, anonymousUser.TokenHash)
	return err
}

func (s *PostgresStore) GetAnonymousUserByTokenHash(ctx context.Context, tokenHash string) (AnonymousUser, error) {
	var id uuid.UUID
	var hash string
	var createdAt time.Time

	err := s.pool.QueryRow(ctx, `
		SELECT id, token_hash, created_at
		FROM anonymous_users
		WHERE token_hash = $1
	`, tokenHash).Scan(&id, &hash, &createdAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return AnonymousUser{}, ErrAnonymousUserNotFound
	}
	if err != nil {
		return AnonymousUser{}, err
	}

	return AnonymousUser{
		ID:        id.String(),
		TokenHash: hash,
		CreatedAt: createdAt,
	}, nil
}

func (s *PostgresStore) GetUserProfile(ctx context.Context, anonymousUserID string) (UserProfile, error) {
	id, err := uuid.Parse(anonymousUserID)
	if err != nil {
		return UserProfile{}, err
	}

	return s.scanUserProfileRow(ctx, `
		SELECT anonymous_user_id, height_cm, weight_kg, COALESCE(age_range, ''), bmi, bmi_category, created_at, updated_at
		FROM user_profiles
		WHERE anonymous_user_id = $1
	`, id)
}

func (s *PostgresStore) UpsertUserProfile(ctx context.Context, profile UserProfile) (UserProfile, error) {
	id, err := uuid.Parse(profile.AnonymousUserID)
	if err != nil {
		return UserProfile{}, err
	}

	return s.scanUserProfileRow(ctx, `
		INSERT INTO user_profiles (anonymous_user_id, height_cm, weight_kg, age_range, bmi, bmi_category)
		VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6)
		ON CONFLICT (anonymous_user_id) DO UPDATE SET
			height_cm = EXCLUDED.height_cm,
			weight_kg = EXCLUDED.weight_kg,
			age_range = EXCLUDED.age_range,
			bmi = EXCLUDED.bmi,
			bmi_category = EXCLUDED.bmi_category,
			updated_at = now()
		RETURNING anonymous_user_id, height_cm, weight_kg, COALESCE(age_range, ''), bmi, bmi_category, created_at, updated_at
	`, id, profile.HeightCm, profile.WeightKg, profile.AgeRange, profile.BMI, profile.BMICategory)
}

func (s *PostgresStore) scanUserProfileRow(ctx context.Context, sql string, args ...any) (UserProfile, error) {
	var anonymousUserID uuid.UUID
	var profile UserProfile

	err := s.pool.QueryRow(ctx, sql, args...).Scan(
		&anonymousUserID,
		&profile.HeightCm,
		&profile.WeightKg,
		&profile.AgeRange,
		&profile.BMI,
		&profile.BMICategory,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return UserProfile{}, ErrUserProfileNotFound
	}
	if err != nil {
		return UserProfile{}, err
	}

	profile.AnonymousUserID = anonymousUserID.String()
	return profile, nil
}
