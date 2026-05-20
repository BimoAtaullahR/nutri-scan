# Dataset Collaboration with Shared Drop Folder

This document explains how NutriScan teammates can share the curated MVP food
dataset without committing image files to Git.

The team uses a shared cloud folder as a **Dataset Drop Folder**. It stores
versioned dataset archives, selected model artifacts, and evaluation reports.
It is not a shared working directory for training.

## Rule

Do not commit dataset images, generated reports, or trained model artifacts.

Commit only:

- scripts
- configs
- manifests
- documentation
- small metadata files

Keep these local or shared outside Git:

- `data/raw/`
- `data/processed/`
- `data/processed-*/`
- `reports/`
- `model-artifacts/`

## Shared Folder Layout

Use a shared folder such as Google Drive, OneDrive, or another team-owned
storage location:

```txt
NutriScan AI Shared/
  datasets/
    nutriscan-mvp-food-dataset-v0.1-2026-05-15.zip
    nutriscan-mvp-food-dataset-v0.2-2026-05-20.zip
  model-artifacts/
    baseline-efficientnet-b0-v1.zip
    baseline-efficientnet-b0-v2.zip
  reports/
    baseline-efficientnet-b0-v1/
      metrics.json
      confusion_matrix.json
  notes/
    dataset-change-log.md
    experiment-log.md
```

Rules:

- `datasets/` stores official versioned dataset ZIP files.
- `model-artifacts/` stores model artifacts selected for team sharing.
- `reports/` stores evaluation outputs for shared model runs.
- `notes/` stores dataset and experiment notes.
- Do not train directly from the shared folder.
- Do not edit dataset files in place inside the shared folder.

## Current Dataset Version

The active dataset version is recorded in:

```txt
services/ai-inference/DATASET_VERSION.md
```

Use this version name for the current reviewed MVP dataset:

```txt
nutriscan-mvp-food-dataset-v0.2-2026-05-20.zip
```

Expected v0.2 counts:

| Class | Train | Validation | Test | Total |
| --- | ---: | ---: | ---: | ---: |
| `nasi_goreng` | 359 | 77 | 72 | 508 |
| `sate` | 352 | 75 | 71 | 498 |
| `rendang` | 229 | 49 | 43 | 321 |
| `bakso` | 306 | 66 | 60 | 432 |
| `gado_gado` | 259 | 56 | 50 | 365 |
| `soto` | 334 | 72 | 58 | 464 |
| `pempek` | 281 | 60 | 55 | 396 |
| `gudeg` | 201 | 43 | 34 | 278 |

Total reviewed images: 3,262.

v0.2 was created from the current `data/processed` folder plus the reviewed
baseline v2 misclassification CSV:

```txt
reports/baseline-food-classifier-v2/misclassified/misclassified_review.csv
```

Review decisions applied:

- 93 reviewed misclassified images
- 39 kept as valid hard examples
- 30 rejected as ambiguous
- 23 rejected as bad quality
- 1 rejected as duplicate
- 0 relabeled

## ZIP Owner Workflow

Run from `services/ai-inference`.

1. Verify the processed source dataset exists:

```bash
find data/processed -maxdepth 2 -type d | sort
```

2. Apply the misclassified review into a new processed dataset folder:

```bash
python scripts/apply_misclassified_review.py \
  --review-csv reports/baseline-food-classifier-v2/misclassified/misclassified_review.csv \
  --source-processed-dir data/processed \
  --output-processed-dir data/processed-v0.2 \
  --report-path reports/dataset-curation/misclassified_review_apply_report.json \
  --force
```

3. Audit the reviewed dataset:

```bash
python scripts/curate_dataset.py \
  --processed-dir data/processed-v0.2 \
  --class-map configs/mvp_food_categories.json \
  --report-path reports/dataset-curation/curation_report_v0.2.json
```

4. Confirm the audit output matches `data/manifests/mvp_food_dataset.md`.

5. Create a dataset version note beside the ZIP source folder:

```txt
DATASET_VERSION.md
```

Suggested contents:

```md
# NutriScan MVP Food Dataset v0.2

Date: 2026-05-20
Total images: 3,262
Source: Indonesian Food Image - Mendeley Data
License: CC BY 4.0
Curator: <name>

Notes:
- Curated for NutriScan MVP 8-class closed-set classification.
- Applied baseline v2 misclassified image review decisions.
- Removed ambiguous, bad-quality, and duplicate reviewed test images.
- Raw and processed images must not be committed to Git.
- Counts are recorded in services/ai-inference/data/manifests/mvp_food_dataset.md.
```

