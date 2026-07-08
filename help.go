package gcli

import (
	"fmt"
	"sort"
	"strings"
	"text/template"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/gevent"
	"github.com/gookit/gcli/v3/internal/helper"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/sysutil"
)

// HelpConfig struct
type HelpConfig struct {
	// AfterCmdText add text after commands list
	AfterCmdText string
	// FooterText add help footer text on help end
	FooterText string
}

// CmdGroup is a group of commands by category. used for render help.
type CmdGroup struct {
	// Name is the category name. "" means the default group.
	Name string
	// Title is the display title for help. eg: "Available Commands"
	Title string
	// Cmds are the visible commands of this group, sorted by name.
	Cmds []*Command
}

// CommandsByGroup group all visible commands by their Category field.
//
// - groups keep the category insertion order(see base.cmdCategories).
// - the default group(empty Category) uses defaultTitle as its title.
// - returns nil when there is no visible command.
func (b *base) CommandsByGroup(defaultTitle string) []*CmdGroup {
	groups := make([]*CmdGroup, 0, len(b.cmdCategories))

	for _, cat := range b.cmdCategories {
		var cmds []*Command
		for _, c := range b.commands {
			if c.Category == cat && c.Visible() {
				cmds = append(cmds, c)
			}
		}
		if len(cmds) == 0 {
			continue
		}

		sort.Slice(cmds, func(i, j int) bool {
			return cmds[i].Name < cmds[j].Name
		})

		title := defaultTitle
		if cat != "" {
			title = strutil.UpperFirst(cat)
		}
		groups = append(groups, &CmdGroup{Name: cat, Title: title, Cmds: cmds})
	}

	return groups
}

/*************************************************************
 * display app help
 *************************************************************/

// display app version info
func (app *App) showVersionInfo() bool {
	Debugf("print application version info")

	// custom color tag, direct print by color
	if strings.Contains(app.Version, "</>") {
		color.Printf("Version: %s\n", app.Version)
	} else {
		color.Printf("Version: <cyan>%s</>\n", app.Version)
	}

	if app.Logo.Text != "" {
		color.Printf("%s\n", color.WrapTag(app.Logo.Text, app.Logo.Style))
	}
	return false
}

// display unknown input command and similar commands tips
func (app *App) showCommandTips(name string) {
	Debugf("will find and show similar command tips")

	color.Error.Tips(`unknown input command "<mga>%s</>"`, name)
	if ns := app.findSimilarCmd(name); len(ns) > 0 {
		color.Printf("\nMaybe you mean:\n  <green>%s</>\n", strings.Join(ns, ", "))
	}

	color.Printf("\nUse <cyan>%s --help</> to see available commands\n", app.Ctx.binName)
}

// AppHelpTemplate help template for app(all commands)
var AppHelpTemplate = `{{.Desc}} (Version: <info>{{.Version}}</>)
<comment>Usage:</>
  {$binName} [global options...] <info>COMMAND</> [--options ...] [arguments ...]{{if .HasSubs }}
  {$binName} [global options...] <info>COMMAND</> [--options ...] <info>SUBCOMMAND</> [--options ...] [arguments ...]{{end}}

<comment>Global Options:</>
{{.GOpts}}
{{range $g := .CmdGroups}}<comment>{{$g.Title}}:</>{{range $c := $g.Cmds}}
  <info>{{$c.Name | paddingName }}</> {{$c.HelpDesc}}{{if $c.Aliases}} (alias: <green>{{ join $c.Aliases ","}}</>){{end}}{{end}}
{{end}}  <info>{{ paddingName "help" }}</> Display help information

{{.Help.AfterCmdText}}Use "<cyan>{$binName} COMMAND -h</>" for more information about a command.{{.Help.FooterText}}
`

// display app help and list all commands. showCommandList()
func (app *App) showApplicationHelp() bool {
	Debugf("render application help and commands list, replaces=%s", maputil.ToString2(app.Replaces()))
	app.Fire(gevent.OnAppHelpBefore, nil)

	// cmdHelpTemplate = color.ReplaceTag(cmdHelpTemplate)
	// render help text template
	s := helper.RenderText(AppHelpTemplate, map[string]any{
		"CmdGroups": app.CommandsByGroup("Available Commands"),
		"GOpts":     app.fs.BuildOptsHelp(),
		// app version
		"Version": app.Version,
		"HasSubs": app.hasSubcommands,
		// always upper first char
		"Desc": strutil.UpperFirst(app.Desc),
		// user custom help vars
		"Vars": app.HelpVars,
		// custom help config
		"Help": app.HelpConfig,
	}, template.FuncMap{
		"paddingName": func(n string) string {
			return strutil.PadRight(n, " ", app.nameMaxWidth)
		},
	})

	// parse help vars and render color tags
	color.Print(app.ReplacePairs(s))
	app.Fire(gevent.OnAppHelpAfter, nil)

	if sysutil.IsLinux() {
		fmt.Println()
	}
	return false
}

// showCommandHelp display help for a command
func (app *App) showCommandHelp(list []string) (code PrepareState) {
	binName := app.Ctx.binName
	// if len(list) == 0 { TODO support multi level sub command?
	if len(list) > 1 {
		color.Error.Tips("Too many arguments given.\n\nUsage: %s help COMMAND", binName)
		return ERR
	}

	// get real name
	name := app.cmdAliases.ResolveAlias(list[0])
	if name == HelpCommand || name == "-h" {
		Debugf("render help command information")

		color.Println("Display help message for application or command.\n")
		color.Printf(`<yellow>Usage:</>
  <cyan>%s COMMAND --help</>
  <cyan>%s COMMAND SUBCOMMAND --help</>
  <cyan>%s COMMAND SUBCOMMAND ... --help</>
  <cyan>%s help COMMAND</>
`, binName, binName, binName, binName)
		return
	}

	cmd, exist := app.Command(name)
	if !exist {
		color.Error.Prompt("Unknown command name '%s'. Run '%s -h' see all commands", name, binName)
		return ERR
	}

	// show help for the give command.
	_ = cmd.ShowHelp()
	return
}

