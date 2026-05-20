#!/usr/bin/env python3
from __future__ import annotations

import argparse
import csv
import json
import shutil
import sys
from collections import Counter
from dataclasses import dataclass
from datetime import UTC, datetime
from pathlib import Path

IMAGE_EXTENSIONS = {".jpg", ".jpeg", ".png", ".webp"}
VALID_DECISIONS = {
    "keep",
    "relabel",
    "reject_ambiguous",
    "reject_bad_quality",
    "duplicate",
}
REMOVE_DECISIONS = {"reject_ambiguous", "reject_bad_quality", "duplicate"}


@dataclass(frozen=True)
class ReviewRow:
    row_number: int
    image_path: Path
    current_label: str
    predicted_label: str
    decision: str
    new_label: str
    reason: str


@dataclass(frozen=True)
class AppliedAction:
    row_number: int
    decision: str
    source_path: str
    output_path: str
    new_output_path: str | None = None


def main(argv: list[str] | None = None) -> int:
    args = parse_args(argv)

    review_rows = load_review_rows(args.review_csv)
    validate_review_rows(review_rows, args.source_processed_dir)

    if args.output_processed_dir.resolve() == args.source_processed_dir.resolve():
        raise SystemExit("output processed dir must be different from source processed dir")

    planned_actions = plan_actions(
        review_rows=review_rows,
        source_processed_dir=args.source_processed_dir,
        output_processed_dir=args.output_processed_dir,
    )

    report = build_report(args, review_rows, planned_actions)

    if not args.dry_run:
        copy_processed_dataset(
            args.source_processed_dir,
            args.output_processed_dir,
            force=args.force,
        )
        apply_actions(planned_actions)
        args.report_path.parent.mkdir(parents=True, exist_ok=True)
        args.report_path.write_text(json.dumps(report, indent=2, sort_keys=True) + "\n")

    print_summary(report, dry_run=args.dry_run)
    return 0


def parse_args(argv: list[str] | None) -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description=(
            "Create a cleaned processed dataset copy by applying "
            "misclassified_review.csv decisions."
        )
    )
    parser.add_argument(
        "--review-csv",
        type=Path,
        required=True,
        help="CSV with image_path,current_label,predicted_label,decision,new_label,reason.",
    )
    parser.add_argument(
        "--source-processed-dir",
        type=Path,
        required=True,
        help="Existing processed dataset folder, e.g. data/processed.",
    )
    parser.add_argument(
        "--output-processed-dir",
        type=Path,
        required=True,
        help="New processed dataset folder to create, e.g. data/processed-v0.2.",
    )
    parser.add_argument(
        "--report-path",
        type=Path,
        default=Path("reports/dataset-curation/misclassified_review_apply_report.json"),
    )
    parser.add_argument(
        "--force",
        action="store_true",
        help="Replace output processed dir if it already exists.",
    )
    parser.add_argument("--dry-run", action="store_true")
    return parser.parse_args(argv)


def load_review_rows(review_csv: Path) -> list[ReviewRow]:
    if not review_csv.exists():
        raise SystemExit(f"review CSV does not exist: {review_csv}")

    rows: list[ReviewRow] = []
    with review_csv.open(newline="", encoding="utf-8-sig") as f:
        reader = csv.DictReader(f)
        required_fields = {
            "image_path",
            "current_label",
            "predicted_label",
            "decision",
            "new_label",
            "reason",
        }
        missing_fields = required_fields - set(reader.fieldnames or [])
        if missing_fields:
            missing = ", ".join(sorted(missing_fields))
            raise SystemExit(f"review CSV is missing required columns: {missing}")

        for row_number, row in enumerate(reader, start=2):
            rows.append(
                ReviewRow(
                    row_number=row_number,
                    image_path=Path(row["image_path"].strip()),
                    current_label=row["current_label"].strip(),
                    predicted_label=row["predicted_label"].strip(),
                    decision=row["decision"].strip(),
                    new_label=row["new_label"].strip(),
                    reason=row["reason"].strip(),
                )
            )

    if not rows:
        raise SystemExit("review CSV has no review rows")
    return rows


