# Use Monorepo for MVP

NutriScan will use a single monorepo for the MVP, with separate top-level areas for the mobile app, backend API, AI/ML inference, shared contracts, and documentation. This keeps rapidly changing scan, inference, nudge, and trend contracts easy to evolve during the 30-day MVP window, while avoiding early coordination overhead from separate repositories.

## Considered Options

- Use three separate repositories for backend, mobile app, and AI/ML.
- Use one monorepo with clear ownership boundaries.

## Consequences

- Teams share one repository workflow and must keep folder ownership clear.
- Shared API contracts can live beside the services that consume them.
- Deployment still needs to treat backend API, AI/ML inference, and mobile app as independently releasable units.
