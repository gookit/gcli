package gcli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/x/assert"
)

func newCompletionApp() *gcli.App {
	app := gcli.NewApp(gcli.NotExitOnEnd())
	app.Add(&gcli.Command{
		Name:    "build",
		Desc:    "compile packages and dependencies",
		Aliases: []string{"b"},
		Config: func(c *gcli.Command) {
			var name string
			c.StrOpt(&name, "name", "n", "", "the name option")
		},
		Func: func(c *gcli.Command, _ []string) error { return nil },
	})
	app.Add(&gcli.Command{
		Name: "clean",
		Desc: "remove object files",
		Config: func(c *gcli.Command) {
			var force bool
			c.BoolOpt(&force, "force", "f", false, "force clean")
		},
		Func: func(c *gcli.Command, _ []string) error { return nil },
	})
	return app
}

// TestApp_GenCompletionScript 默认产**瘦(动态)**脚本: 委托 --in-completion 取候选,
// 脚本里不应硬编码命令名/选项名。
func TestApp_GenCompletionScript(t *testing.T) {
	app := newCompletionApp()
	binName := strings.TrimSuffix(app.BinName(), ".exe")

	t.Run("bash", func(t *testing.T) {
		script, err := app.GenCompletionScript(gcli.BashShell)
		assert.NoErr(t, err)
		assert.NotEmpty(t, script)

		// 瘦脚本特征: bin 名、委托回调、bash 注册指令
		assert.StrContains(t, script, binName)
		assert.StrContains(t, script, "--in-completion")
		assert.StrContains(t, script, "_complete_for_")
		assert.StrContains(t, script, "complete -F")
		// 不应硬编码命令名/选项名(交给 --in-completion 动态计算)
		assert.StrNotContains(t, script, "build")
		assert.StrNotContains(t, script, "clean")
		assert.StrNotContains(t, script, "--name")
		assert.StrNotContains(t, script, "--force")
	})

	t.Run("zsh", func(t *testing.T) {
		script, err := app.GenCompletionScript(gcli.ZshShell)
		assert.NoErr(t, err)
		assert.NotEmpty(t, script)

		assert.StrContains(t, script, binName)
		assert.StrContains(t, script, "--in-completion")
		assert.StrContains(t, script, "compdef")
		// 不应硬编码命令名
		assert.StrNotContains(t, script, "build")
		assert.StrNotContains(t, script, "clean")
	})

	t.Run("pwsh", func(t *testing.T) {
		script, err := app.GenCompletionScript(gcli.PwshShell)
		assert.NoErr(t, err)
		assert.NotEmpty(t, script)

		// pwsh 瘦脚本特征: 注册原生补全 + 委托回调 + bin 名
		assert.StrContains(t, script, "Register-ArgumentCompleter")
		assert.StrContains(t, script, "--in-completion")
		assert.StrContains(t, script, binName)
		// 不应硬编码命令名
		assert.StrNotContains(t, script, "build")
	})

	t.Run("invalid shell", func(t *testing.T) {
		script, err := app.GenCompletionScript("fish")
		assert.Err(t, err)
		assert.Empty(t, script)
	})

	t.Run("override bin name", func(t *testing.T) {
		// 传入自定义 bin 名(对应 genac 的 --bin-name), 脚本正文应使用它
		script, err := app.GenCompletionScript(gcli.BashShell, "./myapp.exe")
		assert.NoErr(t, err)
		// 规整后应为 myapp
		assert.StrContains(t, script, "_complete_for_myapp")
		assert.StrContains(t, script, "complete -F _complete_for_myapp myapp")
		// 委托回调使用规整后的 bin 名
		assert.StrContains(t, script, `"myapp" --in-completion`)
	})
}

