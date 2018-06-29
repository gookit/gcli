package cliapp

import (
	"os"
	"io"
	"html/template"
	"strings"
	"github.com/gookit/cliapp/color"
	"bytes"
	"fmt"
	"github.com/gookit/cliapp/utils"
)

// showVersionInfo display version info
func (app *Application) showVersionInfo() {
	color.Printf(
		"%s\n\nVersion: <suc>%s</>\n",
		utils.UpperFirst(app.Description),
		app.Version,
	)

	os.Exit(0)
}

// help template for all commands
var commandsHelp = `{{.Description|raw}} (Version: <info>{{.Version}}</>)
<comment>Usage:</>
  {$binName} [global options...] <info>{command}</> [--option ...] [argument ...]

<comment>Global Options:</>
      <info>--verbose</>     Set error reporting level(quiet 0 - 4 debug)
      <info>--no-color</>    Disable color when outputting message
  <info>-h, --help</>        Display the help information
  <info>-V, --version</>     Display app version information

<comment>Available Commands:</>{{range .Cs}}{{if .Runnable}}
  <info>{{.Name | printf "%-12s"}}</> {{.Description|colored}}{{if .Aliases}} (alias: <cyan>{{.Aliases.String}}</>){{end}}{{end}}{{end}}
  <info>help</>         Display help information

Use "<cyan>{$binName} help {command}</>" for more information about a command
`

// showCommandsHelp commands list
func showCommandsHelp() {
	// use buffer receive rendered content
	var buf bytes.Buffer

	commandsHelp = color.ReplaceTag(commandsHelp)

	// RenderStrTpl(os.Print, commandsHelp, map[string]interface{}{
	RenderStrTpl(&buf, commandsHelp, map[string]interface{}{
		"Cs":      commands,
		"Version": app.Version,
		// always upper first char
		"Description": utils.UpperFirst(app.Description),
	})

	fmt.Print(ReplaceVars(buf.String(), app.vars))

	os.Exit(0)
}

// RenderStrTpl
func RenderStrTpl(w io.Writer, text string, data interface{}) {
	t := template.New("cli")

	// don't escape content
	t.Funcs(template.FuncMap{"raw": func(s string) template.HTML {
		return template.HTML(s)
	}})

	// upper first char
	t.Funcs(template.FuncMap{"upFirst": func(s string) template.HTML {
		return template.HTML(utils.UpperFirst(s))
	}})

	t.Funcs(template.FuncMap{"colored": func(s string) template.HTML {
		return template.HTML(color.ReplaceTag(s))
	}})

	t.Funcs(template.FuncMap{"coloredHtml": func(h template.HTML) template.HTML {
		return template.HTML(color.ReplaceTag(string(h)))
	}})

	t.Funcs(template.FuncMap{"trim": func(s template.HTML) template.HTML {
		return template.HTML(strings.TrimSpace(string(s)))
	}})

	t.Funcs(template.FuncMap{"joinStrings": func(ss []string) template.HTML {
		return template.HTML(strings.Join(ss, ","))
	}})

	template.Must(t.Parse(text))

	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}
