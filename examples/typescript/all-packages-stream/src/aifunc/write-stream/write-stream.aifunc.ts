const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "write-stream",
    "version": "1.0.0",
    "description": "Stream long-form writing — articles, reports, or documents — from a prompt and optional structure. Returns plain text.",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.2.0",
    "engineOptions": {
      "injectOutputSchema": false
    }
  },
  "api": {
    "name": "write_stream",
    "description": "Stream long-form writing — articles, reports, or documents — from a prompt and optional structure. Returns plain text.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "audience": {
          "default": "general readers",
          "description": "Target audience (e.g. 'executives', 'engineers', 'general public'). Default: 'general readers'.",
          "type": "string"
        },
        "format": {
          "default": "article",
          "description": "Document format: 'article', 'report', 'proposal', 'documentation', 'essay'. Default: 'article'.",
          "enum": [
            "article",
            "report",
            "proposal",
            "documentation",
            "essay"
          ],
          "type": "string"
        },
        "language": {
          "description": "Output language. If omitted, matches the language of the prompt.",
          "type": "string"
        },
        "prompt": {
          "description": "What to write. Can be a title, a brief, a set of requirements, or a full description of the desired content.",
          "minLength": 1,
          "type": "string"
        },
        "structure": {
          "description": "Optional outline, section headings, or structural notes to follow.",
          "type": "string"
        },
        "tone": {
          "default": "professional",
          "description": "Writing tone (e.g. 'formal', 'professional', 'academic', 'casual'). Default: 'professional'.",
          "type": "string"
        },
        "wordCount": {
          "default": 800,
          "description": "Approximate target word count. Default: 800. Range: 300–5000.",
          "maximum": 5000,
          "minimum": 300,
          "type": "integer"
        }
      },
      "required": [
        "prompt"
      ],
      "type": "object"
    },
    "output": {
      "description": "The generated document as plain text.",
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
          "maxTokens": 6144
        }
      }
    ]
  },
  "prompts": {
    "general": "# System\n\nYou are an expert writer. Your task is to produce a complete, well-structured long-form document based on the given prompt.\n\n## Requirements\n\n- Format: {{input.format}} — follow the conventions and structure appropriate for this document type.\n- Tone: {{input.tone}} — maintain this tone consistently throughout.\n- Target audience: {{input.audience}} — calibrate vocabulary, depth, and assumed knowledge accordingly.\n- Write approximately {{input.wordCount}} words.\n- If a structure or outline is provided, follow it faithfully as the backbone of the document.\n- Begin directly with the document content. Do not include meta-commentary, preamble, or a note about what you are writing.\n- Output plain text only — no Markdown formatting, no JSON, no labels.\n- If a language is specified, write in that language. Otherwise, match the language of the prompt.\n- The document must be coherent, complete, and suitable for its intended purpose without further editing.\n\n## Input\n\nPrompt: {{input.prompt}}\n\nFormat: {{input.format}}\n\nStructure / Outline: {{input.structure}}\n\nTone: {{input.tone}}\n\nAudience: {{input.audience}}\n\nLanguage: {{input.language}}\n\nWord count target: {{input.wordCount}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-16T12:05:18Z",
    "contentHash": "sha256:476d3ea71cfc031bd7beb05d349fea03c14f10d8519430ec7df44ed1dce5d9ec"
  }
};

export default artifact;
