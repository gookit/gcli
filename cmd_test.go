package gcli_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/dump"
	"github.com/stretchr/testify/assert"
)

var simpleArgs = []string{"hi"}

func TestNewCommand(t *testing.T) {
	is := assert.New(t)

	c := gcli.NewCommand("test", "desc test", func(c *gcli.Command) {
		c.Aliases = []string{"alias1", "alias2"}
	})

	is.NotEmpty(c)
	is.False(c.Runnable())
	is.Nil(c.App())

	err := c.Run([]string{})
	is.NoError(err)
	is.True(c.IsStandalone())
	is.False(c.NotStandalone())

	is.False(c.IsDisabled())
	c.Disable()
	is.True(c.IsDisabled())

	// is.Equal("", c.ArgLine())
	is.Equal("alias1,alias2", c.AliasesString())
	is.Equal("alias1alias2", c.AliasesString(""))

	c = gcli.NewCommand("test1", "desc test")
	app := gcli.NewApp()
	c.AttachTo(app)
	is.True(app.HasCommand("test1"))
}

func TestCommand_NewErrf(t *testing.T) {
	is := assert.New(t)

	c := gcli.NewCommand("test", "desc test", func(c *gcli.Command) {
		c.AddArg("arg0", "desc message")
	})
	c.SetFunc(func(c *gcli.Command, args []string) error {
		is.Equal([]string{"hi"}, args)
		return c.NewErrf("error message")
	})

	is.NotEmpty(c)
	is.True(c.Runnable())

	err := c.Run(simpleArgs)
	is.Error(err)
	is.Equal("error message", err.Error())
	is.Equal([]string{"hi"}, c.RawArgs())

	is.Panics(func() {
		c.MustRun(simpleArgs)
	})
}

func TestCommand_Run(t *testing.T) {
	is := assert.New(t)

	// use struct
	c := &gcli.Command{
		Name: "test",
		Desc: "test desc",
		Config: func(c *gcli.Command) {
			is.Equal("test", c.Name)
			c.Aliases = []string{"alias1"}
		},
		Func: func(c *gcli.Command, args []string) error {
			return nil
		},
	}

	err := c.Run([]string{})
	is.NoError(err)
	is.True(c.IsStandalone())
	is.False(c.NotStandalone())
	is.Equal("alias1", c.AliasesString(""))

	err = c.Run([]string{"-h"})
	is.NoError(err)
	is.Equal("alias1", c.AliasesString(""))
}

func TestNewCommand_Run(t *testing.T) {
	is := assert.New(t)

	c := gcli.NewCommand("test", "desc test", func(c *gcli.Command) {
		is.Equal("test", c.Name)
		c.Aliases = []string{"alias1"}
	})
	c.SetFunc(func(c *gcli.Command, args []string) error {
		return nil
	})

	is.NotEmpty(c)
	err := c.Run([]string{})
	is.NoError(err)
	is.True(c.IsStandalone())
	is.False(c.NotStandalone())

	is.Equal("alias1", c.AliasesString(""))

	err = c.Run([]string{"-h"})
	is.NoError(err)
	is.Equal("alias1", c.AliasesString(""))

	// error run on app
	g := gcli.NewApp()
	g.AddCommand(c)
	err = c.Run(simpleArgs)
	is.Error(err)

}

var bf = new(bytes.Buffer)

