package gcli

// GlobalOpts global flags
type GlobalOpts struct {
	verbose  uint // message report level
	noColor  bool
	showVer  bool
	showHelp bool
	// dont display progress
	noProgress bool
	// close interactive confirm
	noInteractive bool
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
	// Creator for create new command
	Creator() *Command
	// BindFlags for the command
	BindFlags(c *Command)
	// Execute(c *Command, args []string) error
	Run(c *Command, args []string) error
}
