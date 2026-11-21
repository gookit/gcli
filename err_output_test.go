package gcli_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/byteutil"
	"github.com/gookit/goutil/x/assert"
)

// An error text carrying a percent sign (an IPv6 zone such as [fe80::1%12]:22)
// must be printed verbatim: the error string is data, not a format string.
func TestApp_ErrorTipsKeepsPercentLiteral(t *testing.T) {
	buf := byteutil.NewBuffer()
	color.Disable()
	color.SetOutput(buf)
	defer color.ResetOptions()

	app := gcli.NewApp(gcli.NotExitOnEnd())
	app.Add(gcli.NewCommand("fail", "desc").WithFunc(func(c *gcli.Command, args []string) error {
		return fmt.Errorf("dial tcp [fe80::3b3c:ca9b:837c:e55a%%12]:22: i/o timeout")
	}))

	code := app.Run([]string{"fail"})
	assert.Eq(t, gcli.ERR.ToInt(), code)

	out := buf.String()
	assert.StrContains(t, out, "dial tcp [fe80::3b3c:ca9b:837c:e55a%12]:22: i/o timeout")
	if strings.Contains(out, "%!") {
		t.Fatalf("error text was mangled as a format string, output:\n%s", out)
	}
}

// A command can pick its own process exit code, e.g. 130 for an interruption.
func TestApp_NewExitErrorSetsCode(t *testing.T) {
	t.Run("custom code replaces the default ERR code", func(t *testing.T) {
		buf := byteutil.NewBuffer()
		color.Disable()
		color.SetOutput(buf)
		defer color.ResetOptions()

		sentinel := errors.New("scan interrupted")
		app := gcli.NewApp(gcli.NotExitOnEnd())
		app.Add(gcli.NewCommand("fail", "desc").WithFunc(func(c *gcli.Command, args []string) error {
			return gcli.NewExitError(130, sentinel)
		}))

		assert.Eq(t, 130, app.Run([]string{"fail"}))
		out := buf.String()
		assert.StrContains(t, out, "scan interrupted")
		assert.StrNotContains(t, out, "exit code 130")

		err := app.LastError()
		assert.Err(t, err)
		assert.True(t, errors.Is(err, sentinel))
	})

	t.Run("nil error still exits with the code", func(t *testing.T) {
		app := gcli.NewApp(gcli.NotExitOnEnd())
		app.Add(gcli.NewCommand("fail", "desc").WithFunc(func(c *gcli.Command, args []string) error {
			return gcli.NewExitError(3, nil)
		}))

		assert.Eq(t, 3, app.Run([]string{"fail"}))
	})
}
