# Use Hybrid Local-First Food Inference

NutriScan will treat local AI/ML Inference as the primary source for product-facing food categories, confidence scores, and estimated energy ranges, while allowing external vision-language models to assist with low-confidence scans, ambiguous mixed foods, and dataset labeling. This keeps scan results evaluable and stable while still using external models to accelerate research and handle cases where the local model is not yet strong.

## Considered Options

- Use an external vision-language model as the primary scan authority.
- Use only local custom models.
- Use a hybrid local-first approach with external models as fallback or labeling assistants.

## Consequences

- External model output must be normalized into NutriScan's inference contract before the Backend API uses it.
- External models should not be treated as the direct authority for exact per-food calories.
- The AI/ML service remains responsible for measurable confidence, estimated energy ranges, and review-needed signals.
- The first per-food prototype may use external model candidates plus user review as a data collection loop before local detection or segmentation models are strong enough.
