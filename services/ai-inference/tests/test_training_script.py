import json
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def test_training_config_records_model_and_labels() -> None:
    from scripts.train_classifier import load_config

    config = json.loads((ROOT / "configs" / "baseline_training.json").read_text())
    class_map = json.loads((ROOT / "configs" / "mvp_food_categories.json").read_text())
    resolved = load_config(ROOT / "configs" / "baseline_training.json")

    assert config["model_name"] in {"efficientnet_b0", "mobilenetv3_large_100"}
    assert config["image_size"] > 0
    assert config["batch_size"] > 0
    assert config["epochs"] > 0
    assert resolved.class_names == class_map["classes"]
    assert config["output_dir"] == "model-artifacts/baseline-food-classifier"


def test_training_script_dry_run_validates_config() -> None:
    result = subprocess.run(
        [
            sys.executable,
            "scripts/train_classifier.py",
            "--config",
            "configs/baseline_training.json",
            "--processed-dir",
            "data/processed",
            "--dry-run",
        ],
        cwd=ROOT,
        check=True,
        capture_output=True,
        text=True,
    )

    assert "model=efficientnet_b0" in result.stdout
    assert "classes=8" in result.stdout
    assert "artifactDir=model-artifacts/baseline-food-classifier" in result.stdout


def test_training_script_dry_run_reports_label_smoothing() -> None:
    result = subprocess.run(
        [
            sys.executable,
            "scripts/train_classifier.py",
            "--config",
            "configs/baseline_training_v2.json",
            "--processed-dir",
            "data/processed-v0.2",
            "--dry-run",
        ],
        cwd=ROOT,
        check=True,
        capture_output=True,
        text=True,
    )

    assert "model=efficientnet_b0" in result.stdout
    assert "labelSmoothing=0.1" in result.stdout
    assert "cropScale=0.75-1" in result.stdout
    assert "rotation=10" in result.stdout
    assert "colorJitter=0.15/0.15/0.1" in result.stdout
    assert "randomErasing=0" in result.stdout


def test_training_criterion_uses_configured_label_smoothing() -> None:
    from scripts.train_classifier import create_criterion, load_config

    config = load_config(ROOT / "configs" / "baseline_training_v2.json")
    criterion = create_criterion(config)

    assert criterion.label_smoothing == 0.1


def test_selected_training_config_uses_strong_context_augmentation_settings() -> None:
    from scripts.train_classifier import load_config

    config = load_config(ROOT / "configs" / "selected_mvp_classifier.json")

    assert config.random_resized_crop_scale == (0.55, 1.0)
    assert config.horizontal_flip_p == 0.5
    assert config.rotation_degrees == 15
    assert config.color_jitter_brightness == 0.25
    assert config.color_jitter_contrast == 0.25
    assert config.color_jitter_saturation == 0.20
    assert config.random_erasing_p == 0.0
