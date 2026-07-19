const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "chat-stream",
    "version": "1.0.0",
    "description": "Send a message and stream a plain-text reply. Optionally include context such as prior turns.",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.2.0",
    "engineOptions": {
      "injectOutputSchema": false
    }
  },
  "api": {
    "name": "chat_stream",
    "description": "Send a message and stream a plain-text reply. Optionally include context such as prior turns.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "context": {
          "description": "Optional conversation history or other background text the reply should take into account.",
          "type": "string"
        },
        "message": {
          "description": "The user message.",
          "minLength": 1,
          "type": "string"
        }
      },
      "required": [
        "message"
      ],
      "type": "object"
    },
    "output": {
      "description": "The assistant reply as plain text.",
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
          "temperature": 0.7,
          "maxTokens": 2048
        }
      }
    ]
  },
  "prompts": {
    "general": "# System\n\nYou are a helpful, concise, and friendly conversational assistant.\n\n## Requirements\n\n- Reply directly and helpfully to the user message.\n- If context is provided, use it to tailor the reply.\n- Match the tone and language of the message.\n- Output plain text only — no Markdown formatting, no JSON, no labels.\n- Be direct. Do not add preambles like \"Sure!\" or \"Of course!\".\n\n## User Message\n\n{{message}}\n\n## Context\n\n{{context}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-19T11:51:40Z",
    "contentHash": "sha256:4cb6df8aabe5134abf74c79adc839f36b323d9ec4c05dc8f0c1e4b4bb4235b3c"
  }
};

export default artifact;
