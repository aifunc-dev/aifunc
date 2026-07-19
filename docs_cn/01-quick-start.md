
# Quick Start

> **目标读者**：所有想使用 AIFunc 的开发者
> **本文内容**：5 分钟完成从安装 CLI 到调用 AI 函数拿到结果的完整流程
> **前置条件**：Node.js 18+、Python 3.10+、Go 1.23+、Java 11+ 或 .NET 6+

---

## Step 1：安装 CLI

```bash
# 安装 CLI
brew tap aifunc-dev/aifn && brew install aifn   # macOS/Linux
scoop bucket add aifn https://github.com/aifunc-dev/scoop-aifn && scoop install aifn  # Windows
```

---

## Step 2：安装 AI 函数包

在项目根目录执行：

```bash
aifn install github:aifunc-dev/aifunc-packages/summarize
```

CLI 会自动：
- 识别项目类型（TypeScript / Python / Go / Java / C#）
- 生成可直接 import 的代码（含类型定义和内置 mock 数据）
- 创建配置文件（如果不存在）

---

## Step 3：写代码调用

### TypeScript

```typescript
import { summarize, AIFuncConfig, SummarizeInput } from './aifunc/summarize';

const config: AIFuncConfig = { mock: true };

const text =
  'The James Webb Space Telescope captured its first full-color images in July 2022, ' +
  'revealing thousands of galaxies in a single image.';

async function main() {
  const result = await summarize(config, { text, maxLength: 30 } as SummarizeInput);
  console.log(`Summary   : ${result.summary}`);
  console.log(`Word count: ${result.wordCount}`);
}

main().catch(console.error);
```

### Python

```python
import asyncio
from aifunc.summarize import summarize, AIFuncConfig, SummarizeInput

config = AIFuncConfig(mock=True)

text = (
    'The James Webb Space Telescope captured its first full-color images in July 2022, '
    'revealing thousands of galaxies in a single image.'
)

async def main():
    result = await summarize(config, SummarizeInput(text=text, max_length=30))
    print(f"Summary   : {result.summary}")
    print(f"Word count: {result.word_count}")

asyncio.run(main())
```

### Go

```go
package main

import (
	"context"
	"fmt"
	"log"

	"your-module/aifunc/summarize"
)

func main() {
	config := &summarize.AIFuncConfig{Mock: true}

	text := "The James Webb Space Telescope captured its first full-color images in July 2022, " +
		"revealing thousands of galaxies in a single image."

	maxLen := 30
	result, err := summarize.Summarize(context.Background(), config, summarize.SummarizeInput{
		Text:      text,
		MaxLength: &maxLen,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Summary   : %s\n", result.Summary)
	fmt.Printf("Word count: %d\n", result.WordCount)
}
```

IDE 提供完整的类型提示和自动补全。

### Java

```java
import aifunc.AIFuncConfig;
import aifunc.summarize.Summarize;
import aifunc.summarize.SummarizeTypes.SummarizeInput;

AIFuncConfig config = AIFuncConfig.builder().mock(true).build();

String text = "The James Webb Space Telescope captured its first full-color images in July 2022, " +
              "revealing thousands of galaxies in a single image.";

Summarize.summarize(config, new SummarizeInput(text, 30))
        .thenAccept(result -> {
            System.out.println("Summary   : " + result.getSummary());
            System.out.println("Word count: " + result.getWordCount());
        })
        .join();
```

### C#

```csharp
using Aifunc;
using Aifunc.Summarize;

var config = new AIFuncConfig { Mock = true };

var text =
    "The James Webb Space Telescope captured its first full-color images in July 2022, " +
    "revealing thousands of galaxies in a single image.";

var result = await Summarize.SummarizeAsync(config, new SummarizeTypes.SummarizeInput(text, 30));
Console.WriteLine($"Summary   : {result.Summary}");
Console.WriteLine($"Word count: {result.WordCount}");
```

---

## Step 4：连接真实模型

上面的示例使用 `mock: true`，无需 API Key 即可跑通完整流程。

准备连接真实模型时，config配置真实参数：

### TypeScript

```typescript
import { summarize, AIFuncConfig, SummarizeInput } from './aifunc/summarize';

const config: AIFuncConfig = {
  baseURL: 'https://your-api-endpoint/v1',
  model: 'your-model-name',
  apiKey: 'your-api-key',
};

const text =
  'The James Webb Space Telescope captured its first full-color images in July 2022, ' +
  'revealing thousands of galaxies in a single image.';

async function main() {
  const result = await summarize(config, { text, maxLength: 30 } as SummarizeInput);
  console.log(`Summary   : ${result.summary}`);
  console.log(`Word count: ${result.wordCount}`);
}

main().catch(console.error);
```

### Python

