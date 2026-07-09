const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "classify",
    "version": "1.0.0",
    "description": "Classify text into user-defined categories with confidence scores (zero-shot).",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.1.0"
  },
  "api": {
    "name": "classify",
    "description": "Classify text into user-defined categories with confidence scores.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "allowMultiple": {
          "default": false,
          "description": "If true, the text can be assigned to multiple categories. Defaults to false.",
          "type": "boolean"
        },
        "categories": {
          "description": "List of candidate categories to classify into.",
          "items": {
            "minLength": 1,
            "type": "string"
          },
          "minItems": 2,
          "type": "array"
        },
        "text": {
          "description": "The text to classify.",
          "minLength": 1,
          "type": "string"
        }
      },
      "required": [
        "text",
        "categories"
      ],
      "type": "object"
    },
    "output": {
      "additionalProperties": false,
      "properties": {
        "classifications": {
          "description": "Classification results sorted by confidence (highest first).",
          "items": {
            "additionalProperties": false,
            "properties": {
              "category": {
                "description": "The category label.",
                "type": "string"
              },
              "confidence": {
                "description": "Confidence score between 0 and 1.",
                "maximum": 1,
                "minimum": 0,
                "type": "number"
              }
            },
            "required": [
              "category",
              "confidence"
            ],
            "type": "object"
          },
          "type": "array"
        }
      },
      "required": [
        "classifications"
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
          "temperature": 0.1,
          "maxTokens": 400
        }
      }
    ]
  },
  "prompts": {
    "general": "# System\n\nYou are a zero-shot text classification function. You must only return a JSON object in the following format:\n{\"classifications\": [{\"category\": \"\u003clabel\u003e\", \"confidence\": \u003c0-1\u003e}, ...]}\n\nDo not output Markdown, do not include any extra explanation, and do not add undeclared fields.\n\nRequirements:\n- Classify the input text into the provided candidate categories.\n- Assign a confidence score (0 to 1) to each category indicating how well the text fits.\n- Sort `classifications` by confidence from highest to lowest.\n- If `allowMultiple` is false, only one category should have a high confidence (dominant), and the rest should sum to roughly 1.\n- If `allowMultiple` is true, each category's confidence should independently reflect how well the text fits that category (they need not sum to 1).\n- Only use categories from the provided list — never invent new ones.\n- Base classification on the text's content, meaning, and intent.\n\n# User\n\nText:\n{{text}}\n\nCategories:\n{{categories}}\n\nAllow multiple: {{allowMultiple}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-09T15:12:49Z",
    "contentHash": "sha256:de1695d4164aab5906d66496030b77676dfb754e96971dddebf1f82a9df19b69"
  }
};

export default artifact;
