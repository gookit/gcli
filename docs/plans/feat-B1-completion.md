# 功能实现计划：B1 自动补全（内置补全选项）

> 状态：**待评审**
> 范围：把补全能力做成应用的两个内置全局选项，无需额外注册 `genac` 命令：
>   - `--in-completion`：运行期**动态**补全（隐藏内部选项，供 shell 脚本回调取候选）。
>   - `--gen-completion <shell>`：直接**静态**生成补全脚本并退出（复用现有 genac 逻辑）。

## 现状

- **已实现**：`builtin/gen_auto_complete.go` 的 `GenAutoComplete()`（命令名 `genac`）能按 bash/zsh
  模板生成**静态**补全脚本（`buildForBashShell` / `buildForZshShell`），`resource/auto-completion/`
  有示例脚本。
- **缺失（半成品）**：运行期**动态**补全：
  - `help.go:186` `showAutoCompletion(_ []string)` 为空实现 `// TODO ...`
  - `app.go:184` 隐藏选项 `--in-completion` 绑定 `gOpts.inCompletion`
  - `app.go:~330` 当 `inCompletion` 为真时调用 `app.showAutoCompletion(app.args)`（目前空转）
- **预留**：`gcli.go` 已加 `genCompletion string` 字段（TODO），用于 `--gen-completion bash|zsh|pwsh`
  直接静态生成；目前静态生成仅能通过注册 `builtin.GenAutoComplete()`（`genac` 命令）使用。
- 维护者路线图相关项：`docs/Changelog-TODO.md` 的 “prompt completion by readline”（交互式，属 B3，不在此）。

静态脚本的弊端：命令/选项变化后需重新生成脚本。动态补全让脚本只负责"回调二进制取候选"，零维护。

## 目标

实现 `app --in-completion <当前命令行片段>` → 输出候选列表（每行一个），供 bash `compgen` /
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
   `{{.BinName}} --in-completion "${COMP_WORDS[@]}"` 并把输出喂给 `compgen -W` / zsh `compadd`。
   保留旧静态模式作为 `--static` 备选（兼容）。
4. 为选项提供候选来源：在 `CliOpt` 增加可选 `Choices []string`（或复用 `EnumString` 的枚举），
   动态补全据此给值候选。（小增量，可作为子提交。）

## 静态内置生成 `--gen-completion <shell>`（免注册 genac）

把现有 `genac` 的静态脚本生成能力提升为**应用内置全局选项**，用户零注册即可用。

1. 绑定：`bindAppOpts` 增加 `fs.StrVar(&app.opts.genCompletion, ...)`，选项名 `--gen-completion`，
   取值 `bash|zsh|pwsh`（**非隐藏**，进帮助）。
2. 处理：`parseAppOpts` 中若 `genCompletion != ""` → 调用静态生成并退出（返回 OK），
   与现有 `inCompletion` / `ShowVersion` 的"命中即退出"分支并列。
3. 复用：把 `builtin/gen_auto_complete.go` 的 `buildForBashShell`/`buildForZshShell` 与模板
   抽取/下沉为可被主流程复用的内部函数（避免 app 反向依赖 builtin）；`genac` 命令保留为薄封装。
4. pwsh：✅ 已补 **动态(瘦)** PowerShell 模板(`Register-ArgumentCompleter` 委托 `--in-completion`，
   文件名用 `.ps1`)。`--gen-completion pwsh` / `genac --shell pwsh` 可用；静态嵌入 pwsh 暂不支持
   (`GenStaticCompletionScript(pwsh)` 返回 error)。需人工在 PowerShell 7+ 实测 Tab。

## 兼容性 / 风险

- 新增能力，不改既有 `genac` 静态行为（默认仍可用）→ 低破坏。
- 风险中：shell 引号/分词差异、`COMP_CWORD` 边界。**候选计算本身可纯单测**（给定片段→期望候选），
  shell 胶水层另行手测。

## 测试

- 单测候选计算：顶层命令、子命令、选项名、枚举值、别名、`help` 各场景。
- 手测：`source` 生成脚本后在 bash/zsh 实际 Tab。

## 提交拆分

1. ✅ `feat(completion): 内置 --gen-completion 静态生成(抽取 genac 生成逻辑复用)` — 已完成(commit 3a76280)
2. ✅ `feat(completion): 运行期动态补全候选计算 + showAutoCompletion(--in-completion)` — 已完成(commit 88d70e9)
3. ✅ `feat(completion): genac/动态脚本委托 --in-completion`（生成回调二进制的瘦脚本，零维护）— 已完成
   - 委托式动态为默认: `GenCompletionScript` 产瘦脚本(回调 `bin --in-completion`)；
     原嵌入式静态改名导出 `GenStaticCompletionScript`，作为 opt-in。
   - 全局 `--gen-completion <shell>` 走默认瘦脚本；`genac` 新增 `--static/-S` 切换为静态。
   - 静默模式: 新增 `App.completionMode` + `hasMetaFlag`，补全/生成请求时抑制用户
     init / opts-parsed 生命周期钩子(保留 bindOpts 钩子)，保证 stdout 只剩候选/脚本。
4. `feat(gflag): CliOpt.Choices 提供选项值候选`（可选）

> 进度：静态 `--gen-completion`、动态候选计算 `--in-completion`（`resolveCompletion` +
> `showAutoCompletion`）、脚本委托动态 + 静默模式（子阶段 3）均已落地并测试。
> 剩余：选项值候选 `Choices`（子阶段 4，可选）。
>
> 已知取舍（待后续按需加强）：补全不含隐藏命令的过滤（与静态模板一致，仅过滤了隐藏选项）；
> 选项值词会停止下钻；选项**值**补全未做。

## 体量预估

约 3~4 文件、200~300 行 + 测试。**超过 3 文件/100 行阈值，实施前需确认。**
