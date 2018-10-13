// Package cliapp is a simple to use command line application, written using golang
//
// Source code and other details for the project are available at GitHub:
// 		https://github.com/gookit/cliapp
//
// usage please ref examples and README
package cliapp

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// constants for error level 0 - 4
const (
	VerbQuiet uint = iota // don't report anything
	VerbError             // reporting on error
	VerbWarn
	VerbInfo
	VerbDebug
	VerbCrazy
)

// HelpVar allow var replace in help info.
// default support:
// 	"{$binName}" "{$cmd}" "{$fullCmd}" "{$workDir}"
const HelpVar = "{$%s}"

// constants for hooks event, there are default allowed hook names
const (
	EvtInit   = "init"
	EvtBefore = "before"
	EvtAfter  = "after"
	EvtError  = "error"
)

const (
	// OK success exit code
	OK = 0
	// ERR error exit code
	ERR = 2
)

/*************************************************************
 * Command Line: command data
 *************************************************************/

// CmdLine store common data for CLI
type CmdLine struct {
	// pid for current application
	pid int
	// os name.
	osName string
	// the CLI app work dir path. by `os.Getwd()`
	workDir string
	// bin script name, by `os.Args[0]`. eg "./cliapp"
	binName string
	// os.Args to string, but no binName.
	argsStr string
}

// PID get PID
func (c *CmdLine) PID() int {
	return c.pid
}

// OsName is equals to `runtime.GOOS`
func (c *CmdLine) OsName() string {
	return c.osName
}

// BinName get bin script name
func (c *CmdLine) BinName() string {
	return c.binName
}

// WorkDir get work dir
func (c *CmdLine) WorkDir() string {
	return c.workDir
}

// ArgsString os.Args to string, but no binName.
func (c *CmdLine) ArgsString() string {
	return c.argsStr
}

// CLI create a default instance
var CLI = &CmdLine{
	pid: os.Getpid(),
	// more info
	osName:  runtime.GOOS,
	binName: os.Args[0],
	argsStr: strings.Join(os.Args[1:], " "),
}

/*************************************************************
 * CLI application
 *************************************************************/

// Logo app logo, ASCII logo
type Logo struct {
	Text  string // ASCII logo string
	Style string // eg "info"
}

type appHookFunc func(app *Application, data interface{})

// Application the cli app definition
type Application struct {
	// internal use
	*CmdLine
	// Name app name
	Name string
	// Version app version. like "1.0.1"
	Version string
	// Description app description
	Description string
	// Logo ASCII logo setting
	Logo Logo
	// Hooks can setting some hooks func on running.
	// allow hooks: "init", "before", "after", "error"
	Hooks map[string]appHookFunc
	// Strict use strict mode. short opt must be begin '-', long opt must be begin '--'
	Strict bool
	// vars you can add some vars map for render help info
	vars map[string]string
	// command names. key is name, value is name string length
	// eg. {"test": 4, "example": 7}
	names map[string]int
	// store some runtime errors
	errors []error
	// command aliases map. {alias: name}
	aliases map[string]string
	// all commands for the app
	commands map[string]*Command
	// current command name
	commandName string
	// default command name
	defaultCommand string
}

// GlobalOpts global flags
type GlobalOpts struct {
	noColor  bool
	verbose  uint // message report level
	showVer  bool
	showHelp bool
}

// global options
var gOpts = &GlobalOpts{verbose: VerbError}

// Exit program
func Exit(code int) {
	os.Exit(code)
}

// Verbose returns verbose level
func Verbose() uint {
	return gOpts.verbose
}

// NewApp create new app instance. alias of the New()
// eg:
// 	cliapp.New()
// 	cliapp.New(func(a *Application) {
// 		// do something before init ....
// 		a.Hooks[cliapp.EvtInit] = func () {}
// 	})
func NewApp(fn ...func(a *Application)) *Application {
	return New(fn...)
}

// New create new app instance.
// eg:
// 	cliapp.NewApp()
// 	cliapp.NewApp(func(a *Application) {
// 		// do something before init ....
// 		a.Hooks[cliapp.EvtInit] = func () {}
// 	})
func New(fn ...func(a *Application)) *Application {
	app = &Application{
		Name:  "My CLI Application",
		Logo:  Logo{Style: "info"},
		Hooks: make(map[string]appHookFunc, 0),
		// set a default version
		Version:  "1.0.0",
		CmdLine:  CLI,
		commands: make(map[string]*Command),
	}

	if len(fn) > 0 {
		fn[0](app)
	}

	// init
	app.Initialize()

	return app
}

// Initialize application
func (app *Application) Initialize() {
	app.names = make(map[string]int)

	// init some tpl vars
	app.vars = map[string]string{
		"pid":     fmt.Sprint(CLI.pid),
		"workDir": CLI.workDir,
		"binName": CLI.binName,
	}

	// parse GlobalOpts
	// parseGlobalOpts()

	app.callHook(EvtInit, nil)
}

// SetLogo text and color style
func (app *Application) SetLogo(logo string, style ...string) {
	app.Logo.Text = logo

	if len(style) > 0 {
		app.Logo.Style = style[0]
	}
}

// DebugMode level
func (app *Application) DebugMode() {
	gOpts.verbose = VerbDebug
}

// QuietMode level
func (app *Application) QuietMode() {
	gOpts.verbose = VerbQuiet
}

// SetVerbose level
func (app *Application) SetVerbose(verbose uint) {
	gOpts.verbose = verbose
}

// DefaultCommand set default command name
func (app *Application) DefaultCommand(name string) {
	app.defaultCommand = name
}

// Add add a command
func (app *Application) Add(c *Command, more ...*Command) {
	if app.commands == nil {
		app.commands = make(map[string]*Command)
	}

	app.addCommand(c)

	// if has more
	if len(more) > 0 {
		for _, cmd := range more {
			app.addCommand(cmd)
		}
	}
}

func (app *Application) addCommand(c *Command) {
	c.Name = strings.TrimSpace(c.Name)
	if c.Name == "" {
		exitWithErr("The added command must have a command name")
	}

	if c.IsDisabled() {
		Logf(VerbDebug, "command %s has been disabled, skip add", c.Name)
		return
	}

	app.names[c.Name] = len(c.Name)
	app.commands[c.Name] = c
	// add aliases for the command
	app.AddAliases(c.Name, c.Aliases)
	Logf(VerbDebug, "add command: %s", c.Name)

	// init command
	c.app = app
	c.initialize()
}

func (app *Application) callHook(event string, data interface{}) {
	Logf(VerbDebug, "application trigger the hook: %s", event)

	if handler, ok := app.Hooks[event]; ok {
		handler(app, data)
	}
}

// On add hook handler for a hook event
func (app *Application) On(name string, handler func(a *Application, data interface{})) {
	app.Hooks[name] = handler
}

// AddError to the application
func (app *Application) AddError(err error) {
	app.errors = append(app.errors, err)
}

// AddVar get command name
func (app *Application) AddVar(name, value string) {
	app.vars[name] = value
}

// AddVars add multi tpl vars
func (app *Application) AddVars(vars map[string]string) {
	for n, v := range vars {
		app.AddVar(n, v)
	}
}

// GetVar get a help var by name
func (app *Application) GetVar(name string) string {
	if v, ok := app.vars[name]; ok {
		return v
	}

	return ""
}

// GetVars get all tpl vars
func (app *Application) GetVars(name string, value string) map[string]string {
	return app.vars
}

// Commands get all commands
func (app *Application) Commands() map[string]*Command {
	return app.commands
}
