# AI/ML Inference

AI/ML Inference owns image preprocessing, food recognition, visual dominance detection, confidence scoring, and estimated energy payloads for NutriScan.

## Language

**Food Category**:
A recognized class of food that the model believes appears in a submitted image.
_Avoid_: menu item, dish label

**Food Item**:
A visible food unit in a scan that can reasonably receive its own category, confidence, and energy estimate.
_Avoid_: ingredient, hidden recipe component

**Meal-Level Food Category**:
A recognized whole-dish category used when visible food items cannot be separated reliably.
_Avoid_: per-lauk result, ingredient breakdown

**Single Primary Food Scan**:
An MVP scan result that recognizes one primary food category for the submitted image.
_Avoid_: per-lauk scan, multi-item analysis

**Food Item Evidence**:
Visual support for a detected food item, such as an item list entry, bounding region, or segmentation mask.
_Avoid_: proof of exact serving size, nutrition measurement

**External Vision Candidate**:
A structured candidate food item or portion produced by an external vision-language model before NutriScan normalization.
_Avoid_: final inference result, calorie authority

**External Vision Fallback**:
An optional fallback path used when the local MVP classifier cannot confidently recognize a food category.
_Avoid_: primary classifier, calorie estimator

**Fallback Label Assistant**:
An external vision fallback constrained to suggest normalized MVP food categories or unknown food.
_Avoid_: second authority, free-form food classifier

**Visual Dominance**:
The relative visual prominence of a detected food category within the image.
_Avoid_: portion size, exact serving size

**Confidence Score**:
The model-reported certainty for a recognized **Food Category**.
_Avoid_: accuracy, correctness

**Unknown Food**:
A scan result used when confidence is too low to trust a food category prediction.
_Avoid_: failed scan, other

**Low-Confidence Fallback**:
A fallback path considered when the local classifier confidence is below the normal-result threshold.
_Avoid_: failed inference, always-on external scan

**Dataset Manifest**:
A record of dataset sources, licenses, class coverage, and train/validation/test counts used for model development.
_Avoid_: loose dataset notes

**Dataset Drop Folder**:
A shared storage location for versioned dataset archives and model artifacts, not a place for direct training work.
_Avoid_: shared working dataset, Git-tracked dataset

**Dataset Owner**:
The teammate responsible for publishing official dataset versions to the dataset drop folder.
_Avoid_: sole model owner, Git approver

**Portion Ground Truth Subset**:
A smaller evaluation set with measured portion or weight information for checking energy estimate quality.
_Avoid_: required label for every dataset image, full nutrition dataset

**Estimated Energy**:
An approximate calorie value inferred from food category and visual cues.
_Avoid_: clinical nutrition value, exact calories

**Estimated Energy Range**:
A bounded calorie estimate for a recognized food category and coarse portion level.
_Avoid_: calorie detection, exact calories

**Per-Food Estimated Energy**:
An approximate calorie range for one detected food item within a scan.
_Avoid_: exact calories per lauk, clinical per-item nutrition

**Review-Needed Energy Estimate**:
A per-food energy estimate whose range is too broad or uncertain to support a clear nudge without user correction.
_Avoid_: unusable estimate, failed estimate

**Energy Lookup Table**:
A maintained table that maps MVP food categories and coarse portion levels to calorie ranges.
_Avoid_: calorie regression model

**Curated Energy Source**:
A documented nutrition reference used to maintain energy ranges for MVP food categories.
_Avoid_: model-guessed calories, untracked nutrition numbers

**Coarse Portion Estimate**:
A small, medium, or large portion estimate used to select an energy range.
_Avoid_: exact serving size, gram estimate

**Default Coarse Portion**:
The initial medium portion value used when the MVP does not have reliable portion evidence.
_Avoid_: model-measured serving, inferred grams

**User-Corrected Portion**:
A user-confirmed coarse portion value that overrides or confirms the model's portion estimate.
_Avoid_: exact serving size, nutrition measurement

