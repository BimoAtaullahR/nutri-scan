import subprocess
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SCRIPT = ROOT / "scripts" / "prepare_dataset.py"
CLASS_MAP = ROOT / "configs" / "mvp_food_categories.json"


def make_images(folder: Path, count: int) -> None:
    folder.mkdir(parents=True)
    for index in range(count):
        (folder / f"image_{index}.jpg").write_bytes(b"fake-image")


def test_prepare_dataset_creates_deterministic_splits(tmp_path: Path):
    raw_dir = tmp_path / "raw"
    processed_dir = tmp_path / "processed"
    make_images(raw_dir / "nasi goreng", 10)
    make_images(raw_dir / "gado-gado", 10)

    result = subprocess.run(
        [
            sys.executable,
            str(SCRIPT),
            "--raw-dir",
            str(raw_dir),
            "--processed-dir",
            str(processed_dir),
            "--class-map",
            str(CLASS_MAP),
            "--seed",
            "7",
        ],
        check=False,
        capture_output=True,
        text=True,
    )

    assert result.returncode == 0, result.stderr
    assert len(list((processed_dir / "train" / "nasi_goreng").iterdir())) == 7
    assert len(list((processed_dir / "validation" / "nasi_goreng").iterdir())) == 1
    assert len(list((processed_dir / "test" / "nasi_goreng").iterdir())) == 2
    assert len(list((processed_dir / "train" / "gado_gado").iterdir())) == 7
    assert "nasi_goreng: train=7 validation=1 test=2" in result.stdout


def test_prepare_dataset_dry_run_does_not_write_processed_data(tmp_path: Path):
    raw_dir = tmp_path / "raw"
    processed_dir = tmp_path / "processed"
    make_images(raw_dir / "soto", 10)

    result = subprocess.run(
        [
            sys.executable,
            str(SCRIPT),
            "--raw-dir",
            str(raw_dir),
            "--processed-dir",
            str(processed_dir),
            "--class-map",
            str(CLASS_MAP),
            "--dry-run",
        ],
        check=False,
        capture_output=True,
        text=True,
    )

    assert result.returncode == 0, result.stderr
    assert not processed_dir.exists()
    assert "soto: train=7 validation=1 test=2" in result.stdout
