const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "answer-question",
    "version": "1.0.0",
    "description": "Answer a question based on provided context or general knowledge.",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.1.0"
  },
  "api": {
    "name": "answer-question",
    "description": "Answer a question based on provided context or general knowledge.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "context": {
          "description": "Source text or document to base the answer on. If omitted, uses general knowledge.",
          "type": "string"
        },
        "language": {
          "description": "Answer language. If omitted, matches the question language.",
          "type": "string"
        },
        "maxLength": {
          "description": "Maximum word count for the answer. Default: 100.",
          "maximum": 500,
          "minimum": 20,
          "type": "integer"
        },
        "question": {
          "description": "The question to answer.",
          "minLength": 1,
          "type": "string"
        }
      },
      "required": [
        "question"
      ],
      "type": "object"
    },
    "output": {
      "additionalProperties": false,
      "properties": {
        "answer": {
          "description": "The generated answer.",
          "type": "string"
        },
        "confidence": {
          "description": "Confidence score between 0 and 1.",
          "maximum": 1,
          "minimum": 0,
          "type": "number"
        },
        "grounded": {
          "description": "True if the answer is based on the provided context, false if from general knowledge.",
          "type": "boolean"
        }
      },
      "required": [
        "answer",
        "grounded",
        "confidence"
      ],
      "type": "object"
    }
  },
  "modelParams": {
    "schemaVersion": "0.1.0",
    "rules": [
      {
        "match": {
          "pattern": ".*"
        },
        "params": {
          "temperature": 0.2,
          "maxTokens": 1024
        }
      }
    ]
  },
  "prompts": {
    "general": "# System\n\nYou are a question answering function. You must only return a JSON object in the following format:\n{\"answer\": \"\u003canswer text\u003e\", \"grounded\": \u003ctrue|false\u003e, \"confidence\": \u003c0.0-1.0\u003e}\n\nDo not output Markdown, do not include any extra explanation, and do not add undeclared fields.\n\nRequirements:\n- If context is provided, answer based solely on that context. Set 'grounded' to true.\n- If no context is provided, answer from general knowledge. Set 'grounded' to false.\n- If the context does not contain enough information to answer, say so clearly and set confidence below 0.5.\n- Answer in the requested language; if not specified, match the question language.\n- Keep the answer within maxLength words (default: 100). Be concise and direct.\n- Set 'confidence' to reflect how certain you are about the answer (0.0 = no idea, 1.0 = certain).\n\n# User\n\nQuestion: {{question}}\n\nContext:\n{{context}}\n\nMax length (words): {{maxLength}}\n\nLanguage: {{language}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-09T15:12:49Z",
    "contentHash": "sha256:73e679d31abadd6f21cc8705bf016c8611968ea3c928303cd074459654eb9343"
  }
};

export default artifact;
