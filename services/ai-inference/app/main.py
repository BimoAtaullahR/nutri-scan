from fastapi import Depends, FastAPI, File, Request, UploadFile

from app.classifier import FoodClassifier, get_classifier
from app.payload import build_inference_payload
from app.runtime import RuntimeConfig, validate_model_artifact

app = FastAPI(title="NutriScan AI/ML Inference", version="0.1.0")
CLASSIFIER_DEPENDENCY = Depends(get_classifier)
IMAGE_FILE = File(default=None)


@app.get("/healthz")
def healthz() -> dict[str, str]:
    return {"status": "ok"}


@app.get("/readyz")
def readyz() -> dict[str, object]:
    return validate_model_artifact(RuntimeConfig.from_env())


@app.post("/infer")
async def infer(
    request: Request,
    portion: str = "medium",
    image: UploadFile | None = IMAGE_FILE,
    classifier: FoodClassifier = CLASSIFIER_DEPENDENCY,
) -> dict[str, object]:
    image_bytes = await image.read() if image is not None else await request.body()
    predictions = classifier.predict(image_bytes)
    return build_inference_payload(
        predictions=predictions,
        portion=portion,
        model_version=classifier.model_version,
        confidence_threshold=classifier.confidence_threshold,
    )
