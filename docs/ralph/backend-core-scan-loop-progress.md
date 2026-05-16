# Ralph Progress: Backend Core Scan Loop 50% Target

Source PRD: `docs/prd/backend-core-scan-loop-50.md`
Parent issue: https://github.com/BimoAtaullahR/nutri-scan/issues/16

## Rules For Ralph

- Work backend-first unless issue explicitly requires shared contract updates.
- Do exactly one issue per run.
- Prefer lowest-numbered incomplete issue unless current code makes another issue clearly unblocked and safer.
- Run focused tests for touched backend packages.
- Update this file after each run.
- Commit changes with best-practice conventional commit message.
- Do not implement out-of-scope PRD items.
- Preserve existing ADRs and backend glossary.

## Issue Queue

- [x] #17 Persist Anonymous User Bearer Identity
- [x] #18 Implement User Profile With BMI Category
- [x] #19 Finalize Backend API Contracts For Core Scan Loop
- [x] #20 Create Sync-First Scan With Fakeable AI Client
- [x] #21 Produce Backend-Owned Nudge Decisions
- [x] #22 Record Nudge Responses
- [x] #23 Expose Daily Energy Summary And Meal Energy Summary
- [x] #24 Expose Weekly Energy Trend
- [ ] #25 Wire End-to-End Core Scan Loop Smoke Test

## Run Log

- 2026-05-16: Completed #17. Anonymous User creation now persists a UUID identity and hashed bearer token through a PostgreSQL-backed store, adds bearer-token authentication middleware for `/me/profile`, and verifies token hashing/auth behavior with focused backend tests. Ran `GOCACHE="$PWD/.gocache" GONOSUMDB='*' GOPROXY=off GOTOOLCHAIN=local go test -mod=readonly ./...` in `services/backend`.
- 2026-05-16: Completed #18. User Profile read/update now requires the Anonymous User bearer token, persists height, weight, optional age range, derived BMI, and BMI Category, rejects out-of-range measurements before persistence, and returns profile responses without blocking the Core Scan Loop. Ran `GOCACHE="$PWD/.gocache" GONOSUMDB='*' GOPROXY=off GOTOOLCHAIN=local go test -mod=readonly ./internal/user` and `GOCACHE="$PWD/.gocache" GONOSUMDB='*' GOPROXY=off GOTOOLCHAIN=local go test -mod=readonly ./...` in `services/backend`.
- 2026-05-16: Completed #19. Backend API OpenAPI now defines the Core Scan Loop mobile contract for Anonymous User auth, User Profile, Sync-First Scan upload and polling, Nudge Decision responses, Daily Energy Summary, Meal Energy Summary, and Weekly Energy Trend without adding scan image storage. Ran `ruby -e 'require "psych"; Psych.load_file("packages/contracts/openapi/backend-api.yaml"); puts "openapi yaml ok"'`, `ruby -e 'require "json"; JSON.parse(File.read("packages/contracts/ai-inference/scan-inference.schema.json")); puts "ai schema json ok"'`, and an OpenAPI local `$ref` check.
- 2026-05-16: Completed #20. Sync-First Scan creation now requires the Anonymous User bearer token, validates JPEG/PNG/WebP Scan Images up to 8 MB, assigns Meal Type from the request or fallback windows, persists Scan Lifecycle state, forwards Scan Image bytes to a fakeable AI/ML Inference client without default image storage, returns completed feedback for fast inference, preserves processing state on inference timeout, records failed state for technical or invalid inference failures, and supports owner-scoped Scan retrieval. Ran `GOCACHE="$PWD/.gocache" GONOSUMDB='*' GOPROXY=off GOTOOLCHAIN=local go test -mod=readonly ./internal/scan` and `GOCACHE="$PWD/.gocache" GONOSUMDB='*' GOPROXY=off GOTOOLCHAIN=local go test -mod=readonly ./...` in `services/backend`.
- 2026-05-16: Blocked completing #21 after implementation and tests because `.git` is mounted read-only and `git add` cannot create `.git/index.lock`; leave #21 unchecked until the changes can be committed. Worktree now produces and persists backend-owned Nudge Decisions from inference results, returns Generic Nudge Decisions without a User Profile, returns Personalized Nudge Decisions when profile BMI Category is available, and completes low-confidence scans with a Review Food Nudge instead of treating them as failures. Ran `GOCACHE="$PWD/.gocache" GONOSUMDB='*' GOPROXY=off GOTOOLCHAIN=local go test -mod=readonly ./internal/scan ./internal/nudge` and `GOCACHE="$PWD/.gocache" GONOSUMDB='*' GOPROXY=off GOTOOLCHAIN=local go test -mod=readonly ./...` in `services/backend`.
- 2026-05-16: Unblocked and committed #20 and #21.
- 2026-05-16: Completed #22. Record Nudge Responses endpoint now accepts requests, validates JSON, verifies nudge ownership against the user's completed scans, and persists the response in PostgreSQL `nudge_responses` table. Integrated into the Chi router. Ran `GOCACHE="$PWD/.gocache" GONOSUMDB='*' GOPROXY=off GOTOOLCHAIN=local go test -mod=readonly ./internal/nudge` and `go build ./...` in `services/backend`.
- 2026-05-16: Completed #23. Implemented `internal/summary` module exposing `GET /summaries/daily`. It accepts an optional `date` query param, groups completed scans by meal type from the `scans` table, calculates `eatenEnergyKcal`, defaults `dailyGoalEnergyKcal` to 2000 for MVP, and computes `remainingEnergyKcal`. Registered the route in Chi with Anonymous User authentication. Ran `GOCACHE="$PWD/.gocache" GONOSUMDB='*' GOPROXY=off GOTOOLCHAIN=local go test -mod=readonly ./...` and `go build ./...` in `services/backend`.
- 2026-05-16: Completed #24. Implemented `internal/trend` module exposing `GET /trends/weekly`. It accepts an optional `weekStart` query param (defaulting to the current week's Monday), queries the `scans` table grouping by day using PostgreSQL `date_trunc`, computes `eatenEnergyKcal` and `scanCount` per day, and zero-fills any missing days in the 7-day week window. Registered the route in Chi with Anonymous User authentication. Ran `GOCACHE="$PWD/.gocache" GONOSUMDB='*' GOPROXY=off GOTOOLCHAIN=local go test -mod=readonly ./...` and `go build ./...` in `services/backend`.
