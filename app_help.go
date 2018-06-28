package cliapp

import (
	"os"
	"io"
	"html/template"
	"strings"
	"reflect"
	"flag"
	"github.com/golangkit/cliapp/color"
	"bytes"
	"fmt"
	"github.com/golangkit/cliapp/utils"
)

// showVersionInfo display version info
func showVersionInfo() {
	color.Printf(
		"%s\n\nVersion: <suc>%s</>\n",
		utils.UpperFirst(app.Description),
		app.Version,
	)
	os.Exit(0)
}

// help template for all commands
var commandsHelp = `{{.Description|raw}} (Version: <info>{{.Version}}</>)
<comment>Usage:</>
  {{.Script}} <info>{command}</> [--option ...] [argument ...]

<comment>Options:</>
  -h, --help        Display this help information
  -V, --version     Display this version information

<comment>Commands:</>{{range .Cs}}{{if .Runnable}}
  {{.Name | printf "%-12s"}} {{.Description|colored}}{{if .Aliases}} (alias: <cyan>{{.Aliases.String}}</>){{end}}{{end}}{{end}}
  help         display help information

Use "<cyan>{{.Script}} help [command]</>" for more information about a command
`

// showCommandsHelp commands list
func showCommandsHelp() {
	// use buffer receive rendered content
	var buf bytes.Buffer

	commandsHelp = color.ReplaceTag(commandsHelp)

	//RenderStrTpl(os.Print, commandsHelp, map[string]interface{}{
	RenderStrTpl(&buf, commandsHelp, map[string]interface{}{
		"Cs":          commands,
		"Script":      script,
		"Version":     app.Version,
		"Description": utils.UpperFirst(app.Description),
	})

	fmt.Print(ReplaceVars(buf.String(), app.vars))

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

	// don't escape content
	t.Funcs(template.FuncMap{"raw": func(s string) interface{} {
		return template.HTML(s)
	}})

	t.Funcs(template.FuncMap{"colored": func(s string) interface{} {
		return template.HTML(color.ReplaceTag(s))
	}})

	t.Funcs(template.FuncMap{"coloredHtml": func(h template.HTML) interface{} {
		return template.HTML(color.ReplaceTag(string(h)))
	}})

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
