package gcli

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/structs"
	"github.com/gookit/goutil/strutil"
)

// core definition TODO rename to context ??
type core struct {
	*cmdLine
	// Hooks manage. allowed hooks: "init", "before", "after", "error"
	*Hooks
	// HelpVars help template vars.
	HelpVars
	// global options flag set
	gFlags *Flags
	SimplePrinter
	// GOptsBinder you can custom binding global options
	GOptsBinder func(gf *Flags)
}

// init core
// func (c core) init(cmdName string) {
// 	c.cmdLine = CLI
//
// 	c.AddVars(c.innerHelpVars())
// 	c.AddVars(map[string]string{
// 		"cmd": cmdName,
// 		// binName with command
// 		"binWithCmd": c.binName + " " + cmdName,
// 		// binFile with command
// 		"fullCmd": c.binFile + " " + cmdName,
// 	})
// }

func (c core) doParseGOpts(args []string) (err error) {
	if c.gFlags == nil { // skip on nil
		return
	}

	// parse global options
	err = c.gFlags.Parse(args)
	if err != nil {
		Logf(VerbWarn, "parse global options err: <red>%s</>", err.Error())
	}

	return
}

// GlobalFlags get the app GlobalFlags
func (c core) GlobalFlags() *Flags {
	return c.gFlags
}

// RawOsArgs get the raw os.Args
func (c core) RawOsArgs() []string {
	return os.Args
}

// common basic help vars
func (c core) innerHelpVars() map[string]string {
	return map[string]string{
		"pid":     CLI.PIDString(),
		"workDir": CLI.workDir,
		"binFile": CLI.binFile,
		"binName": CLI.binName,
	}
}

// SimplePrinter struct. for inject struct
type SimplePrinter struct{}

// Print message
func (s SimplePrinter) Print(v ...interface{}) {
	color.Print(v...)
}

// Printf message
func (s SimplePrinter) Printf(format string, v ...interface{}) {
	color.Printf(format, v...)
}

// Println message
func (s SimplePrinter) Println(v ...interface{}) {
	color.Println(v...)
}

// Infoln message
func (s SimplePrinter) Infoln(a ...interface{}) {
	color.Info.Println(a...)
}

// Warnln message
func (s SimplePrinter) Warnln(a ...interface{}) {
	color.Warn.Println(a...)
}

// Errorln message
func (s SimplePrinter) Errorln(a ...interface{}) {
	color.Error.Println(a...)
}

// simple map[string]interface{} struct
type mapData struct {
	data map[string]interface{}
}

// Data get all
func (md *mapData) Data() map[string]interface{} {
	return md.data
}

// SetData set all data
func (md *mapData) SetData(data map[string]interface{}) {
	md.data = data
}

// Value get from data
func (md *mapData) Value(key string) interface{} {
	return md.data[key]
}

// StrValue get from data
func (md *mapData) StrValue(key string) string {
	return strutil.MustString(md.data[key])
}

// IntValue get from data
func (md *mapData) IntValue(key string) int {
	return mathutil.MustInt(md.data[key])
}

// SetValue to data
func (md *mapData) SetValue(key string, val interface{}) {
	if md.data == nil {
		md.data = make(map[string]interface{})
	}
	md.data[key] = val
}

// ClearData all data
func (md *mapData) ClearData() {
	md.data = nil
}

/*************************************************************
 * simple events manage
 *************************************************************/

// HookFunc definition.
// func arguments:
//  in app, like: func(app *App, data ...interface{})
//  in cmd, like: func(cmd *Command, data ...interface{})
// type HookFunc func(obj interface{}, data interface{})
// return:
// - True go on handle. default is True
// - False stop goon handle.
type HookFunc func(data ...interface{}) (stop bool)

// HookCtx struct
type HookCtx struct {
	mapData
	App *App
	Cmd *Command

	name string
	data map[string]interface{}
}

// Name of event
func (hc *HookCtx) Name() string {
	return hc.name
}

// Hooks struct
type Hooks struct {
	// Hooks can setting some hooks func on running.
	hooks map[string]HookFunc
}

