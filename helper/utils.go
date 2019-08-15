package helper

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/gookit/goutil/strutil"
)

// GetScreenSize for current console terminal. TODO ...
func GetScreenSize() (w int, h int) {
	return
}

// RenderText render text template with data
func RenderText(input string, data interface{}, fns template.FuncMap, isFile ...bool) string {
	t := template.New("cli")
	t.Funcs(template.FuncMap{
		// don't escape content
		"raw": func(s string) string {
			return s
		},
		"trim": strings.TrimSpace,
		// join strings. usage {{ join .Strings ","}}
		"join": func(ss []string, sep string) string {
			return strings.Join(ss, sep)
		},
		// lower first char
		"lcFirst": strutil.LowerFirst,
		// upper first char
		"ucFirst": strutil.UpperFirst,
	})

	// custom add template functions
	if len(fns) > 0 {
		t.Funcs(fns)
	}

	if len(isFile) > 0 && isFile[0] {
		template.Must(t.ParseFiles(input))
	} else {
		template.Must(t.Parse(input))
	}

	// use buffer receive rendered content
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		panic(err)
	}

	return buf.String()
}
