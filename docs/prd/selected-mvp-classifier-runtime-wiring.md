# Selected MVP Classifier Runtime Wiring

## Problem Statement

The team has selected a ConvNeXt-Tiny **Model Artifact** for the MVP **Single Primary Food Scan**, but the AI/ML Inference runtime and Backend API scan loop still need to be wired and validated against that selected artifact. From the user's perspective, the risk is that the product appears to scan food while actually running a stale baseline, a mock fallback, or an unverified model package.

The implementation must make the selected **Model Artifact** the runtime source for **Food Category** recognition, preserve the existing **Recognizer Payload** contract, and prove that the Backend API can complete the scan flow using the selected AI/ML Inference service.

## Solution

Wire the selected MVP classifier artifact into AI/ML Inference as the default runtime model, with explicit runtime configuration, artifact validation, model metadata, and readiness checks. AI/ML Inference should return a stable **Recognizer Payload** for confident predictions and **Unknown Food** behavior for low-confidence predictions. Backend API should continue to consume the payload through its existing AI/ML Inference client and complete the **Core Scan Loop** with the appropriate **Nudge Decision**.

This PRD covers runtime wiring and end-to-end validation. It does not cover more model training, new categories, mobile UI work, per-food detection, segmentation, or calorie regression.

## User Stories

1. As a developer, I want AI/ML Inference to load the selected MVP classifier artifact by default, so that runtime predictions use the model selected by evaluation.
2. As a developer, I want the model artifact location to be configurable, so that local development and deployment can use different artifact storage layouts.
3. As a developer, I want AI/ML Inference to validate the required artifact files before serving predictions, so that a missing or stale artifact is caught early.
4. As a developer, I want AI/ML Inference to reject invalid artifact metadata, so that the service does not silently serve a model with the wrong category mapping.
5. As a developer, I want the runtime to expose the active model version, so that scan results can be traced back to the model artifact that produced them.
6. As a backend developer, I want the AI/ML Inference payload to preserve the existing contract, so that Backend API scan orchestration does not need a breaking contract change.
7. As a backend developer, I want confident AI/ML results to include a recognized food category and alternatives, so that Backend API can produce normal scan feedback.
8. As a backend developer, I want low-confidence AI/ML results to return Unknown Food semantics, so that Backend API can complete the scan with a Review Food Nudge.
9. As a backend developer, I want estimated energy ranges to remain lookup-based, so that model predictions do not become a calorie authority.
10. As a backend developer, I want the AI/ML Inference service to fail clearly when the selected artifact is unavailable, so that backend scan failures are technical failures rather than misleading food results.
11. As a product reviewer, I want scan results to include model metadata, so that debugging can distinguish selected model behavior from stale baseline behavior.
12. As an AI developer, I want preprocessing to match the selected model recipe, so that runtime inference is consistent with evaluation.
13. As an AI developer, I want prediction output to include the top food category and top alternatives, so that the product can show confidence and review options.
14. As an AI developer, I want the confidence threshold to be configurable, so that future calibration does not require code changes.
15. As an AI developer, I want the default confidence threshold to remain unchanged for this wiring work, so that threshold calibration does not block runtime integration.
16. As an AI developer, I want the service to expose readiness or model metadata, so that I can verify the active model without sending a scan image.
17. As a developer, I want liveness to remain separate from readiness, so that a running process is not mistaken for a valid loaded model.
18. As a developer, I want fake classifier behavior to be limited to tests, so that production-like runtime cannot accidentally return mock predictions.
19. As a developer, I want the Backend API client to validate representative AI/ML payloads, so that contract drift is caught before scan orchestration fails.
20. As a developer, I want the Backend API scan flow to be smoke-tested against the AI/ML service, so that end-to-end scan behavior is proven with the selected artifact.
21. As a developer, I want a runbook for artifact placement and runtime verification, so that teammates can reproduce the selected runtime setup.
22. As a maintainer, I want model artifacts to remain outside Git, so that large generated files do not enter repository history.
23. As a maintainer, I want the selected runtime recipe documented, so that future model work does not accidentally overwrite the chosen MVP behavior.
24. As a teammate reviewing future errors, I want runtime metadata and model version to be stable, so that error analysis can be tied to the correct model.
25. As a product owner, I want the selected model to be integrated before additional tuning, so that MVP scan feedback can be validated in the real product loop.

