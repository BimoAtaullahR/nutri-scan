#!/usr/bin/env python3
from __future__ import annotations

import argparse
import contextlib
import json
import random
from dataclasses import asdict, dataclass
from pathlib import Path
from typing import Any

SPLITS = ("train", "validation", "test")
SPLIT_ALIASES = {"validation": ("validation", "val")}
IMAGENET_MEAN = [0.485, 0.456, 0.406]
IMAGENET_STD = [0.229, 0.224, 0.225]


@dataclass(frozen=True)
class TrainingConfig:
    model_name: str
    pretrained: bool
    num_classes: int
    class_names: list[str]
    image_size: int
    batch_size: int
    epochs: int
    learning_rate: float
    weight_decay: float
    optimizer: str
    scheduler: str | None
    early_stopping_patience: int | None
    seed: int
    use_amp: bool
    num_workers: int
    label_smoothing: float
    random_resized_crop_scale: tuple[float, float]
    horizontal_flip_p: float
    rotation_degrees: float
    color_jitter_brightness: float
    color_jitter_contrast: float
    color_jitter_saturation: float
    random_erasing_p: float
    output_dir: Path
    report_dir: Path


def load_json(path: Path) -> dict[str, Any]:
    return json.loads(path.read_text())


def resolve_path(value: str | Path, config_path: Path) -> Path:
    path = Path(value)
    if path.is_absolute():
        return path

    next_to_config = config_path.parent / path
    if next_to_config.exists():
        return next_to_config

    return path


def load_class_names(class_config_path: Path) -> list[str]:
    raw = load_json(class_config_path)
    if "classes" in raw:
        class_names = list(raw["classes"])
    else:
        class_names = [category["slug"] for category in raw["categories"]]

    if len(class_names) != len(set(class_names)):
        raise ValueError("Food classes must be unique")
    if not class_names:
        raise ValueError("Food class list cannot be empty")

    return class_names


def load_crop_scale(raw: dict[str, Any]) -> tuple[float, float]:
    value = raw.get("random_resized_crop_scale", raw.get("randomResizedCropScale", [0.75, 1.0]))
    if not isinstance(value, list | tuple) or len(value) != 2:
        raise ValueError("random_resized_crop_scale must contain two numeric values")

    minimum = float(value[0])
    maximum = float(value[1])
    if minimum <= 0 or maximum <= 0 or minimum > maximum or maximum > 1.0:
        raise ValueError("random_resized_crop_scale must satisfy 0 < min <= max <= 1")
    return minimum, maximum


