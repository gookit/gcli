package cli

import (
    "html/template"
    "flag"
    "fmt"
    "os"
    "strings"
)

// Commander
type Commander interface {
    Name() string
    Execute(app *App, args []string) int
}

// CmdExecutor
// type CmdExecutor func(Context) int

// CmdHandler
// type CmdHandler func(app *App, args []string) int
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

    // UsageLine is the one-line usage message.
    // The first word in the line is taken to be the command name.
    UsageLine string

    // Flags is a set of flags specific to this command.
    Flags flag.FlagSet

    // CustomFlags indicates that the command will do its own flag parsing.
    CustomFlags bool

    // Help is the help message shown in the 'go help <this-command>' output.
    Help template.HTML

    // Examples some usage example display
    Examples template.HTML

    // Run runs the command.
    // The args are the arguments after the command name.
    Execute func(cmd *Command, args []string) int

    // Options

    // arguments [name]description
    Args []string

    // application
    app *App
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
    fmt.Fprintf(os.Stderr, "Usage: %s\n\n", c.UsageLine)
    fmt.Fprintf(os.Stderr, "%s\n\n", c.Help)
    fmt.Fprintf(os.Stderr, "Example:%s\n\n", c.Examples)

    os.Exit(0)
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as importpath.
func (c *Command) Runnable() bool {
    return c.Execute != nil
}

// NewFlags
//func (c *Command) NewFlags() *flag.FlagSet {
//   fst := flag.NewFlagSet(c.Name, flag.ContinueOnError)
//
//   c.Flags = *fst
//
//   return  fst
//}

// App
func (c *Command) App() *App {
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
