package gcli

import (
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v2/helper"
)

/*************************************************************
 * command run
 *************************************************************/

// do parse option flags, remaining is cmd args
func (c *Command) parseFlags(args []string) (ss []string, err error) {
	// strict format options
	if gOpts.strictMode && len(args) > 0 {
		args = strictFormatArgs(args)
	}

	// fix and compatible
	args = moveArgumentsToEnd(args)

	Logf(VerbDebug, "flags on after format: %v", args)

	// NOTICE: disable output internal error message on parse flags
	// c.FSet().SetOutput(ioutil.Discard)

	// parse options, don't contains command name.
	if err = c.Flags.Parse(args); err != nil {
		return
	}

	return c.Flags.RawArgs(), nil
}

// prepare: before execute the command
func (c *Command) prepare(_ []string) (status int, err error) {
	return
}

// do execute the command
func (c *Command) execute(args []string) (err error) {
	c.Fire(EvtCmdBefore, args)

	// collect and binding named args
	if err := c.ParseArgs(args); err != nil {
		c.Fire(EvtCmdError, err)
		return err
	}

	// call command handler func
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

// Fire event handler by name
func (c *Command) Fire(event string, data interface{}) {
	Logf(VerbDebug, "command '%s' trigger the event: <mga>%s</>", c.Name, event)

	c.Hooks.Fire(event, c, data)
}

// On add hook handler for a hook event
func (c *Command) On(name string, handler HookFunc) {
	Logf(VerbDebug, "command '%s' add hook: %s", c.Name, name)

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

/*************************************************************
 * alone running
 *************************************************************/

var errCallRun = errors.New("this method can only be called in standalone mode")

// MustRun Alone the current command, will panic on error
func (c *Command) MustRun(inArgs []string) {
	if err := c.Run(inArgs); err != nil {
		color.Error.Println("Run command error: %s", err.Error())
		panic(err)
	}
}

// Run Alone the current command
func (c *Command) Run(inArgs []string) (err error) {
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

	// TODO parse global options

	// init the command
	c.initialize()

	// add default error handler.
	c.Hooks.AddOn(EvtCmdError, defaultErrHandler)

	// check input args
	if len(inArgs) == 0 {
		inArgs = os.Args[1:]
	}

	// if Command.CustomFlags=true, will not run Flags.Parse()
	if !c.CustomFlags {
		// contains keywords "-h" OR "--help" on end
		if c.hasHelpKeywords() {
			c.ShowHelp()
			return
		}

		// if CustomFlags=true, will not run Flags.Parse()
		inArgs, err = c.parseFlags(inArgs)
		if err != nil {
			// ignore flag.ErrHelp error
			if err == flag.ErrHelp {
				err = nil
			}
			return
		}
	}

	return c.execute(inArgs)
}

/*************************************************************
 * display cmd help
 *************************************************************/

// help template for a command
var commandHelp = `{{.UseFor}}
{{if .Cmd.NotAlone}}
<comment>Name:</> {{.Cmd.Name}}{{if .Cmd.Aliases}} (alias: <info>{{.Cmd.AliasesString}}</>){{end}}{{end}}
<comment>Usage:</> {$binName} [Global Options...] {{if .Cmd.NotAlone}}<info>{{.Cmd.Name}}</> {{end}}[--option ...] [arguments ...]

<comment>Global Options:</>
      <info>--verbose</>     Set error reporting level(quiet 0 - 4 debug)
      <info>--no-color</>    Disable color when outputting message
  <info>-h, --help</>        Display this help information
{{if .Options}}
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
	// commandHelp = color.ReplaceTag(commandHelp)
	// clear space and empty new line
	if c.Examples != "" {
		c.Examples = strings.Trim(c.Examples, "\n") + "\n"
	}

	// clear space and empty new line
	if c.Help != "" {
		c.Help = strings.Join([]string{strings.TrimSpace(c.Help), "\n"}, "")
	}

	// render and output help info
	// RenderTplStr(os.Stdout, commandHelp, map[string]interface{}{
	// render but not output
	s := helper.RenderText(commandHelp, map[string]interface{}{
		"Cmd": c,
		// parse options to string
		"Options": c.Flags.String(),
		// always upper first char
		"UseFor": c.UseFor,
	}, nil)

	// parse help vars then print help
	color.Print(c.ReplaceVars(s))
	// fmt.Printf("%#v\n", s)
}