def validate_review_rows(review_rows: list[ReviewRow], source_processed_dir: Path) -> None:
    errors: list[str] = []
    seen_paths: set[Path] = set()
    service_root = source_processed_dir.parent.parent

    for row in review_rows:
        if row.decision not in VALID_DECISIONS:
            errors.append(
                f"row {row.row_number}: invalid decision {row.decision!r}; "
                f"expected one of {sorted(VALID_DECISIONS)}"
            )

        if row.decision == "relabel" and not row.new_label:
            errors.append(f"row {row.row_number}: relabel decision requires new_label")
        if row.decision != "relabel" and row.new_label:
            errors.append(f"row {row.row_number}: new_label should only be set for relabel")

        source_path = service_root / row.image_path
        if not source_path.exists():
            errors.append(f"row {row.row_number}: image_path does not exist: {row.image_path}")
        elif source_path.suffix.lower() not in IMAGE_EXTENSIONS:
            errors.append(
                f"row {row.row_number}: image_path is not a supported image: {row.image_path}"
            )

        path_parts = row.image_path.parts
        if len(path_parts) < 4 or path_parts[0] != "data" or path_parts[1] != "processed":
            errors.append(
                f"row {row.row_number}: image_path must be relative like "
                "data/processed/<split>/<label>/<filename>"
            )
        elif path_parts[-2] != row.current_label:
            errors.append(
                f"row {row.row_number}: path label {path_parts[-2]!r} does not match "
                f"current_label {row.current_label!r}"
            )

        if row.image_path in seen_paths:
            errors.append(f"row {row.row_number}: duplicate review row for {row.image_path}")
        seen_paths.add(row.image_path)

    if errors:
        raise SystemExit("Review CSV validation failed:\n" + "\n".join(errors[:50]))


def plan_actions(
    review_rows: list[ReviewRow],
    source_processed_dir: Path,
    output_processed_dir: Path,
) -> list[AppliedAction]:
    actions: list[AppliedAction] = []
    service_root = source_processed_dir.parent.parent
    for row in review_rows:
        source_path = service_root / row.image_path
        relative_to_processed = source_path.relative_to(source_processed_dir)
        output_path = output_processed_dir / relative_to_processed

        if row.decision in REMOVE_DECISIONS:
            actions.append(
                AppliedAction(
                    row_number=row.row_number,
                    decision=row.decision,
                    source_path=str(source_path),
                    output_path=str(output_path),
                )
            )
        elif row.decision == "relabel":
            split = relative_to_processed.parts[0]
            new_output_path = output_processed_dir / split / row.new_label / source_path.name
            actions.append(
                AppliedAction(
                    row_number=row.row_number,
                    decision=row.decision,
                    source_path=str(source_path),
                    output_path=str(output_path),
                    new_output_path=str(new_output_path),
                )
            )

    return actions


def copy_processed_dataset(
    source_processed_dir: Path,
    output_processed_dir: Path,
    force: bool,
) -> None:
    if output_processed_dir.exists():
        if not force:
            raise SystemExit(
                f"output processed dir already exists: {output_processed_dir}; "
                "pass --force to replace it"
            )
        shutil.rmtree(output_processed_dir)

    shutil.copytree(source_processed_dir, output_processed_dir)


def apply_actions(actions: list[AppliedAction]) -> None:
    for action in actions:
        output_path = Path(action.output_path)
        if action.decision in REMOVE_DECISIONS:
            output_path.unlink(missing_ok=True)
            continue

        if action.decision == "relabel":
            if action.new_output_path is None:
                raise RuntimeError(f"missing new output path for row {action.row_number}")
            new_output_path = Path(action.new_output_path)
            new_output_path.parent.mkdir(parents=True, exist_ok=True)
            if new_output_path.exists():
                raise SystemExit(
                    f"cannot relabel row {action.row_number}; destination already exists: "
                    f"{new_output_path}"
                )
            output_path.replace(new_output_path)


def build_report(
    args: argparse.Namespace,
    review_rows: list[ReviewRow],
    actions: list[AppliedAction],
) -> dict[str, object]:
    decisions = Counter(row.decision for row in review_rows)
    labels = Counter(row.current_label for row in review_rows)
    predictions = Counter(row.predicted_label for row in review_rows)

    return {
        "generatedAt": datetime.now(UTC).isoformat(),
        "reviewCsv": str(args.review_csv),
        "sourceProcessedDir": str(args.source_processed_dir),
        "outputProcessedDir": str(args.output_processed_dir),
        "reviewedImages": len(review_rows),
        "decisionCounts": dict(sorted(decisions.items())),
        "currentLabelCounts": dict(sorted(labels.items())),
        "predictedLabelCounts": dict(sorted(predictions.items())),
        "appliedActionCounts": dict(sorted(Counter(action.decision for action in actions).items())),
        "appliedActions": [action.__dict__ for action in actions],
    }


def print_summary(report: dict[str, object], dry_run: bool) -> None:
    prefix = "DRY RUN - " if dry_run else ""
    print(f"{prefix}reviewed images: {report['reviewedImages']}")
    print(f"{prefix}decision counts:")
    for decision, count in report["decisionCounts"].items():
        print(f"  {decision}: {count}")
    print(f"{prefix}applied action counts:")
    for decision, count in report["appliedActionCounts"].items():
        print(f"  {decision}: {count}")
    if dry_run:
        print("No files were copied or changed.")


if __name__ == "__main__":
    sys.exit(main())
