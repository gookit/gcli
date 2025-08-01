// Package show provides some formatter tools for display data.
package show

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/gookit/color"
)

// Output the global input out stream
var Output io.Writer = os.Stdout

// SetOutput stream
func SetOutput(out io.Writer) { Output = out }

// ResetOutput stream
func ResetOutput() { Output = os.Stdout }

// Error tips message print
func Error(format string, v ...any) int {
	prefix := color.Red.Sprint("ERROR: ")
	_, _ = fmt.Fprintf(Output, prefix+format+"\n", v...)
	return ERR
}

// Success tips message print
func Success(format string, v ...any) int {
	prefix := color.Green.Sprint("SUCCESS: ")
	_, _ = fmt.Fprintf(Output, prefix+format+"\n", v...)
	return OK
}

// JSON print pretty JSON data
func JSON(v any, prefixAndIndent ...string) int {
	prefix := ""
	indent := "    "

	l := len(prefixAndIndent)
	if l > 0 {
		prefix = prefixAndIndent[0]
		if l > 1 {
			indent = prefixAndIndent[1]
		}
	}

	bs, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		panic(err)
	}

	_, _ = fmt.Fprintln(Output, string(bs))
	return OK
}

// AList create a List instance and print. options see: ListOption
//
// Usage:
//
//	show.AList("some info", map[string]string{"name": "tom"})
func AList(title string, data any, fns ...ListOpFunc) {
	NewList(title, data).WithOptionFns(fns).Println()
}

// MList show multi list data. options see: ListOption
//
// Usage:
//
//	show.MList(data)
//	show.MList(data, func(opts *ListOption) {
//		opts.LeftIndent = "    "
//	})
func MList(listMap any, fns ...ListOpFunc) {
	NewLists(listMap).WithOptionFns(fns).Println()
}

// TabWriter create.
// more please see: package text/tabwriter/example_test.go
//
// Usage:
//
//	w := TabWriter([]string{
//		"a\tb\tc\td\t.",
//		"123\t12345\t1234567\t123456789\t."
//	})
//	w.Flush()
func TabWriter(rows []string) *tabwriter.Writer {
	w := tabwriter.NewWriter(Output, 0, 4, 2, ' ', tabwriter.Debug)

	for _, row := range rows {
		if _, err := fmt.Fprintln(w, row); err != nil {
			panic(err)
		}
	}

	return w
}
