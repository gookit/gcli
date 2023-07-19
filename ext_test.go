package gcli_test

import (
	"testing"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/byteutil"
	"github.com/gookit/goutil/testutil/assert"
)

func TestHelpReplacer(t *testing.T) {
	is := assert.New(t)
	vs := gcli.HelpReplacer{}

	vs.AddReplaces(map[string]string{
		"key0": "val0",
		"key1": "val1",
	})

	is.Len(vs.Replaces(), 2)
	is.Contains(vs.Replaces(), "key0")

	vs.AddReplaces(map[string]string{"key2": "val2"})
	vs.AddReplace("key3", "val3")

	is.Eq("val3", vs.GetReplace("key3"))
	is.Eq("", vs.GetReplace("not-exist"))

	is.Eq("hello val0", vs.ReplacePairs("hello {$key0}"))
	is.Eq("hello val0 val2", vs.ReplacePairs("hello {$key0} {$key2}"))
	// invalid input
	is.Eq("hello {key0}", vs.ReplacePairs("hello {key0}"))
}

func TestHooks_Fire(t *testing.T) {
	is := assert.New(t)
	buf := byteutil.NewBuffer()
	hooks := gcli.Hooks{}

	hooks.AddHook("test", func(ctx *gcli.HookCtx) bool {
		buf.WriteString("fire the test hook")
		return false
	})

	hooks.Fire("test", nil)
	is.Eq("fire the test hook", buf.ResetGet())

	hooks.Fire("not-exist", nil)
	hooks.On("*", func(ctx *gcli.HookCtx) bool {
		buf.WriteString("fire the * hook")
		return false
	})

	// add prefix hook
	hooks.On("app.test.*", func(ctx *gcli.HookCtx) bool {
		buf.WriteString("fire the app.test.* hook")
		return false
	})

	hooks.Fire("app.test.init", nil)

	s := buf.ResetGet()
	is.StrContains(s, "fire the app.test.* hook")
	is.StrContains(s, "fire the * hook")
}
