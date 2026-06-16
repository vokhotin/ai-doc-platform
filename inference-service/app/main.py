from fastapi import FastAPI
from app.schemas import PredictionRequest, PredictionResponse
from app.classifier import DocumentClassifier

app = FastAPI()
classifier = DocumentClassifier()

@app.get("/health")
def health():
    return {"status": "ok"}

@app.post("/predict", response_model=PredictionResponse)
def predict(request: PredictionRequest):
    (label, confidence) = classifier.predict(request.text)

    return {"label": label, "confidence": confidence, "document_id": request.document_id}