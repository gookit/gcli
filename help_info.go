package cliapp

import (
	"os"
	"io"
	"html/template"
	"strings"
	"fmt"
	"reflect"
	"flag"
	"github.com/golangkit/cliapp/color"
)

// showVersionInfo display version info
func showVersionInfo() {
	fmt.Printf("%s\n\nVersion: %s\n", app.Description, app.Version)
	os.Exit(0)
}

// help template for all commands
var commandsHelp = `{{.Des}}
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
	renderString(os.Stdout, commandsHelp, map[string]interface{}{
		"Cs":     commands,
		"Des":    app.Description,
		"Script": script,
	})
	os.Exit(0)
}

// help template for a command
var commandHelp = `{{.Cmd.Description}}

Name: {{.Cmd.Name}}{{if .Cmd.Aliases}}(alias: {{.Cmd.Aliases.String}}){{end}}
Usage: {{.Script}} {{.Cmd.Name}} [--option ...] [argument ...]

Global Options:
  -h, --help        Display this help information{{if .Options}}

Options:
{{.Options}}
{{end}}{{if .Cmd.ArgList}}
Arguments:{{range $k,$v := .Cmd.ArgList}}
  {{$k | printf "%-12s"}}{{$v}}{{end}}
{{end}} {{if .Cmd.Examples}}
Examples:
  {{.Cmd.Examples}}{{end}}
`

// showCommandHelp display help for an command
func showCommandHelp(list []string, quit bool) {
	if len(list) != 1 {
		fmt.Fprintf(
			os.Stdout,
			"usage: %s help %s\n\nToo many arguments given.\n",
			script,
			list[0],
		)
		os.Exit(2) // failed at 'bee help'
	}

	// get real name
	name := GetCommandName(list[0])
	cmd, exist := commands[name]

	if !exist {
		Stdoutf(color.Color(color.FgRed).F("Unknown command name %#q.  Run '%s help'.", name, script))
		os.Exit(2)
	}

	// render help info
	renderString(os.Stdout, commandHelp, map[string]interface{}{
		"Cmd":     cmd,
		"Script":  script,
		"Options": cmd.ParseDefaults(),
	})

	if quit {
		os.Exit(0)
	}

	return
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

// PrintDefaults prints, to standard error unless configured otherwise, the
// default values of all defined command-line flags in the set. See the
// documentation for the global function PrintDefaults for more information.
// NOTICE: the func is copied from package 'flag', func 'PrintDefaults'
func (c *Command) ParseDefaults() string {
	var ss []string

	c.Flags.VisitAll(func(fg *flag.Flag) {
		var s string

		// is short option
		if len(fg.Name) == 1 {
			s = fmt.Sprintf("  -%s", fg.Name) // Two spaces before -; see next two comments.
		} else {
			s = fmt.Sprintf("  --%s", fg.Name) // Two spaces before -; see next two comments.
		}

		name, usage := flag.UnquoteUsage(fg)
		if len(name) > 0 {
			s += " " + name
		}
		// Boolean flags of one ASCII letter are so common we
		// treat them specially, putting their usage on the same line.
		if len(s) <= 4 { // space, space, '-', 'x'.
			s += "\t"
		} else {
			// Four spaces before the tab triggers good alignment
			// for both 4- and 8-space tab stops.
			s += "\n    \t"
		}
		s += strings.Replace(usage, "\n", "\n    \t", -1)

		if !isZeroValue(fg, fg.DefValue) {
			if _, ok := fg.Value.(*stringValue); ok {
				// put quotes on the value
				s += fmt.Sprintf(" (default %q)", fg.DefValue)
			} else {
				s += fmt.Sprintf(" (default %v)", fg.DefValue)
			}
		}

		ss = append(ss, s)
		// fmt.Fprint(fgs.Output(), s, "\n")
	})

	return strings.Join(ss, "\n")
}

// renderString
func renderString(w io.Writer, text string, data interface{}) {
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
