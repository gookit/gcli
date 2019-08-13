// Package gcli is a simple to use command line application, written using golang
//
// Source code and other details for the project are available at GitHub:
// 		https://github.com/gookit/gcli
//
// usage please ref examples and README
package gcli

import (
	"fmt"
	"os"
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

// constants for hooks event, there are default allowed event names
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
	argLine string
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

// ArgLine os.Args to string, but no binName.
func (c *CmdLine) ArgLine() string {
	return c.argLine
}

func (c *CmdLine) helpVars() map[string]string {
	return map[string]string{
		"pid":     fmt.Sprint(CLI.pid),
		"workDir": CLI.workDir,
		"binName": CLI.binName,
	}
}

func (c *CmdLine) hasHelpKeywords() bool {
	return strings.HasSuffix(c.argLine, " -h") || strings.HasSuffix(c.argLine, " --help")
}

/*************************************************************
 * CLI application
 *************************************************************/

// Logo app logo, ASCII logo
type Logo struct {
	Text  string // ASCII logo string
	Style string // eg "info"
}

type appHookFunc func(app *App, data interface{})

// App the cli app definition
type App struct {
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
// 	gcli.New()
// 	gcli.New(func(a *App) {
// 		// do something before init ....
// 		a.Hooks[gcli.EvtInit] = func () {}
// 	})
func NewApp(fn ...func(a *App)) *App {
	return New(fn...)
}

// New create new app instance.
// eg:
// 	gcli.NewApp()
// 	gcli.NewApp(func(a *App) {
// 		// do something before init ....
// 		a.Hooks[gcli.EvtInit] = func () {}
// 	})
func New(fn ...func(a *App)) *App {
	defApp = &App{
		Name:  "My CLI App",
		Logo:  Logo{Style: "info"},
		Hooks: make(map[string]appHookFunc, 0),
		// set a default version
		Version:  "1.0.0",
		CmdLine:  CLI,
		commands: make(map[string]*Command),
	}

	if len(fn) > 0 {
		fn[0](defApp)
	}

	// init
	defApp.Initialize()
	return defApp
}

// Config the application.
// Notice: must be called before adding a command
func (app *App) Config(fn func(a *App)) {
	if fn != nil {
		fn(app)
	}
}

// Initialize application
func (app *App) Initialize() {
	app.names = make(map[string]int)

	// init some help tpl vars
	app.vars = CLI.helpVars()

	// parse GlobalOpts
	// parseGlobalOpts()

	app.fireEvent(EvtInit, nil)
}

// SetLogo text and color style
func (app *App) SetLogo(logo string, style ...string) {
	app.Logo.Text = logo
	if len(style) > 0 {
		app.Logo.Style = style[0]
	}
}

// DebugMode level
func (app *App) DebugMode() {
	gOpts.verbose = VerbDebug
}

// QuietMode level
func (app *App) QuietMode() {
	gOpts.verbose = VerbQuiet
}

// SetVerbose level
func (app *App) SetVerbose(verbose uint) {
	gOpts.verbose = verbose
}

// DefaultCommand set default command name
func (app *App) DefaultCommand(name string) {
	app.defaultCommand = name
}

// NewCommand create a new command
func (app *App) NewCommand(name, useFor string, config func(c *Command)) *Command {
	return NewCommand(name, useFor, config)
}

// Add add one or multi command(s)
func (app *App) Add(c *Command, more ...*Command) {
	app.AddCommand(c)

	// if has more command
	if len(more) > 0 {
		for _, cmd := range more {
			app.AddCommand(cmd)
		}
	}
}

// AddCommand add a new command
func (app *App) AddCommand(c *Command) *Command {
	c.Name = strings.TrimSpace(c.Name)
	if c.Name == "" {
		exitWithErr("The added command name can not be empty.")
	}

	if c.IsDisabled() {
		Logf(VerbDebug, "command %s has been disabled, skip add", c.Name)
		return nil
	}

	app.names[c.Name] = len(c.Name)
	app.commands[c.Name] = c
	// add aliases for the command
	app.AddAliases(c.Name, c.Aliases)
	Logf(VerbDebug, "[App.AddCommand] add a new CLI command: %s", c.Name)

	// init command
	c.app = app
	c.initialize()
	return c
}

func (app *App) fireEvent(event string, data interface{}) {
	Logf(VerbDebug, "trigger the application event: %s", event)

	if handler, ok := app.Hooks[event]; ok {
		handler(app, data)
	}
}

// On add hook handler for a hook event
func (app *App) On(name string, handler func(a *App, data interface{})) {
	app.Hooks[name] = handler
}

// AddError to the application
func (app *App) AddError(err error) {
	app.errors = append(app.errors, err)
}

// AddVar get command name
func (app *App) AddVar(name, value string) {
	app.vars[name] = value
}

// AddVars add multi tpl vars
func (app *App) AddVars(vars map[string]string) {
	for n, v := range vars {
		app.AddVar(n, v)
	}
}

// GetVar get a help var by name
func (app *App) GetVar(name string) string {
	if v, ok := app.vars[name]; ok {
		return v
	}

	return ""
}

// GetVars get all tpl vars
func (app *App) GetVars(name string, value string) map[string]string {
	return app.vars
}

// Commands get all commands
func (app *App) Commands() map[string]*Command {
	return app.commands
}
