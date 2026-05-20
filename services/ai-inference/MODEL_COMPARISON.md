# NutriScan Model Comparison

This document tracks the current MVP food recognizer selection and the next
safe experiment path.

## Current Best Model

The current selected MVP classifier is:

```txt
configs/selected_mvp_classifier.json
```

Selected recipe:

| Field | Value |
| --- | --- |
| Model | `convnext_tiny.fb_in1k` |
| Dataset | `data/processed-v0.2` |
| Image size | `256` |
| Batch size | `8` |
| Epoch budget | `20` |
| Optimizer | `adamw` |
| Scheduler | `cosine` |
| Learning rate | `0.0001` |
| Weight decay | `0.0005` |
| Label smoothing | `0.1` |
| Seed | `42` |
| Output dir | `model-artifacts/selected-mvp-classifier` |
| Report dir | `reports/selected-mvp-classifier` |

Current selected metrics:

| Metric | Value |
| --- | ---: |
| Top-1 accuracy | 91.87% |
| Top-3 accuracy | 98.87% |
| Rendang F1 | 83.33% |
| Gado-gado F1 | 88.66% |
| Soto F1 | 89.26% |
| Weak-class avg F1 | 87.08% |
| Total misclassified | 36 |

Decision: keep ConvNeXt-Tiny image size `256` as the fixed selected MVP
classifier until a follow-up experiment improves the defined objective without
materially reducing top-1 accuracy.

Runtime note: model selection is documented here, but FastAPI artifact wiring
and serving smoke tests are a separate follow-up. The runtime service must point
to the selected artifact before product/demo inference uses this model.

## Active Dataset

- Dataset version: `v0.2`
- Local folder: `data/processed-v0.2`
- Archive: `nutriscan-mvp-food-dataset-v0.2-2026-05-20.zip`
- Total images: 3,262
- Classes: 8 MVP food categories
- Source manifest: `data/manifests/mvp_food_dataset.md`

Do not commit dataset images, generated reports, model artifacts, or ZIP files.

## Architecture Screen

Architecture-screen runs compare model families against the EfficientNet-B0 v2
control baseline. The primary selection signal is weak-class F1 for `rendang`,
`gado_gado`, and `soto`, with top-1 accuracy, top-3 accuracy, artifact size, and
serving practicality as guardrails.

| Run | Config | Model | Top-1 | Top-3 | Rendang F1 | Gado-gado F1 | Soto F1 | Weak-class Avg F1 | Decision | Notes |
| --- | --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | --- | --- |
| Control | `configs/baseline_training_v2.json` | `efficientnet_b0` | TBD | TBD | TBD | TBD | TBD | TBD | reference | Baseline v2 on dataset v0.2 |
| 1 | `configs/model_comparison_mobilenetv3_large.json` | `mobilenetv3_large_100.ra_in1k` | 86.68% | 99.10% | 80.95% | 78.65% | 81.08% | 80.23% | keep as lightweight fallback | Smallest artifact; weaker weak-class recall than ConvNeXt-Tiny |
| 2 | `configs/model_comparison_efficientnet_b2.json` | `efficientnet_b2.ra_in1k` | 85.33% | 97.97% | 79.55% | 75.86% | 81.42% | 78.94% | reject for now | Worse than MobileNetV3-Large and ConvNeXt-Tiny in this screen |
| 3 | `configs/model_comparison_convnext_tiny.json` | `convnext_tiny.fb_in1k` | 91.87% | 98.65% | 87.80% | 90.38% | 88.29% | 88.83% | tune next | Best accuracy and weak-class F1, but largest artifact |

Artifact sizes:

| Model | Artifact Size |
| --- | ---: |
| `mobilenetv3_large_100.ra_in1k` | 16.25 MB |
| `efficientnet_b2.ra_in1k` | 29.82 MB |
| `convnext_tiny.fb_in1k` | 106.20 MB |

Architecture decision: ConvNeXt-Tiny is the strongest quality candidate.
MobileNetV3-Large remains the fallback if latency, memory, or artifact size
becomes the binding constraint.

