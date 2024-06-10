from transformers import pipeline


class Processor:
    def __init__(self):
        self.pipeline = pipeline("zero-shot-classification", model="joeddav/xlm-roberta-large-xnli")

    def predict(self, text, candidate_labels, min_score=0.15) -> dict[str, float]:
        result = self.pipeline(text, candidate_labels, multi_label=False)
        return {k: v for k, v in zip(result['labels'], result['scores']) if v > min_score}
