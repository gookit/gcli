package gcli

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/helper"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/structs"
)

/*************************************************************
 * Command Line: command data
 *************************************************************/

// Context struct
type Context struct {
	maputil.Data
	context.Context
	// PID value
	pid int
	// OsName os name.
	osName string
	// WorkDir the CLI app work dir path. by `os.Getwd()`
	workDir string
	// BinFile bin script file, by `os.Args[0]`. eg "./path/to/cliapp"
	binFile string
	// BinDir bin script dir path. eg "./path/to"
	binDir string
	// BinName bin script filename. eg "cliapp"
	binName string
	// ArgLine os.Args to string, but no binName.
	argLine string
}

// NewCtx instance
func NewCtx() *Context {
	return &Context{
		Data:    make(maputil.Data),
		Context: context.Background(),
	}
}

// Value get by key
func (ctx *Context) Value(key any) any {
	return ctx.Data.Get(key.(string))
}

// InitCtx some common info
func (ctx *Context) InitCtx() *Context {
	binFile := os.Args[0]
	workDir, _ := os.Getwd()

	// binName will contain work dir path on Windows
	// if envutil.IsWin() {
	// 	binFile = strings.Replace(CLI.binName, workDir+"\\", "", 1)
	// }

	ctx.pid = os.Getpid()
	// more info
	ctx.osName = runtime.GOOS
	ctx.workDir = workDir
	ctx.binDir = filepath.Dir(binFile)
	ctx.binFile = binFile
	ctx.binName = filepath.Base(binFile)
	ctx.argLine = strings.Join(os.Args[1:], " ")
	return ctx
}

// PID get pid
func (ctx *Context) PID() int {
	return ctx.pid
}

// PIDString get pid as string
func (ctx *Context) PIDString() string {
	return strconv.Itoa(ctx.pid)
}

// OsName is equals to `runtime.GOOS`
func (ctx *Context) OsName() string {
	return ctx.osName
}

// OsArgs is equals to `os.Args`
func (ctx *Context) OsArgs() []string {
	return os.Args
}

// BinFile get bin script file
func (ctx *Context) BinFile() string {
	return ctx.binFile
}

// BinName get bin script name
func (ctx *Context) BinName() string {
	return ctx.binName
}

// BinDir get bin script dirname
func (ctx *Context) BinDir() string {
	return path.Dir(ctx.binFile)
}

// WorkDir get work dirname
func (ctx *Context) WorkDir() string {
	return ctx.workDir
}

// ArgLine os.Args to string, but no binName.
func (ctx *Context) ArgLine() string {
	return ctx.argLine
}

func (ctx *Context) hasHelpKeywords() bool {
	if ctx.argLine == "" {
		return false
	}
	return strings.HasSuffix(ctx.argLine, " -h") || strings.HasSuffix(ctx.argLine, " --help")
}

// SetValue to ctx
func (ctx *Context) SetValue(key string, val any) {
	ctx.Set(key, val)
}

// GetVal from ctx
func (ctx *Context) GetVal(key string) interface{} {
	return ctx.Get(key)
}

// ResetData from ctx
func (ctx *Context) ResetData() {
	ctx.Data = make(maputil.Data)
}

/*************************************************************
 * command Base
 *************************************************************/

// will inject to every Command
type base struct {
	// Hooks manage. allowed hooks: "init", "before", "after", "error"
	*Hooks
	*Context
	color.SimplePrinter
	// HelpVars help message replace vars.
	helper.HelpVars
	// TODO tplVars for render help template text.
	tplVars map[string]any

	// Logo ASCII logo setting
	Logo *Logo
	// Version app version. like "1.0.1"
	Version string
	// ExitOnEnd call os.Exit on running end
	ExitOnEnd bool
	// ExitFunc default is os.Exit
	ExitFunc func(int)

	// all commands for the group
	commands map[string]*Command
	// command names. key is name, value is name string length
	// eg. {"test": 4, "example": 7}
	cmdNames map[string]int
	// sub command aliases map. {alias: name}
	cmdAliases *structs.Aliases

	// raw input command name
	inputName string
	// current command name
	commandName string
	// the max width for added command names. default set 12.
	nameMaxWidth int
	// has sub-commands on the app
	hasSubcommands bool

	// Whether it has been initialized
	initialized bool
	// store some runtime errors
	errors []error
}

