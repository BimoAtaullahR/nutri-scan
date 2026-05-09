# Do Not Store Scan Images by Default

NutriScan will not store uploaded food images long term by default. The Backend API validates the upload, sends it to AI/ML Inference, stores the structured scan result and nudge decision, then discards the image unless an explicit short-retention debugging mode is enabled.

## Considered Options

- Store every scan image for future model improvement.
- Store images only when users opt in.
- Do not store scan images by default, with optional short-retention debugging mode.

## Consequences

- MVP storage and privacy risk stay lower.
- AI/ML debugging cannot assume historical images are always available.
- If image retention is enabled for development or consented workflows, retention and cleanup must be explicit.
