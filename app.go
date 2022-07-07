package gcli

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/helper"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/strutil"
)

/*************************************************************
 * CLI application
 *************************************************************/

// Handler interface definition
type Handler interface {
	// Creator for create new command
	Creator() *Command
	// Config bind Flags or Arguments for the command
	Config(c *Command)
	// Execute the command
	Execute(c *Command, args []string) error
}

// Logo app logo, ASCII logo
type Logo struct {
	Text  string // ASCII logo string
	Style string // eg "info"
}

// App the cli app definition
type App struct {
	// internal use
	// - *cmdLine
	// - HelpVars
	// - Hooks // allow hooks: "init", "before", "after", "error"
	core
	// for manager commands
	commandBase

	// Name app name
	Name string
	// Desc app description
	Desc string
	// Func on run app, if is empty will display help.
	Func func(app *App, args []string) error
	// ExitOnEnd call os.Exit on running end
	// ExitOnEnd bool
	// ExitFunc default is os.Exit
	// ExitFunc func(int)

	// args on after parse global options and command name.
	args []string
	// moduleCommands map[string]map[string]*Command

	// rawFlagArgs []string
	// clean os.args, not contains bin-name and command-name
	cleanArgs []string
	// has sub-commands on the app
	hasSubcommands bool
	// the default command name.
	// if is empty, will render help message.
	defaultCommand string
}

// New alias of the NewApp()
func New(fns ...func(app *App)) *App {
	return NewApp(fns...)
}

// NewApp create new app instance.
//
// Usage:
// 	NewApp()
// 	// Or with a config func
// 	NewApp(func(a *App) {
// 		// do something before init ....
// 		a.Hooks[gcli.EvtInit] = func () {}
// 	})
func NewApp(fns ...func(app *App)) *App {
	app := &App{
		Name: "GCliApp",
		Desc: "This is my console application",
		// set a default version
		// Version: "1.0.0",
		// config
		// ExitOnEnd: true,
		// group
		// moduleCommands: make(map[string]map[string]*Command),
		commandBase: newCommandBase(),
	}

	// internal core
	Logf(VerbCrazy, "create new core on init application")
	app.core = core{
		cmdLine: CLI,
		// init
		Hooks: &Hooks{},
		gFlags: NewFlags("app.GOptions").WithConfigFn(func(opt *FlagsConfig) {
			opt.WithoutType = true
			opt.Alignment = AlignLeft
		}),
	}

	// init commandBase
	Logf(VerbCrazy, "create new commandBase on init application")
	// set a default version
	app.Version = "1.0.0"

	if len(fns) > 0 {
		for _, fn := range fns {
			fn(app)
		}
	}
	return app
}

// NotExitOnEnd for app
func NotExitOnEnd() func(*App) {
	return func(app *App) {
		app.ExitOnEnd = false
	}
}

// Config the application.
// Notice: must be called before adding a command
func (app *App) Config(fn func(a *App)) {
	if fn != nil {
		fn(app)
	}
}

// binding global options
func (app *App) bindingGlobalOpts() {
	Logf(VerbDebug, "will begin binding global options")
	// global options flag
	// gf := flag.NewFlagSet(app.Args[0], flag.ContinueOnError)
	gf := app.GlobalFlags()

	// binding global options
	// bindingCommonGOpts(gf)
	gOpts.bindingFlags(gf)
	// add more ...
	gf.BoolOpt(&gOpts.showVer, "version", "V", false, "Display app version information")
	// This is a internal option
	gf.BoolVar(&gOpts.inCompletion, &FlagMeta{
		Name: "in-completion",
		Desc: "generate completion scripts for bash/zsh",
		// hidden it
		Hidden: true,
	})

	// support binding custom global options
	if app.GOptsBinder != nil {
		app.GOptsBinder(gf)
	}
}

// initialize application
func (app *App) initialize() {
	if app.initialized {
		return
	}

	Logf(VerbCrazy, "initialize the application")

	// init some vars
	if app.core.Hooks == nil {
		app.core.Hooks = &Hooks{}
	}
	app.core.AddVars(app.core.innerHelpVars())

	// binding global options
	app.bindingGlobalOpts()

	// add default error handler.
	if !app.HasHook(EvtAppRunError) {
		app.On(EvtAppRunError, defaultErrHandler)
	}

	app.Fire(EvtAppInit, nil)
	app.initialized = true
}

/*************************************************************
 * register commands
 *************************************************************/

// Add one or multi command(s)
func (app *App) Add(c *Command, more ...*Command) {
	app.AddCommand(c)

	// has more command
	if len(more) > 0 {
		for _, cmd := range more {
			app.AddCommand(cmd)
		}
	}
}

