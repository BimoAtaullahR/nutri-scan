from pathlib import Path
import tomllib


ROOT = Path(__file__).resolve().parents[1]


def test_pyproject_declares_ai_inference_tooling():
    pyproject = ROOT / "pyproject.toml"

    assert pyproject.exists()

    config = tomllib.loads(pyproject.read_text())
    dependencies = set(config["project"]["dependencies"])
    optional_dependencies = config["project"]["optional-dependencies"]

    assert "fastapi" in dependencies
    assert "python-multipart" in dependencies
    assert "torch" in dependencies
    assert "torchvision" in dependencies
    assert "timm" in dependencies
    assert "pillow" in dependencies
    assert "pytest" in optional_dependencies["dev"]
    assert "ruff" in optional_dependencies["dev"]


def test_gitignore_keeps_large_ai_artifacts_local():
    gitignore = (ROOT / ".gitignore").read_text().splitlines()

    assert "data/raw/" in gitignore
    assert "data/processed/" in gitignore
    assert "model-artifacts/*" in gitignore
    assert "!model-artifacts/.gitkeep" in gitignore
    assert "reports/" in gitignore