def load_config(path: Path) -> TrainingConfig:
    raw = load_json(path)

    # Backward compatibility for the earlier camelCase baseline config.
    if "model_name" not in raw and "modelName" in raw:
        class_names = list(raw["labels"])
        output_dir = Path(raw["artifactDir"])
        return TrainingConfig(
            model_name=raw["modelName"],
            pretrained=True,
            num_classes=len(class_names),
            class_names=class_names,
            image_size=int(raw["imageSize"]),
            batch_size=int(raw["batchSize"]),
            epochs=int(raw["epochs"]),
            learning_rate=float(raw["learningRate"]),
            weight_decay=float(raw.get("weightDecay", 0.0001)),
            optimizer=str(raw.get("optimizer", "adamw")).lower(),
            scheduler=raw.get("scheduler"),
            early_stopping_patience=raw.get("earlyStoppingPatience"),
            seed=int(raw.get("seed", 42)),
            use_amp=bool(raw.get("useAmp", True)),
            num_workers=int(raw.get("numWorkers", 0)),
            label_smoothing=float(raw.get("labelSmoothing", raw.get("label_smoothing", 0.0))),
            random_resized_crop_scale=load_crop_scale(raw),
            horizontal_flip_p=float(raw.get("horizontalFlipP", raw.get("horizontal_flip_p", 0.5))),
            rotation_degrees=float(raw.get("rotationDegrees", raw.get("rotation_degrees", 10))),
            color_jitter_brightness=float(
                raw.get("colorJitterBrightness", raw.get("color_jitter_brightness", 0.15))
            ),
            color_jitter_contrast=float(
                raw.get("colorJitterContrast", raw.get("color_jitter_contrast", 0.15))
            ),
            color_jitter_saturation=float(
                raw.get("colorJitterSaturation", raw.get("color_jitter_saturation", 0.10))
            ),
            random_erasing_p=float(raw.get("randomErasingP", raw.get("random_erasing_p", 0.0))),
            output_dir=output_dir,
            report_dir=Path(raw.get("reportDir", "reports/baseline-food-classifier")),
        )

    class_config = resolve_path(raw.get("class_config", "configs/mvp_food_categories.json"), path)
    class_names = load_class_names(class_config)
    num_classes = int(raw.get("num_classes", len(class_names)))
    if num_classes != len(class_names):
        raise ValueError(
            f"num_classes={num_classes} does not match configured classes={len(class_names)}"
        )

    return TrainingConfig(
        model_name=str(raw["model_name"]),
        pretrained=bool(raw.get("pretrained", True)),
        num_classes=num_classes,
        class_names=class_names,
        image_size=int(raw["image_size"]),
        batch_size=int(raw["batch_size"]),
        epochs=int(raw["epochs"]),
        learning_rate=float(raw["learning_rate"]),
        weight_decay=float(raw.get("weight_decay", 0.0001)),
        optimizer=str(raw.get("optimizer", "adamw")).lower(),
        scheduler=raw.get("scheduler"),
        early_stopping_patience=raw.get("early_stopping_patience"),
        seed=int(raw.get("seed", 42)),
        use_amp=bool(raw.get("use_amp", True)),
        num_workers=int(raw.get("num_workers", 0)),
        label_smoothing=float(raw.get("label_smoothing", 0.0)),
        random_resized_crop_scale=load_crop_scale(raw),
        horizontal_flip_p=float(raw.get("horizontal_flip_p", 0.5)),
        rotation_degrees=float(raw.get("rotation_degrees", 10)),
        color_jitter_brightness=float(raw.get("color_jitter_brightness", 0.15)),
        color_jitter_contrast=float(raw.get("color_jitter_contrast", 0.15)),
        color_jitter_saturation=float(raw.get("color_jitter_saturation", 0.10)),
        random_erasing_p=float(raw.get("random_erasing_p", 0.0)),
        output_dir=Path(raw["output_dir"]),
        report_dir=Path(raw["report_dir"]),
    )


def set_seed(seed: int) -> None:
    random.seed(seed)

    try:
        import numpy as np

        np.random.seed(seed)
    except ImportError:
        pass

    import torch

    torch.manual_seed(seed)
    if torch.cuda.is_available():
        torch.cuda.manual_seed_all(seed)
    torch.backends.cudnn.benchmark = False
    torch.backends.cudnn.deterministic = True


def resolve_split_dirs(processed_dir: Path) -> dict[str, Path]:
    split_dirs: dict[str, Path] = {}
    for split in SPLITS:
        candidates = SPLIT_ALIASES.get(split, (split,))
        for candidate in candidates:
            path = processed_dir / candidate
            if path.is_dir():
                split_dirs[split] = path
                if candidate != split:
                    print(f"split_alias={split}:{candidate}")
                break
        else:
            split_dirs[split] = processed_dir / split

    return split_dirs


def validate_processed_dataset(processed_dir: Path, class_names: list[str]) -> dict[str, Path]:
    split_dirs = resolve_split_dirs(processed_dir)
    missing: list[str] = []
    for split in SPLITS:
        split_dir = split_dirs[split]
        if not split_dir.is_dir():
            missing.append(split)
            continue
        for class_name in class_names:
            if not (split_dir / class_name).is_dir():
                missing.append(f"{split}/{class_name}")

    if missing:
        joined = ", ".join(missing[:12])
        suffix = "" if len(missing) <= 12 else f", and {len(missing) - 12} more"
        raise FileNotFoundError(f"Missing processed dataset folders: {joined}{suffix}")

    return split_dirs


