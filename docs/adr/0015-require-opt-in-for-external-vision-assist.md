# Require Opt-In for External Vision Assist

NutriScan may send scan images to an external vision-language model only in development, explicit opt-in, or clearly configured fallback workflows. This allows external models to support ambiguous scans and dataset bootstrapping without making third-party image processing the silent default for every scan.

## Consequences

- External vision assist must be distinguishable from the default local inference path.
- Product and development workflows must not assume scan images are stored or sent externally by default.
- Any external model result must still be normalized into NutriScan's inference and nudge language before user-facing feedback is produced.
