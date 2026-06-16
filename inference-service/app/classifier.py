class DocumentClassifier:
    _KEYWORDS: dict[str, list[str]] = {
        "finance": ["invoice", "payment", "bank", "revenue", "tax"],
        "legal": ["contract", "agreement", "clause", "liability"],
        "medical": ["diagnosis", "patient", "prescription", "symptom"],
    }
    def predict(self, text: str) -> tuple[str, float]:
        text = text.lower()
        for label, keywords in self._KEYWORDS.items():
            for keyword in keywords:
                if keyword in text:
                    return label, 0.8

        return "other", 0.4