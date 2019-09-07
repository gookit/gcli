package gcli

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v2/helper"
	"github.com/gookit/goutil/strutil"
)

// parseGlobalOpts parse global options
func (app *App) parseGlobalOpts() []string {
	// don't display date on print log
	log.SetFlags(0)

	// bind help func
	flag.Usage = app.showApplicationHelp

	// binding global options
	flag.UintVar(&gOpts.verbose, "verbose", gOpts.verbose, "")
	flag.BoolVar(&gOpts.showHelp, "h", false, "")
	flag.BoolVar(&gOpts.showHelp, "help", false, "")
	flag.BoolVar(&gOpts.showVer, "V", false, "")
	flag.BoolVar(&gOpts.showVer, "version", false, "")
	flag.BoolVar(&gOpts.noColor, "no-color", gOpts.noColor, "")
	// this is a internal command
	flag.BoolVar(&inCompletion, "cmd-completion", false, "")

	// parse global options
	flag.Parse()

	// check global options
	if gOpts.showHelp {
		app.showApplicationHelp()
	}

	if gOpts.showVer {
		app.showVersionInfo()
	}

	if gOpts.noColor {
		color.Enable = false
	}

	Logf(VerbDebug, "[App.parseGlobalOpts] console debug is enabled, level is %d", gOpts.verbose)

	return flag.Args()
}

// prepare to running, parse args, get command name and command args
func (app *App) prepareRun() (string, []string) {
	args := app.parseGlobalOpts()
	if inCompletion {
		app.showCompletion(args)
	}

	// if no input command
	if len(args) == 0 {
		// will try run defaultCommand
		defCmd := app.defaultCommand
		if len(defCmd) == 0 {
			app.showApplicationHelp()
		}

		if !app.IsCommand(defCmd) {
			Logf(VerbError, "The default command '%s' is not exists", defCmd)
			app.showApplicationHelp()
		}

		args = []string{defCmd}
	} else if args[0] == "help" { // is help command
		if len(args) == 1 { // like 'help'
			app.showApplicationHelp()
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
	Logf(VerbDebug, "[App.Run] input command is: '%s', real command: '%s', flags: %v", rawName, name, args)

	if !app.IsCommand(name) {
		color.Error.Prompt("unknown input command '%s'", name)
		if ns := app.findSimilarCmd(name); len(ns) > 0 {
			fmt.Println("\nMaybe you mean:\n  ", color.Green.Render(strings.Join(ns, ", ")))
		}

		fmt.Printf("\nUse \"%s\" to see available commands\n", color.Cyan.Render(app.binName+" -h"))
		Exit(ERR)
	}

	code := OK
	cmd := app.commands[name]

	app.cleanArgs = args
	app.commandName = name
	// fmt.Println(cmd.argsStr, len(cmd.argsStr), strings.LastIndex(cmd.argsStr, " -h"))
	app.fireEvent(EvtBefore, cmd.Copy())

	Logf(VerbDebug, "[App.Run] command raw flags: %v", args)

	// parse options, don't contains command name.
	args = cmd.parseFlags(args, true)

	Logf(VerbDebug, "[App.Run] args on parse end: %v", args)

	// do execute command
	if err := cmd.execute(args); err != nil {
		code = ERR
		app.fireEvent(EvtError, err)
	} else {
		app.fireEvent(EvtAfter, nil)
	}

	Logf(VerbDebug, "[App.Run] command '%s' run complete, exit with code: %d", name, code)

	if app.ExitOnEnd {
		Exit(code)
	}
}

// Exec running other command in current command
func (app *App) Exec(name string, args []string) (err error) {
	if !app.IsCommand(name) {
		color.Error.Prompt("unknown command name '%s'", name)
		return
	}

	cmd := app.commands[name]
	args = cmd.parseFlags(args, false)

	// do execute command
	return cmd.execute(args)
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

/*************************************************************
 * display app help
 *************************************************************/

// help template for all commands
var commandsHelp = `{{.Description}} (Version: <info>{{.Version}}</>)
<comment>Usage:</>
  {$binName} [Global Options...] <info>{command}</> [--option ...] [argument ...]

<comment>Global Options:</>
      <info>--verbose</>     Set error reporting level(quiet 0 - 4 debug)
      <info>--no-color</>    Disable color when outputting message
  <info>-h, --help</>        Display the help information
  <info>-V, --version</>     Display app version information

<comment>Available Commands:</>{{range $module, $cs := .Cs}}{{if $module}}
<comment>{{ $module }}</>{{end}}{{ range $cs }}
  <info>{{.Name | paddingName }}</> {{.UseFor}}{{if .Aliases}} (alias: <cyan>{{ join .Aliases ","}}</>){{end}}{{end}}{{end}}

  <info>{{ paddingName "help" }}</> Display help information

Use "<cyan>{$binName} {command} -h</>" for more information about a command
`

// display app version info
func (app *App) showVersionInfo() {
	color.Printf(
		"%s\n\nVersion: <cyan>%s</>\n",
		strutil.UpperFirst(app.Description),
		app.Version,
	)

	if app.Logo.Text != "" {
		color.Printf("%s\n", color.WrapTag(app.Logo.Text, app.Logo.Style))
	}

	Exit(OK)
}

// display app help and list all commands
func (app *App) showApplicationHelp() {
	// commandsHelp = color.ReplaceTag(commandsHelp)
	// render help text template
	s := helper.RenderText(commandsHelp, map[string]interface{}{
		"Cs": app.moduleCommands,
		// app version
		"Version": app.Version,
		// always upper first char
		"Description": strutil.UpperFirst(app.Description),
	}, template.FuncMap{
		"paddingName": func(n string) string {
			return strutil.PadRight(n, " ", app.nameMaxLength)
		},
	})

	// parse help vars and render color tags
	color.Print(app.ReplaceVars(s))
	Exit(OK)
}

// showCommandHelp display help for an command
func (app *App) showCommandHelp(list []string, quit bool) {
	if len(list) != 1 {
		exitWithErr("Too many arguments given.\n\nUsage: %s help COMMAND", app.binName)
	}

	// get real name
	name := app.RealCommandName(list[0])
	if name == HelpCommand || name == "-h" {
		fmt.Printf("Display help message for application or command.\n\n")
		color.Printf("Usage: %s help COMMAND\n", app.binName)
		Exit(0)
	}

	cmd, exist := app.commands[name]
	if !exist {
		color.Error.Prompt("Unknown command name %#q. Run '%s -h'", name, app.binName)
		Exit(ERR)
	}

	// show help for the command.
	cmd.ShowHelp(quit)
}

// show bash/zsh completion
func (app *App) showCompletion(args []string) {
	// TODO ...
}

// findSimilarCmd find similar cmd by input string
func (app *App) findSimilarCmd(input string) []string {
	var ss []string
	// ins := strings.Split(input, "")
	// fmt.Print(input, ins)
	ln := len(input)

	names := app.Names()
	names["help"] = 4 // add 'help' command

	// find from command names
	for name := range names {
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
