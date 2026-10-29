package gflag

import "strings"

// StrictFormatArgs normalizes the long option form of the args, for the strict mode
// of the cli app: "--a" becomes "-a" and "---name" becomes "--name".
//
// The tokens that are not options are kept verbatim:
//
//   - everything after the "--" terminator, it is positional data;
//   - the value of a known value-taking option, even if it looks like an option.
//     eg: "--name ---deep" keeps "---deep" as the value of "name". A bool option
//     does not take a value, so the next token is still an option.
//
// Short aliases are resolved, so "-n ---deep" works like "--name ---deep". An
// unknown option is left as-is, the flag parser reports it.
func StrictFormatArgs(args []string, fs *FlagSet) []string {
	if len(args) == 0 {
		return args
	}

	fmtArgs := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]

		// args after "--" are positional, keep them verbatim
		if arg == "--" {
			return append(fmtArgs, args[i:]...)
		}

		// if contains '=' append self, the form is already clear
		if strings.ContainsRune(arg, '=') {
			fmtArgs = append(fmtArgs, arg)
			continue
		}

		// the value token must be copied before the loop can move on. a "--" is left to
		// the terminator branch above, so the args after it stay positional, as before.
		consumesValue := takesValue(fs, arg) && i+1 < len(args) && args[i+1] != "--"

		// eg: --a ---name 长选项形态规范化
		if strings.HasPrefix(arg, "--") {
			farg := strings.TrimLeft(arg, "-")
			if rl := len(farg); rl == 1 { // fix: "--a" -> "-a"
				arg = "-" + farg
			} else if rl > 1 { // fix: "---name" -> "--name"
				arg = "--" + farg
			}
		}

		fmtArgs = append(fmtArgs, arg)

		// the value of a value-taking option is not an option, keep it verbatim
		if consumesValue {
			fmtArgs = append(fmtArgs, args[i+1])
			i++
		}
	}
	return fmtArgs
}

// takesValue reports whether arg is a known option that takes a value. The short
// aliases are resolved too, so "-n" is recognized like "--name".
func takesValue(fs *FlagSet, arg string) bool {
	if fs == nil || len(arg) < 2 || arg[0] != '-' {
		return false
	}

	name := strings.TrimLeft(arg, "-")
	if name == "" || strings.ContainsRune(name, '=') {
		return false
	}

	known, isBool := fs.optMeta(name)
	return known && !isBool
}
