#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
import random
import shutil
import sys
from dataclasses import dataclass
from pathlib import Path


IMAGE_EXTENSIONS = {".jpg", ".jpeg", ".png", ".webp"}
SPLITS = ("train", "validation", "test")


@dataclass(frozen=True)
class Category:
    slug: str
    source_labels: tuple[str, ...]


@dataclass(frozen=True)
class SplitPlan:
    category: str
    train: tuple[Path, ...]
    validation: tuple[Path, ...]
    test: tuple[Path, ...]


def main(argv: list[str] | None = None) -> int:
    args = parse_args(argv)

    categories = load_categories(args.class_map)
    source_index = index_source_folders(args.raw_dir)
    plans = build_split_plans(categories, source_index, seed=args.seed)

    if not args.dry_run:
        write_splits(plans, args.processed_dir)

    for plan in plans:
        print(
            f"{plan.category}: train={len(plan.train)} "
            f"validation={len(plan.validation)} test={len(plan.test)}"
        )

    return 0


def parse_args(argv: list[str] | None) -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Prepare NutriScan MVP food image dataset splits."
    )
    parser.add_argument("--raw-dir", type=Path, required=True)
    parser.add_argument("--processed-dir", type=Path, required=True)
    parser.add_argument("--class-map", type=Path, required=True)
    parser.add_argument("--seed", type=int, default=42)
    parser.add_argument("--dry-run", action="store_true")
    return parser.parse_args(argv)


def load_categories(class_map_path: Path) -> list[Category]:
    class_map = json.loads(class_map_path.read_text())

    return [
        Category(
            slug=category["slug"],
            source_labels=tuple(normalize_label(label) for label in category["sourceLabels"]),
        )
        for category in class_map["categories"]
    ]


def index_source_folders(raw_dir: Path) -> dict[str, Path]:
    if not raw_dir.exists():
        raise SystemExit(f"raw dir does not exist: {raw_dir}")

    return {
        normalize_label(path.name): path
        for path in raw_dir.iterdir()
        if path.is_dir()
    }


def build_split_plans(
    categories: list[Category],
    source_index: dict[str, Path],
    seed: int,
) -> list[SplitPlan]:
    plans: list[SplitPlan] = []

    for category in categories:
        images = collect_images(category, source_index)
        if not images:
            continue

        rng = random.Random(f"{seed}:{category.slug}")
        shuffled = list(images)
        rng.shuffle(shuffled)

        train_end = int(len(shuffled) * 0.70)
        validation_end = train_end + int(len(shuffled) * 0.15)

        plans.append(
            SplitPlan(
                category=category.slug,
                train=tuple(shuffled[:train_end]),
                validation=tuple(shuffled[train_end:validation_end]),
                test=tuple(shuffled[validation_end:]),
            )
        )

    return plans


def collect_images(category: Category, source_index: dict[str, Path]) -> tuple[Path, ...]:
    images: set[Path] = set()

    for source_label in category.source_labels:
        source_dir = source_index.get(source_label)
        if source_dir is None:
            continue

        images.update(
            sorted(
                path
                for path in source_dir.iterdir()
                if path.is_file() and path.suffix.lower() in IMAGE_EXTENSIONS
            )
        )

    return tuple(sorted(images))


def write_splits(plans: list[SplitPlan], processed_dir: Path) -> None:
    for plan in plans:
        for split in SPLITS:
            split_dir = processed_dir / split / plan.category
            split_dir.mkdir(parents=True, exist_ok=True)

            for image_path in getattr(plan, split):
                shutil.copy2(image_path, split_dir / image_path.name)


def normalize_label(label: str) -> str:
    normalized = label.strip().lower()
    for char in (" ", "-"):
        normalized = normalized.replace(char, "_")
    while "__" in normalized:
        normalized = normalized.replace("__", "_")
    return normalized


if __name__ == "__main__":
    sys.exit(main())
