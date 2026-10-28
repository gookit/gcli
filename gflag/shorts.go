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
// Tokens consumed as option values and tokens after -- or stopAt stay verbatim.
//
//   - fs: 注册的选项及 short=>fullName 映射
//   - level:  EnhanceShort 等级(0 直接返回原 args)
func expandShortArgs(args []string, fs *FlagSet, level uint8, stopAt func(string) bool) []string {
	if level == 0 || len(args) == 0 {
		return args
	}

	out := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			out = append(out, args[i:]...)
			break
		}

		// 仅处理单 '-' 开头、长度>2、不含 '='
		if len(arg) > 2 && arg[0] == '-' && arg[1] != '-' && strings.IndexByte(arg, '=') < 0 {
			body := arg[1:] // 去掉前导 '-'。eg: "aux"

			// level≥1: 全为 bool 短选项 => 逐个拆成 -x
			allBool := true
			for j := 0; j < len(body); j++ {
				if !fs.isBoolShort(string(body[j])) {
					allBool = false
					break
				}
			}
			if allBool {
				for j := 0; j < len(body); j++ {
					out = append(out, "-"+string(body[j]))
				}
				continue
			}

			// level≥2: 首字符是取值型已注册短选项 => -X + 余下作为值
			if level >= 2 {
				first := string(body[0])
				if _, ok := fs.shorts[first]; ok && !fs.isBoolShort(first) {
					out = append(out, "-"+first, body[1:])
					continue
				}
			}
		}

		// Recognize registered short groups before classifying numeric-looking
		// tokens as positional arguments (for example, the bool aliases -1 -2).
		if stopAt != nil && !looksLikeOption(arg) && stopAt(arg) {
			out = append(out, args[i:]...)
			break
		}

		// 不满足 => 原样保留(不猜测)
		out = append(out, arg)

		// A known value-taking option consumes its next token, even if that
		// token resembles a short group or the -- terminator.
		if len(arg) > 1 && arg[0] == '-' && strings.IndexByte(arg, '=') < 0 {
			name := strings.TrimPrefix(arg[1:], "-")
			if known, isBool := fs.optMeta(name); known && !isBool && i+1 < len(args) {
				i++
				out = append(out, args[i])
			}
		}
	}
	return out
}
