package gflag_test

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"testing"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/testutil"
	"github.com/gookit/goutil/x/assert"
)

func TestFlags_Basic(t *testing.T) {
	fs := gflag.New("testFlags")

	assert.Len(t, fs.Opts(), 0)
	assert.Eq(t, 0, fs.Len())
	assert.Eq(t, "testFlags", fs.Name())

	assert.Nil(t, fs.LookupFlag("opt1"))
	assert.Len(t, fs.ShortNames("opt"), 0)
	assert.False(t, fs.HasOption("opt1"))

	var s1, s2 string
	fs.StrOpt(&s1, "str1", "", "", "desc")

	assert.True(t, fs.IsOption("str1"))
	assert.False(t, fs.IsShortOpt("str1"))
	assert.False(t, fs.IsShortName("str1"))
	assert.False(t, fs.IsOption("not-exist"))

	fs.StrOpt(&s2, "str2", "b", "", "desc")
	assert.True(t, fs.IsShortName("b"))

	buf := new(bytes.Buffer)
	fs.IterAll(func(f *gflag.Flag, meta *gcli.CliOpt) {
		_, _ = fmt.Fprintf(buf, "flag: %s, shorts: %s;", f.Name, meta.Shorts2String(","))
	})

	assert.Eq(t, "flag: str1, shorts: ;flag: str2, shorts: b;", buf.String())
}

func TestFlags_BoolOpt(t *testing.T) {
	fs := gcli.NewFlags("testFlags")

	var b1, b2 bool
	b0 := fs.Bool("bl0", "", false, "desc0")
	fs.BoolOpt(&b1, "bl1", "a,b", false, "desc1")
	fs.BoolVar(&b2, &gcli.CliOpt{
		Name: "bl2",
		Desc: "desc2",
	})

	assert.False(t, *b0)
	assert.False(t, b1)
	assert.NoErr(t, fs.Parse([]string{"--bl0", "-a", "--bl2"}))
	assert.True(t, *b0)
	assert.True(t, b1)
}

func TestFlags_StrOpt(t *testing.T) {
	fs := gcli.NewFlags("testFlags")
	assert.Len(t, fs.Opts(), 0)

	var str string
	fs.StrVar(&str, &gcli.CliOpt{
		Name: "test",
		Desc: "test desc",
	})

	assert.True(t, fs.IsOption("test"))
	assert.False(t, fs.IsOption("not-exist"))
	assert.Len(t, fs.Opts(), 1)

	f := fs.LookupFlag("test")
	assert.NotEmpty(t, f)

	assert.Eq(t, "test", f.Name)
	assert.Eq(t, "test desc", f.Usage)

	ns := fs.FlagNames()
	assert.Len(t, ns, 1)

	f = fs.LookupFlag("not-exist")
	assert.Nil(t, f)

	err := fs.Parse([]string{})
	assert.NoErr(t, err)
	assert.Eq(t, "", str)

	err = fs.Parse([]string{"--test", "value"})
	assert.NoErr(t, err)
	assert.Eq(t, "value", str)
	assert.Len(t, fs.ShortNames("test"), 0)
}

func TestFlags_Float64Opt(t *testing.T) {
	fs := gcli.NewFlags("testFlags")

	var f1, f2 float64
	fs.Float64Opt(&f1, "f1", "a", 0, "desc1")
	fs.Float64Var(&f2, &gcli.CliOpt{
		Name:   "f2",
		Desc:   "desc2",
		DefVal: 3.14,
	})

	assert.Eq(t, float64(0), f1)
	assert.Eq(t, 3.14, f2)
	assert.NoErr(t, fs.Parse([]string{"-a", "12.3", "--f2", "1.63"}))
	assert.Eq(t, 12.3, f1)
	assert.Eq(t, 1.63, f2)
}

