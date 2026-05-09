# Separate Backend API and AI Inference

NutriScan will use a Go Backend API as the main orchestration and persistence layer, with AI/ML inference implemented as a separate service. This keeps the Go backend focused on product workflows, data, and latency control while allowing the inference service to use the Python/model tooling that is more natural for computer vision experimentation and scaling.

## Considered Options

- Put AI inference directly inside the Go backend for a simpler MVP.
- Split AI inference into a separate service from the start.

## Consequences

- The MVP has one extra service boundary to operate.
- AI bottlenecks can be scaled and optimized independently from the main backend.
- The backend API remains the stable integration point for the mobile client.