// AddCommand add a new command to the app
func (app *App) AddCommand(c *Command) {
	// initialize application before add command
	app.initialize()

	// init command
	c.app = app
	// inherit global flags from application
	c.core.gFlags = app.gFlags

	// do add command
	app.commandBase.addCommand(app.Name, c)

	if c.HasCommands() {
		app.hasSubcommands = true
	}

	app.Fire(EvtCmdInit, c)
}

// AddHandler to the application
func (app *App) AddHandler(h Handler) {
	c := h.Creator()
	c.Func = h.Execute

	// binding flags
	h.Config(c)

	// add
	app.AddCommand(c)
}

// AddAliases add alias names for a command
func (app *App) AddAliases(name string, aliases ...string) {
	c := app.FindByPath(name)
	if c == nil {
		panicf("the command '%s' is not exists", name)
	}

	// add alias
	for _, alias := range aliases {
		if app.IsCommand(alias) {
			panicf("The name '%s' has been used as an command name", alias)
		}

		app.cmdAliases.AddAlias(name, alias)
	}
}

// On add hook handler for a hook event
// func (app *App) BeforeInit(name string, handler HookFunc) {}

// stop application and exit
// func stop(code int) {
// 	os.Exit(code)
// }

/*************************************************************
 * run command
 *************************************************************/

// parseGlobalOpts parse global options
func (app *App) parseGlobalOpts(args []string) (ok bool) {
	Logf(VerbDebug, "will begin parse global options")

	// parse global options
	err := app.core.doParseGOpts(args)
	if err != nil { // has error.
		color.Error.Tips(err.Error())
		return
	}

	app.args = app.gFlags.FSetArgs()
	app.Fire(EvtGOptionsParsed, app.args)

	// check global options
	if gOpts.showHelp {
		app.showApplicationHelp()
		return
	}

	if gOpts.showVer {
		app.showVersionInfo()
		return
	}

	// disable color
	if gOpts.NoColor {
		color.Enable = false
	}

	Debugf("global option parsed, verbose level: <mgb>%s</>", gOpts.verbose.String())

	// TODO show auto-completion for bash/zsh
	if gOpts.inCompletion {
		app.showAutoCompletion(app.args)
		return
	}

	return true
}

// prepare to running, parse args, get command name and command args
func (app *App) prepareRun() (code int, name string) {
	// find command name.
	name = app.findCommandName()
	// is help command name.
	if name == HelpCommand {
		if len(app.args) == 0 { // like 'help'
			app.showApplicationHelp()
			return
		}

		// like 'help COMMAND'
		code = app.showCommandHelp(app.args)
		return
	}

	// not input and not set defaultCommand
	if name == "" {
		// run app.Func
		if app.Func != nil {
			code = app.doRunFunc(app.args)
			return
		}

		app.showApplicationHelp()
		return
	}

	// name is not empty, but is not command.
	if app.inputName == "" {
		Logf(VerbDebug, "input the command is not an registered: %s", name)

		if stop := app.Fire(EvtAppCmdNotFound, name); stop == false {
			stop = app.Fire(EvtCmdNotFound, name)
			if stop == false {
				app.showCommandTips(name)
			}
		}
		return
	}

	// is valid command name.
	app.commandName = name
	return GOON, name
}

func (app *App) findCommandName() (name string) {
	args := app.args
	// not input command, will try run app.defaultCommand
	if len(args) == 0 {
		name = app.defaultCommand
		if name == "" {
			return
		}

		// It is not an valid command name. TODO default is command ID.
		if false == app.IsCommand(name) {
			Logf(VerbError, "the default command '<cyan>%s</>' is invalid", name)
			return "" // invalid, return empty string.
		}
		return name
	}

	name = strings.TrimSpace(args[0])
	// is empty string or is an option
	if name == "" || name[0] == '-' {
		return ""
	}

	// check is valid ID/name string.
	if !goodCmdId.MatchString(name) {
		Logf(VerbWarn, "the input command name(%s) string is invalid", name)
		return ""
	}

	rawName := name
	nodes := splitPath2names(name)
	// Is command ID. eg: "top:sub"
	if len(nodes) > 1 {
		name = nodes[0]
		name = app.ResolveAlias(name)
		Debugf("input(args[0]) is an command ID, expand it. '%s' -> '%s'", rawName, name)
	} else {
		rName := app.ResolveAlias(name)
		nodes = splitPath2names(rName)
		// Is command ID. eg: "top:sub"
		if len(nodes) > 1 {
			name = nodes[0]
			Debugf("real command is an command ID, expand it. '%s' -> '%s'", rName, name)
		} else {
			name = rName
		}
	}

	// update app.args
	if len(nodes) > 1 {
		app.args = append(nodes[1:], args[1:]...)
	} else {
		app.args = args[1:]
	}

	// it is exists command name.
	if app.IsCommand(name) {
		app.inputName = rawName
		Debugf("the raw input command: '<cyan>%s</>'; real name: '<green>%s</>', args: %v", rawName, name, app.args)
		return name
	}

	// not exists
	Logf(VerbInfo, "the input command name '%s' is not exists. nodes: %v", rawName, nodes)
	return rawName
}

