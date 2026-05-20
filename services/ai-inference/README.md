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

Keep raw and processed datasets local under `data/raw/`, `data/processed/`, and
versioned processed folders such as `data/processed-v0.2/`. Commit only
manifests, scripts, and lightweight config.

For teammate dataset sharing, see `DATASET_COLLABORATION.md`.

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
  validation/<food_category>/  # `val/` is also accepted by training and audit scripts.
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

## Misclassified Review Cleanup

Dataset v0.2 is created from the reviewed misclassification CSV rather than by
editing v0.1 in place. The active reviewed folder is:

```txt
data/processed-v0.2/
```

Apply the review decisions from the baseline v2 misclassification review:

```bash
python scripts/apply_misclassified_review.py \
  --review-csv reports/baseline-food-classifier-v2/misclassified/misclassified_review.csv \
  --source-processed-dir data/processed \
  --output-processed-dir data/processed-v0.2 \
  --report-path reports/dataset-curation/misclassified_review_apply_report.json \
  --force
```

Audit v0.2:

```bash
python scripts/curate_dataset.py \
  --processed-dir data/processed-v0.2 \
  --class-map configs/mvp_food_categories.json \
  --report-path reports/dataset-curation/curation_report_v0.2.json
```

The current v0.2 review removed 54 test images from 93 reviewed
misclassifications and kept 39 hard examples. No relabel decisions were applied.

## Baseline Training

The current MVP baseline v2 uses EfficientNet-B0 with label smoothing and
weight decay from `configs/baseline_training_v2.json`. It writes local artifacts
to `model-artifacts/baseline-food-classifier-v2/`. These artifacts are ignored by git.

Validate training config without a dataset:

```bash
python scripts/train_classifier.py \
  --config configs/baseline_training_v2.json \
  --processed-dir data/processed-v0.2 \
  --dry-run
```

Train once `data/processed-v0.2` exists:

```bash
python scripts/train_classifier.py \
  --config configs/baseline_training_v2.json \
  --processed-dir data/processed-v0.2
```

CPU-only training is slow for EfficientNet-B0. Prefer a GPU runtime such as
Google Colab for baseline v2 retraining.

## Baseline Evaluation

Evaluation reports are local-only and ignored by git. The evaluator writes:

- `reports/baseline-food-classifier-v2/metrics.json`
- `reports/baseline-food-classifier-v2/confusion_matrix.json`

Run evaluation from a predictions file:

```bash
python scripts/evaluate_model.py \
  --predictions-file reports/baseline-food-classifier-v2/predictions.json \
  --report-dir reports/baseline-food-classifier-v2
```

The output states whether the MVP targets are met: top-1 accuracy at least 80% and
top-3 accuracy at least 90%.

## Model Selection and Tuning

Current model-development progress is tracked in `MODEL_COMPARISON.md`.

The selected MVP classifier is `configs/selected_mvp_classifier.json`:

- model: `convnext_tiny.fb_in1k`
- dataset: `data/processed-v0.2`
- image size: `256`
- learning rate: `0.0001`
- weight decay: `0.0005`
- label smoothing: `0.1`

MobileNetV3-Large remains the lightweight fallback if serving latency, memory,
or artifact size becomes more important than recognition quality.

Run the selected model or a tuning candidate with:

```bash
python scripts/train_classifier.py \
  --config configs/selected_mvp_classifier.json \
  --processed-dir data/processed-v0.2
```

The next planned tuning batch is context-robustness augmentation:

```txt
configs/selected_mvp_aug_context_mild.json
configs/selected_mvp_aug_context_strong.json
configs/selected_mvp_aug_random_erasing.json
```

Before adding another tuning axis, review `MODEL_COMPARISON.md` for the current
selection rule and guardrails.

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

The current runtime classifier still defaults to
`model-artifacts/baseline-food-classifier/` and returns a deterministic stub when
local model artifacts are missing. The selected training config writes to
`model-artifacts/selected-mvp-classifier/`; runtime artifact wiring should be
updated separately before product/demo inference uses the selected model.

## Colab GPU Training

When using Google Colab, copy the project and dataset to `/content` before
training. Reading thousands of images directly from Google Drive is slower.

```python
import torch

print(torch.cuda.is_available())
if torch.cuda.is_available():
    print(torch.cuda.get_device_name(0))
```

```bash
cd /content/nutri-scan/services/ai-inference
pip install -q timm scikit-learn pydantic-settings python-multipart "uvicorn[standard]" ruff
pip install -q -e . --no-deps

python scripts/train_classifier.py \
  --config configs/baseline_training_v2.json \
  --processed-dir data/processed-v0.2
```
