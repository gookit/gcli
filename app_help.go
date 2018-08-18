package cliapp

import (
	"fmt"
	"github.com/gookit/cliapp/utils"
	"github.com/gookit/color"
	"os"
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
		utils.UpperFirst(app.Description),
		color.ApplyTag("cyan", app.Version),
	)

	os.Exit(0)
}

// display app commands help
func (app *Application) showCommandsHelp() {
	commandsHelp = color.ReplaceTag(commandsHelp)
	str := utils.RenderTemplate(commandsHelp, map[string]interface{}{
		"Cs": commands,
		// app version
		"Version": app.Version,
		// always upper first char
		"Description": utils.UpperFirst(app.Description),
	}, false)

	// parse help vars
	str = replaceVars(str, app.vars)
	fmt.Print(color.RenderStr(str))

	os.Exit(0)
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

// Print messages
func Print(args ...interface{}) (int, error) {
	return color.Print(args...)
}

// Println messages
func Println(args ...interface{}) (int, error) {
	return color.Println(args...)
}

// Printf messages
func Printf(format string, args ...interface{}) (int, error) {
	return color.Printf(format, args...)
}
