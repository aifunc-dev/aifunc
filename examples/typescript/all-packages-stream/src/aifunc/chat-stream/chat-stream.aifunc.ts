const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "chat-stream",
    "version": "1.0.0",
    "description": "Stream a conversational AI reply from a message history. Returns plain text.",
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
    "description": "Stream a conversational AI reply from a message history. Returns plain text.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "language": {
          "description": "Reply language. If omitted, matches the language of the last user message.",
          "type": "string"
        },
        "messages": {
          "description": "Conversation history. Each item has a 'role' ('user' or 'assistant') and 'content' (string).",
          "items": {
            "additionalProperties": false,
            "properties": {
              "content": {
                "description": "Message text.",
                "minLength": 1,
                "type": "string"
              },
              "role": {
                "description": "Message sender role.",
                "enum": [
                  "user",
                  "assistant"
                ],
                "type": "string"
              }
            },
            "required": [
              "role",
              "content"
            ],
            "type": "object"
          },
          "minItems": 1,
          "type": "array"
        },
        "systemPrompt": {
          "description": "Optional system-level instruction that sets the assistant's persona, role, or constraints.",
          "type": "string"
        }
      },
      "required": [
        "messages"
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
    "general": "# System\n\n{{input.systemPrompt}}\n\nYou are a helpful, concise, and friendly conversational assistant.\n\n## Requirements\n\n- Reply naturally to the most recent user message, taking the full conversation history into account.\n- Be direct and helpful. Match the tone and register of the conversation.\n- Do not summarize, repeat, or acknowledge the conversation history explicitly — just reply.\n- Output plain text only — no Markdown formatting, no JSON, no labels.\n- If a language is specified, reply in that language. Otherwise, match the language of the last user message.\n\n## Conversation History\n\n{{input_json}}\n\nLanguage: {{input.language}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-16T12:05:18Z",
    "contentHash": "sha256:2698bc3b0adfb9dc78a26736d907dbb9ca8e9231825b10aed10b7f2d486d8849"
  }
};

export default artifact;
