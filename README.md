<div align="center">

<h1>AIFunc</h1>

<p><strong>AI as Functions, Prompts as Code</strong></p>

<p>Use AI like calling a function, share AI like publishing an npm package</p>
<p>Strongly Typed · Testable · Cross-Language · Zero Dependencies · Model Agnostic · Git Native</p>

<p>
  <a href="./LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License"></a>
  <img alt="TypeScript / Node.js ≥ 18" src="https://img.shields.io/badge/TypeScript_/_Node.js-≥18-3178C6?logo=typescript&logoColor=white">
  <img alt="Python ≥ 3.10" src="https://img.shields.io/badge/Python-≥3.10-3776AB?logo=python&logoColor=white">
  <img alt="Go ≥ 1.23" src="https://img.shields.io/badge/Go-≥1.23-00ADD8?logo=go&logoColor=white">
  <img alt="Java ≥ 11" src="https://img.shields.io/badge/Java-≥11-ED8B00?logo=openjdk&logoColor=white">
  <img alt="C# / .NET ≥ 6" src="https://img.shields.io/badge/C%23_/_NET-≥6-512BD4?logo=dotnet&logoColor=white">
</p>

</div>

---

## Why AIFunc

Adding an AI feature to your product should be simple. In reality, it means: adopting a new framework, learning prompt engineering, writing glue code, and figuring out how to test it. A week of coding, and the core logic is just 3 lines.

**80% of real-world needs are actually simple: text in, structured data out. Stateless, no memory. It's essentially a function.**

AIFunc's goal is: **Let you use, manage, and share AI capabilities just like ordinary functions.**

### How does it work?

