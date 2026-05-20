from app.payload import build_inference_payload


def test_low_confidence_prediction_returns_unknown_food() -> None:
    payload = build_inference_payload(
        predictions=[
            {"label": "sate", "confidenceScore": 0.42},
            {"label": "rendang", "confidenceScore": 0.31},
            {"label": "bakso", "confidenceScore": 0.17},
        ],
        portion="medium",
        model_version="selected-mvp-classifier",
        confidence_threshold=0.7,
    )

    assert payload["modelVersion"] == "selected-mvp-classifier"
    assert payload["foodCategory"]["slug"] == "unknown_food"
    assert payload["isLowConfidence"] is True
    assert payload["estimatedEnergyRange"] is None
    assert payload["confidenceThreshold"] == 0.7
    assert [item["slug"] for item in payload["alternatives"]] == ["sate", "rendang", "bakso"]


def test_confident_prediction_includes_portion_energy_and_alternatives() -> None:
    payload = build_inference_payload(
        predictions=[
            {"label": "sate", "confidenceScore": 0.81},
            {"label": "rendang", "confidenceScore": 0.11},
            {"label": "bakso", "confidenceScore": 0.08},
        ],
        portion="small",
        model_version="baseline-0.1.0",
        confidence_threshold=0.6,
    )

    assert payload["modelVersion"] == "baseline-0.1.0"
    assert payload["foodCategory"] == {"slug": "sate", "confidenceScore": 0.81}
    assert payload["coarsePortion"] == "small"
    assert payload["estimatedEnergyRange"]["minKcal"] > 0
    assert [item["slug"] for item in payload["alternatives"]] == ["rendang", "bakso"]
