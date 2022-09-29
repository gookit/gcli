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
	cli.On(events.OnAppInit, func(ctx *gcli.HookCtx) bool {
		buf.WriteString("trigger " + events.OnAppInit)
		return false
	})
	cli.Add(simpleCmd)
	assert.Eq(t, "trigger "+events.OnAppInit, buf.String())

	buf.Reset()
	cli.On(events.OnGlobalOptsParsed, func(ctx *gcli.HookCtx) bool {
		buf.WriteString("trigger " + ctx.Name() + ", args:" + fmt.Sprintf("%v", ctx.Strings("args")))
		return false
	})

	cli.Run([]string{"simple"})
	assert.Eq(t, "trigger "+events.OnGlobalOptsParsed+", args:[simple]", buf.String())
}

func TestApp_Hooks_EvtCmdInit(t *testing.T) {
	buf.Reset()

	cli := newNotExitApp()
	cli.On(events.OnCmdInit, func(ctx *gcli.HookCtx) (stop bool) {
		buf.WriteString(events.OnCmdInit)
		buf.WriteString(":")

		buf.WriteString(ctx.Cmd.Name + ";")
		return
	})

	cli.Add(emptyCmd)
	assert.Eq(t, "cmd.init:empty;", buf.String())

	cli.Add(simpleCmd)
	assert.Eq(t, "cmd.init:empty;cmd.init:simple;", buf.String())
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
	assert.Eq(t, nil, simpleCmd.GetVal("simple"))

	cli := newNotExitApp()
	cli.Add(simpleCmd)

	fmt.Println("--------- redirect to run another command ----------")
	cli.On(events.OnCmdNotFound, func(ctx *gcli.HookCtx) bool {
		buf.WriteString("trigger:" + events.OnCmdNotFound)
		buf.WriteString(" - command:" + ctx.Str("name"))
		buf.WriteString("; redirect:simple - ")

		err := ctx.App.Exec(simpleCmd.Name, nil)
		assert.NoErr(t, err)
		buf.WriteString("value:" + simpleCmd.Data.Str("simple"))
		return true
	})

	cli.Run([]string{"top"})
	want := "trigger:cmd.not.found - command:top; redirect:simple - value:simple command"
	assert.Eq(t, want, buf.String())
}
