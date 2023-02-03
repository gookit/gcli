package gflag_test

import (
	"errors"
	"testing"

	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/testutil/assert"
)

func TestCliOpt_basic(t *testing.T) {
	opt := gflag.NewOpt("opt1", "a option", "")
	opt.WithOptFns(gflag.WithDefault("abc"), gflag.WithRequired(), gflag.WithShorts("o", "a"))

	assert.True(t, opt.Required)
	assert.Eq(t, "abc", opt.DefVal)
	assert.Eq(t, "o,a", opt.Shorts2String())

	opt = gflag.NewOpt("opt1", "a option", "")
	opt.WithOptFns(gflag.WithShortcut("o"))

	assert.False(t, opt.Required)
	assert.Eq(t, "o", opt.Shorts2String())
}

func TestCliOpt_Validate(t *testing.T) {
	fm := gflag.CliOpt{Name: "opt1"}
	assert.False(t, fm.Required)

	fm.WithOptFns(gflag.WithRequired(), gflag.WithValidator(func(val string) error {
		if len(val) < 5 {
			return errors.New("flag value min len is 5")
		}
		return nil
	}))

	assert.True(t, fm.Required)

	err := fm.Validate("")
	assert.Err(t, err)
	assert.Eq(t, "flag 'opt1' is required", err.Error())

	err = fm.Validate("val")
	assert.Err(t, err)
	assert.Eq(t, "flag value min len is 5", err.Error())

	err = fm.Validate("value")
	assert.NoErr(t, err)
}
