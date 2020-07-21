package gcli

import (
	"flag"
	"fmt"
	"strings"

	"github.com/gookit/goutil/strutil"
)

// Runner interface
type Runner interface {
	// Config(c *Command)
	Run(c *Command, args []string) error
}

// CmdFunc definition
type CmdFunc func(c *Command, args []string) error

// Run implement the Runner interface
func (f CmdFunc) Run(c *Command, args []string) error {
	return f(c, args)
}

// Command a CLI command structure
type Command struct {
	// cmdLine is internal use
	*cmdLine
	HelpVars
	// Hooks can allow setting some hooks func on running.
	Hooks // allowed hooks: "init", "before", "after", "error"

	// Name is the command name.
	Name string
	// module is the name for grouped commands
	module string
	// UseFor is the command description message.
	UseFor string
	// Aliases is the command name's alias names
	Aliases []string
	// Config func, will call on `initialize`. you can config options and other works
	Config func(c *Command)
	// CustomFlags indicates that the command will do its own flag parsing.
	CustomFlags bool
	// Flags(command options) is a set of flags specific to this command.
	Flags flag.FlagSet
	// Examples some usage example display
	Examples string
	// Func is the command handler func. Func Runner
	Func CmdFunc
	// Help is the long help message text
	Help string

	// arguments definition for the command.
	// eg. {
	// 	{"arg0", "this is first argument", false, false},
	// 	{"arg1", "this is second argument", false, false},
	// }
	// if you want get raw args, can use: c.RawArgs()
	args []*Argument
	// record min length for args
	// argsMinLen int
	// record argument names and defined positional relationships
	// {
	// 	// name: position
	// 	"arg0": 0,
	// 	"arg1": 1,
	// }
	argsIndexes    map[string]int
	hasArrayArg    bool
	hasOptionalArg bool

	// application
	app *App
	// mark is alone running.
	alone bool
	// mark is disabled. if true will skip register to cli-app.
	disabled bool
	// all option names of the command
	optNames map[string]string
	// shortcuts for command options(Flags) {short:name} eg. {"n": "name", "o": "opt"}
	shortcuts map[string]string
}

// NewCommand create a new command instance.
// Usage:
// 	cmd := NewCommand("my-cmd", "description", func(c *Command) { ... })
// 	app.Add(cmd) // OR cmd.AttachTo(app)
func NewCommand(name, useFor string, config func(c *Command)) *Command {
	return &Command{
		Name:   name,
		UseFor: useFor,
		Config: config,
	}
}

// SetFunc Settings command handler func
func (c *Command) SetFunc(fn CmdFunc) *Command {
	c.Func = fn
	return c
}

// AttachTo attach the command to CLI application
func (c *Command) AttachTo(app *App) {
	app.AddCommand(c)
}

// Disable set cmd is disabled
func (c *Command) Disable() {
	c.disabled = true
}

// IsDisabled get cmd is disabled
func (c *Command) IsDisabled() bool {
	return c.disabled
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as import path.
func (c *Command) Runnable() bool {
	return c.Func != nil
}

// initialize command
func (c *Command) initialize() *Command {
	c.cmdLine = CLI

	// format description
	if len(c.UseFor) > 0 {
		c.UseFor = strutil.UpperFirst(c.UseFor)

		// contains help var "{$cmd}". replace on here is for 'app help'
		if strings.Contains(c.UseFor, "{$cmd}") {
			c.UseFor = strings.Replace(c.UseFor, "{$cmd}", c.Name, -1)
		}
	}

	// call Config func
	if c.Config != nil {
		c.Config(c)
	}

	// set help vars
	// c.Vars = c.app.vars // Error: var is map, map is ref addr
	c.AddVars(c.helpVars())
	c.AddVars(map[string]string{
		"cmd": c.Name,
		// full command
		"fullCmd": c.binName + " " + c.Name,
	})

	c.Fire(EvtInit, nil)

	// if not set application instance
	if c.app == nil {
		// mark is alone
		c.alone = true
		// add default error handler.
		c.Hooks.AddOn(EvtError, defaultErrHandler)
	}

	// init for Flags
	c.Flags.Init(c.Name, flag.ContinueOnError)
	c.Flags.Usage = func() { // call on exists "-h" "--help"
		c.ShowHelp(false)
	}

	return c
}

// IsAlone running
func (c *Command) IsAlone() bool {
	return c.alone
}

// NotAlone running
func (c *Command) NotAlone() bool {
	return !c.alone
}

// Module name of the grouped command
func (c *Command) Module() string {
	return c.module
}

// ID get command ID name.
func (c *Command) ID() string {
	if c.module != "" {
		return fmt.Sprintf("%s:%s", c.module, c.Name)
	}

	return c.Name
}

/*************************************************************
 * helper methods
 *************************************************************/

// App returns the CLI application
func (c *Command) App() *App {
	return c.app
}

// Errorf format message and add error to the command
func (c *Command) Errorf(format string, v ...interface{}) error {
	return fmt.Errorf(format, v...)
}

// AliasesString returns aliases string
func (c *Command) AliasesString(sep ...string) string {
	s := ","
	if len(sep) == 1 {
		s = sep[0]
	}

	return strings.Join(c.Aliases, s)
}

// Logf print log message
func (c *Command) Logf(level uint, format string, v ...interface{}) {
	Logf(level, format, v...)
}
