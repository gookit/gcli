package cliapp

import (
	"flag"
	"os"
	"log"
	"github.com/golangkit/cliapp/color"
	"strings"
)

// some constants
const (
	VerbQuiet = iota // don't report anything
	VerbError        // reporting on error
	VerbWarn
	VerbInfo
	VerbDebug

	// command type
	//TypeCommand = iota
	//TypeCommander
)

// GlobalOpts global flags
type GlobalOpts struct {
	showHelp    bool
	showVersion bool
}

// Logo app logo, ASCII logo
type Logo struct {
	Text  string // ASCII logo string
	Style string // eg "info"
}

// Application the cli app definition
type Application struct {
	Name        string
	Version     string
	Description string

	// open debug
	Debug bool
	// debug level
	Verbose int

	// ASCII logo setting
	Logo Logo

	// script name
	script string
	// current command name
	command string
	// work dir path
	workDir string

	// vars you can add some vars map for render help info
	vars map[string]string
}

// the cli app instance
var app *Application

// commands collect all command
var names map[string]int // value
var aliases map[string]string
var commands map[string]*Command
//var commanders  map[string]Commander

// some info
var script = os.Args[0]
var workDir, _ = os.Getwd()
var gOpts = GlobalOpts{}

// init
func init() {
	names = make(map[string]int)
	aliases = make(map[string]string)
	commands = make(map[string]*Command)
}

// NewApp create new app
// settings name,version,description
// cli.NewApp("my cli app", "1.0.1", "The is is my cil application")
func NewApp(settings ...string) *Application {
	app = &Application{Name: "My CLI Application", Version: "1.0.0"}
	app.script = script
	app.workDir = workDir
	app.Verbose = VerbError
	app.Logo.Style = "info"

	for k, v := range settings {
		switch k {
		case 0:
			app.Name = v
		case 1:
			app.Version = v
		case 2:
			app.Description = v
		}
	}

	// init
	app.Init()

	return app
}

// Init
func (app *Application) Init() {
	// init some tpl vars
	app.vars = map[string]string{
		"script":  script,
		"workDir": workDir,
	}
}

// LogoText
func (app *Application) LogoText(logo string) {
	app.Logo.Text = logo
}

// SetDebug
func (app *Application) SetDebug(debug bool) {
	app.Debug = debug
}

// SetVerbose
func (app *Application) SetVerbose(verbose int) {
	app.Verbose = verbose
}

// Add add a command
func (app *Application) Add(c *Command) {
	// add ...
	names[c.Name] = 1
	commands[c.Name] = c

	// will call it on input './cliapp command -h'
	c.Flags.Usage = func() {
		// add app vars to cmd
		c.AddVars(app.vars)
		c.ShowHelp(true)
	}

	// if contains help var "{$cmd}"
	if strings.Contains(c.Description, "{$cmd}") {
		c.Description = strings.Replace(c.Description, "{$cmd}", c.Name, -1)
	}

	// add aliases for the command
	AddAliases(c.Name, c.Aliases)
}

// Run running app
func (app *Application) Run() {
	rawName, args := prepareRun()
	name := FindCommandName(rawName)
	logf("input command name is: %s", name)

	if !IsCommand(name) {
		color.Tips("error").Printf("unknown input command '%s'", name)
		showCommandsHelp()
	}

	cmd := commands[name]
	app.command = name
	//cmd.Flags.Usage = func() { cmd.ShowHelp() }

	// parse args, don't contains command name.
	if !cmd.CustomFlags {
		cmd.Flags.Parse(args)
		args = cmd.Flags.Args()
	}

	// do execute command
	os.Exit(cmd.Execute(app, args))
}

// Run running a sub-command in current command
func (app *Application) SubRun(name string, args []string) int {
	if !IsCommand(name) {
		color.Tips("error").Printf("unknown input command '%s'", name)
		return -2
	}

	cmd := commands[name]

	// parse args, don't contains command name.
	if !cmd.CustomFlags {
		cmd.Flags.Parse(args)
		args = cmd.Flags.Args()
	}

	// do execute command
	return cmd.Execute(app, args)
}

// prepareRun
func prepareRun() (string, []string) {
	flag.Usage = showCommandsHelp

	// some global options
	flag.BoolVar(&gOpts.showHelp, "h", false, "")
	flag.BoolVar(&gOpts.showHelp, "help", false, "")
	flag.BoolVar(&gOpts.showVersion, "V", false, "")
	flag.BoolVar(&gOpts.showVersion, "version", false, "")

	flag.Parse()
	// don't display date on print log
	log.SetFlags(0)

	if gOpts.showHelp {
		showCommandsHelp()
	}

	if gOpts.showVersion {
		showVersionInfo()
	}

	// no command input
	args := flag.Args()
	if len(args) < 1 {
		showCommandsHelp()
	}

	// is help command
	if args[0] == "help" {
		// like 'go help'
		if len(args) == 1 {
			showCommandsHelp()
		}

		// like 'go help COMMAND'
		showCommandHelp(args[1:], true)
	}

	return args[0], args[1:]
}

// Add add a command
//func (app *Application) AddCommander(c Commander) {
//	// run command configure
//	cmd := c.Configure()
//
//	app.Add(cmd)
//}

// WorkDir get work dir
func (app *Application) WorkDir() string {
	return app.workDir
}

// Script get script name
func (app *Application) Script() string {
	return app.command
}

// Command get command name
func (app *Application) Command() string {
	return app.command
}

// Exit
func Exit(code int)  {
	os.Exit(code)
}

// Command get command name
func CommandNames() map[string]int {
	return names
}

// AddAliases add alias names for a command
func AddAliases(command string, names []string)  {
	// add alias
	for _, a := range names {
		if cmd, has := aliases[a]; has {
			panic(color.FgRed.Renderf("The alias '%s' has been used by command '%s'", a, cmd))
		}

		aliases[a] = command
	}
}

// FindCommandName get real command name by alias
func FindCommandName(alias string) string {
	if name, has := aliases[alias]; has {
		return name
	}

	return alias
}

// IsCommand
func IsCommand(name string) bool {
	_, has := names[name]

	return has
}

// print debug logging
func logf(f string, v ...interface{}) {
	if !app.Debug {
		return
	}

	log.Printf("[DEBUG] "+f, v...)
}
