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

## Dataset Preparation

Raw and processed datasets are local-only.

Expected raw layout:

```txt
data/raw/
  nasi_goreng/
  sate/
  rendang/
  bakso/
  gado_gado/
  soto/
  pempek/
  gudeg/
```

Dry run:

```bash
python scripts/prepare_dataset.py \
  --raw-dir data/raw \
  --processed-dir data/processed \
  --class-map configs/mvp_food_categories.json \
  --dry-run
```

Write train/validation/test split:

```bash
python scripts/prepare_dataset.py \
  --raw-dir data/raw \
  --processed-dir data/processed \
  --class-map configs/mvp_food_categories.json \
  --seed 42
```

Output layout:

```txt
data/processed/
  train/<food_category>/
  validation/<food_category>/
  test/<food_category>/
```

## Dataset Curation Audit

After manually removing bad images from `data/processed`, generate a local curation report:

```bash
python scripts/curate_dataset.py \
  --processed-dir data/processed \
  --class-map configs/mvp_food_categories.json \
  --report-path reports/dataset-curation/curation_report.json
```

Use the report to update `data/manifests/mvp_food_dataset.md` with final
train/validation/test counts and weak-class risks. The report is ignored by git;
only reviewed counts and notes should be committed.

## Baseline Training

The MVP baseline uses a lightweight pretrained image classifier and writes local artifacts to
`model-artifacts/baseline-food-classifier/`. These artifacts are ignored by git.

Validate training config without a dataset:

```bash
python scripts/train_classifier.py \
  --config configs/baseline_training.json \
  --processed-dir data/processed \
  --dry-run
```

Train once `data/processed` exists:

```bash
python scripts/train_classifier.py \
  --config configs/baseline_training.json \
  --processed-dir data/processed
```

## Baseline Evaluation

Evaluation reports are local-only and ignored by git. The evaluator writes:

- `reports/baseline-food-classifier/metrics.json`
- `reports/baseline-food-classifier/confusion_matrix.json`

Run evaluation from a predictions file:

```bash
python scripts/evaluate_model.py \
  --predictions-file reports/baseline-food-classifier/predictions.json \
  --report-dir reports/baseline-food-classifier
```

The output states whether the MVP targets are met: top-1 accuracy at least 80% and
top-3 accuracy at least 90%.

## Estimated Energy Ranges

`configs/estimated_energy_ranges.json` maps every MVP food category to approximate
`small`, `medium`, and `large` kcal ranges. These values are lookup estimates for
MVP feedback, not exact calorie detection. They must be reviewed against trusted
nutrition references before production use.

## Inference API

Start the service:

```bash
python -m uvicorn app.main:app --reload --host 0.0.0.0 --port 8000
```

Smoke-test `/infer`:

```bash
curl -X POST http://localhost:8000/infer \
  -F "image=@/path/to/food.jpg" \
  -F "portion=medium"
```

The endpoint looks for local artifacts in `model-artifacts/baseline-food-classifier/`.
If artifacts are missing, it returns a deterministic stub prediction so Backend API
integration can continue before the trained model is available.
