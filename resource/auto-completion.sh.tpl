#!/usr/bin/env {{.Shell}}

#
# usage:
#   run: source ./auto-completion.{{.Shell}}
# run `complete` to see registered complete function.
#

_console_app_{{.BinName}}()
{
	local cur prev
	_get_comp_words_by_ref -n = cur prev

    COMPREPLY=()
	commands="{{.Commands}}"

	case "$prev" in
		example|exp)
			COMPREPLY=($(compgen -W "--id --dir --opt --names" -- "$cur"))
			return 0
			;;
	esac

	COMPREPLY=($(compgen -W "$commands" -- "$cur"))

} &&
# complete -F {auto_complete_func} {bin_filename}
complete -F _console_app_{{.BinName}} {{.BinName}} {{.BinName}}.exe
