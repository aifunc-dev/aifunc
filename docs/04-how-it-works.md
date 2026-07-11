# How AIFunc Works

> **Target audience**: Developers who want to understand the internals
> **Content**: The complete flow from package installation to runtime execution, including CLI compilation, project structure, and runtime call chain

> [!NOTE]
> Package authors only need to read the [Package Format Spec](./06-spec) to create packages — understanding this document is not required.

---

## 1. Components

| Component | Language | Purpose | Do you interact with it? |
|:---|:---|:---|:---|
| **aifn CLI** | Go | Package management, dependency resolution, code generation | ✅ You run commands in the terminal |
| **Engine SDK** | Same as target language (e.g., TypeScript, Python, or Go) | Runtime: loads compiled artifacts, calls model, validates output | ❌ Generated as source code in your project |
| **Generated functions** | Target language | The strongly-typed entry point you import | ✅ You call it in code |

The CLI is a standalone Go binary, independent of your project's language. The Engine SDK is fetched on demand as pure source files, implemented using only the language's standard library with zero external dependencies. It lives alongside the generated function code in your project directory — no need to install anything via npm, pip, or `go get`.

---

## 2. Installation and Compilation Flow

When you run `aifn install`, the CLI automatically performs these core steps:

```text
Parse config ──────► Download source & Engine SDK ──────► Compile & language wrapping
```

1. **Parse & download**: The CLI reads `aifunc.json` from your project, gets the list of packages to install, and downloads each package's raw files to the local cache directory (`.aifunc/packages/`).
2. **Fetch Engine SDK**: The CLI analyzes each package's declared `engineVersion`, automatically downloads the corresponding target-language Engine SDK source files to the cache directory (`.aifunc/_engine/`).
3. **Compile & wrap (core mechanism)**: To avoid issues that traditional bundlers (like Webpack/Vite) or language runtimes may have when loading raw data files, the CLI no longer outputs plain config files separately. Instead, it **wraps API definitions, prompt templates, model parameters, etc. using your target language**.
   * TypeScript projects get `.ts` files containing config objects.
   * Python projects get `.py` files containing config dictionaries.
   * Go projects get `.go` files containing config maps and struct literals.
   * Corresponding strongly-typed interface files and mock data wrapper files are also generated.
4. **Link output**: Finally, the CLI outputs the generated strongly-typed functions, language-wrapped artifacts, and required Engine SDK source files to your specified project directory (e.g., `src/aifunc/`).

---

## 3. Project Structure After Installation

The project is divided into a cache area (not committed to Git) and a compiled artifacts area (committed to Git). Examples for TypeScript, Python, and Go:

### TypeScript Project Example

```text
your-project/
├── aifunc.json                          ← Package management config (you edit this)
├── aifunc-lock.json                     ← Version lock file
│
├── .aifunc/                             ← CLI download cache (add to .gitignore)
│   ├── packages/summarize/              ← Package raw source files
│   └── _engine/v0.1.0/                  ← Downloaded Engine raw files
│
└── src/
    └── aifunc/                          ← Compiled artifact output directory (commit to Git)
        ├── summarize/                   ← AI function package
        │   ├── index.ts                 ← Function entry point
        │   ├── summarize.types.ts       ← I/O type definitions
        │   ├── summarize.aifunc.ts      ← Compiled artifact (prompt & API spec TS wrapper)
        │   └── summarize.mock.ts        ← Offline test mock data (TS wrapper)
        │
        └── _engine/                     ← Runtime SDK (local source-level dependency)
            └── typescript/
                 └── v0.1.0/
                      ├── index.ts
                      ├── runtime.ts
                      └── ...
```

### Python Project Example