1. **Declarative package definition**: An AIFunc package is just a folder containing a package.json, an api.json, and a few prompts/*.md files. No code, no runtime.
2. **CLI cross-language compilation**: With a single `aifn` CLI command, these declaration files are "compiled" into native language code for your current project (TypeScript, Python, Go, Java, C#).
3. **Zero runtime dependencies**: The generated code only depends on the language's native standard library. You don't need to `npm install` or `pip install` any heavy framework — just `import` and use.

### Why is this approach better?

- **No new concepts, no architecture invasion**: What you get are strongly typed native functions. Use `if/else` to chain logic, use native arrays to manage context — say goodbye to complex Agent orchestration frameworks.
- **Testing out of the box**: The generated code comes with mock data — set `mock: true` for offline execution. CI pipelines don't need real API keys.
- **Git-native sharing**: Write a package, push it to a Git repository, and others can install it with a single command. No need for npm/PyPI — fork, PR, and version control reuse your existing Git workflow directly.

---

## Quick Start

```bash
# Install the CLI
brew tap aifunc-dev/aifn && brew install aifn   # macOS/Linux
scoop bucket add aifn https://github.com/aifunc-dev/scoop-aifn && scoop install aifn  # Windows
```

```bash
# Install an AI function
aifn install github:aifunc-dev/aifunc-packages/summarize
```

### TypeScript

```typescript
import { summarize, AIFuncConfig, SummarizeInput } from './aifunc/summarize';

const config: AIFuncConfig = { mock: true };

const text =
  'The James Webb Space Telescope captured its first full-color images in July 2022, ' +
  'revealing thousands of galaxies in a single image.';

async function main() {
  const result = await summarize(config, { text, maxLength: 30 } as SummarizeInput);
  console.log(`Summary   : ${result.summary}`);   // ← IDE autocomplete, type safe
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

<details>
<summary>Go</summary>

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

</details>

<details>
<summary>Java</summary>

```java
import aifunc.summarize.Summarize;
import aifunc.summarize.SummarizeTypes.SummarizeInput;
import aifunc.summarize.SummarizeTypes.SummarizeOutput;
import aifunc._engine.java.v0_1_0.Types.AIFuncConfig;

AIFuncConfig config = AIFuncConfig.builder().mock(true).build();

String text = "The James Webb Space Telescope captured its first full-color images in July 2022, " +
              "revealing thousands of galaxies in a single image.";

SummarizeInput input = new SummarizeInput(text, 30);
SummarizeOutput result = Summarize.summarize(config, input);
System.out.println("Summary   : " + result.getSummary());
System.out.println("Word count: " + result.getWordCount());
```

</details>

<details>
<summary>C#</summary>

```csharp
using Aifunc;
using Aifunc.Summarize;

var config = new AIFuncConfig { Mock = true };

var text =
    "The James Webb Space Telescope captured its first full-color images in July 2022, " +
    "revealing thousands of galaxies in a single image.";

var result = await Summarize.SummarizeAsync(config, new SummarizeTypes.SummarizeInput(text, 30));
Console.WriteLine($"Summary   : {result.Summary}");
Console.WriteLine($"Word count: {result.WordCount}");
```

</details>

To connect to a real model, simply replace `mock: true` with your actual `baseURL`, `model`, and `apiKey`. Any endpoint compatible with the OpenAI protocol is supported.

```typescript
// Same config fields for TypeScript / Python / Go / Java / C#
const config: AIFuncConfig = {
  baseURL: 'https://your-api-endpoint/v1',
  model: 'your-model-name',
  apiKey: 'your-api-key',
};
```

> Full runnable code available at [examples/go/hello-aifunc](./examples/go/hello-aifunc), [examples/typescript/hello-aifunc](./examples/typescript/hello-aifunc), [examples/python/hello-aifunc](./examples/python/hello-aifunc), [examples/java/hello-aifunc](./examples/java/hello-aifunc), [examples/csharp/hello-aifunc](./examples/csharp/hello-aifunc)

---

## Use Cases & Boundaries

AIFunc's design philosophy is **"reduce AI to ordinary functions"**, which gives it very clear boundaries:

✅ **Great for**
- **Integrating into existing systems**: Seamlessly add AI capabilities without changing your current architecture.
- **Relatively fixed business workflows**: Such as data cleaning, information extraction (Text to JSON), intent recognition, text classification, content summarization, etc.
- **Multi-turn conversations**: **Fully supported.** However, you manage context state through your own business code (e.g., using native arrays or a database), passing the conversation history to the AI function each time. No complex Agent memory framework needed.
- **Multi-model switching and adaptation**: Compatible with the OpenAI protocol under the hood — just modify the `config` to switch between different vendors' LLMs at any time, with zero code changes.

❌ **Not suitable for**
- **Highly open-ended autonomous Agents**: Tasks that require AI to independently plan and loop through tool chains.
  AIFunc advocates using code to orchestrate control flow, rather than letting AI act freely.

---

## Orchestrate AI Like Writing Business Code

AI functions can be freely composed just like ordinary functions. Chain business logic using the control flow you're already familiar with:

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

Use `if`, `switch` for control flow, manage context and data with `array` or databases — native language features are all you need to chain logic together. No extra orchestration abstractions, seamless integration into existing systems.

Full examples: [TypeScript](./examples/typescript/customer-support) / [Python](./examples/python/customer-support) / [Go](./examples/go/customer-support) / [Java](./examples/java/customer-support) / [C#](./examples/csharp/customer-support)

Multi-turn conversation with memory examples: [TypeScript](./examples/typescript/chat-with-context) / [Python](./examples/python/chat-with-context) / [Go](./examples/go/chat-with-context) / [Java](./examples/java/chat-with-context) / [C#](./examples/csharp/chat-with-context)

Streaming examples: [TypeScript](./examples/typescript/chat-stream) / [Python](./examples/python/chat-stream) / [Go](./examples/go/chat-stream) / [Java](./examples/java/chat-stream) / [C#](./examples/csharp/chat-stream) — also [all-packages-stream](./examples/typescript/all-packages-stream)

---

## Features

**Strongly Typed** — Full type definitions for inputs and outputs, IDE autocomplete, compile-time errors for misspelled fields.

**Testable** — Every package includes mock data. Set `mock: true` for offline execution. CI requires no API key — zero-cost testing.

**Cross-Language** — One package definition (`api.json` + `package.json` + `prompts/`), compiled to TypeScript, Python, Go, Java, C# with consistent behavior.

**Zero Dependencies** — The runtime engine is pure source code generated into your project. Each language uses only native libraries with no third-party dependencies.

**Model Agnostic** — Supports any endpoint compatible with the OpenAI protocol. Switch models by changing config only — zero code changes.

**Streaming** — Packages can stream plain-text tokens via `"x-delivery-mode": "stream"`, consumed with native async iterators / channels.

**Git Native** — Compiled artifacts are committed to Git. Team members clone and import directly — no CLI installation needed. Version management, code review, and access control all reuse your Git workflow.

> In one sentence: **AIFunc turns AI capabilities into ordinary functions in your codebase — strongly typed, testable, and version-controlled.**

---

## Sharing & Reuse

Write an AI function, push to Git, and anyone can install it with a single command:

```bash
# Others install your published package
aifn install github:your-name/your-packages/summarize

# Local packages work too
aifn install ../shared-packages/translate
```

No need to publish to npm, no registry required. **Your Git repository is the package registry.** Fork, PR, tag, access control — your existing Git workflow is reused directly.

Want to create your own package? All you need is a `package.json` for the package definition, an `api.json` for the interface definition, and a `prompts/` directory for prompts:

```bash
aifn create my-analyzer
# Generates the package skeleton — fill in package.json, api.json, and prompts
```

> See [Creating AIFunc Packages](./docs_cn/05-create-package.md) and [Sharing & Publishing](./docs_cn/07-sharing.md) for details

---

## Team Collaboration

| Role             | Responsibility                                      |
| ---------------- | --------------------------------------------------- |
| AI Engineer      | Create packages, write prompts, tune performance    |
| App Developer    | `import` → call — no prompt knowledge needed        |
| Platform Lead    | Manage package repositories, review PRs, control versions |

---

## Available AI Function Packages

The official package repository provides ready-to-use common AI functions:

| Package | Purpose |
| --- | --- |
| `summarize` | Text summarization |
| `analyze-sentiment` | Sentiment analysis |
| `recognize-intent` | Intent recognition |
| `classify` | General text classification |
| `extract-json` | Extract structured fields from text |
| `extract-entities` | Entity extraction (names, locations, orgs, etc.) |
| `extract-keywords` | Keyword extraction |
| `detect-language` | Language detection |
| `translate` | Translation |
| `rewrite` | Text rewriting |
| `chat` | Single-turn reply with optional context |
| `generate-reply` | Generate replies |
| `generate-title` | Generate titles |
| `generate-slug` | Generate URL slugs |
| `generate-email` | Generate emails |
| `generate-post` | Generate articles |
| `answer-question` | Question answering |
| `score-quality` | Content quality scoring |
| `chat-stream` | Stream a reply with optional context |
| `answer-stream` | Detailed question answering |
| `explain-stream` | Explain a concept, code, or term |
| `article-stream` | Full article from a title and outline |
| `write-stream` | Long-form writing: articles, reports, docs |
| `translate-stream` | Long document translation |
| `review-stream` | Code and document review with findings |

Install any package:

```bash
aifn install github:aifunc-dev/aifunc-packages/summarize  # Shorthand mode
aifn install https://github.com/aifunc-dev/aifunc-packages/tree/main/summarize # Full URL mode
```

> Full package list and documentation at [aifunc-packages](https://github.com/aifunc-dev/aifunc-packages)

---

## Documentation

- [Quick Start](./docs_cn/01-quick-start.md) — Get up and running in 5 minutes
- [Runtime API Reference](./docs_cn/02-api.md)
- [CLI Command Reference](./docs_cn/03-cli.md)
- [How It Works](./docs_cn/04-how-it-works.md)
- [Creating AIFunc Packages](./docs_cn/05-create-package.md)
- [Package Format Specification](./docs_cn/06-spec.md)
- [Sharing & Publishing](./docs_cn/07-sharing.md)
- [Team Collaboration](./docs_cn/08-team-workflow.md)

---

Licensed under [Apache-2.0](./LICENSE).