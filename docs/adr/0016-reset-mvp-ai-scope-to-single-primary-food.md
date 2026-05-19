# Reset MVP AI Scope to Single Primary Food

NutriScan will reset the MVP AI scope to single primary food classification instead of per-food detection, segmentation, or calorie regression.

The MVP classifier will recognize the eight locally curated Indonesian food categories already present in the repository: nasi goreng, sate, rendang, bakso, gado-gado, soto, pempek, and gudeg. AI/ML Inference returns the food category, confidence score, and alternatives. Backend API derives estimated energy ranges from an energy lookup table using the recognized food category and a coarse portion.

The MVP will use medium as the default coarse portion. Users may correct the portion to small, medium, or large. External vision models may assist only when local confidence is low, and their output must be normalized into NutriScan categories or treated as unknown food. External vision models do not provide user-facing calories.

## Considered Options

- Continue with per-food detection or segmentation for the MVP.
- Add more categories, including ayam goreng and ayam bakar, before launch.
- Use external vision-language models as the primary food recognizer.
- Reset to single primary food classification with eight curated categories and lookup-based energy ranges.

## Consequences

- The MVP is deliverable within the shortened development window.
- Per-food detection, segmentation, item-targeted nudges, and calorie regression move to future research.
- Backend nudge decisions can still provide useful portion guidance from food category, coarse portion, and estimated energy range.
- Additional categories require curated image coverage, energy ranges, and evaluation before they become product scope.
- Existing per-food research documentation remains useful as future direction, but it is not the MVP implementation target.
