package cliapp

import (
	"fmt"
	"github.com/gookit/cliapp/utils"
	"github.com/gookit/color"
	"os"
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
  <info>{{.Name | printf "%-12s"}}</> {{.Description}}{{if .Aliases}} (alias: <cyan>{{.Aliases.String}}</>){{end}}{{end}}{{end}}
  <info>help</>         Display help information

Use "<cyan>{$binName} help {command}</>" for more information about a command
`

// showVersionInfo display version info
func (app *Application) showVersionInfo() {
	fmt.Printf(
		"%s\n\nVersion: %s\n",
		utils.UpperFirst(app.Description),
		color.ApplyTag("cyan", app.Version),
	)

	os.Exit(0)
}

// showCommandsHelp commands list
func showCommandsHelp() {
	commandsHelp = color.ReplaceTag(commandsHelp)

	str := utils.RenderTemplate(commandsHelp, map[string]interface{}{
		"Cs":      commands,
		"Version": app.Version,
		// always upper first char
		"Description": utils.UpperFirst(app.Description),
	}, false)

	// parse help vars
	str = replaceVars(str, app.vars)
	fmt.Print(color.RenderStr(str))

	os.Exit(0)
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
