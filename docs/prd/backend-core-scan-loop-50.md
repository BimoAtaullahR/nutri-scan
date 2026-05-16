# PRD: Backend Core Scan Loop 50% Target

## Problem Statement

NutriScan needs enough backend capability to make the designed MVP feel real, not just navigable. The Figma flow shows a user opening the homepage, scanning food, seeing estimated energy and an aura-plate-style nudge, recording a response, then seeing daily and weekly progress reflected in dashboard surfaces.

Current backend state is mostly skeleton: anonymous user creation returns a token shape, but profile, scan, nudge response, daily summary, and weekly trend workflows are not yet usable end to end. This blocks mobile integration and makes the designed product loop impossible to validate.

The 50% backend target is not full backend coverage for every designed screen. It is a usable Core Scan Loop: scan food, receive estimated energy feedback and a Nudge Decision, record response, and see completed scans reflected in daily and weekly summaries.

## Solution

Build the backend Core Scan Loop for the MVP. The backend will support durable Anonymous Users, optional User Profiles, Sync-First Scans, AI/ML Inference orchestration, backend-owned Nudge Decisions, nudge response recording, Daily Energy Summary, Meal Energy Summary, and Weekly Energy Trend data.

The Figma design informs the required backend surfaces:

- Homepage needs daily eaten, remaining, burned placeholder, and meal bucket energy totals.
- Aura plate needs estimated energy feedback and a user-facing nudge, not accurate macro or micronutrient analysis.
- Detail screen needs scan/detail and meal summary data.
- Profile screen needs user profile data for coarse personalization.
- Progress surfaces need completed scan history aggregated into weekly trends.

The backend will allow scanning before profile completion. Without a User Profile, scan feedback returns a Generic Nudge Decision. With a User Profile, scan feedback can return a Personalized Nudge Decision using coarse profile signals.

## User Stories

1. As a new mobile user, I want to create an anonymous identity, so that my scans and progress persist without full registration.
2. As a returning mobile user, I want my bearer token to identify me, so that I can access my own scans, profile, summaries, and trends.
3. As a mobile user, I want to scan food before filling out my profile, so that onboarding does not block the main product loop.
4. As a mobile user, I want to add or update profile attributes, so that future feedback can become more personalized.
5. As a mobile user, I want my profile to show derived BMI category, so that the app can personalize feedback without claiming clinical diagnosis.
6. As a mobile user, I want to upload a food image for scanning, so that NutriScan can estimate energy before I eat.
7. As a mobile user, I want scan creation to return feedback immediately when inference is fast, so that I do not wait through an unnecessary polling flow.
8. As a mobile user, I want slow scans to remain retrievable by ID, so that the app can poll and update when processing finishes.
9. As a mobile user, I want failed scans to preserve failure state, so that the app can show retry behavior instead of losing the attempt.
10. As a mobile user, I want low-confidence inference to complete with a review nudge, so that uncertainty is handled as product feedback rather than a system crash.
11. As a mobile user, I want scan feedback to include estimated energy, so that I understand the likely energy impact of my food.
12. As a mobile user, I want scan feedback to include a clear nudge action, so that I know whether to eat as planned, set aside a portion, or review the result.
13. As a mobile user, I want scan feedback to include estimated prevented energy when relevant, so that the benefit of a nudge feels concrete.
14. As a mobile user, I want scan feedback to say whether it is personalized, so that the app can set expectations accurately.
15. As a mobile user, I want each scan assigned to breakfast, lunch, dinner, or snack, so that my homepage meal sections update correctly.
16. As a mobile user, I want the app to send my intended meal type with a scan, so that meal grouping matches my intent.
17. As a mobile user, I want the backend to assign a meal type if the app omits it, so that scans still appear in summaries.
18. As a mobile user, I want the homepage to show eaten energy for today, so that I can understand my current intake.
19. As a mobile user, I want the homepage to show remaining energy for today, so that I can decide whether to adjust later meals.
20. As a mobile user, I want the homepage to show meal energy totals, so that I can see which meal bucket contributed most.
21. As a mobile user, I want burned energy to be stable even before exercise integration exists, so that the homepage layout can render consistently.
22. As a mobile user, I want to retrieve a scan by ID, so that the app can show detail and polling states.
23. As a mobile user, I want to record whether I followed a nudge, so that NutriScan can measure pre-portioning behavior.
24. As a product team member, I want nudge responses stored, so that Pre-Portioning Compliance Rate can be calculated later.
25. As a mobile user, I want weekly energy trend data from completed scans, so that I can see progress over time.
26. As a backend developer, I want AI/ML Inference isolated behind an adapter, so that scan workflow tests can use fake inference responses.
27. As a backend developer, I want contract updates in OpenAPI, so that mobile and backend agree on request and response shapes.
28. As a backend developer, I want scan image retention avoided by default, so that MVP privacy risk stays low.
29. As a backend developer, I want invalid uploads rejected before inference, so that AI/ML service receives only valid scan images.
30. As a backend developer, I want summary and trend aggregation tested independently, so that dashboard values are trustworthy.

