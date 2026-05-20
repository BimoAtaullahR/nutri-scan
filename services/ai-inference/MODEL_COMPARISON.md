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

Status: ConvNeXt-Tiny was selected for MVP after the tuning screen. The fixed
project model config is `configs/selected_mvp_classifier.json`, based on
`convnext_tiny.fb_in1k` with `image_size=256`, `learning_rate=0.0001`, and active
`label_smoothing=0.1`. Runtime artifact wiring and serving smoke tests remain
separate follow-up work.

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

Important: earlier architecture-screen runs were produced before
`label_smoothing` was wired into the training loss. Treat those results as
pre-fix references. Before comparing ConvNeXt-Tiny tuning runs, run the
ConvNeXt-Tiny control with the patched trainer and a separate output folder so
all tuning results share the same active `label_smoothing` behavior without
overwriting the pre-fix reference.

Do not add `label_smoothing=0.0` to the first tuning batch. The intended recipe
uses active `label_smoothing=0.1`; keep that fixed for the first batch and only
make label smoothing a tuning axis if the first batch gives a specific reason.

| Run | Config | Change | Decision |
| --- | --- | --- | --- |
| Control rerun | `configs/convnext_tiny_control_label_smoothing.json` | ConvNeXt-Tiny with active `label_smoothing` | pending |
| Tune 1 | `configs/convnext_tiny_tune_lr5e5.json` | learning rate `0.0001` -> `0.00005` | pending |
| Tune 2 | `configs/convnext_tiny_tune_img256.json` | image size `224` -> `256` | pending |
| Tune 3 | `configs/convnext_tiny_tune_lr5e5_img256.json` | learning rate `0.00005`, image size `256` | pending |

Tuning results:

| Run | Top-1 | Top-3 | Rendang F1 | Gado-gado F1 | Soto F1 | Weak-class Avg F1 | Priority Confusions | Total Misclassified | Decision |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: | --- |
| Control rerun | 90.97% | 98.65% | 83.95% | 91.67% | 86.67% | 87.43% | 6 | 40 | baseline for tuning |
| Tune 1: LR `0.00005` | 91.20% | 98.19% | 85.37% | 88.66% | 89.26% | 87.76% | 3 | 39 | not selected |
| Tune 2: image size `256` | 91.87% | 98.87% | 83.33% | 88.66% | 89.26% | 87.08% | 5 | 36 | selected |
| Tune 3: LR `0.00005`, image size `256` | 91.20% | 98.65% | 81.48% | 90.91% | 87.80% | 86.73% | 2 | 39 | not selected |

Tune 2 is selected because it has the highest top-1 accuracy and top-3 accuracy,
and the lowest total misclassified count. Its weak-class average F1 is lower
than Tune 1 by less than one percentage point, but the selection rule prioritizes
top-1 unless candidates are within `0.5 pp`.

Select the tuned ConvNeXt-Tiny candidate using this order:

1. Highest top-1 accuracy.
2. If top-1 differs by less than `0.5 pp`, choose the highest weak-class average
   F1.
3. If still close, choose the lower combined confusion count for `soto` as
   `bakso`, `rendang` as `nasi_goreng`, and `gado_gado` as `soto`.
4. If still tied, choose the cheaper config: `224` image size before `256`, and
   `learning_rate=0.0001` before `0.00005`.

## Further Tuning Guardrails

Do not start another large tuning sweep directly from the selected model. The
current held-out test set has already been used to choose the architecture and
the first tuned ConvNeXt-Tiny config. Repeatedly selecting runs from the same
test set risks optimizing the experiment process around that test set rather
than improving generalization.

Before another tuning batch:

1. Freeze `configs/selected_mvp_classifier.json` as the current selected-model
   checkpoint.
2. Review the selected model's misclassified images and summarize whether errors
   are caused by ambiguous images, label issues, bad quality, weak-class visual
   overlap, or likely model capacity/recipe limitations.
3. Define the next tuning objective before running experiments.
4. Prefer validation-driven tuning and reserve the test set for final
   confirmation of a small number of candidates.

Recommended next objective:

- preserve top-1 accuracy around the selected model's `91.87%`
- improve weak-class average F1 above `87.08%`
- reduce priority confusions, especially `soto` as `bakso`, `rendang` as
  `nasi_goreng`, and `gado_gado` as `soto`
- avoid selecting a model for a tiny top-1 gain if weak-class behavior gets worse

Tuning axes to consider after error analysis:

| Axis | Candidate Values | Notes |
| --- | --- | --- |
| `weight_decay` | `0.0003`, `0.0005`, `0.001` | Regularization sweep around current value |
| `label_smoothing` | `0.0`, `0.05`, `0.1` | Now active in the trainer; compare only after deciding it is worth an axis |
| `learning_rate` | `0.0001`, `0.000075`, `0.00005` | Use smaller changes around the selected config |
| augmentation strength | crop scale, rotation, color jitter | Requires trainer support before config-only tuning |
| training budget | epochs, early stopping patience | Use only if validation curves suggest undertraining |
| weak-class handling | class weighting or weighted sampling | Requires trainer support and careful validation |
| test-time augmentation | flip/crop averaging | Runtime trade-off, not a training improvement |

Do not add all axes at once. Add one axis only when the selected model's error
analysis gives a reason for it.

## Selected Model Error Analysis

Reviewed file:

```txt
reports/convnext-tiny-tune-img256/error_analysis.xlsx
```

Summary after second-pass review:

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

Largest repeated confusion:

- `sate` as `rendang`: 5, all `weak_class_overlap`
- `bakso` as `soto`: 3, mostly `background_or_context_bias`
- `gado_gado` as `soto`: 3, all `background_or_context_bias`
- `rendang` errors are spread across several predicted classes and are often
  `background_or_context_bias` or `low_visual_dominance`

Conclusion: the selected model's remaining errors are mostly valid model errors,
not broad dataset corruption. Only three reviewed images are candidates for
exclusion in a future dataset version. The next improvement path should focus on
context robustness and class-overlap handling rather than a large dataset cleanup
or broad hyperparameter sweep.

## Context Robustness Tuning Batch

The next tuning batch targets `background_or_context_bias`, which is the largest
reviewed error type for the selected model. These configs keep the selected model
architecture and core training recipe fixed while changing only train-time
augmentation.

| Run | Config | Change | Decision |
| --- | --- | --- | --- |
| Selected baseline | `configs/selected_mvp_classifier.json` | Current fixed selected model | reference |
| Mild context augmentation | `configs/selected_mvp_aug_context_mild.json` | Crop scale `0.65-1.0`, rotation `12`, color jitter `0.2/0.2/0.15` | pending |
| Strong context augmentation | `configs/selected_mvp_aug_context_strong.json` | Crop scale `0.55-1.0`, rotation `15`, color jitter `0.25/0.25/0.2` | pending |
| Mild context + random erasing | `configs/selected_mvp_aug_random_erasing.json` | Mild context settings plus `random_erasing_p=0.1` | pending |

Selection rule for this batch:

1. Do not accept a run that drops top-1 materially below the selected model's
   `91.87%`.
2. Prefer runs that reduce `background_or_context_bias` misclassifications after
   review.
3. Prefer runs that preserve or improve weak-class average F1 above `87.08%`.
4. If metrics are close, keep the selected baseline rather than adding stronger
   augmentation.

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
