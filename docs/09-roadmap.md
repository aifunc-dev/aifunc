# AIFunc Roadmap

Our goal is to make AI capabilities standard, typed, and predictable across all major programming languages. Here is where we are heading.

---

## Now (In Active Development)

Features currently being worked on and expected in the near future.

- [ ] **Go Language Support**: First-class generation for Go.
- [ ] **Java Support**: Bringing AI functions to the enterprise JVM ecosystem.
- [ ] **C# Support**: Strongly-typed C# wrapper generation.
- [ ] **Robust Structured Output Handling**: Comprehensive parsing and validation for complex structured responses.
- [ ] **Record/Replay Mocking**: Automatically record real API responses into mock files for zero-setup deterministic testing.
- [ ] **Model-Specific Parameter Adaptation**: Progressively adding strongly-typed, provider-specific parameter support with full IDE autocomplete.

> **Philosophy: Encapsulated by Protocol, Overridable by Caller**
> Model-specific parameters (like `temperature`, `top_p`) are **already defined and optimized by the package author** inside the AIFunc protocol. As a caller, you usually don't need to configure anything — the package automatically selects the right parameters based on the model you use. However, if you explicitly override them via the config object, your IDE provides exact, provider-specific autocomplete based on the `model` string.
>
> **For model providers**: We welcome collaboration. If you'd like AIFunc to support your model's full parameter set, open an issue or reach out directly.

---

## Next (Planned)

Confirmed features that are up next on our priority list.

- [ ] **Streaming Output**: Elegant, typed async iterators for streaming responses.
- [ ] **CLI Package Discovery**: `aifn search` to easily find community packages.
- [ ] **CLI Third-Party Plugin System**: Extensible plugin architecture for the CLI.
- [ ] **Rust Core**: Targeted for Protocol v1.0. Rust will serve as the high-performance foundational library for bindings to other languages (e.g., Swift, Ruby, Dart) — not a replacement for the first-class Java/Kotlin, Go, or C# implementations listed above.
- [ ] **Official Package Registry & Management Hub**: A centralized platform for discovering, publishing, versioning, and managing community-driven AI packages with full lifecycle support.

---

## Non-Goals

To maintain the "AI as a Function" philosophy, AIFunc will **NOT** build or support the following features natively.

- **State & Storage (e.g., Vector DBs, Native RAG)**
  AIFunc is strictly stateless. We do not manage memory or databases. If you need RAG, retrieve the context in your business logic and pass it to AIFunc as a standard argument.

- **Agentic Orchestration & Planning**
  We will not build DAG runners, autonomous planning, or complex agent loops. You do not need a new framework for orchestration; orchestrate your AI functions using standard programming constructs like `if`, `switch`, and `for` loops.