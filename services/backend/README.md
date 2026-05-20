# NutriScan Backend API

Go Backend API for NutriScan product workflows, persistence, scan orchestration, nudge decisions, and trend APIs.

## Ownership

Backend owns:

- User/session identity
- Scan lifecycle and persistence
- Image upload validation and orchestration
- AI/ML inference client integration
- Nudge decision rules
- Weekly energy trend APIs

Backend does not own:

- Camera UX
- Model inference internals
- Visual dominance detection implementation

## Structure

```txt
cmd/api/               # API entrypoint
internal/scan/         # scan lifecycle, orchestration, persistence, AI/ML Inference client
internal/nudge/        # preventive nudge decision rules
internal/trend/        # weekly energy trend reporting
internal/summary/      # daily and meal energy summary reporting
internal/user/         # user/session model and persistence
internal/platform/     # http, config, database, storage adapters
migrations/            # database migrations
```

Do not place feature code in global `controllers`, `services`, or `repositories` folders. Keep feature behavior near the domain module that owns it.

## Current Endpoints

```txt
GET  /healthz
POST /anonymous-users
GET  /me/profile
PUT  /me/profile
POST /scans
GET  /scans/{scanId}
GET  /summaries/daily
GET  /summaries/meals
GET  /trends/weekly
POST /nudges/{nudgeId}/responses
```

## Local Commands

```bash
go run ./cmd/api
go test ./...
```

When running inside restricted sandboxes, set local cache directories:

```bash
GOCACHE="$PWD/.gocache" GOMODCACHE="$PWD/.gomodcache" go test ./...
```

## Database

Migrations use Goose SQL files in `migrations/`.

```bash
goose postgres "$DATABASE_URL" status
goose postgres "$DATABASE_URL" up
goose postgres "$DATABASE_URL" down
```

Database access is planned through sqlc with generated code under `internal/platform/database/dbgen`.

```bash
sqlc generate
```
