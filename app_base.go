package cliapp

import (
	"strings"
	"fmt"
	"os"
	"flag"
	"github.com/gookit/color"
)

// constants for error level
const (
	VerbQuiet = iota // don't report anything
	VerbError        // reporting on error
	VerbWarn
	VerbInfo
	VerbDebug
)

// HelpVar allow var replace in help info. like "{$binName}" "{$cmd}"
const HelpVar = "{$%s}"

// Logo app logo, ASCII logo
type Logo struct {
	Text  string // ASCII logo string
	Style string // eg "info"
}

// GlobalOpts global flags
type GlobalOpts struct {
	noColor     bool
	showHelp    bool
	showVersion bool
}

// Application the cli app definition
type Application struct {
	// app name
	Name string
	// app version. like "1.0.1"
	Version string
	// app description
	Description string

	// ASCII logo setting
	Logo Logo

	// use strict mode. short opt must be begin '-', long opt must be begin '--'
	Strict bool

	// current command name
	command string

	// default command name
	defaultCmd string

	// vars you can add some vars map for render help info
	vars map[string]string

	// command names. key is name, value is name string length
	// eg. {"test": 4, "example": 7}
	names map[string]int // value

	// command aliases map {alias: name}
	aliases map[string]string
}

// global options
var gOpts = GlobalOpts{}

// bin script name eg "./cliapp"
var binName = os.Args[0]

// the app work dir path
var workDir, _ = os.Getwd()

// error level
var Verbose = VerbError

// Exit
func Exit(code int) {
	os.Exit(code)
}

// WorkDir get work dir
func WorkDir() string {
	return workDir
}

// BinName get bin script name
func BinName() string {
	return binName
}

// LogoText
func (app *Application) LogoText(logo string) {
	app.Logo.Text = logo
}

// SetDebug
func (app *Application) SetDebug() {
	Verbose = VerbDebug
}

// SetVerbose
func (app *Application) SetVerbose(verbose int) {
	Verbose = verbose
}

// DefaultCmd set default command name
func (app *Application) DefaultCmd(name string) {
	app.defaultCmd = name
}

// parseGlobalOpts parse global options
func parseGlobalOpts() []string {
	// Some global options
	flag.Usage = showCommandsHelp

	flag.IntVar(&Verbose, "verbose", VerbError, "")
	flag.BoolVar(&gOpts.showHelp, "h", false, "")
	flag.BoolVar(&gOpts.showHelp, "help", false, "")
	flag.BoolVar(&gOpts.showVersion, "V", false, "")
	flag.BoolVar(&gOpts.showVersion, "version", false, "")
	flag.BoolVar(&gOpts.noColor, "no-color", false, "")

	flag.Parse()

	// fmt.Printf("verb %v, global opts: %+v\n", Verbose, gOpts)

	return flag.Args()
}

// AddVar get command name
func (app *Application) AddVar(name string, value string) {
	app.vars[name] = value
}

// AddVars add multi tpl vars
func (app *Application) AddVars(vars map[string]string) {
	for n, v := range vars {
		app.AddVar(n, v)
	}
}

// GetVar get a help var by name
func (app *Application) GetVar(name string) string {
	if v, ok := app.vars[name]; ok {
		return v
	}

	return ""
}

// GetVars get all tpl vars
func (app *Application) GetVars(name string, value string) map[string]string {
	return app.vars
}

// ReplaceVars replace vars in the help info
func ReplaceVars(help string, vars map[string]string) string {
	// if not use var
	if !strings.Contains(help, "{$") {
		return help
	}

	var ss []string
	for n, v := range vars {
		ss = append(ss, fmt.Sprintf(HelpVar, n), v)
	}

	return strings.NewReplacer(ss...).Replace(help)
}

// strictFormatArgs '-ab' will split to '-a -b', '--o' -> '-o'
func strictFormatArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}

	var fmtdArgs []string

	for _, arg := range args {
		l := len(arg)

		if strings.Index(arg, "--") == 0 {
			if l == 3 {
				arg = "-" + string(arg[2])
			}

		} else if strings.Index(arg, "-") == 0 {
			if l > 2 {
				bools := strings.Split(strings.Trim(arg, "-"), "")

				for _, s := range bools {
					fmtdArgs = append(fmtdArgs, "-"+s)
				}

				continue
			}
		}

		fmtdArgs = append(fmtdArgs, arg)
	}

	return fmtdArgs
}

// Print
func Print(args ...interface{}) (int, error) {
	return color.Print(args...)
}

// Println
func Println(args ...interface{}) (int, error) {
	return color.Println(args...)
}

// Printf
func Printf(format string, args ...interface{}) (int, error) {
	return color.Printf(format, args...)
}
