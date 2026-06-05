package gcli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/internal/helper"
)

// current supported shell for completion script generate
const (
	ZshShell  = "zsh"
	BashShell = "bash"
	PwshShell = "pwsh" // PowerShell
)

// hasMetaFlag 判断 args 中是否含补全/生成元选项 token(--in-completion / --gen-completion),
// 兼容 --gen-completion 与 --gen-completion=bash 两种写法。
// 命中时本次运行为补全/生成请求, 需进入静默模式(抑制用户生命周期钩子)。
func hasMetaFlag(args []string) bool {
	for _, arg := range args {
		if arg == "--in-completion" || strings.HasPrefix(arg, "--in-completion=") ||
			arg == "--gen-completion" || strings.HasPrefix(arg, "--gen-completion=") {
			return true
		}
	}
	return false
}

// resolveCompletion 计算运行期动态补全候选。
//
// 输入 words 是 shell 传入、已去掉 bin 名的命令行片段(--in-completion 选项本身已被解析消费)。
// 约定: 最后一个元素是"当前正在输入的词"(可能为空), 其余是"已完成的前序词"。
//
// 返回去重、排序后的候选列表(命令名/子命令名/选项名)。
func (app *App) resolveCompletion(words []string) []string {
	// 1. 没有任何片段: 返回所有顶层命令名 + 别名 + help
	if len(words) == 0 {
		return app.topLevelNames()
	}

	// 2. 拆分: cur 为当前正在输入的词, prev 为已完成的前序词
	cur := words[len(words)-1]
	prev := words[:len(words)-1]

	// 3. 用 prev 定位"当前命令上下文": 从 app 顶层开始, 逐个处理非选项词并尝试下钻。
	//    - 选项词(以 - 开头)在定位时直接跳过;
	//    - 非选项词若是当前层的命令/子命令(经 ResolveAlias 解析别名)则下钻;
	//    - 一旦遇到不是命令的非选项词(视为参数), 立即停止下钻。
	var curCmd *Command // 当前定位到的命令节点(nil 表示仍在 app 顶层)
	for _, word := range prev {
		if strings.HasPrefix(word, "-") {
			// 选项词不参与命令定位
			continue
		}

		name := word
		var next *Command
		if curCmd == nil {
			// 仍在 app 顶层: 解析顶层别名后查找命令
			name = app.ResolveAlias(name)
			next, _ = app.Command(name)
		} else {
			// 已在某命令下: 解析该命令的子命令别名后查找子命令
			name = curCmd.ResolveAlias(name)
			next, _ = curCmd.Command(name)
		}

		if next == nil {
			// 当前词不是命令(视为参数), 停止下钻, 上下文保持不变
			break
		}
		curCmd = next
	}

	// 3.5 选项值补全: 当前词不是选项, 且前一个词是"取值型"选项时, 补全该选项的候选值。
	//     - 有 Choices: 给候选值; 无 Choices: 返回空(交给 shell 做文件名补全)。
	//     - 前一个词是 bool 选项(不取值)时: 跳过, 落到下面的命令/子命令补全。
	if !strings.HasPrefix(cur, "-") && len(prev) > 0 {
		last := prev[len(prev)-1]
		if strings.HasPrefix(last, "-") {
			if opt := optByRef(optsOfNode(curCmd, app), last); opt != nil && opt.TakesValue() {
				return filterAndSort(opt.Choices, cur)
			}
		}
	}

	// 4. 根据 cur 产出候选
	var items []string
	if strings.HasPrefix(cur, "-") {
		// cur 以 - 开头: 候选为当前节点的选项名(长选项 --name / 短选项 -x)
		items = completionOptNames(curCmd, app)
	} else if curCmd == nil {
		// 当前节点为 app 顶层: 候选为顶层命令名 + 别名 + help
		items = app.topLevelNames()
	} else {
		// 当前节点为某命令: 候选为其子命令名
		items = completionSubNames(curCmd)
	}

	// 5. 用 cur 做前缀过滤, 去重、排序后返回
	return filterAndSort(items, cur)
}

// topLevelNames 返回顶层补全名: 所有顶层命令名 + 命令别名 + 内置 help(去重、排序)。
func (app *App) topLevelNames() []string {
	var names []string
	for name := range app.CmdNameMap() {
		names = append(names, name)
	}
	// 顶层命令别名
	for alias := range app.AliasesMapping() {
		names = append(names, alias)
	}
	// 内置 help 命令
	names = append(names, HelpCommand)

	return filterAndSort(names, "")
}

// completionSubNames 返回命令 c 的子命令名 + 子命令别名(去重、排序)。
func completionSubNames(c *Command) []string {
	var names []string
	for name := range c.CmdNameMap() {
		names = append(names, name)
	}
	for alias := range c.AliasesMapping() {
		names = append(names, alias)
	}
	return names
}

