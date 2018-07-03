package builtin

import (
	"github.com/gookit/cliapp"
	"github.com/gookit/cliapp/utils"
	"strings"
	"os"
	"github.com/gookit/color"
	"github.com/gookit/cliapp/interact"
)

const (
	ZshShell  = "zsh"
	BashShell = "bash"
)

//
var genOpts = struct {
	shell   string
	binName string
	output  string
}{}

var shellTpls = map[string]string{
	"zsh":  zshCompleteScriptTpl,
	"bash": bashCompleteScriptTpl,
}

func GenShAutoComplete() *cliapp.Command {
	cmd := cliapp.Command{
		Fn:      doGen,
		Name:    "gen-ac",
		Aliases: []string{"genac"},

		Description: "generate script file for command auto complete",
	}

	shell := utils.GetCurShell(true)

	if shell == "" {
		shell = "bash"
	}

	cmd.StrOpt(
		&genOpts.shell,
		"shell",
		"s",
		shell,
		"the shell env name for want generated, allow: zsh,bash",
	).StrOpt(
		&genOpts.binName,
		"bin-name",
		"b",
		"",
		"your packaged application bin file name.",
	).StrOpt(
		&genOpts.output,
		"output",
		"o",
		"auto-completion.{shell}",
		"output shell auto completion script file name.",
	)

	return &cmd
}

func doGen(cmd *cliapp.Command, args []string) int {
	if len(genOpts.binName) == 0 {
		genOpts.binName = cliapp.BinName()
	}

	genOpts.binName = strings.TrimSuffix(strings.Trim(genOpts.binName, "./"), ".exe")

	if len(genOpts.output) == 0 {
		genOpts.output = genOpts.binName + "." + genOpts.shell
	} else {
		genOpts.output = strings.Replace(genOpts.output, "{shell}", genOpts.shell, 1)

		// check suffix
		if !strings.Contains(genOpts.output, ".") {
			genOpts.output += "." + genOpts.shell
		}
	}

	color.LiteTips("info").Printf("\n  %+v\n", genOpts)

	data := map[string]interface{}{
		"Shell":    genOpts.shell,
		"BinName":  genOpts.binName,
		"FileName": genOpts.output,
	}

	if genOpts.shell == BashShell {
		data = buildForBashShell(data)
	} else if genOpts.shell == ZshShell {
		data = buildForZshShell(data)
	} else {
		color.LiteTips("error").Println("--shell option only allow: zsh,bash")

		return -2
	}

	str := utils.RenderTemplate(shellTpls[genOpts.shell], &data)

	color.Infoln("Now, will write content to file ", genOpts.output)
	color.Normal("Continue?")

	if !interact.AnswerIsYes(true) {
		color.Info("\nBye :)\n")

		return 0
	}

	// 以读写方式打开文件，如果不存在，则创建
	file, err := os.OpenFile(genOpts.output, os.O_RDWR|os.O_CREATE, 0766)

	if err != nil {
		color.Errorln("Open file error: ", err.Error())

		return -2
	}

	_, err = file.WriteString(str)

	if err != nil {
		color.Errorln("Write file error: ", err.Error())

		return -2
	}

	color.Sucln("\nOK, auto-complete file generate successful")

	return 0
}

var bashCompleteScriptTpl = `#!/usr/bin/env {{.Shell}}

# ------------------------------------------------------------------------------
#          FILE:  {{.FileName}}
#        AUTHOR:  inhere (https://github.com/inhere)
#       VERSION:  1.0.0
#   DESCRIPTION:  zsh shell complete for cli app: {{.BinName}}
# ------------------------------------------------------------------------------
# usage: source {{.FileName}}
# run 'complete' to see registered complete function.


_complete_for_{{.BinName}} () {
    local cur prev
    _get_comp_words_by_ref -n = cur prev

    COMPREPLY=()
    commands="{{join .CmdNames " "}} help"

    case "$prev" in{{range $k,$v := .NameOpts}}
        {{$k}})
            COMPREPLY=($(compgen -W "{{$v}}" -- "$cur"))
            return 0
            ;;{{end}}
        help)
            COMPREPLY=($(compgen -W "$commands" -- "$cur"))
            return 0
            ;;
    esac

    COMPREPLY=($(compgen -W "$commands" -- "$cur"))

} &&
# complete -F {auto_complete_func} {bin_filename}
# complete -F _complete_for_{{.BinName}} -A file {{.BinName}} {{.BinName}}.exe
complete -F _complete_for_{{.BinName}} {{.BinName}} {{.BinName}}.exe
`