## Implementation Decisions

- The selected ConvNeXt-Tiny **Model Artifact** is the runtime source for MVP food recognition.
- The selected **Model Artifact** remains outside Git and is provided through runtime artifact handoff.
- AI/ML Inference reads the model artifact directory from runtime configuration, with a default local artifact location for development.
- The artifact package must include the trained model, label map, and resolved training metadata.
- AI/ML Inference validates that the artifact represents the selected MVP classifier recipe before it is treated as ready.
- AI/ML Inference validates that the label map covers the eight MVP **Food Categories**.
- AI/ML Inference derives or accepts a configurable model version and returns it in every **Recognizer Payload**.
- The model version must match the model metadata exposed by readiness or model metadata endpoints.
- The confidence threshold remains `0.6` by default for this PRD.
- The confidence threshold is configurable so future calibration can be shipped independently.
- Confident predictions return a recognized **Food Category**, alternatives, coarse portion, estimated energy range, low-confidence flag, threshold, and model version.
- Low-confidence predictions return **Unknown Food**, alternatives, no estimated energy range, low-confidence flag, threshold, and model version.
- Estimated energy ranges remain lookup-based and are not predicted by the model.
- AI/ML Inference must not silently fall back to fake predictions in production-like runtime.
- Fake prediction behavior may remain only as an explicit test fixture or dependency-injected fake.
- AI/ML Inference keeps a liveness endpoint for process health.
- AI/ML Inference adds readiness or model metadata behavior that proves the active artifact is valid.
- Readiness or model metadata should expose enough information to identify the active artifact, model version, model name, image size, label count, device, and confidence threshold.
- Backend API continues to call AI/ML Inference through the existing HTTP client boundary.
- Backend API treats low-confidence inference as a completed scan that can produce a **Review Food Nudge**, not as a technical scan failure.
- Backend API treats missing AI service, invalid payloads, or artifact load failures as technical inference failures.
- The end-to-end validation path should prove scan image upload, AI/ML inference, scan completion, and nudge decision behavior.
- Further model tuning is deferred until after runtime wiring and scan-loop validation.

## Testing Decisions

- Tests should validate observable behavior and contracts rather than internal implementation details.
- AI runtime configuration tests should prove that artifact directory, confidence threshold, and model version configuration are resolved correctly.
- AI artifact validation tests should prove that missing files, wrong labels, and stale metadata fail clearly.
- AI prediction tests should use a small fake or dependency-injected classifier where possible, and avoid depending on heavyweight training artifacts unless explicitly marked as smoke tests.
- AI payload tests should cover confident predictions, low-confidence Unknown Food behavior, alternatives, estimated energy range presence or absence, threshold echoing, and model version echoing.
- AI endpoint tests should cover inference response shape and readiness/model metadata behavior.
- Backend client tests should cover valid AI payload decoding, invalid AI payload rejection, and non-2xx AI service responses.
- Backend scan tests should cover confident result completion and low-confidence Review Food Nudge behavior.
- End-to-end smoke testing should be optional or environment-gated when it requires the real selected model artifact.
- Existing AI/ML tests around inference payload, energy lookup, endpoint behavior, and training config validation are prior art for the AI side.
- Existing Backend API tests around scan handler, fake inference client, nudge decision behavior, and core scan loop smoke testing are prior art for the backend side.

## Out of Scope

- Training or tuning another model.
- Reopening architecture comparison between model families.
- Adding new MVP food categories.
- Calibrating a new confidence threshold.
- Changing the mobile UI.
- Adding per-food detection.
- Adding segmentation masks.
- Estimating grams, volume, or exact calories from pixels.
- Using external vision as a primary scan authority.
- Committing large model artifacts to Git.
- Building a full model registry or experiment tracking system.

## Further Notes

- The latest selected model evaluation reached top-1 accuracy of `93.91%`, top-3 accuracy of `99.10%`, weak-class average F1 of `90.23%`, and `27` misclassified images on the held-out test set.
- The final selected model error review found `25` valid model errors, `1` out-of-scope image, and `1` ambiguous image. The next engineering step is runtime wiring, not broad dataset cleanup.
- Remaining future improvement themes are `background_or_context_bias` and `weak_class_overlap`, but those belong after runtime validation.
- No new ADR is required for this PRD because the relevant architectural decisions are already covered by existing AI/ML Inference and MVP scope ADRs.
