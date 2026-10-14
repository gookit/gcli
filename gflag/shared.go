package gflag

// InheritOptsFrom 将 src 解析器的选项「重注册」进当前解析器 p，用于实现共享/继承选项。
//
// 关键点：选项底层的 flag.Value 包装了用户的 *ptr，把同一个 flag.Value 注册进多个
// FlagSet 后，任意 FlagSet 解析时都写回同一个 ptr。因此父命令定义、子命令继承的选项
// 父子读写的是同一个变量。
//
// 合并规则：
//   - 本地已存在同名选项（局部优先）→ 跳过；
//   - 选项尚未绑定（opt.flag == nil）→ 跳过（理论上绑定后不会出现）；
//   - 短名与 p 已有的选项名/短名冲突时，安全起见整体跳过该选项的继承，避免 Var 内部 panic。
//
// 元数据完整复制（包括 DefVal/defEnvVar/Hidden/flagType 与 Collector/Question/Handler）：
// 直接调用本方法的解析器改完就能用，行为与源选项一致。gcli 的命令分发另有处理：
// 值与表单一时机相关的触发类字段（Required/Validator/Collector/Question/Handler）
// 在合并进命令段时会被清除，延后到实际执行命令时统一处理，见 Command.prepareSharedOpts。
//
// 可选 category：传入时，继承进来的选项会被归入该 help 分组（仅影响帮助展示，不影响解析）。
func (p *Parser) InheritOptsFrom(src *Parser, category ...string) {
	if src == nil {
		return
	}

	var cat string
	if len(category) > 0 {
		cat = category[0]
	}

	for name, opt := range src.Opts() {
		// 子命令局部同名选项优先, 跳过继承
		if p.HasOption(name) {
			continue
		}
		// 选项未绑定 flag.Value, 无法复用, 跳过
		if opt.flag == nil {
			continue
		}
		// 短名冲突检测: 任一短名已被 p 用作选项名或短名时, 跳过该选项, 避免 Var 内部 panic
		if p.shortsConflict(opt.Shorts) {
			continue
		}

		// 复制元数据, 复用同一 flag.Value 重注册到 p.fSet, 实现父子共享同一 ptr
		copt := &CliOpt{
			Name:      opt.Name,
			Shorts:    opt.Shorts,
			Desc:      opt.Desc,
			Required:  opt.Required,
			Hidden:    opt.Hidden,
			Collector: opt.Collector,
			Validator: opt.Validator,
			Handler:   opt.Handler,
			Question:  opt.Question,
			Choices:   opt.Choices,
			Category:  cat, // 继承选项归入指定 help 分组
			// 默认值元数据必须一起复制: 复用的 flag.Value 已带解析后的 ENV 值(即 Flag.DefValue),
			// 元数据丢失后 help 会把它当默认值渲染, 泄露 `${XXX}` 的真实取值。
			// DefVal 已是解析结果, 不会在 initCheck 中二次解析; defEnvVar 供 help/docs 展示原变量名。
			DefVal:    opt.DefVal,
			defEnvVar: opt.defEnvVar,
		}
		p.Var(opt.flag.Value, copt)

		// p.Var 按自定义 var 注册, 会把类型写成 var, 这里恢复源选项的类型(影响 help/completion)
		copt.flagType = opt.flagType
	}
}

// shortsConflict 检测给定的短名集合中是否有任意一个已被 p 用作选项名或已注册短名。
func (p *Parser) shortsConflict(shorts []string) bool {
	for _, short := range shorts {
		if p.HasOption(short) || p.IsShortName(short) {
			return true
		}
	}
	return false
}
