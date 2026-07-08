package gcli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/gevent"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/gcli/v3/internal/helper"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/structs"
	"github.com/gookit/goutil/strutil"
)

// Runner /Executor interface
type Runner interface {
	// Run the command
	//
	// TIP:
	// 	args is the remain arguments after parse flags(options and arguments).
	Run(c *Command, remainArgs []string) error
}

// RunnerFunc definition
//
// TIP:
//
//	args is the remain arguments after parse flags(options and arguments).
type RunnerFunc func(c *Command, remainArgs []string) error

// Run implement the Runner interface
func (f RunnerFunc) Run(c *Command, remainArgs []string) error {
	return f(c, remainArgs)
}

const maxFunc = 64

// HandlersChain middleware handlers chain definition
type HandlersChain []RunnerFunc

// Last returns the last handler in the chain. tip: the last handler is the main own.
func (c HandlersChain) Last() RunnerFunc {
	length := len(c)
	if length > 0 {
		return c[length-1]
	}
	return nil
}

// Command a CLI command structure
type Command struct {
	// internal use
	base

	// Flags cli (options+arguments) parse and manage for the command
	gflag.Flags

	// Name is the command name.
	Name string
	// Desc is the command description message.
	// Can use string-var in contents, eg: {$cmd}
	Desc string

	// Aliases is the command name's alias names
	Aliases arrutil.Strings
	// Category for grouped command display on help
	Category string
	// Config func, will call on `initialize`.
	//
	// - you can config options and other init works
	Config func(c *Command)
	// Hidden the command on render help
	Hidden bool

	// --- for middleware ---
	// run error
	runErr error
	// middleware index number
	middleIdx int8
	// middleware functions
	middles HandlersChain
	// errorHandler // loop find parent.errorHandler

	// path names of the command. 'parent current'
	pathNames []string

	// command is inject to the App
	app  *App
	root bool // is root command

	// Parent parent command
	parent *Command

	// Subs sub commands of the Command
	// NOTICE: if command has been initialized, adding through this field is invalid
	Subs []*Command

	// module is the name for grouped commands
	// subName is the name for grouped commands
	// eg: "sys:info" -> module: "sys", subName: "info"
	// module, subName string

	// Func is the command handler func. Func Runner
	//
	// TIP:
	// 	params: `args` is the remain arguments after parse flags(options and arguments).
	Func RunnerFunc

	// Examples some usage example display.
	//
	// Can use string-var in contents, eg:
	//   {$cmd}, {$binName}, {$binDir}, {$workDir}, {$binWithCmd}, {$binWithPath}, {$fullCmd}
	Examples string
	// Help is the long help message text
	//
	// Can use string-var in contents, eg:
	//   {$cmd}, {$binName}, {$binDir}, {$workDir}, {$binWithCmd}, {$binWithPath}, {$fullCmd}
	Help string
	// HelpRender custom render cmd help message
	HelpRender func(c *Command)

	// mark is disabled. if true will skip register to app.
	disabled bool
	// command is standalone running.
	standalone bool
	// global option binding on standalone. deny error on repeat run.
	gOptBounded bool
	// runOpts is the command's own parse state (help/version) when run standalone.
	runOpts *AppOptions

	// sharedFs holds the command's shared options (≈ cobra PersistentFlags).
	// 共享选项的定义来源: 本命令及其所有子孙命令都会继承这些选项。lazy 创建于 SharedOpts()。
	sharedFs *gflag.Flags
	// sharedMerged marks shared options(self + ancestors) have been merged into c.Flags.
	// 幂等标记, 保证分发时只合并一次共享选项。
	sharedMerged bool
	// localOptNames snapshots the command's own local option names before shared merge.
	// 合并共享选项前的本地选项名快照, 用于区分局部定义与继承副本(保留局部 Required、跳过被局部覆盖的共享必填校验)。
	localOptNames map[string]bool
}

