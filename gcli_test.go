package gcli_test

import (
	"runtime"
	"testing"

	"github.com/gookit/gcli/v2"
	"github.com/stretchr/testify/assert"
)

func TestVerbose(t *testing.T) {
	is := assert.New(t)

	old := gcli.Verbose()
	is.Equal(gcli.VerbQuiet, old)

	gcli.SetDebugMode()
	is.Equal(gcli.VerbDebug, gcli.Verbose())

	gcli.SetQuietMode()
	is.Equal(gcli.VerbQuiet, gcli.Verbose())

	gcli.SetVerbose(gcli.VerbInfo)
	is.Equal(gcli.VerbInfo, gcli.Verbose())

	gcli.SetVerbose(old)
	is.Equal(gcli.VerbQuiet, gcli.Verbose())
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
	is.NotEmpty(gcli.CLI.OsArgs())
	is.NotEmpty(gcli.CLI.WorkDir())
	is.NotEmpty(gcli.CLI.ArgLine())
}
