# Active AI Dataset

Version: v0.2
Archive: `nutriscan-mvp-food-dataset-v0.2-2026-05-20.zip`
Storage: `NutriScan AI Shared/datasets/`
Total images: 3,262
Classes: 8
Status: active

## MVP Classes

- `nasi_goreng`
- `sate`
- `rendang`
- `bakso`
- `gado_gado`
- `soto`
- `pempek`
- `gudeg`

## Notes

- Use this dataset version for the single primary food classifier MVP.
- v0.2 was created from the reviewed processed dataset after applying
  `reports/baseline-food-classifier-v2/misclassified/misclassified_review.csv`.
- The v0.2 cleanup removed 54 reviewed test images:
  - `reject_ambiguous`: 30
  - `reject_bad_quality`: 23
  - `duplicate`: 1
- No `relabel` decisions were applied in this review batch.
- The local working folder is `data/processed-v0.2/`.
- Do not commit dataset images, generated reports, or model artifacts to Git.
- Official dataset versions are published by the Dataset Owner.
- The shared cloud folder is a dataset drop folder, not a training work directory.
