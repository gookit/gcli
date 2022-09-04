package helper_test

import (
	"testing"

	"github.com/gookit/gcli/v3/helper"
	"github.com/gookit/goutil/testutil/assert"
)

func TestHelpVars(t *testing.T) {
	is := assert.New(t)
	vs := helper.HelpVars{
		Vars: map[string]string{
			"key0": "val0",
			"key1": "val1",
		},
	}

	is.Len(vs.GetVars(), 2)
	is.Contains(vs.GetVars(), "key0")

	vs.AddVars(map[string]string{"key2": "val2"})
	vs.AddVar("key3", "val3")

	is.Eq("val3", vs.GetVar("key3"))
	is.Eq("", vs.GetVar("not-exist"))

	is.Eq("hello val0", vs.ReplaceVars("hello {$key0}"))
	is.Eq("hello val0 val2", vs.ReplaceVars("hello {$key0} {$key2}"))
	// invalid input
	is.Eq("hello {key0}", vs.ReplaceVars("hello {key0}"))
}
