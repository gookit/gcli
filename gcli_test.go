package gcli_test

import (
	"runtime"
	"testing"

	"github.com/gookit/gcli/v3"
	"github.com/stretchr/testify/assert"
)

func TestVerbose(t *testing.T) {
	is := assert.New(t)

	old := gcli.Verbose()
	is.Equal(gcli.VerbError, old)
	is.False(gcli.GOpts().NoColor)

	gcli.SetDebugMode()
	is.Equal(gcli.VerbDebug, gcli.Verbose())

	gcli.SetQuietMode()
	is.Equal(gcli.VerbQuiet, gcli.Verbose())

	gcli.SetVerbose(gcli.VerbInfo)
	is.Equal(gcli.VerbInfo, gcli.Verbose())
	is.Equal("info", gcli.Verbose().Name())
	is.Equal("INFO", gcli.Verbose().Upper())

	gcli.SetVerbose(old)
	is.Equal(gcli.VerbError, gcli.Verbose())

	verb := gcli.VerbLevel(23)
	is.Equal("unknown", verb.Name())
	err := verb.Set("2")
	is.NoError(err)
	is.Equal(gcli.VerbWarn, verb)
	is.Equal("warn", verb.Name())

	err = verb.Set("debug")
	is.NoError(err)
	is.Equal(gcli.VerbDebug, verb)
	is.Equal("debug", verb.Name())

	err = verb.Set("30")
	is.NoError(err)
	is.Equal(gcli.VerbCrazy, verb)
	is.Equal("crazy", verb.Name())
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
	is.Equal(runtime.GOOS, gcli.CLI.OsName())
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
