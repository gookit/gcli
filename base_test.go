package gcli_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gookit/gcli/v3"
	"github.com/stretchr/testify/assert"
)

var (
	buf = new(bytes.Buffer)
)

func newNotExitApp(fns ...func(app *gcli.App)) *gcli.App {
	cli := gcli.New(fns...)
	cli.ExitOnEnd = false

	return cli
}

func TestApp_Hooks(t *testing.T) {
	buf.Reset()

	cli := newNotExitApp()
	cli.On(gcli.EvtAppInit, func(data ...interface{}) bool {
		buf.WriteString("trigger " + gcli.EvtAppInit)
		return true
	})
	cli.Add(simpleCmd)
	assert.Equal(t, "trigger "+gcli.EvtAppInit, buf.String())

	buf.Reset()
	cli.On(gcli.EvtGOptionsParsed, func(data ...interface{}) bool {
		buf.WriteString("trigger " + gcli.EvtGOptionsParsed + ", args:" + fmt.Sprintf("%v", data[1]))
		return true
	})
	cli.Run([]string{"simple"})
	assert.Equal(t, "trigger "+gcli.EvtGOptionsParsed+", args:[simple]", buf.String())
}

func TestApp_On_CmdNotFound(t *testing.T) {
	buf.Reset()

	cli := newNotExitApp()
	cli.Add(simpleCmd)

	fmt.Println("--------- will print command tips ----------")
	cli.On(gcli.EvtCmdNotFound, func(data ...interface{}) bool {
		buf.WriteString("trigger: " + gcli.EvtCmdNotFound)
		buf.WriteString("; command: " + fmt.Sprint(data[1]))
		return true
	})

	cli.Run([]string{"top"})
	assert.Equal(t, "trigger: cmd.not.found; command: top", buf.String())
	buf.Reset()

	fmt.Println("--------- dont print command tips ----------")
	cli.On(gcli.EvtCmdNotFound, func(data ...interface{}) bool {
		buf.WriteString("trigger: " + gcli.EvtCmdNotFound)
		buf.WriteString("; command: " + fmt.Sprint(data[1]))
		return false
	})

	cli.Run([]string{"top"})
	assert.Equal(t, "trigger: cmd.not.found; command: top", buf.String())
}

func TestApp_On_CmdNotFound_redirect(t *testing.T) {
	buf.Reset()
	simpleCmd.ClearData()
	assert.Equal(t, nil, simpleCmd.Value("simple"))

	cli := newNotExitApp()
	cli.Add(simpleCmd)

	fmt.Println("--------- redirect to run another command ----------")
	cli.On(gcli.EvtCmdNotFound, func(data ...interface{}) bool {
		buf.WriteString("trigger:" + gcli.EvtCmdNotFound)
		buf.WriteString(" - command:" + fmt.Sprint(data[1]))
		buf.WriteString("; redirect:simple - ")

		app := data[0].(*gcli.App)
		err := app.Exec("simple", nil)
		assert.NoError(t, err)
		buf.WriteString("value:" + simpleCmd.StrValue("simple"))
		return false
	})

	cli.Run([]string{"top"})
	want := "trigger:cmd.not.found - command:top; redirect:simple - value:simple command"
	assert.Equal(t, want, buf.String())
}
