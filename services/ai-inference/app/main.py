from fastapi import Depends, FastAPI, Request

from app.classifier import FoodClassifier, get_classifier
from app.payload import build_inference_payload

app = FastAPI(title="NutriScan AI/ML Inference", version="0.1.0")


@app.get("/healthz")
def healthz() -> dict[str, str]:
    return {"status": "ok"}


@app.post("/infer")
async def infer(
    request: Request,
    portion: str = "medium",
    classifier: FoodClassifier = Depends(get_classifier),
) -> dict[str, object]:
    image_bytes = await request.body()
    predictions = classifier.predict(image_bytes)
    return build_inference_payload(
        predictions=predictions,
        portion=portion,
        model_version=classifier.model_version,
    )
