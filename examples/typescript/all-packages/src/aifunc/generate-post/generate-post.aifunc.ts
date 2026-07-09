const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "generate-post",
    "version": "1.0.0",
    "description": "Generate a social media post or short article from a topic or brief.",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.1.0"
  },
  "api": {
    "name": "generate-post",
    "description": "Generate a social media post or short article from a topic or brief.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "includeHashtags": {
          "description": "Whether to append relevant hashtags. Default: false.",
          "type": "boolean"
        },
        "maxLength": {
          "description": "Maximum character count.",
          "maximum": 2000,
          "minimum": 10,
          "type": "integer"
        },
        "platform": {
          "description": "Target platform: 'twitter', 'linkedin', 'instagram', 'general'. Default: 'general'.",
          "type": "string"
        },
        "tone": {
          "description": "Desired tone (e.g. 'professional', 'casual', 'inspirational'). Default: 'casual'.",
          "type": "string"
        },
        "topic": {
          "description": "The subject or key idea of the post.",
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
      "additionalProperties": false,
      "properties": {
        "charCount": {
          "description": "Character count of the generated post.",
          "type": "integer"
        },
        "hashtags": {
          "description": "Suggested hashtags (empty if includeHashtags is false).",
          "items": {
            "type": "string"
          },
          "type": "array"
        },
        "post": {
          "description": "The generated post content.",
          "type": "string"
        }
      },
      "required": [
        "post",
        "hashtags",
        "charCount"
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
    "general": "# System\n\nYou are a social media post generation function. You must only return a JSON object in the following format:\n{\"post\": \"\u003cpost content\u003e\", \"hashtags\": [\"\u003ctag1\u003e\", \"\u003ctag2\u003e\"], \"charCount\": \u003cinteger\u003e}\n\nDo not output Markdown, do not include any extra explanation, and do not add undeclared fields.\n\nRequirements:\n- Write a compelling post about the given topic.\n- Adapt length and style to the target platform (twitter: ≤280 chars, linkedin: professional medium-length, instagram: engaging with emojis ok, general: flexible).\n- Apply the requested tone. Default is casual.\n- If maxLength is provided, do not exceed it.\n- If includeHashtags is true, add 3–5 relevant hashtags in the 'hashtags' array (without # prefix). Otherwise return an empty array.\n- Set charCount to the character count of the 'post' field only (not including hashtags).\n\n# User\n\nTopic: {{topic}}\n\nPlatform: {{platform}}\n\nTone: {{tone}}\n\nMax length: {{maxLength}}\n\nInclude hashtags: {{includeHashtags}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-09T15:12:49Z",
    "contentHash": "sha256:19bc07226587b2060125840142f8d58cca5bc2a6c8788b284cff285e1dccf76a"
  }
};

export default artifact;
