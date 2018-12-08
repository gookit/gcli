// Package show provides some formatter tools for display data.
package show

import (
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
	"io"
	"text/tabwriter"
)

// Error tips message print
func Error(format string, v ...interface{}) int {
	color.Red.Print("ERROR: ")
	fmt.Printf(format+"\n", v...)
	return ERR
}

// Success tips message print
func Success(format string, v ...interface{}) int {
	color.Green.Print("SUCCESS: ")
	fmt.Printf(format+"\n", v...)
	return OK
}

// JSON print pretty JSON data
func JSON(v interface{}, settings ...string) int {
	prefix := ""
	indent := "    "

	l := len(settings)
	if l > 0 {
		prefix = settings[0]
		if l > 1 {
			indent = settings[1]
		}
	}

	bs, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bs))
	return OK
}

// AList create a List instance and print.
// Usage:
// 	show.AList("some info", map[string]string{"name": "tom"}, nil)
func AList(title string, data interface{}, fn func(opts *ListOption)) {
	NewList(title, data).WithOptions(fn).Println()
}

// MList show multi list data.
// Usage:
// 	show.MList(data, nil)
// 	show.MList(data, func(opts *ListOption) {
// 		opts.LeftIndent = "    "
// 	})
func MList(listMap map[string]interface{}, fn func(opts *ListOption)) {
	NewLists(listMap).WithOptions(fn).Println()
}

// TabWriter create.
// more please see: package text/tabwriter/example_test.go
// Usage:
// 	w := TabWriter(os.Stdout, []string{
// 		"a\tb\tc\td\t.",
// 		"123\t12345\t1234567\t123456789\t."
// 	})
// 	w.Flush()
func TabWriter(outTo io.Writer, rows []string) *tabwriter.Writer {
	w := tabwriter.NewWriter(outTo, 0, 4, 2, ' ', tabwriter.Debug)

	for _, row := range rows {
		if _, err := fmt.Fprintln(w, row); err != nil {
			panic(err)
		}
	}

	return w
}
