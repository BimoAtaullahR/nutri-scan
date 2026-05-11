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

## Model and Dataset Artifacts

Commit code, manifests, and lightweight config only. Do not commit datasets, trained model weights, or generated experiment outputs.

Use `model-artifacts/` only as a local placeholder for runtime files downloaded from an external artifact store.

Keep raw and processed datasets local under `data/raw/` and `data/processed/`. Commit only manifests, scripts, and lightweight config.

## Local Setup

Use Python 3.12 or newer.

```bash
cd services/ai-inference
python -m venv .venv
source .venv/bin/activate
python -m pip install -e ".[dev]"
```

## Local Commands

```bash
python -m pytest
python -m ruff check .
python -m uvicorn app.main:app --reload --host 0.0.0.0 --port 8000
```

Health check:

```bash
curl http://localhost:8000/healthz
```
