package show_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/testutil/assert"
)

func TestList(t *testing.T) {
	// is := assert.New(t)
	l := show.NewList("test list", []string{
		"list item 0",
		"list item 1",
		"list item 2",
	})
	l.Println()

	l = show.NewList("test list1", map[string]string{
		"key0":     "list item 0",
		"the key1": "list item 1",
		"key2":     "list item 2",
		"key3":     "", // empty value
	})
	l.Opts.SepChar = " | "
	l.Println()
}

func TestList_mlevel(t *testing.T) {
	d := map[string]interface{}{
		"key0":     "list item 0",
		"key2":     []string{"abc", "def"},
		"key4":     map[string]int{"abc": 23, "def": 45},
		"the key1": "list item 1",
		"key3":     "", // empty value
	}

	l := show.NewList("test list", d)
	l.Println()

	l = show.NewList("test list2", d).WithOptions(func(opts *show.ListOption) {
		opts.SepChar = " | "
	})
	l.Println()
}

func TestLists(t *testing.T) {
	ls := show.NewLists(map[string]interface{}{
		"test list": []string{
			"list item 0",
			"list item 1",
			"list item 2",
		},
		"test list1": map[string]string{
			"key0":     "list item 0",
			"the key1": "list item 1",
			"key2":     "list item 2",
			"key3":     "", // empty value
		},
	})
	ls.Opts.SepChar = " : "
	ls.Println()
}

func TestTabWriter(t *testing.T) {
	is := assert.New(t)
	ss := []string{
		"a\tb\taligned\t",
		"aa\tbb\taligned\t",
		"aaa\tbbb\tunaligned",
		"aaaa\tbbbb\taligned\t",
	}

	err := show.TabWriter(os.Stdout, ss).Flush()
	is.NoErr(err)
}

func TestSome(t *testing.T) {
	fmt.Printf("|%8s|\n", "text")
	fmt.Printf("|%-8s|\n", "text")
	fmt.Printf("|%8s|\n", "text")
}
