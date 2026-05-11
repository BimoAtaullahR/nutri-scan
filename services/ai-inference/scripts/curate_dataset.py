#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path


IMAGE_EXTENSIONS = {".jpg", ".jpeg", ".png", ".webp"}
SPLITS = ("train", "validation", "test")


def main(argv: list[str] | None = None) -> int:
    args = parse_args(argv)
    categories = load_categories(args.class_map)
    report = build_curation_report(
        processed_dir=args.processed_dir,
        categories=categories,
        minimum_reviewed=args.minimum_reviewed,
    )

    for category, details in report["classes"].items():
        split_counts = details["splits"]
        print(
            f"{category}: train={split_counts['train']} "
            f"validation={split_counts['validation']} test={split_counts['test']} "
            f"total={details['total']}"
        )

    if not args.dry_run:
        args.report_path.parent.mkdir(parents=True, exist_ok=True)
        args.report_path.write_text(json.dumps(report, indent=2, sort_keys=True) + "\n")

    return 0


def parse_args(argv: list[str] | None) -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Audit manually curated NutriScan processed dataset counts."
    )
    parser.add_argument("--processed-dir", type=Path, required=True)
    parser.add_argument("--class-map", type=Path, required=True)
    parser.add_argument(
        "--report-path",
        type=Path,
        default=Path("reports/dataset-curation/curation_report.json"),
    )
    parser.add_argument("--minimum-reviewed", type=int, default=100)
    parser.add_argument("--dry-run", action="store_true")
    return parser.parse_args(argv)


def load_categories(class_map_path: Path) -> list[str]:
    class_map = json.loads(class_map_path.read_text())
    return [category["slug"] for category in class_map["categories"]]


def build_curation_report(
    processed_dir: Path,
    categories: list[str],
    minimum_reviewed: int,
) -> dict[str, object]:
    return {
        "processedDir": str(processed_dir),
        "minimumReviewedPerClass": minimum_reviewed,
        "reviewPolicy": [
            "Remove watermarks, menus, drawings, duplicate images, and irrelevant images.",
            "Review at least the minimum image count per MVP category when available.",
            "Record weak classes and data quality risks before model training.",
        ],
        "classes": {
            category: inspect_category(processed_dir, category, minimum_reviewed)
            for category in categories
        },
    }


def inspect_category(
    processed_dir: Path,
    category: str,
    minimum_reviewed: int,
) -> dict[str, object]:
    split_counts: dict[str, int] = {}
    risks: list[str] = []

    for split in SPLITS:
        split_dir = processed_dir / split / category
        if not split_dir.exists():
            split_counts[split] = 0
            risks.append(f"missing split folder: {split}/{category}")
            continue

        split_counts[split] = count_images(split_dir)

    total = sum(split_counts.values())
    if total < minimum_reviewed:
        risks.append(f"below minimum reviewed image target: {total}/{minimum_reviewed}")
    if split_counts["validation"] == 0:
        risks.append("validation split has no images")
    if split_counts["test"] == 0:
        risks.append("test split has no images")

    return {
        "splits": split_counts,
        "total": total,
        "meetsMinimumReviewed": total >= minimum_reviewed,
        "risks": risks,
    }


def count_images(directory: Path) -> int:
    return sum(
        1
        for path in directory.iterdir()
        if path.is_file() and path.suffix.lower() in IMAGE_EXTENSIONS
    )


if __name__ == "__main__":
    sys.exit(main())
