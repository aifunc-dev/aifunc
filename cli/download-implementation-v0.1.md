# aifn 下载功能实现原理（v0.1）

> 文档性质：内部开发说明  
> 适用范围：`aifn download <source...>` 以及未来 `aifn install <source...>` 复用的下载阶段  
> 当前状态：`download` 是调试命令，未来可以从对外帮助中隐藏，但底层下载函数会继续保留

---

## 1. 功能目标

下载功能负责把一个 AIFunc 包源码从指定来源获取到当前目录的本地缓存中：

```text
.aifunc/packages/{package.name}/
```

它只做三件事：

1. 解析用户输入的 `<source>`。
2. 获取源码到临时目录并校验包规范。
3. 校验通过后写入 `.aifunc/packages/`。

它不做：

- 编译；
- 代码生成；
- 写 lockfile；
- 安装后清理 compiled/generated；
- 初始化完整项目配置。

因此，`download` 可以作为调试命令使用，也可以作为未来 `install` 流程中的第一步。

---

## 2. 当前支持的 Source 格式

当前下载功能支持以下输入格式。

### 2.1 本地目录

```bash
aifn download ./my-package
```

解析结果：

```text
Kind = local
Path = ./my-package
```

本地目录会通过 `source.CopyLocal` 复制到临时 staging 目录，再执行 manifest 校验。

---

### 2.2 GitHub 简写

```bash
aifn download github:aifunc-dev/aifunc-packages/short-summary
```

解析结果：

```text
Kind     = git
Provider = github
Owner    = aifunc-dev
Repo     = aifunc-packages
Ref      = ""
SubPath  = short-summary
```

`Ref` 为空时，Git 使用远程仓库默认分支。

---

### 2.3 Gitee 简写

```bash
aifn download gitee:aifunc-dev/aifunc-packages/short-summary
```

解析结果：

```text
Kind     = git
Provider = gitee
Owner    = aifunc-dev
Repo     = aifunc-packages
Ref      = ""
SubPath  = short-summary
```

---

### 2.4 GitHub tree URL

```bash
aifn download https://github.com/aifunc-dev/aifunc-packages/tree/main/short-summary
```

解析结果：

```text
Kind     = git
Provider = github
Owner    = aifunc-dev
Repo     = aifunc-packages
Ref      = master
SubPath  = short-summary
```

---

### 2.5 Gitee tree URL

```bash
aifn download https://gitee.com/aifunc-dev/aifunc-packages/tree/master/short-summary
```

解析结果：

```text
Kind     = git
Provider = gitee
Owner    = aifunc-dev
Repo     = aifunc-packages
Ref      = master
SubPath  = short-summary
```

---

## 3. 模块分层

当前下载功能拆成四层。

```text
internal/cli/download.go
  CLI 命令层：参数接收、输出展示、批量下载结果汇总

internal/downloader/download.go
  下载编排层：创建 packages 目录、拉取源码、校验 manifest、写入缓存

internal/source/parse.go
  Source 解析层：把字符串解析为 Source 结构体

internal/source/git.go / local.go
  Source 获取层：Git 拉取或本地目录复制
```

这种分层的目的：

- `download` 命令未来可以隐藏，但 `internal/downloader.Download` 仍能被 `install` 复用；
- CLI 层不直接关心 GitHub/Gitee 的解析细节；
- source 层只负责“来源解析与获取”，不关心包规范；
- downloader 层负责把 source、manifest、workspace 串成完整下载流程。

---

## 4. 核心调用链

执行：

```bash
aifn download https://gitee.com/aifunc-dev/aifunc-packages/tree/master/short-summary
```

核心调用链如下：

```text
cmd/aifn/main.go
  -> cli.Execute(args)
    -> cli.NewRootCommand()
      -> newDownloadCommand(opts)
        -> runDownload(args, opts)
          -> workspace.FromCurrentDir()
          -> downloader.Download(raw, ws)
            -> os.MkdirAll(ws.PackagesPath(), 0755)
            -> fetchToStaging(raw)
              -> source.Parse(raw)
              -> source.FetchGit(src, staging)
                -> source.RepoURL(src)
                -> git clone --depth 1 --branch master https://gitee.com/aifunc-dev/aifunc-packages.git <tmp>
                -> source.CopyLocal(<tmp>/short-summary, staging)
            -> manifest.Validate(staging)
            -> source.CopyLocal(staging, .aifunc/packages/{package.name})
```

---

## 5. 下载流程详解

### 5.1 CLI 层：`internal/cli/download.go`

职责：

- 定义 `download <source...>` 命令；
- 支持别名 `dl`；
- 接收一个或多个 source；
- 对每个 source 调用 `downloader.Download`；
- 汇总成功和失败数量；
- 打印用户可读输出。

关键点：

```text
runDownload 不再调用 workspace.Require()
```

也就是说：

```bash
aifn download <source>
```

不要求当前目录已经执行过：

```bash
aifn init
```

