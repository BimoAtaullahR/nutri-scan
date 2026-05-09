# Context Map

## Contexts

- [Backend API](./services/backend/CONTEXT.md) — owns user-facing backend workflows, persistence, scan orchestration, nudge decisions, and trend APIs.
- [AI/ML Inference](./services/ai-inference/CONTEXT.md) — owns image preprocessing, food recognition, visual dominance detection, confidence scoring, and estimated energy payloads.
- Mobile App — owns camera flow, scan feedback, preventive nudge presentation, and Simple Energy Trend UI.
- Shared Product Domain — owns product language that cuts across contexts, including preventive nutrition intelligence and pre-portioning behavior.

## Relationships

- **Mobile App -> Backend API**: submits scan requests, renders scan feedback, records user responses to nudges, and fetches weekly trend data.
- **Backend API -> AI/ML Inference**: sends food images or image references for inference and receives structured scan results.
- **Backend API -> Shared Product Domain**: persists domain events such as scans, nudges, and user responses using the shared language.
- **AI/ML Inference -> Shared Product Domain**: reports food category, visual dominance, confidence, and estimated energy using terms understood by the rest of the product.
