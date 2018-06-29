package cliapp

import (
	"html/template"
	"flag"
	"strings"
	"github.com/golangkit/cliapp/utils"
)

// Commander
type Commander interface {
	Init() *Command
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

	// Flags(Options) is a set of flags specific to this command.
	Flags flag.FlagSet

	// CustomFlags indicates that the command will do its own flag parsing.
	CustomFlags bool

	// A callback func to runs the command.
	// The args are the arguments after the command name.
	Fn func(cmd *Command, args []string) int

	// Hooks can setting some hooks func on running.
	// names: "init", "before", "after"
	Hooks map[string]func(cmd *Command)

	// Help is the help message text
	Help template.HTML

	// Examples some usage example display
	Examples template.HTML

	// vars you can add some vars map for render help info
	Vars map[string]string

	// ArgList arguments description [name]description
	ArgList map[string]string

	// application
	app *Application

	// mark is alone running.
	alone bool

	// shortcuts storage shortcuts for command options(Flags)
	// [short -> lang] eg. {"n": "name", "o": "opt"}
	shortcuts map[string]string
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as import path.
func (c *Command) Runnable() bool {
	return c.Fn != nil
}

// Init
func (c *Command) Init() *Command {
	if len(c.Description) > 0 {
		c.Description = utils.UpperFirst(c.Description)

		// if contains help var "{$cmd}"
		if strings.Contains(c.Description, "{$cmd}") {
			c.Description = strings.Replace(c.Description, "{$cmd}", c.Name, -1)
		}
	}

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

// GetVar get a help var by name
func (c *Command) GetVar(name string) string {
	if v, ok := c.Vars[name]; ok {
		return v
	}

	return ""
}
