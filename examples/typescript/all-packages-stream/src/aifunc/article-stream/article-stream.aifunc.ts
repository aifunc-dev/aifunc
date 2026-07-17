const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "article-stream",
    "version": "1.0.0",
    "description": "Stream a full article from a title and optional outline. Returns plain text.",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.2.0",
    "engineOptions": {
      "injectOutputSchema": false
    }
  },
  "api": {
    "name": "article_stream",
    "description": "Stream a full article from a title and optional outline. Returns plain text.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "audience": {
          "default": "general readers",
          "description": "Target audience (e.g. 'general readers', 'developers', 'executives'). Default: 'general readers'.",
          "type": "string"
        },
        "language": {
          "description": "Output language (e.g. 'English', 'Chinese'). If omitted, matches the title language.",
          "type": "string"
        },
        "outline": {
          "description": "Optional outline or key points to cover, one per line or as free text.",
          "type": "string"
        },
        "style": {
          "default": "informational",
          "description": "Writing style or tone (e.g. 'informational', 'opinion', 'tutorial', 'news'). Default: 'informational'.",
          "type": "string"
        },
        "title": {
          "description": "The article title.",
          "minLength": 1,
          "type": "string"
        },
        "wordCount": {
          "default": 600,
          "description": "Approximate target word count. Default: 600. Range: 200–3000.",
          "maximum": 3000,
          "minimum": 200,
          "type": "integer"
        }
      },
      "required": [
        "title"
      ],
      "type": "object"
    },
    "output": {
      "description": "The generated article as plain text.",
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
          "temperature": 0.75,
          "maxTokens": 4096
        }
      }
    ]
  },
  "prompts": {
    "general": "# System\n\nYou are a professional article writer. Your task is to write a complete, well-structured article based on the given title and optional outline.\n\n## Requirements\n\n- Write approximately {{input.wordCount}} words of article prose.\n- Style: {{input.style}} — match the appropriate tone and structure for this style.\n- Target audience: {{input.audience}} — calibrate vocabulary, depth, and assumed knowledge accordingly.\n- If an outline is provided, follow it as the structural backbone. If not, devise a logical structure yourself.\n- Begin directly with the article body. Do not include a title header, preamble, or any meta-commentary.\n- Output plain text only — no Markdown formatting, no JSON, no labels.\n- If a language is specified, write in that language. Otherwise, match the language of the title.\n- The article must be coherent, informative, and flow naturally from introduction to conclusion.\n\n## Input\n\nTitle: {{input.title}}\n\nOutline: {{input.outline}}\n\nStyle: {{input.style}}\n\nAudience: {{input.audience}}\n\nLanguage: {{input.language}}\n\nWord count target: {{input.wordCount}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-16T12:05:18Z",
    "contentHash": "sha256:9a83ec7607ab16702eb54a265fb7d0a9a7f27b8db590992c798e3c381c5fa3f2"
  }
};

export default artifact;
