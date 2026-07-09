# AIFunc 包格式规范

> **目标读者**：包作者
> **本文内容**：定义 AIFunc 包的文件结构与字段规范。按此规范组织文件，即可被 CLI 识别、编译和生成代码
> **规范版本**：0.1.0

---

## 1. 包是什么

一个 AIFunc 包是一个目录，里面用声明式文件描述一个 AI 函数：

- 它接收什么输入、返回什么输出（`api.json`）
- 它是谁、做什么（`package.json`）
- 它怎么引导 AI 完成任务（`prompts/general.md`）

包内**不含任何可执行代码**。

---

## 2. 目录结构

```text
<package-name>/
├── package.json         [必需] 包身份与元数据
├── api.json             [必需] 函数输入输出契约
├── prompts/
│   └── general.md       [必需] 提示词模板
├── model-params.json    [可选] 建议 Engine 使用的模型调用参数
├── mock.json            [可选] 静态测试数据
├── README.md            [推荐] 说明文档
└── LICENSE              [推荐] 许可证
```

---

## 3. package.json

包的身份信息与运行配置。

| 字段 | 类型 | 必需 | 说明 |
|:---|:---|:---:|:---|
| `name` | string | ✅ | 包名（小写字母、数字、连字符，2-64 字符） |
| `type` | string | ✅ | 固定 `"standalone"` |
| `version` | string | ✅ | SemVer 格式，如 `"1.0.0"` |
| `description` | string | ✅ | 功能描述 |
| `engine` | string | ✅ | 要求的 AIFunc Engine 兼容版本，支持 SemVer 范围（如 `"^0.1.0"`、`">=0.1.0"`） |
| `engineOptions` | object | | 引擎运行时行为配置 |
| `engineOptions.injectOutputSchema` | boolean | | 默认 `true`，见下方说明 |
| `author` | object | | 作者信息 |
| `author.name` | string | | 作者名称 |
| `author.email` | string | | 联系邮箱 |
| `author.url` | string | | 个人主页或项目主页 |

`author` 对象中所有子字段均为可选，但如果提供 `author`，建议至少包含 `name`。

**`engineOptions.injectOutputSchema` 说明**：

- `true`（默认）：Engine 根据 `api.json` 的 output schema 自动向 Prompt 追加输出格式指令。包作者不需要在 Prompt 中描述返回格式。
- `false`：Engine 不注入格式指令。包作者需要在 `prompts/general.md` 中自行编写完整的输出格式引导。

**示例**：

```json
{
  "name": "@text-toolkit/summarizer",
  "type": "standalone",
  "version": "1.0.0",
  "description": "生成简洁摘要，返回摘要内容与大致词数。",
  "engine": "^0.1.0",
  "engineOptions": {
    "injectOutputSchema": true
  },
  "author": {
    "name": "Your Name",
    "email": "you@example.com",
    "url": "https://yoursite.com"
  }
}
```

包名支持 `@scope/name` 格式（如 `@text-toolkit/summarizer`），用于将相关包归组。Scope 仅是命名空间，不改变包的独立性。

---

## 4. api.json

定义函数的输入输出类型。CLI 据此生成强类型代码，Engine 据此做运行时校验。

| 字段 | 必需 | 说明 |
|:---|:---:|:---|
| `name` | ✅ | 函数名（生成代码的函数名来源）。小写字母、数字、下划线，1-64 字符 |
| `description` | ✅ | 功能描述（写入生成代码的注释） |
| `input` | ✅ | 输入 Schema，遵循 JSON Schema Draft 2020-12 |
| `output` | ✅ | 输出 Schema，遵循 JSON Schema Draft 2020-12 |

`input` 和 `output` 使用标准 JSON Schema 类型：`string`、`number`、`integer`、`boolean`、`object`、`array`。

**示例**：

```json
{
  "name": "summarize",
  "description": "生成简洁摘要。",
  "input": {
    "type": "object",
    "properties": {
      "text": { "type": "string", "description": "待摘要文本" },
      "maxLength": { "type": "integer", "description": "摘要最大词数", "default": 80 }
    },
    "required": ["text"]
  },
  "output": {
    "type": "object",
    "properties": {
      "summary": { "type": "string", "description": "生成的摘要" },
      "wordCount": { "type": "integer", "description": "摘要的大致词数" }
    },
    "required": ["summary", "wordCount"]
  }
}
```

---

## 5. prompts/general.md

给 AI 的指令模板。Engine 在运行时渲染模板变量后作为 System Prompt 发送给模型。

**模板变量**：

| 语法 | 说明 |
|:---|:---|
| `{{input.fieldName}}` | 插入 input 对象的指定字段值 |
| `{{input_json}}` | 将整个 input 序列化为 JSON 字符串插入 |

当 `api.json` 中 input 的 `type` 为 `string`（即整个输入就是一个字符串，而非 object）时，使用 `{{input}}` 直接插入原文。

