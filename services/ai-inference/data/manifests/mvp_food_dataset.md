# MVP Food Dataset Manifest

Reviewed: 2026-05-15

This manifest records candidate data sources and the curated local split counts for NutriScan's MVP food classifier and energy estimate workflow.

## MVP Scope

Target task:

- Closed-set food category classification
- 8 Indonesian food categories
- Estimated Energy Range from lookup table, not calorie regression
- Coarse Portion Estimate with user correction

Machine-readable class mapping:

- `configs/mvp_food_categories.json`

MVP food categories:

- `nasi_goreng`
- `sate`
- `rendang`
- `bakso`
- `gado_gado`
- `soto`
- `pempek`
- `gudeg`

## Recommended Dataset Plan

Use **Indonesian Food Image - Mendeley Data** as the primary dataset because it directly covers all 8 MVP categories and is CC BY 4.0.

Defer `nasi_padang` to future development because it is a mixed meal-level category and needs more careful image curation.

Do not use APAC fallback datasets for MVP training unless Indonesian sources fail. Keep APAC datasets as future research for detection, segmentation, or nutrition estimation.

## Candidate Sources

### 1. Indonesian Food Image - Mendeley Data

URL: https://data.mendeley.com/datasets/vtjd68bmwt  
DOI: https://doi.org/10.17632/vtjd68bmwt.1  
License: CC BY 4.0  
Type: image classification, folder/class style  
Collection: scraped Google Images, manually checked, 70/30 train/test split described by author  

Covered MVP categories:

- `bakso`
- `gado_gado`
- `gudeg`
- `nasi_goreng`
- `pempek`
- `rendang`
- `sate`
- `soto`

Missing MVP categories:

- none

Other classes present:

- `bebek_betutu`
- `rawon`

Use:

- Primary MVP training source
- Primary validation source after re-splitting into train/validation/test

Risks:

- Web-scraped images may include duplicates, watermarks, non-food images, or mislabeled images.
- Need manual cleaning before claiming accuracy.
- No calorie or portion ground truth.

### 2. Padang Cuisine - Kaggle

URL: https://www.kaggle.com/datasets/faldoae/padangfood  
License: check Kaggle page before use  
Type: Padang food image classification dataset  
Collection: Bing Image Search scraping with dish keywords  

Potential MVP use:

- Support `rendang` examples if Mendeley data is weak.

Risks:

- Not needed for MVP after replacing `nasi_padang` with `gudeg`.
- May be individual Padang dishes rather than complete nasi padang plates.
- Web-scraped source quality requires manual review.
- License must be verified on Kaggle before redistribution or model publication.

### 3. Nutricheck Image Dataset - Kaggle

URL: https://www.kaggle.com/datasets/dwiiyy/nutricheck-image-dataset  
License: check Kaggle page before use  
Type: food image classification plus nutrition CSV  
Collection: includes 10 Indonesian food classes and mentions 53 total classes after combining sources  

Directly useful MVP category:

- `soto` / `soto_ayam`

Potentially useful future categories:

- `mie_ayam`
- `nasi_uduk`
- `bubur_ayam`
- `ayam_bakar`
- `ikan_bakar`

Use:

- Secondary source for nutrition lookup research
- Fallback source if adding classes after MVP

Risks:

- Does not directly cover most selected MVP classes based on public dataset card.
- Nutrition CSV source depends on FatSecret API-derived data; verify license and terms before use in product claims.

### 4. Indonesian Food 160-Class Paper

URL: https://beei.org/index.php/EEI/article/view/7996/0  
DOI: https://doi.org/10.11591/eei.v13i5.7996  
License: article CC BY-SA 4.0; dataset/model access must be verified separately  
Type: research dataset/model benchmark  
Reported size: 24,427 images, 160 Indonesian food types  
Reported best model: EfficientNetV2-L, 85.44% accuracy  

Use:

- Research reference for expected architecture and benchmark.
- Possible future source if dataset/model download is available.

Risks:

- Dataset/model availability not confirmed from paper landing page.
- EfficientNetV2-L may be too heavy for low-cost inference if used directly.

### 5. hilmansw/resnet18-food-classifier - Hugging Face

URL: https://huggingface.co/hilmansw/resnet18-food-classifier  
License: Apache 2.0  
Base model: microsoft/resnet-18  
Reported training: fine-tuned on Padang Cuisine Kaggle dataset  
Reported best listed accuracy: 89.45% at epoch 15  

Use:

- Reference implementation/model packaging example.
- Possible quick benchmark for Padang cuisine only.

Risks:

- Class label coverage is Padang-specific, not full MVP.
- Dataset quality and train/test split must be checked before trusting reported accuracy.
- Do not use as final MVP model unless labels match NutriScan categories.

### 6. FoodBD - Bangladeshi Cuisine Dataset

URL: https://data.mendeley.com/datasets/xh3ghf3jbg  
License: CC BY 4.0  
Type: APAC fallback; polygon segmentation, object detection, nutrition labels  
Reported size: 3,523 smartphone-captured meal images; 67 categories; 1,837 images with nutrition information  

Use:

- Future research for segmentation/detection/nutrition workflow design.
- Not recommended for MVP Indonesian classifier training.

Risks:

