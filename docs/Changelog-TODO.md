# Changelog

## TODO

- [ ] hook on set flag value
- [x] option support multi shorts names
- [ ] cmd support flag option category
- [ ] app support command category by `c.Category`
- [ ] print parent's options on subcommand help panel
- [ ] prompt completion by readline
- [ ] keyboard and cursor control on terminal
- [ ] refactor gflag.Flags, remove dep the `flag.FlagSet`
- [ ] collect option/argument value by interactive `Option.Question`
- [ ] support all command docs to markdown

readline refers:

- https://github.com/chzyer/readline/tree/master/example
- https://github.com/abiosoft/ishell/blob/master/completer.go

## v3.0.1

**new**

- [x] add some special flag type vars
- [x] support hidden command on render help by `c.Hidden=true`

**fixed**

- [x] alias not works on command ID
- [x] render color on command/option/argument description

## v3.0.0

**new**

- [x] support multi level sub commands
- [x] support parse flags from struct tags
- [x] support flag/argument validate
- [ ] support controller on application `app.controllers []Controller`
  - 独立于commands之外的。Independent of commands.
  - 支持组选项，全部子命令都拥有这些选项 `Config/GroupOptions()` 里绑定组选项。
- [x] 支持单个command、controller独立运行
