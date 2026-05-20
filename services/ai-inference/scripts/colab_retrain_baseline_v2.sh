#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

CONFIG="${CONFIG:-configs/selected_mvp_classifier.json}"
PROCESSED_DIR="${PROCESSED_DIR:-data/processed-v0.2}"
INSTALL_DEPS="${INSTALL_DEPS:-1}"
REQUIRE_CUDA="${REQUIRE_CUDA:-1}"
DRY_RUN_ONLY="${DRY_RUN_ONLY:-0}"

if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
  cat <<'EOF'
Usage:
  REQUIRE_CUDA=1 INSTALL_DEPS=1 bash scripts/colab_retrain_baseline_v2.sh

Environment:
  CONFIG         Training config path. Default: configs/selected_mvp_classifier.json
  PROCESSED_DIR  Processed dataset path. Default: data/processed-v0.2
  INSTALL_DEPS   Install Python dependencies before training. Default: 1
  REQUIRE_CUDA   Fail when CUDA is unavailable. Default: 1
  DRY_RUN_ONLY   Validate config and dataset layout without training. Default: 0
EOF
  exit 0
fi

if [[ "$INSTALL_DEPS" == "1" ]]; then
  echo "[1/6] Installing Python dependencies"
  python - <<'PY'
from __future__ import annotations

import subprocess
import sys

common_packages = [
    "fastapi",
    "python-multipart",
    "uvicorn[standard]",
    "pydantic-settings",
    "timm",
    "pillow",
    "numpy",
    "scikit-learn",
    "pyyaml",
]

if sys.version_info >= (3, 12):
    command = [sys.executable, "-m", "pip", "install", "-e", "."]
else:
    # Colab already ships CUDA-compatible torch/torchvision. Reinstalling them
    # can take a long time and may replace the CUDA build with a CPU-only wheel.
    command = [sys.executable, "-m", "pip", "install", *common_packages]

subprocess.check_call(command)
PY
else
  echo "[1/6] Skipping dependency install"
fi

echo "[2/6] Checking CUDA runtime"
python - <<'PY'
from __future__ import annotations

import os
import sys

import torch

require_cuda = os.environ.get("REQUIRE_CUDA", "1") == "1"
print(f"torch={torch.__version__}")
print(f"cuda_available={torch.cuda.is_available()}")
if torch.cuda.is_available():
    print(f"gpu={torch.cuda.get_device_name(0)}")
elif require_cuda:
    sys.exit("CUDA is not available. In Colab, switch Runtime > Change runtime type > T4 GPU.")
PY

if [[ ! -d "$PROCESSED_DIR" ]]; then
  cat >&2 <<EOF
Processed dataset folder not found: $PROCESSED_DIR

Expected layout:
  $PROCESSED_DIR/train/<class>/
  $PROCESSED_DIR/val/<class>/ or $PROCESSED_DIR/validation/<class>/
  $PROCESSED_DIR/test/<class>/

Set PROCESSED_DIR=/path/to/dataset if your folder is elsewhere.
EOF
  exit 1
fi

echo "[3/6] Validating training config and dataset layout"
python -u scripts/train_classifier.py \
  --config "$CONFIG" \
  --processed-dir "$PROCESSED_DIR" \
  --dry-run

if [[ "$DRY_RUN_ONLY" == "1" ]]; then
  echo "DRY_RUN_ONLY=1, skipping training and evaluation"
  exit 0
fi

echo "[4/6] Training classifier"
python -u scripts/train_classifier.py \
  --config "$CONFIG" \
  --processed-dir "$PROCESSED_DIR"

REPORT_DIR="$(CONFIG_PATH="$CONFIG" python - <<'PY'
import json
import os
from pathlib import Path

config = json.loads(Path(os.environ["CONFIG_PATH"]).read_text())
print(config["report_dir"])
PY
)"

echo "[5/6] Evaluating predictions"
python -u scripts/evaluate_model.py \
  --predictions-file "$REPORT_DIR/predictions.json" \
  --report-dir "$REPORT_DIR"

echo "[6/6] Exporting misclassified images and printing summary"
python -u scripts/export_misclassified.py \
  --predictions-file "$REPORT_DIR/predictions.json" \
  --output-dir "$REPORT_DIR/misclassified"

REPORT_DIR="$REPORT_DIR" python -u - <<'PY'
from __future__ import annotations

import json
import os
from pathlib import Path

report_dir = Path(os.environ["REPORT_DIR"])
metrics = json.loads((report_dir / "metrics.json").read_text())
per_class = json.loads((report_dir / "per_class_metrics.json").read_text())

print("\nSummary")
print(f"top1_accuracy={metrics['top1_accuracy']:.4f}")
print(f"top3_accuracy={metrics['top3_accuracy']:.4f}")
print(f"num_test_samples={metrics['num_test_samples']}")
print(f"meets_mvp_target={metrics['meets_mvp_target']}")
print("\nWeak classes")
for label in ("rendang", "gado_gado", "soto"):
    values = per_class[label]
    print(
        f"{label}: precision={values['precision']:.4f} "
        f"recall={values['recall']:.4f} f1={values['f1']:.4f} "
        f"support={values['support']}"
    )
PY
