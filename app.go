package gcli

import (
	"os"
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
	core
	// *cmdLine
	// HelpVars
	// Hooks // allow hooks: "init", "before", "after", "error"
	commandBase

	// Name app name
	Name string
	// Desc app description
	Desc string
	// Func on run app, if is empty will display help.
	Func func(app *App, args []string) error
	// ExitOnEnd call os.Exit on running end
	ExitOnEnd bool
	// ExitFunc default is os.Exit
	ExitFunc func(int)

	// args on after parse global options and command name.
	args []string
	// all commands by module TODO remove
	moduleCommands map[string]map[string]*Command

	// rawFlagArgs []string
	// clean os.args, not contains bin-name and command-name
	cleanArgs []string
	// current command name
	commandName string
}

// NewApp create new app instance.
// Usage:
// 	NewApp()
// 	// Or with a config func
// 	NewApp(func(a *App) {
// 		// do something before init ....
// 		a.Hooks[gcli.EvtInit] = func () {}
// 	})
func NewApp(fn ...func(app *App)) *App {
	app := &App{
		Name: "GCli App",
		Desc: "This is my console application",
		// set a default version
		// Version: "1.0.0",
		// config
		ExitOnEnd: true,
		// group
		moduleCommands: make(map[string]map[string]*Command),
	}

	// internal core
	app.core = core{
		cmdLine: CLI,
		gFlags: NewFlags("app.GlobalOpts").WithOption(FlagsOption{
			WithoutType: true,
			NameDescOL:  true,
			Alignment:   AlignLeft,
			TagName:     FlagTagName,
		}),
	}
	// init commandBase
	app.commandBase = newCommandBase()
	// set a default version
	app.Version = "1.0.0"
	app.SetLogo("", "info")

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

// Exit get the app GlobalFlags
func (app *App) Exit(code int) {
	if app.ExitFunc == nil {
		os.Exit(code)
	}

	app.ExitFunc(code)
}

// binding global options
func (app *App) bindingGlobalOpts() {
	Logf(VerbDebug, "will begin binding global options")
	// global options flag
	// gf := flag.NewFlagSet(app.Args[0], flag.ContinueOnError)
	gf := app.GlobalFlags()

	// binding global options
	bindingCommonGOpts(gf)
	// add more ...
	gf.BoolOpt(&gOpts.showVer, "version", "V", false, "Display app version information")
	// This is a internal command
	gf.BoolVar(&gOpts.inCompletion, FlagMeta{
		Name: "cmd-completion",
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

	// app.names = make(map[string]int)

	// init some help tpl vars
	app.core.AddVars(app.core.innerHelpVars())

	// binding GlobalOpts
	app.bindingGlobalOpts()
	// parseGlobalOpts()

	// add default error handler.
	app.core.AddOn(EvtAppError, defaultErrHandler)

	app.fireEvent(EvtAppInit, nil)
	app.initialized = true
}

// NewCommand create a new command
// func (app *App) NewCommand(name, useFor string, config func(c *Command)) *Command {
// 	return NewCommand(name, useFor, config)
// }

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

// AddCommand add a new command to the app
func (app *App) AddCommand(c *Command) {
	// initialize application before add command
	app.initialize()

	// init command
	c.app = app
	// inherit global flags from application
	c.core.gFlags = app.gFlags

	// do add
	app.commandBase.addCommand(c)
}

// AddCommander to the application
func (app *App) AddCommander(cmder Commander) {
	c := cmder.Creator()
	c.Func = cmder.Execute

	// binding flags
	cmder.Config(c)
	app.AddCommand(c)
}

// AddAliases add alias names for a command
func (app *App) AddAliases(command string, aliases ...string) {
	app.addAliases(command, aliases, true)
}

// addAliases add alias names for a command
func (app *App) addAliases(command string, aliases []string, sync bool) {
	c, has := app.commands[command]
	if !has {
		panicf("The command '%s' is not exists", command)
	}

	// add alias
	for _, alias := range aliases {
		if app.IsCommand(alias) {
			panicf("The name '%s' has been used as an command name", alias)
		}

		app.cmdAliases.AddAlias(command, alias)

		// sync to Command
		if sync {
			c.Aliases = append(c.Aliases, alias)
		}
	}
}

// Match command by path. eg. ["top", "sub"]
// func (app *App) Match(names []string) *Command {
// 	return app.commandBase.match(names)
// }

// On add hook handler for a hook event
// func (app *App) BeforeInit(name string, handler HookFunc) {}

// On add hook handler for a hook event
func (app *App) On(name string, handler HookFunc) {
	Logf(VerbDebug, "add application hook: %s", name)

	app.core.On(name, handler)
}

func (app *App) fireEvent(event string, data interface{}) {
	Logf(VerbDebug, "trigger the application event: <mga>%s</>", event)

	app.core.Fire(event, app, data)
}

// stop application and exit
// func stop(code int) {
// 	os.Exit(code)
// }

// CleanArgs get clean args
func (app *App) CleanArgs() []string {
	return app.cleanArgs
}
