package cliapp

import (
	"strings"
	"fmt"
	"os"
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

// Application the cli app definition
type Application struct {
	// app name
	Name        string
	// app version. like "1.0.1"
	Version     string
	// app description
	Description string

	// ASCII logo setting
	Logo Logo

	// open debug
	Debug bool

	// current command name
	command string

	// default command name
	defaultCmd string

	// vars you can add some vars map for render help info
	vars map[string]string

	// command names {name: 1}
	names map[string]int // value

	// command aliases map {alias: name}
	aliases map[string]string
}

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

// Print
func Print(msg ...interface{}) (int, error) {
	return fmt.Fprint(os.Stdout, msg...)
}

// Println
func Println(msg ...interface{}) (int, error) {
	return fmt.Fprintln(os.Stdout, msg...)
}

// Printf
func Printf(f string, v ...interface{}) (int, error) {
	return fmt.Fprintf(os.Stdout, f+"\n", v...)
}
