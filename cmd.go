package gcli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/events"
	"github.com/gookit/gcli/v3/helper"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/structs"
	"github.com/gookit/goutil/strutil"
)

// Runner /Executor interface
type Runner interface {
	// Run Config(c *Command)
	Run(c *Command, args []string) error
}

// RunnerFunc definition
type RunnerFunc func(c *Command, args []string) error

// Run implement the Runner interface
func (f RunnerFunc) Run(c *Command, args []string) error {
	return f(c, args)
}

const maxFunc = 64

// HandlersChain middleware handlers chain definition
type HandlersChain []RunnerFunc

// Last returns the last handler in the chain. ie. the last handler is the main own.
func (c HandlersChain) Last() RunnerFunc {
	length := len(c)
	if length > 0 {
		return c[length-1]
	}
	return nil
}

// Command a CLI command structure
type Command struct {
	// internal use
	base

	// --- provide option and argument parse and binding.

	// Flags options for the command
	Flags
	// Arguments for the command
	Arguments

	// Name is the full command name.
	Name string
	// Desc is the command description message.
	Desc string

	// Aliases is the command name's alias names
	Aliases arrutil.Strings
	// Category for the command
	Category string
	// Config func, will call on `initialize`.
	//
	// - you can config options and other init works
	Config func(c *Command)
	// Hidden the command on render help
	Hidden bool

	// --- for middleware ---
	// run error
	runErr error
	// middleware index number
	middleIdx int8
	// middleware functions
	middles HandlersChain
	// errorHandler // loop find parent.errorHandler

	// path names of the command. 'parent current'
	pathNames []string

	// Parent parent command
	parent *Command

	// Subs sub commands of the Command
	// NOTICE: if command has been initialized, adding through this field is invalid
	Subs []*Command

	// module is the name for grouped commands
	// subName is the name for grouped commands
	// eg: "sys:info" -> module: "sys", subName: "info"
	// module, subName string
	// Examples some usage example display
	Examples string
	// Func is the command handler func. Func Runner
	Func RunnerFunc
	// Help is the long help message text
	Help string
	// HelpRender custom render cmd help message
	HelpRender func(c *Command)

	// command is inject to the App
	app  *App
	root bool // is root command
	// mark is disabled. if true will skip register to cli-app.
	disabled bool
	// command is standalone running.
	standalone bool
	// global option binding on standalone.
	goptBounded bool
}

// NewCommand create a new command instance.
//
// Usage:
//
//	cmd := NewCommand("my-cmd", "description")
//	// OR with an config func
//	cmd := NewCommand("my-cmd", "description", func(c *Command) { ... })
//	app.Add(cmd) // OR cmd.AttachTo(app)
func NewCommand(name, desc string, fn ...func(c *Command)) *Command {
	c := &Command{
		Name: name,
		Desc: desc,
		// Flags: *NewFlags(name, desc),
	}

	// has config func
	if len(fn) > 0 {
		c.Config = fn[0]
	}

	// set name
	c.Arguments.SetName(name)
	return c
}

// Init command. only use for tests
func (c *Command) Init() {
	c.initialize()
}

// SetFunc Settings command handler func
func (c *Command) SetFunc(fn RunnerFunc) {
	c.Func = fn
}

// WithFunc Settings command handler func
func (c *Command) WithFunc(fn RunnerFunc) *Command {
	c.Func = fn
	return c
}

