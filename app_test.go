package gcli_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/stretchr/testify/assert"
)

var (
	emptyCmd = &gcli.Command{
		Name: "empty",
		Desc: "an test command",
	}
	simpleCmd = &gcli.Command{
		Name: "simple",
		Desc: "an simple command",
		Func: func(c *gcli.Command, args []string) error {
			fmt.Println(c.Path(), args)
			return nil
		},
	}
	subCmd = &gcli.Command{
		Name: "sub",
		Desc: "an simple sub command",
		Func: func(c *gcli.Command, args []string) error {
			fmt.Println(c.Path(), args)
			return nil
		},
	}

	appWithMl = gcli.New(func(app *gcli.App) {
		app.ExitOnEnd = false
		app.Add(&gcli.Command{
			Name: "top1",
			Desc: "desc for top1",
			Subs: []*gcli.Command{
				{
					Name: "sub1",
					Desc: "desc for top1.sub1",
					Func: func(c *gcli.Command, args []string) error {
						c.SetValue("msg", c.App().Value("top1.sub1"))
						return nil
					},
				},
			},
		})
		app.Add(&gcli.Command{
			Name: "top2",
			Desc: "desc for top2",
		})
	})
)

func TestApp_MatchByPath(t *testing.T) {
	// is := assert.New(t)
	app := gcli.NewApp()
	app.Add(
		gcli.NewCommand("cmd1", "desc"),
		gcli.NewCommand("cmd2", "desc2"),
	)

	simpleCmd.AddCommand(subCmd)
	app.AddCommand(simpleCmd)

	assert.True(t, app.HasCommand(simpleCmd.Name))

	c := app.MatchByPath("simple:sub")
	assert.NotNil(t, c)
	assert.Equal(t, "sub", c.Name)
	assert.Equal(t, "simple", c.ParentName())

	c = appWithMl.FindByPath("top1:sub1")
	assert.Equal(t, "sub1", c.Name)
}

func TestApp_Add(t *testing.T) {
	is := assert.New(t)

	app := gcli.NewApp()
	app.Add(gcli.NewCommand("c1", "c1 desc", func(c *gcli.Command) {
		is.Equal("c1", c.Name)
	}), gcli.NewCommand("c2", "c2 desc", func(c *gcli.Command) {
		is.Equal("c2", c.Name)
	}))
	app.AddCommand(&gcli.Command{
		Name:    "c3",
		Desc:    "{$cmd} desc",
		Aliases: []string{"alias1"},
		Config: func(c *gcli.Command) {
			is.Equal("c3", c.Name)
		},
	})

	is.True(app.IsCommand("c1"))
	is.True(app.IsCommand("c2"))
	is.True(app.HasCommand("c3"))
	is.Len(app.CmdNames(), 3)
	is.Len(app.CmdNameMap(), 3)
	is.NotEmpty(app.CommandNames())

	c := gcli.NewCommand("mdl-test", "desc test2")
	app.AddCommand(c)

	is.Equal("", c.ParentName())
	is.Equal("mdl-test", c.Name)
	is.Equal("c3", app.ResolveAlias("alias1"))
	is.True(app.IsAlias("alias1"))
}

func TestApp_AddCommand(t *testing.T) {
	app := gcli.NewApp(func(a *gcli.App) {
		a.ExitOnEnd = false
	})

	app.AddCommand(emptyCmd)

	cmd1 := gcli.NewCommand("cmd1", "desc")
	cmd1.Disable()

	app.AddCommand(cmd1)
	assert.Len(t, app.Commands(), 1)
	assert.Len(t, app.CommandNames(), 1)
	assert.False(t, app.IsCommand("cmd1"))
	assert.Equal(t, "", app.CommandName())

	assert.PanicsWithValue(t, "GCli: the command name can not be empty", func() {
		app.AddCommand(&gcli.Command{})
	})
	assert.PanicsWithValue(t, "GCli: the command name '+xdd' is invalid, must match: ^[a-zA-Z][\\w-]*$", func() {
		app.AddCommand(&gcli.Command{Name: "+xdd"})
	})
}

