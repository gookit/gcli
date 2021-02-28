package gcli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/helper"
	"github.com/gookit/goutil/structs"
	"github.com/gookit/goutil/strutil"
)

// Runner /Executor interface
type Runner interface {
	// Config(c *Command)
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
	app *App
	// mark is disabled. if true will skip register to cli-app.
	disabled bool
	// command is standalone running.
	standalone bool
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
	// inherit standalone value
	sub.standalone = c.standalone
	// inherit global flags from application
	sub.core.gFlags = c.gFlags

	// initialize command
	c.initialize()

	// extend path names from parent
	sub.pathNames = c.pathNames[0:]

	// do add
	c.commandBase.addCommand(c.Name, sub)
}

// Match sub command by input names
func (c *Command) Match(names []string) *Command {
	// ensure is initialized
	c.initialize()

	ln := len(names)
	if ln == 0 { // return self.
		return c
	}

	return c.commandBase.Match(names)
}

// Match command by path. eg. "top:sub"
func (c *Command) MatchByPath(path string) *Command {
	return c.Match(splitPath2names(path))
}

// initialize works for the command
func (c *Command) initialize() {
	if c.initialized {
		return
	}

	// check command name
	cName := c.goodName()

	Debugf("initialize the command '%s'", cName)

	c.initialized = true
	c.pathNames = append(c.pathNames, cName)

	// init core
	c.initCore(cName)

	// init commandBase
	c.initCommandBase()

	// load common subs
	if len(c.Subs) > 0 {
		for _, sub := range c.Subs {
			c.AddCommand(sub)
		}
	}

	// init for cmd Arguments
	c.Arguments.SetName(cName)
	c.Arguments.SetValidateNum(gOpts.strictMode)

	// init for cmd Flags
	c.Flags.InitFlagSet(cName)
	// c.Flags.SetOption(cName)
	// c.Flags.FSet().SetOutput(c.Flags.out)
	// c.Flags.FSet().Usage = func() { // call on exists "-h" "--help"
	// 	Logf(VerbDebug, "render help on exists '-h|--help' or has unknown flag")
	// 	c.ShowHelp()
	// }

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

// init core
func (c *Command) initCore(cName string) {
	Logf(VerbCrazy, "init command c.core for the command: %s", cName)

	c.core.cmdLine = CLI
	c.AddVars(c.innerHelpVars())
	c.AddVars(map[string]string{
		"cmd": cName,
		// binName with command
		"binWithCmd": c.binName + " " + cName,
		// binFile with command
		"fullCmd": c.binFile + " " + cName,
	})
}

func (c *Command) initCommandBase() {
	Logf(VerbCrazy, "init command c.commandBase for the command: %s", c.Name)

	c.commandBase.cmdNames = make(map[string]int)
	c.commandBase.commands = make(map[string]*Command)
	// set an default value.
	c.commandBase.nameMaxWidth = 12
	// c.commandBase.cmdAliases = make(maputil.Aliases)
	c.commandBase.cmdAliases = structs.NewAliases(aliasNameCheck)
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

var errCallRunOnApp = errors.New("c.Run() method can only be called in standalone mode")
var errCallRunOnSub = errors.New("c.Run() cannot allow call at subcommand")

// MustRun Alone the current command, will panic on error
//
// Usage:
//	// run with os.Args
//	cmd.MustRun(nil)
//	cmd.MustRun(os.Args[1:])
//	// custom args
//	cmd.MustRun([]string{"-a", ...})
func (c *Command) MustRun(args []string) {
	if err := c.Run(args); err != nil {
		color.Error.Println("ERROR:", err.Error())
		panic(err)
	}
}

// Run standalone running the command
//
// Usage:
//	// run with os.Args
//	cmd.Run(nil)
//	cmd.Run(os.Args[1:])
//	// custom args
//	cmd.Run([]string{"-a", ...})
func (c *Command) Run(args []string) (err error) {
	if c.app != nil {
		return errCallRunOnApp
	}

	if c.parent != nil {
		return errCallRunOnSub
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
	if !c.HasHook(EvtCmdError) {
		c.On(EvtCmdError, defaultErrHandler)
	}

	// binding global options
	Debugf("global options will binding to c.Flags on standalone mode")
	bindingCommonGOpts(&c.Flags)

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

		color.Error.Tips("option error - %s", err.Error())
		return nil
	}

	// remaining args
	if c.standalone {
		if gOpts.showHelp {
			c.ShowHelp()
			return
		}

		c.Fire(EvtGOptionsParsed, args)
	}

	c.Fire(EvtCmdOptParsed, c.Name, args)
	Debugf("cmd: %s - remaining args on options parsed: %v", c.Name, args)

	// find sub command
	if len(args) > 0 {
		name := args[0]

		// ensure is not an option
		if name[0] != '-' {
			name = c.ResolveAlias(name)

			// is valid sub command
			if sub, has := c.Command(name); has {
				// loop find sub...command and run it.
				return sub.innerDispatch(args[1:])
			}
		}
	}

	// defaultCommand is not empty.
	name := c.defaultCommand
	if name != "" {
		// is valid sub command
		if sub, has := c.Command(name); has {
			Debugf("will run the default command '%s' of the '%s'", name, c.Name)

			// run the default command
			return sub.innerExecute(args, true)
		}

		return fmt.Errorf("the default command '%s' is invalid", name)
	}

	// not set command func and has sub commands.
	if c.Func == nil && len(c.commands) > 0 {
		Logf(VerbWarn, "cmd: %s - c.Func is empty, but has sub commands, will render help list", c.Name)
		c.ShowHelp()
		return err
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
	// Debugf("cmd: %s - option flags on after format: %v", c.Name, args)

	// NOTICE: disable output internal error message on parse flags
	// c.FSet().SetOutput(ioutil.Discard)
	Debugf("cmd: %s - will parse options by args: %v", c.Name, args)

	// parse options, don't contains command name.
	if err = c.Parse(args); err != nil {
		Logf(VerbCrazy, "'%s' - parse options  err: <red>%s</>", c.Name, err.Error())
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
	Debugf("cmd: %s - collect and binding named argument", c.Name)
	if err := c.ParseArgs(args); err != nil {
		c.Fire(EvtCmdError, err)
		Logf(VerbCrazy, "binding command '%s' arguments err: <red>%s</>", c.Name, err.Error())
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
 * command help
 *************************************************************/

// CmdHelpTemplate help template for a command
var CmdHelpTemplate = `{{.Desc}}
{{if .Cmd.NotStandalone}}
<comment>Name:</> {{.Cmd.Name}}{{if .Cmd.Aliases}} (alias: <info>{{.Cmd.AliasesString}}</>){{end}}{{end}}
<comment>Usage:</> {$binName} [global options] {{if .Cmd.NotStandalone}}<info>{{.Cmd.Path}}</> {{end}}[--option ...] [arguments ...]
{{if .GOpts}}
<comment>Global Options:</>
{{.GOpts}}{{end}}{{if .Options}}
<comment>Options:</>
{{.Options}}{{end}}{{if .Cmd.Args}}
<comment>Arguments:</>{{range $a := .Cmd.Args}}
  <info>{{$a.HelpName | printf "%-12s"}}</>{{$a.Desc | ucFirst}}{{if $a.Required}}<red>*</>{{end}}{{end}}
{{end}}{{ if .Subs }}
<comment>Sub Commands:</>{{range $n,$c := .Subs}}
  <info>{{$c.Name | paddingName }}</> {{$c.Desc}}{{if $c.Aliases}} (alias: <green>{{ join $c.Aliases ","}}</>){{end}}{{end}}
{{end}}{{if .Cmd.Examples}}
<comment>Examples:</>
{{.Cmd.Examples}}{{end}}{{if .Cmd.Help}}
<comment>Help:</>
{{.Cmd.Help}}{{end}}`

// ShowHelp show command help info
func (c *Command) ShowHelp() {
	Debugf("render the command '%s' help information", c.Name)

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

	vars := map[string]interface{}{
		"Cmd":  c,
		"Subs": c.commands,
		// global options
		// - on standalone, will not init c.core.gFlags
		"GOpts": nil,
		// parse options to string
		"Options": c.Flags.String(),
		// always upper first char
		"Desc": c.Desc,
	}

	if c.NotStandalone() {
		vars["GOpts"] = c.GFlags().String()
	}

	// render help message
	str := helper.RenderText(CmdHelpTemplate, vars, template.FuncMap{
		"paddingName": func(n string) string {
			return strutil.PadRight(n, " ", c.nameMaxWidth)
		},
	})

	// parse gcli help vars then print help
	// fmt.Printf("%#v\n", s)
	color.Print(c.ReplaceVars(str))
}

/*************************************************************
 * helper methods
 *************************************************************/

// GFlags get global flags
func (c *Command) GFlags() *Flags {
	// 如果先注册S子命令到一个命令A中，再将A注册到应用App。此时，S.gFlags 就是空的。
	// If you first register the S subcommand to a command A, then register A to the application App.
	// At this time, S.gFlags is empty.
	if c.gFlags == nil {
		if c.parent == nil {
			return nil
		}

		// inherit from parent command.
		c.core.gFlags = c.parent.GFlags()
	}

	return c.gFlags
}

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

	if !goodCmdName.MatchString(name) {
		panicf("the command name '%s' is invalid, must match: %s", name, regGoodCmdName)
	}

	// update name
	c.Name = name
	return name
}

// Fire event handler by name
func (c *Command) Fire(event string, data ...interface{}) {
	Debugf("cmd: %s - trigger the event: <mga>%s</>", c.Name, event)

	c.Hooks.Fire(event, c, data)
}

// On add hook handler for a hook event
func (c *Command) On(name string, handler HookFunc) {
	Debugf("cmd: %s - register hook: %s", c.Name, name)

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
