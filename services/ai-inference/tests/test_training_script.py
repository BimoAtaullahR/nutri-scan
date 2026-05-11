import json
import subprocess
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]


def test_training_config_records_model_and_labels() -> None:
    config = json.loads((ROOT / "configs" / "baseline_training.json").read_text())
    class_map = json.loads((ROOT / "configs" / "mvp_food_categories.json").read_text())

    assert config["modelName"] in {"efficientnet_b0", "mobilenetv3_large_100"}
    assert config["imageSize"] > 0
    assert config["batchSize"] > 0
    assert config["epochs"] > 0
    assert config["labels"] == [category["slug"] for category in class_map["categories"]]
    assert config["artifactDir"] == "model-artifacts/baseline-food-classifier"


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