// NewCommand create a new command instance.
//
// Usage:
//
//	cmd := NewCommand("my-cmd", "description")
//	// OR with a config func
//	cmd := NewCommand("my-cmd", "description", func(c *Command) { ... })
//	app.Add(cmd) // OR cmd.AttachTo(app)
func NewCommand(name, desc string, setFn ...func(c *Command)) *Command {
	c := &Command{
		Name: name,
		Desc: desc,
	}

	// init set name
	c.Flags.SetName(name)

	// has config func
	if len(setFn) > 0 {
		c.Config = setFn[0]
	}
	return c
}

// Init command. only use for tests
func (c *Command) Init() { c.initialize() }

// SetFunc Settings command handler func
func (c *Command) SetFunc(fn RunnerFunc) { c.Func = fn }

// WithFunc Settings command handler func
func (c *Command) WithFunc(fn RunnerFunc) *Command {
	c.Func = fn
	return c
}

// WithHidden Settings command is hidden
func (c *Command) WithHidden() *Command {
	c.Hidden = true
	return c
}

// Use 注册一个或多个中间件, 按注册顺序在命令主函数前依次执行; 返回 c 以便链式调用。
func (c *Command) Use(handlers ...RunnerFunc) *Command {
	c.middles = append(c.middles, handlers...)
	return c
}

// AttachTo attach the command to CLI application
func (c *Command) AttachTo(app *App) { app.AddCommand(c) }

// InheritedOptsCategory 是继承(共享)选项在命令 help 中的分组标题。
const InheritedOptsCategory = "Inherited Options"

// appOpts returns the command's own parse state (used in standalone mode), lazy.
func (c *Command) appOpts() *AppOptions {
	if c.runOpts == nil {
		c.runOpts = newAppOptions()
	}
	return c.runOpts
}

// SharedOpts 返回命令专属的共享选项持有器(惰性创建), 对标 cobra 的 PersistentFlags()。
//
// 在它上面像普通选项一样绑定(BoolOpt/StrOpt/Opt[T]/FromStruct/...), 这些选项会被本命令
// 及其所有子孙命令继承: 父命令定义、子命令也能解析, 且父子读写同一个变量(共享 flag.Value)。
//
// 注意: sharedFs 仅作为「定义来源」, 自身永不单独 Parse; 分发时由 parseOptions 合并进 c.Flags。
func (c *Command) SharedOpts() *gflag.Flags {
	if c.sharedFs == nil {
		c.sharedFs = gflag.New(c.Name)
	}
	return c.sharedFs
}

// Disable set cmd is disabled
func (c *Command) Disable() { c.disabled = true }

// Visible return cmd is visible
func (c *Command) Visible() bool { return c.Hidden == false }

// IsDisabled get cmd is disabled
func (c *Command) IsDisabled() bool { return c.disabled }

// IsRunnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as import path.
func (c *Command) IsRunnable() bool { return c.Func != nil }

// Add one or multi sub-command(s). alias of the AddSubs
func (c *Command) Add(sub *Command, more ...*Command) { c.AddSubs(sub, more...) }

// AddSubs add one or multi sub-command(s)
func (c *Command) AddSubs(sub *Command, more ...*Command) {
	c.AddCommand(sub)

	if len(more) > 0 {
		for _, cmd := range more {
			c.AddCommand(cmd)
		}
	}
}

// AddCommand add a sub command
func (c *Command) AddCommand(sub *Command) {
	// init command
	sub.app = c.app
	sub.parent = c
	// inherit something from parent command
	sub.Ctx = c.Ctx
	sub.standalone = c.standalone

	// initialize command
	c.initialize()

	// extend path names from parent
	sub.pathNames = c.pathNames[0:]
	// update some parser config before initializing sub command
	sub.Flags.WithConfigFn(gflag.WithIndentLongOpt(c.ParserCfg().IndentLongOpt))

	// do add and init sub command
	c.base.addCommand(c.Name, sub)
}

// Match sub command by input names
func (c *Command) Match(names []string) *Command {
	// ensure is initialized
	c.initialize()

	if len(names) == 0 { // return self.
		return c
	}
	return c.base.Match(names)
}