// On register event hook by name
func (h *Hooks) On(name string, handler HookFunc) {
	if handler != nil {
		if h.hooks == nil {
			h.hooks = make(map[string]HookFunc)
		}

		h.hooks[name] = handler
	}
}

// AddOn register on not exists hook.
func (h *Hooks) AddOn(name string, handler HookFunc) {
	if _, ok := h.hooks[name]; !ok {
		h.On(name, handler)
	}
}

// Fire event by name, allow with event data
func (h *Hooks) Fire(event string, data ...interface{}) (stop bool) {
	if handler, ok := h.hooks[event]; ok {
		return handler(data...)
	}

	return false
}

// HasHook register
func (h *Hooks) HasHook(event string) bool {
	_, ok := h.hooks[event]
	return ok
}

// ClearHooks clear hooks data
func (h *Hooks) ClearHooks() {
	h.hooks = nil
}

/*************************************************************
 * Command Line: command data
 *************************************************************/

// cmdLine store common data for CLI
type cmdLine struct {
	// pid for current application
	pid int
	// os name.
	osName string
	// the CLI app work dir path. by `os.Getwd()`
	workDir string
	// bin script file, by `os.Args[0]`. eg "./path/to/cliapp"
	binFile string
	// bin script dir path. eg "./path/to"
	binDir string
	// bin script filename. eg "cliapp"
	binName string
	// os.Args to string, but no binName.
	argLine string
}

func newCmdLine() *cmdLine {
	binFile := os.Args[0]
	workDir, _ := os.Getwd()

	// binName will contains work dir path on windows
	// if envutil.IsWin() {
	// 	binFile = strings.Replace(CLI.binName, workDir+"\\", "", 1)
	// }

	return &cmdLine{
		pid: os.Getpid(),
		// more info
		osName:  runtime.GOOS,
		workDir: workDir,
		binDir:  filepath.Dir(binFile),
		binFile: binFile,
		binName: filepath.Base(binFile),
		argLine: strings.Join(os.Args[1:], " "),
	}
}

// PID get pid
func (c *cmdLine) PID() int {
	return c.pid
}

// PIDString get pid as string
func (c *cmdLine) PIDString() string {
	return strconv.Itoa(c.pid)
}

// OsName is equals to `runtime.GOOS`
func (c *cmdLine) OsName() string {
	return c.osName
}

// OsArgs is equals to `os.Args`
func (c *cmdLine) OsArgs() []string {
	return os.Args
}

// BinName get bin script file
func (c *cmdLine) BinFile() string {
	return c.binFile
}

// BinName get bin script name
func (c *cmdLine) BinName() string {
	return c.binName
}

// BinDir get bin script dirname
func (c *cmdLine) BinDir() string {
	return path.Dir(c.binFile)
}

// WorkDir get work dirname
func (c *cmdLine) WorkDir() string {
	return c.workDir
}

// ArgLine os.Args to string, but no binName.
func (c *cmdLine) ArgLine() string {
	return c.argLine
}

func (c *cmdLine) hasHelpKeywords() bool {
	if c.argLine == "" {
		return false
	}

	return strings.HasSuffix(c.argLine, " -h") || strings.HasSuffix(c.argLine, " --help")
}

/*************************************************************
 * app/cmd help vars
 *************************************************************/

// HelpVarFormat allow var replace on render help info.
// Default support:
// 	"{$binName}" "{$cmd}" "{$fullCmd}" "{$workDir}"
const HelpVarFormat = "{$%s}"

// HelpVars struct. provide string var function for render help template.
type HelpVars struct {
	// varLeft, varRight string
	// varFormat string
	// Vars you can add some vars map for render help info
	Vars map[string]string
}

// AddVar get command name
func (hv *HelpVars) AddVar(name, value string) {
	if hv.Vars == nil {
		hv.Vars = make(map[string]string)
	}

	hv.Vars[name] = value
}

// AddVars add multi tpl vars
func (hv *HelpVars) AddVars(vars map[string]string) {
	for n, v := range vars {
		hv.AddVar(n, v)
	}
}

