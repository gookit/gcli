package gcli

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
)

// white-box test: findCommandName 现为无副作用纯函数，逐分支覆盖。
func TestApp_findCommandName_pure(t *testing.T) {
	is := assert.New(t)

	app := NewApp(func(a *App) { a.ExitOnEnd = false })
	app.Add(NewCommand("top", "top desc", func(c *Command) {
		c.AddSubs(NewCommand("sub", "sub desc"))
	}))
	app.AddAliases("top", "tp")

	t.Run("founded normal", func(t *testing.T) {
		fc := app.findCommandName([]string{"top", "-a", "x"})
		is.Eq(Founded, fc.state)
		is.Eq("top", fc.name)
		is.Eq("top", fc.raw)
		is.Eq([]string{"-a", "x"}, fc.args)
	})

	t.Run("founded via alias", func(t *testing.T) {
		fc := app.findCommandName([]string{"tp", "arg"})
		is.Eq(Founded, fc.state)
		is.Eq("top", fc.name)
		is.Eq("tp", fc.raw)
		is.Eq([]string{"arg"}, fc.args)
	})

	t.Run("founded via command-ID expands sub into args", func(t *testing.T) {
		fc := app.findCommandName([]string{"top:sub", "x"})
		is.Eq(Founded, fc.state)
		is.Eq("top", fc.name)
		is.Eq("top:sub", fc.raw)
		is.Eq([]string{"sub", "x"}, fc.args)
	})

	t.Run("notfound option", func(t *testing.T) {
		fc := app.findCommandName([]string{"-h"})
		is.Eq(NotFound, fc.state)
		is.Eq("", fc.name)
		is.Eq("", fc.raw)
	})

	t.Run("notfound unknown name", func(t *testing.T) {
		fc := app.findCommandName([]string{"nope", "x"})
		is.Eq(NotFound, fc.state)
		is.Eq("nope", fc.name)
		is.Eq("nope", fc.raw)
		is.Eq([]string{"x"}, fc.args)
	})

	t.Run("founded via default command on empty args", func(t *testing.T) {
		app.SetDefaultCommand("top")
		fc := app.findCommandName(nil)
		is.Eq(Founded, fc.state)
		is.Eq("top", fc.name)
		is.Empty(fc.args)
		app.SetDefaultCommand("")
	})

	t.Run("notfound empty on empty args without default", func(t *testing.T) {
		fc := app.findCommandName(nil)
		is.Eq(NotFound, fc.state)
		is.Eq("", fc.name)
	})

	// 关键：上述调用不应改动 app 状态(无副作用)
	is.Empty(app.args)
	is.Eq("", app.inputName)
}
