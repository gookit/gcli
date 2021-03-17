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
	"regexp"
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
	// match an good option, argument name
	regGoodName = `^[a-zA-Z][\w-]*$`
	// match an good command name
	// regGoodCmdName = `^[a-zA-Z][\w:-]*$`
	regGoodCmdName = `^[a-zA-Z][\w-]*$`
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
	EvtAppInit   = "app.init"
	EvtAppBefore = "app.run.before"
	EvtAppAfter  = "app.run.after"
	EvtAppError  = "app.run.error"

	EvtCmdInit = "cmd.init"

	EvtCmdNotFound  = "cmd.not.found"
	EvtCmdOptParsed = "cmd.opts.parsed"

	EvtCmdBefore = "cmd.run.before"
	EvtCmdAfter  = "cmd.run.after"
	EvtCmdError  = "cmd.run.error"

	EvtAppPrepareAfter = "app.prepare.after"

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

	// good name for option and argument
	goodName = regexp.MustCompile(regGoodName)
	// match an good command name
	goodCmdName = regexp.MustCompile(regGoodCmdName)
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
func CommitID() string {
	return commitID
}

// Verbose returns verbose level
func Verbose() VerbLevel {
	return gOpts.Verbose()
}

// SetCrazyMode level
func SetCrazyMode() {
	gOpts.SetVerbose(VerbCrazy)
}

// SetDebugMode level
func SetDebugMode() {
	gOpts.SetVerbose(VerbDebug)
}

// SetQuietMode level
func SetQuietMode() {
	gOpts.SetVerbose(VerbQuiet)
}

// SetVerbose level
func SetVerbose(verbose VerbLevel) {
	gOpts.SetVerbose(verbose)
}

// ResetVerbose level
func ResetVerbose() {
	gOpts.SetVerbose(DefaultVerb)
}

// StrictMode get is strict mode
func StrictMode() bool {
	return gOpts.strictMode
}

// SetStrictMode for parse flags
func SetStrictMode(strict bool) {
	gOpts.SetStrictMode(strict)
}

// binding global options
func bindingCommonGOpts(fs *Flags) {
	// up: allow use int and string.
	fs.VarOpt(&gOpts.verbose, "verbose", "", "Set error reporting level(quiet 0 - 5 crazy)")

	fs.BoolOpt(&gOpts.showHelp, "help", "h", false, "Display the help information")
	fs.BoolOpt(&gOpts.NoColor, "no-color", "", gOpts.NoColor, "Disable color when outputting message")
	fs.BoolOpt(&gOpts.noProgress, "no-progress", "", gOpts.noProgress, "Disable display progress message")
	fs.BoolOpt(&gOpts.noInteractive, "no-interactive", "", gOpts.noInteractive, "Disable interactive confirmation operations")
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
