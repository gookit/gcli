// Package gcli is a simple-to-use command line application and tool library.
//
// Contains: cli app, flags parse, interact, progress, data show tools.
//
// Source code and other details for the project are available at GitHub:
//
//	https://github.com/gookit/gcli
//
// Usage please refer examples and see README
package gcli

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/envutil"
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
	// VerbEnvName for set gcli debug level
	VerbEnvName = "GCLI_VERBOSE"
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

var (
	// DefaultVerb the default Verbose level
	defaultVerb = VerbError
	// Version the gcli version
	version = "3.0.0"
	// CommitID the gcli last commit ID
	commitID = "z20210214"

	// global vars
	gOpts = newGlobalOpts()
	gCtx  = NewCtx().InitCtx()
)

// init
func init() {
	// set Verbose from ENV var.
	if verb := os.Getenv(VerbEnvName); verb != "" {
		_ = gOpts.Verbose.Set(verb)
	}
}

// GCtx get the global ctx
func GCtx() *Context {
	return gCtx
}

// Flags alias of the gflag.Parser
type Flags = gflag.Parser

// FlagMeta alias of the gflag.CliOpt
type FlagMeta = gflag.CliOpt

// CliOpt alias of the gflag.CliOpt
type CliOpt = gflag.CliOpt

// FlagsConfig alias of the gflag.Config
type FlagsConfig = gflag.Config

// NewFlags create new gflag.Flags
func NewFlags(nameWithDesc ...string) *gflag.Flags {
	return gflag.New(nameWithDesc...)
}

// CliArg alias of the gflag.CliArg
type CliArg = gflag.CliArg

// Argument alias of the gflag.CliArg
type Argument = gflag.CliArg

// CliArgs alias of the gflag.CliArgs
type CliArgs = gflag.CliArgs

// Arguments alias of the gflag.CliArgs
type Arguments = gflag.CliArgs

// NewArgument quick create a new command argument
func NewArgument(name, desc string, requiredAndArrayed ...bool) *Argument {
	return gflag.NewArg(name, desc, nil, requiredAndArrayed...)
}

/*************************************************************************
 * global options
 *************************************************************************/

