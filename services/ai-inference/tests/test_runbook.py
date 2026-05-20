from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def test_runbook_documents_full_demo_path() -> None:
    runbook = (ROOT / "RUNBOOK.md").read_text()

    for required_section in (
        "## Setup",
        "## Dataset Prep",
        "## Training",
        "## Evaluation",
        "## Serving",
        "## Smoke Test",
        "## Known Limitations",
    ):
        assert required_section in runbook

    assert "model-artifacts/selected-mvp-classifier" in runbook
    assert "reports/baseline-food-classifier/metrics.json" in runbook
    assert "curl http://localhost:8000/readyz" in runbook
    assert "curl -X POST" in runbook
    assert "MODEL_COMPARISON.md" in runbook
    assert "configs/selected_mvp_classifier.json" in runbook
