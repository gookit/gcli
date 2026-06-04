# 重构实现计划：A / F / G / H

> 状态：**待评审**（评审通过后再实施）
> 范围：`gcli` 包结构性优化。A/F/G 为**纯重构**(对外 API 不变)；**H 含对外破坏性变更**(需配合次版本号)。
> 建议实施顺序：**F → A → G → H**（H 与前三项正交，可独立进行）

本计划对应的结构性问题：

- **A**：`gOpts`（包级单例）与 `app.opts`（实例级副本）双份割裂
- **F**：`CliOpts.ParseOpts` 与 `Parser.Parse` 解析+校验逻辑重复
- **G**：`findCommandName` 以副作用方式就地修改 `app.args`
- **H**：把内部工具 `helper/` 收敛进 `internal/`，并清理无引用的死代码（`internal/help_tpl.go`、`gclicom/`）

---

## A. 全局选项单一数据源

### 问题（现状）

包内存在两份独立的 `*GlobalOpts`：

| 实例 | 创建处 | 谁在用 |
|---|---|---|
| `gOpts`（包级单例） | `gcli.go:79` | `logf` 读 Verbose（`util.go:37`）；`init()` 读 ENV（`gcli.go:87`）；自由函数 `Verbose/SetVerbose/StrictMode/IsDebugMode/...`（`gcli.go:243-273`）；独立命令路径 `cmd.go:390/419/489`；`app.go:182` 的 `inCompletion` |
| `app.opts`（每个 App 一份） | `app.go:115` `newGlobalOpts()` | `bindingOpts` 绑定 `--help/--version`（`app.go:179`）；`parseAppOpts` 读 `ShowHelp/ShowVersion/NoColor/Verbose/inCompletion`（`app.go:309-329`）；`app.Opts()`（`app.go:603`） |

由此产生的割裂：

1. **Verbose**：`logf` 用 `gOpts`，但 `--verbose`（已移除）曾绑定 `app.opts`，`SetVerbose()` 写 `gOpts` —— 三者不一致。当前靠 `app.go:324` 的 `app.opts.Verbose = gOpts.Verbose` 一行硬同步打补丁。
2. **ShowHelp**：App 路径用 `app.opts.ShowHelp`（经 `app.opts.bindingOpts`），独立命令路径用 `gOpts.ShowHelp`（`cmd.go:390` 的 `gOpts.bindingOpts(&c.Flags)`）—— 同一语义绑到了不同实例。
3. **strictMode / inCompletion** 同样分散在两个实例上。

> 这正是「移除 `--verbose`」时遇到的那个反向同步 bug 的根源。

### 目标

全包仅保留**一份** `*GlobalOpts`，消除同步补丁与「绑哪个实例取决于运行路径」的歧义。**不改任何公共 API 签名**。

### 方案（推荐 A1：App 复用 gOpts）

让 `App` 直接复用包级 `gOpts`，不再单独 `newGlobalOpts()`：

1. `app.go:115` `NewApp`：`app.opts = newGlobalOpts()` → `app.opts = gOpts`。
2. 删除 `app.go:323-324` 的 verbose 同步补丁（同一实例后天然一致）。
3. `bindAppOpts`（`app.go:179`）保持 `app.opts.bindingOpts(fs)` —— 此时即绑定到 `gOpts`，与独立命令路径（`cmd.go:390`）**自动一致**。
4. `app.go:182` 的 `&gOpts.inCompletion` 与 `app.go:329` 的 `app.opts.inCompletion` 现在是同一字段，去掉混用，统一写 `app.opts.inCompletion`。
5. `app.Opts()`、`GOpts()`、`Config()`、`ResetGOpts()`、`SetVerbose()` 等签名与语义全部不变。

### 兼容性 / 风险

- **指针有效性**：`ResetGOpts()` 用的是**就地赋值** `*gOpts = *newGlobalOpts()`（`gcli.go:231`），App 持有的同一指针会随之刷新 —— 安全。
- **行为变化（需在 CHANGELOG 注明）**：同一进程内多个 `App` 实例将**共享**全局选项。但 Verbose/StrictMode 本就经自由函数全局生效，多 App 共存属边缘场景，影响可接受。
- 公共 API 无破坏。

### 测试

- 复用现有 `app_test.go` / `cmd_test.go`（含 standalone、`-h`、`--version`、`by_cmd_ID`）。
- 新增：`SetVerbose(VerbDebug)` 后 `app.Opts().Verbose == gOpts.Verbose`；`app.Opts() == gcli.GOpts()` 为同一指针；`ResetGOpts()` 后 App 侧同步刷新。

### 提交拆分（1 个提交）

- `refactor(gcli): 全局选项统一为单一数据源(App 复用 gOpts)`

---

## F. 合并重复的解析+校验逻辑

