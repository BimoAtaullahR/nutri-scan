# NutriScan Mobile App

Flutter mobile app for the camera-first NutriScan experience.

## Ownership

Mobile owns:

- Camera flow and image preparation before upload
- Scan feedback presentation
- Preventive nudge presentation
- Weekly energy trend UI
- Local UI state and mobile-specific UX

Mobile does not own:

- Scan lifecycle persistence
- Nudge decision rules
- AI/ML inference
- Backend data storage

## Planned Structure

```txt
lib/
  main.dart
  main_dev.dart
  main_prod.dart
  app/          # app bootstrap, router, theme
  core/         # config, errors, network, storage, shared widgets
  domain/       # mobile domain models
  data/         # API clients, DTOs, repositories
  ui/           # screens, view models, widgets by feature
  generated/    # generated API/client code when adopted
```

Run `flutter create .` from this directory when the mobile developer is ready to generate the Flutter project files.