// MatchByPath command by path. eg: "top:sub"
func (c *Command) MatchByPath(path string) *Command {
	return c.Match(splitPath2names(path))
}

// initialize works for the command
//
// - ctx
// - sub-cmd
func (c *Command) initialize() {
	if c.initialized {
		return
	}

	// check command name
	cName := c.goodName()
	Debugf("initialize the command '%s': init flags, run config func", cName)

	c.initialized = true
	c.pathNames = append(c.pathNames, cName)

	// init base
	c.initCommandBase(cName)
	c.Fire(gevent.OnCmdInitBefore, nil)

	// init for cmd flags parser.
	// 用命令全路径(如 "git branch")作为 flags 名: flags 名仅用于报错信息(选项/参数
	// 重复绑定 panic 等), 带上全路径便于定位是哪个命令。此刻 pathNames 已包含本命令段。
	c.Flags.Init(c.Path())

	// load common sub commands
	if len(c.Subs) > 0 {
		for _, sub := range c.Subs {
			c.AddCommand(sub)
		}
	}

	// format description
	if len(c.Desc) > 0 {
		c.Desc = strutil.UpperFirst(c.Desc)
		if strings.Contains(c.Desc, "{$") {
			c.Desc = strings.ReplaceAll(c.Desc, "{$cmd}", c.Name)
		}
	}

	// call config func
	if c.Config != nil {
		c.Config(c)
	}

	c.Fire(gevent.OnCmdInitAfter, nil)
}

// init base, ctx
func (c *Command) initCommandBase(cName string) {
	Logf(VerbCrazy, "init command c.base for the command: %s", cName)

	if c.Hooks == nil {
		c.Hooks = &Hooks{}
	}

	if c.Ctx == nil {
		Logf(VerbDebug, "cmd: %s - use the gCtx as command context", cName)
		c.Ctx = gCtx
	}

	binWithPath := c.Ctx.binName + " " + c.Path()

	c.initHelpReplacer()
	c.AddReplaces(map[string]string{
		"cmd": cName,
		// binName with command name
		"binWithCmd": binWithPath,
		// binName with command path
		"binWithPath": binWithPath,
		// binFile with command
		"fullCmd": binWithPath,
	})

	c.base.cmdNames = make(map[string]int)
	c.base.commands = make(map[string]*Command)
	// set an default value.
	c.base.nameMaxWidth = 12
	// c.base.cmdAliases = make(maputil.Aliases)
	c.base.cmdAliases = structs.NewAliases(aliasNameCheck)
}

// Next TODO processing, run all middleware handlers
func (c *Command) Next() {
	c.middleIdx++
	s := int8(len(c.middles))

	for ; c.middleIdx < s; c.middleIdx++ {
		err := c.middles[c.middleIdx](c, c.RawArgs())
		// will abort on error
		if err != nil {
			c.runErr = err
			return
		}
	}
}

/*************************************************************
 * region standalone running
 *************************************************************/

// MustRun Alone the current command, will output message on error
//
// Usage:
//
//	// run with os.Args
//	cmd.MustRun(nil)
//	cmd.MustRun(os.Args[1:])
//	// custom args
//	cmd.MustRun([]string{"-a", ...})
func (c *Command) MustRun(args []string) {
	if err := c.Run(args); err != nil {
		color.Errorln("ERROR:", err.Error())
	}
}

// Run standalone running the command
//
// Usage:
//
//	// run with os.Args
//	cmd.Run(nil)
//	cmd.Run(os.Args[1:])
//	// custom args
//	cmd.Run([]string{"-a", ...})
func (c *Command) Run(args []string) (err error) {
	if c.app != nil || c.parent != nil {
		return c.innerDispatch(args)
	}

	// mark is standalone
	c.standalone = true

	// if not set input args
	if args == nil {
		args = os.Args[1:]
	}

	// init the command
	c.initialize()

	// add default error handler.
	if !c.HasHook(gevent.OnCmdRunError) {
		c.On(gevent.OnCmdRunError, defaultErrHandler)
	}

	// binding global options
	if !c.gOptBounded {
		Debugf("cmd: %s - binding global options on standalone mode", c.Name)
		c.appOpts().bindingOpts(&c.Flags, gOpts)
		c.gOptBounded = true
	}

	// dispatch and parse flags and execute command
	return c.innerDispatch(args)
}

