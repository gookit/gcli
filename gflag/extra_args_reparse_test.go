package gflag_test

import (
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/x/assert"
	"testing"
)

func TestParseArgsClearsPreviousExtras(t *testing.T) {
	t.Run("named consumed", func(t *testing.T) {
		a := gflag.Arguments{}
		a.AddArg("file", "input", true)
		assert.NoErr(t, a.ParseArgs([]string{"first", "old"}))
		assert.Eq(t, []string{"old"}, a.ExtraArgs())
		assert.NoErr(t, a.ParseArgs([]string{"second"}))
		assert.Empty(t, a.ExtraArgs())
		assert.Eq(t, "second", a.Arg("file").String())
	})
	t.Run("no named args", func(t *testing.T) {
		a := gflag.Arguments{}
		assert.NoErr(t, a.ParseArgs([]string{"old"}))
		assert.NoErr(t, a.ParseArgs([]string{}))
		assert.Empty(t, a.ExtraArgs())
	})
	t.Run("missing required", func(t *testing.T) {
		a := gflag.Arguments{}
		a.AddArg("file", "input", true)
		assert.NoErr(t, a.ParseArgs([]string{"first", "old"}))
		assert.Err(t, a.ParseArgs([]string{}))
		assert.Empty(t, a.ExtraArgs())
	})
	t.Run("too many", func(t *testing.T) {
		a := gflag.Arguments{}
		a.AddArg("file", "input", true)
		assert.NoErr(t, a.ParseArgs([]string{"first", "old"}))
		a.SetValidateNum(true)
		assert.Err(t, a.ParseArgs([]string{"second", "new"}))
		assert.Empty(t, a.ExtraArgs())
	})
	t.Run("new extras retained", func(t *testing.T) {
		a := gflag.Arguments{}
		a.AddArg("file", "input", true)
		assert.NoErr(t, a.ParseArgs([]string{"first", "old"}))
		assert.NoErr(t, a.ParseArgs([]string{"second", "new"}))
		assert.Eq(t, []string{"new"}, a.ExtraArgs())
	})
}
