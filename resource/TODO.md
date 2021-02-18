# TODO

- [ ] support multi level sub commands
- [ ] support parse flags from struct tags
- [ ] support flag/argument validate
- [ ] support controller on application `app.controllers []Controller`
  - 独立于commands之外的。Independent of commands.
  - 支持组选项，全部子命令都拥有这些选项 `Config/GroupOptions()` 里绑定组选项。
  - 是否支持附加子独立命令？
- 是否支持单个command、controller独立运行