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

## Planned Structure

```txt
cmd/api/               # API entrypoint
internal/scan/         # scan lifecycle, orchestration, persistence
internal/nudge/        # preventive nudge decision rules
internal/trend/        # weekly energy trend reporting
internal/inference/    # client adapter for AI/ML Inference
internal/user/         # user/session model and persistence
internal/platform/     # http, config, database, storage adapters
migrations/            # database migrations
```

Do not place feature code in global `controllers`, `services`, or `repositories` folders. Keep feature behavior near the domain module that owns it.
