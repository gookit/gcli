package cliapp

import (
	"flag"
	"fmt"
	"github.com/gookit/cliapp/utils"
	"github.com/gookit/color"
	"log"
	"os"
	"strings"
)

// Runner interface
type Runner interface {
	Run(cmd *Command, args []string) int
}

// CmdFunc definition
type CmdFunc func(c *Command, args []string) int

// Run implement the Runner interface
func (f CmdFunc) Run(c *Command, args []string) int {
	return f(c, args)
}

// HookFunc definition
type HookFunc func(c *Command, data interface{})

// Command a CLI command structure
type Command struct {
	// is internal use
	*CmdLine
	// Name is the command name.
	Name string
	// Func is the command handler func.
	// Func Runner
	Func CmdFunc
	// Config func, will call on `initialize`. you can config options and other works
	Config func(c *Command)
	// Hooks can setting some hooks func on running.
	// allow hooks: "init", "before", "after", "error"
	Hooks map[string]HookFunc
	// Aliases is the command name's alias names
	Aliases []string
	// Description is the command description for 'go help'
	Description string
	// Flags(command options) is a set of flags specific to this command.
	Flags flag.FlagSet
	// CustomFlags indicates that the command will do its own flag parsing.
	CustomFlags bool
	// Vars you can add some vars map for render help info
	Vars map[string]string
	// Help is the long help message text
	Help string
	// Examples some usage example display
	Examples string

	// Args definition for the command.
	// eg. {
	// 	{"arg0", "this is first argument", false, false},
	// 	{"arg1", "this is second argument", false, false},
	// }
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
	app *Application
	// mark is alone running.
	alone bool
	// store a command error
	error error
	// mark is disabled. if true will skip register to cli-app.
	disabled bool
	// option names {name:short}
	optNames map[string]string
	// shortcuts for command options(Flags) {short:name} eg. {"n": "name", "o": "opt"}
	shortcuts map[string]string
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
	c.CmdLine = CLI

	// format description
	if len(c.Description) > 0 {
		c.Description = utils.UcFirst(c.Description)

		// contains help var "{$cmd}". replace on here is for 'app help'
		if strings.Contains(c.Description, "{$cmd}") {
			c.Description = strings.Replace(c.Description, "{$cmd}", c.Name, -1)
		}
	}

	// call Config func
	if c.Config != nil {
		c.Config(c)
	}

	// set help vars
	// c.Vars = c.app.vars // Error: var is map, map is ref addr
	c.AddVars(c.app.vars)
	c.AddVars(map[string]string{
		"cmd": c.Name,
		// full command
		"fullCmd": CLI.binName + " " + c.Name,
	})

	if c.Hooks == nil {
		c.Hooks = make(map[string]HookFunc, 1)
	}

	c.callHook(EvtInit, nil)

	// add default error handler
	if _, ok := c.Hooks[EvtError]; !ok {
		c.Hooks[EvtError] = c.defaultErrHandler
	}

	// init for Flags
	c.Flags.Init(c.Name, flag.ContinueOnError)
	c.Flags.Usage = func() { // call on exists "-h" "--help"
		c.ShowHelp(true)
	}

	// bind some global options
	// c.Flags.BoolVar(&gOpts.showHelp, "h", false, "")
	// c.Flags.BoolVar(&gOpts.showHelp, "help", false, "")
	return c
}

/*************************************************************
 * command run
 *************************************************************/

// Execute do execute the command
func (c *Command) Execute(args []string) int {
	// collect named args
	if err := c.collectNamedArgs(args); err != nil {
		fmt.Println(color.FgRed.Render("ERROR:"), err.Error())
		return ERR
	}

	var eCode int
	c.callHook(EvtBefore, args)

	// call command handler func
	if c.Func == nil {
		Logf(VerbWarn, "the command '%s' no handler func to running.", c.Name)
	} else {
		// eCode := c.Func.Run(c, args)
		eCode = c.Func(c, args)
	}

	if c.error != nil {
		c.app.AddError(c.error)
		c.callHook(EvtError, c.error)
	} else {
		c.callHook(EvtAfter, eCode)
	}

	return eCode
}

func (c *Command) collectNamedArgs(inArgs []string) error {
	var num int
	inNum := len(inArgs)

	for i, arg := range c.args {
		num = i + 1      // num is equal index + 1
		if num > inNum { // no enough arg
			if arg.Required {
				return fmt.Errorf("must set a value for the argument: %s (position %d)", arg.Name, arg.index)
			}
			break
		}

		if arg.IsArray {
			arg.Value = inArgs[i:]
			inNum = num // must reset inNum
		} else {
			arg.Value = inArgs[i]
		}
	}

	if c.app.Strict && inNum > num {
		return fmt.Errorf("enter too many arguments: %v", inArgs[num:])
	}

	return nil
}

func (c *Command) callHook(event string, data interface{}) {
	Logf(VerbDebug, "command '%s' trigger the hook: %s", c.Name, event)

	if handler, ok := c.Hooks[event]; ok {
		handler(c, data)
	}
}

func (c *Command) defaultErrHandler(_ *Command, data interface{}) {
	err := data.(error)
	fmt.Println(color.FgRed.Render("ERROR:"), err.Error())
}

// Copy a new command for current
func (c *Command) Copy() *Command {
	nc := *c
	// reset some fields
	nc.Func = nil
	nc.Hooks = nil
	// nc.Flags = flag.FlagSet{}

	return &nc
}

// On add hook handler for a hook event
func (c *Command) On(name string, handler func(c *Command, data interface{})) {
	c.Hooks[name] = handler
}

/*************************************************************
 * alone running
 *************************************************************/

// AloneRun current command
func (c *Command) AloneRun() int {
	// don't display date on print log
	log.SetFlags(0)
	// mark is alone
	c.alone = true
	// args := parseGlobalOpts()
	// init
	c.initialize()
	// parse args and opts
	c.Flags.Parse(os.Args[1:])

	return c.Execute(c.Flags.Args())
}

// IsAlone running
func (c *Command) IsAlone() bool {
	return c.alone
}

// NotAlone running
func (c *Command) NotAlone() bool {
	return !c.alone
}

/*************************************************************
 * helper methods
 *************************************************************/

// Errorf format message and add error to the command
func (c *Command) Errorf(format string, v ...interface{}) int {
	return c.WithError(fmt.Errorf(format, v...))
}

// WithError for the command
func (c *Command) WithError(err error) int {
	c.error = err
	return ERR
}

// Error get error of the command
func (c *Command) Error() error {
	return c.error
}

// App returns the CLI application
func (c *Command) App() *Application {
	return app
}

// AddVars add multi tpl vars
func (c *Command) AddVars(vars map[string]string) {
	if c.Vars == nil {
		c.Vars = make(map[string]string)
	}

	for n, v := range vars {
		c.Vars[n] = v
	}
}

// GetVar get a help var by name
func (c *Command) GetVar(name string) string {
	if v, ok := c.Vars[name]; ok {
		return v
	}

	return ""
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
