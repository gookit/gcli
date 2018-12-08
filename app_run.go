package gcli

import (
	"flag"
	"fmt"
	"github.com/gookit/color"
	"github.com/gookit/gcli/helper"
	"github.com/gookit/goutil/envUtil"
	"github.com/gookit/goutil/strUtil"
	"log"
	"os"
	"runtime"
	"strings"
)

var (
	// global options
	gOpts = &GlobalOpts{}
	// command auto completion mode.
	// eg "./cli --cmd-completion [COMMAND --OPT ARG]"
	inCompletion bool
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
	if envUtil.IsWin() {
		CLI.binName = strings.Replace(CLI.binName, workDir+"\\", "", 1)
	}
}

// parseGlobalOpts parse global options
func (app *App) parseGlobalOpts() []string {
	// don't display date on print log
	log.SetFlags(0)

	// bind help func
	flag.Usage = app.showCommandsHelp

	// binding global options
	flag.UintVar(&gOpts.verbose, "verbose", VerbError, "")
	flag.BoolVar(&gOpts.showHelp, "h", false, "")
	flag.BoolVar(&gOpts.showHelp, "help", false, "")
	flag.BoolVar(&gOpts.showVer, "V", false, "")
	flag.BoolVar(&gOpts.showVer, "version", false, "")
	flag.BoolVar(&gOpts.noColor, "no-color", false, "")
	// this is a internal command
	flag.BoolVar(&inCompletion, "cmd-completion", false, "")

	// parse global options
	flag.Parse()

	// check global options
	if gOpts.showHelp {
		app.showCommandsHelp()
	}

	if gOpts.showVer {
		app.showVersionInfo()
	}

	if gOpts.noColor {
		color.Enable = false
	}

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

		if err := cmd.Flags.Parse(args); err != nil {
			color.Error.Prompt("Flags parse error: %s", err.Error())
			Exit(ERR)
		}

		args = cmd.Flags.Args()
	}

	Logf(VerbDebug, "[App.Run] args for the command '%s': %v", name, args)

	// do execute command
	err := cmd.Execute(args)
	exitCode := 0

	if err != nil {
		exitCode = ERR
		app.fireEvent(EvtError, err)
	} else {
		app.fireEvent(EvtAfter, nil)
	}

	Logf(VerbDebug, "[App.Run] command %s run complete, exit with code: %d", name, exitCode)
	Exit(exitCode)
}

// Exec running other command in current command
func (app *App) Exec(name string, args []string) (err error) {
	if !app.IsCommand(name) {
		color.Error.Prompt("unknown command name '%s'", name)
		return
	}

	cmd := app.commands[name]
	if !cmd.CustomFlags {
		// parse args, don't contains command name.
		if err = cmd.Flags.Parse(args); err != nil {
			color.Error.Prompt("Flags parse error: %s", err.Error())
			return
		}

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

<comment>Available Commands:</>{{range .Cs}}{{if .Runnable}}
  <info>{{.Name | printf "%-12s"}}</> {{.UseFor}}{{if .Aliases}} (alias: <cyan>{{ join .Aliases ","}}</>){{end}}{{end}}{{end}}
  <info>help</>         Display help information

Use "<cyan>{$binName} {command} -h</>" for more information about a command
`

// display app version info
func (app *App) showVersionInfo() {
	fmt.Printf(
		"%s\n\nVersion: %s\n",
		strUtil.UpperFirst(app.Description),
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
	s := helper.RenderText(commandsHelp, map[string]interface{}{
		"Cs": app.commands,
		// app version
		"Version": app.Version,
		// always upper first char
		"Description": strUtil.UpperFirst(app.Description),
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
