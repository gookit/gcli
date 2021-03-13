package gcli_test

import (
	"testing"

	"github.com/gookit/gcli/v3"
	"github.com/stretchr/testify/assert"
)

func TestHelpVars(t *testing.T) {
	is := assert.New(t)
	vs := gcli.HelpVars{
		Vars: map[string]string{
			"key0": "val0",
			"key1": "val1",
		},
	}

	is.Len(vs.GetVars(), 2)
	is.Contains(vs.GetVars(), "key0")

	vs.AddVars(map[string]string{"key2": "val2"})
	vs.AddVar("key3", "val3")

	is.Equal("val3", vs.GetVar("key3"))
	is.Equal("", vs.GetVar("not-exist"))

	is.Equal("hello val0", vs.ReplaceVars("hello {$key0}"))
	is.Equal("hello val0 val2", vs.ReplaceVars("hello {$key0} {$key2}"))
	// invlaid input
	is.Equal("hello {key0}", vs.ReplaceVars("hello {key0}"))
}