func TestFlags_IntOpt(t *testing.T) {
	fs := gcli.NewFlags("testFlags")

	var int1, int2 int
	fs.IntOpt(&int1, "int1", "a,b", 0, "desc1")
	fs.IntVar(&int2, &gcli.CliOpt{
		Name:   "int2",
		Desc:   "desc2",
		DefVal: 314,
	})

	assert.Eq(t, 0, int1)
	assert.Eq(t, 314, int2)
	assert.NoErr(t, fs.Parse([]string{"-a", "123", "--int2", "163"}))
	assert.Eq(t, 123, int1)
	assert.Eq(t, 163, int2)
}

func TestFlags_Int64Opt(t *testing.T) {
	fs := gcli.NewFlags("testFlags")

	var int1, int2 int64
	fs.Int64Opt(&int1, "int1", "a,b", 0, "desc1")
	fs.Int64Var(&int2, &gcli.CliOpt{
		Name:   "int2",
		Desc:   "desc2",
		DefVal: 314,
	})

	assert.Eq(t, int64(0), int1)
	assert.Eq(t, int64(314), int2)
	assert.NoErr(t, fs.Parse([]string{"-a", "12", "--int2", "16"}))
	assert.Eq(t, int64(12), int1)
	assert.Eq(t, int64(16), int2)
}

func TestFlags_UintOpt(t *testing.T) {
	fs := gcli.NewFlags("testFlags")

	var int1, int2 uint
	fs.UintOpt(&int1, "int1", "a", 0, "desc1")
	fs.UintVar(&int2, &gcli.CliOpt{
		Name:   "c",
		Desc:   "desc2",
		DefVal: 314,
	})

	assert.Eq(t, uint(0), int1)
	assert.Eq(t, uint(314), int2)
	assert.NoErr(t, fs.Parse([]string{"-a", "12", "-c", "16"}))
	assert.Eq(t, uint(12), int1)
	assert.Eq(t, uint(16), int2)
}

func TestFlags_Uint64Opt(t *testing.T) {
	fs := gcli.NewFlags("testFlags")

	var uint1, uint2 uint64
	fs.Uint64Opt(&uint1, "uint1", "a", 0, "desc1")
	fs.Uint64Var(&uint2, &gcli.CliOpt{
		Name:   "uint2",
		Desc:   "desc2",
		DefVal: 314,
		Shorts: []string{"c", "", "f"},
	})

	fm2 := fs.Opt("uint2")
	assert.Len(t, fm2.Shorts, 2)
	assert.Eq(t, "c,f", fm2.Shorts2String())

	assert.Eq(t, uint64(0), uint1)
	assert.Eq(t, uint64(314), uint2)
	assert.NoErr(t, fs.Parse([]string{"-a", "12", "--uint2", "16"}))
	assert.Eq(t, uint64(12), uint1)
	assert.Eq(t, uint64(16), uint2)
}

func TestFlags_FuncOpt(t *testing.T) {
	fs := gcli.NewFlags("testFlags")

	var str string
	fs.FuncOpt("str, s", "desc", func(v string) error {
		str = v
		return nil
	})

	assert.Eq(t, "", str)
	assert.NoErr(t, fs.Parse([]string{"-s", "abc"}))
	assert.Eq(t, "abc", str)
}

func TestFlags_VarOpt(t *testing.T) {
	fs := gcli.NewFlags("testFlags")

	var ints gcli.Ints
	fs.Var(&ints, &gcli.CliOpt{Name: "ints", Desc: "desc"})
	assert.NoErr(t, fs.Parse([]string{"--ints", "12", "--ints", "16"}))

	assert.Len(t, ints, 2)
	assert.Eq(t, "[12,16]", ints.String())

	var ss gcli.Strings
	fs.VarOpt(&ss, "names", "n,s", "desc")
	assert.NoErr(t, fs.Parse([]string{"--names", "abc", "-n", "def", "-s", "ghi"}))

	assert.Len(t, ss, 3)
	assert.Eq(t, "abc,def,ghi", ss.String())
}

