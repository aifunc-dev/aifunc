const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "translate",
    "version": "1.0.0",
    "description": "Translate text into a specified target language with automatic source language detection.",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.2.0"
  },
  "api": {
    "name": "translate",
    "description": "Translate text into a specified target language with automatic source language detection.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "sourceLang": {
          "description": "Source language. If omitted, it will be auto-detected.",
          "type": "string"
        },
        "targetLang": {
          "description": "Target language (e.g. 'English', '日本語', 'zh-CN').",
          "minLength": 1,
          "type": "string"
        },
        "text": {
          "description": "The text to be translated.",
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
      "additionalProperties": false,
      "properties": {
        "sourceLang": {
          "description": "The source language (auto-detected if sourceLang was not provided).",
          "type": "string"
        },
        "translation": {
          "description": "The translated text.",
          "type": "string"
        }
      },
      "required": [
        "translation",
        "sourceLang"
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
          "temperature": 0.3,
          "maxTokens": 2048
        }
      }
    ]
  },
  "prompts": {
    "general": "# System\n\nYou are a professional multilingual translation function. You must only return a JSON object in the following format:\n{\"translation\": \"\u003ctranslated text\u003e\", \"sourceLang\": \"\u003clanguage code\u003e\"}\n\nDo not output Markdown, do not include any extra explanation, and do not add undeclared fields.\n\nRequirements:\n- Accurately translate the source text into the target language, preserving the original semantics and tone.\n- If the user specifies a source language, interpret the text in that language; otherwise, auto-detect the source language.\n- `sourceLang` should use a short language identifier (e.g. zh-CN, en, ja, ko, fr, de).\n- The translation should be natural and fluent, following the conventions of the target language. Avoid stiff word-for-word translation.\n- Preserve proper nouns, brand names, and other content that should not be translated.\n- If the source language is the same as the target language, return the original text as the translation result.\n\n# User\n\nText to translate:\n{{text}}\n\nTarget language: {{targetLang}}\n\nSource language: {{sourceLang}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-16T00:03:04Z",
    "contentHash": "sha256:012dc373afe5d5cca4da2bbab66a253096d028702875fafbb526c2f7357fa25d"
  }
};

export default artifact;
