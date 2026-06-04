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

- [ ] A / F / G / H 四项重构，详见 [plans/refactor-A-F-G.md](plans/refactor-A-F-G.md)
  - F：合并 `Parse`/`ParseOpts` 重复校验（低风险）
  - A：全局选项单一数据源（中风险）
  - G：`findCommandName` 去副作用（中风险）
  - H：`helper/`→`internal/helper/`、清理 `gclicom/` 与死代码（机械，但对外破坏性）
