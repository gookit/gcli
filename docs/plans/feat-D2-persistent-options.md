# 功能实现计划：D2 共享选项继承模型（SharedOpts，≈ cobra PersistentFlags）

> 状态：**D2.1 + D2.5 已完成**（提交拆分 1-5 + 复核修复）；D2.6(gOpts per-App) 待评估
> 范围：`cmd.go`（共享选项存储/合并/分发）+ 新增 gflag 合并基元 + 可选的 `gcli.go`（gOpts per-App）。
> 依据：[../compare-with-others.zh-CN.md](../compare-with-others.zh-CN.md) 差距 5；[../TODO.md](../TODO.md) D2。
> 对标：cobra 的 `PersistentFlags` / `LocalFlags` / `InheritedFlags` 三层模型。

## 修订记录

| 日期 | 变更 |
|---|---|
| 2026-06-21 | 初版：三层选项模型(global/shared/local) + `InheritOptsFrom` 合并基元 + 分发合并方案 |
| 2026-06-21 | 增加「现状 vs 共享选项」输入对照表，澄清 cobra `PersistentFlags` 的含义（父选项下放到子命令树，可写在子命令段） |
| 2026-06-21 | API 命名定为 **`SharedOpts()`**（原拟 `PersistentFlags`，改为更短、贴 gcli `Opts` 词汇） |

---

## 现状

- **只有两层选项**：
  1. **进程级全局选项** `gOpts`（`gcli.go:79` 包级单例，`newGlobalOpts()`），通过
     `bindingOpts(fs)`（`gcli.go:211`）绑定 `--help/-h`、`--version/-V`。App 复用它
     （`app.go:128` `app.opts = gOpts`，`app.go:203` 绑到 `app.fs`）。
  2. **命令局部选项**：每个命令在 `Config` 里用 `c.BoolOpt/StrOpt/FromStruct/...` 绑到 `c.Flags`。
- **缺中间层**：父命令定义、其所有子命令自动继承的「共享选项」不存在。多级命令里要让
  `git --git-dir <path> <subcmd> ...` 这种「父级选项对所有子命令生效」，只能在每个子命令重复声明。
- **分发链路**（`cmd.go:409` `innerDispatch`）：`parseOptions` → 找子命令 →
  `sub.innerDispatch(args[1:])` 递归。解析在子命令名处停止（首个位置参数）。
- **父子链路已就绪**：`AddCommand`（`cmd.go:213`）设 `sub.parent = c`、`sub.app = c.app`，
  `c.Root()/c.Parent()` 可用 → 可沿祖先链向上walk。
- **`gOpts` 是进程单例**：`CHANGELOG v3.4.0` 已记此坑（多 App 实例共享全局选项）。
  > 历史注意：v3.4.0 之所以从「每 App 一份拷贝」改回「复用包级 gOpts」，是因为旧的 per-app 拷贝里
  > `--verbose` 绑的是副本、logger 从没读到它而失效。per-App 改造**不能重新引入这个 bug**。

## 目标

确立**三层选项模型**：

| 层 | 定义者 | 作用范围 |
|---|---|---|
| 全局(global) | App（现 `gOpts`） | 所有命令 |
| 共享(shared) | 某命令 | 该命令**及其所有子孙命令** |
| 局部(local) | 某命令 | 仅该命令 |

> API 命名采用 **`SharedOpts()`**（比 `PersistentFlags` 短、贴 gcli 的 `Opts` 词汇）；
> 下文「共享选项」≈ cobra 的 persistent flags。

#### 动机：现状 vs 共享选项（直观对照）

命令树 `app remote(父) → add(子)`，`--verbose` 定义在 `remote` 上。

| 输入 | 现状(remote 局部选项) | 共享选项后 |
|---|---|---|
| `app remote --verbose add` | ✅ 写在子命令名**之前**，由 remote 解析 | ✅ |
| `app remote add --verbose` | ❌ remote 解析到 `add` 即停，`--verbose` 交给 add；**add 不认识** → 报错 | ✅ add **继承**了 `--verbose`，能解析 |
| `app remote add arg --verbose`（配合重排） | ❌ 同上 | ✅ 写在 arguments 后也能解析 |

本质：**把父命令的某个选项「下放」到整棵子命令树**，使其在子命令(叶子)那一段里也能写、能解析；
父子写的是**同一个变量**（共享同一 `flag.Value`/`*ptr`）。对标 cobra 的 `PersistentFlags`，
如 `git --git-dir=X status` 里 `--git-dir` 能在子命令处使用。

