package gcli_test

import (
	"testing"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/testutil/assert"
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

	is.Eq("val3", vs.GetVar("key3"))
	is.Eq("", vs.GetVar("not-exist"))

	is.Eq("hello val0", vs.ReplaceVars("hello {$key0}"))
	is.Eq("hello val0 val2", vs.ReplaceVars("hello {$key0} {$key2}"))
	// invalid input
	is.Eq("hello {key0}", vs.ReplaceVars("hello {key0}"))
}

func Test_strictFormatArgs(t *testing.T) {
	str1 := ""
	t1 := false
	t2 := false
	t3 := false
	//t4 := false
	is := assert.New(t)
	cmd := gcli.NewCommand("init", "test bool pare", func(c *gcli.Command) {
		c.StrOpt(&str1, "name", "n", "", "test string parse")
		c.BoolOpt(&t1, "test1", "t", false, "test bool arse")
		c.BoolOpt(&t2, "test2", "s", false, "test bool arse")
		c.BoolOpt(&t3, "test3", "c", true, "test bool arse")
		//c.BoolOpt(&t4, "test4", "d", false, "test bool arse")
	})
	err := cmd.Run([]string{"-n", "ccc", "-test1=true", "-s", "--test3=false"})
	is.NoErr(err)
	is.Eq("ccc", str1)
	is.Eq(true, t1)
	is.Eq(true, t2)
	is.Eq(false, t3)
}
