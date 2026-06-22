# 功能实现计划：D1 结构体绑定 — 去 unsafe + 类型丰富度

> 状态：**全部完成**（D1.1-D1.5 + 文档）
> 范围：集中在 `gflag/parser.go`（`fromStructValue`）+ `gflag/util.go`（tag 键）+ 新增小适配器/helper。
> 依据：[../compare-with-others.zh-CN.md](../compare-with-others.zh-CN.md) 差距 4；[../TODO.md](../TODO.md) D1。
> 原则：**全部增量、向后兼容**，老 API 与既有标签行为不变。

---

## 现状

- `gflag/parser.go:436` `fromStructValue` 基础类型分支用 `unsafe.Pointer(fv.UnsafeAddr())` 取地址，
  再 `(*bool)(ptr)` 强转后调用 `BoolVar/IntVar/...`。其**上方** `flag.Value` 分支（:421-433）
  已用安全写法 `fv.Addr().Interface().(flag.Value)` —— 说明 unsafe 并非必需。
- 支持的字段类型窄：仅 `bool/int/int64/uint/uint64/float64/string` + 实现了 `flag.Value` 的字段；
  原生 `[]string/[]int`、`time.Duration`、`map[string]string` 不能从标签自动绑定（:444 default 报错）。
- 已具备的基础设施：
  - `FlagSet.DurationVar` / `CliOpts.DurationVar`（`opts.go:360`）—— Duration 绑定现成。
  - `CliOpts.Var(flag.Value)`（`opts.go:390`）—— 任意 flag.Value 绑定入口。
  - `cflag.Strings`(`type Strings []string`)/`Ints`(`[]int`)/`Booleans`(`[]bool`)：方法在指针接收者上，
    底层类型与原生 slice 一致 → `*[]string` **可转换**为 `*Strings`（Go 指针可转换规则）。
  - `cflag.KVString`：是 **struct**（包 `maputil.SMap`+`Sep`），与原生 `map` 底层不同 → 不能指针转换，需小适配器。
- `gflag/util.go:49` `namedTagKeys = {name,desc,required,default,shorts,short}` —— 无 `enum`。
- `CliOpt.Choices`（`opts.go:675`）仅用于**补全候选**，不参与校验；校验靠 `CliOpt.Validator`。

## 目标

1. **去 unsafe**：基础类型分支改用 `fv.Addr().Interface().(*T)`，删除 `unsafe` 依赖，行为不变。
2. **类型丰富度**：标签自动绑定新增支持
   - `[]string` / `[]int` / `[]bool`（指针转换到 `cflag.Strings/Ints/Booleans`，可重复 `--x a --x b`）
   - `time.Duration`（`DurationVar`，支持 `1h30m` 等）
   - `map[string]string`（小适配器，`--x k=v` 重复填充）
3. **enum 标签**：`enum:"a,b,c"` → 设 `opt.Choices`（补全）+ 成员校验（`Validator`）。
4. **（下一步，可选）泛型 API**：`gflag.Opt[T]/BindVar[T]` 类型安全入口，老 API 保留。

---

## 方案

### D1.1 去 unsafe（重构，零行为变化）

`parser.go` 基础类型分支：

```go
// 取代 unsafe.Pointer(fv.UnsafeAddr())
if !fv.CanAddr() { // 结构体经指针传入，字段恒可寻址；防御性判断
    return fmt.Errorf("field: %s - is not addressable for binding", name)
}
addr := fv.Addr().Interface()
switch ft.Kind() {
case reflect.Bool:    p.BoolVar(addr.(*bool), opt)
case reflect.Int:     p.IntVar(addr.(*int), opt)
case reflect.Int64:   p.Int64Var(addr.(*int64), opt)
case reflect.Uint:    p.UintVar(addr.(*uint), opt)
case reflect.Uint64:  p.Uint64Var(addr.(*uint64), opt)
case reflect.Float64: p.Float64Var(addr.(*float64), opt)
case reflect.String:  p.StrVar(addr.(*string), opt)
...
}
```

移除文件顶部 `unsafe` import。指针字段在 :386-388 已 `fv = fv.Elem()`（仍可寻址），不受影响。

### D1.2 slice + duration（在 kind switch 前/内新增分支）

**Duration**（必须在 `case reflect.Int64` 之前判断，因 `time.Duration` 的 Kind 即 Int64）：

```go
durType := reflect.TypeOf(time.Duration(0))
if ft == durType {
    p.DurationVar(addr.(*time.Duration), opt)
    continue
}
```

**slice**（用指针可转换，避免 unsafe；统一一个 helper）：

