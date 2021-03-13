#!/usr/bin/env bash

# ------------------------------------------------------------------------------
#          FILE:  auto-completion.bash
#        AUTHOR:  inhere (https://github.com/inhere)
#       VERSION:  1.0.0
#   DESCRIPTION:  zsh shell complete for cli app: cliapp
# ------------------------------------------------------------------------------
# usage: source auto-completion.bash
# run 'complete' to see registered complete function.


_complete_for_cliapp () {
    local cur prev
    _get_comp_words_by_ref -n = cur prev

    COMPREPLY=()
    commands="exp ex example env-info ei env git-info git clr colors color"

    case "$prev" in
        clr|colors|color)
            COMPREPLY=($(compgen -W "--id -c --dir" -- "$cur"))
            return 0
            ;;
        env-info|ei|env)
            COMPREPLY=($(compgen -W "--id -c -d --dir" -- "$cur"))
            return 0
            ;;
        exp|ex|example)
            COMPREPLY=($(compgen -W "-d --dir -o --opt -n --names" -- "$cur"))
            return 0
            ;;
        git-info|git)
            COMPREPLY=($(compgen -W "--id -c -d --dir" -- "$cur"))
            return 0
            ;;
        help)
            COMPREPLY=($(compgen -W "$commands" -- "$cur"))
            return 0
            ;;
    esac

    COMPREPLY=($(compgen -W "$commands" -- "$cur"))

} &&
# complete -F {auto_complete_func} {bin_filename}
# complete -F _complete_for_cliapp -A file cliapp cliapp.exe
complete -F _complete_for_cliapp cliapp cliapp.exe
