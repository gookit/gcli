package cliapp

import (
	"flag"
	"fmt"
	"github.com/gookit/cliapp/utils"
	"github.com/gookit/color"
	"log"
	"os"
	"strings"
)

// the cli app instance
var app *Application

// commands collect all command
var commands map[string]*Command

// var commanders  map[string]Commander

// init
func init() {
	commands = make(map[string]*Command)

	// binName will contains work dir path on windows
	if utils.IsWin() {
		binName = strings.Replace(binName, workDir+"\\", "", 1)
	}
}

// Run running application
func (app *Application) Run() {
	rawName, args := app.prepareRun()
	name := app.GetNameByAlias(rawName)
	logf(VerbDebug, "input command name is: %s", name)

	if !app.IsCommand(name) {
		color.Tips("error").Printf("unknown input command '%s'", name)

		ns := app.findSimilarCmd(name, app.CommandNames())
		if len(ns) > 0 {
			fmt.Println("\nMaybe you means:\n  ", color.FgGreen.Render(strings.Join(ns, ", ")))
		}

		color.Info("\nPlease use '", binName, " -h' to see available commands\n")
		Exit(-2)
	}

	cmd := commands[name]
	app.command = name

	if app.Strict {
		args = strictFormatArgs(args)
	}

	// parse args, don't contains command name.
	if !cmd.CustomFlags {
		cmd.Flags.Parse(args)
		args = cmd.Flags.Args()
	}

	// do execute command
	os.Exit(cmd.Execute(app, args))
}

// parseGlobalOpts parse global options
func parseGlobalOpts() []string {
	// Some global options
	flag.Usage = app.showCommandsHelp

	flag.UintVar(&gOpts.verbose, "verbose", gOpts.verbose, "")
	flag.BoolVar(&gOpts.showHelp, "h", false, "")
	flag.BoolVar(&gOpts.showHelp, "help", false, "")
	flag.BoolVar(&gOpts.showVer, "V", false, "")
	flag.BoolVar(&gOpts.showVer, "version", false, "")
	flag.BoolVar(&gOpts.noColor, "no-color", false, "")

	flag.Parse()
	return flag.Args()
}

// Run running a sub-command in current command
func (app *Application) SubRun(name string, args []string) int {
	if !app.IsCommand(name) {
		color.Tips("error").Printf("unknown input command '%s'", name)
		return -2
	}

	cmd := commands[name]

	// parse args, don't contains command name.
	if !cmd.CustomFlags {
		cmd.Flags.Parse(args)
		args = cmd.Flags.Args()
	}

	// do execute command
	return cmd.Execute(app, args)
}

// prepare to running
func (app *Application) prepareRun() (string, []string) {
	// don't display date on print log
	log.SetFlags(0)
	args := parseGlobalOpts()

	if gOpts.showHelp {
		app.showCommandsHelp()
	}

	if gOpts.showVer {
		app.showVersionInfo()
	}

	if gOpts.noColor {
		color.Enable = false
	}

	// check args
	if len(args) < 1 {
		defCmd := app.defaultCmd
		if len(defCmd) == 0 {
			app.showCommandsHelp()
		}

		if !app.IsCommand(defCmd) {
			color.Tips("error").Printf("The default command '%s' is not exists", defCmd)
			app.showCommandsHelp()
		}

		args = []string{defCmd}
	}

	// is help command
	if args[0] == "help" {
		// like 'go help'
		if len(args) == 1 {
			app.showCommandsHelp()
		}

		// like 'go help COMMAND'
		app.showCommandHelp(args[1:], true)
	}

	return args[0], args[1:]
}

// Add add a command
// func (app *Application) AddCommander(c Commander) {
// 	// run command configure
// 	cmd := c.Configure()
//
// 	app.Add(cmd)
// }

// Command get command name
func (app *Application) Command() string {
	return app.command
}

// IsCommand name check
func (app *Application) IsCommand(name string) bool {
	_, has := app.names[name]

	return has
}

// CommandNames get all command names
func (app *Application) CommandNames() []string {
	var ss []string

	for n := range app.names {
		ss = append(ss, n)
	}

	return ss
}

// AddAliases add alias names for a command
func (app *Application) AddAliases(command string, names []string) {
	if app.aliases == nil {
		app.aliases = make(map[string]string)
	}

	// add alias
	for _, alias := range names {
		if cmd, has := app.aliases[alias]; has {
			panic(color.FgRed.Renderf("The alias '%s' has been used by command '%s'", alias, cmd))
		}

		app.aliases[alias] = command
	}
}

// GetNameByAlias get real command name by alias
func (app *Application) GetNameByAlias(alias string) string {
	if name, has := app.aliases[alias]; has {
		return name
	}

	return alias
}

// Commands get all commands
func (app *Application) Commands() map[string]*Command {
	return commands
}

// AllCommands get all commands
func AllCommands() map[string]*Command {
	return commands
}
