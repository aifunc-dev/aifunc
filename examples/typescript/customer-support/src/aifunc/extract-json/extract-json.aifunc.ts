const artifact = {
  "schemaVersion": "0.1.0",
  "artifactVersion": "0.1.0",
  "package": {
    "type": "standalone",
    "name": "extract-json",
    "version": "1.0.0",
    "description": "Extract structured JSON data from natural language text based on a user-defined field schema.",
    "author": {
      "name": "GildenEye"
    },
    "engine": "^0.1.0"
  },
  "api": {
    "name": "extract-json",
    "description": "Extract structured JSON from natural language text according to a user-defined field schema.",
    "input": {
      "additionalProperties": false,
      "properties": {
        "fields": {
          "description": "Schema describing the fields to extract.",
          "items": {
            "additionalProperties": false,
            "properties": {
              "description": {
                "description": "What this field represents, to guide extraction.",
                "type": "string"
              },
              "name": {
                "description": "Field name (will be the key in the output JSON).",
                "minLength": 1,
                "type": "string"
              },
              "type": {
                "description": "Expected value type.",
                "enum": [
                  "string",
                  "number",
                  "boolean",
                  "array",
                  "object"
                ],
                "type": "string"
              }
            },
            "required": [
              "name",
              "description",
              "type"
            ],
            "type": "object"
          },
          "minItems": 1,
          "type": "array"
        },
        "text": {
          "description": "The natural language text to extract information from.",
          "minLength": 1,
          "type": "string"
        }
      },
      "required": [
        "text",
        "fields"
      ],
      "type": "object"
    },
    "output": {
      "additionalProperties": false,
      "properties": {
        "extracted": {
          "additionalProperties": true,
          "description": "Extracted key-value pairs matching the requested fields.",
          "type": "object"
        },
        "missing": {
          "description": "Field names that could not be found in the text.",
          "items": {
            "type": "string"
          },
          "type": "array"
        }
      },
      "required": [
        "extracted",
        "missing"
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
    "general": "# System\n\nYou are a structured information extraction function. You must only return a JSON object in the following format:\n{\"extracted\": {\u003cfieldName\u003e: \u003cvalue\u003e, ...}, \"missing\": [\u003cfieldName\u003e, ...]}\n\nDo not output Markdown, do not include any extra explanation, and do not add undeclared fields.\n\nRequirements:\n- Extract information from the text based on the field definitions provided.\n- Each field has a `name`, `description` (what to look for), and `type` (the expected value type).\n- Place successfully extracted values in `extracted` using the field name as key.\n- Values must match the declared type: \"string\" -\u003e string, \"number\" -\u003e number, \"boolean\" -\u003e true/false, \"array\" -\u003e JSON array, \"object\" -\u003e JSON object.\n- If a field's value cannot be determined from the text, do NOT guess — add the field name to the `missing` array and omit it from `extracted`.\n- Do not invent information that is not present or clearly implied in the text.\n- For \"array\" fields, extract all relevant items mentioned in the text.\n- For \"number\" fields, parse numeric values (e.g. \"five years\" -\u003e 5).\n\n# User\n\nText:\n{{text}}\n\nFields to extract:\n{{fields}}\n"
  },
  "metadata": {
    "sourcePackageVersion": "1.0.0",
    "generatedAt": "2026-07-05T03:58:10Z",
    "contentHash": "sha256:9da2d18f5fb3a155da515c417aeb72b33f7f487601ae13723fb1d733ccda158d"
  }
};

export default artifact;
