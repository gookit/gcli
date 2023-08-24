package gflag_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/testutil/assert"
)

func TestArguments_AddArgByRule(t *testing.T) {
	is := assert.New(t)
	ags := gflag.Arguments{}

	arg := ags.AddArgByRule("arg2", "arg2 desc;false;23")
	is.Eq("arg2 desc", arg.Desc)
	is.Eq(23, arg.Int())
	is.Eq(false, arg.Arrayed)
}

func TestArguments_BindArg(t *testing.T) {
	is := assert.New(t)
	ags := gflag.Arguments{}

	ags.BindArg(&gflag.Argument{Name: "ag0"})
	is.True(ags.HasArg("ag0"))
}

func TestArgument(t *testing.T) {
	is := assert.New(t)
	arg := gflag.NewArgument("arg0", "arg desc")

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
	is.Eq("abc", arg.WithDefault("abc").String())

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
	arg.WithArrayed()
	is.NoErr(arg.SetValue([]string{"a", "b"}))
	is.True(arg.Arrayed)
	is.Eq(0, arg.Int())
	is.Eq("[a b]", arg.String())
	is.Eq([]string{"a", "b"}, arg.Array())

	arg = gflag.NewArgument("arg0", "arg desc").SetArrayed()
	is.True(arg.Arrayed)

	// required and is-array
	arg = gflag.NewArgument("arg1", "arg desc", true, true)
	arg.Init()
	is.True(arg.Arrayed)
	is.True(arg.Required)
	is.Eq("arg1...", arg.HelpName())
}

func TestArgument_GetValue(t *testing.T) {
	arg := gflag.NewArgument("arg0", "arg desc")

	// custom handler
	assert.NoErr(t, arg.SetValue("a-b-c"))
	arg.Handler = func(value any) any {
		str := value.(string)
		return strings.SplitN(str, "-", 2)
	}
	assert.Eq(t, []string{"a", "b-c"}, arg.GetValue())
}

var str2int = func(val any) (any, error) {
	return strconv.Atoi(val.(string))
}

func TestArgument_WithConfig(t *testing.T) {
	arg := gflag.NewArgument("arg0", "arg desc").WithFn(func(arg *gflag.CliArg) {
		_ = arg.SetValue(23)
		arg.Init()
	})

	assert.Eq(t, 23, arg.Val())
	assert.Eq(t, "arg0", arg.HelpName())
}

func TestArgument_SetValue(t *testing.T) {
	arg := gflag.NewArgument("arg0", "arg desc")
	// convert "12" to 12
	arg.Validator = str2int

	err := arg.SetValue("12")
	assert.NoErr(t, err)
	assert.Eq(t, 12, arg.Val())
	arg.Set(nil) // reset value

	err = arg.SetValue("abc")
	assert.Err(t, err)
	assert.Nil(t, arg.Val())

	// convert "12" to 12
	arg = gflag.NewArgument("arg0", "arg desc").WithValidator(str2int)
	err = arg.SetValue("12")
	assert.NoErr(t, err)
	assert.Eq(t, 12, arg.Val())
}
