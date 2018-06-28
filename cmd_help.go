package cliapp

import (
	"fmt"
	"os"
	"github.com/golangkit/cliapp/color"
)

// help template for a command
var commandHelp = `{{.Description}}
{{if .Cmd.NotAlone}}
Name: {{.Cmd.Name}}{{if .Cmd.Aliases}}(alias: {{.Cmd.Aliases.String}}){{end}}{{end}}
Usage: {{.Script}} {{if .Cmd.NotAlone}}{{.Cmd.Name}} {{end}}[--option ...] [argument ...]

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
		color.Tips("danger").Printf("Unknown command name %#q.  Run '%s -h'\n", name, script)
		os.Exit(2)
	}

	cmd.ShowHelp(quit)
}
