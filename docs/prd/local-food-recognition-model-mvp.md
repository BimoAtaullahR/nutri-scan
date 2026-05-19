# Local Food Recognition Model MVP

## Summary

Build a local AI/ML inference slice that recognizes one primary food category from a scan image. The MVP does not estimate calories from pixels, detect per-food items, segment lauk, or estimate grams. Energy feedback is derived by Backend API from a lookup table using the recognized food category and a coarse portion.

## Scope

In scope:

- Local image classifier for eight curated Indonesian food categories.
- Recognizer payload with food category, confidence score, alternatives, and low-confidence status.
- Lookup-based energy estimate owned outside the model.
- Default medium portion with user correction to small, medium, or large.
- External vision fallback for low-confidence local predictions.

Out of scope:

- Per-lauk detection.
- Segmentation masks.
- Gram or volume estimation.
- Calorie regression.
- New categories such as ayam goreng and ayam bakar.

## MVP Categories

- `nasi_goreng`
- `sate`
- `rendang`
- `bakso`
- `gado_gado`
- `soto`
- `pempek`
- `gudeg`

## Development Phases

### 1. Dataset Scope Lock

Confirm that the training dataset only contains the eight MVP categories. Remove extra processed class folders such as `bebek_betutu` and `rawon` from `data/processed` so the training script matches `configs/mvp_food_categories.json`.

Output:

- `data/processed` contains only MVP class folders.
- `configs/mvp_food_categories.json` remains the source of truth for class labels.

### 2. Dataset Audit

Audit train, validation, and test counts per class. Spot-check images for bad labels, corrupt files, watermarks, menus, drawings, and obvious duplicates. Update the dataset manifest with the final reviewed counts and known risks.

Output:

- Updated dataset count summary.
- Known dataset risks documented.

### 3. Baseline Training

Train a lightweight pretrained classifier such as EfficientNet-B0 or MobileNetV3 using the local processed dataset. Start with a short training run before tuning.

Output:

- Model artifact.
- Label map.
- Training configuration record.

### 4. Evaluation

Evaluate the classifier on the held-out test split. Report top-1 accuracy, top-3 accuracy, confusion matrix, and per-class performance. Use the result to identify classes that are often confused.

Output:

- Metrics report.
- Confusion matrix.
- Notes on weak classes.

### 5. Confidence Threshold Calibration

Choose thresholds for normal, low-confidence, and unknown results. Initial target:

- `>= 0.75`: normal local result.
- `0.60-0.75`: low-confidence result.
- `< 0.60`: unknown food.

Output:

- Final confidence thresholds for MVP.
- Clear behavior for low-confidence scans.

### 6. Inference API

Serve the trained model through AI/ML Inference. The endpoint should return a recognizer payload only, without calories, portions, nudges, or segmentation data.

Example payload:

```json
{
  "modelVersion": "baseline-efficientnet-b0-v1",
  "foodCategory": {
    "slug": "nasi_goreng",
    "confidenceScore": 0.82
  },
  "alternatives": [
    {
      "slug": "sate",
      "confidenceScore": 0.09
    }
  ],
  "isLowConfidence": false,
  "confidenceThreshold": 0.75
}
```

Output:

- `/infer` endpoint compatible with Backend API.
- Stable response shape for mobile and backend integration.

### 7. External Vision Fallback

Call Gemini or another external vision model only when local confidence is below the normal-result threshold. The fallback must be constrained to the eight MVP categories plus `unknown_food`. It must not return user-facing calories, portions, or nudges.

Output:

- Fallback label candidate for low-confidence scans.
- Timeout and failure handling.
- Normalization from external candidate to MVP category or `unknown_food`.

### 8. Backend Integration

Backend API receives the recognizer payload, assigns the default medium portion, looks up the estimated energy range, and creates the nudge decision. The model remains a recognizer; Backend API owns energy and nudge behavior.

Output:

- End-to-end scan result with category, estimated energy range, and nudge.
- Clear separation between AI recognition and backend product rules.

### 9. Result Correction

Allow the user to correct the effective food category, coarse portion, or both. Backend API recalculates the energy range and nudge decision after correction.

Output:

- Result correction flow for category and portion.
- Recomputed effective scan result.

### 10. End-to-End QA

Test the scan flow with real food photos. Verify normal-confidence behavior, low-confidence behavior, external fallback, energy lookup, result correction, and nudge output.

Output:

- Demo-ready local food recognition flow.
- Known limitations documented before presentation.

## Acceptance Criteria

- The local classifier supports the eight MVP categories.
- The model returns top-1 and top-3 predictions with confidence scores.
- Low-confidence scans are handled without treating them as technical failures.
- Gemini fallback is used only below the confidence threshold.
- Backend energy estimates come from the lookup table, not model-generated calories.
- Users can correct portion and category when needed.

## Future Scope

- Add ayam goreng and ayam bakar after curated data and energy ranges are available.
- Add per-food detection after the primary-food classifier is reliable.
- Add segmentation after detection is valuable enough for the product.
- Explore calorie regression only after reliable portion or weight ground truth exists.
