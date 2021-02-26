package gcli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/helper"
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
	commandBase

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
	Aliases []string
	// Category for the command
	Category string
	// Config func, will call on `initialize`.
	// - you can config options and other init works
	Config func(c *Command)

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
	Subs []*Command

	// module is the name for grouped commands
	// subName is the name for grouped commands
	// eg: "sys:info" -> module: "sys", subName: "info"
	module, subName string
	// Examples some usage example display
	Examples string
	// Func is the command handler func. Func Runner
	Func RunnerFunc
	// Help is the long help message text
	Help string
	// HelpRender custom render cmd help message
	HelpRender func(c *Command)

	// CustomFlags indicates that the command will do its own flag parsing.
	// CustomFlags bool

	// application
	app *App
	// mark is alone running.
	alone bool
	// mark is disabled. if true will skip register to cli-app.
	disabled bool
}

// NewCommand create a new command instance.
// Usage:
// 	cmd := NewCommand("my-cmd", "description")
//	// OR with an config func
// 	cmd := NewCommand("my-cmd", "description", func(c *Command) { ... })
// 	app.Add(cmd) // OR cmd.AttachTo(app)
func NewCommand(name, desc string, fn ...func(c *Command)) *Command {
	c := &Command{
		Name: name,
		Desc: desc,
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
func (c *Command) SetFunc(fn RunnerFunc) {
	c.Func = fn
}

// WithFunc Settings command handler func
func (c *Command) WithFunc(fn RunnerFunc) *Command {
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
	sub.parent = c
	// inherit global flags from application
	sub.core.gFlags = c.gFlags

	// initialize command
	c.initialize()

	// extend path names from parent
	sub.pathNames = c.pathNames[0:]

	// do add
	c.commandBase.addCommand(sub)
}

// Match sub command by input names
func (c *Command) Match(names []string) *Command {
	// must ensure is initialized
	c.initialize()

	ln := len(names)
	if ln == 0 { // return self.
		return c
	}

	return c.commandBase.Match(names)
}

// Match command by path. eg. "top:sub"
// func (c *Command) MatchByPath(path string) *Command {
// 	var names []string
// 	path = strings.TrimSpace(path)
// 	if path != "" {
// 		names = strings.Split(path, CommandSep)
// 	}
//
// 	return c.Match(names)
// }

// init core
func (c *Command) initCore(cmdName string) {
	c.core.cmdLine = CLI

	c.AddVars(c.innerHelpVars())
	c.AddVars(map[string]string{
		"cmd": cmdName,
		// binName with command
		"binWithCmd": c.binName + " " + cmdName,
		// binFile with command
		"fullCmd": c.binFile + " " + cmdName,
	})
}

// initialize works for the command
func (c *Command) initialize() {
	if c.initialized {
		return
	}

	// check command name
	cName := c.goodName()

	// init core
	// c.core.init(cName)
	// c.core = newCore(cName)
	c.initCore(cName)
	c.pathNames = append(c.pathNames, cName)

	// init commandBase
	c.commandBase = newCommandBase()
	c.initialized = true

	// load common subs
	if len(c.Subs) > 0 {
		for _, sub := range c.Subs {
			c.AddCommand(sub)
		}
	}

	// init for cmd Arguments
	c.Arguments.SetName(cName)
	c.Arguments.SetValidateNum(!c.alone && gOpts.strictMode)

	// init for cmd Flags
	c.Flags.InitFlagSet(cName)
	c.Flags.FSet().SetOutput(c.Flags.out)
	c.Flags.FSet().Usage = func() { // call on exists "-h" "--help"
		Logf(VerbDebug, "render help message on exists '-h|--help' or has unknown flag")
		c.ShowHelp()
	}

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

	c.Fire(EvtCmdInit, nil)
}

// IsAlone running
func (c *Command) IsAlone() bool {
	return c.alone
}

// NotAlone running
func (c *Command) NotAlone() bool {
	return !c.alone
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
 * alone running
 *************************************************************/

var errCallRun = errors.New("c.Run() method can only be called in standalone mode")

// MustRun Alone the current command, will panic on error
func (c *Command) MustRun(args []string) {
	if err := c.Run(args); err != nil {
		color.Error.Println("Run command error: %s", err.Error())
		panic(err)
	}
}

// Run Alone running current command
func (c *Command) Run(args []string) (err error) {
	// - Running in application.
	if c.app != nil {
		return errCallRun
	}

	// - Alone running command

	// mark is alone
	c.alone = true
	// only init global flags on alone run.
	c.core.gFlags = NewFlags(c.Name + ".GlobalOpts").WithOption(FlagsOption{
		Alignment: AlignLeft,
	})

	// binding global options
	bindingCommonGOpts(c.gFlags)

	// init the command
	c.initialize()

	// add default error handler.
	c.AddOn(EvtCmdError, defaultErrHandler)

	// if not set input args
	if args == nil {
		args = os.Args[1:]
	}

	// parse global options
	gf := c.GlobalFlags()
	err = gf.Parse(args)
	if err != nil {
		color.Error.Tips(err.Error())
		return
	}

	// remaining args
	// args = gf.fSet.Args()
	args = gf.RawArgs()

	// contains keywords "-h" OR "--help" on end
	// if c.hasHelpKeywords() {
	// 	c.ShowHelp()
	// 	return
	// }

	// dispatch and parse flags and execute command
	return c.innerDispatch(args)
	// parse flags and execute command
	// return c.innerExecute(args, true)
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
			err = nil
			// TODO call show help on there
			// c.ShowHelp()
			return
		}

		color.Error.Tips("Options parse error - %s", err.Error())
		return
	}

	Debugf("cmd: %s - remain args on options parsed: %v", c.Name, args)

	// find sub command
	if len(args) > 0 {
		name := args[0]

		// ensure is not an option
		if name[0] != '-' {
			name = c.ResolveAlias(name)

			// name is an sub command name?
			if c.IsCommand(name) {
				sub := c.Command(args[0])

				// loop find sub...command and run it.
				return sub.innerDispatch(args[1:])
			}
		}
	}

	// do execute command
	return c.doExecute(args)
}

// execute the command
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
	Logf(VerbDebug, "option flags on after format: %v", args)

	// NOTICE: disable output internal error message on parse flags
	// c.FSet().SetOutput(ioutil.Discard)

	// parse options, don't contains command name.
	if err = c.Parse(args); err != nil {
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
	c.Fire(EvtCmdBefore, args)

	// collect and binding named argument
	if err := c.ParseArgs(args); err != nil {
		c.Fire(EvtCmdError, err)
		return err
	}

	// do call command handler func
	if c.Func == nil {
		Logf(VerbWarn, "the command '%s' no handler func to running", c.Name)
	} else {
		// err := c.Func.Run(c, args)
		err = c.Func(c, args)
	}

	if err != nil {
		c.Fire(EvtCmdError, err)
	} else {
		c.Fire(EvtCmdAfter, nil)
	}
	return
}

/*************************************************************
 * command help
 *************************************************************/

// CmdHelpTemplate help template for a command
var CmdHelpTemplate = `{{.Desc}}
{{if .Cmd.NotAlone}}
<comment>Name:</> {{.Cmd.Name}}{{if .Cmd.Aliases}} (alias: <info>{{.Cmd.AliasesString}}</>){{end}}{{end}}
<comment>Usage:</> {$binName} [Global Options...] {{if .Cmd.NotAlone}}<info>{{.Cmd.Name}}</> {{end}}[--option ...] [arguments ...]

<comment>Global Options:</>
{{.GOpts}}{{if .Options}}
<comment>Options:</>
{{.Options}}{{end}}{{if .Cmd.Args}}
<comment>Arguments:</>{{range $a := .Cmd.Args}}
  <info>{{$a.HelpName | printf "%-12s"}}</>{{$a.Desc | ucFirst}}{{if $a.Required}}<red>*</>{{end}}{{end}}
{{end}}{{if .Cmd.Examples}}
<comment>Examples:</>
{{.Cmd.Examples}}{{end}}{{if .Cmd.Help}}
<comment>Help:</>
{{.Cmd.Help}}{{end}}`

// ShowHelp show command help info
func (c *Command) ShowHelp() {
	// custom help render func
	if c.HelpRender != nil {
		c.HelpRender(c)
		return
	}

	// clear space and empty new line
	if c.Examples != "" {
		c.Examples = strings.Trim(c.Examples, "\n") + "\n"
	}

	// clear space and empty new line
	if c.Help != "" {
		c.Help = strings.Join([]string{strings.TrimSpace(c.Help), "\n"}, "")
	}

	// render help message
	s := helper.RenderText(CmdHelpTemplate, map[string]interface{}{
		"Cmd": c,
		// global options
		"GOpts": c.gFlags.String(),
		// parse options to string
		"Options": c.Flags.String(),
		// always upper first char
		"Desc": c.Desc,
	}, nil)

	// parse help vars then print help
	color.Print(c.ReplaceVars(s))
	// fmt.Printf("%#v\n", s)
}

/*************************************************************
 * helper methods
 *************************************************************/

// Fire event handler by name
func (c *Command) Fire(event string, data interface{}) {
	Debugf("command '%s' trigger the event: <mga>%s</>", c.Name, event)

	c.Hooks.Fire(event, c, data)
}

// On add hook handler for a hook event
func (c *Command) On(name string, handler HookFunc) {
	Debugf("command '%s' add hook: %s", c.Name, name)

	c.Hooks.On(name, handler)
}

// Copy a new command for current
func (c *Command) Copy() *Command {
	nc := *c
	// reset some fields
	nc.Func = nil
	nc.Hooks.ClearHooks()
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
// func (c *Command) Logf(level uint, format string, v ...interface{}) {
// 	Logf(level, format, v...)
// }