/*************************************************************
 * command run
 *************************************************************/

// dispatch execute the command
func (c *Command) innerDispatch(args []string) (err error) {
	// parse command flags
	args, err = c.parseOptions(args)
	if err != nil {
		if err == flag.ErrHelp {
			Debugf("cmd: %s - parse opts return flag.ErrHelp, render command help", c.Name)
			return c.ShowHelp()
		}

		Debugf("cmd: %s - command options parse error", c.Name)
		color.Error.Tips("option error - %s", err.Error())
		return nil
	}

	// remaining args
	if c.standalone {
		if c.appOpts().ShowHelp {
			Debugf("cmd: %s - ShowHelp is True, render command help", c.Name)
			return c.ShowHelp()
		}

		c.Fire(gevent.OnGlobalOptsParsed, map[string]any{"args": args})
	}

	c.Fire(gevent.OnCmdOptParsed, map[string]any{"args": args})
	Debugf("cmd: %s - remaining args on options parsed: %v", c.Name, args)

	// find sub command
	if len(args) > 0 {
		name := args[0]

		// ensure is not an option
		if name != "" && name[0] != '-' {
			name = c.ResolveAlias(name)

			// is valid sub command
			if sub, has := c.Command(name); has {
				// TIP: loop find sub...command and run it.
				return sub.innerDispatch(args[1:])
			}

			// is not a sub command and has no arguments -> error
			if !c.HasArguments() {
				// fire events
				hookData := map[string]any{"name": name, "args": args[1:]}
				if c.Fire(gevent.OnCmdSubNotFound, hookData) {
					return
				}
				if c.Fire(gevent.OnCmdNotFound, hookData) {
					return
				}

				color.Error.Tips("%s - subcommand '%s' is not found", c.Name, name)
				return newRunErr(ERR.ToInt(), c.NewErrf("%s - subcommand %q is not found", c.Name, name))
			}
		}
	}

	// not set command func and has sub commands.
	if c.Func == nil && len(c.commands) > 0 {
		Logf(VerbWarn, "cmd: %s - c.Func is empty, but has subcommands, render help", c.Name)
		return c.ShowHelp()
	}

	// do execute current command
	return c.doExecute(args)
}

// execute the current command
func (c *Command) innerExecute(args []string, igrErrHelp bool) (err error) {
	// parse flags
	args, err = c.parseOptions(args)
	if err != nil {
		// whether ignore flag.ErrHelp error
		if igrErrHelp && err == flag.ErrHelp {
			err = nil
		}
		return
	}

	// do execute command
	return c.doExecute(args)
}

// do parse option flags, remaining is cmd args
func (c *Command) parseOptions(args []string) (ss []string, err error) {
	// apply global EnhanceShort to commands that don't set their own (command-level wins)
	if c.ParserCfg().EnhanceShort == EnhanceShortNone && gOpts.enhanceShort > EnhanceShortNone {
		c.ParserCfg().EnhanceShort = gOpts.enhanceShort
	}

	// strict format options
	if gOpts.strictMode && len(args) > 0 {
		args = strictFormatArgs(args) // 长选项形态规范化(--a/---name)
		// 短选项安全拆分交由 gflag EnhanceShort(仅全 bool 才拆，修复盲拆误伤取值短选项)
		if c.ParserCfg().EnhanceShort == EnhanceShortNone {
			c.ParserCfg().EnhanceShort = EnhanceShortMerge
		}
	}

	// args reorder(默认开启)在子命令边界处停止, 确保多级命令时只重排最终执行命令的 args。
	c.Flags.SetReorderStop(func(name string) bool {
		if c.isReorderStopName(name) {
			return true
		}
		return len(c.commands) > 0 && !c.HasArguments()
	})

	// 合并共享选项: 沿祖先链(含自身)从根到叶把共享选项并入 c.Flags, 使其在本命令段可解析。
	// 幂等: sharedMerged 保证只合并一次; 合并后写在叶子段任意位置(配合 reorder)也能被识别。
	c.mergeSharedOpts()

	Debugf("cmd: %s - will parse options from args: %v", c.Name, args)

	// parse options, don't contains command name.
	if err = c.Parse(args); err != nil {
		Logf(VerbCrazy, "cmd: %s - parse options, err: <red>%s</>", c.Name, err.Error())
		return
	}

	// remaining args, next use for parse arguments
	return c.RawArgs(), nil
}

