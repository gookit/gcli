# ptysmoke

在伪终端（PTY）里驱动真实的示例程序（`_examples/cliapp`）做端到端冒烟，覆盖 `go test`
测不到的终端行为：交互提问真的打印出来、用户输入被读成选项值、已给值时不再提问。

思路与驱动库来自 suga 的 `scripts/ptysmoke`（那份 README 记录了踩坑经验）。**仅 Linux/macOS**
（依赖 Python `pty`）；Windows 用 `scripts/wincons`（真实 Windows 控制台，本机实测可用）。

| 文件 | 作用 |
|---|---|
| `ptydrive.py` | PTY 驱动库：`mark()` 后只检查新输出、`until()` 等模式出现再输入、`wait_exit()` 收尾取退出码、`check()` 记录结果；`Checker` 跨多个短进程汇总结果 |
| `ask_demo.py` | 局部选项的交互收集（`CliOpt.Question`，`ask-demo`）：空值提问一次并采用输入、`--token` 给了值不提问、`-h` 不提问 |
| `shared_ask.py` | 共享选项的交互收集（`SharedOpts()` + `WithQuestion`，`shared-ask run`）：只在**实际执行的命令**处提问一次；父命令段或叶子段给了值都不提问；`-h` 不提问。这条用例专门守住"收集延后到执行命令"的语义（`Command.prepareSharedOpts`） |
| `interact_confirm.py` | cliui 的行式提示在 gcli 应用里可用（`interact` 示例）：`confirm` 答 `n`/`y` 分别走 Unconfirmed/Confirmed、`read` 回显输入。cliui v0.5.1 改签名导致编译不过、v0.5.3 恢复旧签名，这里守住运行侧 |

运行（Linux/macOS）：

```bash
go build -o tmp/smoke/cliapp ./_examples/cliapp
python3 -I scripts/ptysmoke/ask_demo.py --bin tmp/smoke/cliapp
python3 -I scripts/ptysmoke/shared_ask.py --bin tmp/smoke/cliapp
python3 -I scripts/ptysmoke/interact_confirm.py --bin tmp/smoke/cliapp
```

运行（Windows，用 Docker 里的 Python）：

```powershell
$env:GOOS="linux"; $env:GOARCH="amd64"
go build -o tmp/smoke/cliapp ./_examples/cliapp
Remove-Item Env:GOOS, Env:GOARCH
docker run --rm -v ${PWD}:/app -w /app python:3.12-slim `
  sh -c "python3 -I scripts/ptysmoke/ask_demo.py --bin tmp/smoke/cliapp; python3 -I scripts/ptysmoke/shared_ask.py --bin tmp/smoke/cliapp; python3 -I scripts/ptysmoke/interact_confirm.py --bin tmp/smoke/cliapp"
```

每个脚本独立打印 `PASS`/`FAIL` 行并以退出码 0/1 报告结果（全通过为 0）。

编写新用例的要点：

- 断言只看 `s.mark()` 之后的输出（`s.until(pattern, m)`、`s.contains(pattern, m)`），不要在整个缓冲里找。
- 等提示真正出现后再输入（`until()`）；太早写入的输入可能被更早的读取吃掉。
- 终端里回车是 `\r`；`text()` 会剥掉 ANSI，只有颜色相关的断言用 `text_raw()`。
- 每个用例用独立进程（`PtySession` 一次跑一个命令），结束时 `close()` 收尾，避免残留进程。
- 用例默认带 `NO_COLOR=1`、`TERM=xterm`，让断言不受颜色影响；需要看颜色时显式去掉。
