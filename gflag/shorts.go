package gflag

import "strings"

// isBoolShort reports whether the given short name maps to a registered bool option.
func (f *FlagSet) isBoolShort(short string) bool {
	name, ok := f.shorts[short]
	if !ok {
		return false
	}
	flg := f.formal[name]
	if flg == nil {
		return false
	}
	fv, ok := flg.Value.(boolFlag)
	return ok && fv.IsBoolFlag()
}

// expandShortArgs 对 args 做 POSIX 短选项规范化预处理(由 EnhanceShort 等级控制)。
// 仅处理: 单 '-' 开头、长度>2、不含 '='、非 '--' 的 token。其余原样保留。
//
//   - level≥1: 当 token 内每个字符都是 bool 短选项时，拆成多个 '-x'。eg: -aux => -a -u -x
//   - level≥2: 首字符是"取值型"已注册短选项且其后有残余时，拆成 '-O' + rest(作为值)。eg: -Ostdout => -O stdout
//
// 不满足条件的 token 原样保留(不做猜测，避免误伤混合写法如 -aO)。
//
//   - shorts: short=>fullName 映射，判断字符是否为已注册短选项
//   - isBool: 判断给定单字符短选项是否为 bool 类型
//   - level:  EnhanceShort 等级(0 直接返回原 args)
func expandShortArgs(args []string, shorts map[string]string, isBool func(short string) bool, level uint8) []string {
	if level == 0 || len(args) == 0 {
		return args
	}

	out := make([]string, 0, len(args))
	for _, arg := range args {
		// 仅处理单 '-' 开头、长度>2、不含 '='
		if len(arg) <= 2 || arg[0] != '-' || arg[1] == '-' || strings.IndexByte(arg, '=') >= 0 {
			out = append(out, arg)
			continue
		}

		body := arg[1:] // 去掉前导 '-'。eg: "aux"

		// level≥1: 全为 bool 短选项 => 逐个拆成 -x
		allBool := true
		for i := 0; i < len(body); i++ {
			if !isBool(string(body[i])) {
				allBool = false
				break
			}
		}
		if allBool {
			for i := 0; i < len(body); i++ {
				out = append(out, "-"+string(body[i]))
			}
			continue
		}

		// level≥2: 首字符是取值型已注册短选项 => -X + 余下作为值
		if level >= 2 {
			first := string(body[0])
			if _, ok := shorts[first]; ok && !isBool(first) {
				out = append(out, "-"+first, body[1:])
				continue
			}
		}

		// 不满足 => 原样保留(不猜测)
		out = append(out, arg)
	}
	return out
}
