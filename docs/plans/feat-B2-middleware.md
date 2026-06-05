# 功能实现计划：B2 命令中间件 middleware

> 状态：**待评审**
> 范围：让命令支持 before/after 切面（鉴权、日志、计时等通用逻辑）。

## 现状

骨架已有但**未接入**：

- `cmd.go:33` `RunnerFunc`、`cmd.go:43` `HandlersChain`、`cmd.go:46` `HandlersChain.Last()`
- `cmd.go:83-85` `middleIdx int8` / `middles HandlersChain` 字段
- `cmd.go:324` `Next()` 标 `// Next TODO`，会遍历 `middles` 但**没有任何地方调用它**，
  且不会执行 `c.Func`
- `cmd.go:522` `doExecute` 里 `cmd.go:556` 直接 `c.Func(c, fnArgs)`，绕过中间件

即：有数据结构和半个 `Next()`，但没有注册 API、没有接入执行流。

## 目标

提供 `c.Use(...RunnerFunc)` 注册中间件；执行时按"洋葱/链式"顺序运行：
`mw1 → mw2 → ... → c.Func`，任一环节返回 error 即中止并向上传递。不破坏现有"无中间件"路径。

## 方案（采用线性链，必要时支持显式 Next）

1. 注册 API（`cmd.go`）：
   ```go
   func (c *Command) Use(handlers ...RunnerFunc) *Command {
       c.middles = append(c.middles, handlers...)
       return c
   }
   ```
2. 执行接入（`doExecute`）：把 `c.Func` 视为链的最后一环，统一通过链执行。
   - 简化方案（推荐）：线性执行——依次跑 `middles`，任一返回 error 即中止；全部通过后跑 `c.Func`。
     语义清晰、无递归。`middleIdx`/`Next()` 改为内部驱动（或保留 `Next()` 供高级用户在中间件内
     手动推进，gin 风格二选一，默认线性）。
   - 关键改动点：`doExecute` 中 `err = c.Func(c, fnArgs)` 替换为 `err = c.runChain(fnArgs)`，
     `runChain` 跑 `middles` 后再跑 `c.Func`。
3. 应用级中间件（可选子项）：`App.Use(...)`，在每个命令执行前统一注入（对齐维护者路线图
   “group options/controller”思路，但仅做中间件，范围更小）。
4. 复位：`Copy()` 已重置部分字段，需确认 `middles`/`middleIdx` 在复用/复制时的语义。

## 兼容性 / 风险

- 未调用 `Use` 时执行路径与现状一致（`middles` 为空 → 直接跑 `c.Func`）→ 低破坏。
- 风险低中：需处理中间件内 panic（沿用 `doExecute` 既有 recover）、错误事件 `OnCmdRunError` 触发时机。

## 测试

- 顺序：多个中间件按注册顺序执行，`c.Func` 最后执行。
- 中止：中间某个返回 error → 后续中间件与 `c.Func` 不执行，错误向上传递并触发 `OnCmdRunError`。
- 无中间件回归：现有命令执行不受影响。
- panic 恢复：中间件 panic 被 `doExecute` recover 成 error。

## 提交拆分

1. ✅ `feat(cmd): 命令级中间件 Use()`（commit 7cec911）—— `Command.Use` + `runWithMiddles`
   线性执行接入 `doExecute`，测试覆盖顺序/中止/无中间件回归/链式。
2. ✅ `feat(app): App.Use() 应用级中间件`（commit 66502f7）

> **遗留/发现**：
> - `doExecute` 里那段 `recover()` 是**内联调用(非 defer)**，实为 no-op，不能捕获 panic
>   （pre-existing）。本轮按现状保留、未做中间件 panic 恢复。如需真正的 panic→error，
>   应单独修该 recover 为 defer 形式（建议另起 fix 提交）。
> - 子命令未继承父命令中间件（仅作用于注册命令本身）；按需可后续支持继承。

## 体量预估

约 1~2 文件、60~100 行 + 测试。接近阈值，命令级（仅 cmd.go）可直接做；含 App 级需确认。
