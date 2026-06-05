# 功能实现计划：B1 自动补全（运行期动态补全）

> 状态：**待评审**
> 范围：补全能力的"运行期动态"部分。静态脚本生成（`genac`）已存在，不在本计划重做。

## 现状

- **已实现**：`builtin/gen_auto_complete.go` 的 `GenAutoComplete()`（命令名 `genac`）能按 bash/zsh
  模板生成**静态**补全脚本（`buildForBashShell` / `buildForZshShell`），`resource/auto-completion/`
  有示例脚本。
- **缺失（半成品）**：运行期**动态**补全：
  - `help.go:186` `showAutoCompletion(_ []string)` 为空实现 `// TODO ...`
  - `app.go:184` 隐藏选项 `--gen-completion` 绑定 `gOpts.inCompletion`
  - `app.go:~330` 当 `inCompletion` 为真时调用 `app.showAutoCompletion(app.args)`（目前空转）
- 维护者路线图相关项：`resource/Changelog-TODO.md` 的 “prompt completion by readline”（交互式，属 B3，不在此）。

静态脚本的弊端：命令/选项变化后需重新生成脚本。动态补全让脚本只负责"回调二进制取候选"，零维护。

## 目标

实现 `app --gen-completion <当前命令行片段>` → 输出候选列表（每行一个），供 bash `compgen` /
zsh `_values` 消费；并让 `genac` 生成的脚本改为**委托**给该动态模式。

## 方案

1. 新增 `completion.go`（避免 help.go 膨胀），实现候选计算：
   - 输入：已解析的 `app.args`（shell 传入的 `COMP_WORDS` 去掉 bin 名）。
   - 上下文判定（按最后一个 token）：
     - 顶层：未输入命令 → 候选 = 顶层命令名 + 别名（+ `help`）。
     - 已定位到命令/子命令：候选 = 其子命令 + 选项名（`--xxx`/`-x`）。
     - 正在输入选项值：若该选项是 `EnumString`/有候选集 → 输出枚举值。
   - 复用现有能力：`app.CommandsByGroup` 之外的 `CmdNameMap`、`cmdAliases`、各命令 `Flags.Opts()`。
2. 实现 `showAutoCompletion(args)`：调用候选计算并逐行打印（纯 stdout，无颜色）。
3. 更新 `gen_auto_complete.go` 的 bash/zsh 模板：由"硬编码命令列表"改为调用
   `{{.BinName}} --gen-completion "${COMP_WORDS[@]}"` 并把输出喂给 `compgen -W` / zsh `compadd`。
   保留旧静态模式作为 `--static` 备选（兼容）。
4. 为选项提供候选来源：在 `CliOpt` 增加可选 `Choices []string`（或复用 `EnumString` 的枚举），
   动态补全据此给值候选。（小增量，可作为子提交。）

## 兼容性 / 风险

- 新增能力，不改既有 `genac` 静态行为（默认仍可用）→ 低破坏。
- 风险中：shell 引号/分词差异、`COMP_CWORD` 边界。**候选计算本身可纯单测**（给定片段→期望候选），
  shell 胶水层另行手测。

## 测试

- 单测候选计算：顶层命令、子命令、选项名、枚举值、别名、`help` 各场景。
- 手测：`source` 生成脚本后在 bash/zsh 实际 Tab。

## 提交拆分

1. `feat(completion): 运行期动态补全候选计算 + showAutoCompletion`
2. `feat(completion): genac 模板委托动态补全(保留 --static)`
3. `feat(gflag): CliOpt.Choices 提供选项值候选`（可选）

## 体量预估

约 3~4 文件、200~300 行 + 测试。**超过 3 文件/100 行阈值，实施前需确认。**