// RunLine manual run an command by command line string.
//
// eg: app.RunLine("top --top-opt val0 sub --sub-opt val1 arg0")
func (app *App) RunLine(argsLine string) int {
	args := cliutil.ParseLine(argsLine)
	return app.Run(args)
}

// Run running application
//
// Usage:
//	// run with os.Args
//	app.Run(nil)
//	app.Run(os.Args[1:])
//	// custom args
//	app.Run([]string{"cmd", ...})
func (app *App) Run(args []string) (code int) {
	// ensure application initialized
	app.initialize()

	// if not set input args
	if args == nil {
		args = os.Args[1:] // exclude first arg, it's binFile.
	}

	Debugf("will begin run cli application. args: %v", args)

	// parse global flags
	if false == app.parseGlobalOpts(args) {
		return app.exitOnEnd(code)
	}

	Logf(VerbCrazy, "begin run console application, PID: %d", app.PID())

	var name string
	code, name = app.prepareRun()
	if code != GOON {
		return app.exitOnEnd(code)
	}

	// trigger event
	app.Fire(EvtAppPrepareAfter, name)

	// do run input command
	code = app.doRunCmd(name, app.args)

	Debugf("command '%s' run complete, exit with code: %d", name, code)
	return app.exitOnEnd(code)
}

// RunCmd running an top command with custom args
//
// Usage:
//	app.Exec("top")
//	app.Exec("top", []string{"-a", "val0", "arg0"})
//	// can add sub command on args
//	app.Exec("top", []string{"sub", "-o", "abc"})
func (app *App) RunCmd(name string, args []string) int {
	return app.doRunCmd(name, args)
}

func (app *App) doRunCmd(name string, args []string) (code int) {
	cmd := app.GetCommand(name)
	app.Fire(EvtAppRunBefore, cmd)

	Debugf("will run app command '%s' with args: %v", name, args)

	// parse command options
	// args, err = cmd.parseOptions(args)

	// do execute command
	// if err := cmd.innerExecute(args, true); err != nil {
	if err := cmd.innerDispatch(args); err != nil {
		code = ERR
		app.Fire(EvtAppRunError, err)
	} else {
		app.Fire(EvtAppRunAfter, nil)
	}
	return
}

func (app *App) doRunFunc(args []string) (code int) {
	// app bind args TODO
	// app.ParseArgs(args)

	// do execute command
	if err := app.Func(app, args); err != nil {
		code = ERR
		app.Fire(EvtAppRunError, err)
	} else {
		app.Fire(EvtAppRunAfter, nil)
	}
	return
}

func (app *App) exitOnEnd(code int) int {
	Debugf("application exit with code: %d", code)
	// if IsGteVerbose(VerbDebug) {
	// 	app.Infoln("[DEBUG] The Runtime Call Stacks:")

	// bts := goutil.GetCallStacks(true)
	// app.Println(string(bts), len(bts))
	// cs := goutil.GetCallersInfo(2, 10)
	// app.Println(strings.Join(cs, "\n"), len(cs))
	// }

	if app.ExitOnEnd {
		app.Exit(code)
	}
	return code
}

// Exec direct exec other command in current command
//
// name can be:
// - top command name in the app. 'top'
// - command path in the app. 'top sub'
//
// Usage:
//	app.Exec("top")
//	app.Exec("top:sub")
//	app.Exec("top sub")
//	app.Exec("top sub", []string{"-a", "val0", "arg0"})
func (app *App) Exec(path string, args []string) error {
	cmd := app.MatchByPath(path)
	if cmd == nil {
		return fmt.Errorf("exec unknown command: '%s'", path)
	}

	Debugf("manual exec the application command: %s", path)

	// parse flags and execute command
	return cmd.innerExecute(args, false)
}

// ExecLine manual execute an command by command line string.
// eg: app.ExecLine("top --top-opt val0 sub --sub-opt val1 arg0")
// func (app *App) ExecLine(argsLine string) error {
// 	args := cliutil.ParseLine(argsLine)
// }

/*************************************************************
 * helper methods
 *************************************************************/

// Exit get the app GlobalFlags
func (app *App) Exit(code int) {
	if app.ExitFunc == nil {
		os.Exit(code)
	}

	app.ExitFunc(code)
}

// CommandName get current command name
func (app *App) CommandName() string {
	return app.commandName
}

// HasSubcommands on the app
func (app *App) HasSubcommands() bool {
	return app.hasSubcommands
}

