package gcli_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/gookit/gcli/v3"
	"github.com/stretchr/testify/assert"
)

func TestCommand_AddArg(t *testing.T) {
	is := assert.New(t)
	c := gcli.NewCommand("test", "test desc", nil)

	arg := c.AddArg("arg0", "arg desc", true)
	is.Equal(0, arg.Index())

	ret := c.ArgByIndex(0)
	is.Equal(ret, arg)
	ret = c.ArgByIndex(1)
	is.True(ret.IsEmpty())

	arg = c.AddArg("arg1", "arg1 desc")
	is.Equal(1, arg.Index())

	ret = c.Arg("arg1")
	is.Equal(ret, arg)

	ret = c.Arg("not-exist")
	is.True(ret.IsEmpty())
	is.False(ret.HasValue())

	is.Len(c.Args(), 2)

	is.PanicsWithValue("GCli: the command argument name cannot be empty", func() {
		c.AddArg("", "desc")
	})

	is.PanicsWithValue("GCli: the command argument name ':)&dfd' is invalid, must match: ^[a-zA-Z][\\w-]*$", func() {
		c.AddArg(":)&dfd", "desc")
	})

	is.PanicsWithValue("GCli: the argument name 'arg1' already exists in command 'test'", func() {
		c.AddArg("arg1", "desc")
	})
	is.PanicsWithValue("GCli: required argument 'arg2' cannot be defined after optional argument", func() {
		c.AddArg("arg2", "arg2 desc", true)
	})

	c.AddArg("arg3", "arg3 desc", false, true)
	is.PanicsWithValue("GCli: have defined an array argument, you cannot add argument 'argN'", func() {
		c.AddArg("argN", "desc", true)
	})
}

func TestArguments_BindArg(t *testing.T) {
	is := assert.New(t)
	ags := gcli.Arguments{}

	ags.BindArg(&gcli.Argument{Name: "ag0"})
	is.True(ags.HasArg("ag0"))
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
	err := arg.SetValue("23")
	is.NoError(err)
	is.Equal(23, arg.Int())
	is.Equal("23", arg.String())

	// array value
	arg.IsArray = true
	arg.Value = []string{"a", "b"}
	is.True(arg.IsArray)
	is.Equal(0, arg.Int())
	is.Equal("", arg.String())
	is.Equal([]string{"a", "b"}, arg.Array())

	// required and is-array
	arg = gcli.NewArgument("arg1", "arg desc", true, true)
	arg.Init()
	is.True(arg.IsArray)
	is.True(arg.Required)
	is.Equal("arg1...", arg.HelpName())
}

func TestArgument_GetValue(t *testing.T) {
	arg := gcli.NewArgument("arg0", "arg desc")

	// custom handler
	arg.Value = "a-b-c"
	arg.Handler = func(value interface{}) interface{} {
		str := value.(string)
		return strings.SplitN(str, "-", 2)
	}
	assert.Equal(t, []string{"a", "b-c"}, arg.GetValue())
}

var str2int = func(val interface{}) (interface{}, error) {
	return strconv.Atoi(val.(string))
}

func TestArgument_WithConfig(t *testing.T) {
	arg := gcli.NewArgument("arg0", "arg desc").With(func(arg *gcli.Argument) {
		arg.Value = 23
		arg.Init()
	})

	assert.Equal(t, 23, arg.Value)
	assert.Equal(t, "arg0", arg.HelpName())
}

func TestArgument_SetValue(t *testing.T) {
	arg := gcli.NewArgument("arg0", "arg desc")
	// convert "12" to 12
	arg.Validator = str2int

	err := arg.SetValue("12")
	assert.NoError(t, err)
	assert.Equal(t, 12, arg.Value)
	arg.Value = nil // reset value

	err = arg.SetValue("abc")
	assert.Error(t, err)
	assert.Nil(t, arg.Value)

	// convert "12" to 12
	arg = gcli.NewArgument("arg0", "arg desc").WithValidator(str2int)
	err = arg.SetValue("12")
	assert.NoError(t, err)
	assert.Equal(t, 12, arg.Value)
}
