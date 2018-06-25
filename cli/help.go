package cli

import (
    "os"
    "io"
    "html/template"
    "strings"
    "fmt"
)

// help template for all commands
var commandsHelp = `{{.Des}}
Usage:
  {{.Script}} command [--options ...] [arguments ...]

Options:
  -h, --help        Display this help information
  -V, --version     Display this version information

Commands:{{range .Cs}}{{if .Runnable}}
  {{.Name | printf "%-12s"}} {{.Description}}{{if .Aliases}}(alias: {{.Aliases.String}}){{end}}{{end}}{{end}}

Use "{{.Script}} help [command]" for more information about a command
`

// help template for a command
var commandHelp = `{{.Cmd.Description}}

Name: {{.Cmd.Name}}{{if .Cmd.Aliases}}(alias: {{.Cmd.Aliases.String}}){{end}}
Usage: 
  {{.Script}} {{.Cmd.Name}} [-opt ...] [arg0 ...]

Options:

`

// showVersionInfo display version info
func showVersionInfo() {
    fmt.Printf(`%s\n\nVersion: %s\n`, app.Description, app.Version)
    os.Exit(0)
}

// showCommandsHelp commands list
func showCommandsHelp() {
    renderString(os.Stdout, commandsHelp, map[string]interface{}{
        "Cs": commands,
        "Des": app.Description,
        "Script": script,
    })
    os.Exit(0)
}

// showCommandHelp display help for an command
func showCommandHelp(list []string, quit bool) {
    if len(list) != 1 {
        fmt.Fprintf(
            os.Stdout,
            "usage: %s help %s\n\nToo many arguments given.\n",
            script,
            list[0],
        )
        os.Exit(2) // failed at 'bee help'
    }

    // get real name
    name := GetCommandName(list[0])

    if cmd, exist := commands[name]; exist {
        renderString(os.Stdout, commandHelp, map[string]interface{}{
            "Cmd": cmd,
            "Script": script,
        })

        if quit {
            os.Exit(0)
        }

        return
    }

    Stdoutf("Unknown command name %#q.  Run '%s help'.", name, script)
    os.Exit(2)
}

// renderString
func renderString(w io.Writer, text string, data interface{}) {
    t := template.New("cli")

    t.Funcs(template.FuncMap{"trim": func(s template.HTML) template.HTML {
        return template.HTML(strings.TrimSpace(string(s)))
    }})

    t.Funcs(template.FuncMap{"joinAlias": func(ss []string) template.HTML {
        return template.HTML(strings.Join(ss, ","))
    }})

    template.Must(t.Parse(text))

    if err := t.Execute(w, data); err != nil {
        panic(err)
    }
}

