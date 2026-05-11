#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
from dataclasses import dataclass
from pathlib import Path


@dataclass(frozen=True)
class TrainingConfig:
    model_name: str
    image_size: int
    batch_size: int
    epochs: int
    learning_rate: float
    labels: list[str]
    artifact_dir: Path


def load_config(path: Path) -> TrainingConfig:
    raw = json.loads(path.read_text())
    labels = raw["labels"]
    if len(labels) != len(set(labels)):
        raise ValueError("Training labels must be unique")

    return TrainingConfig(
        model_name=raw["modelName"],
        image_size=int(raw["imageSize"]),
        batch_size=int(raw["batchSize"]),
        epochs=int(raw["epochs"]),
        learning_rate=float(raw["learningRate"]),
        labels=list(labels),
        artifact_dir=Path(raw["artifactDir"]),
    )


def validate_processed_dataset(processed_dir: Path, labels: list[str]) -> None:
    missing: list[str] = []
    for split in ("train", "validation"):
        for label in labels:
            if not (processed_dir / split / label).is_dir():
                missing.append(f"{split}/{label}")
    if missing:
        joined = ", ".join(missing[:8])
        suffix = "" if len(missing) <= 8 else f", and {len(missing) - 8} more"
        raise FileNotFoundError(f"Missing processed dataset folders: {joined}{suffix}")


def write_label_map(artifact_dir: Path, labels: list[str]) -> None:
    artifact_dir.mkdir(parents=True, exist_ok=True)
    label_map = {"idToLabel": {str(index): label for index, label in enumerate(labels)}}
    (artifact_dir / "label_map.json").write_text(json.dumps(label_map, indent=2) + "\n")


def train_classifier(config: TrainingConfig, processed_dir: Path) -> None:
    validate_processed_dataset(processed_dir, config.labels)

    import timm
    import torch
    from torch import nn
    from torch.utils.data import DataLoader
    from torchvision import datasets, transforms

    train_transform = transforms.Compose(
        [
            transforms.Resize((config.image_size, config.image_size)),
            transforms.RandomHorizontalFlip(),
            transforms.ToTensor(),
        ]
    )
    eval_transform = transforms.Compose(
        [
            transforms.Resize((config.image_size, config.image_size)),
            transforms.ToTensor(),
        ]
    )

    train_data = datasets.ImageFolder(processed_dir / "train", transform=train_transform)
    validation_data = datasets.ImageFolder(processed_dir / "validation", transform=eval_transform)
    if train_data.classes != config.labels:
        raise ValueError(f"Dataset classes {train_data.classes} do not match config labels")

    train_loader = DataLoader(train_data, batch_size=config.batch_size, shuffle=True)
    validation_loader = DataLoader(validation_data, batch_size=config.batch_size)

    device = "cuda" if torch.cuda.is_available() else "cpu"
    model = timm.create_model(config.model_name, pretrained=True, num_classes=len(config.labels))
    model.to(device)
    optimizer = torch.optim.AdamW(model.parameters(), lr=config.learning_rate)
    criterion = nn.CrossEntropyLoss()

    for epoch in range(config.epochs):
        model.train()
        for images, targets in train_loader:
            images = images.to(device)
            targets = targets.to(device)
            optimizer.zero_grad(set_to_none=True)
            loss = criterion(model(images), targets)
            loss.backward()
            optimizer.step()

        model.eval()
        correct = 0
        total = 0
        with torch.no_grad():
            for images, targets in validation_loader:
                images = images.to(device)
                targets = targets.to(device)
                predictions = model(images).argmax(dim=1)
                correct += (predictions == targets).sum().item()
                total += targets.numel()
        accuracy = correct / total if total else 0
        print(f"epoch={epoch + 1} validationTop1={accuracy:.4f}")

    config.artifact_dir.mkdir(parents=True, exist_ok=True)
    torch.save(model.state_dict(), config.artifact_dir / "model.pt")
    write_label_map(config.artifact_dir, config.labels)


def main() -> None:
    parser = argparse.ArgumentParser(description="Train NutriScan MVP food classifier")
    parser.add_argument("--config", type=Path, required=True)
    parser.add_argument("--processed-dir", type=Path, required=True)
    parser.add_argument("--dry-run", action="store_true")
    args = parser.parse_args()

    config = load_config(args.config)
    if args.dry_run:
        print(
            f"model={config.model_name} classes={len(config.labels)} "
            f"imageSize={config.image_size} batchSize={config.batch_size} "
            f"epochs={config.epochs} artifactDir={config.artifact_dir}"
        )
        return

    train_classifier(config, args.processed_dir)


if __name__ == "__main__":
    main()
