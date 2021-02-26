package gcli

import (
	"strings"

	"github.com/gookit/goutil/maputil"
)

// will inject to every Command
type commandBase struct {
	// Logo ASCII logo setting
	Logo Logo
	// Version app version. like "1.0.1"
	Version string

	// Cmds sub commands of the Command
	// Cmds []*Command
	// mapping sub-command.name => Cmds.index of the Cmds
	// name2idx map[string]int

	// all commands for the group
	commands map[string]*Command
	// command names. key is name, value is name string length
	// eg. {"test": 4, "example": 7}
	cmdNames map[string]int
	// sub command aliases map. {alias: name}
	cmdAliases maputil.Aliases

	// raw input command name
	inputName string
	// current command name
	commandName string
	// the max width for added command names. default set 12.
	nameMaxWidth int
	// the default command name.
	// if is empty:
	// - on app, will render help message.
	// - on cmd, it's sub-command, will run current command
	defaultCommand string

	// Whether it has been initialized
	initialized bool

	// store some runtime errors
	errors []error
	// TODO simple context data map
	data map[string]interface{}
}

func newCommandBase() commandBase {
	return commandBase{
		cmdNames: make(map[string]int),
		// name2idx: make(map[string]int),
		commands: make(map[string]*Command),
		// set an default value.
		nameMaxWidth: 12,
		cmdAliases:   make(maputil.Aliases),
	}
}

// Command get an command by name
func (b commandBase) Command(name string) *Command {
	return b.commands[name]
}

// IsAlias name check
func (b commandBase) IsAlias(alias string) bool {
	return b.cmdAliases.HasAlias(alias)
}

// ResolveAlias get real command name by alias
func (b commandBase) ResolveAlias(alias string) string {
	return b.cmdAliases.ResolveAlias(alias)
}

// SetDefaultCommand set default sub-command name
func (b commandBase) SetDefaultCommand(name string) {
	b.defaultCommand = name
}

// HasCommand name check
func (b commandBase) HasCommand(name string) bool {
	_, has := b.cmdNames[name]
	return has
}

// IsCommand name check. alias of the HasCommand()
func (b commandBase) IsCommand(name string) bool {
	_, has := b.cmdNames[name]
	return has
}

// add Command to the group
func (b commandBase) addCommand(c *Command) {
	// init command
	c.initialize()

	cName := c.Name
	if _, ok := b.cmdNames[cName]; ok {
		panicf("The command name '%s' is already added", cName)
	}

	if b.cmdAliases.HasAlias(cName) {
		panicf("The name '%s' is already used as an alias", cName)
	}

	if c.IsDisabled() {
		Logf(VerbDebug, "command '%s' has been disabled, skip add", cName)
		return
	}

	nameLen := len(cName)

	// add command to app
	b.cmdNames[cName] = nameLen

	// record command name max length
	if nameLen > b.nameMaxWidth {
		b.nameMaxWidth = nameLen
	}

	// add aliases for the command
	b.cmdAliases.AddAliases(c.Name, c.Aliases)
	Logf(VerbDebug, "register a new CLI command: %s", cName)

	// c.app = app
	// inherit global flags from application
	// c.core.gFlags = app.gFlags
	// append
	b.commands[cName] = c
}

// Match command by path names. eg. ["top", "sub"]
func (b commandBase) Match(names []string) *Command {
	ln := len(names)
	if ln == 0 {
		panic("the command names is required")
	}

	top := names[0]
	top = b.ResolveAlias(top)

	c, ok := b.commands[top]
	if !ok {
		return nil
	}

	// sub-sub commands
	if ln > 1 {
		return c.Match(names[1:])
	}

	// current command
	return c
}

// Match command by path. eg. "top:sub" or "top sub"
func (b commandBase) MatchByPath(path string) *Command {
	var names []string
	path = strings.TrimSpace(path)
	if path != "" {
		if strings.ContainsRune(path, ' ') {
			names = strings.Split(path, " ")
		} else {
			names = strings.Split(path, CommandSep)
		}
	}

	return b.Match(names)
}

// GetCommand get command by name. eg "sub"
func (b commandBase) GetCommand(name string) *Command {
	return b.commands[name]
}

// SetLogo text and color style
func (b commandBase) SetLogo(logo string, style ...string) {
	b.Logo.Text = logo
	if len(style) > 0 {
		b.Logo.Style = style[0]
	}
}

// AddError to the application
func (b commandBase) AddError(err error) {
	b.errors = append(b.errors, err)
}

// Commands get all commands
func (b commandBase) Commands() map[string]*Command {
	return b.commands
}

// CmdNames get all command names
func (b commandBase) CmdNames() []string {
	return b.CommandNames()
}

// CommandNames get all command names
func (b commandBase) CommandNames() []string {
	var ss []string
	for n := range b.cmdNames {
		ss = append(ss, n)
	}
	return ss
}

// CmdNameMap get all command names
func (b commandBase) CmdNameMap() map[string]int {
	return b.cmdNames
}

// CmdAliases get all aliases
func (b commandBase) CmdAliases() maputil.Aliases {
	return b.cmdAliases
}