// mergeSharedOpts 沿祖先链(含自身)从根到叶把各级的共享选项合并进 c.Flags。
//
// 顺序为 [root,...,parent,self], 这样祖先(更靠近根)的共享选项先注册; 同名时 InheritOptsFrom
// 内部以「已存在则跳过」保证局部/更近的定义优先。幂等: sharedMerged 仅合并一次。
//
// 关于 Required: 共享选项可能写在叶子命令段, 而本命令(若是中间祖先)解析时尚未见到该值,
// 因此合并进 c.Flags 的副本会清除 Required, 避免在「会继续向子命令分发的命令」上误报必填。
// 真正的 Required 校验延后到实际执行命令时由 validateSharedRequired 统一处理。
func (c *Command) mergeSharedOpts() {
	if c.sharedMerged {
		return
	}
	c.sharedMerged = true

	// 沿 parent 链收集 [self, parent, ..., root]
	var chain []*Command
	for cur := c; cur != nil; cur = cur.parent {
		chain = append(chain, cur)
	}

	// 合并前快照: c.Flags 此刻已有的选项均为命令自身的局部定义, 它们的 Required 应保留。
	c.localOptNames = make(map[string]bool, len(c.Flags.Opts()))
	for name := range c.Flags.Opts() {
		c.localOptNames[name] = true
	}

	// 反向遍历(从根到叶)合并, 使祖先共享选项先注册、局部/更近定义优先。
	// help 分组: 祖先继承来的归入 InheritedOptsCategory; 自身定义的共享选项视为本命令选项,
	// 与本地选项同组(不标 Inherited), 对标 cobra(命令自身的 persistent flag 显示在本命令 Flags)。
	for i := len(chain) - 1; i >= 0; i-- {
		anc := chain[i]
		if anc.sharedFs == nil {
			continue
		}
		if anc == c {
			c.Flags.InheritOptsFrom(anc.sharedFs)
		} else {
			c.Flags.InheritOptsFrom(anc.sharedFs, InheritedOptsCategory)
		}
	}

	// 清除「新继承进来」的共享选项副本的 Required: 共享值可能写在叶子命令段, 中间命令解析时
	// 还看不到取值, 不能在那时误报必填。共享选项的 Required 延后到执行命令时由
	// validateSharedRequired 统一校验。局部已定义的同名选项(快照中)保持其 Required 不变。
	for name, opt := range c.Flags.Opts() {
		if !c.localOptNames[name] {
			opt.Required = false
		}
	}
}

// validateSharedRequired 在实际执行命令前, 沿祖先链(含自身)校验所有 Required 共享选项是否已赋值。
//
// 共享选项的 Required 校验延后到此处统一处理: 因为共享选项可写在叶子命令段, 中间命令解析时
// 还看不到取值, 不能在那时报必填。到达执行命令时, 共享值已写回同一 ptr, 此处即可正确判定。
// 若执行命令定义了同名局部选项(共享被跳过继承), 则由局部自身的校验负责, 这里跳过。
func (c *Command) validateSharedRequired() error {
	for cur := c; cur != nil; cur = cur.parent {
		if cur.sharedFs == nil {
			continue
		}
		for name, opt := range cur.sharedFs.Opts() {
			if !opt.Required {
				continue
			}
			// 执行命令有同名局部选项 → 共享被局部覆盖, 由局部 validateAll 负责, 此处跳过
			if c.localOptNames[name] {
				continue
			}

			// 类型感知判空(int 的空是 "0"、float 是 "0.0"、string 是 "")
			if opt.IsEmpty() {
				return fmt.Errorf("option '%s' is required", name)
			}
		}
	}
	return nil
}

