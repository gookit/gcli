package helper

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/gookit/goutil/strutil"
)

const (
	// RegGoodName match a good option, argument name
	RegGoodName = `^[a-zA-Z][\w-]*$`
	// RegGoodCmdName match a good command name
	RegGoodCmdName = `^[a-zA-Z][\w-]*$`
	// RegGoodCmdId match command id. eg: "self:init"
	RegGoodCmdId = `^[a-zA-Z][\w:-]*$`
	// match command path. eg: "self init"
	// RegGoodCmdPath = `^[a-zA-Z][\w -]*$`
)

var (
	// GoodName good name for option and argument
	goodName = regexp.MustCompile(RegGoodName)
	// GoodCmdId match a good command name
	goodCmdId = regexp.MustCompile(RegGoodCmdId)
	// GoodCmdName match a good command name
	goodCmdName = regexp.MustCompile(RegGoodCmdName)
)

// IsGoodName check
func IsGoodName(name string) bool {
	return goodName.MatchString(name)
}

// IsGoodCmdId check
func IsGoodCmdId(name string) bool {
	return goodCmdId.MatchString(name)
}

// IsGoodCmdName check
func IsGoodCmdName(name string) bool {
	return goodCmdName.MatchString(name)
}

// Panicf message
func Panicf(format string, v ...any) {
	panic(fmt.Sprintf("GCli: "+format, v...))
}

// RenderText render text template with data. TODO use strutil.RenderText()
func RenderText(input string, data any, fns template.FuncMap, isFile ...bool) string {
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
