#compdef cliapp
# ------------------------------------------------------------------------------
#          FILE:  auto-completion.zsh
#        AUTHOR:  inhere (https://github.com/inhere)
#       VERSION:  1.0.0
#   DESCRIPTION:  zsh shell complete for cli app: cliapp
# ------------------------------------------------------------------------------
# usage: source auto-completion.zsh

_complete_for_cliapp () {
    typeset -a commands
    commands+=(
        'color[This is a example for cli color usage(alias clr,colors)]'
        'env[Collect project info by git info(alias env-info,ei)]'
        'example[This is a description message(alias exp,ex)]'
        'git[Collect project info by git info(alias git-info)]'
        'test[This is a description message for command test(alias ts)]'
        'help[Display help information]'
    )

    if (( CURRENT == 2 )); then
        # explain commands
        _values 'cliapp commands' ${commands[@]}
        return
    fi

    case ${words[2]} in
    clr|colors|color)
        _values 'command options' \
            '--id[the id option]' \
            '-c[the config option]' \
            '--dir[the dir option]'
        ;;
    env-info|ei|env)
        _values 'command options' \
            '--id[the id option]' \
            '-c[the config option]' \
            {-d,--dir}'[the dir option]'
        ;;
    exp|ex|example)
        _values 'command options' \
            {-n,--names}'[the option message]' \
            {-d,--dir}'[the DIRECTORY option]' \
            {-o,--opt}'[the option message]'
        ;;
    git-info|git)
        _values 'command options' \
            {-d,--dir}'[the dir option]' \
            '--id[the id option]' \
            '-c[the config option]'
        ;;
    help)
        _values "${commands[@]}"
        ;;
    *)
        # use files by default
        _files
        ;;
    esac
}

compdef _complete_for_cliapp cliapp
compdef _complete_for_cliapp cliapp.exe