// prepare: before execute the command
func (c *Command) prepare(_ []string) (status int, err error) {
	return
}

type panicErr struct {
	val any
}

// Error string
func (p panicErr) Error() string { return fmt.Sprint(p.val) }

// do execute the command
func (c *Command) doExecute(args []string) (err error) {
	// 共享 Required 选项的延后校验: 到达实际执行命令时统一检查祖先链(含自身)的必填共享选项
	if err = c.validateSharedRequired(); err != nil {
		c.Fire(gevent.OnCmdRunError, map[string]any{"cmd": c.Name, "err": err})
		Logf(VerbError, "command '%s' shared required option err: <red>%s</>", c.Name, err.Error())
		return err
	}

	// collect and binding named argument
	Debugf("cmd: %s - collect and binding named arguments", c.Name)
	if err := c.ParseArgs(args); err != nil {
		c.Fire(gevent.OnCmdRunError, map[string]any{"cmd": c.Name, "err": err})
		Logf(VerbError, "binding command '%s' arguments err: <red>%s</>", c.Name, err.Error())
		return err
	}

	fnArgs := c.ExtraArgs()
	c.Fire(gevent.OnCmdRunBefore, map[string]any{"args": fnArgs})

	// do call command handler func
	if c.Func == nil {
		Logf(VerbWarn, "the command '%s' no handler func to running", c.Name)
		c.Fire(gevent.OnCmdRunAfter, nil)
		return
	}

	Debugf("cmd: %s - run command func with extra-args %v", c.Name, fnArgs)

	// recover panics from middleware/command func, convert to error.
	// NOTE: recover 必须在 defer 内调用才有效; 仅在确实 panic 时触发 fireAfterExec,
	// 正常路径下方已 fireAfterExec, 不会重复。
	defer func() {
		if re := recover(); re != nil {
			var ok bool
			err, ok = re.(error)
			if !ok {
				err = panicErr{val: re}
			}
			c.fireAfterExec(err)
		}
	}()

	// do run middlewares and command func
	err = c.runWithMiddles(fnArgs)
	c.fireAfterExec(err)
	return
}

// runWithMiddles 依次执行中间件(任一返回 error 即中止, 后续中间件与主函数都不再执行),
// 全部通过后执行命令主函数。未注册中间件时等价于直接调用 c.Func。
func (c *Command) runWithMiddles(fnArgs []string) error {
	// app 级中间件先于命令级执行(独立命令运行时 c.app 为 nil, 跳过)
	if c.app != nil {
		for _, mw := range c.app.middles {
			if err := mw(c, fnArgs); err != nil {
				return err
			}
		}
	}
	for _, mw := range c.middles {
		if err := mw(c, fnArgs); err != nil {
			return err
		}
	}
	return c.Func(c, fnArgs)
}

func (c *Command) fireAfterExec(err error) {
	if err != nil {
		c.Fire(gevent.OnCmdRunError, map[string]any{"err": err})
	} else {
		c.Fire(gevent.OnCmdRunAfter, nil)
	}
}

/*************************************************************
 * parent and subs
 *************************************************************/

// Root get root command
func (c *Command) Root() *Command {
	if c.parent != nil {
		return c.parent.Root()
	}
	return c
}

// IsRoot command
func (c *Command) IsRoot() bool { return c.parent == nil }

// Parent get parent
func (c *Command) Parent() *Command { return c.parent }

// SetParent set parent
func (c *Command) SetParent(parent *Command) { c.parent = parent }

