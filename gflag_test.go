package gcli_test

import (
	"testing"

	"github.com/gookit/gcli/v2"
	"github.com/stretchr/testify/assert"
)

func TestFlags_Basic(t *testing.T) {
	fs := gcli.NewFlags("testFlags")

	assert.Len(t, fs.Metas(), 0)
	assert.Equal(t, 0, fs.Len())
	assert.Equal(t, "testFlags", fs.Name())

	assert.Nil(t, fs.LookupFlag("opt1"))
	assert.Len(t, fs.ShortNames("opt"), 0)
	assert.False(t, fs.HasFlag("opt1"))
	assert.False(t, fs.HasOption("opt1"))
}

func TestFlags_BoolOpt(t *testing.T) {
	fs := gcli.NewFlags("testFlags")

	var b1, b2 bool
	fs.BoolOpt(&b1, "bl1", "ab", false, "desc1")
	fs.BoolVar(&b2, gcli.FlagMeta{
		Name: "bl2",
		Desc: "desc2",
	})
}

func TestFlags_StrOpt(t *testing.T) {
	fs := gcli.NewFlags("testFlags")
	assert.Len(t, fs.Metas(), 0)

	var str string
	fs.StrVar(&str, gcli.FlagMeta{
		Name: "test",
		Desc: "test desc",
	})

	assert.True(t, fs.HasFlagMeta("test"))
	assert.False(t, fs.HasFlagMeta("not-exist"))
	assert.Len(t, fs.Metas(), 1)

	f := fs.LookupFlag("test")
	assert.NotEmpty(t, f)

	assert.Equal(t, "test", f.Name)
	assert.Equal(t, "test desc", f.Usage)

	ns := fs.FlagNames()
	assert.Len(t, ns, 1)

	f = fs.LookupFlag("not-exist")
	assert.Nil(t, f)

	err := fs.Parse([]string{})
	assert.NoError(t, err)
	assert.Equal(t, "", str)

	err = fs.Parse([]string{"--test", "value"})
	assert.NoError(t, err)
	assert.Equal(t, "value", str)
	assert.Len(t, fs.ShortNames("test"), 0)
}

func TestFlags_CheckName(t *testing.T) {
	assert.PanicsWithValue(t, "GCli: redefined option flag 'opt1'", func() {
		var i int64
		fs := gcli.NewFlags()
		fs.Int64Opt(&i, "opt1", "", 0, "desc")
		fs.Int64Opt(&i, "opt1", "", 0, "desc")
	})

	assert.PanicsWithValue(t, "GCli: option flag name cannot be empty", func() {
		var b bool
		fs := gcli.NewFlags()
		fs.BoolOpt(&b, "", "", false, "desc")
	})

	assert.PanicsWithValue(t, "GCli: option flag name '+invalid' is invalid, must match: ^[a-zA-Z][\\w-]*$", func() {
		var fv uint
		fs := gcli.NewFlags()
		fs.UintOpt(&fv, "+invalid", "", 0, "desc")
	})

	assert.PanicsWithValue(t, "GCli: option flag name '78' is invalid, must match: ^[a-zA-Z][\\w-]*$", func() {
		var fv uint64
		fs := gcli.NewFlags()
		fs.Uint64Opt(&fv, "78", "", 0, "desc")
	})
}

func TestFlags_CheckShorts(t *testing.T) {
	assert.NotPanics(t, func() {
		var fv float64
		fs := gcli.NewFlags()

		// "+" has been filtered by func: splitShortStr()
		fs.Float64Opt(&fv, "float", "+", 0, "desc")

		fm := fs.FlagMeta("float")
		assert.Len(t, fm.Shorts, 0)
	})

	assert.PanicsWithValue(t, "GCli: short name only allow: A-Za-z given: '+'", func() {
		var fv float64
		fs := gcli.NewFlags()
		fs.Float64Var(&fv, gcli.FlagMeta{
			Name:   "float",
			Shorts: []string{"+"},
		})
	})

	assert.PanicsWithValue(t, "GCli: short name 'a' has been used as the current option name", func() {
		var i int
		fs := gcli.NewFlags()
		fs.IntVar(&i, gcli.FlagMeta{Name: "a", Shorts: []string{"a"}})
	})

	assert.PanicsWithValue(t, "GCli: short name 's' has been used as an option name", func() {
		var i int
		fs := gcli.NewFlags()
		fs.IntOpt(&i, "s", "", 0, "desc")
		fs.IntOpt(&i, "int1", "s", 0, "desc")
	})

	assert.PanicsWithValue(t, "GCli: short name 's' has been used by option 'str'", func() {
		var str string
		fs := gcli.NewFlags()
		fs.StrOpt(&str, "str", "s", "", "desc")
		fs.StrOpt(&str, "str1", "s", "", "desc")
	})
}

func TestFlags_FromStruct(t *testing.T) {
	gf := gcli.NewFlags("test")

	type userOpts struct {
		Opt1 string `gcli:"name=opt;shorts=oh;required=true;desc=message"`
		// the option Opt2
		Opt2 string `gcli:"name=opt2;required=true;desc=message"`
	}

	err := gf.FromStruct(&userOpts{})

	assert.NoError(t, err)
}

func TestFlags_PrintHelpPanel(t *testing.T) {
	fs := gcli.NewFlags("test")

	testOpts := struct {
		opt1 int
		opt2 bool
		opt3 string
	}{}

	fs.StrVar(&testOpts.opt3, gcli.FlagMeta{
		Name: "test",
		Desc: "test desc",
		//
		Required: true,
	})
	fs.BoolOpt(&testOpts.opt2, "bol", "ab", false, "opt2 desc")
	fs.PrintHelpPanel()
}
