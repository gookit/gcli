package gcli_test

import (
	"bytes"
	"testing"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v2"
	"github.com/stretchr/testify/assert"
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
		Name:   "m1:c3",
		UseFor: "{$cmd} desc",
		Config: func(c *gcli.Command) {
			is.Equal("m1:c3", c.Name)
		},
	})

	is.True(app.IsCommand("c1"))
	is.True(app.IsCommand("c2"))
	is.True(app.IsCommand("m1:c3"))
	is.Len(app.Names(), 3)
	is.NotEmpty(app.Names())

	c := gcli.NewCommand("mdl:test", "desc test2", func(c *gcli.Command) {
	})
	app.AddCommand(c)

	is.Equal("mdl", c.Module())
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

	// show command help
	app.Args = []string{"./myapp", "help", "test"}
	code = app.Run()
	str = buf.String()
	buf.Reset()
	is.Equal(0, code)
	is.Contains(str, "Name: test")
	is.Contains(str, "Desc for test command")

	// run an command
	app.Args = []string{"./myapp", "test"}
	code = app.Run()
	is.Equal(0, code)
	is.Equal("test", cmdRet)

	// other
	// app.AddError(fmt.Errorf("test error"))

	gcli.SetVerbose(gcli.VerbQuiet)
}
