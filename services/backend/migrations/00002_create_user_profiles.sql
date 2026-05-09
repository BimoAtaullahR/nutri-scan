-- +goose Up
CREATE TABLE user_profiles (
    anonymous_user_id UUID PRIMARY KEY REFERENCES anonymous_users(id) ON DELETE CASCADE,
    height_cm NUMERIC(5, 2) NOT NULL CHECK (height_cm > 0),
    weight_kg NUMERIC(5, 2) NOT NULL CHECK (weight_kg > 0),
    age_range TEXT,
    bmi NUMERIC(5, 2) NOT NULL CHECK (bmi > 0),
    bmi_category TEXT NOT NULL CHECK (bmi_category IN ('underweight', 'normal', 'overweight', 'obese')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE user_profiles;
