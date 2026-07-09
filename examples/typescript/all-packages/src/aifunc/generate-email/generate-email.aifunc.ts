const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "generate-email",
    "version": "1.0.0",
    "description": "Generate a complete email from a brief description of intent and context.",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.1.0"
  },
  "api": {
    "name": "generate-email",
    "description": "Generate a complete email from a brief description of intent and context.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "intent": {
          "description": "What the email should accomplish.",
          "minLength": 1,
          "type": "string"
        },
        "keyPoints": {
          "description": "Specific points or details to include in the email body.",
          "items": {
            "type": "string"
          },
          "type": "array"
        },
        "language": {
          "description": "Email language. Default: 'English'.",
          "type": "string"
        },
        "recipientName": {
          "description": "Name or role of the recipient, used in the greeting.",
          "type": "string"
        },
        "senderName": {
          "description": "Name of the sender, used in the sign-off.",
          "type": "string"
        },
        "tone": {
          "description": "Desired tone: 'formal', 'friendly', 'assertive'. Default: 'formal'.",
          "type": "string"
        }
      },
      "required": [
        "intent"
      ],
      "type": "object"
    },
    "output": {
      "additionalProperties": false,
      "properties": {
        "body": {
          "description": "Full email body including greeting and sign-off.",
          "type": "string"
        },
        "subject": {
          "description": "Suggested email subject line.",
          "type": "string"
        }
      },
      "required": [
        "subject",
        "body"
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
          "maxTokens": 2048
        }
      }
    ]
  },
  "prompts": {
    "general": "# System\n\nYou are an email generation function. You must only return a JSON object in the following format:\n{\"subject\": \"\u003csubject line\u003e\", \"body\": \"\u003cfull email body\u003e\"}\n\nDo not output Markdown, do not include any extra explanation, and do not add undeclared fields.\n\nRequirements:\n- Write a complete, professional email that accomplishes the stated intent.\n- Apply the requested tone (default: formal).\n- Use recipientName in the greeting if provided.\n- Use senderName in the sign-off if provided.\n- Incorporate all keyPoints naturally into the body.\n- Write in the requested language (default: English).\n- The subject should be concise and descriptive.\n\n# User\n\nIntent: {{intent}}\n\nTone: {{tone}}\n\nSender name: {{senderName}}\n\nRecipient name: {{recipientName}}\n\nKey points: {{keyPoints}}\n\nLanguage: {{language}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-09T15:12:49Z",
    "contentHash": "sha256:cde93eb8df15df168dc26ad673e60cbf8b4616c6b6b79c25a0ec0e0e762504f1"
  }
};

export default artifact;