// GetVar get a help var by name
func (hv *HelpVars) GetVar(name string) string {
	return hv.Vars[name]
}

// GetVars get all tpl vars
func (hv *HelpVars) GetVars() map[string]string {
	return hv.Vars
}

// ReplaceVars replace vars in the input string.
func (hv *HelpVars) ReplaceVars(input string) string {
	// if not use var
	if !strings.Contains(input, "{$") {
		return input
	}

	var ss []string
	for n, v := range hv.Vars {
		ss = append(ss, fmt.Sprintf(HelpVarFormat, n), v)
	}

	return strings.NewReplacer(ss...).Replace(input)
}

/*************************************************************
 * command Base
 *************************************************************/

// will inject to every Command
type commandBase struct {
	mapData
	// Logo ASCII logo setting
	Logo *Logo
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
	cmdAliases *structs.Aliases

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
}

func newCommandBase() commandBase {
	return commandBase{
		Logo: &Logo{Style: "info"},
		// init mapping
		cmdNames: make(map[string]int),
		// name2idx: make(map[string]int),
		commands: make(map[string]*Command),
		// set an default value.
		nameMaxWidth: 12,
		// cmdAliases:   make(maputil.Aliases),
		cmdAliases: structs.NewAliases(aliasNameCheck),
	}
}

// GetCommand get an command by name
func (b *commandBase) GetCommand(name string) *Command {
	return b.commands[name]
}

// Command get an command by name
func (b *commandBase) Command(name string) (c *Command, exist bool) {
	c, exist = b.commands[name]
	return
}

// IsAlias name check
func (b *commandBase) IsAlias(alias string) bool {
	return b.cmdAliases.HasAlias(alias)
}

// ResolveAlias get real command name by alias
func (b *commandBase) ResolveAlias(alias string) string {
	return b.cmdAliases.ResolveAlias(alias)
}

// SetDefaultCommand set default sub-command name
func (b *commandBase) SetDefaultCommand(name string) {
	b.defaultCommand = name
}

// HasCommands on the cmd/app
func (b *commandBase) HasCommands() bool {
	return len(b.cmdNames) > 0
}

// HasCommand name check
func (b *commandBase) HasCommand(name string) bool {
	_, has := b.cmdNames[name]
	return has
}

// IsCommand name check. alias of the HasCommand()
func (b *commandBase) IsCommand(name string) bool {
	_, has := b.cmdNames[name]
	return has
}

// add Command to the group
func (b *commandBase) addCommand(pName string, c *Command) {
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
func (b *commandBase) Match(names []string) *Command {
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
func (b *commandBase) FindCommand(path string) *Command {
	return b.Match(splitPath2names(path))
}

// FindByPath command by path. eg. "top:sub" or "top sub"
func (b *commandBase) FindByPath(path string) *Command {
	return b.Match(splitPath2names(path))
}

// MatchByPath command by path. eg. "top:sub" or "top sub"
func (b *commandBase) MatchByPath(path string) *Command {
	return b.Match(splitPath2names(path))
}

// SetLogo text and color style
func (b *commandBase) SetLogo(logo string, style ...string) {
	b.Logo.Text = logo
	if len(style) > 0 {
		b.Logo.Style = style[0]
	}
}

// AddError to the application
func (b *commandBase) AddError(err error) {
	b.errors = append(b.errors, err)
}

// Commands get all commands
func (b *commandBase) Commands() map[string]*Command {
	return b.commands
}

// CmdNames get all command names
func (b *commandBase) CmdNames() []string {
	return b.CommandNames()
}

// CommandNames get all command names
func (b *commandBase) CommandNames() []string {
	var ss []string
	for n := range b.cmdNames {
		ss = append(ss, n)
	}
	return ss
}

// CmdNameMap get all command names
func (b *commandBase) CmdNameMap() map[string]int {
	return b.cmdNames
}

// CmdAliases get cmd aliases
func (b *commandBase) CmdAliases() *structs.Aliases {
	return b.cmdAliases
}

// AliasesMapping get cmd aliases mapping
func (b *commandBase) AliasesMapping() map[string]string {
	return b.cmdAliases.Mapping()
}
