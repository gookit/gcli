package cliapp

import (
	"flag"
	"fmt"
	"github.com/gookit/cliapp/utils"
	"github.com/gookit/color"
	"log"
	"os"
	"runtime"
	"strings"
)

var (
	// store current application instance
	defApp *App
	// global options
	gOpts = &GlobalOpts{verbose: VerbError}
	// CLI create a default instance
	CLI = &CmdLine{
		pid: os.Getpid(),
		// more info
		osName:  runtime.GOOS,
		binName: os.Args[0],
		argLine: strings.Join(os.Args[1:], " "),
	}
)

// init
func init() {
	workDir, _ := os.Getwd()
	CLI.workDir = workDir

	// binName will contains work dir path on windows
	if utils.IsWin() {
		CLI.binName = strings.Replace(CLI.binName, workDir+"\\", "", 1)
	}
}

// Instance returns the current application instance
func Instance() *App {
	return defApp
}

// AllCommands returns all commands
func AllCommands() map[string]*Command {
	return defApp.commands
}

// parseGlobalOpts parse global options
func (app *App) parseGlobalOpts() []string {
	// don't display date on print log
	log.SetFlags(0)

	// bind help func
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

// prepare to running, parse args, get command name and command args
func (app *App) prepareRun() (string, []string) {
	args := app.parseGlobalOpts()

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
	if len(args) == 0 {
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
	} else if args[0] == "help" { // is help command
		if len(args) == 1 { // like 'help'
			app.showCommandsHelp()
		}

		// like 'help COMMAND'
		app.showCommandHelp(args[1:], true)
	}

	return args[0], args[1:]
}

// Run running application
func (app *App) Run() {
	rawName, args := app.prepareRun()
	name := app.RealCommandName(rawName)
	Logf(VerbCrazy, "[App.Run] begin run console application, process ID: %d", app.pid)
	Logf(VerbDebug, "[App.Run] input command is: %s, real command: %s, args: %v", rawName, name, args)

	if !app.IsCommand(name) {
		color.Error.Prompt("unknown input command '%s'", name)

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
	app.fireEvent(EvtBefore, cmd.Copy())

	// parse args, don't contains command name.
	if !cmd.CustomFlags {
		if CLI.hasHelpKeywords() { // contains keywords "-h" OR "--help"
			cmd.ShowHelp(true)
		}

		cmd.Flags.Parse(args)
		args = cmd.Flags.Args()
	}

	Logf(VerbDebug, "[App.Run] args for the command '%s': %v", name, args)

	// do execute command
	exitCode := cmd.Execute(args)
	if len(app.errors) > 0 {
		app.fireEvent(EvtError, app.errors)
	} else {
		app.fireEvent(EvtAfter, exitCode)
	}

	Logf(VerbDebug, "[App.Run] command %s run complete, exit with code: %d", name, exitCode)
	Exit(exitCode)
}

// SubRun running other command in current command
func (app *App) SubRun(name string, args []string) int {
	if !app.IsCommand(name) {
		color.Error.Prompt("unknown command name '%s'", name)
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

// IsCommand name check
func (app *App) IsCommand(name string) bool {
	_, has := app.names[name]
	return has
}

// CommandName get current command name
func (app *App) CommandName() string {
	return app.commandName
}

// CommandNames get all command names
func (app *App) CommandNames() []string {
	var ss []string
	for n := range app.names {
		ss = append(ss, n)
	}

	return ss
}

// AddAliases add alias names for a command
func (app *App) AddAliases(command string, names []string) {
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
func (app *App) RealCommandName(alias string) string {
	if name, has := app.aliases[alias]; has {
		return name
	}

	return alias
}
