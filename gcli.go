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

// PrepareState value. 0=ok, 2=error, -1=goon
type PrepareState int8

// ToInt value
func (ps PrepareState) ToInt() int {
	return int(ps)
}

const (
	// OK success exit code. eg: help command, fired event
	OK PrepareState = 0
	// ERR error exit code
	ERR PrepareState = 2
	// GOON prepare run successful, goon run command
	GOON PrepareState = -1
)

// FoundState for match command name. 0=not found, 1=found
type FoundState int8

const (
	// NotFound not found command by input
	NotFound FoundState = iota
	Founded
	// Invalid // invalid name
)

const (
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
	version = "3.4.0"
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

// FlagMeta alias of the gflag.CliOpt.
// Deprecated: use CliOpt instead
type FlagMeta = gflag.CliOpt

// CliOpt alias of the gflag.CliOpt
type CliOpt = gflag.CliOpt

// FlagsConfig alias of the gflag.Config
type FlagsConfig = gflag.Config

// TagRule* alias of the gflag struct-tag rule type consts. see Flags.FromStruct
const (
	// TagRuleNamed struct tag use named k-v rule. eg: `flag:"name=int0;shorts=i;required=true;desc=message"`
	TagRuleNamed = gflag.TagRuleNamed
	// TagRuleSimple struct tag use simple rule. eg: `flag:"name;desc;required;default;shorts"`
	TagRuleSimple = gflag.TagRuleSimple
	// TagRuleField struct tag use field name as option name, read meta from independent tag keys.
	// eg: `flag:"shorts" desc:"message" default:"val" required:"true"`
	TagRuleField = gflag.TagRuleField
)

// EnhanceShort* alias of the gflag short-option enhance level consts. see Flags Config.EnhanceShort
const (
	// EnhanceShortNone do not enhance short option parse. (default)
	EnhanceShortNone = gflag.EnhanceShortNone
	// EnhanceShortMerge merge bool short option group. eg: `-aux` = `-a -u -x`
	EnhanceShortMerge = gflag.EnhanceShortMerge
	// EnhanceShortAttach also support value-attached short. eg: `-Ostdout` = `-O stdout`
	EnhanceShortAttach = gflag.EnhanceShortAttach
)

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

// GlobalOpts process-level config options. shared across all App instances via
// the package singleton gOpts; set by gcli.SetVerbose / SetStrictMode / SetEnhanceShort.
//
// NOTE: per-app parse/run state (help/version/completion) lives in AppOptions, so
// multiple App instances in one process don't share it. see App.AppOpts().
type GlobalOpts struct {
	// Disable auto binding global options
	Disable bool
	NoColor bool
	// set the message report level.
	//
	// can set by env: GCLI_VERBOSE=debug. see VerbEnvName
	Verbose VerbLevel
	// NoProgress dont display progress. env: NO_PROGRESS
	NoProgress bool
	// NoInteractive close interactive confirm. env: NO_INTERACTIVE
	NoInteractive bool
	// TODO Run application an interactive shell environment
	inShell bool
	// StrictMode use strict mode for parse flags. default: false
	//
	// If True:
	// 	- short opt must be begin "-", long opt must be begin "--"
	//	- will convert like "-ab" to "-a -b"
	// 	- will check invalid arguments, like to many arguments
	strictMode bool
	// enhanceShort global POSIX short-option enhance level, applied to every command
	// that does not set its own Config.EnhanceShort. see EnhanceShortNone/Merge/Attach
	enhanceShort uint8
}

// AppOptions per-app(or standalone command) parse & run state. Each App owns its
// own instance, so concurrent App instances in one process don't share these.
type AppOptions struct {
	// ShowHelp show help information, then exit.
	ShowHelp bool
	// ShowVersion show version information, then exit.
	ShowVersion bool
	// inCompletion dynamic command auto completion mode.
	// eg "./cli --in-completion [COMMAND --OPT ARG]"
	inCompletion bool
	// genCompletion direct generate shell auto completion scripts, then exit.
	// eg "./cli --gen-completion bash|zsh|pwsh"
	genCompletion string
}

// newAppOptions create a new per-app options instance.
func newAppOptions() *AppOptions { return &AppOptions{} }

// SetVerbose value
func (g *GlobalOpts) SetVerbose(verbose VerbLevel) {
	g.Verbose = verbose
}

// SetStrictMode option
func (g *GlobalOpts) SetStrictMode(strictMode bool) {
	g.strictMode = strictMode
}

// SetEnhanceShort global level. see EnhanceShortNone/Merge/Attach
func (g *GlobalOpts) SetEnhanceShort(level uint8) {
	g.enhanceShort = level
}

// SetDisable global options
func (g *GlobalOpts) SetDisable() { g.Disable = true }

// bindingOpts binds the per-app --help/-h (and --version/-V unless globally
// disabled) onto the given parser. g provides the process-level Disable toggle.
func (o *AppOptions) bindingOpts(fs *gflag.Parser, g *GlobalOpts) {
	fs.BoolOpt(&o.ShowHelp, "help", "h", false, "Display the help information")
	fs.AfterParse = func(_ *gflag.Parser) error {
		// return ErrHelp on ShowHelp=true
		if o.ShowHelp {
			return flag.ErrHelp
		}
		return nil
	}

	// disabled
	if g.Disable {
		return
	}

	// NOTE: 不再自动绑定全局 --verbose 选项，避免污染上层应用的选项列表。
	// 日志级别请通过环境变量 GCLI_VERBOSE 控制(见 VerbEnvName)，或调用 gcli.SetVerbose()。
	// fs.BoolOpt(&g.NoColor, "no-color", "nc", g.NoColor, "Disable color when outputting message")
	// fs.BoolOpt(&g.NoProgress, "no-progress", "np", g.NoProgress, "Disable display progress message")
	fs.BoolOpt(&o.ShowVersion, "version", "V", false, "Display app version information")
	// fs.BoolOpt(&g.NoInteractive, "no-interactive", "ni", g.NoInteractive, "Disable interactive confirmation operation")
	// fs.BoolOpt(&g.inShell, "ishell", "", false, "Run in an interactive shell environment(`TODO`)")
}

func newGlobalOpts() *GlobalOpts {
	opts := &GlobalOpts{
		strictMode: false,
		// init error level.
		Verbose: defaultVerb,
		NoColor: envutil.GetBool("NO_COLOR", false),
		// more settings by ENV
		NoProgress:    envutil.GetBool("NO_PROGRESS", false),
		NoInteractive: envutil.GetBool("NO_INTERACTIVE", false),
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

// SetVerbose level by name or level
func SetVerbose[T VerbLevel | string](verbose T) {
	if name, ok := any(verbose).(string); ok {
		gOpts.SetVerbose(VerbLevelFrom(name))
	} else {
		gOpts.SetVerbose(any(verbose).(VerbLevel))
	}
}

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

// EnhanceShort get the global POSIX short-option enhance level.
func EnhanceShort() uint8 { return gOpts.enhanceShort }

// SetEnhanceShort set the global POSIX short-option enhance level for all commands.
//
// level: EnhanceShortNone(0) / EnhanceShortMerge(1) / EnhanceShortAttach(2).
//
// NOTE: a command's own Config.EnhanceShort (if set non-zero) takes priority over this.
func SetEnhanceShort(level uint8) { gOpts.SetEnhanceShort(level) }

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

// String type, a special string
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
	*vl = VerbLevelFrom(value)
	return nil
}
