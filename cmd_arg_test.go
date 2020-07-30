package gcli_test

import (
	"strings"
	"testing"

	"github.com/gookit/gcli/v2"
	"github.com/stretchr/testify/assert"
)

func TestCommand_AddArg(t *testing.T) {
	is := assert.New(t)
	c := gcli.NewCommand("test", "test desc", nil)

	arg := c.AddArg("arg0", "arg desc", true)
	is.Equal(0, arg.Index())

	ret := c.ArgByIndex(0)
	is.Equal(ret, arg)

	arg = c.AddArg("arg1", "arg1 desc")
	is.Equal(1, arg.Index())

	ret = c.Arg("arg1")
	is.Equal(ret, arg)

	ret = c.Arg("not-exist")
	is.True(ret.IsEmpty())
	is.False(ret.HasValue())

	is.Len(c.Args(), 2)

	is.PanicsWithValue("GCLI: the command argument name cannot be empty", func() {
		c.AddArg("", "desc")
	})
	is.PanicsWithValue("GCLI: the command argument name ':)&dfd' is invalid, only allow: a-Z 0-9 _ -", func() {
		c.AddArg(":)&dfd", "desc")
	})
	is.PanicsWithValue("GCLI: the argument name 'arg1' already exists in command 'test'", func() {
		c.AddArg("arg1", "desc")
	})
	is.PanicsWithValue("GCLI: required argument 'arg2' cannot be defined after optional argument", func() {
		c.AddArg("arg2", "arg2 desc", true)
	})

	c.AddArg("arg3", "arg3 desc", false, true)
	is.PanicsWithValue("GCLI: have defined an array argument, you cannot add argument 'argN'", func() {
		c.AddArg("argN", "desc", true)
	})
}

func TestArgument(t *testing.T) {
	is := assert.New(t)
	arg := gcli.NewArgument("arg0", "arg desc")

	is.False(arg.IsArray)
	is.False(arg.Required)
	is.False(arg.IsEmpty())
	is.False(arg.HasValue())

	is.Equal("arg0", arg.Name)
	is.Equal("arg desc", arg.Desc)
	is.Equal(0, arg.Index())

	// no value
	is.Nil(arg.Strings())
	is.Nil(arg.GetValue())
	is.Nil(arg.StringSplit())
	is.Equal(0, arg.Int())
	is.Equal(34, arg.Int(34))
	is.Equal("", arg.String())
	is.Equal("ab", arg.String("ab"))

	// add value
	arg.Value = "ab,cd"

	is.Nil(arg.Strings())
	is.Equal(0, arg.Int())
	is.Equal(34, arg.Int(34))

	is.Equal("ab,cd", arg.String())
	is.Equal([]string{"ab", "cd"}, arg.StringSplit())
	is.Equal([]string{"ab", "cd"}, arg.StringSplit(","))

	// int value
	arg.Value = 23
	is.Equal(23, arg.Int())
	is.Equal("", arg.String())

	// string int value
	arg.WithValue("23")
	is.Equal(23, arg.Int())
	is.Equal("23", arg.String())

	// array value
	arg.IsArray = true
	arg.Value = []string{"a", "b"}
	is.True(arg.IsArray)
	is.Equal(0, arg.Int())
	is.Equal("", arg.String())
	is.Equal([]string{"a", "b"}, arg.Array())

	// custom handler
	arg.Value = "a-b-c"
	arg.Handler = func(value interface{}) interface{} {
		str := value.(string)
		return strings.SplitN(str, "-", 2)
	}
	is.Equal([]string{"a", "b-c"}, arg.GetValue())

	// required and is-array
	arg = gcli.NewArgument("arg1", "arg desc", true, true)
	is.True(arg.IsArray)
	is.True(arg.Required)
}
