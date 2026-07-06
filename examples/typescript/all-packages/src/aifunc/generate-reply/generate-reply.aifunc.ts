const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "generate-reply",
    "version": "1.0.0",
    "description": "Generate a contextually appropriate reply to a message or comment.",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.1.0"
  },
  "api": {
    "name": "generate-reply",
    "description": "Generate a contextually appropriate reply to a message or comment.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "context": {
          "description": "Background context to inform the reply (e.g. role, situation).",
          "type": "string"
        },
        "language": {
          "description": "Reply language. If omitted, matches the input message language.",
          "type": "string"
        },
        "message": {
          "description": "The original message or comment to reply to.",
          "minLength": 1,
          "type": "string"
        },
        "tone": {
          "description": "Desired tone: 'friendly', 'formal', 'empathetic', 'concise'. Default: 'friendly'.",
          "type": "string"
        }
      },
      "required": [
        "message"
      ],
      "type": "object"
    },
    "output": {
      "additionalProperties": false,
      "properties": {
        "reply": {
          "description": "The generated reply text.",
          "type": "string"
        }
      },
      "required": [
        "reply"
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
          "temperature": 0.7,
          "maxTokens": 1024
        }
      }
    ]
  },
  "prompts": {
    "general": "# System\n\nYou are a reply generation function. You must only return a JSON object in the following format:\n{\"reply\": \"\u003creply text\u003e\"}\n\nDo not output Markdown, do not include any extra explanation, and do not add undeclared fields.\n\nRequirements:\n- Write a natural, contextually appropriate reply to the given message.\n- Apply the requested tone (default: friendly).\n- If context is provided, use it to tailor the reply.\n- Reply in the requested language; if not specified, match the language of the input message.\n- Keep the reply concise and focused.\n\n# User\n\nOriginal message:\n{{message}}\n\nTone: {{tone}}\n\nContext: {{context}}\n\nLanguage: {{language}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-05T04:22:14Z",
    "contentHash": "sha256:bd15142db6ed0425bb35a1ae82edfeb3700680180b576b8592b4f74ce5b459f1"
  }
};

export default artifact;