func TestApp_AddAliases(t *testing.T) {
	app := gcli.NewApp(func(a *gcli.App) {
		a.ExitOnEnd = false
	})

	cmd := &gcli.Command{
		Name:    "test",
		Desc:    "the desc",
		Aliases: []string{"alias1"},
	}
	app.AddCommand(cmd)
	app.AddCommand(gcli.NewCommand("cmd2", "desc"))
	assert.True(t, app.IsAlias("alias1"))

	assert.PanicsWithValue(t, "GCli: The name 'alias1' is already used as an alias", func() {
		app.AddCommand(gcli.NewCommand("alias1", "desc"))
	})
	assert.PanicsWithValue(t, "GCli: The name 'test' has been used as an command name", func() {
		app.AddAliases("cmd2", "test")
	})
}

func TestApp_Run_noCommands(t *testing.T) {
	is := assert.New(t)

	app := gcli.NewApp(func(a *gcli.App) {
		a.ExitOnEnd = false
		a.Version = "1.3.9"
	})
	app.Run([]string{})

	// disable color code, re-set output for test
	buf := new(bytes.Buffer)
	color.Disable()
	color.SetOutput(buf)
	gcli.SetVerbose(gcli.VerbCrazy)

	defer func() {
		color.ResetOptions()
		gcli.SetVerbose(gcli.VerbError)
	}()

	// run
	code := app.Run([]string{})
	str := buf.String()
	buf.Reset()

	is.Equal(0, code)
	is.Contains(str, "1.3.9")
	is.Contains(str, "Version: 1.3.9")
	is.Contains(str, "This is my console application")
	is.Contains(str, "Display help information")

	err := app.Exec("not-exists", []string{})
	is.Error(err)
}

func TestApp_Run_command_withArguments(t *testing.T) {
	is := assert.New(t)
	app := gcli.NewApp(func(a *gcli.App) {
		a.ExitOnEnd = false
	})

	// run with command
	var argStr, cmdRet string
	app.Add(&gcli.Command{
		Name: "test",
		Desc: "desc for test command",
		Config: func(c *gcli.Command) {
			c.AddArg("arg0", "desc")
			c.BindArg(&gcli.Argument{Name: "arg1", Desc: "desc1"})
		},
		Func: func(c *gcli.Command, args []string) error {
			cmdRet = c.Name
			argStr = strings.Join(args, ",")
			return nil
		},
	})

	// run an command
	code := app.Run([]string{"test"})
	is.Equal(0, code)
	is.Equal("", argStr)
	is.Equal("test", cmdRet)
	is.Equal("test", app.CommandName())
	// clear
	argStr = ""
	cmdRet = ""

	err := app.Exec("test", []string{"val0", "val1"})
	is.NoError(err)
	is.Equal("test", cmdRet)
	is.Equal("val0,val1", argStr)

	err = app.Exec("not-exists", []string{})
	is.Error(err)
	is.Equal("exec unknown command: 'not-exists'", err.Error())
	// other
	// app.AddError(fmt.Errorf("test error"))
}

func TestApp_Run_command_withOptions(t *testing.T) {
	is := assert.New(t)
	app := gcli.NewApp(gcli.NotExitOnEnd())

	// run with command
	var optStr, cmdRet string
	var opt1 string

	app.Add(&gcli.Command{
		Name: "test",
		Desc: "desc for test command",
		Config: func(c *gcli.Command) {
			c.AddArg("arg0", "desc")
			c.BindArg(&gcli.Argument{Name: "arg1", Desc: "desc1"})
			c.StrOpt(&opt1, "opt1", "o", "", "opt desc")
		},
		Func: func(c *gcli.Command, args []string) error {
			cmdRet = c.Name
			optStr = strings.Join(args, ",")
			return nil
		},
	})

	// run command
	code := app.Run([]string{"test"})
	is.Equal(0, code)
	is.Equal("", optStr)
	is.Equal("test", cmdRet)
	is.Equal("test", app.CommandName())

	// help option
	app.Run([]string{"test", "-h"})

	// disable color code, reset output for test
	buf := new(bytes.Buffer)
	color.Disable()
	color.SetOutput(buf)
	gcli.SetVerbose(gcli.VerbCrazy)

	defer func() {
		color.ResetOptions()
		gcli.SetVerbose(gcli.VerbError)
	}()

	app.Run([]string{"test", "-h"})
	is.Contains(buf.String(), "-o, --opt1 string")
}

func TestApp_Run_subcommand(t *testing.T) {
	is := assert.New(t)
	id := "top1:sub1"

	appWithMl.SetValue(id, "TestApp_Run_subcommand")
	appWithMl.Run([]string{"top1", "sub1"})

	c := appWithMl.FindCommand(id)
	is.NotEmpty(c)
	is.Equal("TestApp_Run_subcommand", c.Value("msg"))
}

