# GCli v3.8 新特性一览 —— 自 v3.5 以来的改进之旅

> [GCli](https://github.com/gookit/gcli) 是一个简单易用、功能丰富的 Go 命令行应用与工具库。
> **v3.6 → v3.8** 这一周期是一次聚焦的现代化升级：更宽容的解析、更丰富的类型绑定、类型安全的泛型 API、
> 完整的三层选项模型，以及内置的文档生成能力。本文将带你逐一了解自 `v3.5` 以来的所有新特性——按能力组织，而非按版本号罗列。

如果你用 Go 写命令行工具，这三个版本补齐了相对 cobra / kong 长期存在的若干差距，同时保留了 GCli「电池全包」的特色。
我们先速览亮点，再用可运行的真实示例逐个展开。

## 亮点速览

- 🔀 **参数乱序自动重排**（v3.6）—— 写在位置参数*之后*的选项依然能被解析；默认开启，对多级命令安全。
- 🧰 **更丰富的结构体标签绑定**（v3.7）—— 原生支持 `[]string`/`[]int`/`[]bool`、`time.Duration`、`map[string]string`，以及 `enum:"a,b,c"` 标签——且绑定器**不再使用 `unsafe`**。
- 🧬 **类型安全的泛型 API**（v3.7）—— `gflag.Opt[T]` / `gflag.BindVar[T]` 取代逐类型的 `BoolVar/IntVar/StrVar/...`。
- 🪜 **三层选项模型**（v3.8）—— `Command.SharedOpts()`（≈ cobra 的 `PersistentFlags`）会被命令及其所有子孙命令继承。
- 📄 **命令文档生成**（v3.8）—— 通过新增的 `docgen` 包或内置 `gendoc` 命令导出 **markdown** 与 **man page**。
- 🏷️ **包重命名 `events` → `gevent`**（v3.6），并新增 `gcli.Evt*` 别名，无需 import 事件包即可引用事件名。
- ⚠️ **少量破坏性变更**，文末附清晰的迁移对照表。

---

## 1. 参数乱序自动重排

GCli 过去严格要求 `--options ... arguments` 这种规范顺序，和标准库 `flag` 一样：遇到第一个非 flag 词元就停止解析选项，
于是 `cmd arg --name tom` 会悄悄丢掉 `--name tom`。可实际使用中，用户经常把顺序写乱。从 v3.6 起，GCli 会在解析**之前**
把输入重排为规范顺序——这样写在位置参数之后的选项也能被正常识别。

```bash
# 现在两者行为一致——选项不再丢失
myapp build --name tom src/
myapp build src/ --name tom
```

该特性**默认开启**，且严格地更宽容：原先能解析的输入解析结果不变，只有原先会失败的顺序现在能成功。重排很谨慎——
已知的取值型选项会带上它的值（`--name tom`），而 bool 选项、`--opt=val`、负数词元（`-5`）、单独的 `-`、以及 `--` 之后的
所有内容都会被正确处理。

关键点：在**多级应用中只重排最终要执行的那个命令的 args**——重排会在遇到子命令名时停止，因此父命令与子命令的选项集
绝不会混淆。

更想要标准库那样的严格顺序？可按解析器关闭：

```go
// 为单个命令关闭
c.ParserCfg().DisableReorderArgs = true
// 或通过 config 函数
gflag.WithReorderArgs(false)
```

可参考可运行的 [`reorder-args`](https://github.com/gookit/gcli/tree/master/_examples/cmd/reorder_demo.go) 示例。

## 2. 更丰富的结构体标签绑定（且不再 `unsafe`）

从结构体绑定选项是 GCli 最好用的特性之一。v3.7 让字段类型丰富了许多——常见的集合类型与时间类型现在都能**原生**绑定，
不再需要声明 `gflag.Strings` / `KVString` 之类的特殊辅助类型：

```go
type deployOpts struct {
    Names []string          `flag:"name=names;shorts=n;desc=name list (repeatable)"`
    Ports []int             `flag:"name=ports;shorts=p;desc=port list (repeatable)"`
    TTL   time.Duration     `flag:"name=ttl;desc=time to live, eg: 1h30m"`
    Meta  map[string]string `flag:"name=meta;shorts=m;desc=key=value metadata (repeatable)"`
    Lang  string            `flag:"name=lang;shorts=l;desc=language;enum=go,php,java"`
}

c.MustFromStruct(&deployOpts{})
```

```bash
myapp deploy -n a -n b -p 80 -p 443 --ttl 1h30m -m k1=v1 -m k2=v2 -l go
# slice 可重复(-n a -n b)，map 可重复(-m k=v)，duration 解析 1h30m，lang 必须是 go/php/java 之一
```

绑定器的新增能力：

- **切片** —— `[]string` / `[]int` / `[]bool` 绑定为可重复选项（`--name a --name b`）。
- **时长** —— `time.Duration` 解析 Go 时长字符串，如 `1h30m`。
- **映射** —— `map[string]string` 绑定为可重复的 `--meta k=v` 选项。
- **枚举** —— 新增的 `enum:"a,b,c"` 标签键会设置选项的候选值（用于补全）**并**加入成员校验，越界的取值会被拒绝。

实现层面，结构体绑定器**不再使用 `unsafe`**——改用安全的 `reflect.Value.Addr().Interface()` 取字段指针，
闭合了结构体绑定里最后一处 unsafe 路径。可试用
[`struct-types`](https://github.com/gookit/gcli/tree/master/_examples/cmd/structtypes_demo.go) 示例。

> **匿名嵌套结构体在任意规则下都会展开。** 通过内嵌结构体复用一组共享选项，对 `named` / `simple` / `field`
> 三种标签规则都生效，并非 `field` 专属；内部字段按当前生效的规则读取。v3.8 还修复了内嵌*未导出*类型
> （类型名小写，如 `commonOpts`）此前被跳过、不展开的问题。

## 3. 类型安全的泛型 API

经典的逐类型绑定器（`BoolVar`、`IntVar`、`StrVar`、`Float64Var` ...）依然可用，但 v3.7 新增了一个泛型入口，
会根据你指针的类型推断出对应的绑定器：

```go
var (
    name string
    age  int
    tags []string
    ttl  time.Duration
)

gflag.Opt(&c.Flags, &name, "name", "n", "tom", "the user name")
gflag.Opt(&c.Flags, &age,  "age",  "a", 18,    "the user age")
gflag.Opt(&c.Flags, &tags, "tag",  "t", nil,   "the tags, repeatable")
gflag.Opt(&c.Flags, &ttl,  "ttl",  "",  time.Duration(0), "time to live")
```

`gflag.Opt[T]` 按指针类型分派到匹配的绑定器，覆盖与结构体绑定相同的类型集——标量、`time.Duration`、切片、
`map[string]string`，以及任意 `flag.Value`。若要完全控制选项元数据，可用 `gflag.BindVar[T]` 配 `*gflag.CliOpt`：

```go
var langs []string
gflag.BindVar(&c.Flags, &langs, gflag.NewOpt("langs", "language list", nil))
```

一次调用，不必记逐类型的方法名，从指针即获得完整类型安全。

## 4. 三层选项模型：共享（继承）选项

这是 v3.8 的重头戏。GCli 一直有**全局**（应用级）选项与**局部**（命令级）选项。v3.8 补上了缺失的中间层——
通过 `Command.SharedOpts()` 提供的**共享选项**，正是 cobra `PersistentFlags` 的对应物。

绑定在 `c.SharedOpts()` 上的选项会被该命令**及其所有子孙命令**继承，并共享同一个绑定变量（同一个底层
`flag.Value` / 指针）。于是一个父命令的选项可以写在任意子命令段并被解析：

```go
var gitDir string

top := &gcli.Command{Name: "git", Desc: "git-like demo"}
// 在父命令上绑定一个共享选项——对每个子命令可见
top.SharedOpts().StrOpt(&gitDir, "git-dir", "", ".git", "the git data dir")

top.Add(&gcli.Command{
    Name: "status",
    Func: func(c *gcli.Command, _ []string) error {
        // 无论 --git-dir 写在这里还是写在父命令上，gitDir 都会被赋值
        gcli.Printf("git dir: %s\n", gitDir)
        return nil
    },
})
```

```bash
myapp git --git-dir /x status      # 写在父命令上
myapp git status --git-dir /x      # 写在子命令上——两者都生效
myapp git status arg --git-dir /x  # 甚至（得益于参数重排）写在参数之后
```

语义与 cobra 高度一致：

- 子命令上**同名的局部选项**优先于被继承的那个。
- `Required` 的共享选项在**最终执行的（叶子）命令**处校验，而非在每个中间祖先命令处误报。
- 在子命令的帮助里，从祖先继承来的选项会归入 **`Inherited Options`** 分组；而命令*自身的*共享选项与其局部选项一起呈现。

底层由一个新的 gflag 原语 `Parser.InheritOptsFrom(src, category...)` 驱动，它按底层 `flag.Value` 把另一个解析器的选项
重新注册进来——所以父子命令写的确实是同一个变量。

## 5. 命令文档生成（markdown + man）

v3.8 带来了文档生成器，大致对应 cobra 的 `cobra/doc`。新增的 `docgen` 包可把单个命令或整个应用渲染为
**markdown** 与 **man page（roff）**：

```go
import "github.com/gookit/gcli/v3/docgen"

// 整个应用的文档树
docgen.MarkdownTree(app, "./docs")  // index.md + 每个命令一个 .md
docgen.ManTree(app, "./man")        // 每个命令一个 .1

// 单个命令
md  := docgen.CmdMarkdown(cmd)
man := docgen.CmdMan(cmd)
```

更想从命令行驱动？加上内置命令再运行即可：

```go
import "github.com/gookit/gcli/v3/builtin"

app.Add(builtin.GenDoc())
```

```bash
./cliapp gendoc -f md  -o ./docs   # markdown
./cliapp gendoc -f man -o ./man    # man page
```

内置了几处贴心处理：

- **Examples 会被清洗并渲染** —— 像 `<cyan>...</>` 这样的颜色标签会被清除，`{$fullCmd}` 等内置变量会被展开，
  文档读起来就是干净可运行的命令。
- **多行 Examples 会被保留** —— man 输出用 `.nf/.fi`（no-fill）区块包裹示例，每行示例都会原样保留，而不会被折叠成一行。
- 应用概览（`index.md`）会包含应用**版本**信息，选项表则通过新增的 `gflag.CliOpt.TypeName()` 访问器带上每个选项的类型。

## 6. 包重命名：`events` → `gevent`，并新增 `gcli.Evt*` 别名

为与其他子包（`gflag`、`gevent`）命名保持一致，事件包从 `github.com/gookit/gcli/v3/events` 重命名为
`github.com/gookit/gcli/v3/gevent`。事件名常量本身不变（`OnAppInitAfter`、`OnCmdRunBefore` ...）。

更进一步：每个事件名现在还以 `gcli.Evt*` 常量的形式暴露，于是你可以**完全不 import 事件包**，直接从 `gcli` 包引用事件名：

```go
// 无需 import 事件包
app.On(gcli.EvtAppInit,      func(ctx *gcli.HookCtx) bool { /* ... */ return false })
app.On(gcli.EvtCmdRunBefore, func(ctx *gcli.HookCtx) bool { /* ... */ return false })
```

## ⚠️ 破坏性变更与迁移

本周期有两处变更可能需要你调整：

| 之前 | 现在 |
|---|---|
| `import "github.com/gookit/gcli/v3/events"` | `import "github.com/gookit/gcli/v3/gevent"`（常量不变） |
| `events.OnCmdRunBefore` | `gevent.OnCmdRunBefore` —— 或 `gcli.EvtCmdRunBefore`（无需 import） |
| `GlobalOpts` 上的应用级解析状态（`ShowHelp`/`ShowVersion`/`inCompletion`/`genCompletion`） | 移到新的应用级 `AppOptions`；改用 `app.AppOpts()` |

关于 `AppOptions` 拆分（v3.8）：描述*单个应用*解析状态的运行时字段，从进程级的 `GlobalOpts` 移到了每个 `App` 各自的
`AppOptions`，于是同进程内多个 `App` 实例不再共享这些字段。**`App.Opts()` 仍返回进程级的 `*GlobalOpts`**
（因此 `app.Opts().Verbose`、`app.Opts() == gcli.GOpts()` 都不变）——只有上面那四个应用级字段搬了家。进程级配置
（verbose / strict / `EnhanceShort` 及日志器）刻意保留在包级单例中，所以日志级别行为不受影响。

> 注意：参数自动重排（v3.6）是一处行为变更，但严格地更宽容——若你依赖旧的「遇到首个非 flag 即停止」行为，
> 可用 `Config.DisableReorderArgs = true` 关闭它。

## 升级方式

```bash
go get -u github.com/gookit/gcli/v3@latest
```

随后可以体验 [`_examples/cmd`](https://github.com/gookit/gcli/tree/master/_examples) 下的可运行示例：
`reorder-args`（参数重排）、`struct-types`（slice/duration/map/enum）、`struct-flag`（字段标签 + 匿名字段）。

## 结语

v3.6 → v3.8 的主题，是补齐重度 CLI 用户期待的那些基础能力：宽容的参数顺序、丰富而安全的结构体绑定、整洁的泛型 API、
跨命令树的持久/共享选项，以及一流的文档生成——同时让 GCli 电池全包的颜色/交互/进度栈保持原样。

欢迎试用！如果遇到问题或有想法，非常欢迎到 [GitHub](https://github.com/gookit/gcli) 提 issue 或 PR。祝你写 CLI 愉快！🎉

---

*相关链接：[GitHub](https://github.com/gookit/gcli) ·
[GoDoc](https://pkg.go.dev/github.com/gookit/gcli/v3) ·
[English version](../gcli-v3.8-whats-new.md)*
