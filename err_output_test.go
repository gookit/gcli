package gcli_test

import (
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
