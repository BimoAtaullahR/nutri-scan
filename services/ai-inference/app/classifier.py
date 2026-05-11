from __future__ import annotations

import json
from pathlib import Path

from app.payload import Prediction


DEFAULT_ARTIFACT_DIR = Path("model-artifacts/baseline-food-classifier")


class FoodClassifier:
    def __init__(self, artifact_dir: Path = DEFAULT_ARTIFACT_DIR) -> None:
        self.artifact_dir = artifact_dir
        self.model_path = artifact_dir / "model.pt"
        self.label_map_path = artifact_dir / "label_map.json"
        self.model_version = "baseline-0.1.0"
        self.labels = self._load_labels()

    def _load_labels(self) -> list[str]:
        if self.label_map_path.exists():
            raw = json.loads(self.label_map_path.read_text())
            return [raw["idToLabel"][str(index)] for index in range(len(raw["idToLabel"]))]
        return ["sate", "rendang", "bakso"]

    def predict(self, image_bytes: bytes) -> list[Prediction]:
        if not image_bytes:
            raise ValueError("Image upload cannot be empty")

        if not self.model_path.exists():
            return [
                {"label": self.labels[0], "confidenceScore": 0.61},
                {"label": self.labels[1], "confidenceScore": 0.24},
                {"label": self.labels[2], "confidenceScore": 0.15},
            ]

        raise NotImplementedError("Model artifact inference wiring is planned after training export")


_classifier = FoodClassifier()


def get_classifier() -> FoodClassifier:
    return _classifier