func TestFlags_CheckName(t *testing.T) {
	// 重复绑定 panic 带上 flags/命令名(默认名 gflag), 便于定位
	assert.PanicsMsg(t, func() {
		var i int64
		fs := gcli.NewFlags()
		fs.Int64Opt(&i, "opt1", "", 0, "desc")
		fs.Int64Opt(&i, "opt1", "", 0, "desc")
	}, "gflag: redefined option flag 'opt1' in command 'gflag'")

	// 命名 flags 时, 重复绑定 panic 会带上该名称
	assert.PanicsMsg(t, func() {
		var i int64
		fs := gcli.NewFlags("mycmd")
		fs.Int64Opt(&i, "opt1", "", 0, "desc")
		fs.Int64Opt(&i, "opt1", "", 0, "desc")
	}, "gflag: redefined option flag 'opt1' in command 'mycmd'")

	assert.PanicsMsg(t, func() {
		var b bool
		fs := gcli.NewFlags()
		fs.BoolOpt(&b, "", "", false, "desc")
	}, "gflag: option flag name cannot be empty")

	assert.PanicsMsg(t, func() {
		var fv uint
		fs := gcli.NewFlags()
		fs.UintOpt(&fv, "+invalid", "", 0, "desc")
	}, "gflag: option flag name '+invalid' is invalid, must match: ^[0-9a-zA-Z][\\w-]*$")

	assert.PanicsMsg(t, func() {
		var fv uint64
		fs := gcli.NewFlags()
		fs.Uint64Opt(&fv, "78", "", 0, "desc")
	}, "gflag: option flag name '78' is not a good name")
}

func TestFlags_CheckShorts(t *testing.T) {
	assert.NotPanics(t, func() {
		var fv float64
		fs := gcli.NewFlags()

		// "+" has been filtered by func: splitShortcut()
		fs.Float64Opt(&fv, "float", "+", 0, "desc")

		fm := fs.Opt("float")
		assert.Len(t, fm.Shorts, 0)
	})

	var fv float64
	fs := gcli.NewFlags()
	fs.Float64Var(&fv, &gflag.CliOpt{
		Name:   "float",
		Shorts: []string{"+", "-"},
	})
	fm := fs.Opt("float")
	assert.Empty(t, fm.Shorts)

	assert.PanicsMsg(t, func() {
		var i int
		fs := gcli.NewFlags()
		fs.IntVar(&i, &gcli.CliOpt{Name: "a", Shorts: []string{"a"}})
	}, "gflag: short name 'a' has been used as the current option name")

	assert.PanicsMsg(t, func() {
		var i int
		fs := gcli.NewFlags()
		fs.IntOpt(&i, "s", "", 0, "desc")
		fs.IntOpt(&i, "int1", "s", 0, "desc")
	}, "gflag: short name 's' has been used as an option name")

	assert.PanicsMsg(t, func() {
		var str string
		fs := gcli.NewFlags()
		fs.StrOpt(&str, "str", "s", "", "desc")
		fs.StrOpt(&str, "str1", "s", "", "desc")
	}, "gflag: short name 's' has been used by option 'str'")
}

var flagOpts = struct {
	intv int
	strv string
}{}

func TestFlags_Run(t *testing.T) {
	is := assert.New(t)

	fg := gcli.NewFlags("test", "desc message")
	// fg.ExitFunc = func(code int) {}

	fg.IntOpt(&flagOpts.intv, "intv", "i", 0, "desc message for intv")
	fg.StrOpt(&flagOpts.strv, "strv", "s", "", "desc message for strv")

	// parse
	fg.Run([]string{"./app", "-i", "23", "-s", "inhere"})
	is.Eq(23, flagOpts.intv)
	is.Eq("inhere", flagOpts.strv)

	// help
	fg.Run([]string{"./app", "-h"})
}

func TestFlags_Parse(t *testing.T) {
	var str string

	gf := gcli.NewFlags("test")
	gf.StrVar(&str, &gcli.CliOpt{
		Name:     "opt1",
		Required: true,
		Validator: func(val string) error {
			if len(val) < 5 {
				return errors.New("flag value min len is 5")
			}

			return nil
		},
	})

	err := gf.Parse([]string{})
	assert.Err(t, err)

	err = gf.Parse([]string{"--opt1", ""})
	assert.Err(t, err)

	err = gf.Parse([]string{"--opt1", "val"})
	assert.Err(t, err)
	assert.ErrSubMsg(t, err, "flag value min len is 5")

	err = gf.Parse([]string{"--opt1", "value"})
	assert.NoErr(t, err)
	assert.Eq(t, "value", str)
}

