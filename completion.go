package gcli

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/internal/helper"
)

// current supported shell for completion script generate
const (
	ZshShell  = "zsh"
	BashShell = "bash"
)

// shellTpls 内置的各 shell 补全脚本模板表
var shellTpls = map[string]string{
	ZshShell:  zshCompleteScriptTpl,
	BashShell: bashCompleteScriptTpl,
}

// GenCompletionScript 生成指定 shell 的静态补全脚本文本。
//
//   - shell: 目标 shell, 取值 bash|zsh, 其它返回 error。
//   - binName: 可选, 覆盖脚本中使用的 bin 名(如 genac 的 --bin-name);
//     不传则使用当前应用的 bin 名。
//
// 生成所需的数据(BinName、命令名/描述、各命令选项等)均从 app 当前已注册的命令中取得。
func (app *App) GenCompletionScript(shell string, binName ...string) (string, error) {
	tpl, ok := shellTpls[shell]
	if !ok {
		return "", fmt.Errorf("gcli: unsupported shell %q for completion, only allow: bash, zsh", shell)
	}

	// bin 名称: 允许调用方覆盖, 否则用当前应用 bin 名
	rawBin := app.BinName()
	if len(binName) > 0 && binName[0] != "" {
		rawBin = binName[0]
	}
	// 规整: 去掉 ./ 前缀与 .exe 后缀
	name := strings.TrimSuffix(strings.Trim(rawBin, "./"), ".exe")
	fileName := name + "." + shell

	data := map[string]any{
		"Shell":    shell,
		"BinName":  name,
		"FileName": fileName,
	}

	if shell == BashShell {
		data = buildForBashShell(app, data)
	} else {
		data = buildForZshShell(app, data)
	}

	return helper.RenderText(tpl, &data, nil), nil
}

var bashCompleteScriptTpl = `#!/usr/bin/env {{.Shell}}

# ------------------------------------------------------------------------------
#          FILE:  {{.FileName}}
#        AUTHOR:  inhere (https://github.com/inhere)
#       VERSION:  1.0.0
#   DESCRIPTION:  zsh shell complete for cli app: {{.BinName}}
# ------------------------------------------------------------------------------
# Usage: source {{.FileName}}
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

// buildForBashShell 收集 bash 补全脚本模板所需数据: 命令名列表、各命令(组)对应的选项名。
func buildForBashShell(app *App, data map[string]any) map[string]any {
	var cNames []string

	// {cmd name: opts}
	nameOpts := make(map[string]string)

	for n, c := range app.Commands() {
		ops := c.FlagNames()
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
		for opName := range ops {
			pfx := "--"
			if len(opName) == 1 {
				pfx = "-"
			}

			opList = append(opList, pfx+opName)
		}

		nameOpts[key] = strings.Join(opList, " ")
	}

	data["CmdNames"] = cNames
	data["NameOpts"] = nameOpts

	return data
}

var zshCompleteScriptTpl = `#compdef {{.BinName}}
# ------------------------------------------------------------------------------
#          FILE:  {{.FileName}}
#        AUTHOR:  inhere (https://github.com/inhere)
#       VERSION:  1.0.0
#   DESCRIPTION:  zsh shell complete for cli app: {{.BinName}}
# ------------------------------------------------------------------------------
# Usage: source {{.FileName}}

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
        _values 'command options' \{{range $vs}}
            {{.}}{{end}}
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

// buildForZshShell 收集 zsh 补全脚本模板所需数据: 命令名->描述、命令(组)->选项描述列表。
func buildForZshShell(app *App, data map[string]any) map[string]any {
	type opInfos []string

	// {cmd name: cmd des}. in zsh eg: 'build[compile packages and dependencies]'
	nameDes := make(map[string]string)
	// {cmd name: {opt: opt des}}.
	// in zsh eg:
	// '-x[description]:message:action'
	// {-h,--help}'[Show usage message]' // multi name
	nameOpts := make(map[string]opInfos)

	for n, c := range app.Commands() {
		nameDes[c.Name] = fmtDes(c.Desc) + "(alias " + c.Aliases.String() + ")"

		ops := c.FlagNames()
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

		sfx := " \\"
		var i int
		var opis []string
		for opName := range ops {
			i++
			pfx := "--"
			opDes := fmtDes(c.Flags.LookupFlag(opName).Usage)

			if len(opName) == 1 {
				pfx = "-"
			}

			opKey := pfx + opName
			desTpl := "'%s[%s]'%s"

			if shorts := c.ShortNames(opName); len(shorts) > 0 {
				desTpl = "%s'[%s]'%s"
				opKey = fmt.Sprintf("{-%s,%s}", strings.Join(shorts, ",-"), pfx+opName)
			}

			// latest item
			if oplen == i {
				sfx = ""
			}

			opis = append(opis, fmt.Sprintf(desTpl, opKey, opDes, sfx))
		}

		nameOpts[key] = opis
	}

	data["NameDes"] = nameDes
	data["NameOpts"] = nameOpts

	return data
}

// fmtDes 清理描述中的颜色标签与方括号/反引号, 避免破坏 zsh 补全语法。
func fmtDes(str string) string {
	str = color.ClearTag(str)
	return strings.NewReplacer("`", "", "[", "", "]", "").Replace(str)
}
