package gcli_test

import (
	"testing"

	"github.com/gookit/gcli"
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
		return c.Errorf("error message")
	})

	ris.NotEmpty(c)

	err := c.Run(simpleArgs)
	ris.Error(err)
	ris.Equal("error message", err.Error())

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