- Cuisine mismatch with Indonesian demo scope.
- Adds task complexity beyond 2.5-week MVP.

## Initial Class Coverage

| Class | Primary Source | Status | Notes |
| --- | --- | --- | --- |
| `nasi_goreng` | Mendeley Indonesian Food Image | covered | Manual cleaning needed |
| `sate` | Mendeley Indonesian Food Image | covered | Manual cleaning needed |
| `rendang` | Mendeley Indonesian Food Image + Padang Cuisine | covered | Padang source can augment |
| `bakso` | Mendeley Indonesian Food Image | covered | Manual cleaning needed |
| `gado_gado` | Mendeley Indonesian Food Image | covered | Manual cleaning needed |
| `soto` | Mendeley Indonesian Food Image + Nutricheck | covered | Align `soto` vs `soto_ayam` labels |
| `pempek` | Mendeley Indonesian Food Image | covered | Manual cleaning needed |
| `gudeg` | Mendeley Indonesian Food Image | covered | Manual cleaning needed |

Deferred categories:

- `nasi_padang`: future meal-level category after manual curation of complete Padang meal images.

## Required Dataset Build Steps

1. Download candidate datasets.
2. Verify licenses and attribution requirements.
3. Normalize labels to NutriScan class slugs.
4. Remove duplicates, watermarks, menus, drawings, and irrelevant images.
5. Manually review at least 100 images per class.
6. Re-split into train/validation/test:
   - train: 70%
   - validation: 15%
   - test: 15%
7. Record final counts in this manifest.
8. Train baseline model.
9. Evaluate:
   - top-1 accuracy
   - top-3 accuracy
   - confusion matrix
   - per-class precision/recall

## Curation Audit Workflow

After `data/processed` exists and bad examples have been removed manually, run:

```bash
python scripts/curate_dataset.py \
  --processed-dir data/processed \
  --class-map configs/mvp_food_categories.json \
  --report-path reports/dataset-curation/curation_report.json
```

The generated report records train/validation/test counts, total images per class,
minimum-review status, and weak-class risks. Do not commit the generated report;
copy final reviewed counts into the table below.

The audit script accepts either `validation/` or `val/` for the validation split.
The current local processed dataset uses `val/`.

## Final Reviewed Counts

Status: v0.2 reviewed local processed dataset audited on 2026-05-20.

Source review file:

```txt
reports/baseline-food-classifier-v2/misclassified/misclassified_review.csv
```

v0.2 was created by copying `data/processed` to `data/processed-v0.2` and
applying the misclassified review decisions:

- reviewed images: 93
- kept hard examples: 39
- removed ambiguous images: 30
- removed bad-quality images: 23
- removed duplicate images: 1
- relabeled images: 0

Audit command:

```bash
python scripts/curate_dataset.py \
  --processed-dir data/processed-v0.2 \
  --class-map configs/mvp_food_categories.json \
  --report-path reports/dataset-curation/curation_report_v0.2.json
```

Total curated images: 3,262.

| Class | Train | Validation | Test | Total | Reviewed >= 100 | Risks |
| --- | ---: | ---: | ---: | ---: | --- | --- |
| `nasi_goreng` | 359 | 77 | 72 | 508 | Yes | None from audit |
| `sate` | 352 | 75 | 71 | 498 | Yes | None from audit |
| `rendang` | 229 | 49 | 43 | 321 | Yes | None from audit |
| `bakso` | 306 | 66 | 60 | 432 | Yes | None from audit |
| `gado_gado` | 259 | 56 | 50 | 365 | Yes | None from audit |
| `soto` | 334 | 72 | 58 | 464 | Yes | None from audit |
| `pempek` | 281 | 60 | 55 | 396 | Yes | None from audit |
| `gudeg` | 201 | 43 | 34 | 278 | Yes | None from audit |

## Dataset v0.2 Update Flow

Create or refresh the reviewed dataset copy from `services/ai-inference`:

```bash
python scripts/apply_misclassified_review.py \
  --review-csv reports/baseline-food-classifier-v2/misclassified/misclassified_review.csv \
  --source-processed-dir data/processed \
  --output-processed-dir data/processed-v0.2 \
  --report-path reports/dataset-curation/misclassified_review_apply_report.json \
  --force
```

Run the audit after applying the review:

```bash
python scripts/curate_dataset.py \
  --processed-dir data/processed-v0.2 \
  --class-map configs/mvp_food_categories.json \
  --report-path reports/dataset-curation/curation_report_v0.2.json
```

## MVP Acceptance Targets

- 8-class closed-set classification
- top-1 accuracy >= 80% on cleaned held-out test set
- top-3 accuracy >= 90% on cleaned held-out test set
- confidence threshold:
  - `< 0.60`: return `unknown_food`
  - `0.60-0.75`: return prediction with low confidence
  - `> 0.75`: return normal prediction
- Estimated Energy Range produced by lookup table for every non-unknown MVP class

## Energy Lookup Table Work

Create a separate nutrition source manifest before finalizing kcal ranges. Candidate source:

- Nutritional Analysis and Macro-Micro Nutrient Profiling of Indonesian Culinary Recipes  
  URL: https://data.mendeley.com/datasets/8b4ztns76h

Do not claim clinical calorie accuracy. Use ranges and visible "estimate" language.
