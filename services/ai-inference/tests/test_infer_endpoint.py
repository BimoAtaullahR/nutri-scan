import asyncio

from app.classifier import FoodClassifier
from app.main import app, infer


def test_infer_endpoint_accepts_image_and_returns_payload() -> None:
    class FakeRequest:
        async def body(self) -> bytes:
            return b"not-a-real-image"

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
