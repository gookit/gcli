package cliapp

import (
	"flag"
	"fmt"
	"github.com/gookit/cliapp/utils"
	"github.com/gookit/color"
	"strings"
	"strconv"
)

// CmdRunner interface
type CmdRunner interface {
	Run(cmd *Command, args []string) int
}

// CmdFunc definition
type CmdFunc func(cmd *Command, args []string) int

// Run implement the CmdRunner interface
func (f CmdFunc) Run(cmd *Command, args []string) int {
	return f(cmd, args)
}

// HookFunc definition
type HookFunc func(cmd *Command, data interface{})

// Command a CLI command structure
type Command struct {
	// Name is the command name.
	Name string
	// Func A callback func to runs the command.
	// Func func(cmd *Command, args []string) int
	// Func CmdRunner
	Func CmdFunc
	// Hooks can setting some hooks func on running.
	// allow hooks: "init", "before", "after", "error"
	Hooks map[string]HookFunc
	// Aliases is the command name's alias names
	Aliases []string
	// Description is the command description for 'go help'
	Description string
	// Flags(Options) is a set of flags specific to this command.
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
	// Used to record argument names and defined positional relationships
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

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as import path.
func (c *Command) Runnable() bool {
	return c.Func != nil
}

// Init command
func (c *Command) Init() *Command {
	if len(c.Description) > 0 {
		c.Description = utils.UpperFirst(c.Description)

		// if contains help var "{$cmd}"
		if strings.Contains(c.Description, "{$cmd}") {
			c.Description = strings.Replace(c.Description, "{$cmd}", c.Name, -1)
		}
	}

	c.AddVars(map[string]string{
		"cmd": c.Name,
		// full command
		"fullCmd": binName + " " + c.Name,
		"workDir": workDir,
		"binName": binName,
	})

	if c.Hooks == nil {
		c.Hooks = make(map[string]HookFunc, 1)
	}

	c.callHook(EvtInit, nil)

	// add default error handler
	if _, ok := c.Hooks[EvtError]; !ok {
		c.Hooks[EvtError] = c.defaultErrHandler
	}

	// set help handler
	c.Flags.Usage = func() {
		c.ShowHelp(true)
	}

	return c
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

// Execute do execute the command
func (c *Command) Execute(args []string) int {
	c.callHook(EvtBefore, nil)

	// call command handler func
	// eCode := c.Func.Run(c, args)
	eCode := c.Func(c, args)

	if c.error != nil {
		c.app.AddError(c.error)
		c.callHook(EvtError, c.error)
	} else {
		c.callHook(EvtAfter, eCode)
	}

	return eCode
}

func (c *Command) callHook(event string, data interface{}) {
	Logf(VerbDebug, "command '%s' trigger the hook: %s", c.Name, event)

	if handler, ok := c.Hooks[event]; ok {
		handler(c, data)
	}
}

/*************************************************************
 * command arguments
 *************************************************************/

// Argument a command argument definition
type Argument struct {
	// Name argument name
	Name string
	// Description argument description message
	Description string
	// IsArray if is array, can allow accept multi values, and must in last.
	IsArray bool
	// Required arg is required
	Required bool
	// value store parsed argument data. (string, []string)
	Value interface{}
}

// Int argument value to int
func (a *Argument) Int(defVal ...int) int {
	def := 0
	if len(defVal) == 1 {
		def = defVal[0]
	}

	if a.Value == nil {
		return def
	}

	str := a.Value.(string)
	if str != "" {
		val, err := strconv.Atoi(str)
		if err != nil {
			return val
		}
	}

	return def
}

// String argument value to string
func (a *Argument) String(defVal ...string) string {
	def := ""
	if len(defVal) == 1 {
		def = defVal[0]
	}

	if a.Value == nil {
		return def
	}

	return a.Value.(string)
}

// Strings argument value to string array, if argument isArray = true.
func (a *Argument) Strings() (ss []string) {
	if a.Value != nil {
		ss = a.Value.([]string)
	}

	return
}

// AddArg add a command argument.
// Notice:
// 	- Required argument cannot be defined after optional argument
// 	- The (array) argument of multiple values ​​can only be defined at the end
func (c *Command) AddArg(name, description string, required, isArray bool) {
	if c.hasArrayArg {
		panic("An array argument has been defined and no more argument definitions can be added")
	}

	if required && c.hasOptionalArg {
		panic("Required argument cannot be defined after optional argument")
	}

	if c.argsIndexes == nil {
		c.argsIndexes = make(map[string]int)
	}

	// add argument index record
	c.argsIndexes[name] = len(c.args)

	// add argument
	c.args = append(c.args, &Argument{
		Name: name, Description: description, Required: required, IsArray: isArray,
	})

	if !required {
		c.hasOptionalArg = true
	}

	if isArray {
		c.hasArrayArg = true
	}
}

// Args get all defined argument
func (c *Command) Args() []*Argument {
	return c.args
}

// Arg get arg by defined name.
// usage:
// 	intVal := c.Arg("name").Int()
// 	strVal := c.Arg("name").String()
// 	arrVal := c.Arg("name").Array()
func (c *Command) Arg(name string) *Argument {
	i, ok := c.argsIndexes[name]
	if !ok {
		return &Argument{}
	}

	return c.args[i]
}

// GetArgs get Flags args
func (c *Command) GetArgs() []string {
	return c.Flags.Args()
}

// GetArg get Flags arg
func (c *Command) GetArg(i int) string {
	return c.Flags.Arg(i)
}

/*************************************************************
 * helper methods
 *************************************************************/

// Disable set cmd is disabled
func (c *Command) Disable() {
	c.disabled = true
}

// IsDisabled get cmd is disabled
func (c *Command) IsDisabled() bool {
	return c.disabled
}

// Errorf format message and add error for the command
func (c *Command) Errorf(format string, v ...interface{}) int {
	return c.SetError(fmt.Errorf(format, v...))
}

// SetError for the command
func (c *Command) SetError(err error) int {
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

// WorkDir returns command work dir
func (c *Command) WorkDir() string {
	return workDir
}
