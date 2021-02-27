package gcli

import (
	"fmt"
	"strconv"
	"strings"
)

// VerbLevel type.
type VerbLevel uint

// Int verbose level to int.
func (vl VerbLevel) Int() int {
	return int(vl)
}

// String verbose level to string.
func (vl VerbLevel) String() string {
	return fmt.Sprintf("%d=%s", vl, vl.Name())
}

// Upper verbose level to string.
func (vl VerbLevel) Upper() string {
	return strings.ToUpper(vl.Name())
}

// String verbose level to string.
func (vl VerbLevel) Name() string {
	switch vl {
	case VerbQuiet:
		return "quiet"
	case VerbError:
		return "error"
	case VerbWarn:
		return "warn"
	case VerbInfo:
		return "info"
	case VerbDebug:
		return "debug"
	case VerbCrazy:
		return "crazy"
	}
	return "unknown"
}

// Set value from option binding.
func (vl *VerbLevel) Set(value string) error {
	// int: level value.
	if iv, err := strconv.Atoi(value); err == nil {
		if iv > VerbCrazy.Int() {
			*vl = VerbCrazy
		} else if iv < 0 { // fallback to default level.
			*vl = VerbError
		} else { // 0 - 5
			*vl = VerbLevel(iv)
		}

		return nil
	}

	// string: level name.
	*vl = name2verbLevel(value)
	return nil
}

// constants for error level (quiet 0 - 5 crazy)
const (
	VerbQuiet VerbLevel = iota // don't report anything
	VerbError                  // reporting on error, default level.
	VerbWarn
	VerbInfo
	VerbDebug
	VerbCrazy
)

// constants for hooks event, there are default allowed event names
const (
	EvtAppInit   = "app.init"
	EvtAppBefore = "app.run.before"
	EvtAppAfter  = "app.run.after"
	EvtAppError  = "app.run.error"

	EvtCmdInit   = "cmd.init"
	EvtCmdBefore = "cmd.run.before"
	EvtCmdAfter  = "cmd.run.after"
	EvtCmdError  = "cmd.run.error"

	EvtAppPrepareAfter = "app.prepare.after"
	// EvtStop   = "stop"
)

const maxFunc = 64

// GlobalOpts global flags
type GlobalOpts struct {
	verbose  VerbLevel // message report level
	NoColor  bool
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

// Commander interface definition
type Commander interface {
	// Creator for create new command
	Creator() *Command
	// Config bind Flags or Arguments for the command
	Config(c *Command)
	// Execute the command
	Execute(c *Command, args []string) error
}

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

/*************************************************************************
 * options: some special flag vars
 * - implemented flag.Value interface
 *************************************************************************/

// Ints The int flag list, implemented flag.Value interface
type Ints []int

// String to string
func (s *Ints) String() string {
	return fmt.Sprintf("%v", *s)
}

// Set new value
func (s *Ints) Set(value string) error {
	intVal, err := strconv.Atoi(value)
	if err == nil {
		*s = append(*s, intVal)
	}

	return err
}

// Strings The string flag list, implemented flag.Value interface
type Strings []string

// String to string
func (s *Strings) String() string {
	return fmt.Sprintf("%v", *s)
}

// Set new value
func (s *Strings) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// Booleans The bool flag list, implemented flag.Value interface
type Booleans []bool

// String to string
func (s *Booleans) String() string {
	return fmt.Sprintf("%v", *s)
}

// Set new value
func (s *Booleans) Set(value string) error {
	boolVal, err := strconv.ParseBool(value)
	if err == nil {
		*s = append(*s, boolVal)
	}

	return err
}

// EnumString The string flag list, implemented flag.Value interface
type EnumString struct {
	val  string
	enum []string
}

// String to string
func (s *EnumString) String() string {
	return s.val
}

// Set new value, will check value is right
func (s *EnumString) Set(value string) error {
	var ok bool
	for _, item := range s.enum {
		if value == item {
			ok = true
			break
		}
	}

	if !ok {
		return fmt.Errorf("value must one of the: %v", s.enum)
	}

	return nil
}
