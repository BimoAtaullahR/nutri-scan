# NutriScan AI/ML Inference

Python inference service for image preprocessing, food recognition, visual dominance detection, confidence scoring, and estimated energy payloads.

## Ownership

AI/ML Inference owns:

- Image preprocessing for inference
- Food category recognition
- Visual dominance detection
- Confidence scoring
- Estimated energy payloads
- Model version metadata

AI/ML Inference does not own:

- User identity
- Scan lifecycle persistence
- Preventive nudge decisions
- Weekly trend calculation

## Planned Structure

```txt
app/main.py
app/api/              # internal HTTP routes and request/response schemas
app/inference/        # runtime inference pipeline
app/models/           # model registry and loading code
app/preprocessing/    # image loading and transforms
app/core/             # config, errors, logging
tests/                # API and inference tests
notebooks/            # experiments only, not runtime dependencies
model-artifacts/      # local placeholder; large models should not be committed
```

Keep experiments in `notebooks/` separate from runtime inference code.
