package gcli_test

import (
	"testing"

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
