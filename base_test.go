package gcli_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/events"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/testutil/assert"
)

var (
	buf = new(bytes.Buffer)
)

func newNotExitApp(fns ...func(app *gcli.App)) *gcli.App {
	cli := gcli.New(fns...)
	cli.ExitOnEnd = false

	return cli
}

func TestApp_Hooks_EvtAppInit(t *testing.T) {
	buf.Reset()

	cli := newNotExitApp()
	cli.On(events.OnAppInitAfter, func(ctx *gcli.HookCtx) bool {
		buf.WriteString("trigger " + ctx.Name())
		return false
	})
	cli.Add(simpleCmd)
	assert.Eq(t, "trigger "+events.OnAppInitAfter, buf.String())

	buf.Reset()
	cli.On(events.OnGlobalOptsParsed, func(ctx *gcli.HookCtx) bool {
		buf.WriteString("trigger " + ctx.Name() + ", args:" + fmt.Sprintf("%v", ctx.Strings("args")))
		return false
	})

	cli.Run([]string{"simple"})
	assert.Eq(t, "trigger "+events.OnGlobalOptsParsed+", args:[simple]", buf.String())
}

func TestApp_Hooks_OnAppCmdAdd(t *testing.T) {
	buf.Reset()

	cli := newNotExitApp()
	cli.On(events.OnAppCmdAdd, func(ctx *gcli.HookCtx) (stop bool) {
		buf.WriteString(ctx.Name())
		buf.WriteString(" - ")
		buf.WriteString(ctx.Cmd.Name + ";")
		return
	})

	cli.Add(emptyCmd)
	assert.Eq(t, "app.cmd.add.before - empty;", buf.String())

	cli.Add(simpleCmd)
	assert.Eq(t, "app.cmd.add.before - empty;app.cmd.add.before - simple;", buf.String())
}

func TestCommand_Hooks_EvtCmdOptParsed(t *testing.T) {
	buf.Reset()

	cli := newNotExitApp()
	cli.Add(&gcli.Command{
		Name: "test",
		Desc: "desc",
		Config: func(c *gcli.Command) {
			buf.WriteString("run config;")
			c.On(events.OnCmdOptParsed, func(ctx *gcli.HookCtx) (stop bool) {
				dump.P(ctx.Strings("args"))
				buf.WriteString(ctx.Name())
				return
			})
		},
	})
	assert.Contains(t, buf.String(), "run config;")

	cli.Run([]string{"test"})
	assert.Contains(t, buf.String(), events.OnCmdOptParsed)
}

func TestApp_On_CmdNotFound(t *testing.T) {
	buf.Reset()

	cli := newNotExitApp()
	cli.Add(simpleCmd)

	fmt.Println("--------- will print command tips ----------")
	cli.On(events.OnCmdNotFound, func(ctx *gcli.HookCtx) bool {
		buf.WriteString("trigger: " + events.OnCmdNotFound)
		buf.WriteString("; command: " + ctx.Str("name"))
		return false
	})

	cli.Run([]string{"top"})
	assert.Eq(t, "trigger: cmd.not.found; command: top", buf.String())
	buf.Reset()

	fmt.Println("--------- dont print command tips ----------")
	cli.On(events.OnCmdNotFound, func(ctx *gcli.HookCtx) bool {
		buf.WriteString("trigger: " + events.OnCmdNotFound)
		buf.WriteString("; command: " + ctx.Str("name"))
		return true
	})

	cli.Run([]string{"top"})
	assert.Eq(t, "trigger: cmd.not.found; command: top", buf.String())
}

func TestApp_On_CmdNotFound_redirect(t *testing.T) {
	buf.Reset()
	simpleCmd.Init()
	simpleCmd.ResetData()
	assert.Eq(t, nil, simpleCmd.Ctx.Get("simple"))

	cli := newNotExitApp()
	cli.Add(simpleCmd)

	fmt.Println("--------- redirect to run another command ----------")
	cli.On(events.OnCmdNotFound, func(ctx *gcli.HookCtx) bool {
		buf.WriteString("trigger:" + events.OnCmdNotFound)
		buf.WriteString(" - command:" + ctx.Str("name"))
		buf.WriteString("; redirect:simple - ")

		err := ctx.App.Exec(simpleCmd.Name, nil)
		assert.NoErr(t, err)
		buf.WriteString("value:" + simpleCmd.Ctx.Str("simple"))
		return true
	})

	cli.Run([]string{"top"})
	want := "trigger:cmd.not.found - command:top; redirect:simple - value:simple command"
	assert.Eq(t, want, buf.String())
}
