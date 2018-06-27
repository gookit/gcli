package cliapp

import (
	"fmt"
	"os"
	"github.com/golangkit/cliapp/color"
	"bytes"
)

// help template for a command
var commandHelp = `{{.Description}}

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
	name := FindCommandName(list[0])
	cmd, exist := commands[name]

	if !exist {
		color.Tips("danger").Printf("Unknown command name %#q.  Run '%s help'.", name, script)
		os.Exit(2)
	}

	// use buffer receive rendered content
	var buf bytes.Buffer

	// render and output help info
	//RenderStrTpl(os.Stdout, commandHelp, map[string]interface{}{
	// render but not output
	RenderStrTpl(&buf, commandHelp, map[string]interface{}{
		"Cmd":         cmd,
		"Script":      script,
		"Options":     color.Render(cmd.ParseDefaults()),
		"Description": color.Render(cmd.Description),
	})

	cmd.Vars["cmd"] = name

	// parse help vars
	fmt.Print(ReplaceVars(buf.String(), cmd.Vars))

	if quit {
		os.Exit(0)
	}

	return
}