```go
// bindSliceField: []string->*Strings, []int->*Ints, []bool->*Booleans
func bindSlice(p *Parser, fv reflect.Value, opt *CliOpt) bool {
    var target reflect.Type
    switch fv.Type().Elem().Kind() {
    case reflect.String: target = reflect.TypeOf(Strings(nil))
    case reflect.Int:    target = reflect.TypeOf(Ints(nil))
    case reflect.Bool:   target = reflect.TypeOf(Booleans(nil))
    default: return false
    }
    // *[]T -> *cflag.XXX（底层一致，指针可转换）
    pv := fv.Addr().Convert(reflect.PointerTo(target)).Interface().(flag.Value)
    p.Var(pv, opt)
    return true
}
```

在 switch 的 `case reflect.Slice:` 调用 `bindSlice`，失败则 default 报错。

> 注：`[]int64/[]float64` 等暂不覆盖（cflag 无现成类型），需要再加可后续扩展。

### D1.3 map[string]string（小适配器）

cflag 无指针式 map 适配器，新增内部类型：

```go
// mapStrValue binds a flag.Value to a user's map[string]string field.
type mapStrValue struct {
    ref *map[string]string
    sep string // default "="
}
func (m *mapStrValue) String() string { /* join k=v, 顺序无关 */ }
func (m *mapStrValue) Get() any { return *m.ref }
func (m *mapStrValue) Set(s string) error {
    if *m.ref == nil { *m.ref = make(map[string]string) }
    k, v := strutil.SplitKV(s, m.sep)
    if k != "" { (*m.ref)[k] = v }
    return nil
}
func (m *mapStrValue) IsRepeatable() bool { return true }
```

绑定：`case reflect.Map`（仅 `map[string]string`）→
`p.Var(&mapStrValue{ref: addr.(*map[string]string), sep: "="}, opt)`。其余 map 类型 default 报错。

### D1.4 enum 标签

- `util.go:49` `namedTagKeys` 增加 `"enum"`；`TagRuleField` 读 `sf.Tag.Get("enum")`；
  `TagRuleSimple` 不支持（文档说明）。
- 在 opt 创建后、bind 前：

```go
if enum := mp["enum"]; enum != "" {
    opt.Choices = strutil.Split(enum, ",")        // 补全候选
    if opt.Validator == nil {                      // 成员校验(不覆盖用户自定义)
        opt.Validator = enumValidator(opt.Choices)
    }
}
```

```go
func enumValidator(choices []string) func(string) error {
    return func(val string) error {
        if val == "" || arrutil.StringsHas(choices, val) { return nil }
        return fmt.Errorf("value %q not in allowed: %v", val, choices)
    }
}
```

> 选用「native 字段 + Choices + Validator」而非强制 `EnumString`，保持字段类型对用户透明。

### D1.5 泛型 API（下一步，可选，单独子阶段）

新增类型安全入口，内部按 `any`→具体绑定分发；老 API 不动：

```go
func Opt[T any](co *CliOpts, p *T, name, shorts string, def T, desc string, fns ...CliOptFn)
```

仅在 D1.1-D1.4 落地并验证后再评估，避免一次摊太大。

---

## 风险 / 兼容

- **去 unsafe**：纯内部重构，结构体经指针传入字段恒可寻址，回归测试保证行为一致。
- **新增类型**：均为 default 分支之前的新增 case，对既有「用户已用 `gflag.Strings` 等字段」的写法无影响
  （仍走 flag.Value 分支）。
- **enum**：仅在 tag 出现 `enum` 时生效；不覆盖用户已设的 `Validator`。
- 无任何对外 API 签名变更（D1.5 泛型为新增函数）。

## 测试（`github.com/gookit/goutil/x/assert`）

- **回归**：现有 `FromStruct`/三种 rule 的绑定测试全绿（去 unsafe 不改行为）。
- **slice**：`[]string`/`[]int`/`[]bool` 字段，`--names a --names b` → `["a","b"]` 等。
- **duration**：`time.Duration` 字段，`--ttl 1h30m` → `90m`。
- **map**：`map[string]string` 字段，`--meta k1=v1 --meta k2=v2` → 两键。
- **enum**：`enum:"a,b,c"` 合法值通过、非法值报错；`opt.Choices` 已填充（补全可用）。
- **指针字段**：`*time.Duration` / `*[]string` 等可寻址路径。

## 提交拆分（按 R002，每子阶段一提交）

1. ✅ `refactor(gflag): struct binding 去除 unsafe, 改用 Addr().Interface()`（commit 7830e0e）
2. ✅ `feat(gflag): struct tag 支持 slice([]string/[]int/[]bool) 与 time.Duration`（commit a38d5d5）
3. ✅ `feat(gflag): struct tag 支持 map[string]string(mapStrValue 适配器)`（commit 633e0ff）
4. ✅ `feat(gflag): struct tag 支持 enum 键(Choices + 成员校验)`（commit cd33f51）
5. ✅ `feat(gflag): 泛型选项绑定 API Opt[T]/BindVar[T]`（commit 99ffdfa）
6. ✅ `docs: 更新 README/CHANGELOG 结构体标签新支持类型与 enum`

> 预计每子阶段 < 60 行业务代码、1-2 个文件，满足「先确认再实施」阈值后可逐个推进。