// completionOptNames 返回某节点的选项名候选: 长选项 --name 与短选项 -x。
// node 为 nil 时表示 app 顶层, 取 app 的全局选项。
func completionOptNames(node *Command, app *App) []string {
	var names []string

	addOpts := func(opts map[string]*CliOpt, shortFn func(string) []string) {
		for name, opt := range opts {
			// 跳过隐藏选项(如框架内部的 --in-completion), 不应出现在补全候选中
			if opt.Hidden {
				continue
			}
			names = append(names, "--"+name)
			// 收集该选项的短名
			for _, short := range shortFn(name) {
				names = append(names, "-"+short)
			}
		}
	}

	if node == nil {
		// app 顶层: 全局选项
		addOpts(app.fs.Opts(), app.fs.ShortNames)
	} else {
		addOpts(node.Opts(), node.ShortNames)
	}
	return names
}

// optsOfNode 返回某节点的选项元数据表: node 为 nil 时取 app 全局选项, 否则取命令自身选项。
func optsOfNode(node *Command, app *App) map[string]*CliOpt {
	if node == nil {
		return app.fs.Opts()
	}
	return node.Opts()
}

// optByRef 按命令行里的选项引用(--name / -x)解析出对应的 CliOpt; 找不到返回 nil。
func optByRef(opts map[string]*CliOpt, ref string) *CliOpt {
	name := strings.TrimLeft(ref, "-")
	if name == "" {
		return nil
	}
	// 长名直接命中
	if opt, ok := opts[name]; ok {
		return opt
	}
	// 否则按短名匹配
	for _, opt := range opts {
		for _, s := range opt.Shorts {
			if s == name {
				return opt
			}
		}
	}
	return nil
}

// filterAndSort 用 prefix 做前缀过滤, 并对结果去重、排序。
func filterAndSort(items []string, prefix string) []string {
	seen := make(map[string]struct{}, len(items))
	var out []string
	for _, it := range items {
		if it == "" {
			continue
		}
		if prefix != "" && !strings.HasPrefix(it, prefix) {
			continue
		}
		if _, ok := seen[it]; ok {
			continue
		}
		seen[it] = struct{}{}
		out = append(out, it)
	}
	sort.Strings(out)
	return out
}

// shellTpls 内置的各 shell **静态**补全脚本模板表(嵌入式: 把命令/选项硬编码进脚本)。
var shellTpls = map[string]string{
	ZshShell:  zshCompleteScriptTpl,
	BashShell: bashCompleteScriptTpl,
}

// dynamicShellTpls 内置的各 shell **动态(瘦)**补全脚本模板表。
// 瘦脚本不硬编码命令/选项, 而是回调二进制 `bin --in-completion <已输入词...>` 取候选, 零维护。
var dynamicShellTpls = map[string]string{
	ZshShell:  zshDynamicTpl,
	BashShell: bashDynamicTpl,
	PwshShell: pwshDynamicTpl,
}

// unsupportedShellErr 构造"不支持的 shell"错误, 允许列表由实际可用的模板表动态给出。
func unsupportedShellErr(shell string, tpls map[string]string) error {
	allow := make([]string, 0, len(tpls))
	for name := range tpls {
		allow = append(allow, name)
	}
	sort.Strings(allow)
	return fmt.Errorf("gcli: unsupported shell %q for completion, only allow: %s", shell, strings.Join(allow, ", "))
}

// completionFileName 推荐的脚本文件名: pwsh 用 .ps1 扩展名, 其余用 shell 名作扩展名。
func completionFileName(name, shell string) string {
	if shell == PwshShell {
		return name + ".ps1"
	}
	return name + "." + shell
}

// normalizeBinName 取脚本中使用的 bin 名: 允许调用方覆盖, 否则用当前应用 bin 名;
// 并规整(去掉 ./ 前缀与 .exe 后缀)。
func (app *App) normalizeBinName(binName ...string) string {
	rawBin := app.BinName()
	if len(binName) > 0 && binName[0] != "" {
		rawBin = binName[0]
	}
	return strings.TrimSuffix(strings.Trim(rawBin, "./"), ".exe")
}

// GenCompletionScript 生成指定 shell 的**动态(瘦)**补全脚本文本(默认方式)。
//
// 瘦脚本不硬编码命令/选项, 而是回调 `bin --in-completion <已输入词...>` 动态取候选,
// 命令/选项变化后无需重新生成脚本, 零维护。
//
//   - shell: 目标 shell, 取值 bash|zsh, 其它返回 error。
//   - binName: 可选, 覆盖脚本中使用的 bin 名(如 genac 的 --bin-name);
//     不传则使用当前应用的 bin 名。
func (app *App) GenCompletionScript(shell string, binName ...string) (string, error) {
	tpl, ok := dynamicShellTpls[shell]
	if !ok {
		return "", unsupportedShellErr(shell, dynamicShellTpls)
	}

	name := app.normalizeBinName(binName...)
	data := map[string]any{
		"Shell":    shell,
		"BinName":  name,
		"FileName": completionFileName(name, shell),
	}

	return helper.RenderText(tpl, &data, nil), nil
}

