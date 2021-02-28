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

	EvtCmdInit   = "cmd.init"
	EvtCmdBefore = "cmd.run.before"
	EvtCmdAfter  = "cmd.run.after"
	EvtCmdError  = "cmd.run.error"

	EvtAppPrepareAfter = "app.prepare.after"

	EvtGlobalOptionParsed = "gcli.global.opts.parsed"
	// EvtStop   = "stop"
)

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

var (
	// Version the gCli version
	Version = "3.0.0"

	// stdApp store default application instance
	stdApp *App
	// an empty argument
	emptyArg = &Argument{}
	// good name for option and argument
	goodName = regexp.MustCompile(regGoodName)
	// match an good command name
	goodCmdName = regexp.MustCompile(regGoodCmdName)

	// global options
	gOpts = &GlobalOpts{
		strictMode: true,
		// init error level.
		verbose: VerbError,
	}

	// CLI create an default instance
	CLI = newCmdLine()
)

// init
func init() {
	// don't display date on print log
	// log.SetFlags(0)
	// workDir, _ := os.Getwd()
	// CLI.workDir = workDir

	// // binName will contains work dir path on windows
	// if envutil.IsWin() {
	// 	CLI.binName = strings.Replace(CLI.binName, workDir+"\\", "", 1)
	// }

	// set verbose from ENV var.
	envVerb := os.Getenv("GCLI_VERBOSE")
	if envVerb != "" {
		_= gOpts.verbose.Set(envVerb)
	}
}

// InitStdApp create the default cli app.
func InitStdApp(fn ...func(a *App)) *App {
	stdApp = NewApp(fn...)
	return stdApp
}

// StdApp get the default std app
func StdApp() *App {
	return stdApp
}

// GOpts get the global options
func GOpts() GlobalOpts {
	return *gOpts
}

// Verbose returns verbose level
func Verbose() VerbLevel {
	return gOpts.verbose
}

// SetDebugMode level
func SetDebugMode() {
	gOpts.verbose = VerbDebug
}

// SetQuietMode level
func SetQuietMode() {
	gOpts.verbose = VerbQuiet
}

// SetVerbose level
func SetVerbose(verbose VerbLevel) {
	gOpts.verbose = verbose
}

// StrictMode get is strict mode
func StrictMode() bool {
	return gOpts.strictMode
}

// SetStrictMode for parse flags
func SetStrictMode(strict bool) {
	gOpts.strictMode = strict
}

func bindingCommonGOpts(fs *Flags) {
	// binding global options
	// fs.UintOpt(&gOpts.verbose, "verbose", "", gOpts.verbose, "Set error reporting level(quiet 0 - 5 crazy)")
	// up: allow use int and string.
	fs.VarOpt(&gOpts.verbose, "verbose", "", "Set error reporting level(quiet 0 - 5 crazy)")

	fs.BoolOpt(&gOpts.showHelp, "help", "h", false, "Display the help information")
	fs.BoolOpt(&gOpts.NoColor, "no-color", "", gOpts.NoColor, "Disable color when outputting message")
	fs.BoolOpt(&gOpts.noProgress, "no-progress", "", gOpts.noProgress, "Disable display progress message")
	fs.BoolOpt(&gOpts.noInteractive, "no-interactive", "", gOpts.noInteractive, "Disable interactive confirmation operations")
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