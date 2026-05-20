from __future__ import annotations

import json
from io import BytesIO
from pathlib import Path

from PIL import Image

from app.payload import Prediction
from app.runtime import RuntimeConfig, validate_model_artifact

IMAGENET_MEAN = [0.485, 0.456, 0.406]
IMAGENET_STD = [0.229, 0.224, 0.225]


class FoodClassifier:
    def __init__(
        self,
        config: RuntimeConfig | None = None,
        artifact_dir: Path | None = None,
    ) -> None:
        if artifact_dir is not None:
            base_config = config or RuntimeConfig.from_env()
            config = RuntimeConfig(
                artifact_dir=artifact_dir,
                model_version=base_config.model_version,
                confidence_threshold=base_config.confidence_threshold,
            )
        self.config = config or RuntimeConfig.from_env()
        self.artifact_dir = self.config.artifact_dir
        self.model_path = self.artifact_dir / "model.pt"
        self.label_map_path = self.artifact_dir / "label_map.json"
        self.model_version = self.config.model_version
        self.confidence_threshold = self.config.confidence_threshold
        self.labels = self._load_labels()
        self._model = None
        self._checkpoint: dict[str, object] | None = None
        self._device: str | None = None
        self._transform = None

    def preprocessing_recipe(self) -> dict[str, object]:
        metadata = validate_model_artifact(self.config)
        image_size = int(metadata["imageSize"])
        return {
            "imageSize": image_size,
            "resizeSize": int(image_size * 1.15),
            "normalizationMean": IMAGENET_MEAN,
            "normalizationStd": IMAGENET_STD,
        }

    def _load_labels(self) -> list[str]:
        if self.label_map_path.exists():
            raw = json.loads(self.label_map_path.read_text())
            return [raw["idToLabel"][str(index)] for index in range(len(raw["idToLabel"]))]
        return []

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
        preprocessing = self.preprocessing_recipe()
        image_size = int(preprocessing["imageSize"])
        resize_size = int(preprocessing["resizeSize"])
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
                transforms.Resize(resize_size),
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

        validate_model_artifact(self.config)
        return self._predict_with_model(image_bytes)


_classifier = FoodClassifier()


def get_classifier() -> FoodClassifier:
    return _classifier
