package cliapp

import (
	"fmt"
	"os"
	"github.com/golangkit/cliapp/color"
	"bytes"
)

// help template for a command
var commandHelp = `{{.Description}}
{{if .Cmd.NotAlone}}
<comment>Name:</> {{.Cmd.Name}}{{if .Cmd.Aliases}} (alias: <info>{{.Cmd.Aliases.String}}</>){{end}}{{end}}
<comment>Usage:</> {{.Script}} {{if .Cmd.NotAlone}}{{.Cmd.Name}} {{end}}[--option ...] [argument ...]

<comment>Global Options:</>
  -h, --help        Display this help information{{if .Options}}

<comment>Options:</>
{{.Options}}
{{end}}{{if .Cmd.ArgList}}
<comment>Arguments:</>{{range $k,$v := .Cmd.ArgList}}
  {{$k | printf "%-12s"}}{{$v}}{{end}}
{{end}} {{if .Cmd.Examples}}
<comment>Examples:</>
  {{.Cmd.Examples}}{{end}}{{if .Cmd.Help}}
<comment>Help:</>
  {{.Cmd.Help}}{{end}}
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
	name := FindCommandName(list[0])
	cmd, exist := commands[name]

	if !exist {
		color.Tips("danger").Printf("Unknown command name %#q.  Run '%s -h'\n", name, script)
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
		"Cmd":         c,
		"Script":      script,
		"Options":     color.ReplaceTag(c.ParseDefaults()),
		"Description": color.ReplaceTag(c.Description),
	})

	c.Vars["cmd"] = c.Name

	// parse help vars
	fmt.Print(ReplaceVars(buf.String(), c.Vars))

	if len(quit) > 0 && quit[0] {
		os.Exit(0)
	}
}
