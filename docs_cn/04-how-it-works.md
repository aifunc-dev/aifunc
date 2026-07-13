# AIFunc 工作原理

> **目标读者**：想了解内部机制的开发者
> **本文内容**：说明 AIFunc 包从安装到运行的完整流程，包括 CLI 编译、项目结构与运行时调用链

> [!NOTE]
> 包作者只需要阅读 [包格式规范](./06-spec) 即可创建包，无需了解本文内容。

---

## 1. 组件说明

| 组件 | 语言 | 作用 | 你是否直接接触 |
|:---|:---|:---|:---|
| **aifn CLI** | Go | 包管理、依赖解析、代码生成 | ✅ 你在终端执行命令 |
| **Engine SDK** | 与目标语言一致（如 TypeScript、Python、Go、Java 或 C#） | 运行时加载编译产物、调用模型、校验输出 | ❌ 作为源码生成在你的项目中 |
| **生成的函数** | 目标语言 | 你 import 的强类型入口 | ✅ 你在代码中调用 |

CLI 是一个独立的 Go 二进制文件，与你的项目语言无关。Engine SDK 则是按需拉取的纯源码文件，仅使用语言标准库实现、零外部依赖，与生成的函数代码一起存放在你的项目目录中，无需通过 npm、pip、`go get`、Maven 或 NuGet 额外安装依赖。

---

## 2. 安装与编译流程

执行 `aifn install` 时，CLI 自动完成以下核心流程：

```text
解析配置 ──────► 下载源码与 Engine SDK ──────► 编译与语言包裹
```

1. **解析与下载**：CLI 读取项目中的 `aifunc.json`，获取需要安装的包列表。将各个包的原始文件下载到本地缓存目录（`.aifunc/packages/`）。
2. **拉取 Engine SDK**：CLI 分析各包声明的 `engineVersion`（引擎版本），自动下载对应版本的目标语言原生 Engine SDK 源文件到缓存目录（`.aifunc/_engine/`）。
3. **编译与包裹（核心机制）**：为避免传统应用打包器（如 Webpack/Vite）或语言运行时在加载纯数据文件时可能出现的异常，CLI 不再单独输出纯配置文件。而是**将 API 定义、Prompt 模板、模型参数等内容，使用你的目标语言进行包裹**。
   * TypeScript 项目会生成包含配置对象的 `.ts` 文件。
   * Python 项目会生成包含配置字典的 `.py` 文件。
   * Go 项目会生成包含配置 map 和结构体字面量的 `.go` 文件。
   * Java 项目会生成包含静态 `Map<String,Object>` 初始化器的 `.java` 文件。
   * C# 项目会生成包含对象/字典初始化器的 `.cs` 文件。
   * 同时生成对应的强类型接口文件和 Mock 数据包裹文件。
4. **链接输出**：最后，CLI 将生成的强类型函数、语言包裹产物，以及所需的 Engine SDK 源文件，一并输出到你指定的项目目录（如 `src/aifunc/`）中。

---

## 3. 安装后的项目结构

项目分为缓存区（不提交 Git）和编译产物区（提交 Git）。以下以 TypeScript、Python、Go、Java 和 C# 为例：

### TypeScript 项目示例

```text
your-project/
├── aifunc.json                          ← 包管理配置（你编辑的）
├── aifunc-lock.json                     ← 版本锁定文件
│
├── .aifunc/                             ← CLI 下载缓存（加入 .gitignore）
│   ├── packages/summarize/              ← 包原始源文件
│   └── _engine/v0.1.0/                  ← 下载的 Engine 原始文件
│
└── src/
    └── aifunc/                          ← 编译产物输出目录（提交 Git）
        ├── summarize/                   ← AI 函数包
        │   ├── index.ts                 ← 函数入口
        │   ├── summarize.types.ts       ← I/O 类型定义
        │   ├── summarize.aifunc.ts      ← 编译产物（提示词与 API 规范的 TS 包裹）
        │   └── summarize.mock.ts        ← 离线测试 Mock 数据（TS 包裹）
        │
        └── _engine/                     ← 运行时 SDK（本地源码级依赖）
            └── typescript/
                 └── v0.1.0/
                      ├── index.ts
                      ├── runtime.ts
                      └── ...
```

### Python 项目示例

