package gcli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/gookit/gcli/v3/events"
	"github.com/gookit/goutil/x/assert"
)

// newCompletionApp 构造一个用于补全测试的应用夹具:
//   - 顶层命令: build(别名 b)、run、help 为内置;
//   - build 含子命令 module(别名 mod) 与选项 --output/-o、--verbose/-v;
//   - run 含选项 --name/-n。
func newCompletionApp() *App {
	app := NewApp(func(a *App) { a.ExitOnEnd = false })

	build := NewCommand("build", "build desc", func(c *Command) {
		var out, verbose string
		c.StrOpt(&out, "output", "o", "the output dir")
		c.StrOpt(&verbose, "verbose", "v", "verbose mode")

		c.AddSubs(NewCommand("module", "module desc", func(sc *Command) {
			sc.Aliases = []string{"mod"}
		}))
	})
	app.Add(build)
	app.AddAliases("build", "b")

	run := NewCommand("run", "run desc", func(c *Command) {
		var name string
		c.StrOpt(&name, "name", "n", "the name")
	})
	app.Add(run)

	return app
}

func TestApp_resolveCompletion(t *testing.T) {
	is := assert.New(t)
	app := newCompletionApp()

	t.Run("top level on empty words", func(t *testing.T) {
		got := app.resolveCompletion(nil)
		// 顶层命令名 + 别名 + help, 去重排序
		is.Eq([]string{"b", "build", "help", "run"}, got)
	})

	t.Run("top level prefix filter", func(t *testing.T) {
		got := app.resolveCompletion([]string{"b"})
		// 前缀 b: 命中 b、build
		is.Eq([]string{"b", "build"}, got)
	})

	t.Run("sub commands of a command", func(t *testing.T) {
		got := app.resolveCompletion([]string{"build", ""})
		// build 的子命令名 + 子命令别名
		is.Eq([]string{"mod", "module"}, got)
	})

	t.Run("sub commands prefix filter", func(t *testing.T) {
		got := app.resolveCompletion([]string{"build", "mod"})
		is.Eq([]string{"mod", "module"}, got)
	})

	t.Run("option names when cur is dash", func(t *testing.T) {
		got := app.resolveCompletion([]string{"build", "-"})
		// build 的选项: 长选项 + 短选项(顶层全局选项不参与命令级补全)
		is.Contains(got, "--output")
		is.Contains(got, "-o")
		is.Contains(got, "--verbose")
		is.Contains(got, "-v")
	})

	t.Run("option names with long prefix", func(t *testing.T) {
		got := app.resolveCompletion([]string{"build", "--o"})
		is.Eq([]string{"--output"}, got)
	})

	t.Run("alias resolves then drill down", func(t *testing.T) {
		// 顶层别名 b -> build, 应能下钻并产出 build 的子命令
		got := app.resolveCompletion([]string{"b", ""})
		is.Eq([]string{"mod", "module"}, got)
	})

	t.Run("sub command alias options", func(t *testing.T) {
		// build mod -> build module, mod 别名解析后定位到 module 节点(无子命令、无选项)
		got := app.resolveCompletion([]string{"build", "mod", ""})
		is.Empty(got)
	})

	t.Run("stop after unknown word", func(t *testing.T) {
		// arg 不是命令(视为参数), 上下文停留在 build, cur 为空 -> 仍补全 build 子命令
		got := app.resolveCompletion([]string{"build", "arg", ""})
		is.Eq([]string{"mod", "module"}, got)
	})

	t.Run("option words are skipped on locating", func(t *testing.T) {
		// 中间的选项词 -o val 在定位时被跳过(注意: -o 跳过, val 视为参数会停止下钻),
		// 这里用纯选项词验证跳过逻辑: build --verbose -> 仍在 build, 补全选项
		got := app.resolveCompletion([]string{"build", "--verbose", "-"})
		is.Contains(got, "--output")
		is.Contains(got, "--verbose")
	})

	t.Run("run command options", func(t *testing.T) {
		got := app.resolveCompletion([]string{"run", "-"})
		is.Contains(got, "--name")
		is.Contains(got, "-n")
	})

	t.Run("hidden global option excluded", func(t *testing.T) {
		// 顶层选项补全: 可见全局选项可补全, 但隐藏的内部选项 --in-completion 不应出现
		got := app.resolveCompletion([]string{"-"})
		is.Contains(got, "--help")
		is.Contains(got, "--gen-completion")
		is.NotContains(got, "--in-completion")
	})
}

