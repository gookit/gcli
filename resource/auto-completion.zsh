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
    'test[This is a description <info>message</> for command test(alias ts)]'
   )

  if (( CURRENT == 2 )); then
    # explain commands
    _values 'cliapp commands' ${commands[@]}
    return
  fi

  case ${words[2]} in
  clr|colors|color)
      _arguments -s -w : \
      "--id[the id option]"\ 
      "-c[the config option]"\ 
      "--dir[the dir option]"\ 
      ;;
  env-info|ei|env)
      _arguments -s -w : \
      "--id[the id option]"\ 
      "-c[the config option]"\ 
      "-d[the dir option]"\ 
      "--dir[the dir option]"\ 
      ;;
  exp|ex|example)
      _arguments -s -w : \
      "-n[the option message]"\ 
      "--names[the option message]"\ 
      "-d[the `DIRECTORY` option]"\ 
      "--dir[the `DIRECTORY` option]"\ 
      "-o[the option message]"\ 
      "--opt[the option message]"\ 
      ;;
  git-info|git)
      _arguments -s -w : \
      "-d[the dir option]"\ 
      "--dir[the dir option]"\ 
      "--id[the id option]"\ 
      "-c[the config option]"\ 
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