```text
your-project/
├── aifunc.json
├── aifunc-lock.json
│
├── .aifunc/                             ← 缓存目录（加入 .gitignore）
│   ├── packages/summarize/
│   └── _engine/v0_1_0/
│
└── aifunc/                              ← 编译产物输出目录（提交 Git）
    ├── __init__.py                      ← 空文件，使目录成为 Python 包
    ├── py.typed                         ← 空文件，PEP 561 类型标记
    ├── summarize/
    │   ├── __init__.py                  ← 函数入口
    │   ├── summarize_types.py           ← I/O 类型定义
    │   ├── summarize_aifunc.py          ← 编译产物（提示词与 API 规范的 Python 包裹）
    │   └── summarize_mock.py            ← 离线测试 Mock 数据（Python 包裹）
    │
    └── _engine/
        └── python/
            └── v0_1_0/
                  ├── __init__.py
                  ├── runtime.py
                  └── ...
```

### Go 项目示例

```text
your-project/
├── aifunc.json
├── aifunc-lock.json
├── go.mod
│
├── .aifunc/                             ← 缓存目录（加入 .gitignore）
│   ├── packages/summarize/
│   └── _engine/go/v0.1.0/
│
└── aifunc/                              ← 编译产物输出目录（提交 Git）
    ├── summarize/                       ← AI 函数包（package summarize）
    │   ├── summarize.go                 ← 函数入口
    │   ├── summarize_types.go           ← I/O 结构体定义
    │   └── summarize_aifunc.go          ← 编译产物（提示词与 API 规范的 Go 包裹）
    │
    └── _engine/
        └── go/
            └── v0.1.0/
                  ├── aifunc.go          ← 公开 API
                  ├── runtime.go
                  ├── types.go
                  └── ...
```

### Java 项目示例

```text
your-project/
├── aifunc.json
├── aifunc-lock.json
│
├── .aifunc/                             ← 缓存目录（加入 .gitignore）
│   ├── packages/summarize/
│   └── _engine/java/v0_1_0/
│
└── aifunc/                              ← 编译产物输出目录（提交 Git）
    ├── summarize/                       ← AI 函数包（package aifunc.summarize）
    │   ├── Summarize.java               ← 函数入口（静态方法，返回 CompletableFuture）
    │   ├── SummarizeTypes.java          ← SummarizeInput / SummarizeOutput POJO
    │   ├── SummarizeAifunc.java         ← 编译产物（Map<String,Object> 字面量）
    │   └── SummarizeAifuncMock.java     ← Mock 数据（Map 字面量，无 mock.json 时为 null）
    │
    ├── AIFuncConfig.java                ← 公开配置（package aifunc）
    ├── AIFuncException.java             ← 公开异常（package aifunc）
    └── _engine/
        └── java/
            └── v0_1_0/
                  ├── Json.java          ← 零依赖 JSON 解析/序列化
                  ├── Types.java
                  ├── Runtime.java       ← executeAsync() 入口
                  ├── Validator.java
                  ├── Prompt.java
                  ├── Request.java
                  ├── Mock.java
                  └── Artifact.java
```

无需构建工具 — 生成文件是纯 Java 源码，零外部依赖。将 `aifunc/` 与业务代码放在同级（或加为 source root），用 `javac` 编译即可。

### C# 项目示例

```text
your-project/
├── aifunc.json
├── aifunc-lock.json
│
├── .aifunc/                             ← 缓存目录（加入 .gitignore）
│   ├── packages/summarize/
│   └── _engine/csharp/v0_1_0/
│
└── aifunc/                              ← 编译产物输出目录（提交 Git）
    ├── summarize/                       ← AI 函数包（namespace Aifunc.Summarize）
    │   ├── Summarize.cs                 ← 函数入口（SummarizeAsync）
    │   ├── SummarizeTypes.cs            ← SummarizeInput / SummarizeOutput
    │   ├── SummarizeAifunc.cs           ← 编译产物（字典字面量）
    │   └── SummarizeAifuncMock.cs       ← Mock 数据（字典字面量，无 mock.json 时为 null）
    │
    ├── AIFuncConfig.cs                  ← 公开配置（namespace Aifunc）
    ├── AIFuncException.cs               ← 公开异常（namespace Aifunc）
    └── _engine/
        └── csharp/
            └── v0_1_0/
                  ├── ...                ← Namespace Aifunc.Engine.Csharp.V0_1_0
                  └── ...
```

