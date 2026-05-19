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
- `reports/`
- `model-artifacts/`

## Shared Folder Layout

Use a shared folder such as Google Drive, OneDrive, or another team-owned
storage location:

```txt
NutriScan AI Shared/
  datasets/
    nutriscan-mvp-food-dataset-v0.1-2026-05-15.zip
  model-artifacts/
    baseline-efficientnet-b0-v1.zip
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

Use this version name for the current curated MVP dataset:

```txt
nutriscan-mvp-food-dataset-v0.1-2026-05-15.zip
```

Expected curated counts:

| Class | Train | Validation | Test | Total |
| --- | ---: | ---: | ---: | ---: |
| `nasi_goreng` | 356 | 64 | 87 | 507 |
| `sate` | 353 | 75 | 81 | 509 |
| `rendang` | 231 | 45 | 52 | 328 |
| `bakso` | 309 | 62 | 69 | 440 |
| `gado_gado` | 251 | 33 | 58 | 342 |
| `soto` | 331 | 72 | 109 | 512 |
| `pempek` | 296 | 44 | 84 | 424 |
| `gudeg` | 237 | 43 | 77 | 357 |

Total curated images: 3,419.

## ZIP Owner Workflow

Run from `services/ai-inference`.

1. Verify the processed dataset exists:

```bash
find data/processed -maxdepth 2 -type d | sort
```

2. Audit the curated dataset:

```bash
python scripts/curate_dataset.py \
  --processed-dir data/processed \
  --class-map configs/mvp_food_categories.json \
  --report-path reports/dataset-curation/curation_report.json
```

3. Confirm the audit output matches `data/manifests/mvp_food_dataset.md`.

4. Create a dataset version note beside the ZIP source folder:

```txt
DATASET_VERSION.md
```

Suggested contents:

```md
# NutriScan MVP Food Dataset v0.1

Date: 2026-05-15
Total images: 3,419
Source: Indonesian Food Image - Mendeley Data
License: CC BY 4.0
Curator: <name>

Notes:
- Curated for NutriScan MVP 8-class closed-set classification.
- Raw and processed images must not be committed to Git.
- Counts are recorded in services/ai-inference/data/manifests/mvp_food_dataset.md.
```

5. Create the ZIP from inside `services/ai-inference`:

```bash
zip -r nutriscan-mvp-food-dataset-v0.1-2026-05-15.zip \
  data/processed \
  DATASET_VERSION.md
```

6. Upload the ZIP to the team storage folder, for example Google Drive or OneDrive.

7. Share the download link with teammates.

## Teammate Setup Workflow

Run from `services/ai-inference`.

1. Download the ZIP from team storage.

2. Extract it into `services/ai-inference`:

```bash
unzip /path/to/nutriscan-mvp-food-dataset-v0.1-2026-05-15.zip
```

3. Confirm the folder layout:

```txt
data/processed/
  train/<food_category>/
  validation/<food_category>/
  test/<food_category>/
```

4. Run the audit locally:

```bash
python scripts/curate_dataset.py \
  --processed-dir data/processed \
  --class-map configs/mvp_food_categories.json \
  --report-path reports/dataset-curation/curation_report.json
```

5. Confirm the output matches the manifest counts.

6. Train the baseline classifier:

```bash
python scripts/train_classifier.py \
  --config configs/baseline_training.json \
  --processed-dir data/processed
```

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
- `v0.2`: duplicate cleanup
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
- `train`, `validation`, and `test` folders exist for every class
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
- `reports/`
- `model-artifacts/`
- `*.zip`

If a dataset file appears in Git status, stop and fix `.gitignore` before committing.
