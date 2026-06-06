# GCli v3.4 新特性一览 —— 自 v3.3.1 以来的改进之旅

> [GCli](https://github.com/gookit/gcli) 是一个简单易用、功能丰富的 Go 命令行
> 应用与工具库。**v3.4** 版本带来了一批聚焦**开发体验**与**健壮性**的更新。
> 本文将带你逐一了解自 `v3.3.1` 以来的所有新特性——从日常最常用的，到进阶的能力。

如果你用 Go 写命令行工具，这次更新值得一看。我们先速览亮点，再用可运行的真实示例
逐个展开。

## 亮点速览

- 🧠 **更聪明的 Shell 自动补全** —— 零注册即可生成，并新增**动态补全**模式，
  脚本零维护（支持 bash / zsh / PowerShell）。
- 🧅 **命令中间件** —— `Command.Use()` / `App.Use()`，用于鉴权、日志、计时等横切逻辑。
- 🗂️ **帮助信息分组** —— 通过 `Category` 把命令与选项归入带标题的分组。
- 🏷️ **更灵活的结构体绑定** —— 新增 `field` 标签规则，并支持**匿名嵌套结构体**自动展开。
- 💬 **声明式交互收集** —— 用一个 `Question` 即可在缺值时交互收集输入。
- ➖ **POSIX 短选项合并** —— `-aux` = `-a -u -x`，按需开启且安全。
- 🛡️ **健壮性修复** —— panic 不再被吞掉、`help <命令>` 首次调用即生效等。
- ⚠️ **少量破坏性变更**，文末附清晰的迁移对照表。

---

## 1. 更聪明的 Shell 自动补全

过去做补全意味着生成一份**静态**脚本，把命令名和选项名都硬编码进去。可一旦你新增了
命令，脚本立刻过时，必须重新生成。v3.4 从两端解决了这个痛点。

**零注册的静态生成。** 你不再需要注册 `genac` 命令。现在每个应用都内置了一个全局选项：

```bash
# 为你的 shell 生成补全脚本，然后 source 它
myapp --gen-completion bash > myapp.bash
source myapp.bash
# zsh / PowerShell 同样支持
myapp --gen-completion zsh  > _myapp
myapp --gen-completion pwsh > myapp.ps1
```

**动态补全（零维护）。** 默认生成的是一份*瘦*脚本：它不再硬编码名称，而是在补全时
通过内置的 `--in-completion` 选项**回调你的二进制**来获取候选项。明天你新增一个命令，
Tab 补全立即生效——无需重新生成。

**选项值候选。** 想让补全也提示某个选项的*取值*？给它一个 `Choices` 列表即可：

```go
c.StrOpt2(&format, "format", "output format",
    gflag.WithChoices("json", "yaml", "table"))
// 输入 `--format <Tab>` 现在会提示: json  yaml  table
```

候选计算逻辑有完整的单元测试覆盖；bash/zsh/pwsh 的 shell 胶水层都委托到这同一个动态
入口，因此行为保持一致。

## 2. 命令中间件

想在命令主逻辑之前做鉴权、日志或计时——又不想把这些代码复制到每个命令里？中间件来了。

用 `Use()` 注册一个或多个处理器，它们会**按注册顺序、在**命令主函数 `Func` **之前**执行。
任一处理器返回 error，链条即中止、错误向上传递（主函数 `Func` 被跳过）。

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

`Command.Use()` 和 `App.Use()` 都返回接收者本身，方便链式调用。未使用中间件的应用，
行为与之前完全一致。

## 3. 帮助信息分组

应用变大后，扁平的命令/选项长列表会越来越难找。现在你可以用 `Category` 把它们归入
带标题的分组。

```go
// 命令分组
app.Add(&gcli.Command{Name: "migrate", Desc: "run db migrate", Category: "database"})
app.Add(&gcli.Command{Name: "serve",   Desc: "start http server"}) // 默认分组

// 选项分组
cmd.StrVar(&dsn, &gcli.CliOpt{Name: "db-dsn", Desc: "database dsn", Category: "database"})
cmd.StrOpt2(&port, "port", "bind port", gflag.WithCategory("network"))
```

分组按首次出现的顺序排列，组内条目按名称排序。当未设置分类时，输出与之前完全相同——
完全向后兼容。

## 4. 更灵活的结构体绑定

GCli 一直支持直接从结构体绑定选项。v3.4 新增了第三种标签规则与匿名字段支持。

`FromStruct` 现在支持三种规则，通过 `c.FromStruct(ptr, ruleType)` 选择：

- `gcli.TagRuleNamed`（默认）：`flag:"name=int0;shorts=i;required=true;desc=message"`
- `gcli.TagRuleSimple`：`flag:"desc;required;default;shorts"`
- `gcli.TagRuleField`**（新增）**：用**字段名**(SnakeCase) 做选项名，元数据从独立的
  tag 键读取。并且**自动展开匿名嵌套结构体**——非常适合复用一组通用选项。

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

`field` 规则最为简洁：选项名来自字段名，而 `desc` / `default` / `required` 各自独立成键
——易读也易维护。

## 5. 声明式交互收集

有时某个值是必需的，但用户忘了传。与其手写一个收集器，不如直接挂一个 `Question`：
当选项值为空时，GCli 会交互式地提示输入（内置的默认收集器）。

```go
c.StrOpt2(&token, "token", "the access token",
    gflag.WithQuestion("Please input your access token: "))
```

```text
$ myapp deploy
Please input your access token: ▮
```

如果你同时设置了自定义 `Collector`，则 `Collector` 优先于 `Question`。

## 6. POSIX 短选项合并

经典的 POSIX 工具允许合并短选项：`-a -u -x` 可以写成 `-aux`。GCli 现在也支持了——
而且很谨慎、**按需开启**，默认不改变任何行为。

通过 `Config.EnhanceShort` 开启，并使用新增的自解释常量：

```go
c.ParserCfg().EnhanceShort = gcli.EnhanceShortMerge  // 1: -aux => -a -u -x
c.ParserCfg().EnhanceShort = gcli.EnhanceShortAttach // 2: 额外支持 -Ostdout => -O stdout
```

| 等级 | 常量 | 行为 |
|------|------|------|
| 0 | `EnhanceShortNone` | 关闭（默认）—— 完全兼容 |
| 1 | `EnhanceShortMerge` | **仅当组合全部为 bool 短选项时**才拆分 |
| 2 | `EnhanceShortAttach` | 额外支持取值紧贴写法 `-Ostdout` = `-O stdout` |

关键的安全规则：只有当组合里**每个字符都是 bool 短选项**时才拆分。像 `-aO`
（其中 `O` 需要取值）这样的混合写法会原样保留，绝不会误伤取值型短选项。
（严格模式现在内部也走这条安全路径，取代了旧的"盲拆"逻辑。）

## 7. 健壮性修复

一批长期存在的小瑕疵被一并打磨：

- **panic 不再被吞掉。** `gflag.Parser.Parse` 过去会打印并忽略被 recover 的 panic；
  现在会作为 error 返回，让你的代码能够响应。
- **`help <命令>` 首次调用即生效。** 此前它可能打印 `unknown input command "help"`。
- **`findSimilarCmd` 不再污染命令注册表**——以前在执行未知命令时会写入一个幽灵 `help` 条目。
- **`Command.Copy()` 不再清空源命令的 hooks**（旧实现共享指针，会重置原命令）。

## ⚠️ 破坏性变更与迁移

少量清理工作，如果你依赖了它们则需要调整：

| 之前 | 现在 |
|---|---|
| `import ".../gcli/v3/helper"` | 已转为内部包 —— 请内联自己的 helper |
| `import ".../gcli/v3/gclicom"` | 已移除（cliui 迁移后已无用） |
| 全局 `--verbose 4` 选项 | 环境变量 `GCLI_VERBOSE=debug`，或 `gcli.SetVerbose(gcli.VerbDebug)` / `gcli.SetDebugMode()` |

为什么移除 `--verbose`？它绑定的是一份应用内副本，而日志器从未读取它，因此毫无实际效果
——只会让你应用的选项列表变得杂乱。请改用环境变量或代码来控制日志级别。

> 注意：作为"全局选项统一为单一数据源"的一部分，同一进程内的多个 `App` 实例现在会共享
> 全局选项（verbose / help / version / strict / completion）。

## 升级方式

```bash
go get -u github.com/gookit/gcli/v3@latest
```

随后可以体验 [`_examples/cmd`](https://github.com/gookit/gcli/tree/master/_examples)
下的可运行示例：`struct-flag`（字段标签 + 匿名字段）、`short-merge`（EnhanceShort）、
`ask-demo`（Question）。

## 结语

v3.4 的主题，就是让用 GCli 构建命令行应用变得更舒服——会自我维护的补全、承接那些
"枯燥但必要"的横切逻辑的中间件、更整洁的帮助输出、更灵活的选项绑定——同时在底层悄悄
修复了若干健壮性问题。

欢迎试用！如果遇到问题或有想法，非常欢迎到 [GitHub](https://github.com/gookit/gcli)
提 issue 或 PR。祝你写 CLI 愉快！🎉

---

*相关链接：[GitHub](https://github.com/gookit/gcli) ·
[GoDoc](https://pkg.go.dev/github.com/gookit/gcli/v3) ·
[English version](../gcli-v3.4-whats-new.md)*
