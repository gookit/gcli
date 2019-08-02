package gcli_test

import (
	"fmt"
	"testing"

	"github.com/gookit/gcli/v2"
	"github.com/stretchr/testify/assert"
)

var simpleArgs = []string{"hi"}

func TestNewCommand(t *testing.T) {
	ris := assert.New(t)

	c := gcli.NewCommand("test", "desc test", func(c *gcli.Command) {
		c.Aliases = []string{"alias1", "alias2"}
	})

	ris.NotEmpty(c)
	ris.Nil(c.App())

	err := c.Run(simpleArgs)
	ris.NoError(err)
	ris.True(c.IsAlone())
	ris.False(c.NotAlone())

	// ris.Equal("", c.ArgLine())
	ris.Equal("alias1,alias2", c.AliasesString())
	ris.Equal("alias1alias2", c.AliasesString(""))
}

func TestCommand_Errorf(t *testing.T) {
	ris := assert.New(t)

	c := gcli.NewCommand("test", "desc test", nil)
	c.SetFunc(func(c *gcli.Command, args []string) error {
		ris.Equal([]string{"hi"}, args)
		return c.Errorf("error message")
	})

	ris.NotEmpty(c)

	err := c.Run(simpleArgs)
	ris.Error(err)
	ris.Equal("error message", err.Error())
	ris.Equal([]string{"hi"}, c.RawArgs())

	ris.Panics(func() {
		c.MustRun(simpleArgs)
	})
}

func TestCommand_Run(t *testing.T) {
	ris := assert.New(t)

	c := gcli.NewCommand("test", "desc test", func(c *gcli.Command) {
		ris.Equal("test", c.Name)
		c.Aliases = []string{"alias1"}
	})
	c.SetFunc(func(c *gcli.Command, args []string) error {
		return nil
	})

	ris.NotEmpty(c)
	err := c.Run(simpleArgs)
	ris.NoError(err)

	ris.Equal("alias1", c.AliasesString(""))
}

func TestCommand_ParseFlag(t *testing.T) {
	ris := assert.New(t)

	var int0 int
	var str0 string

	c := gcli.NewCommand("test", "desc test", func(c *gcli.Command) {
		c.IntOpt(&int0, "int", "", 0, "int desc")
		c.StrOpt(&str0, "str", "", "", "str desc")
		ris.Equal("test", c.Name)
		ris.Equal("int desc", c.OptDes("int"))
	})
	c.SetFunc(func(c *gcli.Command, args []string) error {
		ris.Equal("test", c.Name)
		ris.Equal([]string{"txt"}, args)
		return nil
	})

	err := c.Run([]string{"txt", "--int", "10", "--str=abc"})
	ris.NoError(err)
	ris.Equal(10, int0)
	ris.Equal("abc", str0)
	ris.Equal([]string{"txt"}, c.RawArgs())
	ris.Equal("txt", c.RawArg(0))

	// var str0 string
	co := struct {
		maxSteps  int
		overwrite bool
	}{}

	c = gcli.NewCommand("test", "desc test", func(c *gcli.Command) {
		ris.Equal("test", c.Name)
		c.IntOpt(&int0, "int", "", 0, "desc")
		c.IntOpt(&co.maxSteps, "max-step", "", 0, "setting the max step value")
	})
	c.SetFunc(func(c *gcli.Command, args []string) error {
		ris.Equal("[txt]", fmt.Sprint(args))
		return nil
	})

	err = c.Run([]string{"--int", "10", "--max-step=100", "txt"})
	ris.NoError(err)
	ris.Equal(10, int0)
	ris.Equal(100, co.maxSteps)
	ris.Equal("[txt]", fmt.Sprint(c.RawArgs()))
}
