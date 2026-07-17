const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "answer-stream",
    "version": "1.0.0",
    "description": "Stream a detailed answer to a question, optionally grounded in provided context for RAG use cases. Returns plain text.",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.2.0",
    "engineOptions": {
      "injectOutputSchema": false
    }
  },
  "api": {
    "name": "answer_stream",
    "description": "Stream a detailed answer to a question, optionally grounded in provided context for RAG use cases. Returns plain text.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "audience": {
          "default": "general",
          "description": "Target audience (e.g. 'general', 'technical', 'expert'). Default: 'general'.",
          "type": "string"
        },
        "context": {
          "description": "Optional source text, documents, or retrieved passages to ground the answer in. If provided, the answer must be based on this context.",
          "type": "string"
        },
        "depth": {
          "default": "detailed",
          "description": "Answer depth: 'concise' (1-2 paragraphs), 'detailed' (thorough explanation with examples). Default: 'detailed'.",
          "enum": [
            "concise",
            "detailed"
          ],
          "type": "string"
        },
        "language": {
          "description": "Answer language. If omitted, matches the language of the question.",
          "type": "string"
        },
        "question": {
          "description": "The question to answer.",
          "minLength": 1,
          "type": "string"
        }
      },
      "required": [
        "question"
      ],
      "type": "object"
    },
    "output": {
      "description": "The answer as plain text.",
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
          "temperature": 0.3,
          "maxTokens": 4096
        }
      }
    ]
  },
  "prompts": {
    "general": "# System\n\nYou are a knowledgeable and precise question-answering assistant. Your task is to answer the given question accurately and helpfully.\n\n## Requirements\n\n- Audience: {{input.audience}} — calibrate depth, vocabulary, and assumed knowledge accordingly.\n- Depth: {{input.depth}}\n  - \"concise\": answer in 1-2 focused paragraphs, covering just the key points.\n  - \"detailed\": provide a thorough answer with explanation, reasoning, and examples where helpful.\n- If context is provided, base your answer strictly on that context. Do not introduce information not present in it. If the context does not contain enough information to answer, say so clearly.\n- If no context is provided, answer from general knowledge.\n- Begin directly with the answer. Do not restate the question or add preamble.\n- Output plain text only — no Markdown formatting, no JSON, no labels.\n- If a language is specified, answer in that language. Otherwise, match the language of the question.\n- Be accurate, direct, and avoid unnecessary filler or hedging.\n\n## Context (source documents)\n\n{{input.context}}\n\n## Question\n\n{{input.question}}\n\nAudience: {{input.audience}}\n\nLanguage: {{input.language}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-16T12:05:18Z",
    "contentHash": "sha256:cc859b61013ba12f0b9533bd3cd3207eacafb98754fcf7dfbe18c03957554633"
  }
};

export default artifact;
