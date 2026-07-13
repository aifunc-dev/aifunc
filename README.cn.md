<div align="center">

<h1>AIFunc</h1>

<p><strong>AI 即函数，提示词即代码</strong></p>

<p>像调用函数一样使用 AI，像发布 npm 包一样分享 AI</p>
<p>强类型 · 可测试 · 跨语言 · 零依赖 · 模型无关 · Git 原生</p>

<p>
  <a href="./LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License"></a>
  <img alt="TypeScript / Node.js ≥ 18" src="https://img.shields.io/badge/TypeScript_/_Node.js-≥18-3178C6?logo=typescript&logoColor=white">
  <img alt="Python ≥ 3.10" src="https://img.shields.io/badge/Python-≥3.10-3776AB?logo=python&logoColor=white">
  <img alt="Go ≥ 1.23" src="https://img.shields.io/badge/Go-≥1.23-00ADD8?logo=go&logoColor=white">
  <img alt="Java ≥ 11" src="https://img.shields.io/badge/Java-≥11-ED8B00?logo=openjdk&logoColor=white">
  <img alt="C# / .NET ≥ 6" src="https://img.shields.io/badge/C%23_/_NET-≥6-512BD4?logo=dotnet&logoColor=white">
</p>

</div>

---

## 为什么 AIFunc

加一个 AI 功能到产品里，本应很简单。现实却是：引入新框架、学 Prompt 工程、写胶水代码、想办法测试。写了一周代码，核心逻辑只有 3 行。

**80% 的真实需求其实很简单：文本进去，结构化数据出来。无状态，无记忆。它本质上就是一个函数。**

AIFunc 的目标就是：**让你像使用普通函数一样，使用、管理和分享 AI 能力。**

### 它是怎么做到的？

