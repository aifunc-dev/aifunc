const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "extract-keywords",
    "version": "1.0.0",
    "description": "Extract keywords and key phrases from text with relevance scores.",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.2.0"
  },
  "api": {
    "name": "extract-keywords",
    "description": "Extract keywords and key phrases from text, ranked by relevance.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "maxKeywords": {
          "default": 10,
          "description": "Maximum number of keywords to return. Defaults to 10.",
          "maximum": 50,
          "minimum": 1,
          "type": "integer"
        },
        "text": {
          "description": "The text to extract keywords from.",
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
        "keywords": {
          "description": "Extracted keywords sorted by relevance (highest first).",
          "items": {
            "additionalProperties": false,
            "properties": {
              "relevance": {
                "description": "Relevance score between 0 and 1.",
                "maximum": 1,
                "minimum": 0,
                "type": "number"
              },
              "word": {
                "description": "The keyword or key phrase.",
                "type": "string"
              }
            },
            "required": [
              "word",
              "relevance"
            ],
            "type": "object"
          },
          "type": "array"
        }
      },
      "required": [
        "keywords"
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
          "maxTokens": 600
        }
      }
    ]
  },
  "prompts": {
    "general": "# System\n\nYou are a keyword extraction function. You must only return a JSON object in the following format:\n{\"keywords\": [{\"word\": \"\u003ckeyword\u003e\", \"relevance\": \u003c0-1\u003e}, ...]}\n\nDo not output Markdown, do not include any extra explanation, and do not add undeclared fields.\n\nRequirements:\n- Extract the most important keywords and key phrases from the input text.\n- Rank them by relevance to the text's core topics, with the most relevant first.\n- `relevance` should be a float between 0 and 1, where 1 means the keyword is central to the text.\n- Prefer concise phrases (1-3 words) over single characters or overly long phrases.\n- Do not return duplicates or near-duplicates (e.g. singular and plural forms of the same word).\n- Return at most `maxKeywords` results; if not specified, default to 10.\n- Keywords should be in the same language as the input text.\n\n# User\n\nText:\n{{text}}\n\nMaximum keywords: {{maxKeywords}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-16T00:03:04Z",
    "contentHash": "sha256:f9b675229aa079c3225565e76baef7270f03bc30e31ea6c4e17673afbe90ea90"
  }
};

export default artifact;
