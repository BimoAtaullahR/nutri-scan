-- +goose Up
CREATE TABLE nudge_responses (
    id UUID PRIMARY KEY,
    scan_id UUID NOT NULL REFERENCES scans(id) ON DELETE CASCADE,
    anonymous_user_id UUID NOT NULL REFERENCES anonymous_users(id) ON DELETE CASCADE,
    response TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX nudge_responses_anonymous_user_created_at_idx ON nudge_responses (anonymous_user_id, created_at DESC);

-- +goose Down
DROP TABLE nudge_responses;