> 它**只解决父 → 子方向**：不会让子命令自己的局部选项能写在子命令名之前。

1. **D2.1**（核心）：新增共享选项中间层，父命令的共享选项被子孙命令继承，且与 args 重排契合
   （可写在叶子段任意位置）。
2. **D2.2**（可选、有风险）：`gOpts` 由进程单例改为 per-App 实例，修多 App/并发共享。
3. **D2.3**：三层语义文档 + help 渲染（继承选项分组显示）+ 专项测试。

---

## 方案

### D2.1 共享选项 tier（核心）

#### API（复用现有全部绑定方法）

`Command.SharedOpts() *gflag.Flags`：惰性创建并返回命令专属的共享选项持有器 `c.sharedFs`。
用户在它上面像平时一样绑定 —— 自动获得 `BoolOpt/StrOpt/IntOpt/DurationOpt/Opt[T]/FromStruct/...`
全部能力：

```go
var gitDir string
top.SharedOpts().StrOpt(&gitDir, "git-dir", "", "", "the git work dir")
// 或结构体：top.SharedOpts().FromStruct(&sharedOpts)
```

> 选 `SharedOpts() *Flags`（对标 cobra `cmd.PersistentFlags()`）而非 `CliOpt.Shared` 标记，
> 因为前者零成本复用所有绑定方法，且 local/shared 清晰分离。`c.sharedFs` 仅作**定义来源**，
> 自身永不单独 Parse。

#### 存储

`Command` 增加字段 `sharedFs *gflag.Flags`（lazy，`SharedOpts()` 内首次创建）。

#### 合并基元（gflag 新增）

在 gflag 增加「把一个 parser 的选项按 `flag.Value` 重注册到另一个 parser」的能力。
因为 `flag.Value` 是包装了用户 ptr 的接口，**同一个 Value 注册进多个 FlagSet，各 FlagSet 解析时
都写回同一个 ptr**：

```go
// InheritOptsFrom 将 src 的选项重注册进 p(共享同一 flag.Value/ptr)。
// 已存在的同名选项(局部优先)跳过。
func (p *Parser) InheritOptsFrom(src *Parser) {
    for name, opt := range src.Opts() {
        if p.HasOption(name) { // 子命令局部同名选项优先
            continue
        }
        // 复制元数据, 复用同一 flag.Value 重注册到 p.fSet
        p.Var(opt.flag.Value, &CliOpt{
            Name: opt.Name, Shorts: opt.Shorts, Desc: opt.Desc,
            Required: opt.Required, Validator: opt.Validator, Choices: opt.Choices,
        })
    }
}
```

（`opt.flag.Value` 在 gflag 包内可访问；`Var` 已存在 `opts.go:390`。）

#### 合并点 & 幂等

在 `cmd.go` `parseOptions`（`:493`）中、`c.Parse(args)` 之前合并：

```go
if !c.sharedMerged {
    // 沿祖先链(含自身)收集共享选项, 从根到叶顺序合并进 c.Flags
    for _, anc := range c.ancestorsWithSelf() { // [root,...,parent,self]
        if anc.sharedFs != nil {
            c.Flags.InheritOptsFrom(&anc.sharedFs.Flags)
        }
    }
    c.sharedMerged = true
}
```

- 自身的 `sharedFs` 也并入 `c.Flags` → 命令**自己**也能识别它定义的共享选项
  （写在子命令名之前的场景）。
- 幂等 `sharedMerged bool`，避免重复注册。
- **与 reorder 契合**：合并后共享选项已进入叶子的 flag set，`optMeta` 能识别其取值性，
  写在叶子 arguments 之后也会被正确重排+解析。

#### 多级解析走查

`app top sub --git-dir /x arg`（`--git-dir` 是 top 的共享选项）：
- app.fs 解析全局(重排关闭) → 停在 `top` → dispatch top `[sub,--git-dir,/x,arg]`
- top.parseOptions：合并 top.sharedFs 进 top.Flags；top.Parse 停在 `sub` → dispatch sub `[--git-dir,/x,arg]`
- sub.parseOptions：合并祖先(top)的 sharedFs 进 sub.Flags → sub.Parse 识别 `--git-dir`=/x，写回共享 ptr，
  `arg` 余下 → doExecute。✅ 共享选项在叶子段任意位置可用。

#### help 渲染

合并后继承选项会出现在命令 flag set 中，help 自动显示。为避免与局部选项混淆，
可借助已有的 Category 机制把继承选项归到 `Global/Inherited Options` 分组（对标 cobra 的
"Global Flags"）。**此项作为子任务，先保证功能，再优化展示。**

