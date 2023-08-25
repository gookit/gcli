package gflag

import (
	"fmt"

	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/comdef"
	"github.com/gookit/goutil/strutil"
)

func sepStr(seps []string) string {
	if len(seps) > 0 {
		return seps[0]
	}
	return comdef.DefaultSep
}

func requiredMark(must bool) string {
	if must {
		return "<red>*</>"
	}
	return ""
}

// panicf message
func panicf(format string, v ...any) {
	panic(fmt.Sprintf("gflag: "+format, v...))
}

// allowed keys on struct tag.
//
// Parse named rule: parse tag named k-v value. item split by ';'
//
// eg: "name=int0;shorts=i;required=true;desc=int option message"
//
// Supported field name:
//
//	name
//	desc
//	shorts
//	required
//	default
var (
	flagTagKeys  = arrutil.Strings{"name", "desc", "required", "default", "shorts"}
	flagTagKeys1 = arrutil.Strings{"desc", "required", "default", "shorts"}
	flagArgKeys  = arrutil.Strings{"desc", "required", "default"}
)

// struct tag value use simple rule. each item split by ';'
//
//   - format: "name;desc;required;default;shorts"
//   - format: "desc;required;default;shorts"
//
// eg:
//
//	"int option message;required;i"
//	"opt-name;int option message;;a,b"
//	"int option message;;a,b;23"
//
// Returns field name:
//
//	name
//	desc
//	shorts
//	required
//	default
func parseSimpleRule(rule string) (mp map[string]string) {
	ss := strutil.SplitNTrimmed(rule, ";", 5)
	ln := len(ss)
	if ln == 0 {
		return
	}

	mp = make(map[string]string, ln)
	if ln == 1 {
		mp["desc"] = ss[0]
		return
	}

	// first is name
	if cflag.IsGoodName(ss[0]) {
		return arrutil.CombineToSMap(flagTagKeys, ss)
	}
	return arrutil.CombineToSMap(flagTagKeys1, ss)
}

// UnquoteUsage extracts a back-quoted name from the usage
// string for a flag and returns it and the un-quoted usage.
// Given "a `name` to show" it returns ("name", "a name to show").
// If there are no back quotes, the name is an educated guess of the
// type of the flag's value, or the empty string if the flag is boolean.
//
// Note: from go flag.UnquoteUsage()
func UnquoteUsage(flag *Flag) (name string, usage string) {
	// Look for a back-quoted name, but avoid the strings package.
	usage = flag.Usage
	for i := 0; i < len(usage); i++ {
		if usage[i] == '`' {
			for j := i + 1; j < len(usage); j++ {
				if usage[j] == '`' {
					name = usage[i+1 : j]
					usage = usage[:i] + name + usage[j+1:]
					return name, usage
				}
			}
			break // Only one back quote; use type name.
		}
	}

	// No explicit name, so use type if we can find one.
	name = "value"
	switch fv := flag.Value.(type) {
	case boolFlag:
		if fv.IsBoolFlag() {
			name = ""
		}
	case *durationValue:
		name = "duration"
	case *float64Value:
		name = "float"
	case *intValue, *int64Value:
		name = "int"
	case *stringValue:
		name = "string"
	case *uintValue, *uint64Value:
		name = "uint"
	}
	return
}