Historical note: earlier architecture-screen numbers were recorded before this
trainer reliably used every documented tuning field. Current config files now
include active `label_smoothing` and augmentation fields, and
`scripts/train_classifier.py` reads them.

## ConvNeXt-Tiny Tuning Screen

The first tuning screen changed only learning rate and image size after
ConvNeXt-Tiny won the architecture screen.

| Run | Config | Top-1 | Top-3 | Rendang F1 | Gado-gado F1 | Soto F1 | Weak-class Avg F1 | Priority Confusions | Total Misclassified | Decision |
| --- | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| Control rerun | `configs/convnext_tiny_control_label_smoothing.json` | 90.97% | 98.65% | 83.95% | 91.67% | 86.67% | 87.43% | 6 | 40 | baseline for tuning |
| Tune 1 | `configs/convnext_tiny_tune_lr5e5.json` | 91.20% | 98.19% | 85.37% | 88.66% | 89.26% | 87.76% | 3 | 39 | not selected |
| Tune 2 | `configs/convnext_tiny_tune_img256.json` | 91.87% | 98.87% | 83.33% | 88.66% | 89.26% | 87.08% | 5 | 36 | selected |
| Tune 3 | `configs/convnext_tiny_tune_lr5e5_img256.json` | 91.20% | 98.65% | 81.48% | 90.91% | 87.80% | 86.73% | 2 | 39 | not selected |

Tune 2 is selected because it has the highest top-1 accuracy, highest top-3
accuracy, and lowest total misclassified count. Its weak-class average F1 is
slightly below Tune 1, but the selection rule prioritizes top-1 unless
candidates are within `0.5 pp`.

## Selected Model Error Analysis

Reviewed report:

```txt
reports/convnext-tiny-tune-img256/error_analysis.xlsx
```

Summary:

| Category | Count |
| --- | ---: |
| Total reviewed misclassified images | 36 |
| `model_error` | 33 |
| `out_of_scope` | 2 |
| `ambiguous` | 1 |
| `keep_for_tuning_signal` | 33 |
| `exclude_in_next_dataset` | 3 |

Error type breakdown:

| Error Type | Count |
| --- | ---: |
| `background_or_context_bias` | 14 |
| `weak_class_overlap` | 10 |
| `low_visual_dominance` | 6 |
| `mixed_food` | 2 |
| `presentation_variation` | 2 |
| `image_quality` | 1 |
| `label_noise` | 1 |

Conclusion: remaining errors are mostly valid model errors, not broad dataset
corruption. The next improvement path should focus on context robustness and
class-overlap handling rather than a large dataset cleanup.

## Next Tuning Batch

The next tuning batch targets `background_or_context_bias`, the largest reviewed
error type for the selected model. These configs keep architecture and core
recipe fixed while changing train-time augmentation.

| Run | Config | Change | Status |
| --- | --- | --- | --- |
| Selected baseline | `configs/selected_mvp_classifier.json` | Current fixed selected model | reference |
| Mild context augmentation | `configs/selected_mvp_aug_context_mild.json` | Crop scale `0.65-1.0`, rotation `12`, color jitter `0.2/0.2/0.15` | reject |
| Strong context augmentation | `configs/selected_mvp_aug_context_strong.json` | Crop scale `0.55-1.0`, rotation `15`, color jitter `0.25/0.25/0.2` | select |
| Mild context + random erasing | `configs/selected_mvp_aug_random_erasing.json` | Mild context settings plus `random_erasing_p=0.1` | reject |

Batch result:

| Run | Top-1 | Top-3 | Rendang F1 | Gado-gado F1 | Soto F1 | Weak avg F1 | Priority confusions | Misclassified | Decision |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| Selected baseline | 91.87% | 98.87% | 83.95% | 87.93% | 89.37% | 87.08% | 5 | 36 | previous selected recipe |
| Mild context augmentation | 91.87% | 98.19% | 86.08% | 87.62% | 89.66% | 87.78% | 3 | 36 | reject: no top-1 or error-count gain |
| Strong context augmentation | 93.91% | 99.10% | 90.48% | 90.38% | 89.83% | 90.23% | 4 | 27 | select: best global and weak-class result |
| Mild context + random erasing | 91.87% | 98.87% | 88.10% | 89.32% | 88.70% | 88.70% | 3 | 36 | reject: weak-class gain without error-count gain |

