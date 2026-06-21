package gflag

import "strings"

// optMeta reports whether name (a short or full option name, without leading
// dashes) is a known option, and whether it is a bool option.
//
// A bool option does not consume a following value token; a value-taking option
// does. used by rearrangeArgs to keep an option grouped with its value.
func (f *FlagSet) optMeta(name string) (known, isBool bool) {
	if full, ok := f.shorts[name]; ok {
		name = full
	}

	flg, ok := f.formal[name]
	if !ok {
		return false, false
	}
	if fv, ok := flg.Value.(boolFlag); ok && fv.IsBoolFlag() {
		return true, true
	}
	return true, false
}

// looksLikeOption reports whether s is an option token (-x or --xxx).
//
// A lone "-" and negative-number-like tokens (-5, -1.5, -.5) are NOT options,
// they are treated as positional arguments.
func looksLikeOption(s string) bool {
	if len(s) < 2 || s[0] != '-' {
		return false
	}

	c := s[1]
	if c >= '0' && c <= '9' || c == '.' { // -5, -1.5, -.5 => negative numbers
		return false
	}
	return true
}

// rearrangeArgs reorders args into the canonical "options... arguments" form, so
// that options written after positional arguments are still parsed correctly.
//
// Rules:
//   - options (and the value of a value-taking option) move to the front,
//     keeping their original relative order;
//   - positional arguments keep their original relative order, after options;
//   - a known value-taking option groups the next token as its value (-name tom);
//   - bool options and the `--opt=val` form do not consume the next token;
//   - negative-number-like tokens and a lone "-" are arguments (see looksLikeOption);
//   - "--" terminates reordering: it and everything after are kept verbatim;
//   - when stopAt(token) is true (token is a known sub-command name), reordering
//     stops: that token and everything after are kept verbatim. This confines the
//     reorder to the final executed command in a multi-level app.
func rearrangeArgs(args []string, fs *FlagSet, stopAt func(string) bool) []string {
	if len(args) < 2 {
		return args
	}

	opts := make([]string, 0, len(args))
	vals := make([]string, 0, len(args))

	i, n := 0, len(args)
	for i < n {
		s := args[i]

		// "--" terminator: keep it and the rest verbatim.
		if s == "--" {
			vals = append(vals, args[i:]...)
			break
		}

		// not an option-looking token => positional argument
		if !looksLikeOption(s) {
			// stop reordering at a known sub-command name; keep the rest verbatim.
			if stopAt != nil && stopAt(s) {
				vals = append(vals, args[i:]...)
				break
			}
			vals = append(vals, s)
			i++
			continue
		}

		// it's an option token. resolve the name part (strip dashes, cut off =value).
		name := s[1:]
		if name[0] == '-' {
			name = name[1:]
		}

		hasEq := false
		if idx := strings.IndexByte(name, '='); idx >= 0 {
			name = name[:idx]
			hasEq = true
		}

		opts = append(opts, s)
		i++

		// a known value-taking option consumes the next token as its value.
		if !hasEq {
			if known, isBool := fs.optMeta(name); known && !isBool && i < n {
				opts = append(opts, args[i])
				i++
			}
		}
	}

	return append(opts, vals...)
}
