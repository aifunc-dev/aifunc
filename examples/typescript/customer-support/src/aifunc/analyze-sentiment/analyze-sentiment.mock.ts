const mockData = {
  "cases": [
    {
      "description": "Positive sentiment text example.",
      "id": "positive-text",
      "output": {
        "confidence": 0.92,
        "label": "positive",
        "rankings": [
          {
            "label": "positive",
            "score": 0.92
          },
          {
            "label": "neutral",
            "score": 0.05
          },
          {
            "label": "negative",
            "score": 0.03
          }
        ]
      }
    }
  ],
  "delay": {
    "maxMs": 100,
    "minMs": 30
  },
  "random": {
    "enabled": false,
    "seed": "sentiment-analysis"
  },
  "version": "1.0.0"
};

export default mockData;
