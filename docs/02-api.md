# Runtime API Reference

> **Target audience**: Application developers (consumers)
> **Content**: Complete API documentation for calling AI functions in code, including configuration, call patterns, mock mode, and error handling
> **Prerequisites**: At least one package installed via `aifn install`

---

## Function Signature

Each installed AI function package generates a callable function with the same name, with a unified signature:

### TypeScript

```typescript
async function <functionName>(config: AIFuncConfig, input: <Input>): Promise<<Output>>
```

### Python

```python
async def <function_name>(config: AIFuncConfig, input: <Input>) -> <Output>
```

### Go

```go
func <FunctionName>(ctx context.Context, config *AIFuncConfig, input <Input>) (<Output>, error)
```

### Java

```java
public static <Output> <ClassName>.<Output> <methodName>(AIFuncConfig config, <ClassName>.<Input> input)
        throws AIFuncException
```

- `config` controls the runtime mode (mock or real call) and model connection parameters
- `input` / `output` types are defined by the package's `api.json`, with full IDE type hints
- Go functions are synchronous and accept a `context.Context` for timeout and cancellation
- Java methods are synchronous; the `AIFuncConfig` builder pattern provides fluent configuration

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
  topP?: number;
  maxTokens?: number;
  timeout?: number;      // milliseconds, default 7000
  maxRetries?: number;   // default 1, retries on 429/5xx/network errors only
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
    top_p: float | None = None
    max_tokens: int | None = None
    timeout: float | None = None    # seconds, None = use aifunc.json or engine default (7.0)
    max_retries: int | None = None  # None = use aifunc.json or engine default (1)
    mock: bool = False
    mock_data: Any = None
```

### Go

```go
type AIFuncConfig struct {
    BaseURL     string
    APIKey      string
    Model       string
    Temperature *float64
    MaxTokens   *int
    Timeout     int    // milliseconds, 0 = use aifunc.json or engine default (7000)
    MaxRetries  int    // 0 = use aifunc.json or engine default (1)
    Mock        bool
    MockData    any
}
```

### Java

```java
AIFuncConfig config = AIFuncConfig.builder()
    .baseUrl("https://your-api-endpoint/v1")
    .apiKey("your-api-key")
    .model("your-model-name")
    .temperature(0.2)
    .topP(0.9)
    .maxTokens(300)
    .timeoutMs(7000)
    .maxRetries(1)
    .mock(true)
    .mockData(null)
    .build();
```

### Field Descriptions

| Field | Type | Default | Description |
|:---|:---|:---|:---|
| `baseURL` / `base_url` / `BaseURL` | string | — | Model API endpoint (OpenAI-compatible format). Required when mock mode is disabled |
| `apiKey` / `api_key` / `APIKey` | string | — | API Key. Required when mock mode is disabled |
| `model` / `Model` | string | — | Model name. Required when mock mode is disabled |
| `temperature` / `Temperature` | number | Defined by package | Priority: config → model-params.json → engine default |
| `topP` / `top_p` | number | Defined by package | Use instead of temperature for nucleus sampling; same priority (TS/Python only) |
| `maxTokens` / `max_tokens` / `MaxTokens` | integer | Defined by package | Maximum output token count; same priority |
| `timeout` / `Timeout` | number | 7000ms / 7.0s | Request timeout (TS/Go in ms, Python in seconds); priority: config → aifunc.json → engine default |
| `maxRetries` / `max_retries` / `MaxRetries` | integer | 1 | Retry attempts on any failure, throws last error when exhausted; same priority |
| `mock` / `Mock` | boolean | false | Enable mock mode, skips real model calls |




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

```go
import (
    "context"
    "your-module/aifunc/summarize"
)

config := &summarize.AIFuncConfig{
    BaseURL: "https://your-api-endpoint/v1",
    Model:   "your-model-name",
    APIKey:  "your-api-key",
}

result, err := summarize.Summarize(context.Background(), config, summarize.SummarizeInput{
    Text:      "...",
    MaxLength: 50,
})
if err != nil {
    // handle error
}
fmt.Println(result.Summary)
```

```java
import aifunc.summarize.Summarize;
import aifunc.summarize.SummarizeTypes.SummarizeInput;
import aifunc.summarize.SummarizeTypes.SummarizeOutput;
import aifunc._engine.java.v0_1_0.Types.AIFuncConfig;

AIFuncConfig config = AIFuncConfig.builder()
        .baseUrl("https://your-api-endpoint/v1")
        .model("your-model-name")
        .apiKey("your-api-key")
        .build();

SummarizeOutput result = Summarize.summarize(config, new SummarizeInput("...", 50));
System.out.println(result.getSummary());
```

### Overriding Model Parameters

```typescript
const config: AIFuncConfig = {
  baseURL: 'http://localhost:11434/v1',
  model: 'your-local-model',
  temperature: 0.0,
  maxTokens: 200,
  timeout: 60000,
  maxRetries: 0,
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


| Error Scenario              | Example Error Message                                         |
| --------------------------- | ------------------------------------------------------------- |
| Input doesn't match schema  | `Input validation failed: ...`                                |
| Missing required config     | `AIFuncConfig.baseURL is required when mock mode is disabled` |
| Model not specified         | `AIFuncConfig.model is required when mock mode is disabled`   |
| Request timeout             | `Request timeout after 30000ms`                               |
| API returns non-200         | `Model API returned 429: ...`                                 |
| Model returns non-JSON      | `Failed to parse model output as JSON: ...`                   |
| Output doesn't match schema | `Output validation failed: ...`                               |


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

### Go Error Handling

```go
result, err := summarize.Summarize(ctx, config, input)
if err != nil {
    fmt.Fprintf(os.Stderr, "AI function error: %v\n", err)
}
```

### Java Error Handling

```java
try {
    SummarizeOutput result = Summarize.summarize(config, input);
} catch (AIFuncException e) {
    System.err.println("AI function error: " + e.getMessage());
}
```

All errors are thrown as standard Error / Exception / `error` / `AIFuncException` — no custom exception hierarchy beyond `AIFuncException extends RuntimeException`.

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


| Behavior          | Description                                                                                                                  |
| ----------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| Input validation  | Validates input fields and types against the `api.json` input schema                                                         |
| Prompt rendering  | Replaces `{{input.fieldName}}` with actual input values                                                                      |
| Output format     | Always requires model to return JSON (`response_format: { type: "json_object" }`)                                            |
| Output validation | Validates model response against the `api.json` output schema                                                                |
| Retry             | Retries on any error. Retry count controlled by `maxRetries` (default 1). After exhausting retries, the last error is thrown |
| Timeout           | Default 7000ms (Python: 7.0s), configurable via `timeout`                                                                    |


---

## Next Steps

- **First time using AIFunc?** → [Quick Start](./01-quick-start)
- **View CLI commands?** → [CLI Command Reference](./03-cli)
- **Want to create your own package?** → [Create an AIFunc Package](./05-create-package)
- **Understand the internals?** → [How It Works](./04-how-it-works)

