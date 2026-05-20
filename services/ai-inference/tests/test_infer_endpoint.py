import asyncio

from app.main import app, infer, readyz


def test_infer_endpoint_accepts_image_and_returns_payload() -> None:
    class FakeRequest:
        async def body(self) -> bytes:
            return b"not-a-real-image"

    class FakeClassifier:
        model_version = "test-model"
        confidence_threshold = 0.72

        def predict(self, image_bytes: bytes) -> list[dict[str, object]]:
            assert image_bytes == b"not-a-real-image"
            return [
                {"label": "sate", "confidenceScore": 0.81},
                {"label": "rendang", "confidenceScore": 0.11},
            ]

    route_paths = {route.path for route in app.routes}
    assert "/infer" in route_paths

    payload = asyncio.run(
        infer(
            request=FakeRequest(),
            portion="medium",
            classifier=FakeClassifier(),
        )
    )

    assert payload["modelVersion"] == "test-model"
    assert payload["foodCategory"]["slug"]
    assert payload["coarsePortion"] == "medium"
    assert "isLowConfidence" in payload
    assert payload["confidenceThreshold"] == 0.72


def test_readyz_endpoint_reports_selected_model_metadata(monkeypatch, tmp_path) -> None:
    artifact_dir = tmp_path / "selected-mvp-classifier"
    artifact_dir.mkdir()
    (artifact_dir / "model.pt").write_bytes(b"runtime model weights live outside git")
    (artifact_dir / "label_map.json").write_text(
        """
        {
          "idToLabel": {
            "0": "bakso",
            "1": "gado_gado",
            "2": "gudeg",
            "3": "nasi_goreng",
            "4": "pempek",
            "5": "rendang",
            "6": "sate",
            "7": "soto"
          }
        }
        """
    )
    (artifact_dir / "training_config_resolved.json").write_text(
        """
        {
          "model_name": "convnext_tiny.fb_in1k",
          "num_classes": 8,
          "class_names": [
            "bakso",
            "gado_gado",
            "gudeg",
            "nasi_goreng",
            "pempek",
            "rendang",
            "sate",
            "soto"
          ],
          "image_size": 256
        }
        """
    )
    monkeypatch.setenv("NUTRISCAN_MODEL_ARTIFACT_DIR", str(artifact_dir))
    monkeypatch.setenv("NUTRISCAN_MODEL_VERSION", "convnext-tiny-test")

    route_paths = {route.path for route in app.routes}
    assert "/readyz" in route_paths

    payload = readyz()

    assert payload["status"] == "ready"
    assert payload["modelVersion"] == "convnext-tiny-test"
    assert payload["artifactLocation"] == str(artifact_dir)
    assert payload["modelName"] == "convnext_tiny.fb_in1k"
    assert payload["imageSize"] == 256
    assert payload["labelCount"] == 8
    assert payload["confidenceThreshold"] == 0.6
