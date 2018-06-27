package cliapp

import (
    "flag"
    "os"
    "log"
	"github.com/golangkit/cliapp/color"
)

// commands collect all command
var names map[string]int // value
var aliases map[string]string
var commands  map[string]*Command
//var commanders  map[string]Commander

// some info
var script = os.Args[0]
var showHelp bool
var showVersion bool

// some constants
const (
    VerbQuiet = iota // don't report anything
    VerbError // reporting on error
    VerbWarn
    VerbInfo
    VerbDebug

    // command type
    //TypeCommand = iota
    //TypeCommander
)

// App the cli app
type App struct {
    Name string
    Version string
    Description string

    // ASCII logo string
    LogoText string

    // open debug
    Debug bool
    // debug level
    Verbose int

    // script name
	script string
	// current command name
	command string
    // work dir path
    workDir string
}

// the app instance
var app *App

func init()  {
    names = make(map[string]int)
    aliases = make(map[string]string)
    commands = make(map[string]*Command)
}

// NewApp create new app
// settings name,version,description
// cli.NewApp("my cli app", "1.0.1", "The is is my cil application")
func NewApp(settings ...string) *App {
    app = &App{Name:"My CLI App", Version: "1.0.0"}
	app.script = script
    app.workDir, _ = os.Getwd()
    app.Verbose = VerbError

    for k, v := range settings{
        switch k {
        case 0:
            app.Name = v
        case 1:
            app.Version = v
        case 2:
            app.Description = v
        }
    }

    return app
}

// SetDebug
func (app *App) SetDebug(debug bool) {
    app.Debug = debug
}

// SetVerbose
func (app *App) SetVerbose(verbose int) {
    app.Verbose = verbose
}

// Run running app
func (app *App) Run() {
    rawName, args := prepareRun()
    name := FindCommandName(rawName)
    logf("input command name is: %s", name)

    if !IsCommand(name) {
        color.New(color.FgRed).Printf("Error: unknown input command '%s'\n", name)
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
func (app *App) SubRun(name string, args []string) int {
	if !IsCommand(name) {
		color.New(color.FgRed).Printf("Error: unknown input command '%s'", name)
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
    flag.BoolVar(&showHelp, "h", false, "")
    flag.BoolVar(&showHelp, "help", false, "")
    flag.BoolVar(&showVersion, "V", false, "")
    flag.BoolVar(&showVersion, "version", false, "")

    flag.Parse()
    // don't display date on print log
    log.SetFlags(0)

    if showHelp {
        showCommandsHelp()
    }

    if showVersion {
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
func (app *App) Add(c *Command) {
	// add ...
    names[c.Name] = 1
    commands[c.Name] = c

    // c.NewFlagSet()
    // will call it on input './cliapp command -h'
    c.Flags.Usage = func() {
        showCommandHelp([]string{c.Name}, true)
    }

    // add alias
    for _, a := range c.Aliases {
        if cmd, has := aliases[a]; has {
            panic(color.Color(color.FgRed).F("the alias '%s' has been used by command '%s'", a, cmd))
        }

        aliases[a] = c.Name
    }
}

// Add add a command
//func (app *App) AddCommander(c Commander) {
//	// run command configure
//	cmd := c.Configure()
//
//	app.Add(cmd)
//}

// WorkDir get work dir
func (app *App) WorkDir() string {
	return app.workDir
}

// Script get script name
func (app *App) Script() string {
	return app.command
}

// Command get command name
func (app *App) Command() string {
	return app.command
}

// Command get command name
func CommandNames() map[string]int {
	return names
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

	log.Printf("[DEBUG] " + f, v...)
}
