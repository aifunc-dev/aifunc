const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "explain-stream",
    "version": "1.0.0",
    "description": "Stream a clear explanation of a concept, code snippet, or term. Returns plain text.",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.2.0",
    "engineOptions": {
      "injectOutputSchema": false
    }
  },
  "api": {
    "name": "explain_stream",
    "description": "Stream a clear explanation of a concept, code snippet, or term. Returns plain text.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "audience": {
          "default": "intermediate",
          "description": "Target audience level (e.g. 'beginner', 'intermediate', 'expert', 'non-technical'). Default: 'intermediate'.",
          "type": "string"
        },
        "context": {
          "description": "Optional surrounding context (e.g. the file or system the code belongs to, or the domain the term is used in).",
          "type": "string"
        },
        "depth": {
          "default": "standard",
          "description": "Explanation depth: 'brief' (2-3 sentences), 'standard' (1-2 paragraphs), 'detailed' (full breakdown with examples). Default: 'standard'.",
          "enum": [
            "brief",
            "standard",
            "detailed"
          ],
          "type": "string"
        },
        "language": {
          "description": "Output language. If omitted, matches the language of the topic.",
          "type": "string"
        },
        "topic": {
          "description": "The concept, code snippet, or term to explain.",
          "minLength": 1,
          "type": "string"
        }
      },
      "required": [
        "topic"
      ],
      "type": "object"
    },
    "output": {
      "description": "The explanation as plain text.",
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
          "temperature": 0.4,
          "maxTokens": 2048
        }
      }
    ]
  },
  "prompts": {
    "general": "# System\n\nYou are a knowledgeable and clear technical educator. Your task is to explain a concept, code snippet, or term to the reader.\n\n## Requirements\n\n- Audience level: {{input.audience}} — calibrate your vocabulary, assumed knowledge, and use of jargon accordingly.\n- Depth: {{input.depth}}\n  - \"brief\": 2–3 sentences, just the core idea.\n  - \"standard\": 1–2 paragraphs, covering what it is, why it matters, and a brief example if helpful.\n  - \"detailed\": full breakdown — definition, how it works, why it exists, common use cases, pitfalls, and a concrete example.\n- Begin directly with the explanation. Do not restate the topic as a heading or add any preamble.\n- Output plain text only — no Markdown formatting, no JSON, no labels.\n- If a language is specified, write in that language. Otherwise, match the language of the topic input.\n- Be accurate, concise, and avoid unnecessary filler.\n\n## Input\n\nTopic: {{input.topic}}\n\nContext: {{input.context}}\n\nAudience: {{input.audience}}\n\nDepth: {{input.depth}}\n\nLanguage: {{input.language}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-16T12:05:18Z",
    "contentHash": "sha256:ed6ff39e99f5f11ffd5cc18e0011297ef0eb6fc147182c90cd6c4792a2ca7852"
  }
};

export default artifact;
