
# Quick Start

> **目标读者**：所有想使用 AIFunc 的开发者
> **本文内容**：5 分钟完成从安装 CLI 到调用 AI 函数拿到结果的完整流程
> **前置条件**：Node.js 18+、Python 3.10+ 或 Go 1.23+

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
- 识别项目类型（TypeScript / Python / Go）
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

	result, err := summarize.Summarize(context.Background(), config, summarize.SummarizeInput{
		Text:      text,
		MaxLength: 30,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Summary   : %s\n", result.Summary)
	fmt.Printf("Word count: %d\n", result.WordCount)
}
```

IDE 提供完整的类型提示和自动补全。

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

运行时由生成的 TypeScript/Python/Go 代码负责，零外部依赖。

</details>

<details>
<summary><b>mock 模式的返回值从哪来？</b></summary>

每个包目录下有 `mock.json`，定义了预设 case。引擎的查找顺序：

1. 精确匹配输入 → 返回对应 output
2. 无 input 的 fallback case → 返回其 output
3. 都没匹配 → 根据 `api.json` 的 output schema 自动生成零值

</details>