// GenStaticCompletionScript 生成指定 shell 的**静态(嵌入式)**补全脚本文本。
//
// 静态脚本把当前已注册的命令名/描述/选项硬编码进脚本; 命令/选项变化后需重新生成。
// 一般推荐使用 GenCompletionScript(瘦/动态); 仅在无法回调二进制等场景下用此 opt-in 方式。
//
//   - shell: 目标 shell, 取值 bash|zsh, 其它返回 error。
//   - binName: 可选, 覆盖脚本中使用的 bin 名(如 genac 的 --bin-name);
//     不传则使用当前应用的 bin 名。
//
// 生成所需的数据(BinName、命令名/描述、各命令选项等)均从 app 当前已注册的命令中取得。
func (app *App) GenStaticCompletionScript(shell string, binName ...string) (string, error) {
	tpl, ok := shellTpls[shell]
	if !ok {
		return "", unsupportedShellErr(shell, shellTpls)
	}

	name := app.normalizeBinName(binName...)
	data := map[string]any{
		"Shell":    shell,
		"BinName":  name,
		"FileName": completionFileName(name, shell),
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

// bashDynamicTpl 瘦(动态) bash 补全脚本模板:
// 把"光标前的已输入词(含当前词)"传给 `bin --in-completion`, 每行输出作为一个候选喂给 compgen。
var bashDynamicTpl = `#!/usr/bin/env {{.Shell}}

# ------------------------------------------------------------------------------
#          FILE:  {{.FileName}}
#        AUTHOR:  inhere (https://github.com/inhere)
#       VERSION:  1.0.0
#   DESCRIPTION:  dynamic bash complete for cli app: {{.BinName}}
#                 it delegates candidate computing to: {{.BinName}} --in-completion
# ------------------------------------------------------------------------------
# Usage: source {{.FileName}}

_complete_for_{{.BinName}} () {
    local cur words
    cur="${COMP_WORDS[COMP_CWORD]}"
    # 已输入的词(从命令名之后到当前光标处, 含当前正在输入的词), 交给二进制计算候选
    words=("${COMP_WORDS[@]:1:$COMP_CWORD}")

    local IFS=$'\n'
    # 2>/dev/null 丢弃 stderr, 只取 stdout 的候选行
    COMPREPLY=( $(compgen -W "$("{{.BinName}}" --in-completion "${words[@]}" 2>/dev/null)" -- "$cur") )
}

# complete -F {auto_complete_func} {bin_filename}
complete -F _complete_for_{{.BinName}} {{.BinName}} {{.BinName}}.exe
`

// zshDynamicTpl 瘦(动态) zsh 补全脚本模板:
// 把"命令名之后的已输入词"传给 `bin --in-completion`, 按行切分后用 compadd 作为候选。
var zshDynamicTpl = `#compdef {{.BinName}}
# ------------------------------------------------------------------------------
#          FILE:  {{.FileName}}
#        AUTHOR:  inhere (https://github.com/inhere)
#       VERSION:  1.0.0
#   DESCRIPTION:  dynamic zsh complete for cli app: {{.BinName}}
#                 it delegates candidate computing to: {{.BinName}} --in-completion
# ------------------------------------------------------------------------------
# Usage: source {{.FileName}}

_complete_for_{{.BinName}} () {
    local -a candidates
    # ${words[@]:1} 去掉命令名, 余下为已输入词(含当前词); ${(@f)...} 按行切分为数组
    candidates=("${(@f)$("{{.BinName}}" --in-completion "${words[@]:1}" 2>/dev/null)}")
    compadd -- $candidates
}

compdef _complete_for_{{.BinName}} {{.BinName}}
`

// pwshDynamicTpl PowerShell 动态(瘦)补全脚本: 注册原生参数补全, 回调 bin --in-completion 取候选。
// 注: PowerShell 模板里无反引号(避免与 Go 原始字符串冲突), $ 与单 { } 不会被 Go 模板引擎处理。
var pwshDynamicTpl = `# ------------------------------------------------------------------------------
#          FILE:  {{.FileName}}
#        AUTHOR:  inhere (https://github.com/inhere)
#       VERSION:  1.0.0
#   DESCRIPTION:  dynamic powershell complete for cli app: {{.BinName}}
#                 it delegates candidate computing to: {{.BinName}} --in-completion
# ------------------------------------------------------------------------------
# Usage: add to your $PROFILE, or dot-source it: . {{.FileName}}
# NOTE: recommend PowerShell 7+ (better native empty-arg passing).

Register-ArgumentCompleter -Native -CommandName '{{.BinName}}' -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)

    # collect typed words, skip the executable name (CommandElements[0])
    $words = @()
    for ($i = 1; $i -lt $commandAst.CommandElements.Count; $i++) {
        $words += $commandAst.CommandElements[$i].ToString()
    }
    # trailing space(completing a new word): wordToComplete is empty and not in
    # $words, append '' as the "current word" so the binary completes children.
    if ([string]::IsNullOrEmpty($wordToComplete)) {
        $words += ''
    }

    # delegate candidate computing to the binary(one per line), drop stderr.
    & '{{.BinName}}' --in-completion @words 2>$null | ForEach-Object {
        [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
    }
}
`