1. **声明式的包定义**：一个 AIFunc 包，就是一个文件夹，里面只有 package.json、api.json 和几个 prompts/*.md 文件。没有代码，没有运行时。
2. **CLI 跨语言编译**：通过 `aifn` CLI 一行命令，这些声明文件会被“编译”成你当前项目对应的原生语言代码（TypeScript、Python、Go、Java、C#）。
3. **零运行时依赖**：生成的代码只依赖语言原生标准库。你不需要 `npm install` 或 `pip install` 任何重型框架，直接 `import` 就能用。

### 为什么这样做更好？

- **不引入新概念，不入侵架构**：你拿到的就是强类型的原生函数。用 `if/else` 串联逻辑，用原生数组管理上下文，告别复杂的 Agent 编排框架。
- **开箱即用的测试**：生成的代码自带 Mock 数据，`mock: true` 即可离线运行。CI 流水线无需真实 API Key。
- **Git 原生的分享机制**：写好一个包，推到 Git 仓库，别人一行命令安装。不需要 npm/PyPI，Fork、PR、版本控制直接复用你的 Git 工作流。

---

## Quick Start

```bash
# 安装 CLI
brew tap aifunc-dev/aifn && brew install aifn   # macOS/Linux
scoop bucket add aifn https://github.com/aifunc-dev/scoop-aifn && scoop install aifn  # Windows
```

```bash
# 安装一个 AI 函数
aifn install github:aifunc-dev/aifunc-packages/summarize
```

### TypeScript

```typescript
import { summarize, AIFuncConfig, SummarizeInput } from './aifunc/summarize';

const config: AIFuncConfig = { mock: true };

const text =
  'The James Webb Space Telescope captured its first full-color images in July 2022, ' +
  'revealing thousands of galaxies in a single image.';

async function main() {
  const result = await summarize(config, { text, maxLength: 30 } as SummarizeInput);
  console.log(`Summary   : ${result.summary}`);   // ← IDE 自动补全，类型安全
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

<details>
<summary>Go</summary>

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
	fmt.Println("Summary   :", result.Summary)
	fmt.Println("Word count:", result.WordCount)
}
```

</details>

<details>
<summary>Java</summary>

```java
import aifunc.summarize.Summarize;
import aifunc.summarize.SummarizeTypes.SummarizeInput;
import aifunc.summarize.SummarizeTypes.SummarizeOutput;
import aifunc._engine.java.v0_1_0.Types.AIFuncConfig;

AIFuncConfig config = AIFuncConfig.builder().mock(true).build();

String text = "The James Webb Space Telescope captured its first full-color images in July 2022, " +
              "revealing thousands of galaxies in a single image.";

SummarizeInput input = new SummarizeInput(text, 30);
SummarizeOutput result = Summarize.summarize(config, input);
System.out.println("Summary   : " + result.getSummary());
System.out.println("Word count: " + result.getWordCount());
```

</details>

<details>
<summary>C#</summary>

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

</details>

连接真实模型时，把 `mock: true` 替换为实际的 `baseURL`、`model` 和 `apiKey` 即可，支持任何兼容 OpenAI 协议的服务端点。

```typescript
// TypeScript / Python / Go / Java / C# 配置字段相同
const config: AIFuncConfig = {
  baseURL: 'https://your-api-endpoint/v1',
  model: 'your-model-name',
  apiKey: 'your-api-key',
};
```

> 完整可运行代码见 [examples/go/hello-aifunc](./examples/go/hello-aifunc)、[examples/typescript/hello-aifunc](./examples/typescript/hello-aifunc)、[examples/python/hello-aifunc](./examples/python/hello-aifunc)、[examples/java/hello-aifunc](./examples/java/hello-aifunc)、[examples/csharp/hello-aifunc](./examples/csharp/hello-aifunc)

---

## 适用场景与边界

AIFunc 的设计哲学是**“把 AI 降维成普通函数”**，这决定了它有极其明确的适用边界：

✅ **极其适合的场景**
- **集成到既有系统**：在不改变现有架构的前提下，无缝插入 AI 能力。
- **相对固定的业务流**：如数据清洗、信息抽取（Text to JSON）、意图识别、文本分类、内容摘要等。
- **多轮对话场景**：**完全支持。** 但需要依靠你的业务代码（如使用原生数组或数据库）维护上下文状态，每次将历史记录拼接传给 AI 函数。不引入复杂的 Agent 记忆框架。
- **多模型切换与适配**：底层兼容 OpenAI 协议，只需修改 `config` 即可随时切换不同厂商的大模型，零代码改动。

❌ **不适合的场景**
- **高度开放的自主 Agent**：如需要 AI 自行规划、循环调用工具链的开放性任务。
  AIFunc 提倡用代码编排控制流，而非让 AI 自由发挥。

> 💡 **关于流式输出**
> 当前版本专注于“结构化数据提取”（等待完整结果返回），**流式输出支持将在后续版本中添加**。

---

## 像写业务代码一样编排 AI

AI 函数可以像普通函数一样自由组合。用你熟悉的语言控制流串联业务逻辑：

```typescript
import { analyzeSentiment, AIFuncConfig } from './aifunc/analyze-sentiment';
import { recognizeIntent } from './aifunc/recognize-intent';
import { extractJson } from './aifunc/extract-json';

const config: AIFuncConfig = { /* ... */ };

