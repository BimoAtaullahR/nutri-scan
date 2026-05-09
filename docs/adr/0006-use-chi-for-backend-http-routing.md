# Use Chi for Backend HTTP Routing

NutriScan will use Chi with Go's standard `net/http` package for Backend API routing. Chi is lightweight and composable, supports middleware and mounted sub-routers, and fits the backend's modular structure for scan, nudge, trend, and platform HTTP concerns without making the framework the center of the application.

## Considered Options

- Use only `net/http`.
- Use Chi with `net/http`.
- Use a larger web framework such as Gin, Echo, or Fiber.

## Consequences

- Backend handlers remain standard `net/http` handlers.
- Feature routes can be grouped or mounted by module.
- Middleware such as request IDs, recovery, logging, and request timeouts can be applied consistently at the platform HTTP layer.
