package gcli_test

import (
	"fmt"
	"testing"

	"github.com/gookit/gcli/v2"
	"github.com/stretchr/testify/assert"
)

var simpleArgs = []string{"hi"}

func TestNewCommand(t *testing.T) {
	is := assert.New(t)

	c := gcli.NewCommand("test", "desc test", func(c *gcli.Command) {
		c.Aliases = []string{"alias1", "alias2"}
	})

	is.NotEmpty(c)
	is.Nil(c.App())

	err := c.Run(simpleArgs)
	is.NoError(err)
	is.True(c.IsAlone())
	is.False(c.NotAlone())

	is.False(c.IsDisabled())
	c.Disable()
	is.True(c.IsDisabled())

	c.Logf(gcli.VerbInfo, "command log")

	// is.Equal("", c.ArgLine())
	is.Equal("alias1,alias2", c.AliasesString())
	is.Equal("alias1alias2", c.AliasesString(""))
}

func TestCommand_Errorf(t *testing.T) {
	is := assert.New(t)

	c := gcli.NewCommand("test", "desc test", nil)
	c.SetFunc(func(c *gcli.Command, args []string) error {
		is.Equal([]string{"hi"}, args)
		return c.Errorf("error message")
	})

	is.NotEmpty(c)

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

	c := gcli.NewCommand("test", "desc test", func(c *gcli.Command) {
		is.Equal("test", c.Name)
		c.Aliases = []string{"alias1"}
	})
	c.SetFunc(func(c *gcli.Command, args []string) error {
		return nil
	})

	is.NotEmpty(c)
	err := c.Run(simpleArgs)
	is.NoError(err)

	is.Equal("alias1", c.AliasesString(""))

	err = c.Run([]string{"-h"})
	is.NoError(err)
	is.Equal("alias1", c.AliasesString(""))

	g := gcli.NewApp()
	g.AddCommand(c)
	err = c.Run(simpleArgs)
	is.Error(err)
}

func TestCommand_ParseFlag(t *testing.T) {
	is := assert.New(t)

	var int0 int
	var str0 string

	c := gcli.NewCommand("test", "desc test", func(c *gcli.Command) {
		c.IntOpt(&int0, "int", "", 0, "int desc")
		c.StrOpt(&str0, "str", "", "", "str desc")
		is.Equal("test", c.Name)
		is.Equal("int desc", c.FlagMeta("int").Desc)
	})

	c.SetFunc(func(c *gcli.Command, args []string) error {
		is.Equal("test", c.Name)
		is.Equal([]string{"txt"}, args)
		return nil
	})

	err := c.Run([]string{"txt", "--int", "10", "--str=abc"})
	is.NoError(err)
	is.Equal(10, int0)
	is.Equal("abc", str0)
	is.Equal([]string{"txt"}, c.RawArgs())
	is.Equal("txt", c.RawArg(0))

	// var str0 string
	co := struct {
		maxSteps  int
		overwrite bool
	}{}

	c = gcli.NewCommand("test", "desc test", func(c *gcli.Command) {
		is.Equal("test", c.Name)
		c.IntOpt(&int0, "int", "", 0, "desc")
		c.IntOpt(&co.maxSteps, "max-step", "", 0, "setting the max step value")
	})
	c.SetFunc(func(c *gcli.Command, args []string) error {
		is.Equal("[txt]", fmt.Sprint(args))
		return nil
	})

	err = c.Run([]string{"--int", "10", "--max-step=100", "txt"})
	is.NoError(err)
	is.Equal(10, int0)
	is.Equal(100, co.maxSteps)
	is.Equal("[txt]", fmt.Sprint(c.RawArgs()))
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
