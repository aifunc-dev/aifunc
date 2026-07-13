# aifn CLI Command Reference

> **Target audience**: All AIFunc users
> **Content**: Summary of all commands and usage for aifn CLI v0.1.0
> **Prerequisites**: aifn CLI installed ([installation guide](./01-quick-start#step-1-install-cli))

---

## Installation

```bash
# Install CLI
brew tap aifunc-dev/aifn && brew install aifn   # macOS/Linux
scoop bucket add aifn https://github.com/aifunc-dev/scoop-aifn && scoop install aifn  # Windows
```

Verify installation:

```bash
aifn -v
# aifn v0.1.10
```

---

## Command Overview

```text
aifn
│
├── Project Management
│   ├── init                          Initialize project (generate aifunc.json)
│   └── list | ls                     List installed packages
│
├── Package Consumption
│   ├── install | i    <source...>    Download + validate + compile + generate (full flow)
│   ├── uninstall | rm <name...>      Uninstall package (clear cache, artifacts, config)
│   └── build          [name...]      Generate target language code from cached packages
│
├── Package Authoring
│   ├── create         <name>         Create package scaffold
│   └── validate       <path>         Validate package directory against spec
│
└── Global Options
    ├── --help | -h                   Show help
    └── --version | -v                Show version
```

---

## init

Initialize the current directory as an AIFunc project.

```bash
aifn init
```

**Behavior:**

- Interactive language selection (TypeScript / Python / Go / Java / C#), auto-detects project environment and provides recommendations
- Configure artifact output directory (TypeScript defaults to `src/aifunc`; Python, Go, Java, and C# default to `aifunc`)
- TypeScript projects can configure a path alias (defaults to `./aifunc`)
- Generates `aifunc.json`
- Adds the cache directory (default `.aifunc/`) to `.gitignore`

If `aifunc.json` already exists, initialization is skipped.

---

## install

Install AI function packages.

```bash
# Mode 1: Read aifunc.json and install all declared packages
aifn install

# Mode 2: Install from specified source, automatically writes to aifunc.json
aifn install github:owner/repo/path
aifn install https://github.com/owner/repo/tree/main/path
aifn install ./my-package
aifn install ../shared/classifier
aifn install /absolute/path/to/package

# Specify target language and output directory (for CI/CD, skips interactive init)
aifn install --lang typescript --output dist/aifunc
aifn install <source...> -l python -o aifunc
```

Alias: `aifn i`

### Options

| Option | Alias | Description |
|:---|:---|:---|
| `--lang <language>` | `-l` | Override `language` in `aifunc.json`. Values: `typescript`, `python`, `go`, `java`, `csharp` |
| `--output <dir>` | `-o` | Override `outputDir` in `aifunc.json` |

### Supported Package Source Formats

**Shorthand format** (recommended):

| Format | Example |
|:---|:---|
| GitHub shorthand | `github:owner/repo/path` |
| Local path | `./my-package` |

Resolution rule: The first two segments after the prefix are `owner/repo`, the remaining path is the subdirectory containing the package.

**Full URL format** (copy directly from browser address bar):

| Format | Example |
|:---|:---|
| GitHub Tree URL | `https://github.com/owner/repo/tree/ref/path` |

> [!NOTE]
> The `name@version` format is not yet supported.

### Full Flow

1. **Parse config**: Download package source to `.aifunc/packages/` (local packages linked via `file:` path, not copied)
2. **Fetch Engine SDK**: Resolve each package's declared engine version, fetch Engine SDK to `.aifunc/_engine/`
3. **Compile and generate**: Generate target language code to `outputDir`
4. **Update lock file**: Update `aifunc-lock.json`

### Behavior Details

- Without arguments, reads the `packages` field from `aifunc.json` and installs each one
- With arguments, downloads specified packages and automatically writes the source to `aifunc.json`
- If `aifunc.json` is not found and `--lang` is provided, auto-creates the config file (skips interaction)
- If `aifunc.json` is not found and `--lang` is not provided, enters interactive initialization
- `--lang` and `--output` do not modify existing values in `aifunc.json`, they only affect the current execution

---

## build

Recompile artifacts from cached packages without re-downloading.

```bash
# Compile all packages declared in aifunc.json
aifn build

# Compile specific packages only
aifn build summarize

# Specify target language and output directory (for CI/CD)
aifn build --lang python --output dist/aifunc
aifn build <package-name...> -l typescript -o src/aifunc
```

### Options

| Option | Alias | Description |
|:---|:---|:---|
| `--lang <language>` | `-l` | Override `language` in `aifunc.json`. Values: `typescript`, `python`, `go`, `java`, `csharp` |
| `--output <dir>` | `-o` | Override `outputDir` in `aifunc.json` |

> [!NOTE]
> Packages must have been downloaded via `install` and recorded in `aifunc-lock.json`, otherwise you'll be prompted to run `aifn install` first.

**Typical use case:** Generate artifacts for different language targets in CI/CD pipelines:

```bash
aifn build -l typescript -o dist/ts
aifn build -l python -o dist/py
aifn build -l go -o dist/go
aifn build -l java -o dist/java
aifn build -l csharp -o dist/csharp
```

---

## uninstall

Remove a specified package from the project.

```bash
aifn uninstall summarize
aifn rm summarize                         # alias
```

Alias: `aifn rm`

**Behavior:**

- Deletes cached source in `.aifunc/packages/`
- Deletes generated artifacts in `outputDir`
- Removes records from `aifunc.json` and `aifunc-lock.json`
- Cleans up orphaned engines no longer referenced by any package
- Cleans up empty directories

Supports matching by package name or source URL (e.g., `github:owner/repo/path`).

---

## list

List packages declared in the current project and their installation status.

```bash
aifn list
aifn ls    # alias
```

Alias: `aifn ls`

**Output information:**

- Project language and output directory
- Each package's name, version, engine version, and source
- Packages that are only declared in `aifunc.json` but not locked will be marked as `(not installed)`

---

## create

Create a scaffold directory conforming to the AIFunc package spec in the current directory.

```bash
aifn create my-summarizer
```

**Generated directory structure:**

```text
my-summarizer/
├── package.json       ← Package metadata
├── api.json           ← API interface definition
└── prompts/
    └── general.md     ← Prompt template
```

---

## validate

Validate whether a specified directory is a valid AIFunc package.

```bash
aifn validate ./my-summarizer
```

**Checks:**

- Whether required files exist (`package.json`, `api.json`, `prompts/`, etc.)
- Field completeness and format correctness

On success, outputs the package name, version, description, engine version, and function name.

---

## Configuration Files

The CLI reads and writes the following project-level files:

| File | Description |
|:---|:---|
| `aifunc.json` | Project config: language, output directory, path alias, package list |
| `aifunc-lock.json` | Version lock: resolved package versions, engine versions, integrity checks |
| `.aifunc/` | CLI cache directory (default, customizable via `inputDir`) |

`aifunc.json` example:

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

## Typical Workflows

### Consumer: Install and use AI functions in an existing project

```bash
aifn install github:aifunc-dev/aifunc-packages/summarize
```

### Consumer: Restore from existing config (team collaboration / CI)

```bash
aifn install
```

### Author: Develop a local package

```bash
aifn create my-summarizer
# Edit files under my-summarizer/...
aifn validate ./my-summarizer
aifn install ./my-summarizer
```

### CI/CD: Generate code for multiple language targets

```bash
# No interactive init needed, specify language and output directly
aifn install -l typescript -o dist/ts
aifn build -l python -o dist/py
aifn build -l go -o dist/go
aifn build -l java -o dist/java
aifn build -l csharp -o dist/csharp
```

### Maintenance: Recompile / remove packages

```bash
aifn build                       # Recompile without re-downloading
aifn uninstall my-summarizer     # Remove packages no longer needed
```

---

## Next Steps

- **First time using AIFunc?** → [Quick Start](./01-quick-start)
- **View the runtime API?** → [Runtime API Reference](./02-api)
- **Want to create your own package?** → [Create an AIFunc Package](./05-create-package)
- **Understand the internals?** → [How It Works](./04-how-it-works)
