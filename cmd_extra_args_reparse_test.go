package gcli_test

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/x/assert"
	"testing"
)

func TestCommandRunDoesNotReuseExtraArgs(t *testing.T) {
	defer gcli.ResetGOpts()
	var calls [][]string
	c := &gcli.Command{Name: "repeat", Func: func(_ *gcli.Command, args []string) error {
		calls = append(calls, append([]string{}, args...))
		return nil
	}}
	c.AddArg("file", "input", true)
	assert.NoErr(t, c.Run([]string{"first", "old"}))
	assert.NoErr(t, c.Run([]string{"second"}))
	assert.NoErr(t, c.Run([]string{"third", "new"}))
	assert.Eq(t, [][]string{{"old"}, {}, {"new"}}, calls)
}
