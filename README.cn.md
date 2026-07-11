<div align="center">

<h1>AIFunc</h1>

<p><strong>AI 即函数，提示词即代码</strong></p>

<p>像调用函数一样使用 AI，像发布 npm 包一样分享 AI</p>
<p>强类型 · 可测试 · 跨语言 · 零依赖 · 模型无关 · Git 原生</p>

<p>
  <a href="./LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License"></a>
  <a href="https://github.com/aifunc/cli"><img src="https://img.shields.io/badge/CLI-Go-00ADD8?logo=go&logoColor=white" alt="Go CLI"></a>
  <img alt="Node.js ≥ 18" src="https://img.shields.io/badge/Node.js-≥18-339933?logo=nodedotjs&logoColor=white">
  <img alt="Python ≥ 3.10" src="https://img.shields.io/badge/Python-≥3.10-3776AB?logo=python&logoColor=white">
  <img alt="Go ≥ 1.23" src="https://img.shields.io/badge/Go-≥1.23-00ADD8?logo=go&logoColor=white">
  <img alt="TypeScript types" src="https://img.shields.io/badge/TypeScript-types-3178C6?logo=typescript&logoColor=white">
</p>

</div>

---

## 为什么 AIFunc

加一个 AI 功能到产品里，本应很简单。

现实是：引入新框架、理解新的流程抽象、学 Prompt 工程、写胶水代码、想办法测试、担心换模型怎么办。写了一周代码，核心逻辑只有 3 行。

80% 的真实需求其实很简单：**文本进去，结构化数据出来。无状态，无记忆。** 就是一个函数。

AIFunc 把这件事变成了开箱即用。不引入新概念，不入侵你的架构。你只需要 import。

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

// mock 模式：无需 API Key，离线即可运行
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
	fmt.Println("Summary   :", result.Summary)
	fmt.Println("Word count:", result.WordCount)
}
```

连接真实模型时，把 `mock: true` 替换为实际的 `baseURL`、`model` 和 `apiKey` 即可，支持任何兼容 OpenAI 协议的服务端点。

```typescript
// TypeScript / Python / Go 配置字段相同
const config: AIFuncConfig = {
  baseURL: 'https://your-api-endpoint/v1',
  model: 'your-model-name',
  apiKey: 'your-api-key',
};
```

> 完整可运行代码见 [examples/go/hello-aifunc](./examples/go/hello-aifunc)、[examples/typescript/hello-aifunc](./examples/typescript/hello-aifunc)、[examples/python/hello-aifunc](./examples/python/hello-aifunc)

---

## 组合多个 AI 函数

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

`if`、`switch`、`Promise.all` —— 你本来就掌握的控制流，不需要学任何新东西。

> 完整示例见 [examples/typescript/customer-support](./examples/typescript/customer-support)、[examples/python/customer-support](./examples/python/customer-support) 和 [examples/go/customer-support](./examples/go/customer-support)

> **想实现带记忆和滑动窗口的多轮对话？** 无需引入重型 Agent 框架，用原生数组管理上下文即可。
> 示例见 [examples/typescript/chat-with-context](./examples/typescript/chat-with-context)、[examples/python/chat-with-context](./examples/python/chat-with-context)、[examples/go/chat-with-context](./examples/go/chat-with-context)

---

## 特性

**强类型** — 输入输出都有完整类型定义，IDE 自动补全，拼错字段编译期报错。

**可测试** — 每个包自带 Mock 数据，`mock: true` 即可离线运行。CI 无需 API Key，零成本测试。

**跨语言** — 一份包定义（`api.json` + `package.json` + `prompts/`），编译到 TypeScript、Python、Go，行为一致。

**零依赖** — 运行时 Engine 是生成在项目中的纯源码，各语言均只使用原生库，不引入任何第三方依赖。

**模型无关** — 支持任何兼容 OpenAI 协议的端点。切换模型只改 config，零代码改动。

**Git 原生** — 编译产物提交 Git，团队成员 clone 后直接 import 使用，无需安装 CLI。版本管理、Code Review、权限控制全部复用 Git 工作流。

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