## Implementation Decisions

- The 50% backend target means a usable Core Scan Loop, not all-screen backend coverage.
- Backend remains the owner of Anonymous User identity, Scan Lifecycle, Nudge Decision rules, persistence, daily summaries, and weekly trends.
- AI/ML Inference remains a separate service and returns an Inference Payload. Backend converts that payload into a Nudge Decision.
- Mobile may scan before profile completion. Missing User Profile results in a Generic Nudge Decision.
- User Profile enables Personalized Nudge Decisions, but personalization must stay coarse and non-clinical.
- Scan creation uses Sync-First Scan behavior: create processing scan, call inference, return completed feedback when inference finishes inside the request window, otherwise return a retrievable processing or failed scan state.
- Low-confidence inference is a completed Scan with a Review Food Nudge. Failed Scan is reserved for technical failures.
- Scan image upload uses multipart field `image`, accepts JPEG, PNG, and WebP, rejects empty or invalid images, and caps upload size at 8 MB.
- Scan images are forwarded to AI/ML Inference and are not stored long term by default.
- Scan creation accepts optional Meal Type: breakfast, lunch, dinner, or snack.
- If Meal Type is omitted, backend assigns one using local-time fallback windows: breakfast 05:00-10:59, lunch 11:00-15:59, dinner 16:00-20:59, snack otherwise.
- Daily Energy Summary includes eaten energy, remaining energy, burned energy placeholder, and Meal Energy Summary.
- Remaining energy uses a default daily goal when no profile-derived goal exists.
- Burned energy remains an explicit MVP placeholder until steps or exercise integration exists.
- Macro and micronutrient reports are out of scope for the 50% backend target because current AI contract does not provide defensible macro data.
- Aura plate backend output is a Nudge Decision with message, action, estimated prevented energy when relevant, confidence level, personalization flag, basis, and optional display tags.
- Nudge responses are recorded against scan, nudge, and user ownership for later compliance reporting.
- Weekly Energy Trend is calculated from completed scans only.
- OpenAPI must be updated for finalized mobile-facing request and response contracts.
- Deep modules expected: identity/profile, scan orchestration, inference client adapter, nudge decision engine, daily summary aggregation, weekly trend aggregation.

## Testing Decisions

- Tests should focus on external behavior: HTTP status, response shape, persisted state, lifecycle transitions, nudge outputs, and aggregation results.
- Do not test private helper implementation details directly when module behavior can be tested through stable interfaces.
- User module tests cover anonymous token creation, bearer-token ownership, profile create/update/read, and BMI category derivation.
- Scan module tests cover upload validation, meal type handling, sync-first success, slow inference, inference failure, invalid inference payload, and low-confidence completion.
- Nudge module tests cover generic decisions, personalized decisions, review-food decisions, estimated prevented energy, confidence levels, and action selection.
- Daily summary tests cover eaten energy, remaining energy, burned placeholder, default goal fallback, and meal grouping.
- Trend tests cover weekly aggregation from completed scans only and exclusion of processing or failed scans.
- Handler tests cover auth requirements, request validation, error responses, and JSON response shape.
- Use fake AI client responses for scan orchestration tests.
- Use fake clock for meal type fallback and date-bound summary tests.
- DB integration tests are valuable when PostgreSQL test setup is ready; otherwise, keep repository boundaries explicit and test aggregation logic with controlled data.

## Out of Scope

- Full login, registration, social auth, and account linking.
- Long-term scan image storage.
- Water tracking.
- Achievement journal.
- Accurate macro tracking for protein, carbohydrate, and fat.
- Vitamin, mineral, or clinical nutrient breakdown.
- Medical diagnosis or clinical recommendation.
- Exercise, steps, or burned calorie integration.
- Full support for 50 food categories.
- Food correction UX and manually edited food records.
- B2B corporate wellness or insurance workflows.
- Experiment tracking for AI/ML model development.
- Complete backend coverage for every Figma screen.

## Further Notes

- This PRD follows existing ADRs: separate Backend API and AI/ML Inference, monorepo MVP, OpenAPI/JSON Schema contracts, Flutter mobile app, PostgreSQL persistence, Chi routing, sqlc/pgx database access, Goose migrations, anonymous bearer identity, and no default scan image storage.
- Existing product language should be preserved: Scan, Scan Lifecycle, Nudge Decision, Anonymous User, User Profile, BMI Category, Core Scan Loop, Daily Energy Summary, Meal Energy Summary, Meal Type, Sync-First Scan, Review Food Nudge, and Scan Image.
- Figma source context: NutriScan design file `FxYesyYD0KymEyjg7rD0bm`, node `63:1199`, page/canvas `UI/UX`.
