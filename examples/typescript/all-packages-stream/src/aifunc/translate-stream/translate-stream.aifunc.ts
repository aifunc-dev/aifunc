const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "translate-stream",
    "version": "1.0.0",
    "description": "Stream the translation of a long document or text. Returns plain text.",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.2.0",
    "engineOptions": {
      "injectOutputSchema": false
    }
  },
  "api": {
    "name": "translate_stream",
    "description": "Stream the translation of a long document or text. Returns plain text.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "domain": {
          "description": "Subject domain hint to improve terminology accuracy (e.g. 'legal', 'medical', 'technical', 'literary'). Optional.",
          "type": "string"
        },
        "sourceLang": {
          "description": "Source language. If omitted, it is auto-detected.",
          "type": "string"
        },
        "style": {
          "default": "natural",
          "description": "Translation style: 'literal' (close to source), 'natural' (idiomatic), 'formal'. Default: 'natural'.",
          "enum": [
            "literal",
            "natural",
            "formal"
          ],
          "type": "string"
        },
        "targetLang": {
          "description": "Target language (e.g. 'Chinese', 'English', 'French', 'zh-CN').",
          "minLength": 1,
          "type": "string"
        },
        "text": {
          "description": "The text or document to translate.",
          "minLength": 1,
          "type": "string"
        }
      },
      "required": [
        "text",
        "targetLang"
      ],
      "type": "object"
    },
    "output": {
      "description": "The translated text as plain text.",
      "type": "string",
      "x-delivery-mode": "stream"
    },
    "injectOutputSchema": false
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
          "maxTokens": 8192
        }
      }
    ]
  },
  "prompts": {
    "general": "# System\n\nYou are a professional translator. Your task is to translate the provided text into the target language accurately and fluently.\n\n## Requirements\n\n- Translate into: {{input.targetLang}}\n- Source language: {{input.sourceLang}} (if not specified, detect automatically)\n- Translation style: {{input.style}}\n  - \"literal\": stay close to the source structure and wording\n  - \"natural\": produce idiomatic, fluent prose that reads as if originally written in the target language\n  - \"formal\": use formal register and polished language appropriate for official or professional contexts\n- Domain: {{input.domain}} — if specified, apply domain-specific terminology conventions\n- Output only the translated text — no notes, no commentary, no labels, no original text\n- Preserve paragraph breaks and structural whitespace from the original\n- Do not translate proper names, brand names, or code identifiers unless there is a well-established target-language equivalent\n\n## Input\n\n{{input.text}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-19T11:51:40Z",
    "contentHash": "sha256:a267da318ceb7823ce27f36b267848bc42c2320a0b3dcb98c15acc8a07a6a657"
  }
};

export default artifact;