// captureStdout 捕获 fn 执行期间写入 os.Stdout 的内容(showAutoCompletion 用 fmt.Println 直接写 stdout)。
func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = old

	out, _ := io.ReadAll(r)
	return string(out)
}

// TestApp_completionMode_suppressHooks 验证补全模式静默: 当本次运行是补全请求时,
// 用户 init 钩子(OnAppInitAfter)不被触发, 且 stdout 只剩候选, 不含钩子里的噪声输出。
//
// 注意: initialize() 常由首个 Add() 触发(早于 Run), 因此真实判定来自 os.Args。
// 这里临时改写 os.Args 以忠实复现真实执行路径(钩子在 NewApp config 中注册、Add 触发 init)。
func TestApp_completionMode_suppressHooks(t *testing.T) {
	is := assert.New(t)

	const noise = "init app event app.init.after"
	hookFired := false

	// 临时改写 os.Args 模拟真实补全请求: bin --in-completion build ""
	oldArgs := os.Args
	os.Args = []string{"app", "--in-completion", "build", ""}
	defer func() { os.Args = oldArgs }()
	// Run 会改写包级 gOpts 单例, 用例结束后重置避免污染其他用例
	defer ResetGOpts()

	// 在 NewApp 的 config 函数里注册钩子(早于 Add 触发的 initialize)
	app := NewApp(func(a *App) {
		a.ExitOnEnd = false
		a.On(events.OnAppInitAfter, func(_ *HookCtx) bool {
			hookFired = true
			fmt.Println(noise) // 钩子里的噪声输出, 委托式脚本会误解析
			return false
		})
	})
	build := NewCommand("build", "build desc", func(c *Command) {
		c.AddSubs(NewCommand("module", "module desc", func(sc *Command) {
			sc.Aliases = []string{"mod"}
		}))
	})
	// Add 会触发 initialize(); 此时应已根据 os.Args 进入 completion 模式, 钩子被抑制
	app.Add(build)

	// 补全请求: 当前正在补全 build 的子命令(末尾空词)
	out := captureStdout(func() {
		app.Run([]string{"--in-completion", "build", ""})
	})

	// ① 钩子未触发(completion 模式下 init 钩子被抑制)
	is.False(hookFired, "OnAppInitAfter hook should NOT fire in completion mode")
	// ② stdout 不含噪声, 只有候选
	is.False(strings.Contains(out, noise), "stdout should not contain hook noise, got: %q", out)
	// ③ 输出包含真实候选(build 的子命令 module/别名 mod)
	is.True(strings.Contains(out, "module"), "stdout should contain candidate 'module', got: %q", out)
}

// TestApp_completionMode_hooksFireNormally 反向验证: 非补全模式下 init 钩子正常触发。
func TestApp_completionMode_hooksFireNormally(t *testing.T) {
	is := assert.New(t)

	// 确保 os.Args 不含补全元选项
	oldArgs := os.Args
	os.Args = []string{"app", "--version"}
	defer func() { os.Args = oldArgs }()
	defer ResetGOpts()

	hookFired := false
	app := NewApp(func(a *App) {
		a.ExitOnEnd = false
		a.On(events.OnAppInitAfter, func(_ *HookCtx) bool {
			hookFired = true
			return false
		})
	})
	app.Add(NewCommand("build", "build desc", nil))

	// 普通运行(无补全元选项): 钩子应正常触发
	app.Run([]string{"--version"})
	is.True(hookFired, "OnAppInitAfter hook SHOULD fire in normal mode")
}
