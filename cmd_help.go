package cliapp

import (
	"fmt"
	"os"
	"github.com/golangkit/cliapp/color"
	"bytes"
	"github.com/golangkit/cliapp/utils"
	"flag"
	"strings"
)

// help template for a command
var commandHelp = `{{.Description}}{{if .Cmd.NotAlone}}
<comment>Name:</> {{.Cmd.Name}}{{if .Cmd.Aliases}} (alias: <info>{{.Cmd.Aliases.String}}</>){{end}}{{end}}
<comment>Usage:</> 
  {{.Script}} {{if .Cmd.NotAlone}}{{.Cmd.Name}} {{end}}[--option ...] [argument ...]

<comment>Global Options:</>
  -h, --help        Display this help information{{if .Options}}

<comment>Options:</>
{{.Options}}
{{end}}{{if .Cmd.ArgList}}
<comment>Arguments:</>{{range $k,$v := .Cmd.ArgList}}
  {{$k | printf "%-12s"}}{{$v}}{{end}}
{{end}} {{if .Cmd.Examples}}
<comment>Examples:</>
  {{.Cmd.Examples|coloredHtml}}{{end}}{{if .Cmd.Help}}
<comment>Help:</>
  {{.Cmd.Help|coloredHtml}}{{end}}
`

// showCommandHelp display help for an command
func showCommandHelp(list []string, quit bool) {
	if len(list) != 1 {
		color.Tips("error").Printf(
			"Usage: %s help %s\n\nToo many arguments given.",
			script,
			list[0],
		)
		os.Exit(2) // failed at 'bee help'
	}

	// get real name
	name := FindCommandName(list[0])
	cmd, exist := commands[name]

	if !exist {
		color.Tips("error").Printf("Unknown command name %#q.  Run '%s -h'", name, script)
		os.Exit(2)
	}

	cmd.ShowHelp(quit)
}

// ShowHelp @notice not used
func (c *Command) ShowHelp(quit ...bool) {
	// use buffer receive rendered content
	var buf bytes.Buffer

	commandHelp = color.ReplaceTag(commandHelp)

	// render and output help info
	//RenderStrTpl(os.Stdout, commandHelp, map[string]interface{}{
	// render but not output
	RenderStrTpl(&buf, commandHelp, map[string]interface{}{
		"Cmd": c,

		"Script":      script,
		"Options":     color.RenderStr(c.ParseDefaults()),
		"Description": color.RenderStr(utils.UpperFirst(c.Description)),
	})

	c.Vars["cmd"] = c.Name
	c.Vars["fullCmd"] = script + " " + c.Name

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
