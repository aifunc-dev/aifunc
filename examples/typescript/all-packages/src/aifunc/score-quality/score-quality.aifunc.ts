const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "score-quality",
    "version": "1.0.0",
    "description": "Evaluate text quality across multiple dimensions and return structured scores with improvement advice.",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.2.0"
  },
  "api": {
    "name": "score-quality",
    "description": "Evaluate text quality across multiple dimensions and return structured scores with improvement advice.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "maxSuggestions": {
          "default": 3,
          "description": "Maximum number of improvement suggestions to return.",
          "maximum": 10,
          "minimum": 1,
          "type": "integer"
        },
        "purpose": {
          "default": "general communication",
          "description": "Writing purpose (e.g. 'explanation', 'marketing', 'support', 'learning').",
          "type": "string"
        },
        "strictness": {
          "default": 3,
          "description": "Evaluation strictness from 1 (lenient) to 5 (very strict).",
          "maximum": 5,
          "minimum": 1,
          "type": "integer"
        },
        "targetAudience": {
          "default": "general readers",
          "description": "Intended readers (e.g. 'developers', 'customers', 'students', 'executives').",
          "type": "string"
        },
        "text": {
          "description": "The text content to evaluate.",
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
        "actionabilityScore": {
          "description": "How actionable and useful the text is, from 0 to 100.",
          "maximum": 100,
          "minimum": 0,
          "type": "integer"
        },
        "clarityScore": {
          "description": "Clarity score from 0 to 100.",
          "maximum": 100,
          "minimum": 0,
          "type": "integer"
        },
        "level": {
          "description": "Overall quality level.",
          "enum": [
            "excellent",
            "good",
            "fair",
            "poor"
          ],
          "type": "string"
        },
        "overallScore": {
          "description": "Overall quality score from 0 to 100.",
          "maximum": 100,
          "minimum": 0,
          "type": "integer"
        },
        "structureScore": {
          "description": "Organization and logical structure score from 0 to 100.",
          "maximum": 100,
          "minimum": 0,
          "type": "integer"
        },
        "suggestions": {
          "description": "Concise improvement suggestions, ordered by importance.",
          "items": {
            "type": "string"
          },
          "minItems": 0,
          "type": "array"
        },
        "summary": {
          "description": "One short sentence summarizing the evaluation.",
          "type": "string"
        },
        "toneScore": {
          "description": "Tone suitability score from 0 to 100.",
          "maximum": 100,
          "minimum": 0,
          "type": "integer"
        }
      },
      "required": [
        "overallScore",
        "clarityScore",
        "structureScore",
        "toneScore",
        "actionabilityScore",
        "level",
        "summary",
        "suggestions"
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
          "topP": 0.9,
          "maxTokens": 700
        }
      },
      {
        "match": {
          "models": [
            "gpt-4o-mini",
            "gpt-4o"
          ]
        },
        "params": {
          "temperature": 0,
          "maxTokens": 650
        },
        "providerParams": {
          "openai": {
            "reasoningEffort": "low"
          }
        }
      }
    ]
  },
  "prompts": {
    "general": "# System\n\nYou are a strict text quality evaluation function. You must only return a JSON object in EXACTLY this format:\n{\"overallScore\": \u003c0-100\u003e, \"clarityScore\": \u003c0-100\u003e, \"structureScore\": \u003c0-100\u003e, \"toneScore\": \u003c0-100\u003e, \"actionabilityScore\": \u003c0-100\u003e, \"level\": \"\u003cexcellent|good|fair|poor\u003e\", \"summary\": \"\u003cone sentence\u003e\", \"suggestions\": [\"...\"]}\n\nDo not output Markdown, do not include any extra explanation, and do not omit any fields.\n\nScoring dimensions:\n- overallScore: holistic usefulness and quality.\n- clarityScore: ease of understanding, precision, and lack of ambiguity.\n- structureScore: organization, flow, and logical sequencing.\n- toneScore: suitability for the target audience and purpose.\n- actionabilityScore: whether the reader can understand what to do or take away.\n\nLevel mapping:\n- excellent: overallScore \u003e= 90\n- good: overallScore \u003e= 75 and \u003c 90\n- fair: overallScore \u003e= 50 and \u003c 75\n- poor: overallScore \u003c 50\n\nsummary: one concise sentence summarizing the overall evaluation result.\n\nsuggestions: up to maxSuggestions concise, actionable improvement tips (default 3 if not specified).\n\n# User\n\nText to evaluate:\n{{text}}\n\nTarget audience:\n{{targetAudience}}\n\nPurpose:\n{{purpose}}\n\nMaximum suggestions:\n{{maxSuggestions}}\n\nStrictness from 1 to 5:\n{{strictness}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-16T00:03:04Z",
    "contentHash": "sha256:79d6f5a90194dede133d2c89da1e7e35c27718e490eea613107036895d42da73"
  }
};

export default artifact;
