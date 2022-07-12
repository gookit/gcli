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
