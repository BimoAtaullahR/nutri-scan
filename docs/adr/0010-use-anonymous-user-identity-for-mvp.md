# Use Anonymous User Identity for MVP

NutriScan will use anonymous user identity for the MVP instead of full login and registration. This keeps the first release focused on scan, nudge, and weekly trend behavior while still giving the backend a durable user key that can later be linked to a full authentication provider.

## Considered Options

- Require full account registration from the start.
- Use anonymous user identity during MVP.
- Store scans only on the device without backend user ownership.

## Consequences

- Mobile must persist and send a generated anonymous identity.
- Backend data models should include user ownership from the beginning.
- Account linking can be added later without redesigning scan and trend ownership.
