# Dataset Collaboration with ZIP

This document explains how NutriScan teammates can share the curated MVP food
dataset without committing image files to Git.

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

## Current Dataset Version

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

## Updating the Dataset

When someone changes the dataset, create a new version. Do not overwrite the old ZIP.

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
