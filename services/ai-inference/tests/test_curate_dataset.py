import json
import subprocess
import sys
from pathlib import Path

from PIL import Image


ROOT = Path(__file__).resolve().parents[1]


def write_image(path: Path) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    Image.new("RGB", (8, 8), color=(255, 0, 0)).save(path)


def test_curate_dataset_reports_counts_and_risks(tmp_path: Path) -> None:
    processed_dir = tmp_path / "processed"
    write_image(processed_dir / "train" / "sate" / "a.jpg")
    write_image(processed_dir / "validation" / "sate" / "b.jpg")
    write_image(processed_dir / "test" / "sate" / "c.jpg")
    write_image(processed_dir / "train" / "bakso" / "a.jpg")

    report_path = tmp_path / "curation_report.json"

    result = subprocess.run(
        [
            sys.executable,
            "scripts/curate_dataset.py",
            "--processed-dir",
            str(processed_dir),
            "--class-map",
            "configs/mvp_food_categories.json",
            "--report-path",
            str(report_path),
            "--minimum-reviewed",
            "3",
        ],
        cwd=ROOT,
        check=True,
        capture_output=True,
        text=True,
    )

    report = json.loads(report_path.read_text())

    assert report["classes"]["sate"]["total"] == 3
    assert report["classes"]["sate"]["meetsMinimumReviewed"] is True
    assert report["classes"]["bakso"]["total"] == 1
    assert report["classes"]["bakso"]["meetsMinimumReviewed"] is False
    assert "missing split folder" in report["classes"]["bakso"]["risks"][0]
    assert "sate: train=1 validation=1 test=1 total=3" in result.stdout


def test_curate_dataset_dry_run_does_not_write_report(tmp_path: Path) -> None:
    processed_dir = tmp_path / "processed"
    write_image(processed_dir / "train" / "sate" / "a.jpg")
    report_path = tmp_path / "curation_report.json"

    subprocess.run(
        [
            sys.executable,
            "scripts/curate_dataset.py",
            "--processed-dir",
            str(processed_dir),
            "--class-map",
            "configs/mvp_food_categories.json",
            "--report-path",
            str(report_path),
            "--dry-run",
        ],
        cwd=ROOT,
        check=True,
        capture_output=True,
        text=True,
    )

    assert not report_path.exists()