// showAutoCompletion 计算并逐行打印运行期动态补全候选(纯文本, 无颜色), 供 shell 脚本解析。
//
// words 为 shell 传入、已去掉 bin 名的命令行片段; 无候选时不输出。
func (app *App) showAutoCompletion(words []string) {
	for _, item := range app.resolveCompletion(words) {
		fmt.Println(item)
	}
}

// findSimilarCmd find similar cmd by input string
func (app *App) findSimilarCmd(input string) []string {
	var ss []string
	// ins := strings.Split(input, "")
	// fmt.Print(input, ins)
	ln := len(input)

	// NOTE: copy the map. CmdNameMap() returns the real cmdNames map, mutating it
	// here would pollute the command registry(eg add a phantom 'help' command).
	src := app.CmdNameMap()
	names := make(map[string]int, len(src)+1)
	for n, l := range src {
		names[n] = l
	}
	names["help"] = 4 // add built-in 'help' command for matching

	// find from command names
	for name := range names {
		cln := len(name)
		if cln > ln && strings.Contains(name, input) {
			ss = append(ss, name)
		} else if ln > cln && strings.Contains(input, name) {
			// sns := strings.Split(str, "")
			ss = append(ss, name)
		}

		// max find 5 items
		if len(ss) == 5 {
			break
		}
	}

	// find from aliases
	for alias := range app.cmdAliases.Mapping() {
		// max find 5 items
		if len(ss) >= 5 {
			break
		}

		cln := len(alias)
		if cln > ln && strings.Contains(alias, input) {
			ss = append(ss, alias)
		} else if ln > cln && strings.Contains(input, alias) {
			ss = append(ss, alias)
		}
	}

	return ss
}

/*************************************************************
 * display command help
 *************************************************************/

// CmdHelpTemplate help template for a command
var CmdHelpTemplate = `{{.Desc}}
{{if .Cmd.NotStandalone}}
<comment>Name:</> {{.Cmd.Name}}{{if .Cmd.Aliases}} (alias: <info>{{.Cmd.Aliases.String}}</>){{end}}{{end}}
<comment>Usage:</>
  {$binName} [global options] {{if .Cmd.NotStandalone}}<cyan>{{.Cmd.Path}}</> {{end}}[--options ...] [arguments ...]{{ if .SubGroups }}
  {$binName} [global options] {{if .Cmd.NotStandalone}}<cyan>{{.Cmd.Path}}</> {{end}}<cyan>SUBCOMMAND</> [--options ...] [arguments ...]{{end}}
{{if .GOpts}}
<comment>Global Options:</>
{{.GOpts}}{{end}}{{if .Options}}
<comment>Options:</>
{{.Options}}{{end}}{{if .ArgsHelp}}
<comment>Arguments:</>
{{.ArgsHelp}}{{end}}{{range $g := .SubGroups}}
<comment>{{$g.Title}}:</>{{range $c := $g.Cmds}}
  <info>{{$c.Name | paddingName }}</> {{$c.HelpDesc}}{{if $c.Aliases}} (alias: <green>{{ join $c.Aliases ","}}</>){{end}}{{end}}
{{end}}{{.Help.AfterCmdText}}{{if .Cmd.Examples}}
<comment>Examples:</>
{{.Cmd.Examples}}{{end}}{{if .Cmd.Help}}
<comment>Help:</>
{{.Cmd.Help}}{{end}}{{.Help.FooterText}}`

// ShowHelp show command help information
func (c *Command) ShowHelp() (err error) {
	Debugf("render the command '%s' help information", c.Name)

	// custom help render func
	if c.HelpRender != nil {
		c.HelpRender(c)
		return
	}

	// 合并共享(继承)选项, 使 `help CMD` 等不走 parseOptions 的路径也能展示继承选项。幂等。
	c.mergeSharedOpts()

	// clear space and empty new line
	if c.Examples != "" {
		c.Examples = strings.Trim(c.Examples, "\n") + "\n"
	}

	// clear space and empty new line
	if c.Help != "" {
		c.Help = strings.TrimSpace(c.Help) + "\n"
	}

	vars := map[string]any{
		"Cmd":       c,
		"SubGroups": c.CommandsByGroup("Subcommands"),
		// global options
		// - on standalone
		"GOpts": nil,
		// parse options to string
		"Options": c.Flags.BuildOptsHelp(),
		// parse options to string
		"ArgsHelp": c.Flags.BuildArgsHelp(),
		// always upper first char
		"Desc": c.HelpDesc(),
		// user custom help vars
		"Vars": c.HelpVars,
		// custom help config
		"Help": c.HelpConfig,
	}

	// if c.NotStandalone() {
	// 	vars["GOpts"] = c.GFlags().BuildHelp()
	// }

	// render help message
	str := helper.RenderText(CmdHelpTemplate, vars, template.FuncMap{
		"paddingName": func(n string) string {
			return strutil.PadRight(n, " ", c.nameMaxWidth)
		},
	})

	// parse gcli help vars then print help
	// fmt.Printf("%#v\n", s)
	color.Print(c.ReplacePairs(str))
	if sysutil.IsLinux() {
		fmt.Println()
	}
	return
}
