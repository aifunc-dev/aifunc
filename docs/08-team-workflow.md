# Team Workflow

> **Target audience**: Application development teams, enterprise AI teams, DevOps engineers
> **Content**: How to install and use cloud-hosted packages in a team environment, enterprise internal package management and role collaboration, mock testing, and CI/CD integration
> **Prerequisites**: Have read [Quick Start](./01-quick-start) and [Sharing & Publishing](./07-sharing)

---

## Teams Using Cloud Packages

For: Application development teams that directly install published packages.

### Installation

```bash
aifn init
aifn install github:aifunc-dev/aifunc-packages/summarize
```

### Usage

```typescript
import { summarize, AIFuncConfig } from './aifunc/summarize';

const config: AIFuncConfig = { ... };
const result = await summarize(config, { text: "...", maxLength: 20 });
```

### What to Commit to Git

| File | Purpose |
|:---|:---|
| `aifunc.json` | Project config: language, output directory, installed package list |
| `aifunc-lock.json` | Locks each package's source and version for team consistency |
| `src/aifunc/` (or `aifunc/`) | Compiled artifacts with complete runtime code — members can use directly without installing the CLI |

---

## Enterprise Internal Package Management and Roles

For: Larger teams with clear role separation who want to self-host code rather than depend on cloud-hosted open-source packages.

> **When this applies**: If your team is small or using cloud packages directly is sufficient, the previous section "Teams Using Cloud Packages" is more appropriate. This section is for situations where:
> - Multiple project teams (possibly using different languages) need to consume the same set of AI functions
> - Dedicated prompt engineers or AI teams are responsible for maintaining package quality
> - Security, compliance, or network isolation concerns prevent relying on externally hosted packages
> - You want full control over package release cadence
>
> Adapt the specific practices to your team size, tech stack, and management needs.

### Core Pattern: Unified Repository, Multi-project Multi-language Consumption

All AI function packages are maintained centrally in one internal Git repository. Developers across different projects and languages install packages from the same repository. The CLI compiles strongly-typed code for each project's configured target language.

```text
┌───────────────────────────────────────────────────┐
│  Internal unified package repo (Git remote)        │
│  summarize/ · extract-json/ · classify/ · ...      │
│  (Pure declarative files, language-agnostic)       │
└─────────────────────────┬─────────────────────────┘
                          │ git clone / git pull
                          ▼
┌───────────────────────────────────────────────────┐
│  Local clone (aifunc-packages/)                    │
│  Full copy on developer's machine                  │
└─────────────────────────┬─────────────────────────┘
                          │ aifn install ../aifunc-packages/xxx
       ┌──────────────────┼──────────────────┐
       │                  │                  │
       ▼                  ▼                  ▼
┌─────────────┐   ┌─────────────┐   ┌─────────────┐
│ Order svc    │   │ Data pipeline│   │ Risk svc    │
│ TypeScript  │   │ Python      │   │ Go (future) │
│ src/aifunc/ │   │ aifunc/     │   │ ...         │
└─────────────┘   └─────────────┘   └─────────────┘
```

The package repository only stores language-agnostic declarative files (`api.json`, `prompts/`, `package.json`) — no executable code. The CLI handles compilation to the target language during install/build.

### Developer Workflow

Regardless of language, every application developer follows the same flow:

| Step | Action | Frequency |
|:---|:---|:---|
| Clone package repo | `git clone` locally | Once only |
| Initialize project | `aifn init` | Once only (choose language, configure output dir) |
| Local install | `aifn install ./local-path` | First time introducing a package |
| Sync after package update | `git pull` + `aifn build` | When upstream publishes a new version |
| Daily usage | `import` → call | Every day (no CLI needed) |

**Step 1: Clone the internal package repository locally**

```bash
# All developers share the same package repo, clone once
git clone https://gitlab.yourcompany.com/ai-team/aifunc-packages.git
```

Local directory after cloning:

```text
aifunc-packages/
├── summarize/
├── extract-json/
├── classify/
└── generate-title/
```

**First install (TypeScript project example):**

```bash
cd order-service
aifn init                                # Choose typescript, output to src/aifunc
aifn install ../aifunc-packages/summarize
```

**Same package in a Python project:**

```bash
cd data-pipeline
aifn init                                # Choose python, output to aifunc
aifn install ../aifunc-packages/summarize
```

The same `summarize` package generates `.ts` files for the former and `.py` files for the latter. The same applies to Go (`.go`), Java (`.java`), and C# (`.cs`). The package repo needs no language-specific adaptations.

**Long-term usage:** When upstream prompt engineers update a package, application developers just need:

```bash
cd aifunc-packages
git pull                     # Pull latest package content

cd ../order-service
aifn build                   # Recompile from local cache, no re-download needed
```

Compiled artifacts are recommended to be committed to Git so team members can import directly after cloning without needing the CLI. The above flow is only needed when updating or adding packages.

