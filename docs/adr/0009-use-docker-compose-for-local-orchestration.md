# Use Docker Compose for Local Orchestration

NutriScan will use Docker Compose for local dependency and service orchestration, especially PostgreSQL, Backend API, and AI/ML Inference smoke testing. Developers can still run each source stack directly during daily work, with Flutter run through normal Flutter tooling rather than inside Docker.

## Consequences

- Local onboarding can start from a shared Compose workflow.
- Backend and AI/ML integration can be smoke-tested without manual service setup.
- Mobile debugging stays aligned with Flutter's native tooling and device workflow.
