package gflag_test

import (
	"testing"

	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/x/assert"
)

func TestFlags_Parse_reorderArgs(t *testing.T) {
	t.Run("enabled by default", func(t *testing.T) {
		var name string
		var verbose bool
		fs := gflag.New("test")
		fs.StrOpt(&name, "name", "n", "", "name opt")
		fs.BoolOpt(&verbose, "verbose", "v", false, "verbose opt")

		// options written after arguments are still parsed
		err := fs.Parse([]string{"arg1", "--name", "tom", "arg2", "-v"})
		assert.NoErr(t, err)
		assert.Eq(t, "tom", name)
		assert.True(t, verbose)
		assert.Eq(t, []string{"arg1", "arg2"}, fs.RawArgs())
	})

	t.Run("disabled keeps strict behavior", func(t *testing.T) {
		var name string
		fs := gflag.New("test")
		fs.WithConfigFn(gflag.WithReorderArgs(false))
		fs.StrOpt(&name, "name", "n", "", "name opt")

		// strict: parsing stops at first positional, option after arg is lost
		err := fs.Parse([]string{"arg1", "--name", "tom"})
		assert.NoErr(t, err)
		assert.Eq(t, "", name)
		assert.Eq(t, []string{"arg1", "--name", "tom"}, fs.RawArgs())
	})
}
