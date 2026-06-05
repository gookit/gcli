# 功能实现计划：B4+B5 POSIX 短选项增强

> 状态：**待评审**
> 范围：合并 strictFormatArgs（B4）与 EnhanceShort（B5），提供标准 POSIX 短选项解析。

## 现状（含一处文档与实现不符）

- `gflag/flags.go:304` `parseOne()` 只解析**单个**短名（`:341 f.shorts[name]`）。输入 `-ab` 时
  `name="ab"`，不在 shorts 映射 → 当作未知长选项处理。**即组合短选项 `-ab` 默认并不会拆成 `-a -b`。**
- **文档不符**：`AGENTS.md` / README 写着「POSIX 风格的短选项合并（`-a -b` = `-ab`）」，但默认不成立。
- `util.go:104` `strictFormatArgs`：仅在 `strictMode=true`（默认 false）时预处理，且 `-ab` 拆分、
  `-Oval` 等模式都还是 `// TODO`。
- `gflag/gflag.go:72` `EnhanceShort uint8` 配置位已留好（0 none / 1 multi-bool / 2 attached-value），未实现。

## 目标

让以下输入按 POSIX 习惯正确解析（由 `EnhanceShort` 等级控制，默认 0 保持兼容）：

- level ≥ 1：`-aux` → `-a -u -x`（**仅当 a/u/x 均为 bool 短选项**时才拆，避免误伤）
- level ≥ 2：`-O stdout` 的紧贴写法 `-Ostdout` → `-O stdout`（首字符是取值短选项）

## 方案（在解析前做一次"短选项规范化"预处理）

1. 在 `gflag` 增加 `expandShortArgs(args, shorts, optTypes, level)` 纯函数（新文件
   `gflag/shorts.go`），逻辑：
   - 仅处理单 `-` 开头、长度 > 2、不含 `=` 的 token。
   - 逐字符查 `shorts` 映射 + 选项类型：
     - 全是 bool 短选项 → 拆成多个 `-x`（level≥1）。
     - 首字符是取值短选项且后面有残余 → `-O` + `rest` 作为值（level≥2）。
     - 不满足 → 原样保留（不猜测）。
2. 接入点：`Parser.Parse` 中 `p.fSet.Parse(args)` 之前调用，传入 `p.cfg.EnhanceShort`、`co.shorts`、
   各 `opt.flagType`。（不改 forked `FlagSet` 内部，风险可控。）
3. 收敛 `strictFormatArgs`：把 `--a→-a`、`---name→--name` 等规范化也并入同一预处理；
   `gcli` 的 `strictMode` 改为"开启即把 EnhanceShort 设为 ≥1"，消除两套并存。
4. 更新 `AGENTS.md`/README：明确组合短选项需 `EnhanceShort`（或 strict 模式）开启，纠正现有不实表述。

## 兼容性 / 风险

- 默认 `EnhanceShort=0` → 行为完全不变；仅显式开启才改变解析 → 低破坏。
- 风险中：解析是核心路径，需覆盖大量边界（`=` 取值、`--` 终止符、未知短名、bool 与取值混合）。
- "仅当全为 bool 才拆"是关键安全约束，避免把 `-name`（某长选项简写习惯）误拆。

## 测试

- 矩阵用例：`-aux`(全 bool)、`-ab`(含取值→不拆)、`-Ostdout`(level2)、`-O=stdout`、`--a`、`---name`、
  `--`、未知短名、与长选项混合。
- level 0/1/2 三档各跑一遍，确认 0 档与现状一致。

## 提交拆分

1. `feat(gflag): expandShortArgs 短选项规范化(EnhanceShort 1/2)`
2. `refactor(gflag): strictFormatArgs 收敛进短选项规范化`
3. `docs: 修正短选项合并的说明(需 EnhanceShort/strict)`

## 体量预估

约 2~3 文件、120~180 行 + 较多表驱动测试。**超阈值，实施前确认。**