func TestFlags_Int_bindingNilPtr(t *testing.T) {
	type userOpts struct {
		Int *int
	}

	opt := userOpts{}
	dump.P(opt)

	// must init for an ptr value.
	assert.Panics(t, func() {
		fs := gcli.NewFlags("test")
		fs.IntOpt(opt.Int, "int", "i", 20, "")
	})

	aint := 23
	opt = userOpts{Int: &aint}
	dump.P(opt)
	fs := gcli.NewFlags("test")
	fs.IntOpt(opt.Int, "int", "i", 20, "")
	fs.PrintHelpPanel()
}

func TestFlags_FromStruct_simple(t *testing.T) {
	fs := gcli.NewFlags("test")

	type userOpts0 struct {
		Int int    `flag:"name=int0;shorts=i;required=true;desc=int option message"`
		Bol bool   `flag:"name=bol;shorts=b;default=true;desc=bool option message"`
		Str string `flag:"name=str1;shorts=o,h;required=true;desc=str1 message"`
	}

	opt := &userOpts0{}
	assert.False(t, opt.Bol)
	assert.Eq(t, 0, opt.Int)
	assert.Eq(t, "", opt.Str)

	err := fs.FromStruct(opt)
	assert.NoErr(t, err)
	assert.True(t, opt.Bol)
	assert.True(t, fs.HasOption("bol"))

	fs.PrintHelpPanel()

	err = fs.Parse([]string{"--int0", "13", "--str1", "xyz"})
	assert.NoErr(t, err)

	assert.Eq(t, 13, opt.Int)
	assert.Eq(t, "xyz", opt.Str)

	// not use ptr
	opts1 := userOpts0{}
	err = fs.FromStruct(opts1)
	assert.Err(t, err)
}

func TestFlags_FromStruct_ptrField(t *testing.T) {
	type userOpts struct {
		Int *int    `flag:"shorts=i;default=13;desc=int option message"`
		Str *string `flag:"name=str2;required=true;desc=str2 message"`
	}

	opt := &userOpts{}
	dump.P(opt)

	fs := gcli.NewFlags("test1")
	opt = &userOpts{}
	err := fs.FromStruct(opt)
	assert.Err(t, err)
	assert.Eq(t, "field: Int - nil pointer dereference", err.Error())

	aint := 23
	astr := "xyz"
	opt = &userOpts{Int: &aint, Str: &astr}
	dump.P(opt)
	assert.Eq(t, 23, *opt.Int)
	assert.Eq(t, "xyz", *opt.Str)

	fs = gcli.NewFlags("test1")
	err = fs.FromStruct(opt)
	assert.NoErr(t, err)
	assert.Eq(t, 13, *(opt.Int))
	// assert.Eq(t, "", *opt.Str)

	dump.P(opt)
	fmt.Println("Flag Help:")
	fs.PrintHelpPanel()
}

func TestFlags_FromStruct_fieldRule(t *testing.T) {
	type fieldOpts struct {
		UserName string `flag:"n,u" desc:"input user name" required:"true"`
		Age      int    `desc:"input age" default:"18"`
		Verbose  bool   `flag:"V" desc:"verbose mode"`
		Ignored  string `flag:"-"`
		NoTag    string
	}

	fs := gcli.NewFlags("test")
	opt := &fieldOpts{}
	err := fs.FromStruct(opt, gflag.TagRuleField)
	assert.NoErr(t, err)

	assert.True(t, fs.HasOption("user-name"))
	assert.True(t, fs.HasOption("age"))
	assert.True(t, fs.HasOption("verbose"))
	assert.False(t, fs.HasOption("ignored"))
	assert.False(t, fs.HasOption("no-tag"))

	err = fs.Parse([]string{"-u", "tom", "-V"})
	assert.NoErr(t, err)
	assert.Eq(t, "tom", opt.UserName)
	assert.True(t, opt.Verbose)
	assert.Eq(t, 18, opt.Age) // default value take effect

	// verify the other short `n` also works, and gcli.TagRuleField re-export is usable
	fs2 := gcli.NewFlags("test2")
	opt2 := &fieldOpts{}
	assert.NoErr(t, fs2.FromStruct(opt2, gcli.TagRuleField))
	assert.NoErr(t, fs2.Parse([]string{"-n", "jody"}))
	assert.Eq(t, "jody", opt2.UserName)
}

