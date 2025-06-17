package gflag_test

import (
	"errors"
	"testing"

	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/testutil/assert"
)

func newFlagOptions() gflag.CliOpts {
	co := gflag.CliOpts{}
	co.SetName("opts_test")
	co.InitFlagSet()

	return co
}

func TestCliOpts_useShorts(t *testing.T) {
	co := newFlagOptions()

	var str string
	co.StrOpt(&str, "str", "s", "a string option")

	err := co.ParseOpts([]string{"-s", "val"})
	assert.NoErr(t, err)
	assert.Eq(t, "val", str)
}

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

func TestCliOpt_Collector(t *testing.T) {
	opt := gflag.NewOpt("opt1", "a option", "")
	opt.WithOptFns(gflag.WithCollector(func() (string, error) {
		return "abc", nil
	}))

	fo := newFlagOptions()
	var str string
	fo.StrVar(&str, opt)

	assert.NoErr(t, fo.ParseOpts(nil))
	assert.Eq(t, "abc", str)
	assert.Eq(t, "abc", opt.Value().String())

	// test collect error
	t.Run("error_case", func(t *testing.T) {
		var optInt int
		fo1 := newFlagOptions()
		fo1.IntOpt2(&optInt, "opt2", "a int option", gflag.WithCollector(func() (string, error) {
			return "", errors.New("collect error")
		}))
		assert.ErrMsg(t, fo1.ParseOpts(nil), "collect error")
		assert.Eq(t, 0, optInt)
	})
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

	var handledVal string
	fm.WithOptFns(gflag.WithHandler(func(val string) error {
		handledVal = val
		return nil
	}))

	err := fm.Validate("")
	assert.Err(t, err)
	assert.Eq(t, "option 'opt1' is required", err.Error())

	err = fm.Validate("val")
	assert.Err(t, err)
	assert.Eq(t, "option 'opt1': flag value min len is 5", err.Error())
	assert.Empty(t, handledVal)

	err = fm.Validate("value")
	assert.NoErr(t, err)
	assert.Eq(t, "value", handledVal)
}