func TestApp_Run_by_cmd_ID(t *testing.T) {
	is := assert.New(t)

	appWithMl.SetValue("top1:sub1", "TestApp_Run_by_cmd_ID")
	appWithMl.Run([]string{"top1:sub1"})

	c := appWithMl.FindCommand("top1:sub1")
	is.NotEmpty(c)
	is.Equal("TestApp_Run_by_cmd_ID", c.Value("msg"))
}

func TestApp_AddAliases_and_run(t *testing.T) {
	is := assert.New(t)
	id := "top1:sub1"

	appWithMl.AddAliases(id, "ts1")
	appWithMl.SetValue(id, "TestApp_AddAliases_and_run")
	appWithMl.Run([]string{"ts1"})

	c := appWithMl.FindCommand(id)
	is.NotEmpty(c)
	is.Equal("TestApp_AddAliases_and_run", c.Value("msg"))
}

func TestApp_showCommandHelp(t *testing.T) {
	is := assert.New(t)

	app := gcli.NewApp(func(a *gcli.App) {
		a.ExitOnEnd = false
	})

	app.AddCommand(gcli.NewCommand("test", "desc for test command"))

	app.Run([]string{"help", "test"})

	// disable color code, re-set output for test
	buf := new(bytes.Buffer)
	color.Disable()
	color.SetOutput(buf)
	defer color.ResetOptions()

	// show command help
	code := app.Run([]string{"help", "test"})
	str := buf.String()
	buf.Reset()
	is.Equal(0, code)
	is.Contains(str, "Name: test")
	is.Contains(str, "Desc for test command")

	// show command help: arg error
	code = app.Run([]string{"help", "test", "more"})
	str = buf.String()
	buf.Reset()
	is.Equal(gcli.ERR, code)
	is.Contains(str, "ERROR: Too many arguments given.")

	// show command help for 'help'
	code = app.Run([]string{"help", "help"})
	str = buf.String()
	buf.Reset()
	is.Equal(gcli.OK, code)
	is.Contains(str, "Display help message for application or command.")

	// show command help: unknown command
	code = app.Run([]string{"help", "not-exist"})
	str = buf.String()
	buf.Reset()
	is.Equal(gcli.ERR, code)
	is.Contains(str, "Unknown command name 'not-exist'")
}

func TestApp_showVersion(t *testing.T) {
	app := gcli.NewApp(func(a *gcli.App) {
		a.ExitOnEnd = false
		a.Version = "1.3.9"
		a.Desc = "application desc"
		a.Logo.Text = "MY-LOGO"
	})

	app.Run([]string{"--version"})

	// disable color code, re-set output for test
	buf := new(bytes.Buffer)
	color.Disable()
	color.SetOutput(buf)
	defer color.ResetOptions()

	app.Run([]string{"--version"})
	str := buf.String()
	buf.Reset()
	assert.Contains(t, str, "Version: 1.3.9")
	assert.Contains(t, str, "Application desc")
	assert.Contains(t, str, "MY-LOGO")
}

func TestApp_showCommandTips(t *testing.T) {
	app := gcli.NewApp()
	app.ExitOnEnd = false

	app.AddCommand(emptyCmd)
	app.Run([]string{"emp"})

	// disable color code, re-set output for test
	buf := new(bytes.Buffer)
	color.Disable()
	color.SetOutput(buf)
	defer color.ResetOptions()

	app.Run([]string{"emp"})
	str := buf.String()
	buf.Reset()
	assert.Contains(t, str, "ERROR: unknown input command \"emp\"")
	assert.Contains(t, str, `Maybe you mean:
  empty`)
}

// func TestApp_RemoveCommand(t *testing.T) {
// 	app := gcli.NewApp()
//
// 	app.Add(
// 		gcli.NewCommand("cmd1", "desc"),
// 		gcli.NewCommand("cmd2", "desc"),
// 	)
//
// 	assert.Len(t, app.Commands(), 2)
// 	assert.True(t, app.IsCommand("cmd1"))
//
// 	assert.Equal(t, 1, app.RemoveCommand("cmd1"))
// 	assert.Len(t, app.Commands(), 1)
// 	assert.False(t, app.IsCommand("cmd1"))
// }

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

func (uc *UserCommand) Config(c *gcli.Command) {
	c.StrOpt(&uc.opt1, "opt", "o", "", "desc")
}

func (uc *UserCommand) Execute(c *gcli.Command, args []string) error {
	return nil
}
