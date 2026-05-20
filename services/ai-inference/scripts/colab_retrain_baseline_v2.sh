#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

CONFIG="${CONFIG:-configs/baseline_training_v2.json}"
PROCESSED_DIR="${PROCESSED_DIR:-data/processed-v0.2}"
INSTALL_DEPS="${INSTALL_DEPS:-1}"
REQUIRE_CUDA="${REQUIRE_CUDA:-1}"

if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
  cat <<'EOF'
Usage:
  REQUIRE_CUDA=1 INSTALL_DEPS=1 bash scripts/colab_retrain_baseline_v2.sh

Environment:
  CONFIG         Training config path. Default: configs/baseline_training_v2.json
  PROCESSED_DIR  Processed dataset path. Default: data/processed-v0.2
  INSTALL_DEPS   Install Python dependencies before training. Default: 1
  REQUIRE_CUDA   Fail when CUDA is unavailable. Default: 1
EOF
  exit 0
fi

if [[ "$INSTALL_DEPS" == "1" ]]; then
  python - <<'PY'
from __future__ import annotations

import subprocess
import sys

if sys.version_info >= (3, 12):
    command = [sys.executable, "-m", "pip", "install", "-q", "-e", "."]
else:
    command = [
        sys.executable,
        "-m",
        "pip",
        "install",
        "-q",
        "fastapi",
        "python-multipart",
        "uvicorn[standard]",
        "pydantic-settings",
        "torch",
        "torchvision",
        "timm",
        "pillow",
        "numpy",
        "scikit-learn",
        "pyyaml",
    ]

subprocess.check_call(command)
PY
fi

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

python scripts/train_classifier.py \
  --config "$CONFIG" \
  --processed-dir "$PROCESSED_DIR" \
  --dry-run

python scripts/train_classifier.py \
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

python scripts/evaluate_model.py \
  --predictions-file "$REPORT_DIR/predictions.json" \
  --report-dir "$REPORT_DIR"

python scripts/export_misclassified.py \
  --predictions-file "$REPORT_DIR/predictions.json" \
  --output-dir "$REPORT_DIR/misclassified"

REPORT_DIR="$REPORT_DIR" python - <<'PY'
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
