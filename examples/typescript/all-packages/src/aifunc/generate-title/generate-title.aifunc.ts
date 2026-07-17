const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "generate-title",
    "version": "1.0.0",
    "description": "Generate title or headline candidates for a piece of content.",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.2.0"
  },
  "api": {
    "name": "generate-title",
    "description": "Generate title or headline candidates for a piece of content.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "content": {
          "description": "The text, summary, or topic to generate titles for.",
          "minLength": 1,
          "type": "string"
        },
        "count": {
          "description": "Number of title candidates to generate. Default: 3.",
          "maximum": 10,
          "minimum": 1,
          "type": "integer"
        },
        "maxLength": {
          "description": "Maximum character count per title. Default: 80.",
          "type": "integer"
        },
        "style": {
          "description": "Title style: 'neutral', 'clickbait', 'seo', 'academic'. Default: 'neutral'.",
          "type": "string"
        }
      },
      "required": [
        "content"
      ],
      "type": "object"
    },
    "output": {
      "additionalProperties": false,
      "properties": {
        "titles": {
          "description": "Generated title candidates, ordered from most to least recommended.",
          "items": {
            "type": "string"
          },
          "type": "array"
        }
      },
      "required": [
        "titles"
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
          "maxTokens": 512
        }
      }
    ]
  },
  "prompts": {
    "general": "# System\n\nYou are a title generation function. You must only return a JSON object in the following format:\n{\"titles\": [\"\u003ctitle 1\u003e\", \"\u003ctitle 2\u003e\", ...]}\n\nDo not output Markdown, do not include any extra explanation, and do not add undeclared fields.\n\nRequirements:\n- Generate the requested number of title candidates (default: 3).\n- Order titles from most to least recommended.\n- Apply the requested style (default: neutral). 'clickbait' uses curiosity-driven language; 'seo' front-loads keywords; 'academic' is formal and descriptive.\n- Each title must not exceed maxLength characters (default: 80).\n- Titles should be distinct from each other — vary wording and angle.\n\n# User\n\nContent:\n{{content}}\n\nStyle: {{style}}\n\nCount: {{count}}\n\nMax length per title: {{maxLength}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-16T00:03:04Z",
    "contentHash": "sha256:22ec3881fab0d05cf08324780aa9ed1dcddb7ff0cbfca1ab00cb289c0dd35ce4"
  }
};

export default artifact;
