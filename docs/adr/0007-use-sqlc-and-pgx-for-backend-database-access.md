# Use sqlc and pgx for Backend Database Access

NutriScan will use sqlc with PostgreSQL and pgx for Backend API database access. Writing SQL explicitly keeps scan lifecycle, nudge response, and weekly trend queries transparent, while sqlc generates type-safe Go code and avoids hand-written query boilerplate.

## Considered Options

- Use an ORM such as GORM.
- Use hand-written SQL with manual scanning.
- Use sqlc-generated Go code with pgx.

## Consequences

- Backend modules own their SQL queries and repository adapters.
- Database schema changes must be reflected in migrations and regenerated sqlc code.
- Complex trend queries remain visible as SQL rather than hidden behind ORM behavior.
