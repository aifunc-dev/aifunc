# Sharing & Publishing

> **Target audience**: Package authors
> **Content**: How to publish your AIFunc package via Git for others to install and use
> **Prerequisites**: Completed [Create a Package](./05-create-package)

---

## Pre-publish Checklist

Run validation before publishing to ensure your package conforms to the spec:

```bash
aifn validate ./my-summarizer
```

Checklist:

- [ ] `package.json` has `name`, `version`, and `description` filled in completely
- [ ] `api.json` accurately describes the function interface with input/output schemas
- [ ] `prompts/general.md` has been tested to produce results conforming to the schema
- [ ] Added a `README.md` explaining usage and showing example results
- [ ] Added a `LICENSE` file declaring the license

---

## Repository Organization

It's recommended to place multiple packages in the same repository, one subdirectory per package. The official package repository [aifunc-packages](https://github.com/AIfunc-dev/aifunc-packages) is organized this way:

```text
my-packages/
├── summarize/
│   ├── package.json
│   ├── api.json
│   ├── model-params.json
│   ├── mock.json
│   ├── prompts/
│   │   └── general.md
│   ├── README.md
│   └── LICENSE
├── translate/
│   └── ...
├── classify/
│   └── ...
└── README.md
```

Specify the subdirectory path when installing:

```bash
aifn install github:your-name/my-packages/summarize
aifn install github:your-name/my-packages/translate
```

---

## Supported Installation Source Formats

### Shorthand Format (Recommended)

| Format | Example |
|:---|:---|
| GitHub shorthand | `github:owner/repo/path` |
| Local path | `./my-summarizer`, `../shared/my-summarizer` |

### Full URL Format

Copy the repository directory URL directly from your browser's address bar:

```bash
aifn install https://github.com/owner/repo/tree/main/summarize
```

---

## Contributing to the Official Package Repository

Submit a PR to the official package repository to have your package included as a community package:

- GitHub: [aifunc-dev/aifunc-packages](https://github.com/aifunc-dev/aifunc-packages)

> [!NOTE]
> The official repository currently only accepts commonly used, general-purpose packages. Packages that are too domain-specific or experimental are better shared via personal repositories.

Submission process:

1. **Fork the official repository**
2. **Add your package directory**: Add your package directory to the fork
3. **Run validation**: Ensure `aifn validate` passes
4. **Submit a Pull Request**: Include usage instructions and example results

---

## Next Steps

- **Want to manage private packages within a team?** → [Team Workflow](./08-team-workflow)
- **View all CLI commands?** → [CLI Command Reference](./03-cli)
- **Review the package format spec?** → [Package Format Spec](./06-spec)