func TestFlags_FromStruct_anonymous(t *testing.T) {
	type CommonOpts struct {
		Verbose bool `flag:"name=verbose;shorts=V;desc=verbose mode"`
	}

	t.Run("anonymous struct", func(t *testing.T) {
		type cmdOpts struct {
			CommonOpts
			Name string `flag:"name=name;shorts=n;desc=name message"`
		}

		fs := gcli.NewFlags("test")
		opt := &cmdOpts{}
		err := fs.FromStruct(opt)
		assert.NoErr(t, err)

		assert.True(t, fs.HasOption("verbose")) // anonymous expanded
		assert.True(t, fs.HasOption("name"))

		err = fs.Parse([]string{"-V", "-n", "tom"})
		assert.NoErr(t, err)
		assert.True(t, opt.Verbose)
		assert.Eq(t, "tom", opt.Name)
	})

	t.Run("anonymous pointer struct", func(t *testing.T) {
		type cmdOpts struct {
			*CommonOpts
			Name string `flag:"name=name;shorts=n;desc=name message"`
		}

		fs := gcli.NewFlags("test")
		opt := &cmdOpts{CommonOpts: &CommonOpts{}}
		err := fs.FromStruct(opt)
		assert.NoErr(t, err)

		assert.True(t, fs.HasOption("verbose")) // anonymous ptr expanded
		assert.True(t, fs.HasOption("name"))

		err = fs.Parse([]string{"-V", "-n", "tom"})
		assert.NoErr(t, err)
		assert.True(t, opt.Verbose)
		assert.Eq(t, "tom", opt.Name)
	})
}

func TestFlags_FromStruct_noNameStruct(t *testing.T) {
	logOpts := struct {
		Abbrev    bool        `flag:"Only display the abbrev commit ID"`
		NoColor   bool        `flag:"Dont use color render git output"`
		MaxCommit int         `flag:"Max display how many commits;;15"`
		Logfile   string      `flag:"export changelog message to file"`
		Exclude   gcli.String `flag:"exclude contains given sub-string. multi by comma split."`
	}{}

	fs := gcli.NewFlags("test")
	fs.UseSimpleRule()
	// err := fs.FromStruct(logOpts)
	err := fs.FromStruct(&logOpts)

	assert.NoErr(t, err)
	assert.True(t, fs.HasOption("abbrev"))
	assert.True(t, fs.IsOption("abbrev"))
}

func TestExtType_Strings(t *testing.T) {
	var v1 any
	v1 = gcli.Strings{}
	val, ok := v1.(flag.Value)
	assert.False(t, ok)
	assert.Nil(t, val)

	// NOTE: must use ptr
	v1 = &gcli.Strings{}
	val, ok = v1.(flag.Value)
	assert.True(t, ok)
	assert.NotNil(t, val)
}

func TestFlags_FromStruct_var_Strings(t *testing.T) {
	type fileReplaceOpt struct {
		// Dir   string       `flag:"desc=the directory for find and replace"`
		Files gcli.Strings `flag:"desc=the files want replace content;shorts=f"`
	}

	opt := fileReplaceOpt{Files: make(gcli.Strings, 0)}

	fs := gcli.NewFlags("test")
	err := fs.FromStruct(&opt)

	assert.NoErr(t, err)
	assert.True(t, fs.HasOption("files"))

	err = fs.Parse([]string{"--files", "a.txt", "-f", "b.txt"})
	assert.NoError(t, err)
	assert.Len(t, opt.Files, 2)
	assert.Eq(t, "a.txt,b.txt", opt.Files.String())
}