def build_transforms(config: TrainingConfig):
    from torchvision import transforms

    train_steps = [
        transforms.RandomResizedCrop(config.image_size, scale=config.random_resized_crop_scale),
        transforms.RandomHorizontalFlip(p=config.horizontal_flip_p),
        transforms.RandomRotation(degrees=config.rotation_degrees),
        transforms.ColorJitter(
            brightness=config.color_jitter_brightness,
            contrast=config.color_jitter_contrast,
            saturation=config.color_jitter_saturation,
        ),
        transforms.ToTensor(),
        transforms.Normalize(IMAGENET_MEAN, IMAGENET_STD),
    ]
    if config.random_erasing_p > 0:
        train_steps.append(transforms.RandomErasing(p=config.random_erasing_p))

    train_transform = transforms.Compose(train_steps)
    eval_transform = transforms.Compose(
        [
            transforms.Resize(int(config.image_size * 1.15)),
            transforms.CenterCrop(config.image_size),
            transforms.ToTensor(),
            transforms.Normalize(IMAGENET_MEAN, IMAGENET_STD),
        ]
    )
    return train_transform, eval_transform


def build_datasets(split_dirs: dict[str, Path], config: TrainingConfig):
    from torchvision import datasets

    train_transform, eval_transform = build_transforms(config)
    split_transforms = {
        "train": train_transform,
        "validation": eval_transform,
        "test": eval_transform,
    }
    image_datasets = {
        split: datasets.ImageFolder(split_dirs[split], transform=split_transforms[split])
        for split in SPLITS
    }

    expected = config.class_names
    for split, dataset in image_datasets.items():
        if dataset.classes != expected:
            raise ValueError(
                f"{split} classes {dataset.classes} do not match "
                f"configured class order {expected}. "
                "Use class folders that match configs/mvp_food_categories.json exactly."
            )

    return image_datasets


def build_dataloaders(image_datasets, config: TrainingConfig, device: str):
    import torch
    from torch.utils.data import DataLoader

    generator = torch.Generator()
    generator.manual_seed(config.seed)
    pin_memory = device == "cuda"

    return {
        "train": DataLoader(
            image_datasets["train"],
            batch_size=config.batch_size,
            shuffle=True,
            num_workers=config.num_workers,
            pin_memory=pin_memory,
            generator=generator,
        ),
        "validation": DataLoader(
            image_datasets["validation"],
            batch_size=config.batch_size,
            shuffle=False,
            num_workers=config.num_workers,
            pin_memory=pin_memory,
        ),
        "test": DataLoader(
            image_datasets["test"],
            batch_size=config.batch_size,
            shuffle=False,
            num_workers=config.num_workers,
            pin_memory=pin_memory,
        ),
    }


def detect_device() -> str:
    import torch

    device = "cuda" if torch.cuda.is_available() else "cpu"
    print(f"torch={torch.__version__}")
    print(f"device={device}")
    if device == "cuda":
        print(f"cuda={torch.version.cuda}")
        print(f"gpu={torch.cuda.get_device_name(0)}")
    return device


def create_model(config: TrainingConfig):
    import timm

    return timm.create_model(
        config.model_name,
        pretrained=config.pretrained,
        num_classes=config.num_classes,
    )


def create_optimizer(config: TrainingConfig, model):
    import torch

    if config.optimizer != "adamw":
        raise ValueError(f"Unsupported optimizer: {config.optimizer}")
    return torch.optim.AdamW(
        model.parameters(),
        lr=config.learning_rate,
        weight_decay=config.weight_decay,
    )


def create_scheduler(config: TrainingConfig, optimizer):
    import torch

    if config.scheduler in (None, "", "none"):
        return None
    if str(config.scheduler).lower() == "cosine":
        return torch.optim.lr_scheduler.CosineAnnealingLR(optimizer, T_max=config.epochs)
    raise ValueError(f"Unsupported scheduler: {config.scheduler}")


def create_criterion(config: TrainingConfig):
    from torch import nn

    return nn.CrossEntropyLoss(label_smoothing=config.label_smoothing)


def autocast_context(use_amp: bool):
    if not use_amp:
        return contextlib.nullcontext()

    import torch

    return torch.amp.autocast(device_type="cuda", enabled=True)


def create_grad_scaler(use_amp: bool):
    import torch

    try:
        return torch.amp.GradScaler("cuda", enabled=use_amp)
    except TypeError:
        return torch.cuda.amp.GradScaler(enabled=use_amp)


