const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "extract-entities",
    "version": "1.0.0",
    "description": "Extract named entities (people, places, organizations, dates, etc.) from text.",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.2.0"
  },
  "api": {
    "name": "extract-entities",
    "description": "Extract named entities from text, identifying their type and position.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "entityTypes": {
          "description": "Optional filter: only extract entities of these types. If omitted, extract all recognized types.",
          "items": {
            "minLength": 1,
            "type": "string"
          },
          "type": "array"
        },
        "text": {
          "description": "The text to extract entities from.",
          "minLength": 1,
          "type": "string"
        }
      },
      "required": [
        "text"
      ],
      "type": "object"
    },
    "output": {
      "additionalProperties": false,
      "properties": {
        "entities": {
          "description": "List of extracted entities.",
          "items": {
            "additionalProperties": false,
            "properties": {
              "end": {
                "description": "End character offset (exclusive) in the input text.",
                "minimum": 0,
                "type": "integer"
              },
              "start": {
                "description": "Start character offset in the input text (0-based).",
                "minimum": 0,
                "type": "integer"
              },
              "text": {
                "description": "The entity text as it appears in the input.",
                "type": "string"
              },
              "type": {
                "description": "Entity type (e.g. 'person', 'location', 'organization', 'date', 'money', 'product').",
                "type": "string"
              }
            },
            "required": [
              "text",
              "type",
              "start",
              "end"
            ],
            "type": "object"
          },
          "type": "array"
        }
      },
      "required": [
        "entities"
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
          "temperature": 0,
          "maxTokens": 1024
        }
      }
    ]
  },
  "prompts": {
    "general": "# System\n\nYou are a named entity recognition (NER) function. You must only return a JSON object in the following format:\n{\"entities\": [{\"text\": \"\u003centity\u003e\", \"type\": \"\u003ctype\u003e\", \"start\": \u003cint\u003e, \"end\": \u003cint\u003e}, ...]}\n\nDo not output Markdown, do not include any extra explanation, and do not add undeclared fields.\n\nRequirements:\n- Extract all named entities from the input text.\n- Common entity types include: person, location, organization, date, time, money, percentage, product, event. Use lowercase type names.\n- If `entityTypes` is provided, only extract entities matching those types. Otherwise, extract all recognized entities.\n- `start` is the 0-based character offset where the entity begins in the input text.\n- `end` is the exclusive character offset where the entity ends (i.e. text[start:end] == entity text).\n- The `text` field must exactly match the substring in the input at the given offsets.\n- Do not overlap entities. If a span could match multiple types, choose the most specific one.\n- Return entities in the order they appear in the text (by `start` offset).\n\n# User\n\nText:\n{{text}}\n\nEntity types to extract: {{entityTypes}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-16T00:03:04Z",
    "contentHash": "sha256:f302a2b25694fb3cc2cbd56901803b9e9e35e12e2339c9a64aee6e7376a9ddb5"
  }
};

export default artifact;
