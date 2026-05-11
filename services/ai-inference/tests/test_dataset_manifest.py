import json
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
EXPECTED_CLASSES = [
    "nasi_goreng",
    "sate",
    "rendang",
    "bakso",
    "gado_gado",
    "soto",
    "pempek",
    "gudeg",
]


def test_mvp_food_category_mapping_is_machine_readable():
    mapping_path = ROOT / "configs" / "mvp_food_categories.json"

    assert mapping_path.exists()

    mapping = json.loads(mapping_path.read_text())

    assert [category["slug"] for category in mapping["categories"]] == EXPECTED_CLASSES
    assert all(category["displayName"] for category in mapping["categories"])
    assert mapping["deferredCategories"] == ["nasi_padang"]


def test_dataset_manifest_records_primary_source_and_deferred_scope():
    manifest = (ROOT / "data" / "manifests" / "mvp_food_dataset.md").read_text()

    assert "Indonesian Food Image - Mendeley Data" in manifest
    assert "https://data.mendeley.com/datasets/vtjd68bmwt" in manifest
    assert "CC BY 4.0" in manifest
    assert "DOI: https://doi.org/10.17632/vtjd68bmwt.1" in manifest
    assert "`nasi_padang`: future meal-level category" in manifest