async function handleTicket(message: string) {
  const sentiment = await analyzeSentiment(config, {
    text: message,
    labels: ['angry', 'frustrated', 'neutral', 'happy'],
  });

  if (sentiment.label === 'angry' && sentiment.confidence > 0.7) {
    return { action: 'escalate', priority: 'HIGH' };
  }

  const intent = await recognizeIntent(config, {
    text: message,
    intents: ['query_order', 'request_refund', 'technical_support', 'billing_issue'],
  });

  const info = await extractJson(config, {
    text: message,
    fields: [
      { name: 'order_id', description: 'Order number', type: 'string' },
      { name: 'issue', description: 'What the customer wants', type: 'string' },
    ],
  });

  return { action: intent.intent, ...info.extracted };
}
```

使用 `if`、`switch` 处理控制流，通过 `array` 或 db 管理上下文与数据 —— 依靠原生语言特性即可完成逻辑串联。不引入额外的编排抽象，能够无缝集成至现有系统。

完整示例见：[TypeScript](./examples/typescript/customer-support) / [Python](./examples/python/customer-support) / [Go](./examples/go/customer-support) / [Java](./examples/java/customer-support) / [C#](./examples/csharp/customer-support)

带记忆的多轮对话示例见：[TypeScript](./examples/typescript/chat-with-context) / [Python](./examples/python/chat-with-context) / [Go](./examples/go/chat-with-context) / [Java](./examples/java/chat-with-context) / [C#](./examples/csharp/chat-with-context)

---

## 特性

**强类型** — 输入输出都有完整类型定义，IDE 自动补全，拼错字段编译期报错。

**可测试** — 每个包自带 Mock 数据，`mock: true` 即可离线运行。CI 无需 API Key，零成本测试。

**跨语言** — 一份包定义（`api.json` + `package.json` + `prompts/`），编译到 TypeScript、Python、Go、Java、C#，行为一致。

**零依赖** — 运行时 Engine 是生成在项目中的纯源码，各语言均只使用原生库，不引入任何第三方依赖。

**模型无关** — 支持任何兼容 OpenAI 协议的端点。切换模型只改 config，零代码改动。

**Git 原生** — 编译产物提交 Git，团队成员 clone 后直接 import 使用，无需安装 CLI。版本管理、Code Review、权限控制全部复用 Git 工作流。

> 一句话：**AIFunc 把 AI 能力变成了你的代码库里的普通函数——强类型、可测试、可版本管理。**

---

## 分享与复用

写好一个 AI 函数，推到 Git，任何人一行命令安装：

```bash
# 别人安装你发布的包
aifn install github:your-name/your-packages/summarize

# 本地包也行
aifn install ../shared-packages/translate
```

不需要发布到 npm，不需要注册中心。**Git 仓库就是包仓库。** Fork、PR、Tag、权限控制——你已有的 Git 工作流直接复用。

想创建自己的包？只需要一个 `package.json`包定义，一个 `api.json` 定义接口和 一个 `prompts/` 写提示词：

```bash
aifn create my-analyzer
# 生成包骨架，填写 package.json，api.json 和 prompt
```

> 详见 [创建 AIFunc 包](./docs_cn/05-create-package.md) 和 [分享与发布](./docs_cn/07-sharing.md)

---

## 团队协作


| 角色     | 职责                        |
| ------ | ------------------------- |
| AI 工程师 | 创建包、编写 Prompt、调优效果        |
| 应用开发者  | `import` → 调用，不需要懂 Prompt |
| 平台负责人  | 管理包仓库、审核 PR、控制版本          |

---

## 可用的 AI 函数包

官方包仓库提供开箱即用的常用 AI 函数：


| 包名                  | 用途              |
| ------------------- | --------------- |
| `summarize`         | 文本摘要            |
| `analyze-sentiment` | 情感分析            |
| `recognize-intent`  | 意图识别            |
| `classify`          | 通用文本分类          |
| `extract-json`      | 从文本中提取结构化字段     |
| `extract-entities`  | 实体抽取（人名、地点、组织等） |
| `extract-keywords`  | 关键词提取           |
| `detect-language`   | 语言检测            |
| `translate`         | 翻译              |
| `rewrite`           | 文本改写            |
| `generate-reply`    | 生成回复            |
| `generate-title`    | 生成标题            |
| `generate-slug`     | 生成 URL slug     |
| `generate-email`    | 生成邮件            |
| `generate-post`     | 生成文章            |
| `answer-question`   | 问答              |
| `score-quality`     | 内容质量评分          |


安装任意包：

```bash
aifn install github:aifunc-dev/aifunc-packages/summarize  # 简写模式
aifn install https://github.com/aifunc-dev/aifunc-packages/tree/main/summarize # 完整URL模式
```

> 完整包列表与文档见 [aifunc-packages](https://github.com/aifunc-dev/aifunc-packages)

---

## 文档

- [Quick Start](./docs_cn/01-quick-start.md) — 5 分钟跑通完整流程
- [运行时 API 参考](./docs_cn/02-api.md)
- [CLI 命令参考](./docs_cn/03-cli.md)
- [工作原理](./docs_cn/04-how-it-works.md)
- [创建 AIFunc 包](./docs_cn/05-create-package.md)
- [包格式规范](./docs_cn/06-spec.md)
- [分享与发布](./docs_cn/07-sharing.md)
- [团队协作](./docs_cn/08-team-workflow.md)

---

Licensed under [Apache-2.0](./LICENSE).
