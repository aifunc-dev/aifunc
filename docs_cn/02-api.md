# 运行时 API 参考

> **目标读者**：应用开发者（消费者）
> **本文内容**：AI 函数在代码中调用时的完整 API 说明，包括配置、调用模式、Mock 与错误处理
> **前置条件**：已通过 `aifn install` 安装了至少一个包

---

## 函数签名

每个 AI 函数包安装后生成一个同名的可调用函数，签名统一为：

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

- `config` 控制运行模式（Mock 或真实调用）及模型连接参数
- `input` / `output` 的类型由包的 `api.json` 定义，IDE 提供完整的类型提示
- Go 函数为同步调用，通过 `context.Context` 支持超时与取消

---

## AIFuncConfig

控制函数的运行模式和模型连接。

### TypeScript

```typescript
interface AIFuncConfig {
  baseURL?: string;
  apiKey?: string;
  model?: string;
  temperature?: number;
  topP?: number;
  maxTokens?: number;
  timeout?: number;     // 毫秒，默认 7000
  maxRetries?: number;  // 默认 1
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
    timeout: float | None = None   # 秒，默认 7.0
    max_retries: int | None = None # 默认 1
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
    Timeout     int    // 毫秒，默认 7000
    MaxRetries  int    // 默认 1
    Mock        bool
    MockData    any
}
```

### 字段说明

| 字段 | 类型 | 默认值 | 说明 |
|:---|:---|:---|:---|
| `baseURL` / `base_url` / `BaseURL` | string | — | 模型 API 端点（OpenAI 兼容格式）。非 Mock 模式下必填 |
| `apiKey` / `api_key` / `APIKey` | string | — | API Key。非 Mock 模式下必填 |
| `model` / `Model` | string | — | 模型名称。非 Mock 模式下必填 |
| `temperature` / `Temperature` | number | 由包定义 | 生效优先级：config → model-params.json → Engine 默认 |
| `topP` / `top_p` | number | 由包定义 | 与 temperature 二选一；生效优先级同上（仅 TS/Python） |
| `maxTokens` / `max_tokens` / `MaxTokens` | integer | 由包定义 | 最大输出 Token 数；生效优先级同上 |
| `timeout` / `Timeout` | number | 7000ms / 7.0s | 请求超时（TS/Go 毫秒，Python 秒）；生效优先级：config → aifunc.json → Engine 默认 |
| `maxRetries` / `max_retries` / `MaxRetries` | integer | 1 | 失败重试次数，耗尽后抛出最后一次错误；生效优先级同上 |
| `mock` / `Mock` | boolean | false | 启用 Mock 模式，不调用真实模型 |

---

## 调用示例

### Mock 模式（离线开发、测试）

```typescript
import { summarize, AIFuncConfig, SummarizeInput } from './aifunc/summarize';

const config: AIFuncConfig = { mock: true };

async function main() {
  const result = await summarize(config, {
    text: "用了三个月，手感和续航都超出预期，非常满意！",
    maxLength: 20
  } as SummarizeInput);
}

main().catch(console.error);
```

Mock 模式不需要 `baseURL`、`apiKey`、`model`，函数从包内置的 mock 数据或自动生成的伪数据返回结果。

### 真实模型调用

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
    // 处理错误
}
fmt.Println(result.Summary)
```

### 覆盖模型参数

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

## Mock 模式详解

设置 `mock: true` 后，Engine 不调用任何模型 API，按以下顺序查找返回值：

1. **精确匹配**：在 mock 数据中查找 `input` 字段与实际输入完全一致的 case，返回其 `output`
2. **Fallback**：查找没有 `input` 字段的 case（作为默认返回值）
3. **自动生成**：都没匹配时，根据 `api.json` 的 output schema 生成零值（string → `""`，number → `0`，boolean → `false`，enum → 第一个值）

Mock 数据来源： 包内置 `mock.json`

---

## 错误处理

函数在以下情况会抛出异常：


| 错误场景         | 错误信息示例                                                        |
| ------------ | ------------------------------------------------------------- |
| 输入不符合 schema | `Input validation failed: ...`                                |
| 缺少必填配置       | `AIFuncConfig.baseURL is required when mock mode is disabled` |
| 模型未指定        | `AIFuncConfig.model is required when mock mode is disabled`   |
| 请求超时         | `Request timeout after 30000ms`                               |
| API 返回非 200  | `Model API returned 429: ...`                                 |
| 模型返回非 JSON   | `Failed to parse model output as JSON: ...`                   |
| 输出不符合 schema | `Output validation failed: ...`                               |


### TypeScript 错误处理

```typescript
try {
  const result = await summarize(config, input);
} catch (error) {
  if (error instanceof Error) {
    console.error(error.message);
  }
}
```

### Python 错误处理

```python
try:
    result = await summarize(config, input)
except Exception as e:
    print(f"AI function error: {e}")
```

### Go 错误处理

```go
result, err := summarize.Summarize(ctx, config, input)
if err != nil {
    fmt.Fprintf(os.Stderr, "AI function error: %v\n", err)
}
```

所有错误均以标准 Error / Exception / `error` 返回，无自定义异常类型。

---

## 兼容的模型服务

任何兼容 OpenAI Chat Completions API (`/chat/completions`) 的服务均可使用，包括各类云端服务和本地部署方案。

本地部署的推理服务通常不需要 `apiKey`，但字段仍需提供（可传任意非空字符串）。

---

## 运行时行为

### 请求流程

```text
输入校验 → 渲染 Prompt → 构建请求 → 调用模型 API → 解析 JSON → 输出校验 → 返回结果
```

### 关键行为说明


| 行为        | 说明                                                        |
| --------- | --------------------------------------------------------- |
| 输入校验      | 根据 `api.json` 的 input schema 校验输入字段和类型                    |
| Prompt 渲染 | 将 `{{input.fieldName}}` 替换为实际输入值                          |
| 输出格式      | 始终要求模型返回 JSON（`response_format: { type: "json_object" }`） |
| 输出校验      | 根据 `api.json` 的 output schema 校验模型返回值                     |
| 重试        | 任何错误均自动重试，次数由 `maxRetries` 控制（默认 1 次）。耗尽重试后抛出最后一次的错误原因    |
| 超时        | 默认 7000ms（Python: 7.0s），可通过 `timeout` 配置                  |


---

## 接下来

- **第一次使用？** → [Quick Start](./01-quick-start)
- **查看 CLI 命令？** → [CLI 命令参考](./03-cli)
- **想创建自己的包？** → [创建 AIFunc 包](./05-create-package)
- **了解内部机制？** → [工作原理](./04-how-it-works)

