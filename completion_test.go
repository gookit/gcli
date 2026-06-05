package gcli_test

import (
	"strings"
	"testing"

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

func TestApp_GenCompletionScript(t *testing.T) {
	app := newCompletionApp()
	binName := app.BinName()

	t.Run("bash", func(t *testing.T) {
		script, err := app.GenCompletionScript(gcli.BashShell)
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
		script, err := app.GenCompletionScript(gcli.ZshShell)
		assert.NoErr(t, err)
		assert.NotEmpty(t, script)

		assert.StrContains(t, script, binName)
		assert.StrContains(t, script, "compdef")
		assert.StrContains(t, script, "build")
		assert.StrContains(t, script, "clean")
		// zsh 模板包含命令描述(注意 Desc 首字母会被自动大写)
		assert.StrContains(t, script, "packages and dependencies")
	})

	t.Run("invalid shell", func(t *testing.T) {
		script, err := app.GenCompletionScript("pwsh")
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
