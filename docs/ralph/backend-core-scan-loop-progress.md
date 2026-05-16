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
- [ ] #19 Finalize Backend API Contracts For Core Scan Loop
- [ ] #20 Create Sync-First Scan With Fakeable AI Client
- [ ] #21 Produce Backend-Owned Nudge Decisions
- [ ] #22 Record Nudge Responses
- [ ] #23 Expose Daily Energy Summary And Meal Energy Summary
- [ ] #24 Expose Weekly Energy Trend
- [ ] #25 Wire End-to-End Core Scan Loop Smoke Test

## Run Log

- 2026-05-16: Completed #17. Anonymous User creation now persists a UUID identity and hashed bearer token through a PostgreSQL-backed store, adds bearer-token authentication middleware for `/me/profile`, and verifies token hashing/auth behavior with focused backend tests. Ran `GOCACHE="$PWD/.gocache" GONOSUMDB='*' GOPROXY=off GOTOOLCHAIN=local go test -mod=readonly ./...` in `services/backend`.
- 2026-05-16: Completed #18. User Profile read/update now requires the Anonymous User bearer token, persists height, weight, optional age range, derived BMI, and BMI Category, rejects out-of-range measurements before persistence, and returns profile responses without blocking the Core Scan Loop. Ran `GOCACHE="$PWD/.gocache" GONOSUMDB='*' GOPROXY=off GOTOOLCHAIN=local go test -mod=readonly ./internal/user` and `GOCACHE="$PWD/.gocache" GONOSUMDB='*' GOPROXY=off GOTOOLCHAIN=local go test -mod=readonly ./...` in `services/backend`.
