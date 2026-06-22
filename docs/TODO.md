# TODO

## 待修复（重构期间发现的既有问题）

- [x] **`help COMMAND` 失效**：`help` 未注册为命令，`findCommandName` 返回 `NotFound`，
  首次输出「unknown input command help」。已在 `findCommandName` 将 `help` 识别为 `Founded`。
- [x] **`findSimilarCmd` 污染命令注册表**：`CmdNameMap()` 返回真实 `cmdNames` map，
  `names["help"]=4` 直接写进注册表（导致上面的测试靠一次 warmup 污染才"碰巧"通过）。已改为复制 map 后再加 `help`。
- [x] **`TestApp_showCommandHelp` 断言失效**：改为每场景全新 app + `StrContains`，
  并验证「去掉 help 修复则测试 FAIL」，确保断言真正生效；新增防污染回归测试。
- [x] 预存 gofmt 漂移：`gflag/gflag.go`、`gflag/util.go`、`builtin/tcpproxy/tcp_proxy.go`。

## 结构性重构（待评审 → 待实施）

- [x] A / F / G / H 四项重构均已完成，详见 [plans/refactor-A-F-G.md](plans/refactor-A-F-G.md)
  - [x] F：合并 `Parse`/`ParseOpts` 重复校验
  - [x] A：全局选项单一数据源
  - [x] G：`findCommandName` 去副作用
  - [x] H：`helper/`→`internal/helper/`、清理 `gclicom/` 与死代码
- [x] 发版前：CHANGELOG/Release Notes（见 [CHANGELOG.md](../CHANGELOG.md) v3.4.0），并将 `version` 常量更新为 3.4.0

## 功能规划（TODO 扫描后立项，待实施）

- [x] A：清理本轮已实现功能的陈旧 TODO 注释（命令/选项 Category）
- [x] **B1** 内置补全选项 — [plans/feat-B1-completion.md](plans/feat-B1-completion.md)
  - [x] 静态 `--gen-completion`（commit 3a76280）
  - [x] 动态 `--in-completion` 候选计算 + showAutoCompletion（commit 88d70e9）
  - [x] 脚本委托动态(默认瘦脚本) + 静默模式 + genac --static opt-in（子阶段 3）
  - [x] PowerShell(pwsh) 动态补全（`--gen-completion pwsh`，`.ps1`；静态 pwsh 不支持）
  - [x] 选项值候选 `CliOpt.Choices`（子阶段 4，commit 75bc031）
  - [x] 收尾(按需)：bash/zsh/pwsh 实际 Tab 交互的人工验证；隐藏命令过滤
- [x] **B2** 命令中间件 — [plans/feat-B2-middleware.md](plans/feat-B2-middleware.md)
  - [x] 命令级 `Command.Use()`（commit 7cec911）
  - [x] 修 doExecute 内联 recover(no-op) 为 defer，真正捕获 panic（commit 467d883）
  - [x] App 级 `App.Use()`（commit 66502f7）
- [x] **B4+B5** POSIX 短选项增强 — [plans/feat-B4-B5-posix-options.md](plans/feat-B4-B5-posix-options.md)
- [x] **B6+B7** 结构体标签增强 + 声明式交互收集 — [plans/feat-B6-B7-struct-tag-question.md](plans/feat-B6-B7-struct-tag-question.md)

### 其余 TODO（低优先，按需）

- [x] **测试隔离债**：已改为每用例独立工厂函数 + 补 gOpts 重置，`-shuffle` 40 次 0 失败（commit 2a6af97）
- B3 交互式 shell `inShell`；C 类边角：`parser.go:223` ParseArgs 未接、`help.go:153` 多级子命令帮助、
  `cmd.go:510` `prepare()` 空桩、`ext.go:165` `HookCtx.err` 未用、`builtin/gen_emoji_codeMap.go:59` 打印 "TODO"。

## 对标改进（对比主流 Go CLI 库后立项）

> 背景与依据见 [compare-with-others.zh-CN.md](compare-with-others.zh-CN.md) 的「gcli 的差距」。

- [ ] 支持 man page / markdown 命令文档生成
- [x] **D1 结构体绑定：去 unsafe + 类型丰富度**（增量、低风险）— [plans/feat-D1-struct-binding.md](plans/feat-D1-struct-binding.md)
  - [x] `gflag/parser.go` `fromStructValue` 基础类型分支去掉 `unsafe.Pointer`/`UnsafeAddr()`，
    改用安全的 `fv.Addr().Interface().(*T)`（commit 7830e0e）
  - [x] 扩充标签自动绑定的类型：原生 `[]string/[]int/[]bool` → `cflag.Strings/Ints/Booleans`，
    `time.Duration` → `DurationVar`，`map[string]string` → `mapStrValue` 适配器（commit a38d5d5/633e0ff）
  - [x] 标签打通 `enum`：`enum:"a,b,c"` → `CliOpt.Choices` + 成员校验（commit cd33f51）
  - [x] 泛型 API：`gflag.Opt[T]/BindVar[T]`，类型安全、可扩展，老 API 保留（commit 99ffdfa）
  - 对标：kong / go-arg / go-flags 原生支持 slice/map/enum/Duration（go-arg 连 Duration 都内置）

- [ ] **D2 持久化选项继承模型**（伤筋动骨、建议独立里程碑）
  - [ ] 新增「持久化选项」中间层：`Command.PersistentFlags()` 或 `CliOpt.Persistent` 标记；
    dispatch 到子命令时把祖先链持久选项 merge 进叶子 flag set（与 args 重排契合：可写在叶子段任意位置）
  - [ ] 进程级单例 `gOpts` 改为 per-App 实例，解决多 App/并发共享（CHANGELOG v3.4.0 已记此坑）
  - [ ] 明确「全局(App) / 持久(命令子树) / 局部(命令)」三层语义 + 文档 + 专项测试（持久选项 × 多级）
  - 对标：cobra 的 `PersistentFlags` / `LocalFlags` / `InheritedFlags` 三层模型
