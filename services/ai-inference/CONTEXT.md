# AI/ML Inference

AI/ML Inference owns image preprocessing, food recognition, visual dominance detection, confidence scoring, and estimated energy payloads for NutriScan.

## Language

**Food Category**:
A recognized class of food that the model believes appears in a submitted image.
_Avoid_: menu item, dish label

**Visual Dominance**:
The relative visual prominence of a detected food category within the image.
_Avoid_: portion size, exact serving size

**Confidence Score**:
The model-reported certainty for a recognized **Food Category**.
_Avoid_: accuracy, correctness

**Unknown Food**:
A scan result used when confidence is too low to trust a food category prediction.
_Avoid_: failed scan, other

**Dataset Manifest**:
A record of dataset sources, licenses, class coverage, and train/validation/test counts used for model development.
_Avoid_: loose dataset notes

**Estimated Energy**:
An approximate calorie value inferred from food category and visual cues.
_Avoid_: clinical nutrition value, exact calories

**Estimated Energy Range**:
A bounded calorie estimate for a recognized food category and coarse portion level.
_Avoid_: calorie detection, exact calories

**Energy Lookup Table**:
A maintained table that maps MVP food categories and coarse portion levels to calorie ranges.
_Avoid_: calorie regression model

**Coarse Portion Estimate**:
A small, medium, or large portion estimate used to select an energy range.
_Avoid_: exact serving size, gram estimate

**Inference Payload**:
The structured result returned by AI/ML Inference to the Backend API.
_Avoid_: AI recommendation, nudge

## Relationships

- An **Inference Payload** contains one or more **Food Categories**
- A **Food Category** has a **Confidence Score**
- A **Food Category** may have **Visual Dominance**
- **Estimated Energy** is approximate and belongs to the inference result, not a medical diagnosis
- An **Estimated Energy Range** is derived from a **Food Category** and coarse portion estimate
- An **Energy Lookup Table** provides MVP **Estimated Energy Range** values
- A **Coarse Portion Estimate** selects the relevant **Estimated Energy Range**
- **Unknown Food** is returned when the **Confidence Score** is below the MVP threshold

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

## Future Research

- Train or fine-tune a calorie regression model if NutriScan later has food images paired with reliable weight, portion, or calorie ground truth.
- Add nasi padang as a meal-level category after MVP if enough complete-meal images are available.

## MVP Model Strategy

- Fine-tune a pretrained image classifier for the MVP food categories.
- Do not train a model from scratch for the MVP.
- Use FastAPI for the inference HTTP service.
- Use PyTorch and timm for classifier training and inference.
- Prefer lightweight pretrained baselines such as EfficientNet-B0 or MobileNetV3.
- Evaluate with top-1 accuracy, top-3 accuracy, and a confusion matrix on a held-out test set.
- Store MVP evaluation outputs as local metrics files, not in an experiment tracking service.
- Maintain a **Dataset Manifest** for source, license, class coverage, and split counts.
- Deliver an end-to-end inference slice: dataset preparation, baseline training, evaluation, exported model artifact, and an inference API usable by the Backend API.
