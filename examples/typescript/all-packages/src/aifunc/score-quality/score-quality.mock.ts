const mockData = {
  "cases": [
    {
      "description": "Evaluate quality of a developer-facing documentation paragraph.",
      "id": "developer-doc",
      "output": {
        "actionabilityScore": 62,
        "clarityScore": 78,
        "level": "good",
        "overallScore": 72,
        "structureScore": 70,
        "suggestions": [
          "Add a concrete code example to illustrate the main concept.",
          "Break the paragraph into shorter sections with subheadings.",
          "End with a clear next step or call to action."
        ],
        "summary": "The text is clear and well-toned but could be more actionable and better structured.",
        "toneScore": 80
      }
    }
  ],
  "delay": {
    "maxMs": 150,
    "minMs": 50
  },
  "random": {
    "enabled": false,
    "seed": "score-quality"
  },
  "version": "1.0.0"
};

export default mockData;
