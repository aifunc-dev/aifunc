const mockData = {
  "cases": [
    {
      "description": "Extract structured info from a resume snippet.",
      "id": "resume-extraction",
      "output": {
        "extracted": {
          "email": "zhangwei@example.com",
          "name": "Zhang Wei",
          "skills": [
            "Python",
            "Machine Learning",
            "SQL"
          ],
          "yearsOfExperience": 5
        },
        "missing": []
      }
    }
  ],
  "delay": {
    "maxMs": 150,
    "minMs": 30
  },
  "random": {
    "enabled": false,
    "seed": "json-extract"
  },
  "version": "1.0.0"
};

export default mockData;
