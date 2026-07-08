# Runtime API Reference

> **Target audience**: Application developers (consumers)
> **Content**: Complete API documentation for calling AI functions in code, including configuration, call patterns, mock mode, and error handling
> **Prerequisites**: At least one package installed via `aifn install`

---

## Function Signature

Each installed AI function package generates an async function with the same name, with a unified signature:

### TypeScript

```typescript
async function <functionName>(config: AIFuncConfig, input: <Input>): Promise<<Output>>
```

### Python

```python
async def <function_name>(config: AIFuncConfig, input: <Input>) -> <Output>
```

- `config` controls the runtime mode (mock or real call) and model connection parameters
- `input` / `output` types are defined by the package's `api.json`, with full IDE type hints

---

## AIFuncConfig

Controls the function's runtime mode and model connection.

### TypeScript

```typescript
interface AIFuncConfig {
  baseURL?: string;
  apiKey?: string;
  model?: string;
  temperature?: number;
  maxTokens?: number;
  timeout?: number;
  mock?: boolean;
}
```

### Python

```python
@dataclass
class AIFuncConfig:
    base_url: str | None = None
    api_key: str | None = None
    model: str | None = None
    temperature: float | None = None
    max_tokens: int | None = None
    timeout: int = 30000
    mock: bool = False
    mock_data: Any = None
```

### Field Descriptions

| Field | Type | Default | Description |
|:---|:---|:---|:---|
| `baseURL` / `base_url` | string | — | Model API endpoint (OpenAI-compatible format). Required when mock mode is disabled |
| `apiKey` / `api_key` | string | — | API Key. Required when mock mode is disabled |
| `model` | string | — | Model name. Required when mock mode is disabled |
| `temperature` | number | Defined by package | Overrides the value suggested by the package author in `model-params.json` |
| `maxTokens` / `max_tokens` | integer | Defined by package | Maximum output token count, overrides the package's suggested value |
| `timeout` | integer | 30000 | Request timeout in milliseconds |
| `mock` | boolean | false | Enable mock mode, skips real model calls |


### Parameter Priority

Priority order for temperature and maxTokens (highest to lowest):

1. **Values explicitly passed in config** (your code)
2. **Rules matching the current model in the package's model-params.json**
3. **Engine defaults**

---

## Usage Examples

### Mock Mode (offline development, testing)

```typescript
import { summarize, AIFuncConfig, SummarizeInput } from './aifunc/summarize';

const config: AIFuncConfig = { mock: true };

async function main() {
  const result = await summarize(config, {
    text: "Used it for three months, the feel and battery life exceeded expectations, very satisfied!",
    maxLength: 20
  } as SummarizeInput);
}

main().catch(console.error);
```

Mock mode doesn't require `baseURL`, `apiKey`, or `model`. The function returns results from built-in mock data or auto-generated pseudo data.

### Real Model Call

```typescript
import { summarize, AIFuncConfig, SummarizeInput } from './aifunc/summarize';

const config: AIFuncConfig = {
  baseURL: 'https://your-api-endpoint/v1',
  model: 'your-model-name',
  apiKey: 'your-api-key',
};

async function main() {
  const result = await summarize(config, { text: "...", maxLength: 50 } as SummarizeInput);
  console.log(result.summary);
}

main().catch(console.error);
```

```python
import asyncio
from aifunc.summarize import summarize, AIFuncConfig, SummarizeInput

config = AIFuncConfig(
    base_url="https://your-api-endpoint/v1",
    model="your-model-name",
    api_key="your-api-key",
)

async def main():
    result = await summarize(config, SummarizeInput(text="...", max_length=50))
    print(result.summary)

asyncio.run(main())
```

### Overriding Model Parameters

```typescript
const config: AIFuncConfig = {
  baseURL: 'http://localhost:11434/v1',
  model: 'your-local-model',
  temperature: 0.0,
  maxTokens: 200,
  timeout: 60000,
};
```

## Mock Mode Details

When `mock: true` is set, the Engine does not call any model API. It looks for return values in the following order:

1. **Exact match**: Finds a case in mock data where the `input` field exactly matches the actual input, returns its `output`
2. **Fallback**: Finds a case without an `input` field (serves as a default return value)
3. **Auto-generate**: When nothing matches, generates zero values based on the `api.json` output schema (string → `""`, number → `0`, boolean → `false`, enum → first value)

Mock data source: Package built-in `mock.json`

---

## Error Handling

The function throws exceptions in the following scenarios:

| Error Scenario | Example Error Message |
|:---|:---|
| Input doesn't match schema | `Input validation failed: ...` |
| Missing required config | `AIFuncConfig.baseURL is required when mock mode is disabled` |
| Model not specified | `AIFuncConfig.model is required when mock mode is disabled` |
| Request timeout | `Request timeout after 30000ms` |
| API returns non-200 | `Model API returned 429: ...` |
| Model returns non-JSON | `Failed to parse model output as JSON: ...` |
| Output doesn't match schema | `Output validation failed: ...` |

### TypeScript Error Handling

```typescript
try {
  const result = await summarize(config, input);
} catch (error) {
  if (error instanceof Error) {
    console.error(error.message);
  }
}
```

### Python Error Handling

```python
try:
    result = await summarize(config, input)
except Exception as e:
    print(f"AI function error: {e}")
```

All errors are thrown as standard Error / Exception — no custom exception types.

---

## Compatible Model Services

Any service compatible with the OpenAI Chat Completions API (`/chat/completions`) can be used, including various cloud services and local deployment solutions.

Locally deployed inference services typically don't require an `apiKey`, but the field still needs to be provided (you can pass any non-empty string).

---

## Runtime Behavior

### Request Flow

```text
Input validation → Render Prompt → Build request → Call Model API → Parse JSON → Output validation → Return result
```

### Key Behaviors

| Behavior | Description |
|:---|:---|
| Input validation | Validates input fields and types against the `api.json` input schema |
| Prompt rendering | Replaces `{{input.fieldName}}` with actual input values |
| Output format | Always requires model to return JSON (`response_format: { type: "json_object" }`) |
| Output validation | Validates model response against the `api.json` output schema |
| Retry | Current version has no built-in retry mechanism; failures throw immediately |
| Timeout | Default 30 seconds, configurable via `timeout` |

---

## Next Steps

- **First time using AIFunc?** → [Quick Start](./01-quick-start)
- **View CLI commands?** → [CLI Command Reference](./03-cli)
- **Want to create your own package?** → [Create an AIFunc Package](./05-create-package)
- **Understand the internals?** → [How It Works](./04-how-it-works)
