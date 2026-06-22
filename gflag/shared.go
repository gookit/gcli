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
func (p *Parser) InheritOptsFrom(src *Parser) {
	if src == nil {
		return
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
		p.Var(opt.flag.Value, &CliOpt{
			Name:      opt.Name,
			Shorts:    opt.Shorts,
			Desc:      opt.Desc,
			Required:  opt.Required,
			Validator: opt.Validator,
			Choices:   opt.Choices,
		})
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
