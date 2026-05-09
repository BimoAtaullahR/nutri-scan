# Use Bearer Token for Anonymous User Identity

NutriScan will let the Backend API create anonymous users and issue a bearer token for MVP identity. The Mobile App stores that token and sends it through the `Authorization` header so scan, profile, nudge, and trend data can be owned by a durable anonymous user without introducing full login and registration yet.

## Considered Options

- Send anonymous user IDs in request bodies.
- Use custom identity headers.
- Issue backend-owned bearer tokens for anonymous users.

## Consequences

- Backend API must provide an anonymous user creation endpoint.
- Mobile must persist the token and attach it to authenticated MVP requests.
- A future full-auth flow can reuse the same authorization pattern and link anonymous data to registered accounts.