def train_one_epoch(
    model,
    loader,
    criterion,
    optimizer,
    scaler,
    device: str,
    use_amp: bool,
) -> float:
    model.train()
    total_loss = 0.0
    total_examples = 0

    for images, targets in loader:
        images = images.to(device, non_blocking=True)
        targets = targets.to(device, non_blocking=True)

        optimizer.zero_grad(set_to_none=True)
        with autocast_context(use_amp):
            logits = model(images)
            loss = criterion(logits, targets)

        scaler.scale(loss).backward()
        scaler.step(optimizer)
        scaler.update()

        batch_size = targets.size(0)
        total_loss += float(loss.detach().item()) * batch_size
        total_examples += batch_size

    return total_loss / total_examples if total_examples else 0.0


def evaluate_loader(model, loader, criterion, device: str) -> dict[str, float]:
    import torch

    model.eval()
    total_loss = 0.0
    total_examples = 0
    top1_correct = 0
    top3_correct = 0

    with torch.no_grad():
        for images, targets in loader:
            images = images.to(device, non_blocking=True)
            targets = targets.to(device, non_blocking=True)
            logits = model(images)
            loss = criterion(logits, targets)
            top_k = min(3, logits.size(1))
            _, top_indices = logits.topk(top_k, dim=1)

            batch_size = targets.size(0)
            total_loss += float(loss.detach().item()) * batch_size
            total_examples += batch_size
            top1_correct += int((top_indices[:, 0] == targets).sum().item())
            top3_correct += int((top_indices == targets.unsqueeze(1)).any(dim=1).sum().item())

    return {
        "loss": total_loss / total_examples if total_examples else 0.0,
        "top1": top1_correct / total_examples if total_examples else 0.0,
        "top3": top3_correct / total_examples if total_examples else 0.0,
    }


def write_label_map(output_dir: Path, class_to_idx: dict[str, int]) -> None:
    idx_to_class = {str(index): label for label, index in class_to_idx.items()}
    label_map = {
        "class_to_idx": class_to_idx,
        "idx_to_class": idx_to_class,
        "idToLabel": idx_to_class,
    }
    (output_dir / "label_map.json").write_text(json.dumps(label_map, indent=2) + "\n")


def checkpoint_payload(
    config: TrainingConfig,
    class_to_idx: dict[str, int],
    model,
) -> dict[str, Any]:
    idx_to_class = {str(index): label for label, index in class_to_idx.items()}
    return {
        "model_name": config.model_name,
        "num_classes": config.num_classes,
        "image_size": config.image_size,
        "class_to_idx": class_to_idx,
        "idx_to_class": idx_to_class,
        "state_dict": model.state_dict(),
    }


def resolved_config_payload(
    config: TrainingConfig,
    device: str,
    class_to_idx: dict[str, int],
) -> dict[str, Any]:
    payload = asdict(config)
    payload["output_dir"] = str(config.output_dir)
    payload["report_dir"] = str(config.report_dir)
    payload["device"] = device
    payload["class_to_idx"] = class_to_idx
    return payload


def generate_predictions(
    model,
    dataset,
    loader,
    class_names: list[str],
    device: str,
) -> dict[str, Any]:
    import torch

    model.eval()
    predictions: list[dict[str, Any]] = []
    sample_index = 0

    with torch.no_grad():
        for images, targets in loader:
            images = images.to(device, non_blocking=True)
            logits = model(images)
            probabilities = torch.softmax(logits, dim=1)
            top_k = min(3, probabilities.size(1))
            top_probabilities, top_indices = probabilities.topk(top_k, dim=1)

            for row in range(images.size(0)):
                path, _ = dataset.samples[sample_index]
                true_index = int(targets[row].item())
                top3 = [
                    {
                        "label": class_names[int(top_indices[row, col].item())],
                        "confidence": float(top_probabilities[row, col].item()),
                    }
                    for col in range(top_k)
                ]
                predictions.append(
                    {
                        "path": str(path),
                        "true_label": class_names[true_index],
                        "top1_label": top3[0]["label"],
                        "top1_confidence": top3[0]["confidence"],
                        "top3": top3,
                    }
                )
                sample_index += 1

    return {"classes": class_names, "predictions": predictions}


