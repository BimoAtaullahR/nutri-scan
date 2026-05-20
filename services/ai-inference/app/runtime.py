from __future__ import annotations

import json
import os
from dataclasses import dataclass
from pathlib import Path
from typing import Any

DEFAULT_ARTIFACT_DIR = Path("model-artifacts/selected-mvp-classifier")
DEFAULT_MODEL_VERSION = "selected-mvp-classifier"
DEFAULT_CONFIDENCE_THRESHOLD = 0.6
EXPECTED_MVP_LABELS = {
    "bakso",
    "gado_gado",
    "gudeg",
    "nasi_goreng",
    "pempek",
    "rendang",
    "sate",
    "soto",
}
REQUIRED_ARTIFACT_FILES = ("model.pt", "label_map.json", "training_config_resolved.json")


class ArtifactValidationError(RuntimeError):
    pass


@dataclass(frozen=True)
class RuntimeConfig:
    artifact_dir: Path
    model_version: str
    confidence_threshold: float

    @classmethod
    def from_env(cls) -> RuntimeConfig:
        artifact_dir = Path(os.getenv("NUTRISCAN_MODEL_ARTIFACT_DIR", DEFAULT_ARTIFACT_DIR))
        model_version = os.getenv("NUTRISCAN_MODEL_VERSION", DEFAULT_MODEL_VERSION)
        confidence_threshold = float(
            os.getenv("NUTRISCAN_CONFIDENCE_THRESHOLD", str(DEFAULT_CONFIDENCE_THRESHOLD))
        )
        return cls(
            artifact_dir=artifact_dir,
            model_version=model_version,
            confidence_threshold=confidence_threshold,
        )


def validate_model_artifact(config: RuntimeConfig) -> dict[str, object]:
    missing_files = [
        file_name
        for file_name in REQUIRED_ARTIFACT_FILES
        if not (config.artifact_dir / file_name).is_file()
    ]
    if missing_files:
        raise ArtifactValidationError(
            "Model Artifact is missing required files: " + ", ".join(missing_files)
        )

    label_map = _read_json_object(config.artifact_dir / "label_map.json")
    training_config = _read_json_object(config.artifact_dir / "training_config_resolved.json")
    labels = _labels_from_label_map(label_map)
    class_names = training_config.get("class_names")
    if not isinstance(class_names, list) or not all(isinstance(item, str) for item in class_names):
        raise ArtifactValidationError("training_config_resolved.json must include class_names")

    if set(labels) != EXPECTED_MVP_LABELS:
        raise ArtifactValidationError("label_map.json does not cover the eight MVP Food Categories")
    if set(class_names) != EXPECTED_MVP_LABELS:
        raise ArtifactValidationError(
            "training_config_resolved.json does not cover the eight MVP Food Categories"
        )
    if labels != class_names:
        raise ArtifactValidationError("label_map.json and training_config_resolved.json disagree")

    model_name = training_config.get("model_name")
    image_size = training_config.get("image_size")
    num_classes = training_config.get("num_classes")
    if not isinstance(model_name, str) or not model_name:
        raise ArtifactValidationError("training_config_resolved.json must include model_name")
    if not isinstance(image_size, int) or image_size <= 0:
        raise ArtifactValidationError("training_config_resolved.json must include image_size")
    if num_classes != len(EXPECTED_MVP_LABELS):
        raise ArtifactValidationError("training_config_resolved.json has an invalid num_classes")

    return {
        "status": "ready",
        "modelVersion": config.model_version,
        "artifactLocation": str(config.artifact_dir),
        "modelName": model_name,
        "imageSize": image_size,
        "labelCount": len(labels),
        "device": _runtime_device(),
        "confidenceThreshold": config.confidence_threshold,
    }


def _read_json_object(path: Path) -> dict[str, Any]:
    try:
        value = json.loads(path.read_text())
    except json.JSONDecodeError as exc:
        raise ArtifactValidationError(f"{path.name} must be valid JSON") from exc
    if not isinstance(value, dict):
        raise ArtifactValidationError(f"{path.name} must contain a JSON object")
    return value


def _labels_from_label_map(label_map: dict[str, Any]) -> list[str]:
    id_to_label = label_map.get("idToLabel")
    if not isinstance(id_to_label, dict):
        raise ArtifactValidationError("label_map.json must include idToLabel")
    labels: list[str] = []
    for index in range(len(id_to_label)):
        label = id_to_label.get(str(index))
        if not isinstance(label, str):
            raise ArtifactValidationError("label_map.json idToLabel must use contiguous string ids")
        labels.append(label)
    return labels


def _runtime_device() -> str:
    try:
        import torch
    except ImportError:
        return "cpu"
    return "cuda" if torch.cuda.is_available() else "cpu"