### Internal Package Repository Organization

A monorepo is recommended for managing all internal AI function packages:

```text
gitlab.yourcompany.com/ai-team/aifunc-packages/
├── summarize/
├── extract-json/
├── classify/
├── generate-title/
└── README.md
```

### Role Separation

| Role | Responsibility | Work Content |
|:---|:---|:---|
| AI Engineer | Create and maintain packages | Write `api.json`, `prompts/`, `model-params.json`, tune performance |
| Application Developer | Consume packages | `aifn install` → import → call |
| Platform Owner | Manage package repo | Review PRs, manage versions, maintain permissions |

### Collaboration Flow

```text
    AI Engineer                   Application Developer
     │                                │
     ├── Create/modify package         │
     ├── Local test (mock mode)        │
     ├── Submit PR                     │
     │        ↓                        │
     │   [Platform Owner Review]       │
     │        ↓                        │
     │   Merge → tag                   │
     │                                 ├── git pull + aifn build
     │                                 └── Integrate into business code
```

### Versioning and Stability

Use Git Tags to manage stable versions, preventing application teams from being affected by unverified changes:

```bash
# Prompt engineer releases new version
git tag summarize/v1.2.0
git push origin summarize/v1.2.0

# Application developer checks out specific version before installing
cd aifunc-packages
git checkout summarize/v1.2.0

cd ../order-service
aifn install ../aifunc-packages/summarize
```

### Access Control

Leverage your Git platform's permission mechanisms:

| Goal | Approach |
|:---|:---|
| Only AI team can write to package repo | Grant write access only to the AI team |
| Application teams can only read/install | Grant read / clone permissions |
| Releases require review | Implement via PR + protected branches |
| Isolate by department | Use different Groups with separate repos |

### Multi-language Projects / Monorepo Scenarios

Generate artifacts for multiple sub-services in different languages within one repository:

```bash
aifn build -l typescript -o services/order-api/src/aifunc
aifn build -l python -o services/data-pipeline/aifunc
aifn build -l go -o services/risk-svc/aifunc
aifn build -l java -o services/billing/aifunc
aifn build -l csharp -o services/admin/aifunc
```

### Internal Package Directory (Optional)

Maintain a `README.md` or internal doc listing all available packages:

| Package | Purpose | Owner | Latest Version |
|:---|:---|:---|:---|
| summarize | Text summarization | Alice | v1.2.0 |
| extract-json | Information extraction | Bob | v2.0.1 |
| classify | Text classification | Charlie | v1.0.3 |

---

## Testing

For: Package authors and application developers.

### mock.json: VCR-based Data

Record real model responses as mock.json to serve as test baselines:

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
        "summary": "After three months of use, very satisfied with the feel and battery life.",
        "wordCount": 20
      }
    }
  ]
}
```

### Without mock.json

When `mock.json` doesn't exist in the package directory, the Engine auto-generates pseudo data based on the `api.json` output schema in mock mode:

You can start developing without real data — the Engine guarantees return values conform to the type definitions.

### Recommendations

| Scenario | Recommendation |
|:---|:---|
| CI needs to assert specific values | Provide mock.json (VCR data) |
| Rapid prototyping, integration testing | Don't provide, rely on auto-generated pseudo data |
| Acceptance testing, effect comparison | Provide mock.json (baseline data) |
| Package just created, not yet working | Don't provide |

### Enabling Mock in Code

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

```csharp
using Aifunc;
using Aifunc.Summarize;

var config = new AIFuncConfig { Mock = true };
var result = await Summarize.SummarizeAsync(config, new SummarizeTypes.SummarizeInput("test input", 20));
```

---

## CI/CD Integration

For: DevOps and pipeline maintainers.

`aifn build` supports `--lang` and `--output` to directly specify the compilation target and output path, suitable for non-interactive execution in pipelines.

### Compile and Output

```bash
# Read aifunc.json config to compile
aifn build

# Specify language and output directory directly (independent of aifunc.json config)
aifn build -l typescript -o src/aifunc
aifn build -l python -o aifunc
aifn build -l go -o aifunc
aifn build -l java -o aifunc
aifn build -l csharp -o aifunc
```

### Multi-project Build

```bash
aifn build -l typescript -o services/order-api/src/aifunc
aifn build -l python -o services/data-pipeline/aifunc
aifn build -l go -o services/risk-svc/aifunc
aifn build -l java -o services/billing/aifunc
aifn build -l csharp -o services/admin/aifunc
```

---

## Next Steps

- **View the runtime API?** → [Runtime API Reference](./02-api)
- **View all CLI commands?** → [CLI Command Reference](./03-cli)
- **Understand the internals?** → [How It Works](./04-how-it-works)
- **View the full package format definition?** → [Package Format Spec](./06-spec)
