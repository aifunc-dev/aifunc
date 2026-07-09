const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "rewrite",
    "version": "1.0.0",
    "description": "Rewrite text in a specified style or tone, such as formal, casual, concise, or expanded.",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.1.0"
  },
  "api": {
    "name": "rewrite",
    "description": "Rewrite text according to a specified style or instruction while preserving the original meaning.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "instructions": {
          "description": "Additional rewriting instructions or constraints.",
          "type": "string"
        },
        "style": {
          "description": "Target style or tone (e.g. 'formal', 'casual', 'concise', 'expanded', 'academic', 'humorous').",
          "minLength": 1,
          "type": "string"
        },
        "text": {
          "description": "The original text to rewrite.",
          "minLength": 1,
          "type": "string"
        }
      },
      "required": [
        "text",
        "style"
      ],
      "type": "object"
    },
    "output": {
      "additionalProperties": false,
      "properties": {
        "rewritten": {
          "description": "The rewritten text.",
          "type": "string"
        }
      },
      "required": [
        "rewritten"
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
    "general": "# System\n\nYou are a text rewriting function. You must only return a JSON object in the following format:\n{\"rewritten\": \"\u003crewritten text\u003e\"}\n\nDo not output Markdown, do not include any extra explanation, and do not add undeclared fields.\n\nRequirements:\n- Rewrite the input text according to the specified style/tone.\n- Preserve the original meaning and key information — do not add or remove facts.\n- The output language should match the input language unless the style implies otherwise.\n- Common styles include: formal, casual, concise, expanded, academic, humorous, professional, poetic, simplified.\n- If additional instructions are provided, follow them as constraints on the rewrite.\n- Produce natural, fluent text that reads as if originally written in the target style.\n\n# User\n\nOriginal text:\n{{text}}\n\nTarget style: {{style}}\n\nAdditional instructions: {{instructions}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-09T15:12:49Z",
    "contentHash": "sha256:d730cdf6c2da955b293cb369f21c58c745d22e5f1676db19b025d70e152adcbc"
  }
};

export default artifact;
