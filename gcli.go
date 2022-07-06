// Package gcli is a simple to use command line application and tool library.
//
// Contains: cli app, flags parse, interact, progress, data show tools.
//
// Source code and other details for the project are available at GitHub:
// 		https://github.com/gookit/gcli
//
// Usage please refer examples and see README
package gcli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	// OK success exit code
	OK = 0
	// ERR error exit code
	ERR = 2
	// GOON prepare run successful, goon run command
	GOON = -1
	// CommandSep char
	CommandSep = ":"
	// HelpCommand name
	HelpCommand = "help"
)

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
	EvtAppInit = "app.init"

	EvtAppPrepareAfter = "app.prepare.after"

	EvtAppRunBefore = "app.run.before"
	EvtAppRunAfter  = "app.run.after"
	EvtAppRunError  = "app.run.error"

	EvtCmdInit = "cmd.init"

	// EvtCmdNotFound app or sub command not found
	EvtCmdNotFound = "cmd.not.found"
	// EvtAppCmdNotFound app command not found
	EvtAppCmdNotFound = "app.cmd.not.found"
	// EvtCmdSubNotFound sub command not found
	EvtCmdSubNotFound = "cmd.sub.not.found"

	EvtCmdOptParsed = "cmd.opts.parsed"

	// EvtCmdRunBefore cmd run
	EvtCmdRunBefore = "cmd.run.before"
	EvtCmdRunAfter  = "cmd.run.after"
	EvtCmdRunError  = "cmd.run.error"

	// EvtCmdExecBefore cmd exec
	EvtCmdExecBefore = "cmd.exec.before"
	EvtCmdExecAfter  = "cmd.exec.after"
	EvtCmdExecError  = "cmd.exec.error"

	EvtGOptionsParsed = "gcli.gopts.parsed"
	// EvtStop   = "stop"
)

var (
	// CLI create an default instance
	CLI = newCmdLine()
	// DefaultVerb the default verbose level
	DefaultVerb = VerbError
	// global options
	gOpts = newDefaultGOptions()
	// Version the gcli version
	version = "3.0.0"
	// CommitID the gcli last commit ID
	commitID = "z20210214"
)

// init
func init() {
	// set verbose from ENV var.
	envVerb := os.Getenv("GCLI_VERBOSE")
	if envVerb != "" {
		_ = gOpts.verbose.Set(envVerb)
	}
}

// GOpts get the global options
func GOpts() *GOptions {
	return gOpts
}

// ResetGOpts instance
func ResetGOpts() {
	gOpts = newDefaultGOptions()
}

// Version of the gcli
func Version() string {
	return version
}

// CommitID of the gcli
func CommitID() string { return commitID }

// Verbose returns verbose level
func Verbose() VerbLevel { return gOpts.Verbose() }

// SetCrazyMode level
func SetCrazyMode() { gOpts.SetVerbose(VerbCrazy) }

// SetDebugMode level
func SetDebugMode() { gOpts.SetVerbose(VerbDebug) }

// SetQuietMode level
func SetQuietMode() { gOpts.SetVerbose(VerbQuiet) }

// SetVerbose level
func SetVerbose(verbose VerbLevel) { gOpts.SetVerbose(verbose) }

// ResetVerbose level
func ResetVerbose() { gOpts.SetVerbose(DefaultVerb) }

// StrictMode get is strict mode
func StrictMode() bool { return gOpts.strictMode }

// SetStrictMode for parse flags
func SetStrictMode(strict bool) { gOpts.SetStrictMode(strict) }

// IsGteVerbose get is strict mode
func IsGteVerbose(verb VerbLevel) bool { return gOpts.verbose >= verb }

// IsDebugMode get is debug mode
func IsDebugMode() bool { return gOpts.verbose >= VerbDebug }

// Commander interface
type Commander interface {
	Value(string) (interface{}, bool)
	SetValue(string, interface{})
}

/*************************************************************************
 * global options
 *************************************************************************/

// GOptions global flag options
type GOptions struct {
	NoColor  bool
	verbose  VerbLevel // message report level
	showVer  bool
	showHelp bool
	// TODO Run application an interactive shell environment
	inShell bool
	// dont display progress
	noProgress bool
	// close interactive confirm
	noInteractive bool
	// TODO auto format shorts `-a` to POSIX or UNIX style.
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

func newDefaultGOptions() *GOptions {
	return &GOptions{
		strictMode: true,
		// init error level.
		verbose: DefaultVerb,
	}
}

// Verbose value
func (g *GOptions) Verbose() VerbLevel {
	return g.verbose
}

// SetVerbose value
func (g *GOptions) SetVerbose(verbose VerbLevel) {
	g.verbose = verbose
}

// SetStrictMode option
func (g *GOptions) SetStrictMode(strictMode bool) {
	g.strictMode = strictMode
}

// NoInteractive value
func (g *GOptions) NoInteractive() bool {
	return g.noInteractive
}

// NoProgress value
func (g *GOptions) NoProgress() bool {
	return g.noProgress
}

func (g *GOptions) bindingFlags(fs *Flags) {
	// up: allow use int and string.
	fs.VarOpt(&g.verbose, "verbose", "", "Set error reporting level(quiet 0 - 5 crazy)")

	fs.BoolOpt(&g.inShell, "ishell", "", false, "Run in an interactive shell environment(`TODO`)")
	fs.BoolOpt(&g.showHelp, "help", "h", false, "Display the help information")
	fs.BoolOpt(&g.NoColor, "no-color", "", g.NoColor, "Disable color when outputting message")
	fs.BoolOpt(&g.noProgress, "no-progress", "", g.noProgress, "Disable display progress message")
	fs.BoolOpt(&g.noInteractive, "no-interactive", "", g.noInteractive, "Disable interactive confirmation operations")
}

/*************************************************************************
 * verbose level
 *************************************************************************/

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

// Name verbose level to string.
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
			*vl = DefaultVerb
		} else { // 0 - 5
			*vl = VerbLevel(iv)
		}

		return nil
	}

	// string: level name.
	*vl = name2verbLevel(value)
	return nil
}
