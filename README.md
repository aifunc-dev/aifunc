<div align="center">
<!---->
<h1>AIFunc</h1>
<!---->
<p><strong>AI as Function. Prompt as Code.</strong></p>
<!---->
<p>Call AI like a function. Share AI like an npm package.</p>
<p>Typed · Testable · Cross-language · Zero-dependency · Model-agnostic · Git-native</p>
<!---->
<p>
  <a href="./LICENSE"><img src="https://img.shields.io/badge/License-Apache_2.0-blue.svg" alt="License"></a>
  <a href="https://github.com/aifunc/cli"><img src="https://img.shields.io/badge/CLI-Go-00ADD8?logo=go&logoColor=white" alt="Go CLI"></a>
  <img alt="Node.js ≥ 18" src="https://img.shields.io/badge/Node.js-≥18-339933?logo=nodedotjs&logoColor=white">
  <img alt="Python ≥ 3.10" src="https://img.shields.io/badge/Python-≥3.10-3776AB?logo=python&logoColor=white">
  <img alt="Go ≥1.23" src="https://img.shields.io/badge/Go-≥1.23-00ADD8?logo=go&logoColor=white">
  <img alt="TypeScript types" src="https://img.shields.io/badge/TypeScript-types-3178C6?logo=typescript&logoColor=white">
</p>
<!---->
</div>

---

## Why AIFunc

Adding an AI feature to your product should be simple.

The reality: adopt a new framework, learn new orchestration abstractions, study prompt engineering, write glue code, figure out how to test, worry about switching models. A week of work for 3 lines of core logic.

80% of real-world AI needs are straightforward: **text in, structured data out. Stateless. No memory.** That's a function.

AIFunc makes this work out of the box. No new concepts, no architecture invasion. Just import.

---

## Quick Start

```bash
# Install CLI
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

// Mock mode: no API key needed, works offline
const config: AIFuncConfig = { mock: true };

const text =
  'The James Webb Space Telescope captured its first full-color images in July 2022, ' +
  'revealing thousands of galaxies in a single image.';

async function main() {
  const result = await summarize(config, { text, maxLength: 30 } as SummarizeInput);
  console.log(`Summary   : ${result.summary}`);   // ← IDE autocomplete, type-safe
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

To connect a real model, replace `mock: true` with your actual `baseURL`, `model`, and `apiKey`. Any OpenAI-compatible endpoint works.

```typescript
// TypeScript / Python / Go — same config fields
config = { baseURL: "https://your-api-endpoint/v1", model: "your-model-name", apiKey: "your-api-key" }
```

> See fully runnable examples: [examples/go/hello-aifunc](./examples/go/hello-aifunc), [examples/typescript/hello-aifunc](./examples/typescript/hello-aifunc), [examples/python/hello-aifunc](./examples/python/hello-aifunc)

---

## Compose Multiple AI Functions

AI functions compose like regular functions. Use familiar control flow to wire business logic:

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

`if`, `switch`, `Promise.all` — control flow you already know. Nothing new to learn.

> Full examples: [examples/typescript/customer-support](./examples/typescript/customer-support), [examples/python/customer-support](./examples/python/customer-support) and [examples/go/customer-support](./examples/go/customer-support)

> **Want multi-turn conversations with memory and sliding window?** No heavy Agent framework needed — just use native arrays to manage context.
> See [examples/typescript/chat-with-context](./examples/typescript/chat-with-context), [examples/python/chat-with-context](./examples/python/chat-with-context), [examples/go/chat-with-context](./examples/go/chat-with-context)

---

## Features

**Typed** — Full type definitions for inputs and outputs. IDE autocomplete. Misspell a field? Caught at compile time.

**Testable** — Every package ships with mock data. `mock: true` runs offline. CI without API keys, zero cost.

**Cross-language** — One package definition (`api.json` + `package.json` + `prompts/`) compiles to TypeScript, Python, and Go with identical behavior.

**Zero-dependency** — The runtime engine is pure source code generated into your project. Every language runtime uses only its native standard library — no third-party dependencies of any kind.

**Model-agnostic** — Works with any OpenAI-compatible endpoint. Switch models by changing config, zero code changes.

**Git-native** — Compiled output commits to Git. Team members clone and import directly, no CLI needed. Version control, code review, and access control all reuse your existing Git workflow.

---

## Share and Reuse

Write an AI function, push to Git, anyone installs with one command:

```bash
# Others install your published package
aifn install github:your-name/your-packages/summarize

# Local packages work too
aifn install ../shared-packages/translate
```

No npm publish, no registry. **A Git repo is a package registry.** Fork, PR, Tag, access control — your existing Git workflow, directly reused.

Want to create your own package? Just a `package.json` for metadata, an `api.json` for the interface, and a `prompts/` directory for the prompt:

```bash
aifn create my-analyzer
# Generates package skeleton — fill in package.json, api.json, and prompt
```

> See [Create an AIFunc Package](./docs/05-create-package.md) and [Sharing & Publishing](./docs/07-sharing.md)

---

## Team Workflow

| Role | Responsibility |
|:---|:---|
| AI Engineer | Create packages, write prompts, tune quality |
| App Developer | `import` → call, no prompt knowledge needed |
| Platform Lead | Manage package repo, review PRs, control versions |

---

## Available Packages

The official package registry provides ready-to-use AI functions:

| Package | Purpose |
|:---|:---|
| `summarize` | Text summarization |
| `analyze-sentiment` | Sentiment analysis |
| `recognize-intent` | Intent recognition |
| `classify` | General text classification |
| `extract-json` | Extract structured fields from text |
| `extract-entities` | Entity extraction (names, places, orgs) |
| `extract-keywords` | Keyword extraction |
| `detect-language` | Language detection |
| `translate` | Translation |
| `rewrite` | Text rewriting |
| `generate-reply` | Reply generation |
| `generate-title` | Title generation |
| `generate-slug` | URL slug generation |
| `generate-email` | Email generation |
| `generate-post` | Post generation |
| `answer-question` | Question answering |
| `score-quality` | Content quality scoring |

Install any package:

```bash
aifn install github:aifunc-dev/aifunc-packages/summarize  # Short form
aifn install https://github.com/aifunc-dev/aifunc-packages/tree/main/summarize  # Full URL
```

> Full package list and docs: [aifunc-packages](https://github.com/aifunc-dev/aifunc-packages)

---

## Documentation

- [Quick Start](./docs/01-quick-start.md) — Up and running in 5 minutes
- [Runtime API Reference](./docs/02-api.md)
- [CLI Reference](./docs/03-cli.md)
- [How It Works](./docs/04-how-it-works.md)
- [Create a Package](./docs/05-create-package.md)
- [Package Spec](./docs/06-spec.md)
- [Sharing & Publishing](./docs/07-sharing.md)
- [Team Workflow](./docs/08-team-workflow.md)

---

Licensed under [Apache-2.0](./LICENSE).
