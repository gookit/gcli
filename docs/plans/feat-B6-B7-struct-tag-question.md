# 功能实现计划：B6 结构体标签增强 + B7 声明式交互收集

> 状态：B6 **待评审** / B7 **已完成**
> 范围：两个独立小特性，可分别实施。

---

## B6. TagRuleField 标签规则 + 匿名字段支持

### 现状

- `gflag/gflag.go:33` `TagRuleField` 常量已定义但 `// TODO`：期望"用字段名做选项名，
  `flag:"name,n" desc:"..." required:"true" default:"0"` 这种**分键** tag 风格"。
- `gflag/parser.go:328` `FromStruct` 仅实现 `TagRuleNamed` 与 `TagRuleSimple` 两种；
  传 `TagRuleField` 会走到 `errTagRuleType`。
- `gflag/parser.go:301` `// TODO support anonymous field by sf.Anonymous`：匿名嵌套 struct 不展开。

### 目标

1. 实现 `TagRuleField`：以字段名（SnakeCase）为选项名，从**独立 tag 键**读取元数据：
   `flag:"<shorts>" desc:"..." default:"..." required:"true"`（或约定：`flag` 仅放 name/shorts，
   其余用 `desc`/`default`/`required` 独立 tag）。
2. 支持匿名字段：`sf.Anonymous` 时递归解析其内部字段（实现选项分组/复用的基础）。

### 方案

- `FromStruct` 增加 `TagRuleField` 分支：用 `sf.Tag.Get("desc"/"default"/"required")` 组装 `mp`，
  `optName` 默认 `SnakeCase(sf.Name)`，`flag` tag 仅解析 shorts（可空）。
- 匿名字段：在字段循环里，若 `sf.Anonymous && ft.Kind()==Struct`，递归调用一个内部
  `fromStructValue(v.Field(i))`，把内层字段并入同一 parser。注意指针匿名字段与导出性判断。

### 风险 / 测试

- 低中风险，纯绑定期逻辑。测试：三种 rule 的等价绑定、匿名字段展开、指针字段、未导出字段跳过。

### 提交拆分

1. `feat(gflag): FromStruct 支持 TagRuleField 标签规则`
2. `feat(gflag): FromStruct 支持匿名嵌套字段展开`

---

## B7. 声明式交互收集 `CliOpt.Question`

### 现状

- 命令式 `Collector` **已实现**：`gflag/opts.go:656` 字段 + `:608 WithCollector` + `:787` 在
  `Validate` 中"值为空则调用 Collector 取值"。
- 声明式入口缺失：`gflag/opts.go:664` `// Question string` 被注释，未启用。

### 目标

提供 `CliOpt.Question`（+ `WithQuestion`）：当 required/空值时，自动用该问题向用户提问收集输入，
无需手写 Collector。本质是"内置一个基于 Question 的默认 Collector"。

### 方案

- 启用字段：`Question string`（opts.go），新增 `WithQuestion(q string) CliOptFn`。
- 在 `Validate` 的空值分支：若 `Collector == nil && Question != ""`，用内置交互读取
  （复用 `gookit/goutil` 的交互或 `cliui`，输出问题、读一行、trim）。Collector 优先级高于 Question。
- 复用现有空值/required 判定逻辑（`valIsEmpty`），不改变其它路径。

### 风险 / 测试

- 低风险。交互读取需可注入输入源以便测试（用 `goutil` 的可替换 stdin，或抽一个包级 reader 变量）。
- 测试：空值触发提问→设值成功；Collector 存在时忽略 Question；非空值不提问。

### 提交拆分

1. ✅ `feat(gflag): CliOpt.Question 声明式交互收集(内置默认 Collector)`

> 实现要点（已落地）：启用 `Question` 字段 + `WithQuestion()`；`Validate` 空值分支统一为
> "Collector 优先，否则用基于 `cliutil.ReadLine(Question)` 的内置默认 collector"；收集后重算
> `valEmpty`，顺带修复"Required+收集成功后仍误报 required"的潜在问题。测试覆盖三态。

---

## 体量预估

- B6：约 1 文件（parser.go）、80~120 行 + 测试。
- B7：约 1 文件（opts.go）、40~60 行 + 测试。
- 各自接近/低于阈值，可独立实施；B6 含匿名字段递归改动稍大，建议分两个提交。
