# NutriScan AI Inference Runbook

This runbook covers the MVP food recognition flow: prepare data, apply the
misclassified review cleanup, train, evaluate, serve, and smoke-test.

## Setup

```bash
cd services/ai-inference
python -m venv .venv
source .venv/bin/activate
python -m pip install -e ".[dev]"
```

Windows PowerShell:

```powershell
cd services/ai-inference
python -m venv .venv
.\.venv\Scripts\Activate.ps1
python -m pip install -e ".[dev]"
```

## Dataset Prep

Keep raw and processed datasets local. Do not commit `data/raw/`,
`data/processed/`, `data/processed-v0.2/`, generated reports, or model artifacts.

```bash
python scripts/prepare_dataset.py \
  --raw-dir data/raw \
  --processed-dir data/processed \
  --class-map configs/mvp_food_categories.json \
  --seed 42
```

Expected processed layout:

```txt
data/processed/
  train/<food_category>/
  validation/<food_category>/  # `val/` is also accepted.
  test/<food_category>/
```

## Dataset v0.2 Review Cleanup

The current active dataset is `data/processed-v0.2`. It is derived from
`data/processed` by applying the reviewed baseline v2 misclassification CSV.

```bash
python scripts/apply_misclassified_review.py \
  --review-csv reports/baseline-food-classifier-v2/misclassified/misclassified_review.csv \
  --source-processed-dir data/processed \
  --output-processed-dir data/processed-v0.2 \
  --report-path reports/dataset-curation/misclassified_review_apply_report.json \
  --force
```

Audit the reviewed dataset:

```bash
python scripts/curate_dataset.py \
  --processed-dir data/processed-v0.2 \
  --class-map configs/mvp_food_categories.json \
  --report-path reports/dataset-curation/curation_report_v0.2.json
```

Current v0.2 count: 3,262 images.

Review effect:

- `keep`: 39
- `reject_ambiguous`: 30
- `reject_bad_quality`: 23
- `duplicate`: 1
- `relabel`: 0

## Training

```bash
python scripts/train_classifier.py \
  --config configs/baseline_training_v2.json \
  --processed-dir data/processed-v0.2
```

Artifacts are written to `model-artifacts/baseline-food-classifier-v2/`:

- `model.pt`
- `label_map.json`
- `training_config_resolved.json`

Do not commit model artifacts.

For EfficientNet-B0 baseline v2, prefer a GPU runtime. CPU-only retraining can
take long enough that Colab or another NVIDIA GPU machine is the practical path.

Minimal Colab command sequence:

```bash
cd /content/nutri-scan/services/ai-inference
pip install -q timm scikit-learn pydantic-settings python-multipart "uvicorn[standard]" ruff
pip install -q -e . --no-deps

python scripts/train_classifier.py \
  --config configs/baseline_training_v2.json \
  --processed-dir data/processed-v0.2
```
The script selects CUDA automatically when available and falls back to CPU.

## Model Comparison Workflow

Use `MODEL_COMPARISON.md` as the source of truth for selected-model status,
comparison metrics, and next tuning guardrails.

Selected MVP classifier:

```txt
configs/selected_mvp_classifier.json
```

Current selected recipe:

```txt
convnext_tiny.fb_in1k
image_size=256
learning_rate=0.0001
weight_decay=0.0005
label_smoothing=0.1
```

Run the selected model:

```bash
python scripts/train_classifier.py \
  --config configs/selected_mvp_classifier.json \
  --processed-dir data/processed-v0.2
```

Next context-robustness tuning configs:

```txt
configs/selected_mvp_aug_context_mild.json
configs/selected_mvp_aug_context_strong.json
configs/selected_mvp_aug_random_erasing.json
```

After each run, copy only the metric summary into `MODEL_COMPARISON.md`. Keep
`reports/`, `model-artifacts/`, dataset images, and ZIP files out of Git.

## Evaluation

```bash
python scripts/evaluate_model.py \
  --predictions-file reports/baseline-food-classifier-v2/predictions.json \
  --report-dir reports/baseline-food-classifier-v2
```

Reports are written to:

- `reports/baseline-food-classifier/predictions.json`
- `reports/baseline-food-classifier/metrics.json`
- `reports/baseline-food-classifier/per_class_metrics.json`
- `reports/baseline-food-classifier/confusion_matrix.json`

MVP target:

- top-1 accuracy >= 80%
- top-3 accuracy >= 90%

For current selected-model metrics and model-comparison status, see
`MODEL_COMPARISON.md`.

## Export Misclassified Images

After evaluation, export misclassified test images for manual review:

```cmd
python scripts\export_misclassified.py --predictions-file reports\baseline-food-classifier\predictions.json --output-dir reports\baseline-food-classifier\misclassified
```

The script copies misclassified images into folders grouped by
`<true_label>_as_<predicted_label>` so the team can inspect weak classes and decide whether
to clean the dataset or tune the model.

## Single Image Prediction

```bash
python scripts/predict_image.py \
  --model-path model-artifacts/baseline-food-classifier/model.pt \
  --image-path /path/to/food.jpg
```

Windows one-line command:

```cmd
python scripts/predict_image.py --model-path model-artifacts/baseline-food-classifier/model.pt --image-path C:\path\to\food.jpg
```

## Serving

```bash
python -m uvicorn app.main:app --reload --host 0.0.0.0 --port 8000
```

Health check:

```bash
curl http://localhost:8000/healthz
```

## Smoke Test

```bash
curl -X POST "http://localhost:8000/infer?portion=medium" \
  -F "image=@/path/to/food.jpg"
```

Example response:

```json
{
  "modelVersion": "baseline-0.1.0",
  "foodCategory": {
    "slug": "sate",
    "confidenceScore": 0.61
  },
  "alternatives": [
    {
      "slug": "rendang",
      "confidenceScore": 0.24
    },
    {
      "slug": "bakso",
      "confidenceScore": 0.15
    }
  ],
  "coarsePortion": "medium",
  "estimatedEnergyRange": {
    "minKcal": 400,
    "maxKcal": 600
  },
  "isLowConfidence": false,
  "confidenceThreshold": 0.6
}
```

## Known Limitations

- Estimated energy is a lookup range, not exact calorie detection.
- Real model inference depends on the selected local model artifact being present.
- The service returns a deterministic stub prediction when local model artifacts are missing.
- Nasi padang is deferred from the MVP class set.
