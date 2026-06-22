# 功能实现计划：D3 命令文档生成（man / markdown）

> 状态：**待实施**
> 范围：新增 `docgen` 包 + builtin `GenDoc` 命令 + gflag 小访问器；约 4-6 个文件。
> 依据：[../compare-with-others.zh-CN.md](../compare-with-others.zh-CN.md) 差距「文档生成」；[../TODO.md](../TODO.md) 对标改进。
> 对标：cobra 的 `cobra/doc`（`GenMarkdownTree` / `GenManTree`）。

## 现状

- gcli 能渲染交互式 help（`Command.ShowHelp` / `gflag.BuildOptsHelp`），但**无法导出**命令文档
  （markdown / man page）。cobra、go-flags 有，是对比文档里列出的差距项之一。
- 可用内省面（已核对）：
  - `Command`：`Name`/`Desc`/`Help`/`Examples`/`Aliases`/`Path()`/`ID()`/`Commands()`(子命令)/`App()`/`Root()`。
  - 选项：`c.Opts()` → `map[string]*gflag.CliOpt`；`CliOpt` 有 `Name/Desc/Shorts/DefVal/Required/Hidden/Choices/Category`。
    类型名 `flagType` 是私有 → 需补公开访问器 `CliOpt.TypeName()`。
  - 参数：`c.Args()` → `[]*gflag.CliArg`；`CliArg` 有 `Name/Desc/Required/Arrayed/HelpName()`。
  - builtin 命令范式：`GenAutoComplete(fns...) *gcli.Command`，`doGen` 内用 `c.App()` 内省。

## 目标

1. 新增 `docgen` 包：把单个命令 / 整个 App 渲染为 **markdown** 与 **man page(roff)**。
2. 新增 builtin `GenDoc()` 命令：app 添加后即可 `./cliapp gendoc -f md -o ./docs` 导出。
3. 不改变现有 help/解析行为；纯新增。

## 方案

### D3.0 gflag 小访问器（前置）
- `gflag/opts.go` 增 `func (m *CliOpt) TypeName() string { return m.flagType }`（公开类型名，供文档渲染）。

### D3.1 docgen 包：markdown
新包 `github.com/gookit/gcli/v3/docgen`（import gcli，无环）：
- `func CmdMarkdown(c *gcli.Command) string`：单命令 → markdown。结构(cobra 风格)：
  - `# <full path>`、Desc、`## Synopsis`(若有 Help)、用法代码块、
  - `## Options` 表(`Option | Type | Default | Required | Description`，跳过 Hidden)、
  - `## Arguments` 表(`Argument | Required | Description`)、
  - `## Examples`(若有)、`## SubCommands`(列子命令 + 相对链接)、底部生成标记。
- `func AppMarkdown(app *gcli.App) string`：App 概览(命令列表 + 链接)。
- `func MarkdownTree(app *gcli.App, dir string) error`：每命令一个 `.md`(含子命令递归)，
  文件名 `path` 以 `_` 连接(如 `cliapp_remote_add.md`)，外加一个 `index.md`(= AppMarkdown)。

### D3.2 docgen 包：man page
- `func CmdMan(c *gcli.Command) string`：roff 格式。最少包含 `.TH`、`.SH NAME`、`.SH SYNOPSIS`、
  `.SH DESCRIPTION`、`.SH OPTIONS`(+ ARGUMENTS/EXAMPLES 若有)。注意转义 `-` → `\-`、`.`/`\` 行首处理。
- `func ManTree(app *gcli.App, dir string) error`：每命令一个 `.1` 文件。

### D3.3 builtin GenDoc 命令
- `func GenDoc(fns ...func(c *gcli.Command)) *gcli.Command`(参照 GenAutoComplete)：
  - name `gendoc`，options：`--format/-f`(md|man，默认 md)、`--output/-o`(目录，默认 `./docs`)。
  - `Func` 内：按 format 调 `docgen.MarkdownTree(c.App(), dir)` / `ManTree(...)`，打印结果路径。
- 在 `_examples/cliapp/main.go` 注册（可选演示）。

## 风险 / 兼容
- 纯新增包/命令，零破坏。man 的 roff 转义是主要细节点，用测试固化。
- `TypeName()` 仅暴露已有私有字段，安全。

## 测试（`github.com/gookit/goutil/x/assert`）
- 构造一个含「选项+参数+示例+子命令」的小 app：
  - markdown：断言包含命令全名、Desc、`## Options`、某选项行(`-n, --name`)、`## Arguments`、`## Examples`、子命令链接。
  - man：断言包含 `.TH`、`.SH NAME`、`.SH OPTIONS`，且 `--name` 被正确转义。
  - tree：写入临时目录(项目 tmp 或 t.TempDir())，断言生成了预期文件名集合。
- gendoc 命令：run 后目录下生成文件、退出码 0。

## 提交拆分（按 R002）
1. `feat(gflag): 增加 CliOpt.TypeName() 类型名访问器`
2. `feat(docgen): markdown 命令文档生成(CmdMarkdown/AppMarkdown/MarkdownTree)`
3. `feat(docgen): man page 命令文档生成(CmdMan/ManTree)`
4. `feat(builtin): GenDoc 命令导出 md/man 文档`
5. `docs: CHANGELOG/README/TODO 记录文档生成特性`