def train_classifier(config: TrainingConfig, processed_dir: Path) -> None:
    split_dirs = validate_processed_dataset(processed_dir, config.class_names)

    import torch
    set_seed(config.seed)
    device = detect_device()
    use_amp = config.use_amp and device == "cuda"
    print(f"amp={str(use_amp).lower()}")

    image_datasets = build_datasets(split_dirs, config)
    dataloaders = build_dataloaders(image_datasets, config, device)
    model = create_model(config).to(device)
    optimizer = create_optimizer(config, model)
    scheduler = create_scheduler(config, optimizer)
    criterion = create_criterion(config)
    scaler = create_grad_scaler(use_amp)

    config.output_dir.mkdir(parents=True, exist_ok=True)
    config.report_dir.mkdir(parents=True, exist_ok=True)

    class_to_idx = image_datasets["train"].class_to_idx
    write_label_map(config.output_dir, class_to_idx)
    (config.output_dir / "training_config_resolved.json").write_text(
        json.dumps(resolved_config_payload(config, device, class_to_idx), indent=2) + "\n"
    )

    best_validation_top1 = -1.0
    best_epoch = 0
    epochs_without_improvement = 0
    model_path = config.output_dir / "model.pt"

    for epoch in range(1, config.epochs + 1):
        train_loss = train_one_epoch(
            model,
            dataloaders["train"],
            criterion,
            optimizer,
            scaler,
            device,
            use_amp,
        )
        validation_metrics = evaluate_loader(
            model,
            dataloaders["validation"],
            criterion,
            device,
        )
        if scheduler is not None:
            scheduler.step()

        validation_top1 = validation_metrics["top1"]
        improved = validation_top1 > best_validation_top1
        if improved:
            best_validation_top1 = validation_top1
            best_epoch = epoch
            epochs_without_improvement = 0
            torch.save(checkpoint_payload(config, class_to_idx, model), model_path)
        else:
            epochs_without_improvement += 1

        print(
            f"epoch={epoch} train_loss={train_loss:.4f} "
            f"validation_loss={validation_metrics['loss']:.4f} "
            f"validation_top1={validation_metrics['top1']:.4f} "
            f"validation_top3={validation_metrics['top3']:.4f} "
            f"best_epoch={best_epoch}"
        )

        patience = config.early_stopping_patience
        if patience is not None and epochs_without_improvement >= patience:
            print(f"early_stopping=true patience={patience}")
            break

    checkpoint = torch.load(model_path, map_location=device)
    model.load_state_dict(checkpoint["state_dict"])
    predictions = generate_predictions(
        model,
        image_datasets["test"],
        dataloaders["test"],
        config.class_names,
        device,
    )
    (config.report_dir / "predictions.json").write_text(json.dumps(predictions, indent=2) + "\n")
    print(f"saved_model={model_path}")
    print(f"saved_predictions={config.report_dir / 'predictions.json'}")


def main() -> None:
    parser = argparse.ArgumentParser(description="Train NutriScan MVP food classifier")
    parser.add_argument("--config", type=Path, required=True)
    parser.add_argument("--processed-dir", type=Path, required=True)
    parser.add_argument("--dry-run", action="store_true")
    args = parser.parse_args()

    config = load_config(args.config)
    if args.dry_run:
        print(
            f"model={config.model_name} classes={len(config.class_names)} "
            f"imageSize={config.image_size} batchSize={config.batch_size} "
            f"epochs={config.epochs} labelSmoothing={config.label_smoothing:g} "
            f"cropScale={config.random_resized_crop_scale[0]:g}-{config.random_resized_crop_scale[1]:g} "
            f"rotation={config.rotation_degrees:g} "
            f"colorJitter={config.color_jitter_brightness:g}/{config.color_jitter_contrast:g}/{config.color_jitter_saturation:g} "
            f"randomErasing={config.random_erasing_p:g} "
            f"artifactDir={config.output_dir.as_posix()} "
            f"reportDir={config.report_dir.as_posix()}"
        )
        return

    train_classifier(config, args.processed_dir)


if __name__ == "__main__":
    main()