```python
import asyncio
from aifunc.summarize import summarize, AIFuncConfig, SummarizeInput

config = AIFuncConfig(
    base_url="https://your-api-endpoint/v1",
    model="your-model-name",
    api_key="your-api-key",
)

text = (
    'The James Webb Space Telescope captured its first full-color images in July 2022, '
    'revealing thousands of galaxies in a single image.'
)

async def main():
    result = await summarize(config, SummarizeInput(text=text, max_length=30))
    print(f"Summary   : {result.summary}")
    print(f"Word count: {result.word_count}")

asyncio.run(main())
```

### Go

```go
package main

import (
	"context"
	"fmt"
	"log"

	"your-module/aifunc/summarize"
)

func main() {
	config := &summarize.AIFuncConfig{
		BaseURL: "https://your-api-endpoint/v1",
		Model:   "your-model-name",
		APIKey:  "your-api-key",
	}

	text := "The James Webb Space Telescope captured its first full-color images in July 2022, " +
		"revealing thousands of galaxies in a single image."

	maxLen := 30
	result, err := summarize.Summarize(context.Background(), config, summarize.SummarizeInput{
		Text:      text,
		MaxLength: &maxLen,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Summary   :", result.Summary)
	fmt.Println("Word count:", result.WordCount)
}
```

### Java

```java
import aifunc.AIFuncConfig;
import aifunc.summarize.Summarize;
import aifunc.summarize.SummarizeTypes.SummarizeInput;

AIFuncConfig config = AIFuncConfig.builder()
        .baseUrl("https://your-api-endpoint/v1")
        .model("your-model-name")
        .apiKey("your-api-key")
        .build();

String text = "The James Webb Space Telescope captured its first full-color images in July 2022, " +
              "revealing thousands of galaxies in a single image.";

Summarize.summarize(config, new SummarizeInput(text, 30))
        .thenAccept(result -> {
            System.out.println("Summary   : " + result.getSummary());
            System.out.println("Word count: " + result.getWordCount());
        })
        .join();
```

### C#

```csharp
using Aifunc;
using Aifunc.Summarize;

var config = new AIFuncConfig
{
    BaseUrl = "https://your-api-endpoint/v1",
    Model = "your-model-name",
    ApiKey = "your-api-key",
};

var text =
    "The James Webb Space Telescope captured its first full-color images in July 2022, " +
    "revealing thousands of galaxies in a single image.";

var result = await Summarize.SummarizeAsync(config, new SummarizeTypes.SummarizeInput(text, 30));
Console.WriteLine($"Summary   : {result.Summary}");
Console.WriteLine($"Word count: {result.WordCount}");
```

支持任何兼容 OpenAI 协议的服务端点。

运行代码：

```bash
# TypeScript：安装依赖后编译运行
npm install
npm run build
npm run start

# 或使用 dev 模式（ts-node + 路径别名解析）
npm run dev

# Python
python main.py

# Go
go run main.go

# Java
javac Main.java
java Main

# C#
dotnet run
```

---

## 完成

就像使用普通 npm/pip 包一样：安装 → import → 调用 → 拿结果。

无需接触 Prompt、Token、Temperature 这些概念。

---

## 关于 Mock 模式

每个 AI 函数包自带 `mock.json`，定义了预设的输入输出映射。在代码中设置 `mock: true` 即可激活：

- 优先匹配 mock 数据中与输入完全一致的 case
- 无匹配时，根据输出 schema 自动生成符合类型的返回值
- 适用于离线开发、单元测试、CI 流水线

无需环境变量，完全由代码控制。

---

## 接下来

- **查看完整 API 参考？** → [运行时 API 参考](./02-api)
- **创建自己的 AI 函数？** → [创建包教程](./05-create-package)
- **发布包给团队？** → [分享与发布](./07-sharing)

---

## 常见问题

<details>
<summary><b>怎么拿到流式输出？</b></summary>

