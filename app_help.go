package cliapp

import (
	"fmt"
	"github.com/gookit/cliapp/utils"
	"github.com/gookit/color"
	"strings"
)

// help template for all commands
var commandsHelp = `{{.Description}} (Version: <info>{{.Version}}</>)
<comment>Usage:</>
  {$binName} [global options...] <info>{command}</> [--option ...] [argument ...]

<comment>Global Options:</>
      <info>--verbose</>     Set error reporting level(quiet 0 - 4 debug)
      <info>--no-color</>    Disable color when outputting message
  <info>-h, --help</>        Display the help information
  <info>-V, --version</>     Display app version information

<comment>Available Commands:</>{{range .Cs}}{{if .Runnable}}
  <info>{{.Name | printf "%-12s"}}</> {{.Description}}{{if .Aliases}} (alias: <cyan>{{ join .Aliases ","}}</>){{end}}{{end}}{{end}}
  <info>help</>         Display help information

Use "<cyan>{$binName} {command} -h</>" for more information about a command
`

// display app version info
func (app *Application) showVersionInfo() {
	fmt.Printf(
		"%s\n\nVersion: %s\n",
		utils.UcFirst(app.Description),
		color.ApplyTag("cyan", app.Version),
	)
	Exit(OK)
}

// display app commands help
func (app *Application) showCommandsHelp() {
	commandsHelp = color.ReplaceTag(commandsHelp)
	// render help text template
	str := utils.RenderTemplate(commandsHelp, map[string]interface{}{
		"Cs": commands,
		// app version
		"Version": app.Version,
		// always upper first char
		"Description": utils.UcFirst(app.Description),
	}, false)

	// parse help vars and render color tags
	fmt.Print(color.RenderStr(replaceVars(str, app.vars)))
	Exit(OK)
}

// showCommandHelp display help for an command
func (app *Application) showCommandHelp(list []string, quit bool) {
	if len(list) != 1 {
		color.Tips("error").Printf(
			"Usage: %s help %s\n\nToo many arguments given.",
			binName,
			list[0],
		)
		Exit(ERR)
	}

	// get real name
	name := app.RealCommandName(list[0])
	cmd, exist := commands[name]
	if !exist {
		color.Tips("error").Printf("Unknown command name %#q. Run '%s -h'", name, binName)
		Exit(ERR)
	}

	cmd.ShowHelp(quit)
}

// findSimilarCmd find similar cmd by input string
func (app *Application) findSimilarCmd(input string) []string {
	var ss []string
	// ins := strings.Split(input, "")
	// fmt.Print(input, ins)
	ln := len(input)

	// find from command names
	for name := range app.names {
		cln := len(name)
		if cln > ln && strings.Contains(name, input) {
			ss = append(ss, name)
		} else if ln > cln && strings.Contains(input, name) {
			// sns := strings.Split(str, "")
			ss = append(ss, name)
		}

		// max find 5 items
		if len(ss) == 5 {
			break
		}
	}

	// find from aliases
	for alias := range app.aliases {
		// max find 5 items
		if len(ss) >= 5 {
			break
		}

		cln := len(alias)
		if cln > ln && strings.Contains(alias, input) {
			ss = append(ss, alias)
		} else if ln > cln && strings.Contains(input, alias) {
			ss = append(ss, alias)
		}
	}

	return ss
}