### 问题（现状）

两处各实现一遍「`fSet.Parse` + 遍历 opts 校验」：

- `CliOpts.ParseOpts`（`gflag/opts.go:509`）：`fSet.Parse` → 遍历 `co.opts` 调 `opt.Validate(opt.flag.Value.String())`。
- `Parser.Parse`（`gflag/parser.go:197`）：`prepare` → `fSet.Parse` → `AfterParse` → 遍历 `p.opts` 调 `opt.Validate(fSet.Lookup(name).Value.String())`（外加 `recover`）。

两套校验循环已出现分叉，未来易只改一处。`ParseOpts` 目前**仅被测试**直接调用（`opts_test.go:25/78/89`）；生产代码走 `Parse`（`app.go` `fs.Parse`、`cmd.go:500` `c.Parse`）。

### 目标

校验逻辑只留一份，`Parse` 与 `ParseOpts` 共用，行为不变。保留 `ParseOpts` 公共方法（测试与潜在外部调用）。

### 方案

1. 抽私有助手：
   ```go
   // 遍历所有选项做 Validate，集中一处
   func (co *CliOpts) validateAll() error {
       for name, opt := range co.opts {
           if err := opt.Validate(co.fSet.Lookup(name).Value.String()); err != nil {
               return err
           }
       }
       return nil
   }
   ```
2. `ParseOpts` 改为：`fSet.Parse(args)` → `co.validateAll()`。
3. `Parser.Parse` 改为：`prepare` → `fSet.Parse` → `AfterParse` → `p.validateAll()`（删掉内联循环，`recover` 保留 —— 见提交 `d4c8b25`）。

> 注：`opt.flag.Value` 与 `fSet.Lookup(name).Value` 是同一 `flag.Value`，统一取后者即可，语义等价。

### 兼容性 / 风险

- 纯内部重构，签名不变，低风险。

### 测试

- 现有 `gflag/opts_test.go`、`gflag/gflag_test.go`（含 required / validator / collector / panic 回归 `TestParser_Parse_recoverPanic`）已覆盖两条路径。

### 提交拆分（1 个提交）

- `refactor(gflag): 抽取 validateAll，合并 Parse/ParseOpts 重复校验逻辑`

---

## G. findCommandName 去副作用

### 问题（现状）

`findCommandName() (name string, fState FoundState)`（`app.go:401`）在解析过程中**就地修改** `app.args`（`app.go:426/453/455`）并读取（`402`），还顺带写 `app.inputName`（`425`）。调用方 `prepareRun`（`app.go:347`）随后依赖被改过的 `app.args`。职责不单一、隐式耦合，且包内已有**未使用**的 `foundCmd` 结构（`app.go:394-400`）显然是为此预留。

### 目标

把「输入参数 → 命令名 + 剩余参数」做成**输入输出明确**的纯函数式解析，由调用方在一处显式落地到 `app.args`。语义保持完全一致。

### 方案

1. 启用并完善 `foundCmd` 结构作为返回值：
   ```go
   type foundCmd struct {
       state FoundState
       name  string   // 解析出的(顶层)命令名
       raw   string   // 原始输入名，用于提示
       args  []string // 去掉命令名后的剩余参数
   }
   ```
2. 重写为 `func (app *App) findCommandName(args []string) foundCmd`：
   - 全程基于入参 `args` 与局部变量计算，**不写** `app.args` / `app.inputName`。
   - 各分支把结果（state/name/raw/args）填进 `foundCmd` 返回。
3. `prepareRun` 调整为：
   ```go
   fc := app.findCommandName(app.args)
   app.args = fc.args
   app.inputName = fc.raw
   name = fc.name
   // 按 fc.state 走 Founded / NotFound 分支（逻辑同现状）
   ```
   保持 `help` 命令、默认命令、command-ID(`top:sub`)、别名展开等所有现有分支不变。

### 兼容性 / 风险

- **中等风险**：命令解析是核心路径，需逐分支对照保证等价（尤其 `top:sub` 展开与别名→ID 的 `app.args` 拼接，见 `app.go:451-457`）。
- 无公共 API 变化。

### 测试

- 强回归：`TestApp_Run_subcommand`、`TestApp_Run_by_cmd_ID`、`TestApp_AddAliases_and_run`、`TestApp_default_command`、`TestApp_Run_noCommands`。
- 新增针对 `findCommandName` 各分支的纯函数级用例（Founded / 默认命令 / command-ID / 别名 / NotFound）—— 因其将变为无副作用，便于单测。

### 提交拆分（1 个提交）

- `refactor(app): findCommandName 改为无副作用返回 foundCmd`

---

## H. 收敛内部 API 到 internal + 清理死代码

> 范围：**最小 + 一并清 gclicom**（评审结论）。**含对外破坏性变更**，需配合次版本号与 CHANGELOG。