6. Create the ZIP from inside `services/ai-inference`:

```bash
zip -r nutriscan-mvp-food-dataset-v0.2-2026-05-20.zip \
  data/processed-v0.2 \
  DATASET_VERSION.md
```

7. Upload the ZIP to the team storage folder, for example Google Drive or OneDrive.

8. Share the download link with teammates.

## Teammate Setup Workflow

Run from `services/ai-inference`.

1. Download the ZIP from team storage.

2. Extract it into `services/ai-inference`:

```bash
unzip /path/to/nutriscan-mvp-food-dataset-v0.2-2026-05-20.zip
```

3. Confirm the folder layout:

```txt
data/processed-v0.2/
  train/<food_category>/
  validation/<food_category>/  # `val/` is also accepted.
  test/<food_category>/
```

4. Run the audit locally:

```bash
python scripts/curate_dataset.py \
  --processed-dir data/processed-v0.2 \
  --class-map configs/mvp_food_categories.json \
  --report-path reports/dataset-curation/curation_report_v0.2.json
```

5. Confirm the output matches the manifest counts.

6. Train the baseline classifier:

```bash
python scripts/train_classifier.py \
  --config configs/baseline_training_v2.json \
  --processed-dir data/processed-v0.2
```

Prefer a GPU runtime such as Google Colab for baseline v2 retraining. CPU-only
EfficientNet-B0 training is slow.

## Sharing Model Artifacts

Do not upload every local experiment to the shared folder. Share only selected
model artifacts that the team may use for backend, mobile, or demo testing.

Rules:

- Put shared model artifacts in `model-artifacts/`.
- Put the matching evaluation report in `reports/<model-version>/`.
- Keep local experiments on the local machine unless the team decides to share them.
- Do not replace a shared model artifact in place; publish a new version instead.

Example:

```txt
NutriScan AI Shared/
  model-artifacts/
    baseline-efficientnet-b0-v1.zip
  reports/
    baseline-efficientnet-b0-v1/
      metrics.json
      confusion_matrix.json
```

## Progress Tracking

Track shared model-development progress in GitHub Issues and PRs. Chat is fine
for coordination, but decisions and completed work should land in the repo.

Suggested issue slices:

- dataset cleanup
- dataset audit
- baseline training
- evaluation report
- confidence threshold calibration
- external vision fallback
- backend integration
- mobile result correction flow

Every PR that changes model behavior, dataset metadata, or evaluation results
should mention the dataset version and model artifact version it used.

## Updating the Dataset

When someone changes the dataset, create a new version. Do not overwrite the old ZIP.

Only the Dataset Owner publishes official dataset ZIP files to `datasets/`.
Other teammates may propose changes, clean data locally, or share candidate
changes, but the official version is the ZIP published by the Dataset Owner.

Version examples:

- `v0.1`: first curated dataset
- `v0.2`: misclassified review cleanup for EfficientNet-B0 baseline v2
- `v0.3`: weak class improvement

Each version update must include:

- new ZIP file
- updated `DATASET_VERSION.md`
- updated `data/manifests/mvp_food_dataset.md`
- audit command output checked locally
- short note in the PR describing what changed

Suggested owner flow:

1. Review the proposed dataset change.
2. Apply accepted changes locally.
3. Run the dataset audit.
4. Update the manifest counts and dataset change log.
5. Create a new ZIP version.
6. Upload the ZIP to `datasets/`.
7. Keep older ZIP versions available until the team agrees they are obsolete.

## What to Check Before Sharing

Before sharing a ZIP, check:

- all 8 MVP classes exist
- `train`, validation (`validation` or `val`), and `test` folders exist for every class
- every class has at least 100 reviewed images when available
- no obvious wrong-label images remain
- no menus, posters, drawings, heavy watermarks, or duplicates remain
- audit risks are either empty or documented in the manifest

## Git Checklist

Before committing:

```bash
git status --short
```

Make sure Git does not show:

- `data/raw/`
- `data/processed/`
- `data/processed-*/`
- `reports/`
- `model-artifacts/`
- `*.zip`

If a dataset file appears in Git status, stop and fix `.gitignore` before committing.