// WithHidden Settings command is hidden
func (c *Command) WithHidden() *Command {
	c.Hidden = true
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

// Visible return cmd is visible
func (c *Command) Visible() bool {
	return c.Hidden == false
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

// Add one or multi sub-command(s). alias of the AddSubs
func (c *Command) Add(sub *Command, more ...*Command) {
	c.AddSubs(sub, more...)
}

// AddSubs add one or multi sub-command(s)
func (c *Command) AddSubs(sub *Command, more ...*Command) {
	c.AddCommand(sub)

	if len(more) > 0 {
		for _, cmd := range more {
			c.AddCommand(cmd)
		}
	}
}

// AddCommand add a sub command
func (c *Command) AddCommand(sub *Command) {
	// init command
	sub.app = c.app
	sub.parent = c
	// inherit standalone value
	sub.standalone = c.standalone
	// inherit something from parent
	sub.Context = c.Context

	// initialize command
	c.initialize()

	// extend path names from parent
	sub.pathNames = c.pathNames[0:]

	// do add
	c.base.addCommand(c.Name, sub)
}

// Match sub command by input names
func (c *Command) Match(names []string) *Command {
	// ensure is initialized
	c.initialize()

	if len(names) == 0 { // return self.
		return c
	}
	return c.base.Match(names)
}

// MatchByPath command by path. eg: "top:sub"
func (c *Command) MatchByPath(path string) *Command {
	return c.Match(splitPath2names(path))
}

// initialize works for the command
//
// - sub-cmd
func (c *Command) initialize() {
	if c.initialized {
		return
	}

	// check command name
	cName := c.goodName()
	Debugf("initialize the command '%s': init flags, run config func", cName)

	c.initialized = true
	c.pathNames = append(c.pathNames, cName)

	// init base
	c.initCommandBase(cName)

	// load common subs
	if len(c.Subs) > 0 {
		for _, sub := range c.Subs {
			c.AddCommand(sub)
		}
	}

	// init for cmd Arguments
	c.Arguments.SetName(cName)
	// c.Arguments.SetValidateNum(gOpts.strictMode)

	// init for cmd Flags
	c.Flags.InitFlagSet(cName)

	// format description
	if len(c.Desc) > 0 {
		c.Desc = strutil.UpperFirst(c.Desc)
		// contains help var "{$cmd}". replace on here is for 'app help'
		if strings.Contains(c.Desc, "{$cmd}") {
			c.Desc = strings.Replace(c.Desc, "{$cmd}", c.Name, -1)
		}
	}

	// call Config func
	if c.Config != nil {
		c.Config(c)
	}

	c.Fire(events.OnCmdInit, nil)
}

// init core
func (c *Command) initCommandBase(cName string) {
	Logf(VerbCrazy, "init command c.base for the command: %s", cName)

	if c.Hooks == nil {
		c.Hooks = &Hooks{}
	}

	binWithPath := c.binName + " " + c.Path()

	c.initHelpVars()
	c.AddVars(map[string]string{
		"cmd": cName,
		// binName with command name
		"binWithCmd": binWithPath,
		// binName with command path
		"binWithPath": binWithPath,
		// binFile with command
		"fullCmd": c.binFile + " " + cName,
	})

	c.base.cmdNames = make(map[string]int)
	c.base.commands = make(map[string]*Command)
	// set an default value.
	c.base.nameMaxWidth = 12
	// c.base.cmdAliases = make(maputil.Aliases)
	c.base.cmdAliases = structs.NewAliases(aliasNameCheck)
}

// Next TODO processing, run all middleware handlers
func (c *Command) Next() {
	c.middleIdx++
	s := int8(len(c.middles))

	for ; c.middleIdx < s; c.middleIdx++ {
		err := c.middles[c.middleIdx](c, c.RawArgs())
		// will abort on error
		if err != nil {
			c.runErr = err
			return
		}
	}
}

/*************************************************************
 * standalone running
 *************************************************************/

// var errCallRunOnApp = errors.New("c.Run() method can only be called in standalone mode")
// var errCallRunOnSub = errors.New("c.Run() cannot allow call at subcommand")

// MustRun Alone the current command, will panic on error
//
// Usage:
//
//	// run with os.Args
//	cmd.MustRun(nil)
//	cmd.MustRun(os.Args[1:])
//	// custom args
//	cmd.MustRun([]string{"-a", ...})
func (c *Command) MustRun(args []string) {
	if err := c.Run(args); err != nil {
		color.Errorln("ERROR:", err.Error())
	}
}

// Run standalone running the command
//
// Usage:
//
//	// run with os.Args
//	cmd.Run(nil)
//	cmd.Run(os.Args[1:])
//	// custom args
//	cmd.Run([]string{"-a", ...})
func (c *Command) Run(args []string) (err error) {
	if c.app != nil || c.parent != nil {
		return c.innerDispatch(args)
	}

	// mark is standalone
	c.standalone = true

	// if not set input args
	if args == nil {
		args = os.Args[1:]
	}

	// init the command
	c.initialize()

	// add default error handler.
	if !c.HasHook(events.OnCmdRunError) {
		c.On(events.OnCmdRunError, defaultErrHandler)
	}

	// binding global options
	if !c.goptBounded {
		c.goptBounded = true
		Debugf("global options will binding to c.Flags on standalone mode")
		// bindingCommonGOpts(&c.Flags)
		gOpts.bindingFlags(&c.Flags)
	}

	// dispatch and parse flags and execute command
	return c.innerDispatch(args)
}

/*************************************************************
 * command run
 *************************************************************/

// dispatch execute the command
func (c *Command) innerDispatch(args []string) (err error) {
	// parse command flags
	args, err = c.parseOptions(args)
	if err != nil {
		// ignore flag.ErrHelp error
		if err == flag.ErrHelp {
			c.ShowHelp()
			return nil
		}

		Debugf("cmd: %s - command options parse error", c.Name)
		color.Error.Tips("option error - %s", err.Error())
		return nil
	}

	// remaining args
	if c.standalone {
		if gOpts.ShowHelp {
			c.ShowHelp()
			return
		}

		c.Fire(events.OnGlobalOptsParsed, map[string]any{"args": args})
	}

	c.Fire(events.OnCmdOptParsed, map[string]any{"args": args})
	Debugf("cmd: %s - remaining args on options parsed: %v", c.Name, args)

	// find sub command
	if len(args) > 0 {
		name := args[0]

		// ensure is not an option
		if name != "" && name[0] != '-' {
			name = c.ResolveAlias(name)

			// is valid sub command
			if sub, has := c.Command(name); has {
				// loop find sub...command and run it.
				return sub.innerDispatch(args[1:])
			}

			// is not a sub command and has no arguments -> error
			if !c.HasArguments() {
				// fire events
				if stop := c.Fire(events.OnCmdSubNotFound, map[string]any{"name": name}); stop {
					return
				}
				if stop := c.Fire(events.OnCmdNotFound, map[string]any{"name": name}); stop {
					return
				}

				color.Error.Tips("subcommand '%s' - not found on the command", name)
			}
		}
	}

	// not set command func and has sub commands.
	if c.Func == nil && len(c.commands) > 0 {
		Logf(VerbWarn, "cmd: %s - c.Func is empty, but has subcommands, render help", c.Name)
		c.ShowHelp()
		return err
	}

	// do execute current command
	return c.doExecute(args)
}

// execute the current command
func (c *Command) innerExecute(args []string, igrErrHelp bool) (err error) {
	// parse flags
	args, err = c.parseOptions(args)
	if err != nil {
		// whether ignore flag.ErrHelp error
		if igrErrHelp && err == flag.ErrHelp {
			err = nil
		}
		return
	}

	// do execute command
	return c.doExecute(args)
}

// do parse option flags, remaining is cmd args
func (c *Command) parseOptions(args []string) (ss []string, err error) {
	// strict format options
	if gOpts.strictMode && len(args) > 0 {
		args = strictFormatArgs(args)
	}

	// fix and compatible
	// args = moveArgumentsToEnd(args)
	// Debugf("cmd: %s - option flags on after format: %v", c.Name, args)

	Debugf("cmd: %s - will parse options from args: %v", c.Name, args)

	// parse options, don't contains command name.
	if err = c.Parse(args); err != nil {
		Logf(VerbCrazy, "cmd: %s - parse options, err: <red>%s</>", c.Name, err.Error())
		return
	}

	// remaining args
	return c.Flags.RawArgs(), nil
}

// prepare: before execute the command
func (c *Command) prepare(_ []string) (status int, err error) {
	return
}

// do execute the command
func (c *Command) doExecute(args []string) (err error) {
	// collect and binding named argument
	Debugf("cmd: %s - collect and binding named argument", c.Name)
	if err := c.ParseArgs(args); err != nil {
		c.Fire(events.OnCmdRunError, map[string]any{"err": err})
		Logf(VerbCrazy, "binding command '%s' arguments err: <red>%s</>", c.Name, err.Error())
		return err
	}

	c.Fire(events.OnCmdRunBefore, map[string]any{"args": args})

	// do call command handler func
	if c.Func == nil {
		Logf(VerbWarn, "the command '%s' no handler func to running", c.Name)
	} else {
		// err := c.Func.Run(c, args)
		err = c.Func(c, args)
	}

	if err != nil {
		c.Fire(events.OnCmdRunError, map[string]any{"err": err})
	} else {
		c.Fire(events.OnCmdRunAfter, nil)
	}
	return
}

/*************************************************************
 * parent and subs
 *************************************************************/

// Root get root command
func (c *Command) Root() *Command {
	if c.parent != nil {
		return c.parent.Root()
	}

	return c
}

// IsRoot command
func (c *Command) IsRoot() bool {
	return c.parent == nil
}

// Parent get parent
func (c *Command) Parent() *Command {
	return c.parent
}

// SetParent set parent
func (c *Command) SetParent(parent *Command) {
	c.parent = parent
}

// Module name of the grouped command
func (c *Command) ParentName() string {
	if c.parent != nil {
		return c.parent.Name
	}

	return ""
}

// Sub get sub command by name. eg "sub"
func (c *Command) Sub(name string) *Command {
	return c.GetCommand(name)
}

// SubCommand get sub command by name. eg "sub"
func (c *Command) SubCommand(name string) *Command {
	return c.GetCommand(name)
}

// IsSubCommand name check. alias of the HasCommand()
func (c *Command) IsSubCommand(name string) bool {
	return c.IsCommand(name)
}

// find sub command by name
// func (c *Command) findSub(name string) *Command {
// 	if index, ok := c.subName2index[name]; ok {
// 		return c.Subs[index]
// 	}
//
// 	return nil
// }

/*************************************************************
 * helper methods
 *************************************************************/

// IsStandalone running
func (c *Command) IsStandalone() bool {
	return c.standalone
}

// NotStandalone running
func (c *Command) NotStandalone() bool {
	return !c.standalone
}

// ID get command ID name.
func (c *Command) goodName() string {
	name := strings.Trim(strings.TrimSpace(c.Name), ": ")
	if name == "" {
		panicf("the command name can not be empty")
	}

	if !helper.IsGoodCmdName(name) {
		panicf("the command name '%s' is invalid, must match: %s", name, helper.RegGoodCmdName)
	}

	// update name
	c.Name = name
	return name
}

// Fire event handler by name
func (c *Command) Fire(event string, data map[string]any) (stop bool) {
	Debugf("cmd: %s - trigger the event: <mga>%s</>", c.Name, event)

	return c.Hooks.Fire(event, newHookCtx(event, c, data))
}

// On add hook handler for a hook event
func (c *Command) On(name string, handler HookFunc) {
	Debugf("cmd: %s - register hook: <cyan>%s</>", c.Name, name)

	if c.Hooks == nil {
		c.Hooks = &Hooks{}
	}
	c.Hooks.On(name, handler)
}

// Copy a new command for current
func (c *Command) Copy() *Command {
	nc := *c
	// reset some fields
	nc.Func = nil
	nc.Hooks.ResetHooks() // TODO bug, will clear c.Hooks
	// nc.Flags = flag.FlagSet{}

	return &nc
}

// App returns the CLI application
func (c *Command) App() *App {
	return c.app
}

// ID get command ID string
func (c *Command) ID() string {
	return strings.Join(c.pathNames, CommandSep)
}

// Path get command full path
func (c *Command) Path() string {
	return strings.Join(c.pathNames, " ")
}

// PathNames get command path names
func (c *Command) PathNames() []string {
	return c.pathNames
}

// Errorf format message and add error to the command
func (c *Command) Errorf(format string, v ...any) error {
	return fmt.Errorf(format, v...)
}

// NewErr format message and add error to the command
func (c *Command) NewErr(msg string) error { return errors.New(msg) }

// NewErrf format message and add error to the command
func (c *Command) NewErrf(format string, v ...any) error {
	return fmt.Errorf(format, v...)
}

// HelpDesc format desc string for render help
func (c *Command) HelpDesc() (desc string) {
	if len(c.Desc) == 0 {
		return
	}

	// dump.P(desc)
	desc = strutil.UpperFirst(c.Desc)
	// contains help var "{$cmd}". replace on here is for 'app help'
	if strings.Contains(desc, "{$cmd}") {
		desc = strings.Replace(desc, "{$cmd}", color.WrapTag(c.Name, "mga"), -1)
	}

	return wrapColor2string(desc)
}

// Logf print log message
// func (c *Command) Logf(level uint, format string, v ...any) {
// 	Logf(level, format, v...)
// }
