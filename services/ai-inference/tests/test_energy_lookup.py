import json
from pathlib import Path

from app.energy import lookup_energy_range


ROOT = Path(__file__).resolve().parents[1]
PORTIONS = ("small", "medium", "large")


def test_energy_lookup_covers_all_mvp_categories_and_portions() -> None:
    categories = json.loads((ROOT / "configs" / "mvp_food_categories.json").read_text())[
        "categories"
    ]

    for category in categories:
        previous_max = 0
        for portion in PORTIONS:
            energy = lookup_energy_range(category["slug"], portion)
            assert energy["minKcal"] > 0
            assert energy["maxKcal"] >= energy["minKcal"]
            assert energy["minKcal"] >= previous_max
            previous_max = energy["maxKcal"]


def test_energy_lookup_rejects_unknown_category_or_portion() -> None:
    try:
        lookup_energy_range("pizza", "medium")
    except KeyError as exc:
        assert "Unknown food category" in str(exc)
    else:
        raise AssertionError("Expected unknown category to fail")

    try:
        lookup_energy_range("sate", "extra_large")
    except KeyError as exc:
        assert "Unknown portion" in str(exc)
    else:
        raise AssertionError("Expected unknown portion to fail")
