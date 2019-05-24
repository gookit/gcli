package show_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gookit/gcli/v2/show"
	"github.com/stretchr/testify/assert"
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
	is.NoError(err)
}

func TestSome(t *testing.T) {
	fmt.Printf("|%8s|\n", "text")
	fmt.Printf("|%-8s|\n", "text")
	fmt.Printf("|%8s|\n", "text")
}
