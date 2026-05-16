-- +goose Up
ALTER TABLE scans
    ADD COLUMN meal_type TEXT NOT NULL DEFAULT 'snack' CHECK (meal_type IN ('breakfast', 'lunch', 'dinner', 'snack')),
    ADD COLUMN estimated_energy_min_kcal INTEGER CHECK (estimated_energy_min_kcal >= 0),
    ADD COLUMN estimated_energy_max_kcal INTEGER CHECK (estimated_energy_max_kcal >= 0),
    ADD CONSTRAINT scans_estimated_energy_range_valid CHECK (
        estimated_energy_min_kcal IS NULL
        OR estimated_energy_max_kcal IS NULL
        OR estimated_energy_min_kcal <= estimated_energy_max_kcal
    );

-- +goose Down
ALTER TABLE scans
    DROP CONSTRAINT scans_estimated_energy_range_valid,
    DROP COLUMN estimated_energy_max_kcal,
    DROP COLUMN estimated_energy_min_kcal,
    DROP COLUMN meal_type;