### 前置结论（已核实）

- **`base` 不迁、不动**：`base` 已是包私有类型（`base.go:154` `type base`，小写），外部无法引用——真正对外的只是它提升到 `App`/`Command` 上的方法。把 `base` 抽到独立 `internal` 包会触发 **import 环**：`base` 持有 `map[string]*Command`、返回 `*Command/*App`（引用 10 处），`internal` 需 `import gcli`，而 `gcli` 又要 `import internal` 内嵌 `base`。除非泛型/接口大改否则不可行，且收益≈0。**本项不处理 base。**
- `internal/` 目前仅有 `help_tpl.go`，内容是 3 行空壳（`// var AppHelp` 注释），**无人引用**。
- `gclicom/`（`Output/SetOutput/ResetOutput` + `TextPos/BorderPos`，共 39 行）**全仓零引用**（go 源码、`_examples`、README、docs 均无），疑为 show/progress 移除后的遗留。

### H1. `helper/` → `internal/helper/`

`helper` 导出 `IsGoodName/IsGoodCmdId/IsGoodCmdName/Panicf/RenderText` + 正则常量，纯内部工具。模块内 7 个调用方：
`help.go`、`app.go`、`cmd.go`、`util.go`、`builtin/gen_auto_complete.go`、`gflag/args.go`、`gflag/opts.go`。

做法（一次原子提交）：

1. `git mv helper internal/helper`（包名保持 `package helper` 不变）。
2. 7 个文件的 import 路径：`gcli/v3/helper` → `gcli/v3/internal/helper`；调用处 `helper.X` **无需改**（包名不变）。
3. `internal/helper` 对模块内全部可见（internal 规则允许模块根下任意包导入），`gflag`/`builtin` 等照常引用。

> **破坏性**：外部若直接 `import gcli/v3/helper` 将失效。这些是内部校验/渲染工具，不属常规对外用法，可接受。

### H2. 删除死代码 `internal/help_tpl.go`

3 行空壳、无引用 → 直接 `git rm internal/help_tpl.go`（H1 后 `internal/` 由 `internal/helper/` 承载，不留空包）。

### H3. 清理未引用的 `gclicom/`

全仓零引用。两种处理（**推荐删除**，移动死代码意义不大）：

- **推荐**：`git rm -r gclicom`（破坏性移除）。
- 备选：若想保留这些类型供日后用，`git mv gclicom internal/gclicom`（同样对外破坏，但至少不再属公共 API）。

> 删除/迁移均为对外破坏性变更，需在 CHANGELOG 注明。

### 兼容性 / 风险

- **无行为变化**，纯包归属调整；构建在 import 全部改完前会断，故每步**单提交内保持可编译**。
- 破坏性仅限外部直接 import `helper/`、`gclicom/` 的用户 → **建议随次版本号发布**并在 CHANGELOG/Release Notes 列明。

### 测试

- 无需新增用例；以 `go build ./...`、`go test ./. ./gflag`、`go build ./_examples/cliapp` 全绿为准。

### 提交拆分（2 个提交）

- `refactor(helper): 迁移 helper 至 internal/helper，删除空壳 internal/help_tpl.go`
- `chore: 移除未引用的 gclicom 包`

---

## 总览与排期

| 项 | 风险 | 影响文件 | 公共 API | 预估改动 |
|---|---|---|---|---|
| F | 低 | `gflag/opts.go`、`gflag/parser.go` | 不变 | ~30 行 |
| A | 中 | `app.go`、（`gcli.go` 核对）| 不变 | ~15 行 + 测试 |
| G | 中 | `app.go` | 不变 | ~60 行 + 测试 |
| H | 低(机械)，但**破坏性** | `helper/`→`internal/helper/`、7 处 import、删 `gclicom/`、`internal/help_tpl.go` | **破坏**(移除 `helper`/`gclicom` 公共包) | 机械移动 |

**建议顺序 F → A → G → H**，每项独立提交、独立跑全量 `go test . ./gflag` 后再进入下一项。
H 与 A/F/G 正交，可单独/最后做；因含破坏性，**建议随次版本号一起发布**。

### 实施清单（实施时勾选）

- [x] F：抽 `validateAll`，`Parse`/`ParseOpts` 复用；全量测试通过
- [ ] A：`App` 复用 `gOpts`，删同步补丁，统一 `inCompletion`；补单测；全量测试通过
- [ ] G：`findCommandName` 改 `foundCmd` 无副作用；补分支单测；全量回归通过
- [ ] H：`helper`→`internal/helper`、删 `internal/help_tpl.go` 与 `gclicom/`；`go build ./... && go test` 全绿
- [ ] 在 CHANGELOG/README 注明：A 的「多 App 共享全局选项」、H 的「移除 helper/gclicom 公共包」
