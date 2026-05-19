#!/usr/bin/env python3
from __future__ import annotations

import argparse
import json
from pathlib import Path

from PIL import Image

IMAGENET_MEAN = [0.485, 0.456, 0.406]
IMAGENET_STD = [0.229, 0.224, 0.225]


def build_transform(image_size: int):
    from torchvision import transforms

    return transforms.Compose(
        [
            transforms.Resize(int(image_size * 1.15)),
            transforms.CenterCrop(image_size),
            transforms.ToTensor(),
            transforms.Normalize(IMAGENET_MEAN, IMAGENET_STD),
        ]
    )


def load_model(model_path: Path, device: str):
    import timm
    import torch

    checkpoint = torch.load(model_path, map_location=device)
    model = timm.create_model(
        checkpoint["model_name"],
        pretrained=False,
        num_classes=int(checkpoint["num_classes"]),
    )
    model.load_state_dict(checkpoint["state_dict"])
    model.to(device)
    model.eval()
    return model, checkpoint


def predict_image(model_path: Path, image_path: Path) -> dict[str, object]:
    import torch

    device = "cuda" if torch.cuda.is_available() else "cpu"
    model, checkpoint = load_model(model_path, device)
    idx_to_class = checkpoint["idx_to_class"]
    image_size = int(checkpoint["image_size"])
    transform = build_transform(image_size)

    image = Image.open(image_path).convert("RGB")
    tensor = transform(image).unsqueeze(0).to(device)

    with torch.no_grad():
        probabilities = torch.softmax(model(tensor), dim=1)[0]
        top_k = min(3, probabilities.numel())
        top_probabilities, top_indices = probabilities.topk(top_k)

    top3 = [
        {
            "label": idx_to_class[str(int(index.item()))],
            "confidence": float(probability.item()),
        }
        for probability, index in zip(top_probabilities, top_indices, strict=True)
    ]
    return {
        "label": top3[0]["label"],
        "confidence": top3[0]["confidence"],
        "top3": top3,
    }


def main() -> None:
    parser = argparse.ArgumentParser(description="Predict one NutriScan food image")
    parser.add_argument("--model-path", type=Path, required=True)
    parser.add_argument("--image-path", type=Path, required=True)
    args = parser.parse_args()

    print(json.dumps(predict_image(args.model_path, args.image_path), indent=2))


if __name__ == "__main__":
    main()
