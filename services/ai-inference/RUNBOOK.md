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

For the selected ConvNeXt-Tiny MVP classifier, prefer a GPU runtime. CPU-only
retraining can take long enough that Colab or another NVIDIA GPU machine is the
practical path.

Minimal Colab command sequence:

```bash
cd /content/nutri-scan/services/ai-inference
REQUIRE_CUDA=1 INSTALL_DEPS=1 bash scripts/colab_retrain_baseline_v2.sh
```

Before training in Colab, verify that the cloned branch contains the selected
augmentation recipe in `configs/selected_mvp_classifier.json`:

```txt
random_resized_crop_scale = [0.55, 1.0]
rotation_degrees = 15
color_jitter = 0.25 / 0.25 / 0.2
random_erasing_p = 0.0
```

The helper installs training dependencies, fails fast when CUDA is unavailable,
trains with `configs/selected_mvp_classifier.json` by default, evaluates
predictions, exports misclassified images, and prints top-1, top-3, and
weak-class metrics. Set `CONFIG=...` only when intentionally rerunning an older
baseline or comparison config.

For config-only validation without starting training:

```bash
REQUIRE_CUDA=0 INSTALL_DEPS=0 DRY_RUN_ONLY=1 bash scripts/colab_retrain_baseline_v2.sh
```

## Model Comparison Workflow

Use `MODEL_COMPARISON.md` to track classifier architecture screens and tuning
screens. The comparison flow has two stages:

1. Architecture screen: compare model families against the EfficientNet-B0 v2
   control baseline using the same dataset and training recipe.
2. Tuning screen: tune only the strongest architecture from the architecture
   screen.

The ConvNeXt-Tiny tuning configs are not required before comparing ConvNeXt-Tiny
with the original EfficientNet-B0 v2 baseline. They are only needed after
ConvNeXt-Tiny has already been selected as the tuning candidate.

Architecture-screen configs:

```txt
configs/model_comparison_mobilenetv3_large.json
configs/model_comparison_efficientnet_b2.json
configs/model_comparison_convnext_tiny.json
```

ConvNeXt-Tiny tuning configs:

```txt
configs/convnext_tiny_tune_lr5e5.json
configs/convnext_tiny_tune_img256.json
configs/convnext_tiny_tune_lr5e5_img256.json
```

Selected MVP classifier config:

```txt
configs/selected_mvp_classifier.json
```

This selected config uses `convnext_tiny.fb_in1k`, `image_size=256`,
`learning_rate=0.0001`, active `label_smoothing=0.1`, and the selected strong
context augmentation settings.
Before running more tuning, follow the further tuning guardrails in
`MODEL_COMPARISON.md`: review selected-model errors, define the next objective,
and avoid repeatedly selecting models from the same held-out test set.

Run any comparison config in Colab by overriding `CONFIG`:

```bash
cd /content/nutri-scan/services/ai-inference
CONFIG=configs/model_comparison_convnext_tiny.json \
PROCESSED_DIR=data/processed-v0.2 \
REQUIRE_CUDA=1 \
INSTALL_DEPS=1 \
bash scripts/colab_retrain_baseline_v2.sh
```

After each run, copy only the generated metric values into `MODEL_COMPARISON.md`.
Keep `reports/`, `model-artifacts/`, dataset images, and ZIP files out of Git.

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

Current v0.2 metrics are not final until baseline v2 is retrained on
`data/processed-v0.2` and evaluated against the cleaned held-out test split.

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
