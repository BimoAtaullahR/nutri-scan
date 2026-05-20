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

## Comparison Stages

There are two separate comparison stages:

1. **Architecture screen** compares different model families against the
   EfficientNet-B0 v2 control baseline. This stage does not require ConvNeXt-Tiny
   tuning runs. It answers: "Which architecture is worth tuning?"
2. **Tuning screen** compares small hyperparameter changes for the winning
   architecture. This stage answers: "Which version of the selected architecture
   should become the MVP candidate?"

Do not wait for the three ConvNeXt-Tiny tuning configs before comparing the
architecture-screen ConvNeXt-Tiny result with the project baseline. The current
ConvNeXt-Tiny architecture-screen result is already valid for deciding that
ConvNeXt-Tiny is the next tuning candidate.

Run the three ConvNeXt-Tiny tuning configs only after the architecture screen
selects ConvNeXt-Tiny. Use those tuning results to choose the final ConvNeXt-Tiny
candidate, then compare that tuned candidate back against:

- the EfficientNet-B0 v2 control baseline
- the untuned ConvNeXt-Tiny architecture-screen result
- MobileNetV3-Large as the lightweight fallback

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
| 1 | `configs/model_comparison_mobilenetv3_large.json` | `mobilenetv3_large_100.ra_in1k` | 86.68% | 99.10% | 80.95% | 78.65% | 81.08% | 80.23% | keep as lightweight fallback | Smallest artifact; weaker weak-class recall than ConvNeXt-Tiny |
| 2 | `configs/model_comparison_efficientnet_b2.json` | `efficientnet_b2.ra_in1k` | 85.33% | 97.97% | 79.55% | 75.86% | 81.42% | 78.94% | reject for now | Worse than MobileNetV3-Large and ConvNeXt-Tiny in this screen |
| 3 | `configs/model_comparison_convnext_tiny.json` | `convnext_tiny.fb_in1k` | 91.87% | 98.65% | 87.80% | 90.38% | 88.29% | 88.83% | tune next | Best accuracy and weak-class F1, but largest artifact |

## Architecture Screen Result

ConvNeXt-Tiny is the strongest candidate from this architecture screen. It has
the best top-1 accuracy, keeps top-3 accuracy above the MVP target, and improves
all weak-class F1 values compared with the other comparison candidates.

Artifact sizes:

| Model | Artifact Size |
| --- | ---: |
| `mobilenetv3_large_100.ra_in1k` | 16.25 MB |
| `efficientnet_b2.ra_in1k` | 29.82 MB |
| `convnext_tiny.fb_in1k` | 106.20 MB |

Key confusion-pair counts:

| Pair | MobileNetV3-Large | EfficientNet-B2 | ConvNeXt-Tiny |
| --- | ---: | ---: | ---: |
| `soto` as `bakso` | 3 | 3 | 1 |
| `bakso` as `soto` | 1 | 2 | 2 |
| `rendang` as `nasi_goreng` | 2 | 4 | 1 |
| `rendang` as `gado_gado` | 0 | 0 | 2 |
| `gado_gado` as `nasi_goreng` | 2 | 2 | 1 |
| `gado_gado` as `soto` | 4 | 5 | 1 |

Next recommended step: tune ConvNeXt-Tiny first, while keeping MobileNetV3-Large
as the lightweight deployment fallback if serving size or latency becomes the
binding constraint.

For the current MVP demo phase, prioritize recognition quality over model size.
Use ConvNeXt-Tiny as the primary tuning candidate and keep MobileNetV3-Large as
the fallback if FastAPI serving latency, memory use, or artifact size becomes too
costly for the demo environment.

## ConvNeXt-Tiny Tuning Plan

Run a small tuning screen before changing the selected MVP classifier. Keep the
dataset, split, seed, optimizer, scheduler, epoch budget, weight decay, and label
smoothing fixed unless a later result gives a specific reason to change them.

| Run | Config | Change | Decision |
| --- | --- | --- | --- |
| Control | `configs/model_comparison_convnext_tiny.json` | Architecture screen result | current best |
| Tune 1 | `configs/convnext_tiny_tune_lr5e5.json` | learning rate `0.0001` -> `0.00005` | pending |
| Tune 2 | `configs/convnext_tiny_tune_img256.json` | image size `224` -> `256` | pending |
| Tune 3 | `configs/convnext_tiny_tune_lr5e5_img256.json` | learning rate `0.00005`, image size `256` | pending |

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