**Inference Payload**:
The structured result returned by AI/ML Inference to the Backend API.
_Avoid_: AI recommendation, nudge

**Recognizer Payload**:
An MVP inference payload limited to food category, confidence score, and alternatives.
_Avoid_: calorie payload, nudge payload, segmentation result

## Relationships

- An **Inference Payload** contains one or more **Food Categories**
- A **Recognizer Payload** is an **Inference Payload**
- A **Single Primary Food Scan** contains one primary **Food Category**
- A **Food Category** may describe a **Food Item** or a **Meal-Level Food Category**
- A **Food Item** may have **Food Item Evidence**
- An **External Vision Candidate** may become a **Food Item** only after NutriScan normalization
- An **External Vision Fallback** may produce **External Vision Candidates**
- A **Fallback Label Assistant** is an **External Vision Fallback**
- A **Food Category** has a **Confidence Score**
- A **Low-Confidence Fallback** may use a **Fallback Label Assistant**
- A **Food Category** may have **Visual Dominance**
- **Estimated Energy** is approximate and belongs to the inference result, not a medical diagnosis
- **Per-Food Estimated Energy** belongs to one detected **Food Category** within an **Inference Payload**
- A **Review-Needed Energy Estimate** should lead the Backend API toward a review-oriented nudge rather than a confident portion-reduction nudge
- An **Estimated Energy Range** is derived from a **Food Category** and coarse portion estimate
- An **Energy Lookup Table** provides MVP **Estimated Energy Range** values
- An **Energy Lookup Table** should be backed by one or more **Curated Energy Sources**
- A **Coarse Portion Estimate** selects the relevant **Estimated Energy Range**
- A **Default Coarse Portion** may be used until the user provides a **User-Corrected Portion**
- A **User-Corrected Portion** may replace a **Coarse Portion Estimate** for product feedback
- A **Portion Ground Truth Subset** is used to evaluate **Per-Food Estimated Energy** quality
- **Unknown Food** is returned when the **Confidence Score** is below the MVP threshold
- A **Dataset Drop Folder** distributes dataset versions while local machines run training work
- A **Dataset Owner** publishes official dataset versions to the **Dataset Drop Folder**

## MVP Food Categories

- nasi goreng
- sate
- rendang
- bakso
- gado-gado
- soto
- pempek
- gudeg

## Example Dialogue

> **Dev:** "Can AI/ML return 'sisihkan 1/4 porsi' directly?"
> **Domain expert:** "No — AI/ML returns an **Inference Payload**; the Backend API turns that into a nudge decision."

## Flagged Ambiguities