The selected training recipe remains `convnext_tiny.fb_in1k`,
`image_size=256`, `learning_rate=0.0001`, and active `label_smoothing=0.1`,
but should adopt the strong context augmentation settings:

- `random_resized_crop_scale=[0.55, 1.0]`
- `rotation_degrees=15`
- `color_jitter_brightness=0.25`
- `color_jitter_contrast=0.25`
- `color_jitter_saturation=0.2`
- `random_erasing_p=0.0`

This is not evidence for a broader augmentation sweep yet. Strong context
augmentation improved top-1 by `+2.03pp`, top-3 by `+0.23pp`, weak-class
average F1 by `+3.15pp`, and reduced total misclassified images from `36` to
`27`. The next work should promote this recipe, validate serving/runtime
compatibility, then run focused error analysis on the new `27` failures.

## Final Selected Model Error Analysis

Reviewed file:

```txt
reports/selected-mvp-classifier/error_analysis.xlsx
```

Parsed summary:

| Category | Count |
| --- | ---: |
| Total reviewed misclassified images | 27 |
| `model_error` | 25 |
| `out_of_scope` | 1 |
| `ambiguous` | 1 |
| `keep_for_tuning_signal` | 25 |
| `exclude_in_next_dataset` | 2 |

Error type breakdown:

| Error Type | Count |
| --- | ---: |
| `background_or_context_bias` | 13 |
| `weak_class_overlap` | 7 |
| `image_quality` | 4 |
| `low_visual_dominance` | 2 |
| `mixed_food` | 1 |

Largest repeated confusion:

- `bakso` as `soto`: 3, mixed causes across context, visual dominance, and
  weak-class overlap
- `pempek` as `gado_gado`: 3, two valid tuning signals and one out-of-scope
  exclusion candidate
- `sate` as `rendang`: 3, all `weak_class_overlap`
- `gado_gado` as `soto`: 2, both `background_or_context_bias`
- `soto` as `gudeg`: 2, both `background_or_context_bias`

Conclusion: the final selected model's remaining failures are mostly valid
model errors. Only two images should be excluded in a future dataset version, so
do not run a broad dataset cleanup before MVP integration. The next engineering
step is to wire and validate this artifact in the inference service. Further
model improvement, if needed after serving validation, should target
`background_or_context_bias` and `weak_class_overlap` with more representative
examples or a small, validation-driven tuning batch.

Selection rule for this batch:

1. Do not accept a run that drops top-1 materially below the selected model's
   `91.87%`.
2. Prefer runs that reduce `background_or_context_bias` misclassifications after
   review.
3. Prefer runs that preserve or improve weak-class average F1 above `87.08%`.
4. If metrics are close, keep the selected baseline rather than adding stronger
   augmentation.

## Guardrails

Do not run a broad tuning sweep directly from the selected model. The held-out
test set has already been used to choose architecture and the first tuned
ConvNeXt-Tiny config, so repeated selection against the same test set risks
overfitting the experiment process.

Before another tuning batch:

1. Define the next objective before training.
2. Tune one axis at a time.
3. Prefer validation-driven tuning and reserve the test set for final
   confirmation of a small number of candidates.
4. Keep generated `reports/`, `model-artifacts/`, dataset images, and ZIP files
   out of Git.

## Colab Command

Run from `services/ai-inference` in Colab:

```bash
CONFIG=configs/selected_mvp_aug_context_mild.json \
PROCESSED_DIR=data/processed-v0.2 \
REQUIRE_CUDA=1 \
INSTALL_DEPS=1 \
bash scripts/colab_retrain_baseline_v2.sh
```

Change `CONFIG` for the next candidate.
