package gflag

import (
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

func getRequiredMark(must bool) string {
	if must {
		return "<red>*</>"
	}
	return ""
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
var flagTagKeys = arrutil.Strings{"name", "desc", "required", "default", "shorts"}
var flagTagKeys1 = arrutil.Strings{"desc", "required", "default", "shorts"}
var flagArgKeys = arrutil.Strings{"desc", "required", "default"}

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
