# Use Goose for Backend Database Migrations

NutriScan will use Goose for Backend API PostgreSQL migrations. Goose keeps migration files as explicit SQL with `-- +goose Up` and `-- +goose Down` sections, provides simple CLI commands for status, applying, and rolling back migrations, and is lightweight enough for the MVP team workflow.

## Considered Options

- Run schema changes manually.
- Use Goose.
- Use Atlas.
- Use golang-migrate.

## Consequences

- Backend schema changes must include versioned SQL migration files.
- Local and deployment workflows need a Goose command against the configured PostgreSQL database.
- Migration files remain readable without hiding schema changes behind application code.
