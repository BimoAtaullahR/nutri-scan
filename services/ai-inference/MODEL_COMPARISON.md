# NutriScan Model Comparison

This document tracks lightweight classifier architecture screens for the MVP
single primary food recognizer.

## Scope

- Dataset: `data/processed-v0.2`
- Control config: `configs/baseline_training_v2.json`
- Comparison configs:
  - `configs/model_comparison_mobilenetv3_large.json`
  - `configs/model_comparison_efficientnet_b2.json`
  - `configs/model_comparison_convnext_tiny.json`
- Primary selection signal: weak-class F1 for `rendang`, `gado_gado`, and `soto`
- Guardrails: top-1 accuracy, top-3 accuracy, model size, and serving practicality

## Run Order

1. Train the control baseline if no current v0.2 reference exists.
2. Train `mobilenetv3_large_100.ra_in1k`.
3. Train `efficientnet_b2.ra_in1k`.
4. Train `convnext_tiny.fb_in1k` only if the first two comparisons do not give a
   clear enough signal.

Use the same dataset, split, seed, optimizer, scheduler, epoch budget, and
image size for the first pass. Tune only the strongest candidate after the
architecture screen.

## Colab Command

Run from `services/ai-inference` in Colab:

```bash
CONFIG=configs/model_comparison_mobilenetv3_large.json \
PROCESSED_DIR=data/processed-v0.2 \
REQUIRE_CUDA=1 \
INSTALL_DEPS=1 \
bash scripts/colab_retrain_baseline_v2.sh
```

Change `CONFIG` for the next candidate.

## Comparison Table

Fill this table after each Colab run. Do not commit generated `reports/`,
`model-artifacts/`, dataset images, or ZIP files.

| Run | Config | Model | Top-1 | Top-3 | Rendang F1 | Gado-gado F1 | Soto F1 | Weak-class Avg F1 | Decision | Notes |
| --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Control | `configs/baseline_training_v2.json` | `efficientnet_b0` | TBD | TBD | TBD | TBD | TBD | TBD | reference | Colab baseline v0.2 |
| 1 | `configs/model_comparison_mobilenetv3_large.json` | `mobilenetv3_large_100.ra_in1k` | TBD | TBD | TBD | TBD | TBD | TBD | pending | Lightweight MVP candidate |
| 2 | `configs/model_comparison_efficientnet_b2.json` | `efficientnet_b2.ra_in1k` | TBD | TBD | TBD | TBD | TBD | TBD | pending | Stronger EfficientNet candidate |
| 3 | `configs/model_comparison_convnext_tiny.json` | `convnext_tiny.fb_in1k` | TBD | TBD | TBD | TBD | TBD | TBD | optional | Heavier modern classifier |

## Decision Rule

Prefer a candidate that improves weak-class average F1 without meaningfully
reducing global accuracy. Treat top-1 drops greater than one percentage point
from the control baseline as a warning unless the weak-class gain is large and
product-relevant.

A candidate is eligible for tuning when:

- weak-class average F1 improves by at least `+0.03` over the control baseline
- top-1 accuracy drops by no more than one percentage point
- top-3 accuracy remains at or above `90%`

If weak-class average F1 improves but top-1 drops by one to two percentage
points, treat the result as a trade-off and inspect the confusion pairs before
tuning. If no candidate beats the control baseline, tune the EfficientNet-B0 v2
baseline first.

Inspect these confusion pairs before choosing a model:

- `soto` predicted as `bakso`
- `bakso` predicted as `soto`
- `rendang` predicted as `nasi_goreng`
- `rendang` predicted as `gado_gado`
- `gado_gado` predicted as `nasi_goreng`
- `gado_gado` predicted as `soto`
