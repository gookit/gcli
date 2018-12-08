// Package gcli is a simple to use command line application and tool library.
//
// Contains: cli app, flags parse, interact, progress, data show tools.
//
// Source code and other details for the project are available at GitHub:
// 		https://github.com/gookit/gcli
//
// Usage please refer examples and README
package gcli

import (
	"fmt"
	"github.com/gookit/goutil/envUtil"
	"os"
	"runtime"
	"strings"
)

var (
	// DefaultApp store default application instance
	DefaultApp *App
	// global options
	gOpts = &GlobalOpts{}
	// command auto completion mode.
	// eg "./cli --cmd-completion [COMMAND --OPT ARG]"
	inCompletion bool
	// CLI create a default instance
	CLI = &CmdLine{
		pid: os.Getpid(),
		// more info
		osName:  runtime.GOOS,
		binName: os.Args[0],
		argLine: strings.Join(os.Args[1:], " "),
	}
)

// GlobalOpts global flags
type GlobalOpts struct {
	noColor  bool
	verbose  uint // message report level
	showVer  bool
	showHelp bool
}

// init
func init() {
	workDir, _ := os.Getwd()
	CLI.workDir = workDir

	// binName will contains work dir path on windows
	if envUtil.IsWin() {
		CLI.binName = strings.Replace(CLI.binName, workDir+"\\", "", 1)
	}
}

// NewDefaultApp create the default cli app.
func NewDefaultApp(fn ...func(a *App)) *App {
	DefaultApp = NewApp(fn...)
	return DefaultApp
}

// AllCommands returns all commands in the default app
func AllCommands() map[string]*Command {
	return DefaultApp.Commands()
}

// Exit program
func Exit(code int) {
	os.Exit(code)
}

// Verbose returns verbose level
func Verbose() uint {
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
func SetVerbose(verbose uint) {
	gOpts.verbose = verbose
}

/*************************************************************
 * Command Line: command data
 *************************************************************/

// CmdLine store common data for CLI
type CmdLine struct {
	// pid for current application
	pid int
	// os name.
	osName string
	// the CLI app work dir path. by `os.Getwd()`
	workDir string
	// bin script name, by `os.Args[0]`. eg "./cliapp"
	binName string
	// os.Args to string, but no binName.
	argLine string
}

// PID get PID
func (c *CmdLine) PID() int {
	return c.pid
}

// OsName is equals to `runtime.GOOS`
func (c *CmdLine) OsName() string {
	return c.osName
}

// BinName get bin script name
func (c *CmdLine) BinName() string {
	return c.binName
}

// WorkDir get work dir
func (c *CmdLine) WorkDir() string {
	return c.workDir
}

// ArgLine os.Args to string, but no binName.
func (c *CmdLine) ArgLine() string {
	return c.argLine
}

func (c *CmdLine) helpVars() map[string]string {
	return map[string]string{
		"pid":     fmt.Sprint(CLI.pid),
		"workDir": CLI.workDir,
		"binName": CLI.binName,
	}
}

func (c *CmdLine) hasHelpKeywords() bool {
	return strings.HasSuffix(c.argLine, " -h") || strings.HasSuffix(c.argLine, " --help")
}
