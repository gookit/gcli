package gcli_test

import (
	"runtime"
	"testing"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/testutil/assert"
)

func TestGcliBasic(t *testing.T) {
	is := assert.New(t)
	is.NotEmpty(gcli.Version())
	is.NotEmpty(gcli.CommitID())
}

func TestVerbose(t *testing.T) {
	is := assert.New(t)

	old := gcli.Verbose()
	is.Eq(gcli.VerbError, old)
	is.False(gcli.GOpts().NoColor)

	gcli.SetDebugMode()
	is.Eq(gcli.VerbDebug, gcli.Verbose())

	gcli.SetQuietMode()
	is.Eq(gcli.VerbQuiet, gcli.Verbose())

	gcli.SetVerbose(gcli.VerbInfo)
	is.Eq(gcli.VerbInfo, gcli.Verbose())
	is.Eq("info", gcli.Verbose().Name())
	is.Eq("INFO", gcli.Verbose().Upper())

	gcli.SetVerbose(old)
	is.Eq(gcli.VerbError, gcli.Verbose())

	verb := gcli.VerbLevel(23)
	is.Eq("unknown", verb.Name())
	err := verb.Set("2")
	is.NoErr(err)
	is.Eq(gcli.VerbWarn, verb)
	is.Eq("warn", verb.Name())

	err = verb.Set("debug")
	is.NoErr(err)
	is.Eq(gcli.VerbDebug, verb)
	is.Eq("debug", verb.Name())

	err = verb.Set("30")
	is.NoErr(err)
	is.Eq(gcli.VerbCrazy, verb)
	is.Eq("crazy", verb.Name())
}

func TestStrictMode(t *testing.T) {
	is := assert.New(t)

	is.True(gcli.StrictMode())

	gcli.SetStrictMode(false)
	is.False(gcli.StrictMode())

	gcli.SetStrictMode(true)
	is.True(gcli.StrictMode())
}

func TestCmdLine(t *testing.T) {
	is := assert.New(t)

	is.True(gcli.CLI.PID() > 0)
	is.Eq(runtime.GOOS, gcli.CLI.OsName())
	is.NotEmpty(gcli.CLI.BinName())
	is.NotEmpty(gcli.CLI.WorkDir())

	args := gcli.CLI.OsArgs()
	is.NotEmpty(args)

	if len(args) > 1 {
		is.NotEmpty(gcli.CLI.ArgLine())
	} else {
		is.Empty(gcli.CLI.ArgLine())
	}
}

func TestSetStrictMode(t *testing.T) {
	stm := gcli.StrictMode()
	defer gcli.SetStrictMode(stm)

	opts := struct {
		name   string
		ok, bl bool
	}{}

	// gcli.SetVerbose(gcli.VerbDebug)
	app := gcli.NewApp(gcli.NotExitOnEnd())
	app.Add(&gcli.Command{
		Name: "test",
		Config: func(c *gcli.Command) {
			c.StrOpt(&opts.name, "name", "n", "", "1")
			c.BoolOpt(&opts.ok, "ok", "o", true, "2")
			c.BoolOpt(&opts.bl, "bl", "b", false, "3")
		},
		Func: func(c *gcli.Command, _ []string) error {
			return nil
		},
	})

	app.Run([]string{"test", "-o", "-n", "inhere"})
	assert.Eq(t, "inhere", opts.name)
	assert.True(t, opts.ok)

	app.Run([]string{"test", "-o=false", "-n=inhere"})
	assert.Eq(t, "inhere", opts.name)
	assert.False(t, opts.ok)

	app.Run([]string{"test", "-ob"})
	// assert.StrContains(t, errMsg, "ddd")

	gcli.SetStrictMode(true)
	app.Run([]string{"test", "-ob"})
	assert.True(t, opts.ok)
	assert.True(t, opts.bl)

}
