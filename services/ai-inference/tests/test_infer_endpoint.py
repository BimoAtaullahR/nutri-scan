import asyncio
from io import BytesIO

from PIL import Image

from app.classifier import FoodClassifier
from app.main import app, infer


def image_bytes() -> bytes:
    buffer = BytesIO()
    Image.new("RGB", (32, 32), color=(255, 0, 0)).save(buffer, format="JPEG")
    return buffer.getvalue()


def test_infer_endpoint_accepts_image_and_returns_payload() -> None:
    class FakeRequest:
        async def body(self) -> bytes:
            return image_bytes()

    route_paths = {route.path for route in app.routes}
    assert "/infer" in route_paths

    payload = asyncio.run(
        infer(
            request=FakeRequest(),
            portion="medium",
            classifier=FoodClassifier(),
        )
    )

    assert payload["modelVersion"]
    assert payload["foodCategory"]["slug"]
    assert payload["coarsePortion"] == "medium"
    assert "isLowConfidence" in payload
