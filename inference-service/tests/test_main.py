from fastapi.testclient import TestClient
from app.main import app

client = TestClient(app)

def test_health():
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json() == {"status": "ok"}

def test_predict():
    response = client.post("/predict", json={"document_id": "123", "text": "invoice for payment"})
    assert response.status_code == 200
    data = response.json()
    assert data["document_id"] == "123"
    assert data["label"] == "finance"
    assert data["confidence"] == 0.8