package database

import (
	"context"
	"errors"
)

var ErrMissingDatabaseURL = errors.New("database url is required")

type DB struct {
	URL string
}

func Open(ctx context.Context, databaseURL string) (*DB, error) {
	if databaseURL == "" {
		return nil, ErrMissingDatabaseURL
	}

	return &DB{URL: databaseURL}, nil
}
