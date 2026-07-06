const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "summarize",
    "version": "1.0.0",
    "description": "Generate a concise summary of the input text.",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.1.0"
  },
  "api": {
    "name": "summarize",
    "description": "Generate a concise summary of the input text.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "maxLength": {
          "default": 80,
          "description": "Maximum word count for the summary. Defaults to 80.",
          "maximum": 300,
          "minimum": 20,
          "type": "integer"
        },
        "text": {
          "description": "The text to summarize.",
          "minLength": 1,
          "type": "string"
        }
      },
      "required": [
        "text"
      ],
      "type": "object"
    },
    "output": {
      "additionalProperties": false,
      "properties": {
        "summary": {
          "description": "The generated summary.",
          "type": "string"
        },
        "wordCount": {
          "description": "Approximate word count of the summary.",
          "minimum": 0,
          "type": "integer"
        }
      },
      "required": [
        "summary",
        "wordCount"
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
          "maxTokens": 300
        }
      }
    ]
  },
  "prompts": {
    "general": "# System\n\nYou are a strict summarization function. You must only return a JSON object that conforms to the output schema. Do not output Markdown, do not include any extra explanation, and do not add undeclared fields.\n\nRequirements:\n- The summary language must match the language of the input text — do not translate.\n- The summary should be concise, accurate, and fluent.\n- Preserve the most essential information; do not fabricate anything not present in the original.\n- If the input covers multiple points, prioritize the most important 1 to 3.\n- `summary` must not exceed the word/character count specified by `maxLength`; default to 80 if not provided.\n- `wordCount` should reflect the approximate length of the summary.\n\n# User\n\nText:\n{{text}}\n\nMaximum length:\n{{maxLength}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-05T03:56:46Z",
    "contentHash": "sha256:48cfa8b2183761308b7c77d705f9a785122077790f7267c8275c5e0452b38f6f"
  }
};

export default artifact;
