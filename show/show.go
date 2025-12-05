// Package show provides some formatter tools for display data.
package show

import (
	"encoding/json"
	"fmt"
	"text/tabwriter"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/gclicom"
	"github.com/gookit/gcli/v3/show/banner"
	"github.com/gookit/gcli/v3/show/lists"
	"github.com/gookit/gcli/v3/show/title"
)

const (
	// OK success exit code
	OK = 0
	// ERR error exit code
	ERR = 2
)

// Error tips message print
func Error(format string, v ...any) int {
	prefix := color.Red.Sprint("ERROR: ")
	_, _ = fmt.Fprintf(gclicom.Output, prefix+format+"\n", v...)
	return ERR
}

// Success tips message print
func Success(format string, v ...any) int {
	prefix := color.Green.Sprint("SUCCESS: ")
	_, _ = fmt.Fprintf(gclicom.Output, prefix+format+"\n", v...)
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

	_, _ = fmt.Fprintln(gclicom.Output, string(bs))
	return OK
}

// ATitle create a Title instance and print. options see: TitleOption
func ATitle(titleText string, fns ...title.OptionFunc) {
	title.New(titleText).WithOptionFns(fns).Println()
}

type ListOpFunc = lists.ListOpFunc

// NewList create a List instance. options see: ListOption
func NewList(title string, data any, fns ...ListOpFunc) *lists.List {
	return lists.NewList(title, data).WithOptionFns(fns)
}

// NewLists create a Lists instance and print. options see: ListOption
func NewLists(listMap any, fns ...ListOpFunc) *lists.Lists {
	return lists.NewLists(listMap).WithOptionFns(fns)
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

// NewBanner create a Banner instance. options see: banner.Options
func NewBanner(content any, fns ...banner.OptionFunc) *banner.Banner {
	return banner.New(content, fns...)
}

// Banner create a Banner instance and print. options see: banner.Options
func Banner(content any, fns ...banner.OptionFunc) {
	banner.New(content, fns...).Println()
}

// TabWriter create. more please see: package text/tabwriter/example_test.go
//
// Usage:
//
//	w := TabWriter([]string{
//		"a\tb\tc\td\t.",
//		"123\t12345\t1234567\t123456789\t."
//	})
//	w.Flush()
func TabWriter(rows []string) *tabwriter.Writer {
	w := tabwriter.NewWriter(gclicom.Output, 0, 4, 2, ' ', tabwriter.Debug)

	for _, row := range rows {
		if _, err := fmt.Fprintln(w, row); err != nil {
			panic(err)
		}
	}

	return w
}