- "portion size" was used for image-based estimates — resolved: AI/ML reports **Visual Dominance**, not exact serving size.
- "calories" was used as if exact — resolved: NutriScan uses **Estimated Energy** for approximate preventive feedback.
- "calorie detection" was used as if exact — resolved: MVP uses **Estimated Energy Range**, not exact calorie detection.
- "nasi padang" can mean many separate lauk — resolved: it is deferred to future development rather than included in MVP.
- "calorie regression model" is future research — resolved: MVP uses an **Energy Lookup Table** until measured calorie ground truth is available.
- "portion detection" was used as if automatic and exact — resolved: MVP uses **Coarse Portion Estimate** with user correction.
- "other" was used for low confidence predictions — resolved: MVP returns **Unknown Food** so the app can ask for user confirmation or correction.
- "calories per lauk" was used as if each item could be measured exactly from an image — resolved: NutriScan reports **Per-Food Estimated Energy** as a range with confidence and user correction.
- "useful estimate" was vague — resolved: normal estimates should target roughly ±20–30%, while broader or uncertain per-food ranges become **Review-Needed Energy Estimates**.
- "lauk" was ambiguous between visible food item, whole dish, and hidden ingredient — resolved: NutriScan targets **Food Item** detection, with **Meal-Level Food Category** fallback for mixed foods that cannot be separated reliably.
- "segmentation" was considered mandatory for the first per-food version — resolved: **Food Item Evidence** may start as item-level or region-level evidence, while segmentation remains the target for stronger portion estimates.
- "Gemini result" was ambiguous between candidate and product result — resolved: external vision models produce **External Vision Candidates** that must be normalized before use.
- "use Gemini" was ambiguous between primary inference and fallback — resolved: external vision is an **External Vision Fallback** only for low-confidence local results.
- "Gemini integration" risked becoming a second food authority — resolved: Gemini acts only as a **Fallback Label Assistant** constrained to MVP categories or unknown food.
- "when to call Gemini" was ambiguous — resolved: Gemini is considered only through **Low-Confidence Fallback**, not for normal-confidence local results.
- "calorie source" was ambiguous — resolved: user-facing energy ranges should come from **Curated Energy Sources**, not free-form model guesses.
- "per-lauk MVP" exceeded the available delivery window — resolved: MVP uses a **Single Primary Food Scan**, while per-food detection remains future research.
- "add chicken categories" would improve Indonesian food coverage but expands dataset work — resolved: MVP uses the eight locally curated food categories, while ayam goreng and ayam bakar are future categories.
- "small/medium/large from AI" was too strong for available evidence — resolved: MVP starts with a **Default Coarse Portion** and lets the user correct it.
- "AI contract" was too broad for the shortened MVP — resolved: MVP uses a **Recognizer Payload** and Backend API adds portion, energy range, and nudge decisions.
- "shared dataset folder" was ambiguous between a working directory and an archive exchange — resolved: the team uses a **Dataset Drop Folder** for versioned archives, not direct training work.
- "who updates the dataset" was ambiguous — resolved: teammates may propose changes, but a **Dataset Owner** publishes official dataset versions.

## Future Research

- Train or fine-tune a calorie regression model if NutriScan later has food images paired with reliable weight, portion, or calorie ground truth.
- Add nasi padang as a meal-level category after MVP if enough complete-meal images are available.
- Add per-food detection or segmentation after the MVP can reliably classify a primary food category.
- Add ayam goreng and ayam bakar after MVP if curated image coverage and energy ranges are available.

## MVP Model Strategy

- Fine-tune a pretrained image classifier for the MVP food categories.
- Do not train a model from scratch for the MVP.
- Use a **Single Primary Food Scan** for the MVP.
- Return a **Recognizer Payload** without calories, nudges, or segmentation masks.
- Use a **Low-Confidence Fallback** with a **Fallback Label Assistant** when local confidence is below the normal-result threshold.
- Use FastAPI for the inference HTTP service.
- Use PyTorch and timm for classifier training and inference.
- Prefer lightweight pretrained baselines such as EfficientNet-B0 or MobileNetV3.
- Use ConvNeXt-Tiny with the selected strong context augmentation recipe as the
  selected MVP classifier after the model comparison and tuning screens, pending
  runtime artifact wiring and serving validation.
- Use `nutriscan-mvp-food-dataset-v0.2` / `data/processed-v0.2` as the active
  reviewed dataset for EfficientNet-B0 baseline v2 retraining.
- Apply misclassified review decisions through `scripts/apply_misclassified_review.py`;
  do not edit the source processed dataset in place.
- Evaluate with top-1 accuracy, top-3 accuracy, and a confusion matrix on a held-out test set.
- Evaluate per-food research with item match rate, top-3 coverage, coarse portion accuracy, energy range quality, review-needed rate, direct nudge rate, and user correction rate.
- Store MVP evaluation outputs as local metrics files, not in an experiment tracking service.
- Maintain a **Dataset Manifest** for source, license, class coverage, and split counts.
- Deliver an end-to-end inference slice: dataset preparation, baseline training, evaluation, exported model artifact, and an inference API usable by the Backend API.
