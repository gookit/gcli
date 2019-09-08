package gcli_test

import (
	"strings"
	"testing"

	"github.com/gookit/gcli/v2"
	"github.com/stretchr/testify/assert"
)

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

func TestArgument(t *testing.T) {
	is := assert.New(t)
	arg := gcli.NewArgument("arg0", "arg desc")

	is.False(arg.IsArray)
	is.False(arg.Required)
	is.False(arg.HasValue())

	is.Equal("arg0", arg.Name)
	is.Equal("arg desc", arg.Description)
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
	is.Equal([]string{"ab","cd"}, arg.StringSplit())
	is.Equal([]string{"ab","cd"}, arg.StringSplit(","))

	// int value
	arg.Value = 23
	is.Equal(23, arg.Int())
	is.Equal("", arg.String())

	// string int value
	arg.Value = "23"
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
