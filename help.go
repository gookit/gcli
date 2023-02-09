package gcli

import (
	"strings"
	"text/template"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/helper"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
)

/*************************************************************
 * display app help
 *************************************************************/

// display app version info
func (app *App) showVersionInfo() bool {
	Debugf("print application version info")

	color.Printf(
		"%s\n\nVersion: <cyan>%s</>\n",
		strutil.UpperFirst(app.Desc),
		app.Version,
	)

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
<comment>Available Commands:</>{{range $cmdName, $c := .Cs}}{{if $c.Visible}}
  <info>{{$c.Name | paddingName }}</> {{$c.HelpDesc}}{{if $c.Aliases}} (alias: <green>{{ join $c.Aliases ","}}</>){{end}}{{end}}{{end}}
  <info>{{ paddingName "help" }}</> Display help information

Use "<cyan>{$binName} COMMAND -h</>" for more information about a command
`

// display app help and list all commands. showCommandList()
func (app *App) showApplicationHelp() bool {
	Debugf("render application help and commands list, replaces=%s", maputil.ToString2(app.Replaces()))

	// cmdHelpTemplate = color.ReplaceTag(cmdHelpTemplate)
	// render help text template
	s := helper.RenderText(AppHelpTemplate, map[string]any{
		"Cs":    app.commands,
		"GOpts": app.fs.BuildOptsHelp(),
		// app version
		"Version": app.Version,
		"HasSubs": app.hasSubcommands,
		// always upper first char
		"Desc": strutil.UpperFirst(app.Desc),
		// user custom help vars
		"Vars": app.helpVars,
	}, template.FuncMap{
		"paddingName": func(n string) string {
			return strutil.PadRight(n, " ", app.nameMaxWidth)
		},
	})

	// parse help vars and render color tags
	color.Print(app.ReplacePairs(s))
	return false
}

// showCommandHelp display help for a command
func (app *App) showCommandHelp(list []string) (code int) {
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

// show bash/zsh completion
func (app *App) showAutoCompletion(_ []string) {
	// TODO ...
}

// findSimilarCmd find similar cmd by input string
func (app *App) findSimilarCmd(input string) []string {
	var ss []string
	// ins := strings.Split(input, "")
	// fmt.Print(input, ins)
	ln := len(input)

	names := app.CmdNameMap()
	names["help"] = 4 // add 'help' command

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
  {$binName} [global options] {{if .Cmd.NotStandalone}}<cyan>{{.Cmd.Path}}</> {{end}}[--options ...] [arguments ...]{{ if .Subs }}
  {$binName} [global options] {{if .Cmd.NotStandalone}}<cyan>{{.Cmd.Path}}</> {{end}}<cyan>SUBCOMMAND</> [--options ...] [arguments ...]{{end}}
{{if .GOpts}}
<comment>Global Options:</>
{{.GOpts}}{{end}}{{if .Options}}
<comment>Options:</>
{{.Options}}{{end}}{{if .ArgsHelp}}
<comment>Arguments:</>
{{.ArgsHelp}}{{end}}{{ if .Subs }}
<comment>Subcommands:</>{{range $n,$c := .Subs}}
  <info>{{$c.Name | paddingName }}</> {{$c.HelpDesc}}{{if $c.Aliases}} (alias: <green>{{ join $c.Aliases ","}}</>){{end}}{{end}}
{{end}}{{if .Cmd.Examples}}
<comment>Examples:</>
{{.Cmd.Examples}}{{end}}{{if .Cmd.Help}}
<comment>Help:</>
{{.Cmd.Help}}{{end}}`

// ShowHelp show command help info
func (c *Command) ShowHelp() (err error) {
	Debugf("render the command '%s' help information", c.Name)

	// custom help render func
	if c.HelpRender != nil {
		c.HelpRender(c)
		return
	}

	// clear space and empty new line
	if c.Examples != "" {
		c.Examples = strings.Trim(c.Examples, "\n") + "\n"
	}

	// clear space and empty new line
	if c.Help != "" {
		c.Help = strings.TrimSpace(c.Help) + "\n"
	}

	vars := map[string]any{
		"Cmd":  c,
		"Subs": c.commands,
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
		"Vars": c.helpVars,
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
	return
}
