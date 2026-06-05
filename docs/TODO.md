# TODO

## 待修复（重构期间发现的既有问题）

- [x] **`help COMMAND` 失效**：`help` 未注册为命令，`findCommandName` 返回 `NotFound`，
  首次输出「unknown input command help」。已在 `findCommandName` 将 `help` 识别为 `Founded`。
- [x] **`findSimilarCmd` 污染命令注册表**：`CmdNameMap()` 返回真实 `cmdNames` map，
  `names["help"]=4` 直接写进注册表（导致上面的测试靠一次 warmup 污染才"碰巧"通过）。
  已改为复制 map 后再加 `help`。
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

> 扫描明细见 `tmp/todo-scan.md`（工作产物，未入库）。

- [x] A：清理本轮已实现功能的陈旧 TODO 注释（命令/选项 Category）
- [ ] **B1** 内置补全选项 — [plans/feat-B1-completion.md](plans/feat-B1-completion.md)
  - [x] 静态 `--gen-completion`（commit 3a76280）
  - [x] 动态 `--in-completion` 候选计算 + showAutoCompletion（commit 88d70e9）
  - [x] 脚本委托动态(默认瘦脚本) + 静默模式 + genac --static opt-in（子阶段 3）
  - [x] PowerShell(pwsh) 动态补全（`--gen-completion pwsh`，`.ps1`；静态 pwsh 不支持）
  - [ ] 选项值候选 `CliOpt.Choices`（子阶段 4，可选）
  - [ ] bash/zsh/pwsh 实际 Tab 交互的人工验证（shell 胶水层无法自动测）
- [ ] **B2** 命令中间件 — [plans/feat-B2-middleware.md](plans/feat-B2-middleware.md)
- [ ] **B4+B5** POSIX 短选项增强 — [plans/feat-B4-B5-posix-options.md](plans/feat-B4-B5-posix-options.md)
- [ ] **B6+B7** 结构体标签增强 + 声明式交互收集 — [plans/feat-B6-B7-struct-tag-question.md](plans/feat-B6-B7-struct-tag-question.md)

### 其余 TODO（低优先，按需）

- B3 交互式 shell `inShell`；C 类边角：`parser.go:223` ParseArgs 未接、`help.go:153` 多级子命令帮助、
  `cmd.go:510` `prepare()` 空桩、`ext.go:165` `HookCtx.err` 未用、`builtin/gen_emoji_codeMap.go:59` 打印 "TODO"。
