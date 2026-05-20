import json
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def test_export_misclassified_groups_copies_and_summary(tmp_path: Path) -> None:
    image_dir = tmp_path / "data" / "processed" / "test"
    bakso_image = image_dir / "bakso" / "example.jpg"
    soto_image = image_dir / "soto" / "missing.jpg"
    bakso_image.parent.mkdir(parents=True)
    bakso_image.write_bytes(b"fake-image")

    predictions_file = tmp_path / "predictions.json"
    predictions_file.write_text(
        json.dumps(
            {
                "predictions": [
                    {
                        "path": str(bakso_image),
                        "true_label": "bakso",
                        "top1_label": "soto",
                        "top1_confidence": 0.74,
                    },
                    {
                        "path": str(soto_image),
                        "true_label": "soto",
                        "top1_label": "bakso",
                        "top1_confidence": 0.61,
                    },
                    {
                        "path": str(bakso_image),
                        "true_label": "bakso",
                        "top1_label": "bakso",
                        "top1_confidence": 0.98,
                    },
                ]
            }
        )
    )
    output_dir = tmp_path / "reports" / "misclassified"

    result = subprocess.run(
        [
            sys.executable,
            "scripts/export_misclassified.py",
            "--predictions-file",
            str(predictions_file),
            "--output-dir",
            str(output_dir),
        ],
        cwd=ROOT,
        check=True,
        capture_output=True,
        text=True,
    )

    exported = list((output_dir / "bakso_as_soto").glob("*.jpg"))
    summary = json.loads((output_dir / "summary.json").read_text())

    assert len(exported) == 1
    assert exported[0].name.startswith("0001__true-bakso__pred-soto__conf-0.740")
    assert summary["total_predictions"] == 3
    assert summary["total_misclassified"] == 2
    assert summary["groups"]["bakso_as_soto"]["count"] == 1
    assert summary["missing_files"][0]["group"] == "soto_as_bakso"
    assert "bakso_as_soto: 1" in result.stdout


def test_export_misclassified_rejects_invalid_items(tmp_path: Path) -> None:
    predictions_file = tmp_path / "predictions.json"
    predictions_file.write_text(json.dumps({"predictions": [{"path": "image.jpg"}]}))

    result = subprocess.run(
        [
            sys.executable,
            "scripts/export_misclassified.py",
            "--predictions-file",
            str(predictions_file),
            "--output-dir",
            str(tmp_path / "out"),
        ],
        cwd=ROOT,
        check=False,
        capture_output=True,
        text=True,
    )

    assert result.returncode != 0
    assert "Invalid prediction item at index 0" in result.stderr
