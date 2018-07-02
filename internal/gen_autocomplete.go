package internal

import (
	"github.com/gookit/cliapp"
	"github.com/gookit/cliapp/utils"
	"strings"
	"os"
	"github.com/gookit/color"
)

//
var genOpts = struct {
	shell   string
	binName string
	output string
}{}

func GenShAutoComplete() *cliapp.Command {
	cmd := cliapp.Command{
		Fn: doGen,
		Name:    "gen-ac",
		Aliases: []string{"gen"},

		Description: "generate script file for command auto complete",
	}

	cmd.StrOpt(
		&genOpts.shell,
		"shell",
		"",
		"bash",
		"the shell env name for want generated, allow: sh,zsh,bash",
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

	color.LiteTips("info").Printf("%+v\n", genOpts)

	var cNames []string
	fNameOpts := make(map[string]string)

	for n, c := range cliapp.AllCommands() {
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
				opList = append(opList, "-" + st)
			}

			opList = append(opList, "--" + op)
		}

		fNameOpts[key] = strings.Join(opList, " ")
	}

	str := utils.RenderTemplate(autoCompleteScriptTpl, map[string]interface{}{
		"Shell": genOpts.shell,
		"BinName": genOpts.binName,
		"CmdNames": cNames,
		"NameOpts": fNameOpts,
	})

	color.Infoln("Now, will write content to file ", genOpts.output)

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

var autoCompleteScriptTpl = `#!/usr/bin/env {{.Shell}}

#
# usage: source ./auto-completion.{{.Shell}}
# run 'complete' to see registered complete function.
#

_complete_for_{{.BinName}}()
{
    local cur prev
    _get_comp_words_by_ref -n = cur prev

    COMPREPLY=()
    commands="{{join .CmdNames " "}}"

    case "$prev" in{{range $k,$v := .NameOpts}}
        {{$k}})
            COMPREPLY=($(compgen -W "{{$v}}" -- "$cur"))
            return 0
            ;;{{end}}
    esac

    COMPREPLY=($(compgen -W "$commands" -- "$cur"))

} &&
# complete -F {auto_complete_func} {bin_filename}
complete -F _complete_for_{{.BinName}} {{.BinName}} {{.BinName}}.exe
`
