from app.classifier import DocumentClassifier

def test_predict_finance():
    classifier = DocumentClassifier()
    label, confidence = classifier.predict("Please process the invoice payment")
    assert label == "finance"
    assert confidence == 0.8

def test_predict_other():
    classifier = DocumentClassifier()
    label, confidence = classifier.predict("A simple test message")
    assert label == "other"
    assert confidence == 0.4
