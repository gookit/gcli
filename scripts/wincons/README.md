# wincons

在 **真实 Windows 控制台**上驱动示例程序（`_examples/cliapp`），是 `scripts/ptysmoke`（PTY，仅
Linux/macOS）的 Windows 对应物。Python 的 `pty` 模块在 Windows 上不存在；ConPTY 是一条可走的路
（`gookit/gopty` 在本机 gofer job 里实测能起 `cmd.exe`，而 suga 的 wincons 在同类 job 里报过
`0xC0000142`——取决于 job 设置），但**控制台**才给出忠实的 cooked 行输入与 Ctrl-C 语义，
所以这里用 Windows 自己的控制台：

1. `FreeConsole` 脱离启动它的控制台（gofer job 里没有可保留的），`AllocConsole` 另建一个私有控制台；
2. 应用以该控制台的 `CONIN$` 作为 stdin 启动（**cooked 行输入是真的**，提示符也是真的），
   stdout/stderr 走管道 —— 这样 harness 仍能读到它打印了什么；
3. 等提示符在输出里出现后，用 `WriteConsoleInput` 写 **KEY_EVENT** 记录（等价真人敲键 + 回车）；
4. 每个用例都新建一个控制台，避免上一次残留的按键事件串到下一次。

```powershell
go build -o tmp/wincons/cliapp.exe ./_examples/cliapp
go run ./scripts/wincons --bin tmp/wincons/cliapp.exe          # --verbose 打印每个用例的完整输出与控制台状态
```

该文件带 `//go:build windows`：在 Linux/macOS 上 `go build ./...` / `go test ./...` 会跳过它（不报错），
但**指名**构建这个包会报 `build constraints exclude all Go files`（平台专用工具的正常行为，同 suga）。

运行时会为每个用例弹出（并很快关闭）一个控制台窗口；输出与结果打印在本进程的 stdout（管道），
退出码 0 表示全部通过。

| 用例 | 检查 |
|---|---|
| `ask-demo`（局部选项 `CliOpt.Question`） | 空值 → 提示符出现一次、键入的内容成为选项值、退出码 0 |
| `ask-demo --token=cli-token` | 不提问、值被采用 |
| `shared-ask run`（共享选项 + `WithQuestion`） | 只在**实际执行的命令**处提问一次，键入值被叶子命令读到 |
| `shared-ask run --token=leaf-token` | 不提问、叶子用给定值 |
| `shared-ask run -h` | help 列出继承选项且从不提问 |

实测（本机 Windows 10 IoT LTSC 2024，Go 1.25.10）：15/15 通过；把 `doExecute` 里的
`prepareSharedOpts()` 调用去掉后重建，`shared-ask run` 两条检查立即 FAIL、退出码 1（与 PTY 侧
`shared_ask.py` 的反向验证一致）。

## 踩坑记录

- **虚拟键码必须真实**：`KEY_EVENT` 里的 `wVirtualKeyCode` 用字符本身是错的 —— `-` 的 0x2D 恰好是
  `VK_INSERT`，控制台会把它当功能键（插入模式切换）而丢掉字符，`win-token` 变成空/`wintoken`。
  改用 `VkKeyScanW` 取真实的 VK + 移位状态（高位字节）。
- 回车要显式发 `VK_RETURN`（`VkKeyScan` 对 CR 拿不到）。
- 子进程的 stdout 是管道 → 颜色/高亮默认关闭（`NO_COLOR=1`，断言看文本）；要测彩色效果就把它接到
  `CONOUT$`，但那样 harness 读不到输出（suga 的 `--console-output` 是同一权衡）。
- 控制台窗口会真出现（AllocConsole），CI/无人值守时注意；本 harness 每次用例后 `FreeConsole` 释放。

## 与 ptysmoke 的分工

同一批场景两边都跑了 `ask-demo` / `shared-ask`：PTY 侧（Linux/macOS）还覆盖 cliui 的行式提示
（`interact_confirm.py`），Windows 侧覆盖的是控制台路径。真正跨平台不变的语义（收集只在执行命令处、
已给值不提问）用 `go test`（`shared_opts_test.go`）固化，这两个 harness 负责"在真终端里确实如此"。
