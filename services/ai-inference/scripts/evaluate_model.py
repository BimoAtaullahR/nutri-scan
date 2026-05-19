#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
from pathlib import Path
from typing import Any


def load_prediction_file(path: Path) -> tuple[list[str], list[dict[str, Any]]]:
    raw = json.loads(path.read_text())

    if isinstance(raw, list):
        predictions = raw
        labels = infer_labels(predictions)
        return labels, predictions

    if "predictions" in raw:
        predictions = list(raw["predictions"])
        labels = list(raw.get("classes") or infer_labels(predictions))
        return labels, predictions

    # Backward compatibility for the earlier test/report shape.
    labels = list(raw["labels"])
    predictions = [
        {
            "true_label": example["actual"],
            "top1_label": example["predicted"][0],
            "top3": [{"label": label, "confidence": None} for label in example["predicted"][:3]],
        }
        for example in raw["examples"]
    ]
    return labels, predictions


def infer_labels(predictions: list[dict[str, Any]]) -> list[str]:
    labels: list[str] = []
    for prediction in predictions:
        true_label = str(prediction["true_label"])
        if true_label not in labels:
            labels.append(true_label)
        for item in prediction.get("top3", []):
            label = str(item["label"])
            if label not in labels:
                labels.append(label)
    return labels


def top3_labels(prediction: dict[str, Any]) -> list[str]:
    return [str(item["label"]) for item in prediction.get("top3", [])]


def evaluate_predictions(labels: list[str], predictions: list[dict[str, Any]]) -> dict[str, Any]:
    label_to_index = {label: index for index, label in enumerate(labels)}
    matrix = [[0 for _ in labels] for _ in labels]
    top1_correct = 0
    top3_correct = 0

    for prediction in predictions:
        true_label = str(prediction["true_label"])
        top1_label = str(prediction.get("top1_label") or top3_labels(prediction)[0])
        top3 = top3_labels(prediction)

        if true_label not in label_to_index:
            raise ValueError(f"Unknown true label: {true_label}")
        if top1_label not in label_to_index:
            raise ValueError(f"Unknown predicted label: {top1_label}")

        matrix[label_to_index[true_label]][label_to_index[top1_label]] += 1
        if top1_label == true_label:
            top1_correct += 1
        if true_label in top3[:3]:
            top3_correct += 1

    total = len(predictions)
    per_class_metrics = compute_per_class_metrics(labels, matrix)
    top1_accuracy = top1_correct / total if total else 0.0
    top3_accuracy = top3_correct / total if total else 0.0

    metrics = {
        "top1_accuracy": top1_accuracy,
        "top3_accuracy": top3_accuracy,
        "num_test_samples": total,
        "classes": labels,
        "mvp_targets": {
            "top1_accuracy": 0.8,
            "top3_accuracy": 0.9,
        },
        "meets_mvp_target": top1_accuracy >= 0.8 and top3_accuracy >= 0.9,
        # Backward-compatible aliases for earlier repo tests/reports.
        "top1Accuracy": top1_accuracy,
        "top3Accuracy": top3_accuracy,
        "meetsMvpTarget": top1_accuracy >= 0.8 and top3_accuracy >= 0.9,
    }

    return {
        "metrics": metrics,
        "per_class_metrics": per_class_metrics,
        "confusion_matrix": {"classes": labels, "labels": labels, "matrix": matrix},
    }


def compute_per_class_metrics(
    labels: list[str],
    matrix: list[list[int]],
) -> dict[str, dict[str, float | int]]:
    per_class: dict[str, dict[str, float | int]] = {}

    for index, label in enumerate(labels):
        true_positive = matrix[index][index]
        false_negative = sum(matrix[index]) - true_positive
        false_positive = sum(row[index] for row in matrix) - true_positive
        support = sum(matrix[index])

        precision = (
            true_positive / (true_positive + false_positive)
            if true_positive + false_positive
            else 0.0
        )
        recall = (
            true_positive / (true_positive + false_negative)
            if true_positive + false_negative
            else 0.0
        )
        f1 = 2 * precision * recall / (precision + recall) if precision + recall else 0.0
        per_class[label] = {
            "precision": precision,
            "recall": recall,
            "f1": f1,
            "support": support,
        }

    return per_class


def write_reports(report_dir: Path, evaluation: dict[str, Any]) -> None:
    report_dir.mkdir(parents=True, exist_ok=True)
    (report_dir / "metrics.json").write_text(
        json.dumps(evaluation["metrics"], indent=2, sort_keys=True) + "\n"
    )
    (report_dir / "per_class_metrics.json").write_text(
        json.dumps(evaluation["per_class_metrics"], indent=2, sort_keys=True) + "\n"
    )
    (report_dir / "confusion_matrix.json").write_text(
        json.dumps(evaluation["confusion_matrix"], indent=2) + "\n"
    )


def main() -> None:
    parser = argparse.ArgumentParser(description="Evaluate NutriScan MVP food classifier")
    parser.add_argument("--predictions-file", type=Path, required=True)
    parser.add_argument("--report-dir", type=Path, default=Path("reports/baseline-food-classifier"))
    args = parser.parse_args()

    labels, predictions = load_prediction_file(args.predictions_file)
    evaluation = evaluate_predictions(labels, predictions)
    write_reports(args.report_dir, evaluation)

    metrics = evaluation["metrics"]
    print(
        f"top1={metrics['top1_accuracy']:.4f} top3={metrics['top3_accuracy']:.4f} "
        f"top1>=80%={str(metrics['top1_accuracy'] >= 0.8).lower()} "
        f"top3>=90%={str(metrics['top3_accuracy'] >= 0.9).lower()}"
    )


if __name__ == "__main__":
    main()
