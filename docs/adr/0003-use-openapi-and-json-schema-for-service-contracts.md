# Use OpenAPI and JSON Schema for Service Contracts

NutriScan will keep cross-context service contracts in `packages/contracts`, using OpenAPI for the Mobile App to Backend API boundary and JSON Schema for the Backend API to AI/ML Inference boundary. OpenAPI gives the mobile client and Go backend a neutral HTTP contract, while JSON Schema keeps the internal inference payload lightweight for a Python-based AI/ML service during the MVP.

## Considered Options

- Use manually written documentation only.
- Use TypeScript types as the source of truth.
- Use protobuf for all service boundaries.
- Use OpenAPI for the public backend API and JSON Schema for internal inference payloads.

## Consequences

- Contract changes must be reviewed as shared product changes, not hidden inside one service.
- The backend remains the stable API surface for the mobile app.
- AI inference can evolve quickly while still returning a structured payload the backend can validate.
