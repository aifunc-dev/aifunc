const mockData = {
  "cases": [
    {
      "description": "Generate slug from a Chinese article title.",
      "id": "chinese-article",
      "output": {
        "metaDescription": "A beginner-friendly guide to understanding the fundamentals of machine learning, covering key concepts and practical applications.",
        "slug": "getting-started-with-machine-learning",
        "tags": [
          "machine learning",
          "beginner",
          "AI",
          "tutorial"
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
    "seed": "generate-slug"
  },
  "version": "1.0.0"
};

export default mockData;
