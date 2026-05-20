import json
from pathlib import Path

import pytest

from app.classifier import FoodClassifier
from app.runtime import ArtifactValidationError, RuntimeConfig, validate_model_artifact

MVP_LABELS = [
    "bakso",
    "gado_gado",
    "gudeg",
    "nasi_goreng",
    "pempek",
    "rendang",
    "sate",
    "soto",
]


def test_runtime_config_resolves_selected_artifact_defaults(monkeypatch) -> None:
    monkeypatch.delenv("NUTRISCAN_MODEL_ARTIFACT_DIR", raising=False)
    monkeypatch.delenv("NUTRISCAN_MODEL_VERSION", raising=False)
    monkeypatch.delenv("NUTRISCAN_CONFIDENCE_THRESHOLD", raising=False)

    config = RuntimeConfig.from_env()

    assert config.artifact_dir == Path("model-artifacts/selected-mvp-classifier")
    assert config.model_version == "selected-mvp-classifier"
    assert config.confidence_threshold == 0.6


def test_runtime_config_allows_env_overrides(monkeypatch, tmp_path) -> None:
    artifact_dir = tmp_path / "artifact"
    monkeypatch.setenv("NUTRISCAN_MODEL_ARTIFACT_DIR", str(artifact_dir))
    monkeypatch.setenv("NUTRISCAN_MODEL_VERSION", "convnext-tiny-test")
    monkeypatch.setenv("NUTRISCAN_CONFIDENCE_THRESHOLD", "0.72")

    config = RuntimeConfig.from_env()

    assert config.artifact_dir == artifact_dir
    assert config.model_version == "convnext-tiny-test"
    assert config.confidence_threshold == 0.72


def test_valid_model_artifact_returns_readiness_metadata(tmp_path) -> None:
    artifact_dir = tmp_path / "selected-mvp-classifier"
    artifact_dir.mkdir()
    (artifact_dir / "model.pt").write_bytes(b"runtime model weights live outside git")
    (artifact_dir / "label_map.json").write_text(
        json.dumps(
            {
                "idToLabel": {str(index): label for index, label in enumerate(MVP_LABELS)},
                "class_to_idx": {label: index for index, label in enumerate(MVP_LABELS)},
            }
        )
    )
    (artifact_dir / "training_config_resolved.json").write_text(
        json.dumps(
            {
                "model_name": "convnext_tiny.fb_in1k",
                "num_classes": 8,
                "class_names": MVP_LABELS,
                "image_size": 256,
            }
        )
    )
    config = RuntimeConfig(
        artifact_dir=artifact_dir,
        model_version="convnext-tiny-test",
        confidence_threshold=0.6,
    )

    metadata = validate_model_artifact(config)

    assert metadata == {
        "status": "ready",
        "modelVersion": "convnext-tiny-test",
        "artifactLocation": str(artifact_dir),
        "modelName": "convnext_tiny.fb_in1k",
        "imageSize": 256,
        "labelCount": 8,
        "device": metadata["device"],
        "confidenceThreshold": 0.6,
    }
    assert metadata["device"] in {"cpu", "cuda"}


def test_missing_artifact_files_fail_readiness(tmp_path) -> None:
    artifact_dir = tmp_path / "selected-mvp-classifier"
    artifact_dir.mkdir()
    config = RuntimeConfig(
        artifact_dir=artifact_dir,
        model_version="convnext-tiny-test",
        confidence_threshold=0.6,
    )

    with pytest.raises(ArtifactValidationError, match="missing required files"):
        validate_model_artifact(config)


def test_stale_artifact_labels_fail_readiness(tmp_path) -> None:
    artifact_dir = tmp_path / "selected-mvp-classifier"
    artifact_dir.mkdir()
    (artifact_dir / "model.pt").write_bytes(b"runtime model weights live outside git")
    stale_labels = MVP_LABELS[:-1] + ["nasi_padang"]
    (artifact_dir / "label_map.json").write_text(
        json.dumps({"idToLabel": {str(index): label for index, label in enumerate(stale_labels)}})
    )
    (artifact_dir / "training_config_resolved.json").write_text(
        json.dumps(
            {
                "model_name": "convnext_tiny.fb_in1k",
                "num_classes": 8,
                "class_names": stale_labels,
                "image_size": 256,
            }
        )
    )
    config = RuntimeConfig(
        artifact_dir=artifact_dir,
        model_version="convnext-tiny-test",
        confidence_threshold=0.6,
    )

    with pytest.raises(ArtifactValidationError, match="eight MVP Food Categories"):
        validate_model_artifact(config)


def test_classifier_does_not_fallback_to_hardcoded_labels_when_label_map_missing(tmp_path) -> None:
    artifact_dir = tmp_path / "selected-mvp-classifier"
    artifact_dir.mkdir()
    (artifact_dir / "model.pt").write_bytes(b"runtime model weights live outside git")

    classifier = FoodClassifier(
        config=RuntimeConfig(
            artifact_dir=artifact_dir,
            model_version="convnext-tiny-test",
            confidence_threshold=0.6,
        )
    )

    assert classifier.labels == [], (
        "Classifier must not fall back to hardcoded labels when label_map.json is missing"
    )


def test_classifier_fails_when_artifact_is_missing(tmp_path) -> None:
    config = RuntimeConfig(
        artifact_dir=tmp_path / "missing-selected-artifact",
        model_version="convnext-tiny-test",
        confidence_threshold=0.6,
    )
    classifier = FoodClassifier(config=config)

    with pytest.raises(ArtifactValidationError, match="missing required files"):
        classifier.predict(b"image bytes")


def test_classifier_preprocessing_recipe_matches_selected_artifact_metadata(tmp_path) -> None:
    artifact_dir = tmp_path / "selected-mvp-classifier"
    artifact_dir.mkdir()
    (artifact_dir / "model.pt").write_bytes(b"runtime model weights live outside git")
    (artifact_dir / "label_map.json").write_text(
        json.dumps({"idToLabel": {str(index): label for index, label in enumerate(MVP_LABELS)}})
    )
    (artifact_dir / "training_config_resolved.json").write_text(
        json.dumps(
            {
                "model_name": "convnext_tiny.fb_in1k",
                "num_classes": 8,
                "class_names": MVP_LABELS,
                "image_size": 256,
            }
        )
    )
    classifier = FoodClassifier(
        config=RuntimeConfig(
            artifact_dir=artifact_dir,
            model_version="convnext-tiny-test",
            confidence_threshold=0.6,
        )
    )

    assert classifier.preprocessing_recipe() == {
        "imageSize": 256,
        "resizeSize": 294,
        "normalizationMean": [0.485, 0.456, 0.406],
        "normalizationStd": [0.229, 0.224, 0.225],
    }
