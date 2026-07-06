const mockData = {
  "cases": [
    {
      "description": "Answer a question based on provided context.",
      "id": "context-based-answer",
      "output": {
        "answer": "The project was completed three weeks ahead of schedule due to the team's efficient sprint planning and early resolution of key technical blockers.",
        "confidence": 0.92,
        "grounded": true
      }
    }
  ],
  "delay": {
    "maxMs": 120,
    "minMs": 30
  },
  "random": {
    "enabled": false,
    "seed": "answer-question"
  },
  "version": "1.0.0"
};

export default mockData;