// ParentName name of the parent command
func (c *Command) ParentName() string {
	if c.parent != nil {
		return c.parent.Name
	}
	return ""
}

// Sub get sub command by name. eg "sub"
func (c *Command) Sub(name string) *Command { return c.GetCommand(name) }

// SubCommand get sub command by name. eg "sub"
func (c *Command) SubCommand(name string) *Command { return c.GetCommand(name) }

// IsSubCommand name check. alias of the HasCommand()
func (c *Command) IsSubCommand(name string) bool { return c.IsCommand(name) }

// find sub command by name
// func (c *Command) findSub(name string) *Command {
// 	if index, ok := c.subName2index[name]; ok {
// 		return c.Subs[index]
// 	}
//
// 	return nil
// }

/*************************************************************
 * helper methods
 *************************************************************/

// IsStandalone running
func (c *Command) IsStandalone() bool { return c.standalone }

// NotStandalone running
func (c *Command) NotStandalone() bool { return !c.standalone }

// ID get command ID name.
func (c *Command) goodName() string {
	name := strings.Trim(strings.TrimSpace(c.Name), ": ")
	if name == "" {
		panicf("the command name can not be empty")
	}

	if !helper.IsGoodCmdName(name) {
		panicf("the command name '%s' is invalid, must match: %s", name, helper.RegGoodCmdName)
	}

	// update name
	c.Name = name
	return name
}

// Fire event handler by name
func (c *Command) Fire(event string, data map[string]any) (stop bool) {
	hookCtx := newHookCtx(event, c, data)
	Debugf("cmd: %s - trigger the event: <mga>%s</>", c.Name, event)

	// notify all parent commands
	p := c.parent
	for p != nil {
		if p.Hooks.Fire(event, hookCtx) {
			return true
		}
		p = p.parent
	}

	// notify to app
	if c.app != nil && c.app.Hooks.Fire(event, hookCtx) {
		return true
	}

	return c.Hooks.Fire(event, hookCtx)
}

// On add hook handler for a hook event
func (c *Command) On(name string, handler HookFunc) {
	Debugf("cmd: %s - register hook: <cyan>%s</>", c.Name, name)

	if c.Hooks == nil {
		c.Hooks = &Hooks{}
	}
	c.Hooks.On(name, handler)
}

// Copy a new command for current
func (c *Command) Copy() *Command {
	nc := *c
	// reset some fields
	nc.Func = nil
	// 副本需独立的 Hooks：base 内嵌的是 *Hooks 指针，浅拷贝会与原命令共享，
	// 直接 ResetHooks() 会清空原命令的钩子，这里改为赋予一个全新的空 Hooks。
	nc.Hooks = &Hooks{}

	return &nc
}

// App returns the CLI application
func (c *Command) App() *App { return c.app }

// ID get command ID string. return like "git:branch:create"
func (c *Command) ID() string { return strings.Join(c.pathNames, CommandSep) }

// Path get command full path, joined by space. eg: "git branch create"
func (c *Command) Path() string { return strings.Join(c.pathNames, " ") }

// PathNames get command path names
func (c *Command) PathNames() []string { return c.pathNames }

// NewErr format message and add error to the command
func (c *Command) NewErr(msg string) error { return errors.New(msg) }

// NewErrf format message and add error to the command
func (c *Command) NewErrf(format string, v ...any) error {
	return fmt.Errorf(format, v...)
}

// HelpDesc format desc string for render help
func (c *Command) HelpDesc() (desc string) {
	if len(c.Desc) == 0 {
		return
	}

	// dump.P(desc)
	desc = strutil.UpperFirst(c.Desc)
	// contains help var "{$cmd}". replace on here is for 'app help'
	if strings.Contains(desc, "{$") {
		desc = strings.ReplaceAll(desc, "{$cmd}", color.WrapTag(c.Name, "mga"))
	}

	return wrapColor2string(desc)
}

// Logf print log message
// func (c *Command) Logf(level uint, format string, v ...any) {
// 	Logf(level, format, v...)
// }
