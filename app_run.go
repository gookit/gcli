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

// store current application instance
var app *Application

// init
func init() {
	workDir, _ := os.Getwd()
	CLI.workDir = workDir

	// binName will contains work dir path on windows
	if utils.IsWin() {
		CLI.binName = strings.Replace(CLI.binName, workDir+"\\", "", 1)
	}
}

// App returns the current application instance
func App() *Application {
	return app
}

// AllCommands returns all commands
func AllCommands() map[string]*Command {
	return app.commands
}

// parseGlobalOpts parse global options
func parseGlobalOpts() []string {
	flag.Usage = app.showCommandsHelp
	// binding global options
	flag.UintVar(&gOpts.verbose, "verbose", gOpts.verbose, "")
	flag.BoolVar(&gOpts.showHelp, "h", false, "")
	flag.BoolVar(&gOpts.showHelp, "help", false, "")
	flag.BoolVar(&gOpts.showVer, "V", false, "")
	flag.BoolVar(&gOpts.showVer, "version", false, "")
	flag.BoolVar(&gOpts.noColor, "no-color", false, "")

	flag.Parse()
	return flag.Args()
}

// prepare to running
func (app *Application) prepareRun() (string, []string) {
	// don't display date on print log
	log.SetFlags(0)
	args := parseGlobalOpts()
	// args := flag.Args()

	if gOpts.showHelp {
		app.showCommandsHelp()
	}

	if gOpts.showVer {
		app.showVersionInfo()
	}

	if gOpts.noColor {
		color.Enable = false
	}

	// if no input command
	if len(args) < 1 {
		// will try run defaultCommand
		defCmd := app.defaultCommand
		if len(defCmd) == 0 {
			app.showCommandsHelp()
		}

		if !app.IsCommand(defCmd) {
			Logf(VerbError, "The default command '%s' is not exists", defCmd)
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

// Run running application
func (app *Application) Run() {
	rawName, args := app.prepareRun()
	name := app.RealCommandName(rawName)
	Logf(VerbDebug, "input command is: %s, real command: %s, args: %v", rawName, name, args)

	if !app.IsCommand(name) {
		color.Tips("error").Printf("unknown input command '%s'", name)

		ns := app.findSimilarCmd(name)
		if len(ns) > 0 {
			fmt.Println("\nMaybe you mean:\n  ", color.Green.Render(strings.Join(ns, ", ")))
		}

		fmt.Printf("\nUse \"%s\" to see available commands\n", color.Cyan.Render(CLI.binName+" -h"))
		Exit(ERR)
	}

	cmd := app.commands[name]
	app.commandName = name
	if app.Strict {
		args = strictFormatArgs(args)
	}

	// fmt.Println(cmd.argsStr, len(cmd.argsStr), strings.LastIndex(cmd.argsStr, " -h"))
	app.callHook(EvtBefore, cmd.Copy())

	// parse args, don't contains command name.
	if !cmd.CustomFlags {
		cmd.Flags.Parse(args)
		args = cmd.Flags.Args()
	}

	Logf(VerbDebug, "args for the command %s: %v", name, args)

	// do execute command
	exitCode := cmd.Execute(args)

	if len(app.errors) > 0 {
		app.callHook(EvtError, app.errors)
	} else {
		app.callHook(EvtAfter, exitCode)
	}

	Logf(VerbDebug, "command %s run complete, exit with code: ", name, exitCode)
	Exit(exitCode)
}

// SubRun running other command in current command
func (app *Application) SubRun(name string, args []string) int {
	if !app.IsCommand(name) {
		color.Tips("error").Printf("unknown command name '%s'", name)
		return ERR
	}

	cmd := app.commands[name]
	if !cmd.CustomFlags {
		// parse args, don't contains command name.
		cmd.Flags.Parse(args)
		args = cmd.Flags.Args()
	}

	// do execute command
	return cmd.Execute(args)
}

// Add add a command
// func (app *Application) AddCommander(c Commander) {
// 	// run command configure
// 	cmd := c.Configure()
//
// 	app.Add(cmd)
// }

// IsCommand name check
func (app *Application) IsCommand(name string) bool {
	_, has := app.names[name]
	return has
}

// CommandName get current command name
func (app *Application) CommandName() string {
	return app.commandName
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
			exitWithErr("The alias '%s' has been used by command '%s'", alias, cmd)
		}

		app.aliases[alias] = command
	}
}

// RealCommandName get real command name by alias
func (app *Application) RealCommandName(alias string) string {
	if name, has := app.aliases[alias]; has {
		return name
	}

	return alias
}
