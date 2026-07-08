# 团队协作

> **目标读者**：应用开发团队、企业内部 AI 团队、DevOps 工程师
> **本文内容**：团队环境下如何安装使用云端包、企业内部包管理与分工协作、Mock 测试以及 CI/CD 集成
> **前置条件**：已阅读 [Quick Start](./01-quick-start)、[分享与发布](./07-sharing)

---

## 使用云端包的团队

面向：直接安装已发布包的应用开发团队。

### 安装

```bash
aifn init
aifn install github:aifunc-dev/aifunc-packages/summarize
```

### 使用

```typescript
import { summarize, AIFuncConfig } from './aifunc/summarize';

const config: AIFuncConfig = { ... };
const result = await summarize(config, { text: "...", maxLength: 20 });
```

### 提交到 Git 的内容

| 文件 | 作用 |
|:---|:---|
| `aifunc.json` | 项目配置：语言、输出目录、已安装的包列表 |
| `aifunc-lock.json` | 锁定每个包的来源和版本，保证团队一致性 |
| `src/aifunc/`（或 `aifunc/`） | 编译产物，包含完整运行时代码，成员无需安装 CLI 即可直接使用 |

---

## 企业内部包管理与分工

面向：规模较大、分工明确，希望自持代码而非依赖云端开源包的团队。

> **适用场景判断**：如果你的团队只有几个人、或者直接使用云端包已经够用，上一节"使用云端包的团队"更适合你。本节适合以下情况：
> - 多个项目组（可能使用不同语言）需要消费同一批 AI 函数
> - 有专职的 Prompt 工程师或 AI 团队负责维护包质量
> - 出于安全、合规或网络隔离的考虑，不愿依赖外部托管的包
> - 希望完全掌控包的版本发布节奏
>
> 具体实践请结合自身团队规模、技术栈和管理需求调整。

### 核心模式：统一仓库，多项目多语言消费

所有 AI 函数包集中在一个内部 Git 仓库中维护。不同项目、不同语言的开发者从同一仓库安装包，CLI 根据各自项目配置的目标语言编译出对应的强类型代码。

```text
┌───────────────────────────────────────────────────┐
│  内部统一包仓库（Git 远程）                         │
│  summarize/ · extract-json/ · classify/ · ...      │
│  （纯声明式文件，语言无关）                         │
└─────────────────────────┬─────────────────────────┘
                          │ git clone / git pull
                          ▼
┌───────────────────────────────────────────────────┐
│  本地克隆（aifunc-packages/）                      │
│  开发者机器上的完整副本                             │
└─────────────────────────┬─────────────────────────┘
                          │ aifn install ../aifunc-packages/xxx
       ┌──────────────────┼──────────────────┐
       │                  │                  │
       ▼                  ▼                  ▼
┌─────────────┐   ┌─────────────┐   ┌─────────────┐
│ 订单服务     │   │ 数据管道     │   │ 风控服务     │
│ TypeScript  │   │ Python      │   │ Go (未来)   │
│ src/aifunc/ │   │ aifunc/     │   │ ...         │
└─────────────┘   └─────────────┘   └─────────────┘
```

包仓库只存放语言无关的声明文件（`api.json`、`prompts/`、`package.json`），不含可执行代码。CLI 在 install/build 阶段完成到目标语言的编译。

### 开发者工作流

无论使用什么语言，每个应用开发者的流程一致：

| 步骤 | 操作 | 频率 |
|:---|:---|:---|
| 克隆包仓库 | `git clone` 到本地 | 仅一次 |
| 初始化项目 | `aifn init` | 仅一次（选择语言、配置输出目录） |
| 本地安装 | `aifn install ./本地路径` | 第一次引入包时 |
| 包更新后同步 | `git pull` + `aifn build` | 上游发布新版本时 |
| 日常使用 | `import` → 调用 | 每天（无需 CLI） |

**第一步：克隆内部包仓库到本地**

```bash
# 所有开发者共享同一个包仓库，克隆一次即可
git clone https://gitlab.yourcompany.com/ai-team/aifunc-packages.git
```

克隆后的本地目录：

```text
aifunc-packages/
├── summarize/
├── extract-json/
├── classify/
└── generate-title/
```

**首次安装（以 TypeScript 项目为例）：**

```bash
cd order-service
aifn init                                # 选择 typescript，输出到 src/aifunc
aifn install ../aifunc-packages/summarize
```

**同一个包在 Python 项目中安装：**

```bash
cd data-pipeline
aifn init                                # 选择 python，输出到 aifunc
aifn install ../aifunc-packages/summarize
```

同一个 `summarize` 包，前者生成 `.ts` 文件，后者生成 `.py` 文件。包仓库无需为不同语言做适配。

**长期使用：** 当上游 Prompt 工程师更新了包，应用开发者只需：

