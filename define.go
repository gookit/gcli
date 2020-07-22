package gcli

import (
	"fmt"

	"github.com/gookit/goutil/mathutil"
)

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
	// Creator for create new command
	Creator() *Command
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
	opt1 string
}

func (uc *UserCommand) 	Creator() *Command {
	return NewCommand("test", "desc message")
}

func (uc *UserCommand) BindFlags(c *Command)  {
	c.StrOpt(&uc.opt1, "opt", "o", "", "desc")
}

func (uc *UserCommand) Run(c *Command, args []string) error {
	return nil
}


// Value data store
type Value struct {
	// V value
	V interface{}
}

// Reset value
func (v *Value) Reset() {
	v.V = nil
}

// Val get
func (v Value) Val() interface{} {
	return v.V
}

// Int value
func (v Value) Int() int {
	if v.V == nil {
		return 0
	}

	return mathutil.MustInt(v.V)
}

// Int64 value
func (v Value) Int64() int64 {
	if v.V == nil {
		return 0
	}

	return mathutil.MustInt64(v.V)
}

// Bool value
func (v Value) Bool() bool {
	if v.V == nil {
		return false
	}

	if bl, ok := v.V.(bool); ok {
		return bl
	}
	return false
}

// String value
func (v Value) String() string {
	if v.V == nil {
		return ""
	}

	if str, ok := v.V.(string); ok {
		return str
	}

	return fmt.Sprintf("%v", v.V)
}

// Strings value
func (v Value) Strings() (ss []string) {
	if v.V == nil {
		return
	}

	if ss, ok := v.V.([]string); ok {
		return ss
	}
	return
}

// IsEmpty value
func (v Value) IsEmpty() bool {
	return v.V == nil
}
