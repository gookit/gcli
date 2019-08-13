package gcli

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/gcli/utils"
	"github.com/gookit/goutil/strutil"
)

// help template for all commands
var commandsHelp = `{{.Description}} (Version: <info>{{.Version}}</>)
<comment>Usage:</>
  {$binName} [Global Options...] <info>{command}</> [--option ...] [argument ...]

<comment>Global Options:</>
      <info>--verbose</>     Set error reporting level(quiet 0 - 4 debug)
      <info>--no-color</>    Disable color when outputting message
  <info>-h, --help</>        Display the help information
  <info>-V, --version</>     Display app version information

<comment>Available Commands:</>{{range .Cs}}{{if .Runnable}}
  <info>{{.Name | printf "%-12s"}}</> {{.UseFor}}{{if .Aliases}} (alias: <cyan>{{ join .Aliases ","}}</>){{end}}{{end}}{{end}}
  <info>help</>         Display help information

Use "<cyan>{$binName} {command} -h</>" for more information about a command
`

// display app version info
func (app *App) showVersionInfo() {
	fmt.Printf(
		"%s\n\nVersion: %s\n",
		strutil.UpperFirst(app.Description),
		color.ApplyTag("cyan", app.Version),
	)

	if app.Logo.Text != "" {
		fmt.Printf("%s\n", color.ApplyTag(app.Logo.Style, app.Logo.Text))
	}

	Exit(OK)
}

// display app commands help
func (app *App) showCommandsHelp() {
	commandsHelp = color.ReplaceTag(commandsHelp)
	// render help text template
	s := utils.RenderTemplate(commandsHelp, map[string]interface{}{
		"Cs": app.commands,
		// app version
		"Version": app.Version,
		// always upper first char
		"Description": strutil.UpperFirst(app.Description),
	}, false)

	// parse help vars and render color tags
	fmt.Print(color.String(replaceVars(s, app.vars)))
	Exit(OK)
}

// showCommandHelp display help for an command
func (app *App) showCommandHelp(list []string, quit bool) {
	if len(list) != 1 {
		color.Error.Tips(
			"Usage: %s help %s\n\nToo many arguments given.",
			CLI.binName,
			list[0],
		)
		Exit(ERR)
	}

	// get real name
	name := app.RealCommandName(list[0])
	cmd, exist := app.commands[name]
	if !exist {
		color.Error.Prompt("Unknown command name %#q. Run '%s -h'", name, CLI.binName)
		Exit(ERR)
	}

	cmd.ShowHelp(quit)
}

func (app *App) showCompletion(args []string) {

}

// findSimilarCmd find similar cmd by input string
func (app *App) findSimilarCmd(input string) []string {
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