```bash
cd aifunc-packages
git pull                     # 拉取最新包内容

cd ../order-service
aifn build                   # 从本地缓存重新编译，无需重新下载
```

安装后的编译产物建议提交到 Git，团队成员 clone 后可直接 import 使用，无需安装 CLI。仅当需要更新或新增包时，才执行上述流程。

### 内部包仓库组织

推荐使用 Monorepo 管理内部所有 AI 函数包：

```text
gitlab.yourcompany.com/ai-team/aifunc-packages/
├── summarize/
├── extract-json/
├── classify/
├── generate-title/
└── README.md
```

### 角色划分

| 角色 | 职责 | 工作内容 |
|:---|:---|:---|
| AI 工程师 | 创建和维护包 | 编写 `api.json`、`prompts/`、`model-params.json`、调优效果 |
| 应用开发者 | 消费包 | `aifn install` → import → 调用 |
| 平台负责人 | 管理包仓库 | 审核 PR、管理版本、维护权限 |

### 分工协作流程

```text
    AI 工程师                    应用开发者
     │                                │
     ├── 创建/修改包                   │
     ├── 本地测试 (mock mode)          │
     ├── 提交 PR                       │
     │        ↓                        │
     │   [平台负责人 Review]            │
     │        ↓                        │
     │   合并 → 打 tag                  │
     │                                 ├── git pull + aifn build
     │                                 └── 集成到业务代码
```

### 版本与稳定性

使用 Git Tag 管理稳定版本，避免应用团队被未经验证的变更影响：

```bash
# Prompt 工程师发布新版本
git tag summarize/v1.2.0
git push origin summarize/v1.2.0

# 应用开发者切换到指定版本后安装
cd aifunc-packages
git checkout summarize/v1.2.0

cd ../order-service
aifn install ../aifunc-packages/summarize
```

### 权限控制

利用 Git 平台的权限机制：

| 目标 | 做法 |
|:---|:---|
| 包仓库只有 AI 团队可写 | 仓库写权限仅授予 AI 团队 |
| 应用团队只读安装 | 授予 read / clone 权限 |
| 发布需审核 | 通过 PR + 保护分支实现 |
| 按部门隔离 | 不同 Group 放不同仓库 |

### 多语言项目 / Monorepo 场景

一个仓库内同时为多个子服务生成不同语言的产物：

```bash
aifn build -l typescript -o services/order-api/src/aifunc
aifn build -l python -o services/data-pipeline/aifunc
```

### 内部包目录（可选）

维护一个 `README.md` 或内部文档列出所有可用包：

| 包名 | 用途 | 负责人 | 最新版本 |
|:---|:---|:---|:---|
| summarize | 文本摘要 | 张三 | v1.2.0 |
| extract-json | 信息抽取 | 李四 | v2.0.1 |
| classify | 文本分类 | 王五 | v1.0.3 |

---

## 测试

面向：包作者和应用开发者。

### mock.json：基于 VCR 数据

将真实模型响应录制为 mock.json，作为测试基准：

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

### 不提供 mock.json

包目录中不存在 `mock.json` 时，Engine 在 Mock 模式下根据 `api.json` 的 output schema 自动生成伪数据：

无需真实数据即可开始开发，Engine 保证返回值符合类型定义。

### 选择建议

| 场景 | 建议 |
|:---|:---|
| CI 需要断言具体值 | 提供 mock.json（VCR 数据） |
| 快速原型、联调 | 不提供，依赖自动伪数据 |
| 验收测试、效果对比 | 提供 mock.json（基准数据） |
| 包刚创建、尚未调通 | 不提供 |

### 代码中启用 Mock

```typescript
import { summarize, AIFuncConfig } from './aifunc/summarize';

const config: AIFuncConfig = { mock: true };
const result = await summarize(config, { text: "test input", maxLength: 20 });
```

```python
from aifunc.summarize import summarize, AIFuncConfig, SummarizeInput

config = AIFuncConfig(mock=True)
result = await summarize(config, SummarizeInput(text="test input", max_length=20))
```

---

## CI/CD 集成

面向：DevOps、流水线维护者。

`aifn build` 支持通过 `--lang` 和 `--output` 直接指定编译目标和输出路径，适合在流水线中无交互执行。

### 编译并输出

```bash
# 读取 aifunc.json 配置编译
aifn build

# 直接指定语言和输出目录（不依赖 aifunc.json 中的配置）
aifn build -l typescript -o src/aifunc
aifn build -l python -o aifunc
```

### 多项目构建

```bash
aifn build -l typescript -o services/order-api/src/aifunc
aifn build -l python -o services/data-pipeline/aifunc
```

---

## 接下来

- **查看运行时 API？** → [运行时 API 参考](./02-api)
- **查看 CLI 所有命令？** → [CLI 命令参考](./03-cli)
- **了解内部运行机制？** → [工作原理](./04-how-it-works)
- **查看包格式完整定义？** → [包格式规范](./06-spec)

