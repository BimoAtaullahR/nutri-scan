from __future__ import annotations

import json
from functools import lru_cache
from pathlib import Path
from typing import TypedDict


CONFIG_PATH = Path(__file__).resolve().parents[1] / "configs" / "estimated_energy_ranges.json"


class EnergyRange(TypedDict):
    minKcal: int
    maxKcal: int


@lru_cache(maxsize=1)
def load_energy_table() -> dict[str, object]:
    return json.loads(CONFIG_PATH.read_text())


def lookup_energy_range(food_category: str, portion: str) -> EnergyRange:
    table = load_energy_table()
    ranges = table["ranges"]
    if food_category not in ranges:
        raise KeyError(f"Unknown food category: {food_category}")

    category_ranges = ranges[food_category]
    if portion not in category_ranges:
        raise KeyError(f"Unknown portion: {portion}")

    return category_ranges[portion]
