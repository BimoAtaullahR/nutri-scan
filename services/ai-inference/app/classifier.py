from __future__ import annotations

import json
from io import BytesIO
from pathlib import Path

from PIL import Image

from app.payload import Prediction

DEFAULT_ARTIFACT_DIR = Path("model-artifacts/baseline-food-classifier")
IMAGENET_MEAN = [0.485, 0.456, 0.406]
IMAGENET_STD = [0.229, 0.224, 0.225]


class FoodClassifier:
    def __init__(self, artifact_dir: Path = DEFAULT_ARTIFACT_DIR) -> None:
        self.artifact_dir = artifact_dir
        self.model_path = artifact_dir / "model.pt"
        self.label_map_path = artifact_dir / "label_map.json"
        self.model_version = "baseline-0.1.0"
        self.labels = self._load_labels()
        self._model = None
        self._checkpoint: dict[str, object] | None = None
        self._device: str | None = None
        self._transform = None

    def _load_labels(self) -> list[str]:
        if self.label_map_path.exists():
            raw = json.loads(self.label_map_path.read_text())
            return [raw["idToLabel"][str(index)] for index in range(len(raw["idToLabel"]))]
        return ["sate", "rendang", "bakso"]

    def _load_model(self):
        if self._model is not None:
            return self._model

        import timm
        import torch
        from torchvision import transforms

        checkpoint = torch.load(self.model_path, map_location="cpu")
        if "state_dict" not in checkpoint:
            raise ValueError("Model checkpoint must include a state_dict")

        model_name = str(checkpoint.get("model_name", "efficientnet_b0"))
        num_classes = int(checkpoint.get("num_classes", len(self.labels)))
        image_size = int(checkpoint.get("image_size", 224))
        idx_to_class = checkpoint.get("idx_to_class")
        if isinstance(idx_to_class, dict):
            self.labels = [str(idx_to_class[str(index)]) for index in range(len(idx_to_class))]

        device = "cuda" if torch.cuda.is_available() else "cpu"
        model = timm.create_model(model_name, pretrained=False, num_classes=num_classes)
        model.load_state_dict(checkpoint["state_dict"])
        model.to(device)
        model.eval()

        self._checkpoint = checkpoint
        self._device = device
        self._transform = transforms.Compose(
            [
                transforms.Resize(int(image_size * 1.15)),
                transforms.CenterCrop(image_size),
                transforms.ToTensor(),
                transforms.Normalize(IMAGENET_MEAN, IMAGENET_STD),
            ]
        )
        self._model = model
        return model

    def _predict_with_model(self, image_bytes: bytes) -> list[Prediction]:
        import torch

        model = self._load_model()
        assert self._device is not None
        assert self._transform is not None

        image = Image.open(BytesIO(image_bytes)).convert("RGB")
        tensor = self._transform(image).unsqueeze(0).to(self._device)

        with torch.no_grad():
            probabilities = torch.softmax(model(tensor), dim=1)[0]
            top_k = min(3, probabilities.numel())
            top_probabilities, top_indices = probabilities.topk(top_k)

        return [
            {
                "label": self.labels[int(index.item())],
                "confidenceScore": float(probability.item()),
            }
            for probability, index in zip(top_probabilities, top_indices, strict=True)
        ]

    def predict(self, image_bytes: bytes) -> list[Prediction]:
        if not image_bytes:
            raise ValueError("Image upload cannot be empty")

        if not self.model_path.exists():
            return [
                {"label": self.labels[0], "confidenceScore": 0.61},
                {"label": self.labels[1], "confidenceScore": 0.24},
                {"label": self.labels[2], "confidenceScore": 0.15},
            ]

        return self._predict_with_model(image_bytes)


_classifier = FoodClassifier()


def get_classifier() -> FoodClassifier:
    return _classifier
