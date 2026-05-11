import json
import subprocess
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]


def test_evaluation_script_writes_metrics_and_confusion_matrix(tmp_path: Path) -> None:
    predictions = {
        "labels": ["nasi_goreng", "sate", "rendang"],
        "examples": [
            {"actual": "nasi_goreng", "predicted": ["nasi_goreng", "sate", "rendang"]},
            {"actual": "sate", "predicted": ["rendang", "sate", "nasi_goreng"]},
            {"actual": "rendang", "predicted": ["rendang", "sate", "nasi_goreng"]},
        ],
    }
    predictions_path = tmp_path / "predictions.json"
    predictions_path.write_text(json.dumps(predictions))
    report_dir = tmp_path / "reports"

    result = subprocess.run(
        [
            sys.executable,
            "scripts/evaluate_model.py",
            "--predictions-file",
            str(predictions_path),
            "--report-dir",
            str(report_dir),
        ],
        cwd=ROOT,
        check=True,
        capture_output=True,
        text=True,
    )

    metrics = json.loads((report_dir / "metrics.json").read_text())
    confusion_matrix = json.loads((report_dir / "confusion_matrix.json").read_text())

    assert metrics["top1Accuracy"] == 2 / 3
    assert metrics["top3Accuracy"] == 1.0
    assert metrics["meetsMvpTarget"] is False
    assert "top1>=80%=false" in result.stdout
    assert confusion_matrix["labels"] == predictions["labels"]
    assert confusion_matrix["matrix"] == [[1, 0, 0], [0, 0, 1], [0, 0, 1]]
