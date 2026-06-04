# TODO

## 结构性重构（待评审 → 待实施）

- [ ] A / F / G / H 四项重构，详见 [plans/refactor-A-F-G.md](plans/refactor-A-F-G.md)
  - F：合并 `Parse`/`ParseOpts` 重复校验（低风险）
  - A：全局选项单一数据源（中风险）
  - G：`findCommandName` 去副作用（中风险）
  - H：`helper/`→`internal/helper/`、清理 `gclicom/` 与死代码（机械，但对外破坏性）
