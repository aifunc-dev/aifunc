# Quick Start

> **Target audience**: All developers who want to use AIFunc
> **Content**: Complete workflow from installing the CLI to calling an AI function and getting results in 5 minutes
> **Prerequisites**: Node.js 18+, Python 3.10+, Go 1.23+, or Java 11+

---

## Step 1: Install CLI

```bash
# Install CLI
brew tap aifunc-dev/aifn && brew install aifn   # macOS/Linux
scoop bucket add aifn https://github.com/aifunc-dev/scoop-aifn && scoop install aifn  # Windows
```

---

## Step 2: Install an AI Function Package

Run in your project root:

```bash
aifn install github:aifunc-dev/aifunc-packages/summarize
```

The CLI will automatically:
- Detect your project type (TypeScript / Python / Go / Java)
- Generate importable code (with type definitions and built-in mock data)
- Create a configuration file (if one doesn't exist)

---

## Step 3: Write Code to Call It

### TypeScript

```typescript
import { summarize, AIFuncConfig, SummarizeInput } from './aifunc/summarize';

const config: AIFuncConfig = { mock: true };

const text =
  'The James Webb Space Telescope captured its first full-color images in July 2022, ' +
  'revealing thousands of galaxies in a single image.';

async function main() {
  const result = await summarize(config, { text, maxLength: 30 } as SummarizeInput);
  console.log(`Summary   : ${result.summary}`);
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
	fmt.Printf("Summary   : %s\n", result.Summary)
	fmt.Printf("Word count: %d\n", result.WordCount)
}
```

Your IDE provides full type hints and autocompletion.

### Java

```java
import aifunc.AIFuncConfig;
import aifunc.summarize.Summarize;
import aifunc.summarize.SummarizeTypes.SummarizeInput;

AIFuncConfig config = AIFuncConfig.builder().mock(true).build();

String text = "The James Webb Space Telescope captured its first full-color images in July 2022, " +
              "revealing thousands of galaxies in a single image.";

Summarize.summarize(config, new SummarizeInput(text, 30))
        .thenAccept(result -> {
            System.out.println("Summary   : " + result.getSummary());
            System.out.println("Word count: " + result.getWordCount());
        })
        .join();
```

---

## Step 4: Connect to a Real Model

The examples above use `mock: true`, which runs the full workflow without needing an API Key.

When you're ready to connect to a real model, configure real parameters:

### TypeScript

```typescript
import { summarize, AIFuncConfig, SummarizeInput } from './aifunc/summarize';

const config: AIFuncConfig = {
  baseURL: 'https://your-api-endpoint/v1',
  model: 'your-model-name',
  apiKey: 'your-api-key',
};

const text =
  'The James Webb Space Telescope captured its first full-color images in July 2022, ' +
  'revealing thousands of galaxies in a single image.';

async function main() {
  const result = await summarize(config, { text, maxLength: 30 } as SummarizeInput);
  console.log(`Summary   : ${result.summary}`);
  console.log(`Word count: ${result.wordCount}`);
}

main().catch(console.error);
```

### Python

```python
import asyncio
from aifunc.summarize import summarize, AIFuncConfig, SummarizeInput

config = AIFuncConfig(
    base_url="https://your-api-endpoint/v1",
    model="your-model-name",
    api_key="your-api-key",
)

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
	config := &summarize.AIFuncConfig{
		BaseURL: "https://your-api-endpoint/v1",
		Model:   "your-model-name",
		APIKey:  "your-api-key",
	}

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

### Java

```java
import aifunc.AIFuncConfig;
import aifunc.summarize.Summarize;
import aifunc.summarize.SummarizeTypes.SummarizeInput;

AIFuncConfig config = AIFuncConfig.builder()
        .baseUrl("https://your-api-endpoint/v1")
        .model("your-model-name")
        .apiKey("your-api-key")
        .build();

String text = "The James Webb Space Telescope captured its first full-color images in July 2022, " +
              "revealing thousands of galaxies in a single image.";

Summarize.summarize(config, new SummarizeInput(text, 30))
        .thenAccept(result -> {
            System.out.println("Summary   : " + result.getSummary());
            System.out.println("Word count: " + result.getWordCount());
        })
        .join();
```
Works with any OpenAI-compatible API endpoint.

Run the code:

```bash
# TypeScript: install deps, compile and run
npm install
npm run build
npm run start

# Or use dev mode (ts-node + path alias resolution)
npm run dev

# Python
python main.py

# Go
go run main.go

# Java
javac Main.java
java Main
```

---

## Done

Just like using a regular npm/pip package: install → import → call → get results.

No need to deal with Prompts, Tokens, or Temperature.

---

## About Mock Mode

Every AI function package comes with a `mock.json` that defines preset input-output mappings. Set `mock: true` in your code to activate:

- Prioritizes matching mock data entries that exactly match the input
- When no match is found, automatically generates type-conforming return values based on the output schema
- Suitable for offline development, unit testing, and CI pipelines

No environment variables needed — fully controlled by code.

---

## Next Steps

- **View the full API reference?** → [Runtime API Reference](./02-api)
- **Create your own AI function?** → [Create a Package](./05-create-package)
- **Publish a package to your team?** → [Sharing & Publishing](./07-sharing)

---

## FAQ

<details>
<summary><b>Using a different service endpoint?</b></summary>

Specify the corresponding `baseURL` in the config:

```typescript
const config: AIFuncConfig = {
  baseURL: 'https://your-provider-api/v1',
  model: 'your-model-name',
  apiKey: 'your-api-key',
};
```

```typescript
// Locally deployed inference services usually don't need an apiKey
const config: AIFuncConfig = {
  baseURL: 'http://localhost:11434/v1',
  model: 'your-local-model',
  apiKey: '',
};
```
</details>

<details>
<summary><b>Should generated code be committed to Git?</b></summary>

Yes, recommended. Team members can use it directly after cloning without regenerating.

The `.aifunc/` cache directory is automatically added to `.gitignore`.
</details>

<details>
<summary><b>The CLI is written in Go — will it conflict with my project?</b></summary>

No. The CLI is only used when running `aifn` commands.

The runtime is handled by generated TypeScript/Python/Go/Java code with zero external dependencies.
</details>

<details>
<summary><b>Where do mock mode return values come from?</b></summary>

Each package directory has a `mock.json` with preset cases. The engine's lookup order:

1. Exact input match → return corresponding output
2. Fallback case with no input field → return its output
3. No match → auto-generate zero values based on `api.json` output schema

</details>
