#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
import shutil
from collections import Counter
from dataclasses import dataclass
from pathlib import Path
from typing import Any

PATH_KEYS = ("path", "image_path", "file_path")
TRUE_LABEL_KEYS = ("true_label", "label", "target")
PREDICTED_LABEL_KEYS = ("top1_label", "predicted_label", "prediction")
CONFIDENCE_KEYS = ("top1_confidence", "confidence", "confidence_score")


@dataclass(frozen=True)
class PredictionItem:
    image_path: Path
    true_label: str
    predicted_label: str
    confidence: float | None


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Export misclassified NutriScan test images for manual review."
    )
    parser.add_argument("--predictions-file", type=Path, required=True)
    parser.add_argument("--output-dir", type=Path, required=True)
    return parser.parse_args()


def load_prediction_items(predictions_file: Path) -> list[dict[str, Any]]:
    raw = json.loads(predictions_file.read_text())

    if isinstance(raw, list):
        return raw
    if isinstance(raw, dict) and "predictions" in raw:
        return list(raw["predictions"])
    if isinstance(raw, dict) and "examples" in raw:
        return [
            {
                "path": example.get("path")
                or example.get("image_path")
                or example.get("file_path"),
                "true_label": example.get("actual"),
                "top1_label": example.get("predicted", [None])[0],
            }
            for example in raw["examples"]
        ]

    raise ValueError(
        "Unsupported predictions format. Expected a list or an object with 'predictions'."
    )


def first_present(item: dict[str, Any], keys: tuple[str, ...]) -> Any:
    for key in keys:
        value = item.get(key)
        if value is not None and value != "":
            return value
    return None


def parse_prediction(item: dict[str, Any], index: int) -> PredictionItem:
    image_path = first_present(item, PATH_KEYS)
    true_label = first_present(item, TRUE_LABEL_KEYS)
    predicted_label = first_present(item, PREDICTED_LABEL_KEYS)
    confidence = first_present(item, CONFIDENCE_KEYS)

    if image_path is None or true_label is None or predicted_label is None:
        raise ValueError(
            "Invalid prediction item at index "
            f"{index}. Required image path, true label, and predicted label. "
            f"First invalid item: {json.dumps(item, ensure_ascii=False)}"
        )

    return PredictionItem(
        image_path=Path(str(image_path)),
        true_label=str(true_label),
        predicted_label=str(predicted_label),
        confidence=float(confidence) if confidence is not None else None,
    )


def safe_name(value: str) -> str:
    allowed = []
    for char in value.strip():
        if char.isalnum() or char in ("-", "_", "."):
            allowed.append(char)
        else:
            allowed.append("_")
    return "".join(allowed).strip("._") or "unknown"


def resolve_image_path(image_path: Path, cwd: Path) -> Path:
    if image_path.is_absolute():
        return image_path
    return cwd / image_path


def ensure_inside_directory(path: Path, directory: Path) -> None:
    resolved_path = path.resolve()
    resolved_directory = directory.resolve()
    if resolved_path != resolved_directory and not resolved_path.is_relative_to(resolved_directory):
        raise ValueError(f"Refusing to write outside output directory: {resolved_path}")


def build_destination(
    output_dir: Path,
    group_name: str,
    sequence: int,
    prediction: PredictionItem,
) -> Path:
    source_name = safe_name(prediction.image_path.name)
    confidence = prediction.confidence
    confidence_text = "na" if confidence is None else f"{confidence:.3f}"
    filename = (
        f"{sequence:04d}__true-{safe_name(prediction.true_label)}"
        f"__pred-{safe_name(prediction.predicted_label)}"
        f"__conf-{confidence_text}__original-{source_name}"
    )
    return output_dir / group_name / filename


def export_misclassified(predictions_file: Path, output_dir: Path) -> dict[str, Any]:
    raw_items = load_prediction_items(predictions_file)
    predictions = [parse_prediction(item, index) for index, item in enumerate(raw_items)]
    output_dir.mkdir(parents=True, exist_ok=True)

    groups: Counter[str] = Counter()
    missing_files: list[dict[str, str]] = []
    total_misclassified = 0
    cwd = Path.cwd()

    for prediction in predictions:
        if prediction.true_label == prediction.predicted_label:
            continue

        total_misclassified += 1
        group_name = (
            f"{safe_name(prediction.true_label)}"
            f"_as_{safe_name(prediction.predicted_label)}"
        )
        source_path = resolve_image_path(prediction.image_path, cwd)

        if not source_path.is_file():
            missing_files.append(
                {
                    "path": str(prediction.image_path),
                    "true_label": prediction.true_label,
                    "top1_label": prediction.predicted_label,
                    "group": group_name,
                }
            )
            continue

        groups[group_name] += 1
        destination = build_destination(output_dir, group_name, groups[group_name], prediction)
        ensure_inside_directory(destination, output_dir)
        destination.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(source_path, destination)

    summary = {
        "predictions_file": str(predictions_file),
        "output_dir": str(output_dir),
        "total_predictions": len(predictions),
        "total_misclassified": total_misclassified,
        "missing_files": missing_files,
        "groups": {group: {"count": count} for group, count in sorted(groups.items())},
    }
    summary_path = output_dir / "summary.json"
    ensure_inside_directory(summary_path, output_dir)
    summary_path.write_text(json.dumps(summary, indent=2, sort_keys=True) + "\n")
    return summary


def print_summary(summary: dict[str, Any]) -> None:
    print("Exported misclassified images:")
    groups = summary["groups"]
    if groups:
        for group, values in groups.items():
            print(f"- {group}: {values['count']}")
    else:
        print("- none")

    missing_count = len(summary["missing_files"])
    if missing_count:
        print(f"Missing source files: {missing_count}")

    print(f"Summary written to {Path(summary['output_dir']) / 'summary.json'}")


def main() -> None:
    args = parse_args()
    summary = export_misclassified(args.predictions_file, args.output_dir)
    print_summary(summary)


if __name__ == "__main__":
    main()
