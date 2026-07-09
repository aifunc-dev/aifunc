const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "detect-language",
    "version": "1.0.0",
    "description": "Detect the language of input text, returning a language code and confidence score.",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.1.0"
  },
  "api": {
    "name": "detect-language",
    "description": "Detect the language of input text, returning a language code, language name, and confidence score.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "text": {
          "description": "The text whose language should be detected.",
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
        "confidence": {
          "description": "Confidence score between 0 and 1.",
          "maximum": 1,
          "minimum": 0,
          "type": "number"
        },
        "language": {
          "description": "Detected language code (e.g. 'en', 'zh-CN', 'ja', 'fr').",
          "type": "string"
        },
        "languageName": {
          "description": "Human-readable name of the detected language (e.g. 'English', '中文', '日本語').",
          "type": "string"
        }
      },
      "required": [
        "language",
        "languageName",
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
          "temperature": 0,
          "maxTokens": 200
        }
      }
    ]
  },
  "prompts": {
    "general": "# System\n\nYou are a language detection function. You must only return a JSON object in the following format:\n{\"language\": \"\u003ccode\u003e\", \"languageName\": \"\u003cname\u003e\", \"confidence\": \u003c0-1\u003e}\n\nDo not output Markdown, do not include any extra explanation, and do not add undeclared fields.\n\nRequirements:\n- Identify the primary language of the input text.\n- `language`: standard language code (BCP 47 / ISO 639, e.g. \"en\", \"zh-CN\", \"zh-TW\", \"ja\", \"ko\", \"fr\", \"de\", \"es\", \"pt-BR\").\n- `languageName`: human-readable language name, written in that language itself (e.g. \"English\", \"中文\", \"日本語\").\n- `confidence`: a float between 0 and 1 indicating how certain you are.\n- If the text contains multiple languages, detect the dominant one.\n- For very short or ambiguous text, lower the confidence accordingly.\n\n# User\n\nText:\n{{text}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-09T15:12:49Z",
    "contentHash": "sha256:a1de9b8449d941e016d225a854871a9a8a7f73b0ce9d45a37449559c37794c5f"
  }
};

export default artifact;
