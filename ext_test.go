package gcli_test

import (
	"testing"

	"github.com/gookit/gcli/v3"
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
