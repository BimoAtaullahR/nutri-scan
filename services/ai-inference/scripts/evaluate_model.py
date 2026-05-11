#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
from pathlib import Path


def load_prediction_file(path: Path) -> tuple[list[str], list[dict[str, object]]]:
    raw = json.loads(path.read_text())
    labels = list(raw["labels"])
    examples = list(raw["examples"])
    return labels, examples


def evaluate_predictions(labels: list[str], examples: list[dict[str, object]]) -> dict[str, object]:
    label_to_index = {label: index for index, label in enumerate(labels)}
    matrix = [[0 for _ in labels] for _ in labels]
    per_class = {
        label: {"total": 0, "top1Correct": 0, "top3Correct": 0}
        for label in labels
    }
    top1_correct = 0
    top3_correct = 0

    for example in examples:
        actual = str(example["actual"])
        predicted = list(example["predicted"])
        if actual not in label_to_index:
            raise ValueError(f"Unknown actual label: {actual}")
        if not predicted:
            raise ValueError("Prediction list cannot be empty")

        first_prediction = str(predicted[0])
        if first_prediction not in label_to_index:
            raise ValueError(f"Unknown predicted label: {first_prediction}")

        actual_index = label_to_index[actual]
        predicted_index = label_to_index[first_prediction]
        matrix[actual_index][predicted_index] += 1
        per_class[actual]["total"] += 1

        if first_prediction == actual:
            top1_correct += 1
            per_class[actual]["top1Correct"] += 1
        if actual in predicted[:3]:
            top3_correct += 1
            per_class[actual]["top3Correct"] += 1

    total = len(examples)
    top1_accuracy = top1_correct / total if total else 0
    top3_accuracy = top3_correct / total if total else 0
    class_metrics = {}
    for label, counts in per_class.items():
        class_total = counts["total"]
        class_metrics[label] = {
            "top1Accuracy": counts["top1Correct"] / class_total if class_total else 0,
            "top3Accuracy": counts["top3Correct"] / class_total if class_total else 0,
            "support": class_total,
        }

    return {
        "metrics": {
            "top1Accuracy": top1_accuracy,
            "top3Accuracy": top3_accuracy,
            "top1Target": 0.8,
            "top3Target": 0.9,
            "meetsMvpTarget": top1_accuracy >= 0.8 and top3_accuracy >= 0.9,
            "perClass": class_metrics,
        },
        "confusionMatrix": {"labels": labels, "matrix": matrix},
    }


def write_reports(report_dir: Path, evaluation: dict[str, object]) -> None:
    report_dir.mkdir(parents=True, exist_ok=True)
    (report_dir / "metrics.json").write_text(
        json.dumps(evaluation["metrics"], indent=2, sort_keys=True) + "\n"
    )
    (report_dir / "confusion_matrix.json").write_text(
        json.dumps(evaluation["confusionMatrix"], indent=2) + "\n"
    )


def main() -> None:
    parser = argparse.ArgumentParser(description="Evaluate NutriScan MVP food classifier")
    parser.add_argument("--predictions-file", type=Path, required=True)
    parser.add_argument("--report-dir", type=Path, default=Path("reports/baseline-food-classifier"))
    args = parser.parse_args()

    labels, examples = load_prediction_file(args.predictions_file)
    evaluation = evaluate_predictions(labels, examples)
    write_reports(args.report_dir, evaluation)

    metrics = evaluation["metrics"]
    print(
        f"top1={metrics['top1Accuracy']:.4f} top3={metrics['top3Accuracy']:.4f} "
        f"top1>=80%={str(metrics['top1Accuracy'] >= 0.8).lower()} "
        f"top3>=90%={str(metrics['top3Accuracy'] >= 0.9).lower()}"
    )


if __name__ == "__main__":
    main()
