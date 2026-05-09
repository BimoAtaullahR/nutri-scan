# Use PostgreSQL for Backend Persistence

NutriScan will use PostgreSQL as the Backend API persistence store for the MVP. Scan lifecycle records, nudge responses, user profiles, and weekly energy trends are relational enough to benefit from SQL constraints and queryability, while keeping deployment and operational complexity reasonable.

## Considered Options

- Use PostgreSQL.
- Use a document database for flexible scan payloads.
- Use only local/in-memory storage during MVP development.

## Consequences

- Backend changes that alter persisted product data need migrations.
- Weekly trend queries can be built directly from completed scan records and nudge response history.
- AI inference remains stateless from a product persistence perspective.
