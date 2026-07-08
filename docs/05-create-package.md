# Create an AIFunc Package

> **Target audience**: Package authors
> **Content**: Create an AI function package from scratch that can be installed in any project and called like a regular function
> **Prerequisites**: aifn CLI installed ([installation guide](./01-quick-start#step-1-install-cli))

---

## Step 1: Create the Scaffold

```bash
aifn create my-summarizer
```

This generates the following directory:

```text
my-summarizer/
├── package.json       ← Package metadata
├── api.json           ← API interface definition
└── prompts/
    └── general.md     ← Prompt template
```

---

## Step 2: Edit package.json

Fill in the package's basic information:

```json
{
  "type": "standalone",
  "name": "my-summarizer",
  "displayName": "Text Summarization",
  "description": "Generate a concise summary of the input text.",
  "version": "1.0.0",
  "engine": "^0.1.0",
  "author": {
    "name": "YourName"
  },
  "license": "MIT",
  "categories": ["text", "summary", "productivity"],
  "tags": ["summary", "short", "text"]
}
```

See [Package Format Spec - package.json](./spec#3-packagejson) for full field descriptions.

---

## Step 3: Define the Interface in api.json

Think about two things: what the function receives and what it returns.

```json
{
  "version": "1.0.0",
  "name": "summarize",
  "description": "Generate a concise summary of the input text.",
  "input": {
    "type": "object",
    "properties": {
      "text": {
        "type": "string",
        "description": "The text to summarize.",
        "minLength": 1
      },
      "maxLength": {
        "type": "integer",
        "description": "Maximum word count for the summary. Defaults to 80.",
        "minimum": 20,
        "maximum": 300,
        "default": 80
      }
    },
    "required": ["text"],
    "additionalProperties": false
  },
  "output": {
    "type": "object",
    "properties": {
      "summary": {
        "type": "string",
        "description": "The generated summary."
      },
      "wordCount": {
        "type": "integer",
        "description": "Approximate word count of the summary.",
        "minimum": 0
      }
    },
    "required": ["summary", "wordCount"],
    "additionalProperties": false
  }
}
```

The schema follows JSON Schema Draft 2020-12, supporting `string`, `number`, `integer`, `boolean`, `object`, and `array` types.

---

## Step 4: Write the Prompt in prompts/general.md

Tell the AI how to complete the task:

```markdown
# System

You are a concise and accurate summarization assistant.

Requirements:
- The summary language must match the language of the input text — do not translate.
- The summary should be concise, accurate, and fluent.
- Preserve the most essential information; do not fabricate anything not present in the original.
- If the input covers multiple points, prioritize the most important 1 to 3.
- `summary` must not exceed the word count specified by `maxLength`; default to 80 if not provided.
- `wordCount` should reflect the approximate length of the summary.

# User

Text:
{{input.text}}

Maximum length:
{{input.maxLength}}
```

By default, you don't need to write "please return JSON" or similar format instructions in the prompt — the Engine automatically injects output format directives based on the `api.json` output schema.

If you want full control over the output format, you can disable auto-injection in `package.json`:

```json
{
  "engineOptions": {
    "injectOutputSchema": false
  }
}
```

When disabled, the Engine won't compile the output schema into the prompt, and the package author must explicitly describe how the model should output in the prompt.

### Template Variables

| Syntax | Description |
|:---|:---|
| `{{input.fieldName}}` | Insert the specified field value from the input object |
| `{{input_json}}` | Serialize the entire input as a JSON string and insert |
| `{{input}}` | When input is a string type, inserts the raw text; when object, equivalent to `{{input_json}}` |

---

## Step 5: Validate the Package Format

```bash
aifn validate ./my-summarizer
```

On success, outputs the package name, version, description, engine version, and function name. On errors, it clearly indicates which fields are missing or non-compliant.

---

## Step 6: Install Locally and Test

Install this local package in your project:

```bash
aifn install ./my-summarizer
```

Then use it in your code:

```typescript
import { summarize, AIFuncConfig } from './aifunc/my-summarizer';

const config: AIFuncConfig = { mock: true };

async function main() {
  const result = await summarize(config, {
    text: "Artificial intelligence is transforming industries worldwide. From healthcare to finance, AI-powered solutions are improving efficiency and enabling new capabilities that were previously impossible.",
    maxLength: 30
  });

  console.log(result.summary);    // "AI is transforming industries like healthcare and finance by improving efficiency."
  console.log(result.wordCount);  // 12
}

main();
```

Local packages are linked via `file:` path. After modifying the package source files, re-run `aifn install` to update the artifacts.

---

## Optional: Add Model Parameter Suggestions

If you want to suggest a lower temperature for more stable outputs, create `model-params.json`:

```json
{
  "presets": [
    {
      "match": { "pattern": ".*" },
      "params": { "temperature": 0.1, "maxTokens": 300 }
    }
  ]
}
```

Parameter priority (lowest to highest): Engine defaults → this file's config → values explicitly passed in calling code.

---

## Optional: Add Mock Data

Create `mock.json` to provide static responses for testing, useful for offline development and CI tests:

```json
{
  "version": "1.0.0",
  "delay": {
    "minMs": 30,
    "maxMs": 100
  },
  "random": {
    "enabled": false,
    "seed": "summarize"
  },
  "cases": [
    {
      "id": "basic-summary",
      "description": "Basic summary example.",
      "output": {
        "summary": "After three months of use, very satisfied with the feel and battery life.",
        "wordCount": 20
      }
    }
  ]
}
```

---

## Tips for Writing Good Prompts

| Tip | Description |
|:---|:---|
| Define the role | Start by clearly stating "you are what kind of expert" |
| List rules | Use a checklist to specify quality, length, and boundary requirements — don't be vague |
| Set boundaries | Explain how to handle edge cases (e.g., "if input is too short, summarize the core message directly") |
| Don't write format requirements | Output format is auto-injected by the Engine — focus on task logic only |
| Put variables last | Place `{{input.text}}` at the end of the prompt so the model can locate input content easily |

---

## Complete Package Directory Reference

Minimal package (required files only):

```text
my-package/
├── package.json
├── api.json
└── prompts/
    └── general.md
```

Complete package with optional files:

```text
my-package/
├── package.json
├── api.json
├── prompts/
│   └── general.md
├── model-params.json
├── mock.json
├── README.md
└── LICENSE
```

---

## Next Steps

| Goal | Documentation |
|:---|:---|
| Share your package or publish to a public repository | [Sharing & Publishing](./07-sharing) |
| Manage private packages within a team | [Team Workflow](./08-team-workflow) |
| View the complete field definitions | [Package Format Spec](./06-spec) |
| View CLI command reference | [CLI Command Reference](./03-cli) |
