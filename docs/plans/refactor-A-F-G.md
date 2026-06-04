# 重构实现计划：A / F / G

> 状态：**待评审**（评审通过后再实施）
> 范围：`gcli` 包结构性优化，均为**纯重构**，对外行为/公共 API 尽量保持不变
> 建议实施顺序：**F → A → G**（按风险从低到高）

本计划对应前期分析中的三个结构性问题：

- **A**：`gOpts`（包级单例）与 `app.opts`（实例级副本）双份割裂
- **F**：`CliOpts.ParseOpts` 与 `Parser.Parse` 解析+校验逻辑重复
- **G**：`findCommandName` 以副作用方式就地修改 `app.args`

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

## 总览与排期

| 项 | 风险 | 影响文件 | 公共 API | 预估改动 |
|---|---|---|---|---|
| F | 低 | `gflag/opts.go`、`gflag/parser.go` | 不变 | ~30 行 |
| A | 中 | `app.go`、（`gcli.go` 核对）| 不变 | ~15 行 + 测试 |
| G | 中 | `app.go` | 不变 | ~60 行 + 测试 |

**建议顺序 F → A → G**，每项独立提交、独立跑全量 `go test . ./gflag` 后再进入下一项。

### 实施清单（实施时勾选）

- [ ] F：抽 `validateAll`，`Parse`/`ParseOpts` 复用；全量测试通过
- [ ] A：`App` 复用 `gOpts`，删同步补丁，统一 `inCompletion`；补单测；全量测试通过
- [ ] G：`findCommandName` 改 `foundCmd` 无副作用；补分支单测；全量回归通过
- [ ] 在 CHANGELOG/README 注明 A 的「多 App 共享全局选项」行为说明