如果当前目录没有 `.aifunc/packages/`，下载流程会自动创建。

这样设计的原因：

- `download` 是调试命令；
- 下载源码本身不依赖 `.aifunc/config.json`；
- 可以让用户快速验证远程包是否能被解析、拉取和校验。

---

### 5.2 下载编排层：`internal/downloader/download.go`

核心函数：

```go
func Download(raw string, ws workspace.Workspace) (Result, error)
```

职责：

1. 确保 `.aifunc/packages/` 存在；
2. 调用 `fetchToStaging(raw)` 获取源码到临时目录；
3. 调用 `manifest.Validate(staging)` 校验包规范；
4. 根据 manifest 中的 `name` 决定最终缓存目录；
5. 把 staging 目录复制到 `.aifunc/packages/{package.name}`。

为什么先 staging 再写入 packages：

- 避免下载失败时污染 `.aifunc/packages/`；
- 避免非法包覆盖已有缓存；
- manifest 校验通过后才落盘，行为更安全。

结果结构：

```go
type Result struct {
    Name string
}
```

当前只返回包名，后续可扩展：

- source 信息；
- resolved ref；
- commit hash；
- cache path；
- 是否覆盖已有包。

---

### 5.3 Source 解析层：`internal/source/parse.go`

核心函数：

```go
func Parse(raw string) (Source, error)
```

核心结构：

```go
type Source struct {
    Raw string

    Kind Kind
    Path string

    Provider Provider
    Owner    string
    Repo     string
    Ref      string
    SubPath  string
}
```

`Kind` 当前有两种：

```text
local
  本地目录

git
  GitHub/Gitee Git 仓库
```

`Provider` 当前有两种：

```text
github
gitee
```

解析顺序：

```text
1. trim 空白
2. 空字符串报错
3. github: 前缀 -> GitHub 简写
4. gitee: 前缀 -> Gitee 简写
5. http/https URL -> 尝试解析 GitHub/Gitee tree URL
6. name@version -> v0.1 暂不支持，报错
7. 其他输入 -> 本地目录
```

对于 tree URL，当前只支持：

```text
https://github.com/owner/repo/tree/ref/path
https://gitee.com/owner/repo/tree/ref/path
```

其中 `ref` 暂时只取 `tree` 后面的第一个路径段。

例如：

```text
/tree/master/short-summary
```

解析为：

```text
Ref     = master
SubPath = short-summary
```

当前不支持多段 branch 自动识别，例如：

```text
/tree/feature/foo/short-summary
```

因为仅从 URL 路径无法可靠判断：

```text
Ref     = feature/foo
SubPath = short-summary
```

还是：

```text
Ref     = feature
SubPath = foo/short-summary
```

这类能力可以在后续版本通过显式语法或远程 API 解决。

---

### 5.4 Git 获取层：`internal/source/git.go`

核心函数：

```go
func FetchGit(src Source, dstPath string) error
```

执行步骤：

1. 检查本机是否安装 `git`；
2. 创建临时目录；
3. 根据 provider 生成仓库 URL；
4. 执行浅克隆；
5. 如果存在 `SubPath`，从仓库子目录复制到目标 staging；
6. 清理临时目录。

Provider 到仓库 URL 的映射：

```text
github -> https://github.com/{owner}/{repo}.git
gitee  -> https://gitee.com/{owner}/{repo}.git
```

Git clone 参数：

```text
无 Ref：
git clone --depth 1 <repoURL> <tmp>

有 Ref：
git clone --depth 1 --branch <ref> <repoURL> <tmp>
```

当前采用浅克隆整个仓库，再复制子目录。

这样做的原因：

- v0.1 优先保证稳定和跨 Git 版本兼容；
- 实现简单，错误路径清晰；
- 后续可以替换为 sparse checkout 优化大仓库下载速度。

后续可优化为：

```bash
git clone --depth 1 --filter=blob:none --sparse <repoURL> <tmp>
git -C <tmp> sparse-checkout set <subpath>
```

---

### 5.5 本地复制层：`internal/source/local.go`

核心函数：

```go
func CopyLocal(srcPath, dstPath string) error
```

职责：

- 确认源路径存在；
- 确认源路径是目录；
- 删除目标目录；
- 递归复制源目录到目标目录。

复制时会跳过：

```text
.git
node_modules
*.tmp
```

这个函数被两类流程复用：

1. 本地 source 下载；
2. Git 拉取后从仓库目录复制子包；
3. staging 校验通过后复制到 `.aifunc/packages/{name}`。

---

## 6. 缓存目录规则

下载结果最终写入：

```text
.aifunc/packages/{package.name}/
```

注意：最终目录名来自 manifest 校验得到的 `package.name`，不是 URL 里的目录名。

例如：

```bash
aifn download https://gitee.com/aifunc-dev/aifunc-packages/tree/master/short-summary
```

如果包 manifest 中的 name 是：

```text
short-summary
```

最终目录为：

