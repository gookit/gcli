package gcli

import (
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

// constants for hooks event, there are default allowed event names
const (
	EvtInit   = "init"
	EvtBefore = "before"
	EvtAfter  = "after"
	EvtError  = "error"
	// EvtStop   = "stop"
)

const (
	// OK success exit code
	OK = 0
	// ERR error exit code
	ERR = 2
	// GOON prepare run successful, goon run command
	GOON = -1
	// HelpCommand name
	HelpCommand = "help"
)

/*************************************************************
 * CLI application
 *************************************************************/

// HookFunc definition.
// func arguments:
//  in app, like: func(app *App, data interface{})
//  in cmd, like: func(cmd *Command, data interface{})
// type HookFunc func(obj interface{}, data interface{})
type HookFunc func(obj ...interface{})

// Logo app logo, ASCII logo
type Logo struct {
	Text  string // ASCII logo string
	Style string // eg "info"
}

// App the cli app definition
type App struct {
	// internal use
	*CmdLine
	HelpVars
	SimpleHooks // allow hooks: "init", "before", "after", "error"

	// Name app name
	Name string
	// Version app version. like "1.0.1"
	Version string
	// Description app description
	Description string
	// Logo ASCII logo setting
	Logo Logo
	// Args default is equals to os.args
	Args []string
	// ExitOnEnd call os.Exit on running end
	ExitOnEnd bool
	// vars you can add some vars map for render help info
	// vars map[string]string
	// command names. key is name, value is name string length
	// eg. {"test": 4, "example": 7}
	names map[string]int
	// store some runtime errors
	errors []error
	// command aliases map. {alias: name}
	aliases map[string]string
	// all commands for the app
	commands map[string]*Command
	// all commands by module
	moduleCommands map[string]map[string]*Command
	// the max length for added command names. default set 12.
	nameMaxLen int
	// default command name
	defaultCommand string
	// raw input command name
	rawName string
	rawFlagArgs []string
	// clean os.args, not contains bin-name and command-name
	cleanArgs []string
	// current command name
	commandName string
	// Whether it has been initialized
	initialized bool
}

// NewApp create new app instance.
// Usage:
// 	NewApp()
// 	// Or with a config func
// 	NewApp(func(a *App) {
// 		// do something before init ....
// 		a.Hooks[gcli.EvtInit] = func () {}
// 	})
func NewApp(fn ...func(a *App)) *App {
	app := &App{
		Args: os.Args,
		Name: "GCLI App",
		Logo: Logo{Style: "info"},
		// set a default version
		Version: "1.0.0",
		CmdLine: CLI,
		// config
		ExitOnEnd: true,
		// commands
		commands:       make(map[string]*Command),
		moduleCommands: make(map[string]map[string]*Command),
		// some default values
		nameMaxLen:  12,
		Description: "This is my CLI application",
	}

	if len(fn) > 0 {
		fn[0](app)
	}

	return app
}

// Config the application.
// Notice: must be called before adding a command
func (app *App) Config(fn func(a *App)) {
	if fn != nil {
		fn(app)
	}
}

// initialize application
func (app *App) initialize() {
	app.names = make(map[string]int)

	// init some help tpl vars
	app.AddVars(app.helpVars())

	// parse GlobalOpts
	// parseGlobalOpts()

	// add default error handler.
	app.SimpleHooks.AddOn(EvtError, defaultErrHandler)

	app.fireEvent(EvtInit, nil)
	app.initialized = true
}

// SetLogo text and color style
func (app *App) SetLogo(logo string, style ...string) {
	app.Logo.Text = logo
	if len(style) > 0 {
		app.Logo.Style = style[0]
	}
}

// SetDebugMode level
func (app *App) SetDebugMode() {
	SetDebugMode()
}

// SetQuietMode level
func (app *App) SetQuietMode() {
	SetQuietMode()
}

// SetVerbose level
func (app *App) SetVerbose(verbose uint) {
	SetVerbose(verbose)
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
	c.Name = strings.Trim(strings.TrimSpace(c.Name), ": ")
	if c.Name == "" {
		panicf("the added command name can not be empty.")
	}

	if c.IsDisabled() {
		Logf(VerbDebug, "command %s has been disabled, skip add", c.Name)
		return c
	}

	// initialize application
	if !app.initialized {
		app.initialize()
	}

	// check and find module name
	if i := strings.IndexByte(c.Name, ':'); i > 0 {
		c.module = c.Name[:i]
	}

	nameLen := len(c.Name)

	// add command to app
	app.names[c.Name] = nameLen
	app.commands[c.Name] = c

	// record command name max length
	if nameLen > app.nameMaxLen {
		app.nameMaxLen = nameLen
	}

	if _, ok := app.moduleCommands[c.module]; !ok {
		app.moduleCommands[c.module] = make(map[string]*Command)
	}
	app.moduleCommands[c.module][c.Name] = c

	// add aliases for the command
	app.AddAliases(c.Name, c.Aliases)
	Logf(VerbDebug, "[App.AddCommand] add a new CLI command: %s", c.Name)

	// init command
	c.app = app
	c.initialize()
	return c
}

// AddAliases add alias names for a command
func (app *App) AddAliases(command string, names []string) {
	if app.aliases == nil {
		app.aliases = make(map[string]string)
	}

	// add alias
	for _, alias := range names {
		if cmd, has := app.aliases[alias]; has {
			panicf("The alias '%s' has been used by command '%s'", alias, cmd)
		}

		app.aliases[alias] = command
	}
}

// On add hook handler for a hook event
func (app *App) On(name string, handler HookFunc) {
	Logf(VerbDebug, "[App.On] add application hook: %s", name)

	app.SimpleHooks.On(name, handler)
}

func (app *App) fireEvent(event string, data interface{}) {
	Logf(VerbDebug, "[App.Fire] trigger the application event: %s", event)

	app.SimpleHooks.Fire(event, app, data)
}

// stop application and exit
// func stop(code int) {
// 	os.Exit(code)
// }

// AddError to the application
func (app *App) AddError(err error) {
	app.errors = append(app.errors, err)
}

// Names get all command names
func (app *App) Names() map[string]int {
	return app.names
}

// Commands get all commands
func (app *App) Commands() map[string]*Command {
	return app.commands
}

// CleanArgs get clean args
func (app *App) CleanArgs() []string {
	return app.cleanArgs
}