func buildForBashShell(data map[string]interface{}) map[string]interface{} {
	var cNames []string

	// {cmd name: opts}
	nameOpts := make(map[string]string)

	for n, c := range cliapp.AllCommands() {
		// skip self
		if n == "genac" || n == "gen-ac" {
			continue
		}

		ops := c.OptNames()
		if len(ops) == 0 {
			continue
		}

		ns := c.Aliases
		key := n

		if len(ns) > 0 {
			ns = append(ns, n)
			key = strings.Join(ns, "|")
			cNames = append(cNames, ns...)
		} else {
			cNames = append(cNames, n)
		}

		var opList []string
		for op, st := range ops {
			if st != "" {
				opList = append(opList, "-"+st)
			}

			pfx := "--"

			if len(op) == 1 {
				pfx = "-"
			}

			opList = append(opList, pfx+op)
		}

		nameOpts[key] = strings.Join(opList, " ")
	}

	data["CmdNames"] = cNames
	data["NameOpts"] = nameOpts

	return data
}

var zshCompleteScriptTpl = `# ------------------------------------------------------------------------------
#          FILE:  {{.FileName}}
#        AUTHOR:  inhere (https://github.com/inhere)
#       VERSION:  1.0.0
#   DESCRIPTION:  zsh shell complete for cli app: {{.BinName}}
# ------------------------------------------------------------------------------
# usage: source {{.FileName}}

_complete_for_{{.BinName}} () {
   typeset -a commands
   commands+=({{range $k,$v := .NameDes}}
    '{{$k}}[{{$v}}]'{{end}}
    'help[Display help information]'
   )

  if (( CURRENT == 2 )); then
    # explain commands
    _values 'cliapp commands' ${commands[@]}
    return
  fi

  case ${words[2]} in{{range $k,$vs := .NameOpts}}
  {{$k}})
      _arguments -s -w \{{range $vs}}
        "{{.N}}[{{.V}}]" {{.Sfx}}{{end}}
      ;;{{end}}
  help)
      _values "${commands[@]}"
      ;;
  *)
      # use files by default
      _files
      ;;
  esac
}

compdef _complete_for_{{.BinName}} {{.BinName}}
compdef _complete_for_{{.BinName}} {{.BinName}}.exe
`

func buildForZshShell(data map[string]interface{}) map[string]interface{} {
	type opInfo struct{ N, V, Sfx string }
	type opInfos []opInfo

	// {cmd name: cmd des}. in zsh eg: 'build[compile packages and dependencies]'
	nameDes := make(map[string]string)
	// {cmd name: {opt: opt des}}. in zsh eg: '-n[print the commands but do not run them]'
	nameOpts := make(map[string]opInfos)

	for n, c := range cliapp.AllCommands() {
		// skip self
		if n == "genac" || n == "gen-ac" {
			continue
		}
		nameDes[c.Name] = fmtDes(c.Description) + "(alias " + c.Aliases.String() + ")"

		ops := c.OptNames()
		oplen := len(ops)
		if oplen == 0 {
			continue
		}

		ns := c.Aliases
		key := n

		if len(ns) > 0 {
			ns = append(ns, n)
			key = strings.Join(ns, "|")
		}

		sfx := "\\"
		var i int
		var opis opInfos
		for op, st := range ops {
			i++
			opDes := fmtDes(c.Flags.Lookup(op).Usage)

			if st != "" {
				opis = append(opis, opInfo{"-" + st, opDes, sfx})
			}

			pfx := "--"

			if len(op) == 1 {
				pfx = "-"
			}

			if oplen == i {
				sfx = ""
			}

			opis = append(opis, opInfo{pfx + op, opDes, sfx})
		}

		nameOpts[key] = opis
	}

	data["NameDes"] = nameDes
	data["NameOpts"] = nameOpts

	return data
}

func fmtDes(str string) string {
	str = color.ClearTag(str)

	return strings.NewReplacer("`", "").Replace(str)
}