func newBase() base {
	return base{
		Hooks: &Hooks{},
		Logo:  &Logo{Style: "info"},
		// init mapping
		cmdNames: make(map[string]int),
		// name2idx: make(map[string]int),
		commands: make(map[string]*Command),
		// set an default value.
		nameMaxWidth: 12,
		// cmdAliases:   make(maputil.Aliases),
		cmdAliases: structs.NewAliases(aliasNameCheck),
		// ExitOnEnd:  false,
		tplVars: make(map[string]any),
		Context: NewCtx(),
	}
}

// init common basic help vars
func (b *base) initHelpVars() {
	b.AddVars(map[string]string{
		"pid":     b.PIDString(),
		"workDir": b.workDir,
		"binFile": b.binFile,
		"binName": b.binName,
	})
}

// ResetData from ctx
func (b *base) ResetData() {
	if b.Context != nil {
		b.Context.ResetData()
	}
}

// GetCommand get a command by name
func (b *base) GetCommand(name string) *Command {
	return b.commands[name]
}

// Command get a command by name
func (b *base) Command(name string) (c *Command, exist bool) {
	c, exist = b.commands[name]
	return
}

// IsAlias name check
func (b *base) IsAlias(alias string) bool {
	return b.cmdAliases.HasAlias(alias)
}

// ResolveAlias get real command name by alias
func (b *base) ResolveAlias(alias string) string {
	return b.cmdAliases.ResolveAlias(alias)
}

// HasSubcommands on the app
func (b *base) HasSubcommands() bool {
	return b.hasSubcommands
}

// HasCommands on the cmd/app
func (b *base) HasCommands() bool {
	return len(b.cmdNames) > 0
}

// HasCommand name check
func (b *base) HasCommand(name string) bool {
	_, has := b.cmdNames[name]
	return has
}

// IsCommand name check. alias of the HasCommand()
func (b *base) IsCommand(name string) bool {
	_, has := b.cmdNames[name]
	return has
}

// add Command to the group
func (b *base) addCommand(pName string, c *Command) {
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
		Debugf("command '%s' has been disabled, skip add", cName)
		return
	}

	nameLen := len(cName)

	// add command to app
	b.cmdNames[cName] = nameLen
	if c.HasCommands() {
		b.hasSubcommands = true
	}

	// record command name max length
	if nameLen > b.nameMaxWidth {
		b.nameMaxWidth = nameLen
	}

	// add aliases for the command
	Logf(VerbCrazy, "register command '%s'(parent: %s), aliases: %v", cName, pName, c.Aliases)
	b.cmdAliases.AddAliases(c.Name, c.Aliases)

	// c.app = app
	// inherit global flags from application
	// c.core.gFlags = app.gFlags
	// append
	b.commands[cName] = c
}

// Match command by path names. eg. ["top", "sub"]
func (b *base) Match(names []string) *Command {
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

// FindCommand command by path. eg. "top:sub" or "top sub"
func (b *base) FindCommand(path string) *Command {
	return b.Match(splitPath2names(path))
}

// FindByPath command by path. eg. "top:sub" or "top sub"
func (b *base) FindByPath(path string) *Command {
	return b.Match(splitPath2names(path))
}

// MatchByPath command by path. eg: "top:sub" or "top sub"
func (b *base) MatchByPath(path string) *Command {
	return b.Match(splitPath2names(path))
}

// SetLogo text and color style
func (b *base) SetLogo(logo string, style ...string) {
	b.Logo.Text = logo
	if len(style) > 0 {
		b.Logo.Style = style[0]
	}
}

// AddError to the application
func (b *base) AddError(err error) {
	b.errors = append(b.errors, err)
}

// Commands get all commands
func (b *base) Commands() map[string]*Command {
	return b.commands
}

// CmdNames get all command names
func (b *base) CmdNames() []string {
	return b.CommandNames()
}

// CommandNames get all command names
func (b *base) CommandNames() []string {
	var ss []string
	for n := range b.cmdNames {
		ss = append(ss, n)
	}
	return ss
}

// CmdNameMap get all command names
func (b *base) CmdNameMap() map[string]int {
	return b.cmdNames
}

// CmdAliases get cmd aliases
func (b *base) CmdAliases() *structs.Aliases {
	return b.cmdAliases
}

// AliasesMapping get cmd aliases mapping
func (b *base) AliasesMapping() map[string]string {
	return b.cmdAliases.Mapping()
}

// AddTplVar to instance.
func (b *base) AddTplVar(key string, val any) {
	b.tplVars[key] = val
}