// l0: root command
var r = &gcli.Command{
	Name: "git",
	Desc: "git usage",
	Subs: []*gcli.Command{
		// l1: sub command 1
		{
			Name: "add",
			Desc: "the add command for git",
			Config: func(c *gcli.Command) {
				c.AddArg("files", "added files", true)
			},
			Func: func(c *gcli.Command, args []string) error {
				bf.WriteString("command path: " + c.Path())
				dump.Println(c.Name, args)
				return nil
			},
		},
		// l1: sub command 2
		{
			Name:    "pull",
			Desc:    "the pull command for git",
			Aliases: []string{"pul"},
			Func: func(c *gcli.Command, args []string) error {
				bf.WriteString("command path: " + c.Path())
				dump.Println(c.Name, args)
				return nil
			},
		},
		// l1: sub command 3
		{
			Name:    "remote",
			Desc:    "remote command for git",
			Aliases: []string{"rmt"},
			Func: func(c *gcli.Command, args []string) error {
				dump.Println(c.Path())
				return nil
			},
			Subs: []*gcli.Command{
				// l2: sub command 4
				{
					Name: "add",
					Desc: "add command for git remote",
					Config: func(c *gcli.Command) {
						c.AddArg("name", "the remote name", true)
						c.AddArg("address", "the remote address", true)
					},
					Func: func(c *gcli.Command, args []string) error {
						bf.WriteString("command path: " + c.Path())
						dump.Println(c.Path(), args)
						return nil
					},
				},
				// l2: sub command 5
				{
					Name:    "set-url",
					Desc:    "set-url command for git remote",
					Aliases: []string{"su"},
					Func: func(c *gcli.Command, args []string) error {
						bf.WriteString("command path: " + c.Path())
						dump.Println(c.Path(), args)
						return nil
					},
				},
			},
		},
	},
	Func: func(c *gcli.Command, args []string) error {
		bf.WriteString("command path: " + c.Path())
		dump.Println(c.Path(), args)
		return nil
	},
}

func TestCommand_MatchByPath(t *testing.T) {
	c := r.MatchByPath("add")

	assert.NotNil(t, c)
	assert.Equal(t, "add", c.Name)
	assert.Equal(t, "git", c.ParentName())

	c = r.MatchByPath("remote:add")
	assert.NotNil(t, c)
	assert.Equal(t, "add", c.Name)
	assert.Equal(t, "Add command for git remote", c.Desc)
	assert.Equal(t, "remote", c.Parent().Name)
	assert.Equal(t, "git", c.Root().Name)

	// empty will return self
	c = r.MatchByPath("")
	assert.NotNil(t, c)
	assert.Equal(t, "git", c.Name)

	c = r.MatchByPath("not-exist")
	assert.Nil(t, c)
}

func TestCommand_Sub(t *testing.T) {
	r.MatchByPath("") // use for init

	assert.True(t, r.IsRoot())
	assert.True(t, r.IsCommand("remote"))
	assert.True(t, r.IsCommand("remote"))

	c := r.Sub("add")
	assert.NotNil(t, c)
	assert.False(t, c.IsRoot())
	assert.Equal(t, "add", c.Name)
}

func TestCommand_Run_top(t *testing.T) {
	bf.Reset() // reset buffer

	err := r.Run([]string{})
	assert.NoError(t, err)
	assert.Equal(t, "command path: git", bf.String())
}

func TestCommand_Run_oneLevelSub(t *testing.T) {
	bf.Reset() // reset buffer

	err := r.Run([]string{"add", "./"})
	assert.NoError(t, err)
}

func TestCommand_Run_moreLevelSub(t *testing.T) {
	bf.Reset() // reset buffer
	err := r.Run([]string{
		"remote",
		"add",
		"origin",
		"https://github.com/inhere/console",
	})

	assert.NoError(t, err)
	assert.True(t, r.IsAlias("rmt"))
	assert.True(t, r.IsAlias("pul"))
	assert.False(t, r.IsAlias("not-exist"))
	assert.Equal(t, "remote", r.ResolveAlias("rmt"))
	assert.Equal(t, "command path: git remote add", bf.String())
}

var int0 int
var str0 string

var c0 = gcli.NewCommand("test", "desc test", func(c *gcli.Command) {
	c.IntOpt(&int0, "int", "", 0, "int desc")
	c.StrOpt(&str0, "str", "", "", "str desc")
	c.AddArg("arg0", "arg0 desc")
	c.AddArg("arg1", "arg1 desc")
	c.Func = func(c *gcli.Command, args []string) error {
		bf.WriteString("name=" + c.Name)
		c.SetValue("name", c.Name)
		c.SetValue("args", args)
		dump.P(c.ID(), "command Func is exec")
		return nil
	}
})

