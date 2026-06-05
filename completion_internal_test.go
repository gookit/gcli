package gcli

import (
	"testing"

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
