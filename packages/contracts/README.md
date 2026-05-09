# NutriScan Contracts

Shared service contracts for boundaries between NutriScan contexts.

## Ownership

Contract changes are shared product changes. They should be reviewed by the owners of every context that consumes the contract.

## Structure

```txt
openapi/
  backend-api.yaml              # Mobile App -> Backend API
ai-inference/
  scan-inference.schema.json    # Backend API -> AI/ML Inference
```
