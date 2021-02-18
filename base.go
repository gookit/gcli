package gcli

import "github.com/gookit/goutil/maputil"

// router struct definition TODO refactoring
type router struct {

}

// will inject to every Command
type commandBase struct {
	// Logo ASCII logo setting
	Logo Logo
	// Version app version. like "1.0.1"
	Version string

	// Cmds sub commands of the Command
	Cmds []*Command
	// mapping sub-command.name => Cmds.index of the Cmds
	name2idx map[string]int
	// command names. key is name, value is name string length
	// eg. {"test": 4, "example": 7}
	cmdNames map[string]int

	data map[string]interface{} // TODO simple data

	nameMaxWidth int
	defaultCommand string // default sub command name for run.

	// Whether it has been initialized
	initialized bool

	// the max length for added command names. default set 12.
	// the default command name. default is empty, will render help message.
	// all commands for the group
	// commands map[string]*Command
	// sub command aliases map. {alias: name}
	cmdAliases maputil.Aliases

	// store some runtime errors
	errors []error
}

func newCommandBase() commandBase {
	return commandBase{
		cmdNames: make(map[string]int),
		name2idx: make(map[string]int),
		cmdAliases: make(maputil.Aliases),
	}
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
func (b commandBase) addCommand(c *Command)  {
	// validate command name
	cName := c.goodName()
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

	// init command
	// c.app = app
	// inherit global flags from application
	// c.core.gFlags = app.gFlags
	c.initialize()
	// append
	b.Cmds = append(b.Cmds, c)
}

// AddError to the application
func (b commandBase) AddError(err error) {
	b.errors = append(b.errors, err)
}