func TestFlags_FromStruct(t *testing.T) {
	type userOpts struct {
		Int  int    `flag:"name=int0;shorts=i;required=true;desc=int option message"`
		Bol  bool   `flag:"name=bol;shorts=b;desc=bool option message"`
		Str1 string `flag:"name=str1;shorts=o,h;required=true;desc=str1 message"`
		// use ptr
		Str2 *string `flag:"name=str2;required=true;desc=str2 message"`
		// custom type and implement flag.Value
		Verb0 gcli.VerbLevel `flag:"name=verb0;shorts=V;desc=verb0 message"`
		// use ptr
		Verb1 *gcli.VerbLevel `flag:"name=verb1;desc=verb1 message"`
	}

	astr := "xyz"
	verb := gcli.VerbWarn
	fs := gcli.NewFlags("test")
	err := fs.FromStruct(&userOpts{
		Str2:  &astr,
		Verb1: &verb,
	})
	assert.NoErr(t, err)

	help := fs.String()
	assert.Contains(t, help, "-h, -o, --str1")

	fmt.Println("Flag Help:")
	fs.PrintHelpPanel()
}

// func TestFlags_FromText(t *testing.T) {
// }

func TestFlags_PrintHelpPanel(t *testing.T) {
	fs := gcli.NewFlags("test")

	testOpts := struct {
		opt1      int
		opt2      bool
		opt3      string
		optByEnv1 string
		optByEnv2 string
		optByEnv3 int
		optByEnv4 string
	}{
		optByEnv1: "${TEST_OPT_ENV1}",
	}

	// dont set env value for optByEnv4: TEST_OPT_ENV4
	tmpKey := testutil.SetOsEnvs(map[string]string{
		"TEST_OPT_ENV1": "test_value_env1",
		"TEST_OPT_ENV2": "test_value_env2",
		"TEST_OPT_ENV3": "345",
	})
	defer testutil.RemoveTmpEnvs(tmpKey)

	fs.IntVar(&testOpts.opt1, &gcli.CliOpt{Name: "opt1"})
	fs.StrVar(&testOpts.opt3, &gcli.CliOpt{
		Name: "test, t",
		Desc: "test desc",
		// required
		Required: true,
	})
	fs.BoolOpt(&testOpts.opt2, "bol", "ab", false, "opt2 desc")

	// set default from ENV
	fs.StrOpt2(&testOpts.optByEnv1, "optByEnv1", "optByEnv1 desc")
	fs.StrOpt2(&testOpts.optByEnv2, "optByEnv2", "optByEnv2 desc", gflag.WithDefault("${TEST_OPT_ENV2}"))
	fs.StrOpt(&testOpts.optByEnv4, "optByEnv4", "", "${TEST_OPT_ENV4}", "optByEnv4 desc")
	fs.IntVar(&testOpts.optByEnv3, &gcli.CliOpt{
		Name:   "optByEnv3",
		Desc:   "optByEnv3 desc",
		DefVal: "${TEST_OPT_ENV3}",
	})

	fs.PrintHelpPanel()
}

func TestFlags_PrintHelpPanel_IndentLongOpt(t *testing.T) {
	fs := gflag.New("test").WithConfigFn(func(c *gflag.Config) {
		c.IndentLongOpt = true
	})

	testOpts := struct {
		opt1 int
		opt2 bool
		opt3 string
	}{}

	fs.IntVar(&testOpts.opt1, &gcli.CliOpt{Name: "opt1"})
	fs.StrVar(&testOpts.opt3, &gcli.CliOpt{
		Name: "test, t",
		Desc: "test desc",
		// required
		Required: true,
	})
	fs.BoolOpt(&testOpts.opt2, "bol", "ab", false, "opt2 desc")
	fmt.Println("Flag Help - enable IndentLongOpt:")
	fs.PrintHelpPanel()
}

