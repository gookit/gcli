package cliapp

import (
	"os"
	"io"
	"html/template"
	"strings"
	"fmt"
	"reflect"
	"flag"
)

// showVersionInfo display version info
func showVersionInfo() {
	fmt.Printf("%s\n\nVersion: %s\n", app.Description, app.Version)
	os.Exit(0)
}

// help template for all commands
var commandsHelp = `{{.Description}}
Usage:
  {{.Script}} command [--option ...] [argument ...]

Options:
  -h, --help        Display this help information
  -V, --version     Display this version information

Commands:{{range .Cs}}{{if .Runnable}}
  {{.Name | printf "%-12s"}} {{.Description}}{{if .Aliases}}(alias: {{.Aliases.String}}){{end}}{{end}}{{end}}
  help         display help information

Use "{{.Script}} help [command]" for more information about a command
`

// showCommandsHelp commands list
func showCommandsHelp() {
	RenderStrTpl(os.Stdout, commandsHelp, map[string]interface{}{
		"Cs":          commands,
		"Script":      script,
		"Description": app.Description,
	})

	os.Exit(0)
}

// isZeroValue guesses whether the string represents the zero
// value for a flag. It is not accurate but in practice works OK.
// NOTICE: the func is copied from package 'flag', func 'isZeroValue'
func isZeroValue(fg *flag.Flag, value string) bool {
	// Build a zero value of the flag's Value type, and see if the
	// result of calling its String method equals the value passed in.
	// This works unless the Value type is itself an interface type.
	typ := reflect.TypeOf(fg.Value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	if value == z.Interface().(flag.Value).String() {
		return true
	}

	switch value {
	case "false", "", "0":
		return true
	}
	return false
}

// -- string Value
// NOTICE: the var is copied from package 'flag'
type stringValue string

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}
func (s *stringValue) Get() interface{} { return string(*s) }
func (s *stringValue) String() string   { return string(*s) }

// RenderStrTpl
func RenderStrTpl(w io.Writer, text string, data interface{}) {
	t := template.New("cli")

	t.Funcs(template.FuncMap{"trim": func(s template.HTML) template.HTML {
		return template.HTML(strings.TrimSpace(string(s)))
	}})

	t.Funcs(template.FuncMap{"joinStrings": func(ss []string) template.HTML {
		return template.HTML(strings.Join(ss, ","))
	}})

	template.Must(t.Parse(text))

	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}
