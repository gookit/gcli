package cliapp

import (
	"html/template"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Commander
type Commander interface {
	Configure() *Command
	Execute(app *Application, args []string) int
	//Fn(cmd *Command, args []string) int
}

// CmdExecutor
// type CmdExecutor func(Context) int

// CmdHandler
// type CmdHandler func(app *Application, args []string) int
type CmdHandler Command

// CmdAliases
type CmdAliases []string

// to string
func (a *CmdAliases) String() string {
	return strings.Join(*a, ",")
}

// Command a cli command
type Command struct {
	// Name is the command name.
	Name string

	// Aliases is the command name's alias names
	Aliases CmdAliases

	// Description is the command description for 'go help'
	Description string

	// Flags is a set of flags specific to this command.
	Flags flag.FlagSet

	// CustomFlags indicates that the command will do its own flag parsing.
	CustomFlags bool

	// A callback func to runs the command.
	// The args are the arguments after the command name.
	Fn func(cmd *Command, args []string) int

	// Help is the help message shown in the 'go help <this-command>' output.
	Help template.HTML

	// Examples some usage example display
	Examples template.HTML

	// vars you can add some vars map for render help info
	Vars map[string]string

	// Options

	// ArgList arguments description [name]description
	ArgList map[string]string

	// application
	app *Application
}

// Option a command option
type Option struct {
	// Name is the Option name. eg 'name' -> '--name'
	Name string

	// Short is the Option short name. eg 'n' -> '-n'
	Short string

	// Description is the option description message
	Description string
}

// ShowHelp @notice not used
func (c *Command) ShowHelp() {
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(string(c.Description)))
	// fmt.Fprintf(os.Stderr, "Usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n\n", c.Help)
	fmt.Fprintf(os.Stderr, "Example:%s\n\n", c.Examples)

	os.Exit(0)
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as import path.
func (c *Command) Runnable() bool {
	return c.Fn != nil
}

// Configure
func (c *Command) Configure() *Command {
	return c
}

// Execute do execute the command
func (c *Command) Execute(app *Application, args []string) int {
	return c.Fn(c, args)
}

// Application
func (c *Command) App() *Application {
	return app
}

// GetArgs get args
func (c *Command) GetArgs() []string {
	return c.Flags.Args()
}

// Arg get arg
func (c *Command) Arg(i int) string {
	return c.Flags.Arg(i)
}

// IntOption
func (c *Command) IntOpt(p *int, name string, defaultValue int, description string) *Command {
	c.Flags.IntVar(p, name, defaultValue, description)
	return c
}

// StrOption
func (c *Command) StrOpt(p *string, name string, defaultValue string, description string) *Command {
	c.Flags.StringVar(p, name, defaultValue, description)
	return c
}

// AliasesStr
func (c *Command) AliasesStr() string {
	return c.Aliases.String()
}

// AddVars add multi tpl vars
func (c *Command) AddVars(vars map[string]string) {
	// first init
	if c.Vars == nil {
		c.Vars = vars
		return
	}

	for n, v := range vars {
		c.Vars[n] = v
	}
}

// PrintDefaults prints, to standard error unless configured otherwise, the
// default values of all defined command-line flags in the set. See the
// documentation for the global function PrintDefaults for more information.
// NOTICE: the func is copied from package 'flag', func 'PrintDefaults'
func (c *Command) ParseDefaults() string {
	var ss []string

	c.Flags.VisitAll(func(fg *flag.Flag) {
		var s string

		// is short option
		if len(fg.Name) == 1 {
			s = fmt.Sprintf("  -%s", fg.Name) // Two spaces before -; see next two comments.
		} else {
			s = fmt.Sprintf("  --%s", fg.Name) // Two spaces before -; see next two comments.
		}

		name, usage := flag.UnquoteUsage(fg)
		if len(name) > 0 {
			s += " " + name
		}
		// Boolean flags of one ASCII letter are so common we
		// treat them specially, putting their usage on the same line.
		if len(s) <= 4 { // space, space, '-', 'x'.
			s += "\t"
		} else {
			// Four spaces before the tab triggers good alignment
			// for both 4- and 8-space tab stops.
			s += "\n    \t"
		}
		s += strings.Replace(usage, "\n", "\n    \t", -1)

		if !isZeroValue(fg, fg.DefValue) {
			if _, ok := fg.Value.(*stringValue); ok {
				// put quotes on the value
				s += fmt.Sprintf(" (default %q)", fg.DefValue)
			} else {
				s += fmt.Sprintf(" (default %v)", fg.DefValue)
			}
		}

		ss = append(ss, s)
		// fmt.Fprint(fgs.Output(), s, "\n")
	})

	return strings.Join(ss, "\n")
}
