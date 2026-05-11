from __future__ import annotations

from typing import TypedDict

from app.energy import EnergyRange, lookup_energy_range


class Prediction(TypedDict):
    label: str
    confidenceScore: float


def _prediction_item(prediction: Prediction) -> dict[str, object]:
    return {
        "slug": prediction["label"],
        "confidenceScore": float(prediction["confidenceScore"]),
    }


def build_inference_payload(
    predictions: list[Prediction],
    portion: str,
    model_version: str,
    confidence_threshold: float = 0.6,
) -> dict[str, object]:
    if not predictions:
        raise ValueError("At least one prediction is required")

    ranked = sorted(predictions, key=lambda item: item["confidenceScore"], reverse=True)
    top_prediction = ranked[0]
    is_low_confidence = top_prediction["confidenceScore"] < confidence_threshold
    alternatives_source = ranked[:3] if is_low_confidence else ranked[1:3]
    alternatives = [_prediction_item(prediction) for prediction in alternatives_source]

    energy_range: EnergyRange | None = None
    food_category = _prediction_item(top_prediction)
    if is_low_confidence:
        food_category = {"slug": "unknown_food", "confidenceScore": top_prediction["confidenceScore"]}
    else:
        energy_range = lookup_energy_range(top_prediction["label"], portion)

    return {
        "modelVersion": model_version,
        "foodCategory": food_category,
        "alternatives": alternatives,
        "coarsePortion": portion,
        "estimatedEnergyRange": energy_range,
        "isLowConfidence": is_low_confidence,
        "confidenceThreshold": confidence_threshold,
    }
