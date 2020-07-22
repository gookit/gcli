package gcli

import "github.com/gookit/gcli/v2"

// GlobalOpts global flags
type GlobalOpts struct {
	verbose  uint // message report level
	noColor  bool
	showVer  bool
	showHelp bool
	// StrictMode use strict mode for parse flags
	// If True(default):
	// 	- short opt must be begin "-", long opt must be begin "--"
	//	- will convert like "-ab" to "-a -b"
	// 	- will check invalid arguments, like to many arguments
	strictMode bool
	// command auto completion mode.
	// eg "./cli --cmd-completion [COMMAND --OPT ARG]"
	inCompletion bool
}

// Runner interface
type Runner interface {
	// Config(c *Command)
	Run(c *Command, args []string) error
}

// CmdFunc definition
type CmdFunc func(c *Command, args []string) error

// Run implement the Runner interface
func (f CmdFunc) Run(c *Command, args []string) error {
	return f(c, args)
}

// Commander interface definition
type Commander interface {
	// BindFlags for the command
	BindFlags(c *Command)
	// Execute(c *Command, args []string) error
	Run(c *Command, args []string) error
}

func newUserCommand()  {
	app := NewApp()

	app.AddCommander(&UserCommand{})
}

type UserCommand struct {
	cmd *gcli.Command
	opt1 string
}

func (uc *UserCommand) BindFlags(c *Command)  {
	c.StrOpt(&uc.opt1, "opt", "o", "", "desc")
}

func (uc *UserCommand) Run(c *Command, args []string) error {
	return nil
}