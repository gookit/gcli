package cliapp

import (
	"github.com/gookit/cliapp/utils"
	"github.com/gookit/color"
	"log"
	"os"
	"strings"
)

// the cli app instance
var app *Application

// commands collect all command
var commands map[string]*Command

// var commanders  map[string]Commander

// init
func init() {
	commands = make(map[string]*Command)

	// binName will contains work dir path on windows
	if utils.IsWin() {
		binName = strings.Replace(binName, workDir+"\\", "", 1)
	}
}

// NewApp create new app
// settings (name, version, description)
// cli.NewApp("my cli app", "1.0.1", "The is is my cil application")
func NewApp(settings ...string) *Application {
	app = &Application{Name: "My CLI Application", Version: "1.0.0"}
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
	app.names = make(map[string]int)

	// init some tpl vars
	app.vars = map[string]string{
		"workDir": workDir,
		"binName": binName,
	}
}

// Add add a command
func (app *Application) Add(c *Command) {
	if c.IsDisabled() {
		return
	}

	app.names[c.Name] = len(c.Name)
	commands[c.Name] = c

	// will call it on input './cliapp command -h'
	c.Flags.Usage = func() {
		c.ShowHelp(true)
	}

	// init
	c.Init()

	// add app vars to cmd
	c.AddVars(app.vars)

	// add aliases for the command
	app.AddAliases(c.Name, c.Aliases)
}

// Run running app
func (app *Application) Run() {
	rawName, args := app.prepareRun()
	name := app.GetNameByAlias(rawName)
	debugf(VerbDebug, "input command name is: %s", name)

	if !app.IsCommand(name) {
		color.Tips("error").Printf("unknown input command '%s'", name)
		// utils.FindSimilar(name, app.CommandNames())

		color.Info("\nPlease use '", binName, " -h' to see available commands\n")
		Exit(-2)
	}

	cmd := commands[name]
	app.command = name

	if app.Strict {
		args = strictFormatArgs(args)
	}

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
	if !app.IsCommand(name) {
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

// prepare to running
func (app *Application) prepareRun() (string, []string) {
	// don't display date on print log
	log.SetFlags(0)

	args := parseGlobalOpts()

	if gOpts.showHelp {
		showCommandsHelp()
	}

	if gOpts.showVersion {
		app.showVersionInfo()
	}

	if gOpts.noColor {
		color.Enable = false
	}

	// check args
	if len(args) < 1 {
		defCmd := app.defaultCmd

		if len(defCmd) == 0 {
			showCommandsHelp()
		}

		if !app.IsCommand(defCmd) {
			color.Tips("error").Printf("The default command '%s' is not exists", defCmd)
			showCommandsHelp()
		}

		args = []string{defCmd}
	}

	// is help command
	if args[0] == "help" {
		// like 'go help'
		if len(args) == 1 {
			showCommandsHelp()
		}

		// like 'go help COMMAND'
		app.showCommandHelp(args[1:], true)
	}

	return args[0], args[1:]
}

// Add add a command
// func (app *Application) AddCommander(c Commander) {
// 	// run command configure
// 	cmd := c.Configure()
//
// 	app.Add(cmd)
// }

// Command get command name
func (app *Application) Command() string {
	return app.command
}

// IsCommand
func (app *Application) IsCommand(name string) bool {
	_, has := app.names[name]

	return has
}

// Command get all command names
func (app *Application) CommandNames() []string {
	var ss []string

	for n := range app.names {
		ss = append(ss, n)
	}

	return ss
}

// AddAliases add alias names for a command
func (app *Application) AddAliases(command string, names []string) {
	if app.aliases == nil {
		app.aliases = make(map[string]string)
	}

	// add alias
	for _, alias := range names {
		if cmd, has := app.aliases[alias]; has {
			panic(color.FgRed.Renderf("The alias '%s' has been used by command '%s'", alias, cmd))
		}

		app.aliases[alias] = command
	}
}

// GetNameByAlias get real command name by alias
func (app *Application) GetNameByAlias(alias string) string {
	if name, has := app.aliases[alias]; has {
		return name
	}

	return alias
}

// Commands get all commands
func (app *Application) Commands() map[string]*Command {
	return commands
}

// AllCommands get all commands
func AllCommands() map[string]*Command {
	return commands
}
