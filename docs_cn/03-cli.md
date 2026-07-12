# aifn CLI 命令参考

> **目标读者**：所有 AIFunc 用户
> **本文内容**：汇总 aifn CLI v0.1.0 当前已实现的全部命令与用法
> **前置条件**：已安装 aifn CLI（[安装方式](./01-quick-start#step-1安装-cli)）

---

## 安装

```bash
# 安装 CLI
brew tap aifunc-dev/aifn && brew install aifn   # macOS/Linux
scoop bucket add aifn https://github.com/aifunc-dev/scoop-aifn && scoop install aifn  # Windows
```

验证安装：

```bash
aifn -v
# aifn v0.1.10
```

---

## 命令概览

```text
aifn
│
├── 项目管理
│   ├── init                          初始化项目（生成 aifunc.json）
│   └── list | ls                     列出已安装的包
│
├── 包消费
│   ├── install | i    <source...>    下载 + 校验 + 编译 + 生成（完整流程）
│   ├── uninstall | rm <name...>      卸载包（清除缓存、产物、配置记录）
│   └── build          [name...]      从已有缓存包生成目标语言代码
│
├── 包创作
│   ├── create         <name>         创建包脚手架
│   └── validate       <path>         校验包目录是否符合规范
│
└── 全局选项
    ├── --help | -h                   显示帮助信息
    └── --version | -v                显示版本号
```

---

## init

初始化当前目录为 AIFunc 项目。

```bash
aifn init
```

**行为：**

- 交互式选择项目语言（TypeScript / Python / Go / Java），自动检测项目环境并给出推荐
- 配置产物输出目录（TypeScript 默认 `src/aifunc`；Python、Go 和 Java 默认 `aifunc`）
- TypeScript 项目可配置路径别名（默认 `./aifunc`）
- 生成 `aifunc.json`
- 将缓存目录（默认 `.aifunc/`）写入 `.gitignore`

若 `aifunc.json` 已存在，则跳过初始化。

---

## install

安装 AI 函数包。

```bash
# 模式 1：读取 aifunc.json，安装所有已声明的包
aifn install

# 模式 2：指定来源安装，并自动写入 aifunc.json
aifn install github:owner/repo/path
aifn install https://github.com/owner/repo/tree/main/path
aifn install ./my-package
aifn install ../shared/classifier
aifn install /absolute/path/to/package

# 指定目标语言和输出目录（适用于 CI/CD，跳过交互式初始化）
aifn install --lang typescript --output dist/aifunc
aifn install <source...> -l python -o aifunc
```

别名：`aifn i`

### 选项

| 选项 | 别名 | 说明 |
|:---|:---|:---|
| `--lang <language>` | `-l` | 覆盖 `aifunc.json` 中的 `language`。可选值：`typescript`、`python`、`go`、`java` |
| `--output <dir>` | `-o` | 覆盖 `aifunc.json` 中的 `outputDir` |

### 支持的包来源格式

**简写模式**（推荐）：

| 格式 | 示例 |
|:---|:---|
| GitHub 简写 | `github:owner/repo/path` |
| 本地路径 | `./my-package` |

解析规则：前缀后的前两段为 `owner/repo`，后续路径为包所在子目录。

**完整 URL 模式**（直接从浏览器地址栏复制）：

| 格式 | 示例 |
|:---|:---|
| GitHub Tree URL | `https://github.com/owner/repo/tree/ref/path` |

> [!NOTE]
> 暂不支持 `name@version` 格式。

### 完整流程

1. **解析配置**：下载包源码到 `.aifunc/packages/`（本地包以 `file:` 路径链接，不复制）
2. **拉取 Engine SDK**：解析各包声明的 engine 版本，拉取 Engine SDK 到 `.aifunc/_engine/`
3. **编译生成**：生成目标语言代码到 `outputDir`
4. **更新锁文件**：更新 `aifunc-lock.json`

### 行为细节

- 无参数时读取 `aifunc.json` 的 `packages` 字段，逐个下载并安装
- 有参数时下载指定包，自动将来源写入 `aifunc.json`
- 若未找到 `aifunc.json` 且提供了 `--lang`，自动创建配置文件（跳过交互）
- 若未找到 `aifunc.json` 且未提供 `--lang`，进入交互式初始化
- `--lang` 和 `--output` 不会修改 `aifunc.json` 中已有的值，仅影响本次执行

---

## build

从已缓存的包重新编译产物，无需重新下载。

```bash
# 编译 aifunc.json 中声明的所有包
aifn build

# 仅编译指定包
aifn build summarize

# 指定目标语言和输出目录（适用于 CI/CD）
aifn build --lang python --output dist/aifunc
aifn build <package-name...> -l typescript -o src/aifunc
```

### 选项

| 选项 | 别名 | 说明 |
|:---|:---|:---|
| `--lang <language>` | `-l` | 覆盖 `aifunc.json` 中的 `language`。可选值：`typescript`、`python`、`go`、`java` |
| `--output <dir>` | `-o` | 覆盖 `aifunc.json` 中的 `outputDir` |

> [!NOTE]
> 包须已通过 `install` 下载并记录在 `aifunc-lock.json` 中，否则会提示先运行 `aifn install`。

**典型用途：** CI/CD 流水线中为不同语言目标分别生成产物：

```bash
aifn build -l typescript -o dist/ts
aifn build -l python -o dist/py
aifn build -l go -o dist/go
aifn build -l java -o dist/java
```

---

## uninstall

从项目中移除指定包。

```bash
aifn uninstall summarize
aifn rm summarize                         # 别名
```

别名：`aifn rm`

**行为：**

- 删除 `.aifunc/packages/` 中的缓存源码
- 删除 `outputDir` 中的生成产物
- 从 `aifunc.json` 和 `aifunc-lock.json` 中移除记录
- 清理不再被任何包引用的孤立 engine
- 清理空目录

支持按包名或来源 URL 匹配（如 `github:owner/repo/path`）。

---

## list

列出当前项目中声明的包及其安装状态。

```bash
aifn list
aifn ls    # 别名
```

别名：`aifn ls`

**输出信息：**

- 项目语言与输出目录
- 每个包的名称、版本、engine 版本、来源
- 未安装（仅在 `aifunc.json` 中声明但未 lock）的包会标注 `(未安装)`

---

## create

在当前目录下创建符合 AIFunc 包规范的脚手架目录。

```bash
aifn create my-snewummarizer
```

**生成的目录结构：**

```text
my-summarizer/
├── package.json       ← 包元信息
├── api.json           ← API 接口定义
└── prompts/
    └── general.md     ← 提示词模板
```

---

## validate

校验指定目录是否为合法的 AIFunc 包。

```bash
aifn validate ./my-summarizer
```

**检查内容：**

- 必需文件是否存在（`package.json`、`api.json`、`prompts/` 等）
- 字段完整性与格式正确性

校验通过后输出包名、版本、描述、engine 版本及函数名。

---

## 配置文件

CLI 读写以下项目级文件：

| 文件 | 说明 |
|:---|:---|
| `aifunc.json` | 项目配置：语言、输出目录、路径别名、包列表 |
| `aifunc-lock.json` | 版本锁定：已解析的包版本、engine 版本、完整性校验 |
| `.aifunc/` | CLI 缓存目录（默认，可通过 `inputDir` 自定义） |

`aifunc.json` 示例：

```json
{
  "configVersion": 1,
  "language": "typescript",
  "outputDir": "src/aifunc",
  "alias": "@aifunc",
  "packages": {
    "summarize": "github:aifunc-dev/aifunc-packages/summarize",
    "my-summarizer": "file:./packages/my-summarizer"
  }
}
```

---

## 典型工作流

### 消费者：在已有项目中安装并使用 AI 函数

```bash
aifn install github:aifunc-dev/aifunc-packages/summarize
```

### 消费者：从已有配置恢复（团队协作 / CI）

```bash
aifn install
```

### 创作者：开发本地包

```bash
aifn create my-summarizer
# 编辑 my-summarizer/ 下的文件...
aifn validate ./my-summarizer
aifn install ./my-summarizer
```

### CI/CD：为多语言目标生成代码

```bash
# 无需交互式初始化，直接指定语言和输出目录
aifn install -l typescript -o dist/ts
aifn build -l python -o dist/py
```

### 日常维护：重新编译 / 移除包

```bash
aifn build                    # 重新编译，不重新下载
aifn uninstall my-summarizer     # 移除不再需要的包
```

---

## 接下来

- **第一次使用？** → [Quick Start](./01-quick-start)
- **查看运行时 API？** → [运行时 API 参考](./02-api)
- **想创建自己的包？** → [创建 AIFunc 包](./05-create-package)
- **了解内部运行机制？** → [工作原理](./04-how-it-works)
