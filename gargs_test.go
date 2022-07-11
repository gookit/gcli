package gcli_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/gookit/gcli/v3"
	assert2 "github.com/gookit/goutil/testutil/assert"
	"github.com/stretchr/testify/assert"
)

func TestCommand_AddArg(t *testing.T) {
	is := assert.New(t)
	c := gcli.NewCommand("test", "test desc", nil)

	arg := c.AddArg("arg0", "arg desc", true)
	is.Equal(0, arg.Index())

	ret := c.ArgByIndex(0)
	is.Equal(ret, arg)

	assert2.PanicsMsg(t, func() {
		c.ArgByIndex(1)
	}, "GCli: get not exists argument #1")

	arg = c.AddArg("arg1", "arg1 desc")
	is.Equal(1, arg.Index())

	ret = c.Arg("arg1")
	is.Equal(ret, arg)

	assert2.PanicsMsg(t, func() {
		c.Arg("not-exist")
	}, "GCli: get not exists argument 'not-exist'")

	is.Len(c.Args(), 2)

	is.PanicsWithValue("GCli: the command argument name cannot be empty", func() {
		c.AddArg("", "desc")
	})

	assert2.PanicsMsg(t, func() {
		c.AddArg(":)&dfd", "desc")
	}, "GCli: the argument name ':)&dfd' is invalid, must match: ^[a-zA-Z][\\w-]*$")

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

func TestArguments_AddArgByRule(t *testing.T) {
	is := assert.New(t)
	ags := gcli.Arguments{}

	arg := ags.AddArgByRule("arg2", "arg2 desc;false;23")
	is.Equal("arg2 desc", arg.Desc)
	is.Equal(23, arg.Int())
	is.Equal(false, arg.Arrayed)
}

func TestArguments_BindArg(t *testing.T) {
	is := assert.New(t)
	ags := gcli.Arguments{}

	ags.BindArg(&gcli.Argument{Name: "ag0"})
	is.True(ags.HasArg("ag0"))
}

func TestArgument(t *testing.T) {
	is := assert2.New(t)
	arg := gcli.NewArgument("arg0", "arg desc")

	is.False(arg.Arrayed)
	is.False(arg.Required)
	is.True(arg.IsEmpty())
	is.False(arg.HasValue())

	is.Eq("arg0", arg.Name)
	is.Eq("arg desc", arg.Desc)
	is.Eq(0, arg.Index())

	// no value
	is.Nil(arg.Strings())
	is.Nil(arg.GetValue())
	is.Nil(arg.SplitToStrings())
	is.Eq(0, arg.Int())
	is.Eq("", arg.String())
	is.Eq("ab", arg.WithValue("ab").String())

	// add value
	err := arg.SetValue("ab,cd")
	is.NoErr(err)

	is.Eq(0, arg.Int())
	is.Eq("ab,cd", arg.String())
	is.Eq([]string{"ab", "cd"}, arg.Array())
	is.Eq([]string{"ab", "cd"}, arg.SplitToStrings(","))

	// int value
	is.NoErr(arg.SetValue(23))
	is.Eq(23, arg.Int())
	is.Eq("23", arg.String())

	// string int value
	err = arg.SetValue("23")
	is.NoErr(err)
	is.Eq(23, arg.Int())
	is.Eq("23", arg.String())

	// array value
	arg.Arrayed = true
	is.NoErr(arg.SetValue([]string{"a", "b"}))
	is.True(arg.Arrayed)
	is.Eq(0, arg.Int())
	is.Eq("[a b]", arg.String())
	is.Eq([]string{"a", "b"}, arg.Array())

	// required and is-array
	arg = gcli.NewArgument("arg1", "arg desc", true, true)
	arg.Init()
	is.True(arg.Arrayed)
	is.True(arg.Required)
	is.Eq("arg1...", arg.HelpName())
}

func TestArgument_GetValue(t *testing.T) {
	arg := gcli.NewArgument("arg0", "arg desc")

	// custom handler
	assert.NoError(t, arg.SetValue("a-b-c"))
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
	arg := gcli.NewArgument("arg0", "arg desc").WithFn(func(arg *gcli.Argument) {
		arg.SetValue(23)
		arg.Init()
	})

	assert.Equal(t, 23, arg.Val())
	assert.Equal(t, "arg0", arg.HelpName())
}

func TestArgument_SetValue(t *testing.T) {
	arg := gcli.NewArgument("arg0", "arg desc")
	// convert "12" to 12
	arg.Validator = str2int

	err := arg.SetValue("12")
	assert.NoError(t, err)
	assert.Equal(t, 12, arg.Val())
	arg.Set(nil) // reset value

	err = arg.SetValue("abc")
	assert.Error(t, err)
	assert.Nil(t, arg.Val())

	// convert "12" to 12
	arg = gcli.NewArgument("arg0", "arg desc").WithValidator(str2int)
	err = arg.SetValue("12")
	assert.NoError(t, err)
	assert.Equal(t, 12, arg.Val())
}
