package gcli

import (
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/events"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/helper"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/cliutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/maputil"
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

// AppConfig struct
type AppConfig struct {
	BeforeRun     func() bool
	AfterRun      func() bool
	BeforeAddOpts func(opts *Flags)
	AfterAddOpts  func(app *App) bool
}

// App the cli app definition
type App struct {
	// internal use
	// for manage commands
	base

	// AppConfig

	fs *Flags
	// cli input options for app
	opts *GlobalOpts

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
//
//	NewApp()
//	// Or with a config func
//	NewApp(func(a *App) {
//		// do something before init ....
//		a.Hooks[events.OnAppInitAfter] = func () {}
//	})
func NewApp(fns ...func(app *App)) *App {
	app := &App{
		Name: "GCliApp",
		Desc: "This is my console application",
	}

	app.fs = gflag.New(app.Name).WithConfigFn(func(opt *gflag.Config) {
		opt.WithoutType = true
		opt.IndentLongOpt = true
		opt.Alignment = gflag.AlignLeft
	})

	Logf(VerbCrazy, "create a new cli application, and create base ")

	// init
	app.base = newBase()
	app.opts = newGlobalOpts()

	// set a default value
	app.Version = "1.0.0"
	app.base.Ctx = gCtx

	for _, fn := range fns {
		fn(app)
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
//
// Notice: must be called before add command
func (app *App) Config(fns ...func(a *App)) {
	for _, fn := range fns {
		if fn != nil {
			fn(app)
		}
	}
}

/*************************************************************
 * app initialize
 *************************************************************/

// initialize application on: add, run
func (app *App) initialize() {
	if app.initialized {
		return
	}

	app.initialized = true
	app.Fire(events.OnAppInitBefore, nil)
	Logf(VerbCrazy, "initialize the cli application")

	// init some info
	app.initHelpReplacer()
	app.bindAppOpts()

	// add default error handler.
	if !app.HasHook(events.OnAppRunError) {
		app.On(events.OnAppRunError, defaultErrHandler)
	}

	app.Fire(events.OnAppInitAfter, nil)
}

// binding app options
func (app *App) bindAppOpts() {
	Logf(VerbDebug, "will begin binding app global options")
	// global options flag
	fs := app.fs
	app.Fire(events.OnAppBindOptsBefore, nil)

	// binding global options
	app.opts.bindingOpts(fs)
	// add more ...
	// This is an internal option
	fs.BoolVar(&gOpts.inCompletion, &gflag.CliOpt{
		Name: "in-completion",
		Desc: "generate completion scripts for bash/zsh",
		// hidden it
		Hidden: true,
	})

	// support binding custom global options
	app.Fire(events.OnAppBindOptsAfter, nil)
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
	app.fireWithCmd(events.OnAppCmdAdd, c, nil)

	// init command
	c.app = app
	// inherit some from application
	c.Ctx = app.Ctx
	// init for cmd flags parser
	c.Flags.Init(c.Name)

	// inherit flag parser config
	fsCfg := app.fs.ParserCfg()
	c.Flags.WithConfigFn(gflag.WithIndentLongOpt(fsCfg.IndentLongOpt))

	// do add command
	app.addCommand(app.Name, c)

	app.fireWithCmd(events.OnAppCmdAdded, c, nil)
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
 * parse global options
 *************************************************************/

// parseAppOpts parse global options
func (app *App) doParseOpts(args []string) error {
	err := app.fs.Parse(args)
	if err != nil {
		if cflag.IsFlagHelpErr(err) {
			return nil
		}
		Logf(VerbWarn, "parse global options err: <red>%s</>", err.Error())
	}

	return err
}

// parseAppOpts parse global options
func (app *App) parseAppOpts(args []string) (ok bool) {
	Logf(VerbDebug, "will begin parse app options, input-args: %v", args)

	// parse global options
	if err := app.doParseOpts(args); err != nil { // has error.
		color.Error.Tips(err.Error())
		return
	}

	app.args = app.fs.FSetArgs()
	evtData := map[string]any{"args": app.args}
	if app.Fire(events.OnAppOptsParsed, evtData) {
		Logf(VerbDebug, "stop running on the event %s return True", events.OnGlobalOptsParsed)
		return
	}

	if app.Fire(events.OnGlobalOptsParsed, evtData) {
		Logf(VerbDebug, "stop running on the event %s return True", events.OnGlobalOptsParsed)
		return
	}

	// check global options
	if app.opts.ShowHelp {
		return app.showApplicationHelp()
	}
	if app.opts.ShowVersion {
		return app.showVersionInfo()
	}

	// disable color
	if app.opts.NoColor {
		color.Enable = false
	}

	Debugf("app options parsed, verbose: <mgb>%s</>, options: %#v", app.opts.Verbose.String(), app.opts)

	// TODO show auto-completion for bash/zsh
	if app.opts.inCompletion {
		app.showAutoCompletion(app.args)
		return
	}

	return true
}

/*************************************************************
 * prepare run
 *************************************************************/

// prepare to running
//
//	parse args
//	check global options
//	get command name and command args
func (app *App) prepareRun() (code int, name string) {
	// find command name.
	name = app.findCommandName()
	if name == HelpCommand {
		if len(app.args) == 0 { // like 'help'
			app.showApplicationHelp()
		} else {
			// like 'help COMMAND'
			code = app.showCommandHelp(app.args)
		}
		return
	}

	// not input and not set defaultCommand
	if name == "" {
		if app.Func != nil {
			code = app.doRunFunc(app.args)
		} else {
			app.showApplicationHelp()
		}
		return
	}

	// name is not empty, but is not command.
	if app.inputName == "" {
		Logf(VerbDebug, "input the command is not an registered: %s", name)
		hookData := map[string]any{"name": name, "args": app.args}

		// fire events
		if stop := app.Fire(events.OnAppCmdNotFound, hookData); stop {
			return
		}
		if stop := app.Fire(events.OnCmdNotFound, hookData); stop {
			return
		}

		app.showCommandTips(name)
		return
	}

	// is valid command name.
	app.commandName = name
	return GOON, name
}

func (app *App) findCommandName() (name string) {
	args := app.args
	if len(args) == 0 {
		// not input command, will try to run app.defaultCommand
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
	if !helper.IsGoodCmdId(name) {
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

/*************************************************************
 * prepare run
 *************************************************************/

// QuickRun the application with os.Args
func (app *App) QuickRun() int {
	return app.Run(os.Args[1:])
}

// Run the application with input args
//
// Usage:
//
//	// run with os.Args
//	app.Run(nil)
//	app.Run(os.Args[1:])
//
//	// custom args
//	app.Run([]string{"cmd", "--name", "inhere"})
func (app *App) Run(args []string) (code int) {
	// ensure application initialized
	app.initialize()

	// if not set input args
	if args == nil {
		args = os.Args[1:] // exclude first arg, it's binFile.
	}

	Debugf("will begin run application. input-args: %v", args)

	// parse global flags
	if false == app.parseAppOpts(args) {
		return app.exitOnEnd(code)
	}

	Logf(VerbCrazy, "begin run console application, PID: %d", app.Ctx.PID())

	var name string
	code, name = app.prepareRun()
	if code != GOON {
		return app.exitOnEnd(code)
	}

	app.Fire(events.OnAppPrepared, map[string]any{"name": name})

	// do run input command
	var exCode int
	err := app.doRunCmd(name, app.args)
	if err != nil {
		if ec, ok := err.(errorx.ErrorCoder); ok {
			exCode = ec.Code()
		}
	}

	Debugf("command '%s' run complete, exit with code: %d", name, exCode)
	return app.exitOnEnd(exCode)
}

// RunLine manual run a command by command line string.
//
// eg: app.RunLine("top --top-opt val0 sub --sub-opt val1 arg0")
func (app *App) RunLine(argsLine string) int {
	args := cliutil.ParseLine(argsLine)
	return app.Run(args)
}

// RunCmd running a top command with custom args
//
// Usage:
//
//	app.Exec("top")
//	app.Exec("top", []string{"-a", "val0", "arg0"})
//	// can add sub command on args
//	app.Exec("top", []string{"sub", "-o", "abc"})
func (app *App) RunCmd(name string, args []string) error {
	if !app.HasCommand(name) {
		return errorx.Failf(ERR, "command %q not exists", name)
	}
	return app.doRunCmd(name, args)
}

func (app *App) doRunCmd(name string, args []string) (err error) {
	cmd := app.GetCommand(name)
	app.fireWithCmd(events.OnAppRunBefore, cmd, map[string]any{"args": args})
	Debugf("will run app command '%s' with args: %v", name, args)

	// do execute command
	if err = cmd.innerDispatch(args); err != nil {
		err = newRunErr(ERR, err)
		app.Fire(events.OnAppRunError, map[string]any{"err": err})
	} else {
		app.Fire(events.OnAppRunAfter, map[string]any{"cmd": name})
	}
	return
}

func (app *App) doRunFunc(args []string) (code int) {
	// do execute command
	if err := app.Func(app, args); err != nil {
		code = ERR
		app.Fire(events.OnAppRunError, map[string]any{"err": err})
	} else {
		app.Fire(events.OnAppRunAfter, nil)
	}
	return
}

// Exec direct exec other command in current command
//
// Name can be:
//   - top command name in the app. 'top'
//   - command path in the app. 'top sub'
//
// Usage:
//
//	app.Exec("top")
//	app.Exec("top:sub")
//	app.Exec("top sub")
//	app.Exec("top sub", []string{"-a", "val0", "arg0"})
func (app *App) Exec(path string, args []string) error {
	cmd := app.MatchByPath(path)
	if cmd == nil {
		return fmt.Errorf("exec unknown command %q", path)
	}

	Debugf("manual exec the application command: %q", path)

	// parse flags and execute command
	return cmd.innerExecute(args, false)
}

/*************************************************************
 * helper methods
 *************************************************************/

// Opts get the app GlobalOpts
func (app *App) Opts() *GlobalOpts {
	return app.opts
}

// Flags get
func (app *App) Flags() *Flags {
	return app.fs
}

// Exit get the app GlobalFlags
func (app *App) Exit(code int) {
	if app.ExitFunc == nil {
		os.Exit(code)
	}
	app.ExitFunc(code)
}

func (app *App) exitOnEnd(code int) int {
	Debugf("application exit with code: %d", code)
	app.Fire(events.OnAppExit, map[string]any{"code": code})

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

// CommandName get current command name
func (app *App) CommandName() string {
	return app.commandName
}

// SetDefaultCommand set default command name
func (app *App) SetDefaultCommand(name string) {
	app.defaultCommand = name
}

// On add hook handler for a hook event
func (app *App) On(name string, handler HookFunc) {
	Debugf("register application hook: %s", name)
	app.Hooks.On(name, handler)
}

// fire hook on the app. returns True for stop continue run.
func (app *App) fireWithCmd(event string, cmd *Command, data map[string]any) bool {
	Debugf("trigger the application event: <green>%s</>, data: %s", event, maputil.ToString(data))

	ctx := newHookCtx(event, cmd, data).WithApp(app)
	return app.Hooks.Fire(event, ctx)
}

// Fire hook on the app. returns True for stop continue run.
func (app *App) Fire(event string, data map[string]any) bool {
	Debugf("trigger the application event: <green>%s</>, data: %s", event, maputil.ToString(data))

	ctx := newHookCtx(event, nil, data).WithApp(app)
	return app.Hooks.Fire(event, ctx)
}