```text
your-project/
├── aifunc.json
├── aifunc-lock.json
│
├── .aifunc/                             ← Cache directory (add to .gitignore)
│   ├── packages/summarize/
│   └── _engine/v0_1_0/
│
└── aifunc/                              ← Compiled artifact output directory (commit to Git)
    ├── __init__.py                      ← Empty file, makes directory a Python package
    ├── py.typed                         ← Empty file, PEP 561 type marker
    ├── summarize/
    │   ├── __init__.py                  ← Function entry point
    │   ├── summarize_types.py           ← I/O type definitions
    │   ├── summarize_aifunc.py          ← Compiled artifact (prompt & API spec Python wrapper)
    │   └── summarize_mock.py            ← Offline test mock data (Python wrapper)
    │
    └── _engine/
        └── python/
            └── v0_1_0/
                  ├── __init__.py
                  ├── runtime.py
                  └── ...
```

### Go Project Example

```text
your-project/
├── aifunc.json
├── aifunc-lock.json
├── go.mod
│
├── .aifunc/                             ← Cache directory (add to .gitignore)
│   ├── packages/summarize/
│   └── _engine/go/v0.1.0/
│
└── aifunc/                              ← Compiled artifact output directory (commit to Git)
    ├── summarize/                       ← AI function package
    │   ├── summarize.go                 ← Function entry point (package summarize)
    │   ├── summarize_types.go           ← I/O struct definitions
    │   └── summarize_aifunc.go          ← Compiled artifact (prompt & API spec Go wrapper)
    │
    └── _engine/
        └── go/
            └── v0.1.0/
                  ├── aifunc.go          ← Public API
                  ├── runtime.go
                  ├── types.go
                  └── ...
```

> [!NOTE]
> Files within packages use different naming conventions depending on the target language (e.g., TS uses `.` separators, Python and Go use `_` separators). All code dependencies are self-contained within the generated directory.

---

## 4. File Responsibilities

Using the generated `summarize` package as an example, internal files have clearly defined roles:

| File (logical name) | Responsibility |
|:---|:---|
| `entry` (e.g., `index.ts`, `summarize.go`) | **Function entry point**: Creates and exports the AI function instance. Your business code imports this directly. |
| `types` (e.g., `.types.ts`, `_types.go`) | **Interface file**: Input and output structure type definitions, providing strong typing support. |
| `aifunc` (e.g., `.aifunc.ts`, `_aifunc.go`) | **Core artifact**: Merged result of prompt template, API spec, and model config, wrapped in the target language to ensure safe runtime loading. |
| `mock` (e.g., `.mock.ts`, `_mock.py`) | **Mock data**: Input-to-output mapping data, also wrapped in the target language, used for offline test mode. (Go embeds mock data directly in the entry file.) |

---

## 5. Runtime Call Chain

```text
Your application code
  │  await summarize(config, { text: "...", maxLength: 20 })        ← TypeScript/Python
  │  summarize.Summarize(ctx, config, SummarizeInput{...})          ← Go
  ▼
Generated function entry (aifunc/summarize/index.ts or summarize.go)
  │  Strongly-typed entry, no business logic
  │  References types, aifunc artifact from same directory, and _engine
  ▼
Engine SDK (aifunc/_engine/vX.Y.Z/runtime.ts|runtime.py|runtime.go)
  │  ① Parse the language-native wrapped .aifunc config object
  │  ② Validate input data structure
  │  ③ Render prompt template (substitute variables)
  │  ④ Call AI model API (supports any OpenAI-compatible endpoint)
  │  ⑤ Parse model response
  │  ⑥ Validate output data structure
  │  ⑦ Validation passes → return strongly-typed result
  ▼
AI Model API
```

---

## 6. Configuration File

`aifunc.json` is the project-level config that controls code generation language and output path:

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

| Field | Description |
|:---|:---|
| `language` | Target language, determines the language wrapper format and Engine SDK language (currently supports `typescript`, `python`, `go`) |
| `outputDir` | Output path for compiled artifacts (including generated functions and Engine SDK) |
| `alias` | (TS only) Used for tsconfig paths alias configuration |
| `packages` | Mapping of package names to their installation sources (supports github paths, local paths) |

---

## Next Steps

- **Want to create your own package?** → [Create an AIFunc Package](./05-create-package)
- **View the runtime API?** → [Runtime API Reference](./02-api)
- **View the full package format definition?** → [Package Format Spec](./06-spec)
- **View all CLI commands?** → [CLI Command Reference](./03-cli)
