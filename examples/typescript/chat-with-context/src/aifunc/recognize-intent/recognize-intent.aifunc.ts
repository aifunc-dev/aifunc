const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "recognize-intent",
    "version": "1.0.0",
    "description": "Recognize user intent from conversational text with confidence scores (zero-shot).",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.1.0"
  },
  "api": {
    "name": "recognize-intent",
    "description": "Recognize user intent from conversational text with confidence scores.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "context": {
          "description": "Optional conversation context or system description to improve recognition accuracy.",
          "type": "string"
        },
        "intents": {
          "description": "List of candidate intents to recognize from.",
          "items": {
            "minLength": 1,
            "type": "string"
          },
          "minItems": 2,
          "type": "array"
        },
        "text": {
          "description": "The user message to recognize intent from.",
          "minLength": 1,
          "type": "string"
        }
      },
      "required": [
        "text",
        "intents"
      ],
      "type": "object"
    },
    "output": {
      "additionalProperties": false,
      "properties": {
        "confidence": {
          "description": "Confidence score of the top intent, between 0 and 1.",
          "maximum": 1,
          "minimum": 0,
          "type": "number"
        },
        "intent": {
          "description": "The highest-confidence recognized intent.",
          "type": "string"
        },
        "rankings": {
          "description": "All intents ranked by confidence (highest first).",
          "items": {
            "additionalProperties": false,
            "properties": {
              "confidence": {
                "description": "Confidence score between 0 and 1.",
                "maximum": 1,
                "minimum": 0,
                "type": "number"
              },
              "intent": {
                "description": "Intent label.",
                "type": "string"
              }
            },
            "required": [
              "intent",
              "confidence"
            ],
            "type": "object"
          },
          "type": "array"
        }
      },
      "required": [
        "intent",
        "confidence",
        "rankings"
      ],
      "type": "object"
    }
  },
  "modelParams": {
    "schemaVersion": "",
    "rules": null
  },
  "prompts": {
    "general": "# System\n\nYou are a conversational intent recognition function. You must only return a JSON object in the following format:\n{\"intent\": \"\u003ctop_intent\u003e\", \"confidence\": \u003c0-1\u003e, \"rankings\": [{\"intent\": \"\u003clabel\u003e\", \"confidence\": \u003c0-1\u003e}, ...]}\n\nDo not output Markdown, do not include any extra explanation, and do not add undeclared fields.\n\nRequirements:\n- Analyze the user's message to determine their underlying intent (what they want to accomplish).\n- Only use intents from the provided candidate list — never invent new ones.\n- Assign a confidence score (0 to 1) to each intent indicating how likely the user's message maps to that intent.\n- Sort `rankings` by confidence from highest to lowest.\n- The `intent` field should contain the highest-ranked intent label.\n- Focus on the user's goal and action, not the topic or sentiment.\n- If context is provided, use it to disambiguate between similar intents.\n\n# User\n\nMessage:\n{{text}}\n\nCandidate intents:\n{{intents}}\n\n{{#if context}}\nContext:\n{{context}}\n{{/if}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-05T03:57:59Z",
    "contentHash": "sha256:447ebba1e669ee49b44239ad85f6cae2de2eaa4ea7ed1e658c1d4e3c0f909d91"
  }
};

export default artifact;