// SetDefaultCommand set default command name
func (app *App) SetDefaultCommand(name string) {
	app.defaultCommand = name
}

// On add hook handler for a hook event
func (app *App) On(name string, handler HookFunc) {
	Debugf("register application hook: %s", name)

	app.core.On(name, handler)
}

// Fire hook on the app
func (app *App) Fire(event string, data interface{}) bool {
	Debugf("trigger the application event: <green>%s</>", event)

	return app.core.Fire(event, app, data)
}

/*************************************************************
 * display app help
 *************************************************************/

// display app version info
func (app *App) showVersionInfo() {
	Debugf("print application version info")

	color.Printf(
		"%s\n\nVersion: <cyan>%s</>\n",
		strutil.UpperFirst(app.Desc),
		app.Version,
	)

	if app.Logo.Text != "" {
		color.Printf("%s\n", color.WrapTag(app.Logo.Text, app.Logo.Style))
	}
}

// display unknown input command and similar commands tips
func (app *App) showCommandTips(name string) {
	Debugf("will find and show similar command tips")

	color.Error.Tips(`unknown input command "<mga>%s</>"`, name)
	if ns := app.findSimilarCmd(name); len(ns) > 0 {
		color.Printf("\nMaybe you mean:\n  <green>%s</>\n", strings.Join(ns, ", "))
	}

	color.Printf("\nUse <cyan>%s --help</> to see available commands\n", app.binName)
}

// AppHelpTemplate help template for app(all commands)
// TODO {{if .HasSubcommands }}  {$binName} [global options...] <info>COMMAND</> [--options ...] <info>SUBCOMMAND</> [--options ...] [arguments ...]
// {{end}}
var AppHelpTemplate = `{{.Desc}} (Version: <info>{{.Version}}</>)
<comment>Usage:</>
  {$binName} [global options...] <info>COMMAND</> [--options ...] [arguments ...]
  {$binName} [global options...] <info>COMMAND</> [--options ...] <info>SUBCOMMAND</> [--options ...] [arguments ...]

<comment>Global Options:</>
{{.GOpts}}
<comment>Available Commands:</>{{range $cmdName, $c := .Cs}}{{if $c.Visible}}
  <info>{{$c.Name | paddingName }}</> {{$c.HelpDesc}}{{if $c.Aliases}} (alias: <green>{{ join $c.Aliases ","}}</>){{end}}{{end}}{{end}}
  <info>{{ paddingName "help" }}</> Display help information

Use "<cyan>{$binName} COMMAND -h</>" for more information about a command
`

// display app help and list all commands. showCommandList()
func (app *App) showApplicationHelp() {
	Debugf("render application help and commands list")

	// cmdHelpTemplate = color.ReplaceTag(cmdHelpTemplate)
	// render help text template
	s := helper.RenderText(AppHelpTemplate, map[string]interface{}{
		"Cs":    app.commands,
		"GOpts": app.gFlags.String(),
		// app version
		"Version": app.Version,
		// always upper first char
		"Desc": strutil.UpperFirst(app.Desc),
	}, template.FuncMap{
		"paddingName": func(n string) string {
			return strutil.PadRight(n, " ", app.nameMaxWidth)
		},
	})

	// parse help vars and render color tags
	color.Print(app.ReplaceVars(s))
}

// showCommandHelp display help for an command
func (app *App) showCommandHelp(list []string) (code int) {
	binName := app.binName
	// if len(list) == 0 { TODO support multi level sub command?
	if len(list) > 1 {
		color.Error.Tips("Too many arguments given.\n\nUsage: %s help COMMAND", binName)
		return ERR
	}

	// get real name
	name := app.cmdAliases.ResolveAlias(list[0])
	if name == HelpCommand || name == "-h" {
		Debugf("render help command information")

		color.Println("Display help message for application or command.\n")
		color.Printf(`<yellow>Usage:</>
  <cyan>%s COMMAND --help</>
  <cyan>%s COMMAND SUBCOMMAND --help</>
  <cyan>%s COMMAND SUBCOMMAND ... --help</>
  <cyan>%s help COMMAND</>
`, binName, binName, binName, binName)
		return
	}

	cmd, exist := app.Command(name)
	if !exist {
		color.Error.Prompt("Unknown command name '%s'. Run '%s -h' see all commands", name, binName)
		return ERR
	}

	// show help for the give command.
	cmd.ShowHelp()
	return
}

// show bash/zsh completion
func (app *App) showAutoCompletion(_ []string) {
	// TODO ...
}

// findSimilarCmd find similar cmd by input string
func (app *App) findSimilarCmd(input string) []string {
	var ss []string
	// ins := strings.Split(input, "")
	// fmt.Print(input, ins)
	ln := len(input)

	names := app.CmdNameMap()
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
	for alias := range app.cmdAliases.Mapping() {
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