使用官方名称后缀为 `-stream` 的包（如 `chat-stream`、`answer-stream`）。流式由包定义声明——见协议中的 [`x-delivery-mode`](./06-spec.md#output-扩展字段)。

先安装 `chat-stream`：

```bash
aifn install github:aifunc-dev/aifunc-packages/chat-stream
```

各语言调用示例：

<details open>
<summary>Python</summary>

```python
import asyncio, sys
from aifunc.chat_stream import chat_stream, AIFuncConfig, ChatStreamInput

config = AIFuncConfig(
    base_url="https://your-api-endpoint/v1",
    model="your-model-name",
    api_key="your-api-key",
)

input = ChatStreamInput(
    message="进程和线程的区别是什么？用三句话回答。",
    # context="对话历史或背景（可选）",
)

async def main():
    async for token in await chat_stream(config, input):
        sys.stdout.write(token)
        sys.stdout.flush()
    sys.stdout.write("\n")

asyncio.run(main())
```

</details>

<details>
<summary>TypeScript</summary>

```typescript
import { chatStream, AIFuncConfig, ChatStreamInput } from './aifunc/chat-stream';

const config: AIFuncConfig = {
  baseURL: 'https://your-api-endpoint/v1',
  model: 'your-model-name',
  apiKey: 'your-api-key',
};

const input: ChatStreamInput = {
  message: '进程和线程的区别是什么？用三句话回答。',
  // context: '对话历史或背景（可选）',
};

async function main() {
  for await (const token of chatStream(config, input)) {
    process.stdout.write(token);
  }
  process.stdout.write('\n');
}

main().catch(console.error);
```

</details>

<details>
<summary>Go</summary>

```go
package main

import (
    "context"
    "fmt"
    "os"

    "your-module/aifunc/chat_stream"
)

func main() {
    config := &chat_stream.AIFuncConfig{
        BaseURL: "https://your-api-endpoint/v1",
        Model:   "your-model-name",
        APIKey:  "your-api-key",
    }

    input := chat_stream.ChatStreamInput{
        Message: "进程和线程的区别是什么？用三句话回答。",
        // Context: strPtr("对话历史或背景（可选）"),
    }

    tokens, errc := chat_stream.ChatStream(context.Background(), config, input)
    for token := range tokens {
        fmt.Print(token)
    }
    if err := <-errc; err != nil {
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    }
    fmt.Println()
}
```

</details>

<details>
<summary>Java</summary>

```java
import aifunc.AIFuncConfig;
import aifunc.chat_stream.ChatStream;
import aifunc.chat_stream.ChatStreamTypes.ChatStreamInput;

AIFuncConfig config = AIFuncConfig.builder()
        .baseUrl("https://your-api-endpoint/v1")
        .model("your-model-name")
        .apiKey("your-api-key")
        .build();

ChatStreamInput input = new ChatStreamInput(
        "进程和线程的区别是什么？用三句话回答。",
        null  // context 可选
);

try (var tokens = ChatStream.chatStream(config, input)) {
    while (tokens.hasNext()) {
        System.out.print(tokens.next());
    }
}
System.out.println();
```

</details>

<details>
<summary>C#</summary>

```csharp
using Aifunc;
using Aifunc.ChatStream;

var config = new AIFuncConfig
{
    BaseUrl = "https://your-api-endpoint/v1",
    Model   = "your-model-name",
    ApiKey  = "your-api-key",
};

var input = new ChatStreamTypes.ChatStreamInput(
    message: "进程和线程的区别是什么？用三句话回答。"
    // context: "对话历史或背景（可选）"
);

await foreach (var token in ChatStream.ChatStreamAsync(config, input))
{
    Console.Write(token);
}
Console.WriteLine();
```

</details>

可运行完整例程：

- chat-stream：[Python](../examples/python/chat-stream) / [TypeScript](../examples/typescript/chat-stream) / [Go](../examples/go/chat-stream) / [Java](../examples/java/chat-stream) / [C#](../examples/csharp/chat-stream)
- all-packages-stream：[Python](../examples/python/all-packages-stream) / [TypeScript](../examples/typescript/all-packages-stream) / [Go](../examples/go/all-packages-stream) / [Java](../examples/java/all-packages-stream) / [C#](../examples/csharp/all-packages-stream)


</details>

<details>
<summary><b>用其他服务端点？</b></summary>

在 config 中指定对应的 `baseURL`：

```typescript
const config: AIFuncConfig = {
  baseURL: 'https://your-provider-api/v1',
  model: 'your-model-name',
  apiKey: 'your-api-key',
};
```

```typescript
// 本地部署的推理服务通常不需要 apiKey
const config: AIFuncConfig = {
  baseURL: 'http://localhost:11434/v1',
  model: 'your-local-model',
  apiKey: '',
};
```
</details>

<details>
<summary><b>生成的代码要提交到 Git 吗？</b></summary>

建议提交。团队成员 clone 后可直接使用，无需重新生成。

`.aifunc/` 缓存目录会自动加入 `.gitignore`。
</details>

<details>
<summary><b>CLI 是 Go 写的，会和我的项目冲突吗？</b></summary>

不会。CLI 只在执行 `aifn` 命令时使用。

运行时由生成的 TypeScript/Python/Go/Java/C# 代码负责，零外部依赖（无需 npm / pip / `go get` / Maven / NuGet，仅生成源码）。

</details>

<details>
<summary><b>mock 模式的返回值从哪来？</b></summary>

每个包目录下有 `mock.json`，定义了预设 case。引擎的查找顺序：

1. 精确匹配输入 → 返回对应 output
2. 无 input 的 fallback case → 返回其 output
3. 都没匹配 → 根据 `api.json` 的 output schema 自动生成零值

</details>