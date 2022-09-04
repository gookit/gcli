package gcli_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gookit/gcli/v3"
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
	cli.On(gcli.EvtAppInit, func(data ...any) bool {
		buf.WriteString("trigger " + gcli.EvtAppInit)
		return false
	})
	cli.Add(simpleCmd)
	assert.Eq(t, "trigger "+gcli.EvtAppInit, buf.String())

	buf.Reset()
	cli.On(gcli.EvtGOptionsParsed, func(data ...any) bool {
		buf.WriteString("trigger " + gcli.EvtGOptionsParsed + ", args:" + fmt.Sprintf("%v", data[1]))
		return false
	})
	cli.Run([]string{"simple"})
	assert.Eq(t, "trigger "+gcli.EvtGOptionsParsed+", args:[simple]", buf.String())
}

func TestApp_Hooks_EvtCmdInit(t *testing.T) {
	buf.Reset()

	cli := newNotExitApp()
	cli.On(gcli.EvtCmdInit, func(data ...any) (stop bool) {
		buf.WriteString(gcli.EvtCmdInit)
		buf.WriteString(":")

		c := data[1].(*gcli.Command)
		buf.WriteString(c.Name + ";")
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
			c.On(gcli.EvtCmdOptParsed, func(data ...any) (stop bool) {
				dump.P(data[1])
				buf.WriteString(gcli.EvtCmdOptParsed)
				return
			})
		},
	})
	assert.Contains(t, buf.String(), "run config;")

	cli.Run([]string{"test"})
	assert.Contains(t, buf.String(), gcli.EvtCmdOptParsed)
}

func TestApp_On_CmdNotFound(t *testing.T) {
	buf.Reset()

	cli := newNotExitApp()
	cli.Add(simpleCmd)

	fmt.Println("--------- will print command tips ----------")
	cli.On(gcli.EvtCmdNotFound, func(data ...any) bool {
		buf.WriteString("trigger: " + gcli.EvtCmdNotFound)
		buf.WriteString("; command: " + fmt.Sprint(data[1]))
		return false
	})

	cli.Run([]string{"top"})
	assert.Eq(t, "trigger: cmd.not.found; command: top", buf.String())
	buf.Reset()

	fmt.Println("--------- dont print command tips ----------")
	cli.On(gcli.EvtCmdNotFound, func(data ...any) bool {
		buf.WriteString("trigger: " + gcli.EvtCmdNotFound)
		buf.WriteString("; command: " + fmt.Sprint(data[1]))
		return true
	})

	cli.Run([]string{"top"})
	assert.Eq(t, "trigger: cmd.not.found; command: top", buf.String())
}

func TestApp_On_CmdNotFound_redirect(t *testing.T) {
	buf.Reset()
	simpleCmd.ClearData()
	assert.Eq(t, nil, simpleCmd.GetVal("simple"))

	cli := newNotExitApp()
	cli.Add(simpleCmd)

	fmt.Println("--------- redirect to run another command ----------")
	cli.On(gcli.EvtCmdNotFound, func(data ...any) bool {
		buf.WriteString("trigger:" + gcli.EvtCmdNotFound)
		buf.WriteString(" - command:" + fmt.Sprint(data[1]))
		buf.WriteString("; redirect:simple - ")

		app := data[0].(*gcli.App)
		err := app.Exec("simple", nil)
		assert.NoErr(t, err)
		buf.WriteString("value:" + simpleCmd.StrValue("simple"))
		return true
	})

	cli.Run([]string{"top"})
	want := "trigger:cmd.not.found - command:top; redirect:simple - value:simple command"
	assert.Eq(t, want, buf.String())
}
