# GCli v3.5 更新：自 v3.3.1 以来的改动汇总

> [GCli](https://github.com/gookit/gcli) 是一个用 Go 编写的命令行应用与工具库。
> 本文汇总了从 `v3.3.1` 到最近发布的 `v3.5`（包含 v3.4 周期）的主要改动。这次更新主要集中在开发体验优化和底层健壮性修复上。

如果你在用 Go 写 CLI 工具，这次更新里有几个比较实用的功能。下面直接看重点。

## 核心更新概览

- **Shell 自动补全**：支持零注册生成及动态补全模式（bash / zsh / PowerShell）。
- **命令中间件**：通过 `Command.Use()` / `App.Use()` 处理鉴权、日志等通用逻辑。
- **帮助信息分组**：使用 `Category` 对命令和选项进行分类展示。
- **结构体绑定优化**：新增 `field` 标签规则，支持匿名嵌套结构体自动展开。
- **交互式收集**：缺值时可通过 `Question` 自动提示输入。
- **POSIX 短选项合并**：支持 `-aux` 拆分为 `-a -u -x`。
- **健壮性修复**：panic 处理优化、`help` 命令行为修正等。
- 少量破坏性变更（文末附迁移指南）。

## 1. Shell 自动补全改进

以前生成补全脚本需要硬编码命令和选项，新增命令后得重新生成。从 v3.5 开始，GCli 改进了这一流程。

不需要再手动注册 `genac` 命令，直接使用内置全局选项即可生成静态脚本：

```bash
# 为你的 shell 生成补全脚本，然后 source 它
myapp --gen-completion bash > myapp.bash
source myapp.bash

# zsh / PowerShell 同样支持
myapp --gen-completion zsh  > _myapp
myapp --gen-completion pwsh > myapp.ps1
```

此外，新增了**动态补全**模式。默认生成的脚本不再硬编码名称，而是在补全时通过内置的 `--in-completion` 选项动态调用你的程序获取候选项。以后新增命令，Tab 补全立即生效。

对于选项的取值，也可以通过 `Choices` 设置候选列表：

```go
c.StrOpt2(&format, "format", "output format",
    gflag.WithChoices("json", "yaml", "table"))
// 输入 `--format <Tab>` 会提示: json  yaml  table
```

## 2. 命令中间件

如果需要在执行命令前做鉴权或记录日志，又不想把代码塞进每个命令里，可以使用中间件。

通过 `Use()` 注册的处理器会按顺序在命令主函数 `Func` 之前执行。如果某个处理器返回 error，执行链会中止，错误向上传递。

```go
// 命令级中间件
cmd.Use(func(c *gcli.Command, args []string) error {
    if os.Getenv("TOKEN") == "" {
        return c.NewErrf("missing TOKEN env")
    }
    return nil // 返回 nil 继续执行链
})

// 应用级中间件：在每个命令执行前生效
app.Use(func(c *gcli.Command, args []string) error {
    gcli.Debugf("running command: %s", c.Name)
    return nil
})
```

`Command.Use()` 和 `App.Use()` 都返回接收者本身，支持链式调用。未使用中间件的代码行为与之前完全一致。

## 3. 帮助信息分组

当应用的命令和选项变多时，帮助信息会显得杂乱。现在可以通过设置 `Category` 字段将它们归类到带标题的分组中。

```go
// 命令分组
app.Add(&gcli.Command{Name: "migrate", Desc: "run db migrate", Category: "database"})
app.Add(&gcli.Command{Name: "serve",   Desc: "start http server"}) // 默认分组

// 选项分组
cmd.StrVar(&dsn, &gcli.CliOpt{Name: "db-dsn", Desc: "database dsn", Category: "database"})
cmd.StrOpt2(&port, "port", "bind port", gflag.WithCategory("network"))
```

分组按首次出现的顺序排列，组内条目按名称排序。未设置分类时，输出格式与旧版相同。

## 4. 结构体绑定优化

`FromStruct` 新增了第三种标签规则（`TagRuleField`），并支持自动展开匿名嵌套结构体。

现在支持的三种规则：

- `gcli.TagRuleNamed`（默认）：`flag:"name=int0;shorts=i;required=true;desc=message"`
- `gcli.TagRuleSimple`：`flag:"desc;required;default;shorts"`
- `gcli.TagRuleField`（新增）：使用字段名的 SnakeCase 作为选项名，元数据从独立的 tag 键读取。

```go
type commonOpts struct {
    Verbose bool `flag:"v" desc:"enable verbose output"`
}

type demoOpts struct {
    commonOpts        // 匿名嵌套：展开为 --verbose/-v 选项
    UserName string `flag:"u" desc:"the user name" required:"true"`
    Age      int    `desc:"the user age" default:"18"`
}

c.MustFromStruct(&demoOpts{}, gcli.TagRuleField)
// => 选项: --user-name/-u (必填), --age (默认 18), --verbose/-v
```

`field` 规则的好处是选项名直接取自字段名，`desc`、`default` 等属性各自独立成键，可读性更好。

## 5. 声明式交互收集

如果某个选项必填但用户没传，可以挂载一个 `Question`。GCli 会在检测到空值时交互式提示输入。

```go
c.StrOpt2(&token, "token", "the access token",
    gflag.WithQuestion("Please input your access token: "))
```

运行效果：
```text
$ myapp deploy
Please input your access token: ▮
```

如果同时设置了自定义 `Collector`，则 `Collector` 逻辑优先。

## 6. POSIX 短选项合并

支持了 `-a -u -x` 合并写成 `-aux` 的 POSIX 风格。该功能默认关闭，需要通过 `Config.EnhanceShort` 手动开启。

```go
c.ParserCfg().EnhanceShort = gcli.EnhanceShortMerge  // 1: -aux => -a -u -x
c.ParserCfg().EnhanceShort = gcli.EnhanceShortAttach // 2: 额外支持 -Ostdout => -O stdout
```

也可以全局开启：
```go
gcli.SetEnhanceShort(gcli.EnhanceShortMerge)
```

| 等级 | 常量 | 行为 |
|------|------|------|
| 0 | `EnhanceShortNone` | 关闭（默认），完全兼容旧行为 |
| 1 | `EnhanceShortMerge` | 仅当组合全部为 bool 短选项时才拆分 |
| 2 | `EnhanceShortAttach` | 额外支持取值紧贴写法 `-Ostdout` = `-O stdout` |

这里做了一个安全限制：只有当组合里的每个字符都是 bool 短选项时才会拆分。像 `-aO`（其中 `O` 需要取值）这种混合写法会原样保留，避免误解析。

## 7. 健壮性修复

除了新功能，也修了几个历史遗留问题：

- **panic 不再被吞掉**：`gflag.Parser.Parse` 过去会忽略 recover 的 panic，现在会作为 error 返回，方便上层处理。
- **`help <命令>` 首次调用生效**：修复了以前可能提示 `unknown input command "help"` 的问题。
- **`findSimilarCmd` 逻辑修正**：不再在执行未知命令时往注册表写入幽灵 `help` 条目。
- **`Command.Copy()` 修复**：不再因为共享指针而清空源命令的 hooks。

## 破坏性变更与迁移

有少量清理工作，如果依赖了这些内部实现需要做调整：

| 之前 | 现在 |
|---|---|
| `import ".../gcli/v3/helper"` | 已转为内部包，请内联自己的 helper |
| `import ".../gcli/v3/gclicom"` | 已移除（cliui 迁移后已无用） |
| 全局 `--verbose 4` 选项 | 环境变量 `GCLI_VERBOSE=debug`，或 `gcli.SetVerbose(gcli.VerbDebug)` / `gcli.SetDebugMode()` |

关于移除 `--verbose` 的原因：它绑定的是一份应用内副本，底层的日志器并没有读取它，实际上不起作用，反而会让应用的选项列表变乱。控制日志级别请改用环境变量或代码。

另外，同一进程内的多个 `App` 实例现在会共享全局选项（verbose / help / version / strict / completion）。

## 升级与示例

```bash
go get -u github.com/gookit/gcli/v3@latest
```

仓库的 `_examples/cmd` 目录下新增了几个可运行的示例，包括 `struct-flag`（字段标签+匿名字段）、`short-merge`（短选项合并）和 `ask-demo`（交互收集）。

如果遇到问题或有建议，欢迎在 [GitHub](https://github.com/gookit/gcli) 提 issue 或 PR。完整 API 文档参考 [GoDoc](https://pkg.go.dev/github.com/gookit/gcli/v3)。
