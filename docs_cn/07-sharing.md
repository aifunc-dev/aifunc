# 分享与发布

> **目标读者**：包作者
> **本文内容**：如何将创建好的 AIFunc 包通过 Git 发布，供他人安装使用
> **前置条件**：已完成 [创建包教程](./05-create-package)

---

## 发布前检查

发布前运行校验，确保包格式合规：

```bash
aifn validate ./my-summarizer
```

检查清单：

- [ ] `package.json` 中的 `name`、`version`、`description` 填写完整
- [ ] `api.json` 中的 input/output schema 准确描述函数接口
- [ ] `prompts/general.md` 经过测试能产出符合 schema 的结果
- [ ] 添加 `README.md` 说明使用方法和效果示例
- [ ] 添加 `LICENSE` 文件声明许可证

---

## 仓库组织方式

推荐将多个包放在同一仓库中，每个包一个子目录。官方包仓库 [aifunc-packages](https://github.com/AIfunc-dev/aifunc-packages) 就是这样组织的：

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

安装时指定子目录路径：

```bash
aifn install github:your-name/my-packages/summarize
aifn install github:your-name/my-packages/translate
```

---

## 支持的安装来源格式

### 简写模式（推荐）

| 格式 | 示例 |
|:---|:---|
| GitHub 简写 | `github:owner/repo/path` |
| 本地路径 | `./my-summarizer`、`../shared/my-summarizer` |

### 完整 URL 模式

直接从浏览器地址栏复制仓库目录的 URL：

```bash
aifn install https://github.com/owner/repo/tree/main/summarize
```

---

## 提交到官方包仓库

向官方包仓库提 PR，经审核后收录为社区包：

- GitHub: [aifunc-dev/aifunc-packages](https://github.com/aifunc-dev/aifunc-packages)

> [!NOTE]
> 官方包仓库目前只接收常用、通用的包。领域过于狭窄或实验性质的包建议通过个人仓库分享。

提交流程：

1. **Fork 官方仓库**
2. **添加包目录**：在 fork 中添加你的包目录
3. **运行校验**：确保 `aifn validate` 通过
4. **提交 Pull Request**：附带使用说明和效果示例

---

## 接下来

- **想在团队内管理私有包？** → [团队协作](./08-team-workflow)
- **查看 CLI 所有命令？** → [CLI 命令参考](./03-cli)
- **回顾包格式规范？** → [包格式规范](./06-spec)
