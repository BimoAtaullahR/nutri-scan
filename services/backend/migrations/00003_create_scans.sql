-- +goose Up
CREATE TABLE scans (
    id UUID PRIMARY KEY,
    anonymous_user_id UUID NOT NULL REFERENCES anonymous_users(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('created', 'processing', 'completed', 'failed')),
    estimated_energy_kcal NUMERIC(8, 2),
    inference_payload JSONB,
    nudge_decision JSONB,
    failure_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX scans_anonymous_user_created_at_idx ON scans (anonymous_user_id, created_at DESC);
CREATE INDEX scans_status_idx ON scans (status);

-- +goose Down
DROP TABLE scans;