### D2.2 gOpts per-App（可选、有风险）

- 现状：`gOpts` 包级单例，`app.opts = gOpts`。
- 方案：`NewApp` 时 `app.opts = newGlobalOpts()`（每 App 独立）；`bindingOpts` 绑到 app 自己的实例。
- **风险（重点）**：勿重蹈 v3.4.0 覆辙 —— verbose/日志级别当年因 per-app 副本不被 logger 读取而失效。
  - 处理：日志级别继续以 `GCLI_VERBOSE` 环境变量 / `gcli.SetVerbose()`（写包级 logger）为准，
    per-App 的 `GlobalOpts.Verbose` 仅作展示/读取，不作为 logger 的唯一来源；或在 App.Run 时
    把 `app.opts.Verbose` 显式注入 logger。
  - 单独提交 + 多 App 并发测试，确认 verbose/strict/completion 互不串扰。
- 若风险评估偏高，可**暂缓**，仅交付 D2.1/D2.3。

### D2.3 文档 + 测试

- README/CHANGELOG 增「三层选项模型」说明与示例。
- 测试见下。

---

## 风险 / 兼容

- **D2.1 增量、无破坏**：不碰现有 local/global 绑定路径；命令未用 `SharedOpts()` 时
  `sharedFs==nil`，合并是 no-op，行为完全不变。
- **同名冲突**：局部同名选项优先（合并时 `HasOption` 跳过），语义明确且安全。
- **Required 共享选项**：在「实际执行命令」处校验（合并后随 `validateAll` 生效），需在测试中固化。
- **D2.2 有破坏风险**：改单例为 per-App 可能影响依赖 `gcli.GOpts()` 全局读取的代码；
  单独评估、独立提交、可暂缓。

## 测试（`github.com/gookit/goutil/x/assert`）

- **继承**：top 定义共享 `--git-dir`，`app top sub --git-dir /x` → sub 内读到 `/x`。
- **任意位置**：`app top sub arg --git-dir /x`（重排）→ 同样生效，arg 仍为参数。
- **多级继承**：三层 `app a b c`，a 的共享选项在 c 可用（沿祖先链合并）。
- **局部优先**：子命令定义同名局部选项时，覆盖继承（互不写串）。
- **自身可用**：`app top --git-dir /x sub` → top 段即识别共享选项。
- **Required 共享**：缺失时在执行命令处报错。
- **未用共享**：现有用例全绿（sharedFs==nil，零行为变化）。
- **(D2.2)** 两个 App 实例的 verbose/strict 互不影响；并发 `-shuffle` 稳定。

## 提交拆分（按 R002）

1. ✅ `feat(gflag): 增加 Parser.InheritOptsFrom 选项合并基元`（commit 323219f）
2. ✅ `feat(gcli): Command.SharedOpts 共享选项定义 + 分发时合并(含幂等)`（commit 77f2391）
3. ✅ `test(gcli): 共享选项继承/任意位置/局部优先/Required 用例`（commit bf4ef0c）
4. ✅ `docs: 三层选项模型说明(README/CHANGELOG)`（commit a42019e）
   - ✅ 复核修复：共享 Required 改用类型感知 `CliOpt.IsEmpty()`（commit 87951f0）
5. ✅ `feat(gcli): help 渲染继承选项分组(Inherited Options)`（commit b704f29）
6. ⏳（可选/独立评估）`refactor(gcli): gOpts 改 per-App 实例 + 多 App 测试` — 风险最高，需先确认

> 1-5（核心 + 文档 + help 分组）已落地并验证；6 风险最高，需先确认再做。
> 说明：内置 `help` 命令仅支持单级（help.go 既有 TODO），故 `help top sub` 暂不可用；
> `top sub -h` 正常。help 分组对两条 ShowHelp 路径均生效（ShowHelp 开头补了幂等合并）。

### 实施备注（复核发现）

- **Required 共享选项校验时机**：因「自身也并入」会让必填共享选项进入中间祖先命令的 flag set，
  而共享值常写在叶子段——祖先解析时尚未见到取值会误报必填、中止分发。修正：合并进 `c.Flags` 的
  **继承副本清除 Required**，必填校验延后到 `doExecute`（执行命令）由 `validateSharedRequired` 沿
  祖先链统一判定（用类型感知 `CliOpt.IsEmpty()`）。与 cobra persistent-flag 语义一致。
- **短名冲突**：`InheritOptsFrom` 在短名与目标已有名/短名冲突时**整体跳过**该选项继承，避免 panic。