func TestApp_GenCompletionHelp(t *testing.T) {
	app := newCompletionApp()

	help := app.GenCompletionHelp("./myapp.exe")
	assert.StrContains(t, help, "myapp --gen-completion bash")
	assert.StrContains(t, help, "myapp --gen-completion zsh")
	assert.StrContains(t, help, "myapp --gen-completion pwsh")
	assert.StrContains(t, help, `eval "$(myapp --gen-completion zsh)"`)
	assert.StrContains(t, help, "$PROFILE")
	assert.StrNotContains(t, help, "./myapp.exe")
}

// TestApp_GenStaticCompletionScript 静态(嵌入式)脚本应把命令名/选项名硬编码进脚本。
func TestApp_GenStaticCompletionScript(t *testing.T) {
	app := newCompletionApp()
	binName := strings.TrimSuffix(app.BinName(), ".exe")

	t.Run("bash", func(t *testing.T) {
		script, err := app.GenStaticCompletionScript(gcli.BashShell)
		assert.NoErr(t, err)
		assert.NotEmpty(t, script)

		// 关键串: bin 名、补全函数、注册的命令名/选项
		assert.StrContains(t, script, binName)
		assert.StrContains(t, script, "_complete_for_")
		assert.StrContains(t, script, "complete -F")
		assert.StrContains(t, script, "build")
		assert.StrContains(t, script, "clean")
		assert.StrContains(t, script, "--name")
		assert.StrContains(t, script, "--force")
	})

	t.Run("zsh", func(t *testing.T) {
		script, err := app.GenStaticCompletionScript(gcli.ZshShell)
		assert.NoErr(t, err)
		assert.NotEmpty(t, script)

		assert.StrContains(t, script, binName)
		assert.StrContains(t, script, "compdef")
		assert.StrContains(t, script, "build")
		assert.StrContains(t, script, "clean")
		// zsh 模板包含命令描述(注意 Desc 首字母会被自动大写)
		assert.StrContains(t, script, "packages and dependencies")
	})

	t.Run("pwsh not supported for static", func(t *testing.T) {
		// pwsh 仅支持动态(瘦)脚本, 静态嵌入不支持 -> 返回 error
		script, err := app.GenStaticCompletionScript(gcli.PwshShell)
		assert.Err(t, err)
		assert.Empty(t, script)
	})

	t.Run("invalid shell", func(t *testing.T) {
		script, err := app.GenStaticCompletionScript("fish")
		assert.Err(t, err)
		assert.Empty(t, script)
	})

	t.Run("override bin name", func(t *testing.T) {
		script, err := app.GenStaticCompletionScript(gcli.BashShell, "./myapp.exe")
		assert.NoErr(t, err)
		assert.StrContains(t, script, "_complete_for_myapp")
		assert.StrContains(t, script, "complete -F _complete_for_myapp myapp")
	})
}

func TestApp_genCompletionOpt(t *testing.T) {
	// App 复用包级 gOpts 单例, 用例结束后需重置, 避免污染其他用例
	defer gcli.ResetGOpts()

	app := newCompletionApp()

	// --gen-completion 命中即生成并退出, 不再继续运行命令
	code := app.Run([]string{"--gen-completion", "bash"})
	assert.Eq(t, 0, code)
	assert.Eq(t, "", app.CommandName())

	// 帮助信息中应包含该(非隐藏)选项
	help := app.Flags().BuildHelp()
	assert.True(t, strings.Contains(help, "gen-completion"))
}

func TestApp_genCompletionOptHelp(t *testing.T) {
	defer gcli.ResetGOpts()
	defer color.ResetOptions()

	app := newCompletionApp()
	buf := new(bytes.Buffer)
	color.Disable()
	color.SetOutput(buf)

	code := app.Run([]string{"--gen-completion", "help"})
	out := buf.String()

	assert.Eq(t, 0, code)
	assert.StrContains(t, out, "Generate shell completion script")
	assert.StrContains(t, out, "Supported shells: bash, zsh, pwsh")
	assert.StrContains(t, out, "--gen-completion bash")
	assert.StrContains(t, out, `eval "$(`)
	assert.StrContains(t, out, "$PROFILE")
	assert.Eq(t, "", app.CommandName())
}
