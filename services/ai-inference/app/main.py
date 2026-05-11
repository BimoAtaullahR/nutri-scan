from fastapi import FastAPI


app = FastAPI(title="NutriScan AI/ML Inference", version="0.1.0")


@app.get("/healthz")
def healthz() -> dict[str, str]:
    return {"status": "ok"}
