package gflag_test

import (
	"testing"
	"time"

	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/x/assert"
)

func TestGeneric_Opt(t *testing.T) {
	var name string
	var age int
	var verbose bool
	var tags []string
	var ttl time.Duration
	var meta map[string]string

	fs := gflag.New("test")
	gflag.Opt(fs, &name, "name", "n", "tom", "user name")
	gflag.Opt(fs, &age, "age", "a", 18, "user age")
	gflag.Opt(fs, &verbose, "verbose", "v", false, "verbose")
	gflag.Opt(fs, &tags, "tag", "t", nil, "tags, repeatable")
	gflag.Opt(fs, &ttl, "ttl", "", time.Duration(0), "time to live")
	gflag.Opt(fs, &meta, "meta", "m", nil, "metadata, repeatable")

	// scalar defaults are applied at bind time
	assert.Eq(t, "tom", name)
	assert.Eq(t, 18, age)

	err := fs.Parse([]string{
		"--name", "jerry", "-a", "30", "-v",
		"-t", "x", "-t", "y",
		"--ttl", "1h30m", "-m", "k=v",
	})
	assert.NoErr(t, err)
	assert.Eq(t, "jerry", name)
	assert.Eq(t, 30, age)
	assert.True(t, verbose)
	assert.Eq(t, []string{"x", "y"}, tags)
	assert.Eq(t, 90*time.Minute, ttl)
	assert.Eq(t, "v", meta["k"])
}

// BindVar accepts a custom flag.Value pointer via the fallback branch
func TestGeneric_BindVar_flagValue(t *testing.T) {
	var langs gflag.Strings // *gflag.Strings implements flag.Value
	fs := gflag.New("test")
	gflag.BindVar(fs, &langs, gflag.NewOpt("langs", "lang list", nil))

	assert.NoErr(t, fs.Parse([]string{"--langs", "go", "--langs", "php"}))
	assert.Eq(t, "go,php", langs.String())
}

// unsupported type panics with a clear message
func TestGeneric_BindVar_unsupported(t *testing.T) {
	var f float32
	fs := gflag.New("test")
	assert.PanicsMsg(t, func() {
		gflag.BindVar(fs, &f, gflag.NewOpt("ff", "ff desc", nil))
	}, `gflag: BindVar: unsupported type *float32 for option "ff"`)
}