func TestCommand_Run_emptyArgs(t *testing.T) {
	bf.Reset()
	is := assert.New(t)

	gcli.SetCrazyMode()
	defer gcli.ResetVerbose()

	is.Equal("test", c0.Name)

	err := c0.Run([]string{})

	is.NoError(err)
	is.Equal("name=test", bf.String())
	is.Equal("int desc", c0.FlagMeta("int").Desc)
	is.NotEmpty(c0.Args())
	is.Equal("arg0", c0.Arg("arg0").Name)
}

func TestCommand_Run_parseHelp(t *testing.T) {
	bf.Reset()
	is := assert.New(t)

	err := c0.Run([]string{"-h"})
	is.NoError(err)

	// no color
	color.Disable()
	color.SetOutput(bf)
	defer color.ResetOptions()

	err = c0.Run([]string{"--help"})
	str := bf.String()
	is.Contains(str, "Int desc")
	is.Contains(str, "--str string")
	is.Contains(str, "Str desc")
	is.Contains(str, "Display the help information")
	is.Contains(str, "arg0        Arg0 desc")
	is.Contains(str, "arg1        Arg1 desc")
}

func TestCommand_Run_parseOptions(t *testing.T) {
	bf.Reset()
	is := assert.New(t)

	gcli.ResetGOpts()
	gcli.SetDebugMode()
	defer gcli.ResetVerbose()

	is.Equal("test", c0.Name)

	err := c0.Run([]string{"--int", "10", "--str=abc", "txt"})

	is.NoError(err)
	is.Equal("test", c0.Value("name"))
	is.Equal([]string{"txt"}, c0.Value("args"))

	is.Equal(10, int0)
	is.Equal("abc", str0)
	is.Equal([]string{"txt"}, c0.FSetArgs())
	is.Equal("txt", c0.RawArg(0))

	// var str0 string
	co := struct {
		maxSteps  int
		overwrite bool
	}{}

	c1 := gcli.NewCommand("test1", "desc test", func(c *gcli.Command) {
		c.IntOpt(&int0, "int", "", 0, "desc")
		c.IntOpt(&co.maxSteps, "max-step", "", 0, "setting the max step value")
		c.AddArg("arg0", "arg0 desc")
	})
	c1.SetFunc(func(c *gcli.Command, args []string) error {
		is.Equal("[txt]", fmt.Sprint(args))
		return nil
	})

	is.Equal("test1", c1.Name)
	err = c1.Run([]string{"--int", "10", "--max-step=100", "txt"})
	is.NoError(err)
	is.Equal(10, int0)
	is.Equal(100, co.maxSteps)
	is.Equal("[txt]", fmt.Sprint(c0.RawArgs()))
}

func TestInts(t *testing.T) {
	is := assert.New(t)
	ints := gcli.Ints{}

	err := ints.Set("1")
	is.NoError(err)
	err = ints.Set("3")
	is.NoError(err)
	is.Equal("[1 3]", ints.String())
	err = ints.Set("abc")
	is.Error(err)

	ints = gcli.Ints{1, 3}
	is.Equal("[1 3]", ints.String())
}

func TestStrings(t *testing.T) {
	is := assert.New(t)
	ss := gcli.Strings{}

	err := ss.Set("1")
	is.NoError(err)
	err = ss.Set("3")
	is.NoError(err)
	err = ss.Set("abc")
	is.NoError(err)
	is.Equal("[1 3 abc]", ss.String())
}

func TestBooleans(t *testing.T) {
	is := assert.New(t)
	val := gcli.Booleans{}

	err := val.Set("false")
	is.NoError(err)
	is.False(val[0])
	is.Equal("[false]", val.String())

	err = val.Set("True")
	is.NoError(err)
	is.Equal("[false true]", val.String())

	err = val.Set("abc")
	is.Error(err)
}
