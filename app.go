package cliapp

import (
    "flag"
    "os"
    "fmt"
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
    app = &App{Name:"My App", Version: "1.0.0"}
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

    name := GetCommandName(rawName)
    logf("input command name is: %s", name)

    if !IsCommand(name) {
        Stdoutf(color.Color(color.FgRed).F("Error: unknown input command '%s'", name))
        showCommandsHelp()
    }

    cmd := commands[name]
    //cmd.Flags.Usage = func() { cmd.ShowHelp() }

    // parse args, don't contains command name.
    if !cmd.CustomFlags {
        cmd.Flags.Parse(args)
        args = cmd.Flags.Args()
    }

    // do execute command
    os.Exit(cmd.Execute(cmd, args))
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

// print debug logging
func logf(f string, v ...interface{}) {
    if !app.Debug {
        return
    }

    log.Printf("[DEBUG] " + f, v...)
}

// Add add a command
func (app *App) Add(c *Command) {
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

// Add add a simple command
//func (app *App) AddSimple(name string, des string, handler func() int) {
//
//}

// Add add a command
//func (app *App) AddCommander(c Commander) {
//    n := c.Name()
//    names[n] = 1
//    commanders[n] = c

    // add alias
    //for _, a := range c.Aliases {
    //    if cmd, has := aliases[a]; has {
    //        panic(fmt.Sprintf("the alias '%s' has been used by command '%s'", a, cmd))
    //    }
	//
    //    aliases[a] = c.Name
    //}
//}

// Run running a sub-command in current command
func (app *App) SubRun(name string, args []string) int {
    if !IsCommand(name) {
        Stdoutf(color.Color(color.FgRed).F("Error: unknown input command '%s'", name))
        return -2
    }

    cmd := commands[name]

    // parse args, don't contains command name.
    if !cmd.CustomFlags {
        cmd.Flags.Parse(args)
        args = cmd.Flags.Args()
    }

    // do execute command
    return cmd.Execute(cmd, args)
}

// GetCommandName get real command name by alias
func GetCommandName(alias string) string {
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

func Stdout(msg ...interface{})  {
    fmt.Fprint(os.Stdout, msg...)
}

func Stdoutf(f string, v ...interface{})  {
    fmt.Fprintf(os.Stdout, f + "\n", v...)
}

func Stderr(msg ...interface{})  {
    fmt.Fprint(os.Stderr, msg...)
}

func Stderrf(f string, v ...interface{})  {
    fmt.Fprintf(os.Stderr, f + "\n", v...)
}
