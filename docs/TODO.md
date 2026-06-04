# TODO

## 待修复（重构期间发现的既有问题，本轮未动）

- [ ] **`help COMMAND` 失效**：`help` 未注册为命令，`findCommandName` 对它返回 `NotFound`，
  实际输出「unknown input command help」而非命令帮助。`prepareRun` 中 `name == HelpCommand`
  的分支因此为死分支。修法：在 `findCommandName` 里将 `help` 识别为 `Founded`（或注册内置 help 命令）。
- [ ] **`TestApp_showCommandHelp` 断言失效**：输出明显不含 "Name: test" 仍 PASS，
  说明被捕获的 buffer 与实际输出流不一致，断言形同虚设。需修正测试使其真正校验。
- [ ] 预存 gofmt 漂移：`gflag/gflag.go`、`gflag/util.go`、`builtin/tcpproxy/tcp_proxy.go`。

## 结构性重构（待评审 → 待实施）

- [ ] A / F / G / H 四项重构，详见 [plans/refactor-A-F-G.md](plans/refactor-A-F-G.md)
  - F：合并 `Parse`/`ParseOpts` 重复校验（低风险）
  - A：全局选项单一数据源（中风险）
  - G：`findCommandName` 去副作用（中风险）
  - H：`helper/`→`internal/helper/`、清理 `gclicom/` 与死代码（机械，但对外破坏性）