```text
.aifunc/packages/short-summary/
```

如果 source 路径名和 manifest name 不一致，以 manifest name 为准。

---

## 7. 错误处理策略

### 7.1 单个 source 失败

`downloader.Download` 遇到错误会直接返回 error。

常见错误：

- source 为空；
- 简写格式非法；
- GitHub/Gitee tree URL 格式非法；
- 未安装 Git；
- git clone 失败；
- subpath 不存在；
- manifest 校验失败；
- 复制文件失败。

### 7.2 多个 source 批量下载

CLI 层支持：

```bash
aifn download <source1> <source2> <source3>
```

处理策略：

- 每个 source 独立处理；
- 某个 source 失败，不阻止后续 source；
- 最后如果存在失败项，返回汇总错误。

示例：

```text
2 failed, 1 cached
```

---

## 8. 为什么 download 不要求 init

`aifn init` 的职责是创建完整项目结构和配置：

```text
.aifunc/config.json
.aifunc/packages/
.aifunc/compiled/
.aifunc/generated/
```

但下载功能只需要：

```text
.aifunc/packages/
```

因此 `download` 不调用：

```go
ws.Require()
```

而是在 `downloader.Download` 中直接执行：

```go
os.MkdirAll(ws.PackagesPath(), 0755)
```

这样用户可以在任意目录直接验证下载：

```bash
aifn download gitee:aifunc-dev/aifunc-packages/short-summary
```

这也符合 `download` 作为调试命令的定位。

---

## 9. 与 install 的关系

未来 `install` 应该复用：

```go
downloader.Download(raw, ws)
```

推荐流程：

```text
aifn install <source...>
  -> workspace.Require() 或 workspace.Ensure()
  -> downloader.Download(source, ws)
  -> compiler.Compile(package)
  -> codegen.Generate(artifact)
  -> lockfile.Update(...)
```

区别：

| 命令 | 是否要求 init | 是否下载 | 是否校验 | 是否编译 | 是否生成代码 |
|:---|:---:|:---:|:---:|:---:|:---:|
| download | 否 | 是 | 是 | 否 | 否 |
| install | 建议是 | 是 | 是 | 是 | 是 |

`download` 未来可以隐藏，但 `downloader` 包不应删除。

---

## 10. 当前测试覆盖

当前新增了 source 层单元测试：

```text
internal/source/parse_test.go
internal/source/git_test.go
```

覆盖内容：

- `gitee:owner/repo/path` 简写解析；
- `github:owner/repo/path` 简写解析；
- GitHub tree URL 解析；
- Gitee tree URL 解析；
- 本地路径解析；
- 非法 source 报错；
- provider 到 repo URL 的映射；
- clone 参数生成。

没有在单元测试中执行真实 `git clone`。

原因：

- 真实下载依赖网络；
- 依赖本机 Git 环境；
- 容易造成测试不稳定；
- 更适合作为手动集成测试。

---

## 11. 手动验证命令

在 `cli` 目录下构建：

```bash
go build -o aifn.exe ./cmd/aifn
```

验证 Gitee 简写：

```bash
.\aifn.exe download gitee:aifunc-dev/aifunc-packages/short-summary
```

验证 GitHub tree URL：

```bash
.\aifn.exe download https://github.com/aifunc-dev/aifunc-packages/tree/main/short-summary
```

验证 Gitee tree URL：

```bash
.\aifn.exe download https://gitee.com/aifunc-dev/aifunc-packages/tree/master/short-summary
```

预期输出类似：

```text
Downloading 1 package(s)...

  short-summary        cached

Cached to .aifunc/packages/
Run 'aifn build' to compile and build code.
```

---

## 12. 后续优化方向

### 12.1 隐藏 download 命令

未来可以在 cobra 命令上设置：

```go
Hidden: true
```

这样用户帮助中不展示 `download`，但内部仍可用于调试或测试。

---

### 12.2 Git sparse checkout

当前是浅克隆整个仓库，再复制子目录。

后续可以优化为 sparse checkout：

```bash
git clone --depth 1 --filter=blob:none --sparse <repoURL> <tmp>
git -C <tmp> sparse-checkout set <subpath>
```

适合包仓库变大以后降低下载成本。

---

### 12.3 记录 resolved commit

当前下载没有记录具体 commit。

未来为了可复现安装，可以在下载后记录：

```bash
git -C <tmp> rev-parse HEAD
```

并写入 lockfile：

```json
{
  "name": "short-summary",
  "source": "gitee:aifunc-dev/aifunc-packages/short-summary",
  "provider": "gitee",
  "repo": "aifunc-dev/aifunc-packages",
  "ref": "master",
  "commit": "..."
}
```

---

### 12.4 支持多段 ref

当前 tree URL 的 ref 只支持单段，例如：

```text
master
main
v0.1.0
```

后续可引入显式语法支持多段 branch：

```text
gitee:owner/repo/path#feature/foo
```

或通过远程 API 判断真实 ref。