无需 NuGet 包 — 生成文件是纯 C# 源码，零外部依赖。将 `aifunc/` 与业务代码放在同级（或加为 source root），用 `dotnet run` 运行即可。

> [!NOTE]
> 包内的各个文件会根据目标语言的规范采用不同的命名风格（如 TS 使用 `.` 分隔，Python 和 Go 使用 `_` 分隔，Java 和 C# 使用 PascalCase）。所有的代码依赖都在生成的目录内部闭环完成。

---

## 4. 各文件职责说明

以生成的 `summarize` 包为例，内部文件分工明确：

| 文件（逻辑名） | 职责 |
|:---|:---|
| `entry` (如 `index.ts`、`summarize.go`、`Summarize.java`、`Summarize.cs`) | **函数入口**：创建并导出 AI 函数实例，你的业务代码直接引用它。 |
| `types` (如 `.types.ts`、`_types.go`、`SummarizeTypes.java`、`SummarizeTypes.cs`) | **接口文件**：函数的输入、输出结构类型定义，提供强类型支持。 |
| `aifunc` (如 `.aifunc.ts`、`_aifunc.go`、`SummarizeAifunc.java`、`SummarizeAifunc.cs`) | **核心产物**：Prompt 模板、API 规范、模型配置的合并结果，使用目标语言包裹以保证运行时安全加载。 |
| `mock` (如 `.mock.ts`、`_mock.py`、`SummarizeAifuncMock.java`、`SummarizeAifuncMock.cs`) | **Mock 数据**：输入到输出的映射数据，同样使用目标语言包裹，用于离线测试模式。（Go 直接内嵌在入口文件中。） |

---

## 5. 运行时调用链

```text
你的应用代码
  │  await summarize(config, { text: "...", maxLength: 20 })   ← TypeScript/Python
  │  summarize.Summarize(ctx, config, SummarizeInput{...})     ← Go
  │  Summarize.summarize(config, new SummarizeInput("...", 20)) ← Java（CompletableFuture）
  │  await Summarize.SummarizeAsync(config, new SummarizeTypes.SummarizeInput("...", 20))  ← C#
  ▼
生成的函数入口（aifunc/summarize/index.ts | summarize.go | Summarize.java | Summarize.cs）
  │  强类型入口，无业务逻辑
  │  引用同目录下的类型文件、产物文件，以及 _engine
  ▼
Engine SDK（aifunc/_engine/.../runtime.ts|runtime.py|runtime.go|Runtime.java|...）
  │  ① 解析由语言原生包裹的 .aifunc 配置对象
  │  ② 校验输入数据结构
  │  ③ 渲染 Prompt 模板（替换变量）
  │  ④ 调用 AI 模型 API（支持任何兼容 OpenAI 协议的端点）
  │  ⑤ 解析模型返回结果
  │  ⑥ 校验输出数据结构
  │  ⑦ 校验通过 → 返回强类型结果
  ▼
AI 模型 API
```

---

## 6. 配置文件

`aifunc.json` 是项目级配置，控制代码生成的语言与输出路径：

```json
{
  "version": "0.1",
  "language": "typescript",
  "outputDir": "src/aifunc",
  "alias": "@aifunc",
  "packages": {
    "summarize": "github:aifunc-dev/aifunc-packages/summarize"
  }
}
```

| 字段 | 说明 |
|:---|:---|
| `language` | 目标语言，决定了生成的语言包裹格式及引擎 SDK 语言（当前支持 `typescript`, `python`, `go`, `java`, `csharp`） |
| `outputDir` | 编译产物（包含生成的函数与 Engine SDK）的输出路径 |
| `alias` | （仅 TS）用于 tsconfig paths 别名的设定 |
| `packages` | 包名与对应安装源（支持 github 路径、本地路径）的映射关系 |

---

## 接下来

- **想创建自己的包？** → [创建 AIFunc 包](./05-create-package)
- **查看运行时 API？** → [运行时 API 参考](./02-api)
- **查看包格式的完整定义？** → [包格式规范](./06-spec)
- **查看 CLI 所有命令？** → [CLI 命令参考](./03-cli)