当 input 为 object 时，`{{input}}` 等同于 `{{input_json}}`。

**示例**：

```markdown
# System

你是一个简洁、准确的摘要助手。

## 任务

为用户提供的文本生成摘要。

## 要求

- 摘要语言必须与输入文本一致，不要翻译。
- 摘要应简洁、准确、流畅。
- 保留最重要的信息，不要编造原文没有的内容。
- `summary` 不应超过 `maxLength` 指定的词数；未提供时默认 80。
- `wordCount` 应反映摘要的大致词数。

## 输入

待摘要文本：{{input.text}}

最大长度：{{input.maxLength}}
```

如果 `package.json` 中 `engineOptions.injectOutputSchema` 为 `true`（默认），你不需要在 Prompt 中描述输出格式。如果设为 `false`，你需要在 Prompt 中自行说明期望的返回 JSON 结构。

---

## 6. model-params.json（可选）

建议 Engine 使用的模型调用参数。如果不提供此文件，Engine 使用自身默认值。

**优先级**（从低到高）：

1. Engine 默认值
2. 本文件的配置
3. 调用方代码中显式传入的值

| 字段 | 类型 | 说明 |
|:---|:---|:---|
| `rules` | array | 参数规则列表 |
| `rules[].match` | object | 匹配条件 |
| `rules[].params.temperature` | number | 温度 |
| `rules[].params.topP` | number | 采样 top-p 值 |
| `rules[].params.maxTokens` | integer | 最大输出 Token 数 |

**match 条件**（三选一）：

| 字段 | 类型 | 说明 |
|:---|:---|:---|
| `match.model` | string | 精确匹配模型名 |
| `match.models` | string[] | 匹配模型名列表中的任意一个 |
| `match.pattern` | string | 正则表达式匹配模型名 |

**示例**：

```json
{
  "rules": [
    {
      "match": { "pattern": ".*" },
      "params": { "temperature": 0.1, "maxTokens": 500 }
    }
  ]
}
```

---

## 7. 命名规则汇总

| 约束项 | 规则 | 示例 |
|:---|:---|:---|
| 字符编码 | 所有文本文件 UTF-8，无 BOM | — |
| 包名 | 小写字母、数字、连字符，2-64 字符 | `summarizer` |
| Scope 包名 | `@scope/name` 格式 | `@text-toolkit/summarizer` |
| 函数名 | 小写字母、数字、下划线，1-64 字符 | `summarize` |
| 版本 | SemVer 2.0.0 | `1.0.0` |

---

## 8. 完整示例

一个最小完整包（仅必需文件）：

```text
summarizer/
├── package.json
├── api.json
└── prompts/
    └── general.md
```

**package.json**：

```json
{
  "name": "summarizer",
  "type": "standalone",
  "version": "1.0.0",
  "description": "生成简洁摘要，返回摘要内容与大致词数。",
  "engine": "^0.1.0",
  "author": { "name": "Your Name" }
}
```

**api.json**：

```json
{
  "name": "summarize",
  "description": "生成简洁摘要。",
  "input": {
    "type": "object",
    "properties": {
      "text": { "type": "string", "description": "待摘要文本" },
      "maxLength": { "type": "integer", "description": "摘要最大词数", "default": 80 }
    },
    "required": ["text"]
  },
  "output": {
    "type": "object",
    "properties": {
      "summary": { "type": "string" },
      "wordCount": { "type": "integer" }
    },
    "required": ["summary", "wordCount"]
  }
}
```

**prompts/general.md**：

```markdown
# System

你是一个简洁、准确的摘要助手。

## 任务

为用户提供的文本生成摘要。

## 要求

- 摘要语言必须与输入文本一致，不要翻译。
- 摘要应简洁、准确、流畅。
- 保留最重要的信息，不要编造原文没有的内容。

## 输入

待摘要文本：{{input.text}}

最大长度：{{input.maxLength}}
```

一个带可选文件与完整配置的包：

```text
summarizer/
├── package.json
├── api.json
├── prompts/
│   └── general.md
├── model-params.json
├── README.md
└── LICENSE
```

**package.json**（完整配置）：

```json
{
  "name": "@text-toolkit/summarizer",
  "type": "standalone",
  "version": "1.0.0",
  "description": "生成简洁摘要，返回摘要内容与大致词数。",
  "engine": "^0.1.0",
  "engineOptions": {
    "injectOutputSchema": true
  },
  "author": {
    "name": "Your Name",
    "email": "you@example.com",
    "url": "https://yoursite.com"
  }
}
```

---

## 接下来

- **想动手创建一个包？** → [创建 AIFunc 包](./05-create-package)
- **想了解 CLI 命令细节？** → [CLI 命令参考](./03-cli)
- **想发布给他人使用？** → [分享与发布](./07-sharing)
