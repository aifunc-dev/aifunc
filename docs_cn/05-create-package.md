# 创建 AIFunc 包

> **目标读者**：包作者
> **本文内容**：从零创建一个 AI 函数包，完成后可被安装到任何项目中像普通函数一样调用
> **前置条件**：已安装 aifn CLI ([安装方式](./01-quick-start#step-1安装-cli))

---

## Step 1：创建脚手架

```bash
aifn create my-summarizer
```

生成以下目录：

```text
my-summarizer/
├── package.json       ← 包元信息
├── api.json           ← API 接口定义
└── prompts/
    └── general.md     ← 提示词模板
```

---

## Step 2：编辑 package.json

填写包的基本信息：

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

完整字段说明见 [包格式规范 - package.json](./spec#3-packagejson)。

---

## Step 3：定义接口 api.json

想清楚两件事：函数接收什么、返回什么。

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

Schema 遵循 JSON Schema Draft 2020-12，支持 `string`、`number`、`integer`、`boolean`、`object`、`array` 类型。

---

## Step 4：编写提示词 prompts/general.md

告诉 AI 怎么完成这个任务：

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

默认情况下，你不需要在 Prompt 中写"请返回 JSON"之类的格式要求——Engine 会根据 `api.json` 的 output schema 自动注入输出格式指令。

如果你希望完全自己控制输出格式，也可以在 `package.json` 中关闭自动注入：

```json
{
  "engineOptions": {
    "injectOutputSchema": false
  }
}
```

关闭后，Engine 不再把 output schema 编译进提示词，包作者需要自行在 Prompt 中明确说明模型应该如何输出。

### 模板变量

| 语法 | 说明 |
|:---|:---|
| `{{input.fieldName}}` | 插入 input 对象的指定字段值 |
| `{{input_json}}` | 将整个 input 序列化为 JSON 字符串插入 |
| `{{input}}` | input 为 string 类型时插入原文；为 object 时等同于 `{{input_json}}` |

---

## Step 5：校验包格式

```bash
aifn validate ./my-summarizer
```

校验通过后输出包名、版本、描述、engine 版本及函数名。如有错误会明确指出缺失或不合规的字段。

---

## Step 6：本地安装并测试

在你的项目中安装这个本地包：

```bash
aifn install ./my-summarizer
```

然后在代码中使用：

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

本地包以 `file:` 路径链接，修改包源文件后重新 `aifn install` 即可更新产物。

---

## 可选：添加模型参数建议

如果你希望建议较低的 temperature 让输出更稳定，创建 `model-params.json`：

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

参数优先级（从低到高）：Engine 默认值 → 本文件配置 → 调用方代码中显式传入的值。

---

## 可选：添加 Mock 数据

创建 `mock.json` 提供测试用的静态响应，用于离线开发和 CI 测试：

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
        "summary": "用户使用三个月后，对手感和续航都很满意。",
        "wordCount": 20
      }
    }
  ]
}
```

---

## 写好 Prompt 的建议

| 建议 | 说明 |
|:---|:---|
| 明确角色 | 开头说清楚"你是什么专家" |
| 列出规则 | 用清单列出摘要质量、长度和边界要求，不要含糊 |
| 给出边界 | 说清楚边界情况怎么处理（如“输入过短时直接概括核心信息”） |
| 不写格式要求 | 输出格式由 Engine 自动注入，你只管任务逻辑 |
| 变量放最后 | `{{input.text}}` 放在 Prompt 末尾，模型更容易定位输入内容 |

---

## 完整的包目录参考

最小包（仅必需文件）：

```text
my-package/
├── package.json
├── api.json
└── prompts/
    └── general.md
```

带可选文件的完整包：

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

## 接下来

| 目标 | 文档 |
|:---|:---|
| 将包分享给他人或发布到公开仓库 | [分享与发布](./07-sharing) |
| 在团队内管理私有包 | [团队协作](./08-team-workflow) |
| 查看所有字段的完整定义 | [包格式规范](./06-spec) |
| 查看 CLI 命令参考 | [CLI 命令参考](./03-cli) |
