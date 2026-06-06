package gflag_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/x/assert"
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

func TestCliOpts_category(t *testing.T) {
	is := assert.New(t)
	co := newFlagOptions()

	var s string
	var b bool
	co.StrVar(&s, &gflag.CliOpt{Name: "host", Desc: "bind host"})
	co.StrVar(&s, &gflag.CliOpt{Name: "port", Desc: "bind port", Category: "network"})
	co.BoolVar(&b, &gflag.CliOpt{Name: "verbose", Desc: "verbose log"})
	co.StrVar(&s, &gflag.CliOpt{Name: "db-dsn", Desc: "database dsn", Category: "database"})

	cats := co.OptCategories()
	is.Len(cats, 3)
	// insertion order: ""(default) -> network -> database
	is.Eq("", cats[0].Name)
	is.Eq("network", cats[1].Name)
	is.Eq("database", cats[2].Name)

	is.Eq([]string{"host", "verbose"}, cats[0].OptNames)
	is.Eq([]string{"port"}, cats[1].OptNames)
	is.Eq([]string{"db-dsn"}, cats[2].OptNames)
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

func TestCliOpt_Question(t *testing.T) {
	oldIn, oldOut := cliutil.Input, cliutil.Output
	defer func() { cliutil.Input, cliutil.Output = oldIn, oldOut }()
	cliutil.Output = io.Discard // silence question prompt in test output

	t.Run("collect_on_empty", func(t *testing.T) {
		cliutil.Input = strings.NewReader("tom\n")

		opt := gflag.NewOpt("opt1", "a option", "", gflag.WithQuestion("your name? "))
		fo := newFlagOptions()
		var str string
		fo.StrVar(&str, opt)

		assert.NoErr(t, fo.ParseOpts(nil))
		assert.Eq(t, "tom", str)
		assert.Eq(t, "tom", opt.Value().String())
	})

	t.Run("collector_priority_over_question", func(t *testing.T) {
		cliutil.Input = strings.NewReader("byQuestion\n")

		opt := gflag.NewOpt("opt1", "a option", "",
			gflag.WithCollector(func() (string, error) { return "byCollector", nil }),
			gflag.WithQuestion("ignored? "),
		)
		fo := newFlagOptions()
		var str string
		fo.StrVar(&str, opt)

		assert.NoErr(t, fo.ParseOpts(nil))
		assert.Eq(t, "byCollector", str)
	})

	t.Run("no_ask_when_not_empty", func(t *testing.T) {
		cliutil.Input = strings.NewReader("shouldNotUse\n")

		opt := gflag.NewOpt("opt1", "a option", "", gflag.WithQuestion("q? "))
		fo := newFlagOptions()
		var str string
		fo.StrVar(&str, opt)

		assert.NoErr(t, fo.ParseOpts([]string{"--opt1", "given"}))
		assert.Eq(t, "given", str)
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

func TestCliOpt_WithChoices(t *testing.T) {
	opt := gflag.NewOpt("format", "output format", "", gflag.WithChoices("json", "yaml", "text"))
	assert.Eq(t, []string{"json", "yaml", "text"}, opt.Choices)
}

func TestCliOpt_TakesValue(t *testing.T) {
	fo := newFlagOptions()

	var s string
	var b bool
	fo.StrVar(&s, &gflag.CliOpt{Name: "str"})
	fo.BoolVar(&b, &gflag.CliOpt{Name: "bl"})

	// 取值型选项 TakesValue 为 true; bool 选项为 false
	assert.True(t, fo.Opt("str").TakesValue())
	assert.False(t, fo.Opt("bl").TakesValue())
}
