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

**Estimated Energy**:
An approximate calorie value inferred from food category and visual cues.
_Avoid_: clinical nutrition value, exact calories

**Inference Payload**:
The structured result returned by AI/ML Inference to the Backend API.
_Avoid_: AI recommendation, nudge

## Relationships

- An **Inference Payload** contains one or more **Food Categories**
- A **Food Category** has a **Confidence Score**
- A **Food Category** may have **Visual Dominance**
- **Estimated Energy** is approximate and belongs to the inference result, not a medical diagnosis

## Example Dialogue

> **Dev:** "Can AI/ML return 'sisihkan 1/4 porsi' directly?"
> **Domain expert:** "No — AI/ML returns an **Inference Payload**; the Backend API turns that into a nudge decision."

## Flagged Ambiguities

- "portion size" was used for image-based estimates — resolved: AI/ML reports **Visual Dominance**, not exact serving size.
- "calories" was used as if exact — resolved: NutriScan uses **Estimated Energy** for approximate preventive feedback.
