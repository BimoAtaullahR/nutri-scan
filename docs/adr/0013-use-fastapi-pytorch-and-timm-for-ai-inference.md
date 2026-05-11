# Use FastAPI, PyTorch, and timm for AI Inference

NutriScan will use FastAPI for the AI/ML Inference HTTP service, with PyTorch and timm for MVP food classifier training and inference. This keeps the service simple to call from the Backend API while using common computer vision tooling for fine-tuning lightweight pretrained models such as EfficientNet-B0 or MobileNetV3.

## Considered Options

- Use FastAPI with PyTorch and timm.
- Use TensorFlow/Keras for training and serving.
- Use a notebook-only model workflow without an HTTP inference service.

## Consequences

- The AI/ML service owns Python runtime dependencies for model loading, preprocessing, and inference.
- The MVP should prefer lightweight pretrained models over large models such as EfficientNetV2-L.
- Backend integration depends on a stable `/infer` HTTP contract rather than notebook outputs.
