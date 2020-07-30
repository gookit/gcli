package gcli_test

import (
	"bytes"
	"testing"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v2"
	"github.com/stretchr/testify/assert"
)

var (
	emptyCmd = &gcli.Command{
		Name:   "empty",
		UseFor: "an test command",
	}
)

func TestStdApp(t *testing.T) {
	is := assert.New(t)

	gcli.InitStdApp(func(a *gcli.App) {
		a.Name = "test-name"
	})

	app := gcli.StdApp()
	app.Config(func(a *gcli.App) {
		a.Logo = gcli.Logo{
			Text:  "logo1",
			Style: "warn",
		}
	})

	is.Equal("test-name", app.Name)
	is.Empty(app.Commands())
	is.Equal("logo1", app.Logo.Text)
	is.Equal("warn", app.Logo.Style)

	app.SetLogo("logo2", "info")
	is.Equal("logo2", app.Logo.Text)
	is.Equal("info", app.Logo.Style)
}

func TestApp_Add(t *testing.T) {
	is := assert.New(t)

	app := gcli.NewApp()
	app.Add(app.NewCommand("c1", "c1 desc", func(c *gcli.Command) {
		is.Equal("c1", c.Name)
	}), gcli.NewCommand("c2", "c2 desc", func(c *gcli.Command) {
		is.Equal("c2", c.Name)
	}))
	app.AddCommand(&gcli.Command{
		Name:    "m1:c3",
		UseFor:  "{$cmd} desc",
		Aliases: []string{"alias1"},
		Config: func(c *gcli.Command) {
			is.Equal("m1:c3", c.Name)
		},
	})

	is.True(app.IsCommand("c1"))
	is.True(app.IsCommand("c2"))
	is.True(app.HasCommand("m1:c3"))
	is.Len(app.Names(), 3)
	is.NotEmpty(app.Names())

	c := gcli.NewCommand("mdl:test", "desc test2")
	app.AddCommand(c)

	is.Equal("mdl", c.Module())
	is.Equal("m1:c3", app.ResolveName("alias1"))
	is.True(app.IsAlias("alias1"))
}

func TestApp_Run(t *testing.T) {
	is := assert.New(t)

	gcli.SetVerbose(gcli.VerbCrazy)
	app := gcli.NewApp(func(a *gcli.App) {
		a.ExitOnEnd = false
		a.Version = "1.3.9"
	})

	// disable color code, re-set output for test
	buf := new(bytes.Buffer)
	color.Disable()
	color.SetOutput(buf)

	// run
	app.Args = []string{"./myapp"}
	code := app.Run()
	str := buf.String()
	buf.Reset()

	is.Equal(0, code)
	is.Contains(str, "1.3.9")
	is.Contains(str, "Version: 1.3.9")
	is.Contains(str, "This is my CLI application")
	is.Contains(str, "Display help information")

	cmdRet := ""
	app.Add(&gcli.Command{
		Name:   "test",
		UseFor: "desc for test command",
		Func: func(c *gcli.Command, args []string) error {
			cmdRet = c.Name
			return nil
		},
	})

	// run an command
	app.Args = []string{"./myapp", "test"}
	code = app.Run()
	is.Equal(0, code)
	is.Equal("test", cmdRet)

	// other
	// app.AddError(fmt.Errorf("test error"))

	gcli.SetVerbose(gcli.VerbQuiet)
}

func TestApp_showCommandHelp(t *testing.T) {
	is := assert.New(t)

	// disable color code, re-set output for test
	buf := new(bytes.Buffer)
	color.Disable()
	color.SetOutput(buf)

	app := gcli.NewApp(func(a *gcli.App) {
		a.ExitOnEnd = false
	})

	app.AddCommand(gcli.NewCommand("test", "desc for test command", func(c *gcli.Command) {
		//
	}))

	// show command help
	app.Args = []string{"./myapp", "help", "test"}
	code := app.Run()
	str := buf.String()
	buf.Reset()
	is.Equal(0, code)
	is.Contains(str, "Name: test")
	is.Contains(str, "Desc for test command")

	// show command help: arg error
	app.Args = []string{"./myapp", "help", "test", "more"}
	code = app.Run()
	str = buf.String()
	buf.Reset()
	is.Equal(gcli.ERR, code)
	is.Contains(str, "ERROR: Too many arguments given.")

	// show command help for 'help'
	app.Args = []string{"./myapp", "help", "help"}
	code = app.Run()
	str = buf.String()
	buf.Reset()
	is.Equal(gcli.OK, code)
	is.Contains(str, "Display help message for application or command.")

	// show command help: unknown command
	app.Args = []string{"./myapp", "help", "not-exist"}
	code = app.Run()
	str = buf.String()
	buf.Reset()
	is.Equal(gcli.ERR, code)
	is.Contains(str, "Unknown command name 'not-exist'")
}

func TestApp_AddCommander(t *testing.T) {
	app := gcli.NewApp()

	app.AddCommander(&UserCommand{})

	assert.True(t, app.HasCommand("test"))
}

// UserCommand for tests
type UserCommand struct {
	opt1 string
}

func (uc *UserCommand) Creator() *gcli.Command {
	return gcli.NewCommand("test", "desc message")
}

func (uc *UserCommand) Prepare(c *gcli.Command) {
	c.StrOpt(&uc.opt1, "opt", "o", "", "desc")
}

func (uc *UserCommand) Run(c *gcli.Command, args []string) error {
	return nil
}