// GlobalOpts global flag options
type GlobalOpts struct {
	Disable  bool
	NoColor  bool
	Verbose  VerbLevel // message report level
	ShowHelp bool
	// ShowVersion show version information
	ShowVersion bool
	// NoProgress dont display progress
	NoProgress bool
	// NoInteractive close interactive confirm
	NoInteractive bool
	// TODO Run application an interactive shell environment
	inShell bool
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

// SetVerbose value
func (g *GlobalOpts) SetVerbose(verbose VerbLevel) {
	g.Verbose = verbose
}

// SetStrictMode option
func (g *GlobalOpts) SetStrictMode(strictMode bool) {
	g.strictMode = strictMode
}

// SetDisable global options
func (g *GlobalOpts) SetDisable() {
	g.Disable = true
}

func (g *GlobalOpts) bindingOpts(fs *gflag.Parser) {
	fs.BoolOpt(&g.ShowHelp, "help", "h", false, "Display the help information")
	fs.AfterParse = func(_ *gflag.Parser) error {
		// return ErrHelp on ShowHelp=true
		if g.ShowHelp {
			return flag.ErrHelp
		}
		return nil
	}

	// disabled
	if g.Disable {
		return
	}

	// up: allow use int and string.
	fs.VarOpt(&g.Verbose, "verbose", "verb", "Set logs reporting level(quiet 0 - 5 crazy)")
	fs.BoolOpt(&g.NoColor, "no-color", "nc", g.NoColor, "Disable color when outputting message")
	fs.BoolOpt(&g.NoProgress, "no-progress", "np", g.NoProgress, "Disable display progress message")
	fs.BoolOpt(&g.ShowVersion, "version", "V", false, "Display app version information")
	fs.BoolOpt(&g.NoInteractive, "no-interactive", "ni", g.NoInteractive, "Disable interactive confirmation operation")
	// fs.BoolOpt(&g.inShell, "ishell", "", false, "Run in an interactive shell environment(`TODO`)")
}

func newGlobalOpts() *GlobalOpts {
	opts := &GlobalOpts{
		strictMode: false,
		// init error level.
		Verbose: defaultVerb,
		NoColor: envutil.GetBool("NO_COLOR", false),
	}

	return opts
}

// GOpts get the global options
func GOpts() *GlobalOpts {
	return gOpts
}

// Config global options
func Config(fn func(opts *GlobalOpts)) {
	if fn != nil {
		fn(gOpts)
	}
}

// ResetGOpts instance
func ResetGOpts() {
	*gOpts = *newGlobalOpts()
}

// Version of the gcli
func Version() string {
	return version
}

// CommitID of the gcli
func CommitID() string { return commitID }

// Verbose returns Verbose level
func Verbose() VerbLevel { return gOpts.Verbose }

// SetVerbose level
func SetVerbose(verbose VerbLevel) { gOpts.SetVerbose(verbose) }

// SetDebugMode level
func SetDebugMode() { SetVerbose(VerbDebug) }

// SetQuietMode level
func SetQuietMode() { SetVerbose(VerbQuiet) }

// ResetVerbose level
func ResetVerbose() { SetVerbose(defaultVerb) }

// StrictMode get is strict mode
func StrictMode() bool { return gOpts.strictMode }

// SetStrictMode for parse flags
func SetStrictMode(strict bool) { gOpts.SetStrictMode(strict) }

// IsGteVerbose get is strict mode
func IsGteVerbose(verb VerbLevel) bool { return gOpts.Verbose >= verb }

// IsDebugMode get is debug mode
func IsDebugMode() bool { return gOpts.Verbose >= VerbDebug }

/*************************************************************************
 * options: some special flag vars
 * - implemented flag.Value interface
 *************************************************************************/

// Ints The int flag list, implemented flag.Value interface
type Ints = cflag.Ints

// Strings The string flag list, implemented flag.Value interface
type Strings = cflag.Strings

// Booleans The bool flag list, implemented flag.Value interface
type Booleans = cflag.Booleans

// EnumString The string flag list, implemented flag.Value interface
type EnumString = cflag.EnumString

// ConfString The config-string flag, INI format, like nginx-config
type ConfString = cflag.ConfString

// String type, a special string
//
// Usage:
//
//	// case 1:
//	var names gcli.String
//	c.VarOpt(&names, "names", "", "multi name by comma split")
//
//	--names "tom,john,joy"
//	 names.Split(",") -> []string{"tom","john","joy"}
//
//	// case 2:
//	var ids gcli.String
//	c.VarOpt(&ids, "ids", "", "multi id by comma split")
//
//	--names "23,34,56"
//	 names.Ints(",") -> []int{23,34,56}
type String = cflag.String

/*************************************************************************
 * Verbose level
 *************************************************************************/

// VerbLevel type.
type VerbLevel uint

// Int Verbose level to int.
func (vl *VerbLevel) Int() int {
	return int(*vl)
}

// String Verbose level to string.
func (vl *VerbLevel) String() string {
	return fmt.Sprintf("%d=%s", *vl, vl.Name())
}

// Upper Verbose level to string.
func (vl *VerbLevel) Upper() string {
	return strings.ToUpper(vl.Name())
}

// Name Verbose level to string.
func (vl *VerbLevel) Name() string {
	switch *vl {
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
		if iv > int(VerbCrazy) {
			*vl = VerbCrazy
		} else if iv < 0 { // fallback to default level.
			*vl = defaultVerb
		} else { // 0 - 5
			*vl = VerbLevel(iv)
		}

		return nil
	}

	// string: level name.
	*vl = name2verbLevel(value)
	return nil
}
