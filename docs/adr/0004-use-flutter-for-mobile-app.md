# Use Flutter for Mobile App

NutriScan will use Flutter for the Mobile App because the mobile developer owns that stack and the MVP is camera-first. The app will live under `apps/mobile` and organize Dart code around `app`, `core`, `domain`, `data`, and `ui` so camera scanning, nudge feedback, and weekly trend views stay separated while sharing backend API access.

## Consequences

- Mobile implementation uses Flutter and Dart rather than React Native or Next.js.
- The mobile app consumes Backend API contracts through OpenAPI-derived or manually maintained Dart client code.
- Native platform folders such as `android` and `ios` remain owned by the Mobile App context.
