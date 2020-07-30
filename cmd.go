package gcli

import (
	"fmt"
	"strings"

	"github.com/gookit/goutil/strutil"
)

// Command a CLI command structure
type Command struct {
	// core is internal use
	core
	// cmdLine is internal use
	// *cmdLine
	// HelpVars
	// // Hooks can allow setting some hooks func on running.
	// Hooks // allowed hooks: "init", "before", "after", "error"

	// Name is the command name.
	Name string
	// module is the name for grouped commands
	module string
	// UseFor is the command description message.
	UseFor string
	// Aliases is the command name's alias names
	Aliases []string
	// Config func, will call on `initialize`.
	// - you can config options and other init works
	Config func(c *Command)
	// Flags(command options) is a set of flags specific to this command.
	// Flags flag.FlagSet
	// Examples some usage example display
	Examples string
	// Func is the command handler func. Func Runner
	Func CmdFunc
	// Help is the long help message text
	Help string

	// CustomFlags indicates that the command will do its own flag parsing.
	CustomFlags bool
	// Flags options for the command.
	Flags
	// Arguments for the command
	Arguments

	// application
	app *App
	// mark is alone running.
	alone bool
	// mark is disabled. if true will skip register to cli-app.
	disabled bool
	// all option names of the command
	optNames map[string]string
}

// NewCommand create a new command instance.
// Usage:
// 	cmd := NewCommand("my-cmd", "description")
//	// OR with an config func
// 	cmd := NewCommand("my-cmd", "description", func(c *Command) { ... })
// 	app.Add(cmd) // OR cmd.AttachTo(app)
func NewCommand(name, useFor string, fn ...func(c *Command)) *Command {
	c := &Command{
		Name:   name,
		UseFor: useFor,
	}

	// has config func
	if len(fn) > 0 {
		c.Config = fn[0]
	}

	// set name
	c.Arguments.SetName(name)
	return c
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

// initialize works for the command
func (c *Command) initialize() *Command {
	c.core.cmdLine = CLI

	// init for cmd Arguments
	c.Arguments.SetName(c.Name)
	c.Arguments.SetValidateNum(!c.alone && gOpts.strictMode)

	// init for cmd Flags
	c.Flags.InitFlagSet(c.Name)
	c.Flags.FSet().SetOutput(c.Flags.out)
	c.Flags.FSet().Usage = func() { // call on exists "-h" "--help"
		Logf(VerbDebug, "render help message on exists '-h|--help' or has unknown flag")
		c.ShowHelp()
	}

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
	c.AddVars(c.core.innerHelpVars())
	c.AddVars(map[string]string{
		"cmd": c.Name,
		// full command
		"fullCmd": c.binName + " " + c.Name,
	})

	c.Fire(EvtCmdInit, nil)

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

// ID get command ID name.
func (c *Command) goodName() string {
	name := strings.Trim(strings.TrimSpace(c.Name), ": ")
	if name == "" {
		panicf("the command name can not be empty")
	}

	if !goodCmdName.MatchString(name) {
		panicf("the command name '%s' is invalid, must match: %s", name, regGoodCmdName)
	}

	// update name
	c.Name = name
	return name
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
