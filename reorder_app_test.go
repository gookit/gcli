package gcli_test

import (
	"testing"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/x/assert"
)

// options written after arguments are still parsed on the executed command.
func TestApp_Run_reorderArgs(t *testing.T) {
	is := assert.New(t)
	app := gcli.NewApp(gcli.NotExitOnEnd())

	var name string
	var force bool
	var a0, a1 string
	app.Add(&gcli.Command{
		Name: "test",
		Desc: "desc for test command",
		Config: func(c *gcli.Command) {
			c.StrOpt(&name, "name", "n", "", "name opt")
			c.BoolOpt(&force, "force", "f", false, "force opt")
			c.AddArg("arg0", "arg0 desc")
			c.AddArg("arg1", "arg1 desc")
		},
		Func: func(c *gcli.Command, _ []string) error {
			a0 = c.Arg("arg0").String()
			a1 = c.Arg("arg1").String()
			return nil
		},
	})

	code := app.Run([]string{"test", "v0", "--name", "tom", "v1", "-f"})
	is.Eq(0, code)
	is.Eq("tom", name)
	is.True(force)
	is.Eq("v0", a0)
	is.Eq("v1", a1)
}

// in a multi-level app, only the final executed command's args are reordered.
func TestApp_Run_reorderArgs_multiLevel(t *testing.T) {
	is := assert.New(t)
	app := gcli.NewApp(gcli.NotExitOnEnd())

	var subName string
	var a0, a1 string
	app.Add(&gcli.Command{
		Name: "top",
		Desc: "top command",
		Subs: []*gcli.Command{
			{
				Name: "sub",
				Desc: "sub command",
				Config: func(c *gcli.Command) {
					c.StrOpt(&subName, "name", "n", "", "name opt")
					c.AddArg("arg0", "arg0 desc")
					c.AddArg("arg1", "arg1 desc")
				},
				Func: func(c *gcli.Command, _ []string) error {
					a0 = c.Arg("arg0").String()
					a1 = c.Arg("arg1").String()
					return nil
				},
			},
		},
	})

	code := app.Run([]string{"top", "sub", "f1", "--name", "tom", "f2"})
	is.Eq(0, code)
	is.Eq("tom", subName)
	is.Eq("f1", a0)
	is.Eq("f2", a1)
}