func TestParser_Parse_recoverPanic(t *testing.T) {
	is := assert.New(t)

	var name string
	fs := gflag.New("test")
	fs.StrVar(&name, &gcli.CliOpt{Name: "name"})
	// AfterParse 内部 panic，Parse 必须把它转成 error 返回，而非静默吞掉。
	fs.AfterParse = func(_ *gflag.Flags) error { panic("boom on after parse") }

	err := fs.Parse([]string{"--name", "inhere"})
	is.Err(err)
	is.StrContains(err.Error(), "panic")
}

func TestParser_HelperMethods(t *testing.T) {
	t.Run("config required output and raw args", func(t *testing.T) {
		var name string
		buf := new(bytes.Buffer)
		fs := gflag.New("helpers")
		fs.SetConfig(&gflag.Config{TagName: "flagx", IndentLongOpt: true})
		fs.StrOpt(&name, "name", "n", "", "name option")
		fs.Required("name")
		fs.SetOutput(buf)

		assert.Eq(t, "flagx", fs.ParserCfg().TagName)
		assert.PanicsMsg(t, func() {
			fs.Required("missing")
		}, "gflag: config undefined option '--missing'")

		assert.NoErr(t, fs.Parse([]string{"--name", "inhere", "raw"}))
		assert.Eq(t, "inhere", name)
		assert.Eq(t, "raw", fs.RawArg(0))
		assert.Eq(t, []string{"raw"}, fs.FSetArgs())

		fs.PrintHelpPanel()
		assert.StrContains(t, buf.String(), "--name")
	})

	t.Run("handle and init", func(t *testing.T) {
		called := false
		fs := &gflag.Flags{}
		fs.Init("manual")
		fs.SetHandle(func(p *gflag.Flags) error {
			called = true
			assert.Eq(t, "manual", p.Name())
			return nil
		}).Run([]string{"bin"})

		assert.True(t, called)
	})

	t.Run("reorder stop keeps child args verbatim", func(t *testing.T) {
		var verbose bool
		fs := gflag.New("reorder")
		fs.BoolOpt(&verbose, "verbose", "v", false, "verbose mode")
		fs.SetReorderStop(func(name string) bool {
			return name == "child"
		})

		assert.NoErr(t, fs.Parse([]string{"child", "-v"}))
		assert.False(t, verbose)
		assert.Eq(t, []string{"child", "-v"}, fs.FSetArgs())
	})
}

func TestParser_InheritOptsFrom(t *testing.T) {
	t.Run("inherits option and shares bound value", func(t *testing.T) {
		var token string
		parent := gflag.New("parent")
		parent.StrOpt(&token, "token", "t", "", "api token")

		child := gflag.New("child")
		child.InheritOptsFrom(parent, "Global Options")

		assert.True(t, child.HasOption("token"))
		assert.True(t, child.IsShortName("t"))
		assert.NoErr(t, child.Parse([]string{"--token", "abc"}))
		assert.Eq(t, "abc", token)
		assert.Eq(t, "Global Options", child.Opt("token").Category)
	})

	t.Run("skips local name and short conflicts", func(t *testing.T) {
		var parentToken, childToken string
		parent := gflag.New("parent")
		parent.StrOpt(&parentToken, "token", "t", "", "parent token")

		child := gflag.New("child")
		child.StrOpt(&childToken, "token", "", "", "child token")
		child.InheritOptsFrom(parent)

		assert.Len(t, child.ShortNames("token"), 0)
		assert.NoErr(t, child.Parse([]string{"--token", "local"}))
		assert.Eq(t, "", parentToken)
		assert.Eq(t, "local", childToken)

		var mode string
		conflictParent := gflag.New("conflict-parent")
		conflictParent.StrOpt(&mode, "mode", "m", "", "mode")

		var mine bool
		conflictChild := gflag.New("conflict-child")
		conflictChild.BoolOpt(&mine, "mine", "m", false, "mine")
		conflictChild.InheritOptsFrom(conflictParent)
		assert.False(t, conflictChild.HasOption("mode"))
	})
}
