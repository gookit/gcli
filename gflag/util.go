package gflag

import (
	"strings"

	"github.com/gookit/gcli/v3/helper"
	"github.com/gookit/goutil/arrutil"
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
var flagTagKeys = arrutil.Strings{"name", "desc", "required", "default", "shorts"}

// parse tag named k-v value. item split by ';'
//
// eg: "name=int0;shorts=i;required=true;desc=int option message"
//
// supported field name:
//
//	name
//	desc
//	shorts
//	required
//	default
//
// TODO use structs.ParseTagValueNamed()
func parseNamedRule(name, rule string) (mp map[string]string) {
	ss := strutil.Split(rule, ";")
	if len(ss) == 0 {
		return
	}

	mp = make(map[string]string, len(flagTagKeys))
	for _, s := range ss {
		if strings.ContainsRune(s, '=') == false {
			helper.Panicf("parse tag error on field '%s': item must match `KEY=VAL`", name)
		}

		kvNodes := strings.SplitN(s, "=", 2)
		key, val := kvNodes[0], strings.TrimSpace(kvNodes[1])
		if !flagTagKeys.Has(key) {
			helper.Panicf("parse tag error on field '%s': invalid key name '%s'", name, key)
		}

		mp[key] = val
	}
	return
}

// ParseSimpleRule struct tag value use simple rule. each item split by ';'
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
// returns field name:
//
//	name
//	desc
//	shorts
//	required
//	default
//
// TODO use structs.ParseTagValueDefine() and support name.
func ParseSimpleRule(name, rule string) (mp map[string]string) {
	ss := strutil.SplitNTrimmed(rule, ";", 4)
	ln := len(ss)
	if ln == 0 {
		return
	}

	mp = make(map[string]string, ln)
	mp["desc"] = ss[0]
	if ln == 1 {
		return
	}

	required := ss[1]
	if required == "required" {
		required = "true"
	}

	mp["required"] = required

	// has shorts and default
	if ln > 3 {
		mp["default"], mp["shorts"] = ss[2], ss[3]
	} else if ln > 2 {
		mp["default"] = ss[2]
	}
	return
}
