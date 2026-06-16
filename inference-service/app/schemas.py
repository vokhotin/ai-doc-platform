from pydantic import BaseModel, Field


class PredictionRequest(BaseModel):
    document_id: str
    text: str

class PredictionResponse(BaseModel):
    document_id: str
    label: str
    confidence: float = Field(ge=0.0, le=1.0)