# NutriScan AI Inference Runbook

This runbook covers the MVP food recognition flow: prepare data, train, evaluate, serve, and smoke-test.

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

Keep raw and processed datasets local. Do not commit `data/raw/` or `data/processed/`.

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
  validation/<food_category>/
  test/<food_category>/
```

The training script also accepts legacy `val/<food_category>/` as a fallback when
`validation/<food_category>/` is not present.

Issue #5 still requires manual curation before metrics should be trusted.

## Training

```bash
python scripts/train_classifier.py \
  --config configs/baseline_training.json \
  --processed-dir data/processed
```

Artifacts are written to `model-artifacts/baseline-food-classifier/`:

- `model.pt`
- `label_map.json`
- `training_config_resolved.json`

Do not commit model artifacts.

The script selects CUDA automatically when available and falls back to CPU.

## Evaluation

```bash
python scripts/evaluate_model.py \
  --predictions-file reports/baseline-food-classifier/predictions.json \
  --report-dir reports/baseline-food-classifier
```

Reports are written to:

- `reports/baseline-food-classifier/predictions.json`
- `reports/baseline-food-classifier/metrics.json`
- `reports/baseline-food-classifier/per_class_metrics.json`
- `reports/baseline-food-classifier/confusion_matrix.json`

MVP target:

- top-1 accuracy >= 80%
- top-3 accuracy >= 90%

Current metrics are not available until the curated dataset and trained baseline artifact exist.

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
- Real model inference is blocked until dataset curation, training, and artifact export are complete.
- The service returns a deterministic stub prediction when local model artifacts are missing.
- Nasi padang is deferred from the MVP class set.
