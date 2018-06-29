package cliapp

import (
	"fmt"
	"os"
	"github.com/golangkit/cliapp/color"
	"bytes"
	"flag"
	"strings"
	"reflect"
)

// help template for a command
var commandHelp = `{{.Description}}
{{if .Cmd.NotAlone}}
<comment>Name:</> {{.Cmd.Name}}{{if .Cmd.Aliases}} (alias: <info>{{.Cmd.Aliases.String}}</>){{end}}{{end}}
<comment>Usage:</> 
  {$binName} {{if .Cmd.NotAlone}}{{.Cmd.Name}} {{end}}[--option ...] [argument ...]

<comment>Global Options:</>
  <info>-h, --help</>        Display this help information{{if .Options}}

<comment>Options:</>
{{.Options}}{{end}}{{if .Cmd.ArgList}}
<comment>Arguments:</>{{range $k,$v := .Cmd.ArgList}}
  {{$k | printf "%-12s"}}{{$v}}{{end}}
{{end}} {{if .Cmd.Examples}}
<comment>Examples:</>
  {{.Cmd.Examples|coloredHtml}}{{end}}{{if .Cmd.Help}}
<comment>Help:</>
  {{.Cmd.Help|coloredHtml}}{{end}}
`

// showCommandHelp display help for an command
func (app *Application) showCommandHelp(list []string, quit bool) {
	if len(list) != 1 {
		color.Tips("error").Printf(
			"Usage: %s help %s\n\nToo many arguments given.",
			binName,
			list[0],
		)
		os.Exit(2) // failed at 'bee help'
	}

	// get real name
	name := app.GetNameByAlias(list[0])
	cmd, exist := commands[name]

	if !exist {
		color.Tips("error").Printf("Unknown command name %#q.  Run '%s -h'", name, binName)
		os.Exit(2)
	}

	cmd.ShowHelp(quit)
}

// ShowHelp show command help info
func (c *Command) ShowHelp(quit ...bool) {
	// use buffer receive rendered content
	var buf bytes.Buffer

	commandHelp = color.ReplaceTag(commandHelp)

	// render and output help info
	//RenderStrTpl(os.Stdout, commandHelp, map[string]interface{}{
	// render but not output
	RenderStrTpl(&buf, commandHelp, map[string]interface{}{
		"Cmd":     c,
		"Options": color.RenderStr(c.ParseDefaults()),

		// always upper first char
		"Description": color.RenderStr(c.Description),
	})

	c.Vars["cmd"] = c.Name
	c.Vars["fullCmd"] = binName + " " + c.Name

	// parse help vars
	fmt.Print(ReplaceVars(buf.String(), c.Vars))

	if len(quit) > 0 && quit[0] {
		os.Exit(0)
	}
}

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
			s = fmt.Sprintf("  <info>-%s</>", fg.Name) // Two spaces before -; see next two comments.
		} else {
			s = fmt.Sprintf("  <info>--%s</>", fg.Name)
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
				s += fmt.Sprintf(" (default <cyan>%q</>)", fg.DefValue)
			} else {
				s += fmt.Sprintf(" (default <cyan>%v</>)", fg.DefValue)
			}
		}

		ss = append(ss, s)
		// fmt.Fprint(fgs.Output(), s, "\n")
	})

	return strings.Join(ss, "\n")
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
