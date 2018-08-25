package utils

import (
	"fmt"
	"strings"
)

// ParseBool parse string to bool
func ParseBool(str string) (bool, error) {
	lower := strings.ToLower(str)
	switch lower {
	case "1", "on", "yes", "true":
		return true, nil
	case "0", "off", "no", "false":
		return false, nil
	}

	return false, fmt.Errorf("'%s' cannot convert to bool", str)
}

// ResultSet the values for parsed arguments and options
type ResultSet struct {
	// value allow: bool, string, array
	longOpts  map[string]interface{}
	shortOpts map[string]interface{}
	// args list
	args []string
}

func (r *ResultSet) OptString(name string) string {
	return ""
}

// ArgsParser definition. a simple command line args parser
type ArgsParser struct {
	// BoolOpts define. list all bool value options.
	// eg. "--debug -h" -> []string{"debug", "h"}
	BoolOpts []string
	// ArrayOpts list all array value options.
	// eg. "--name tom --name john" should add []string{"name"}
	ArrayOpts []string
	// ValidOpts list all valid option names.
	ValidOpts []string
	//
	args []string
	// current index for loop
	index int
	// args length
	length int
	// boolOpts {"name": false}
	boolOpts map[string]bool
	// arrayOpts {"name": false}
	arrayOpts map[string]bool
	// result
	result *ResultSet
}

// Parse args list to options
//
// Supports options format:
// 	-e  // bool, short option
// 	-e <value> // short option
// 	-e=<value>
// 	-aux // multi short bool options
// 	--bool-opt // bool, lang option
// 	--long-opt <value> // lang option
// 	--long-opt=<value>
func (p *ArgsParser) Parse(args []string) *ResultSet {
	var val string
	p.args = args
	p.prepare()
	p.length = len(args)

	for p.index < p.length {
		cur := args[p.index]
		p.index++

		if cur[0] == '-' { // option
			val = "true"
			noVal := true // mark current option is no value assigned
			opt := cur[1:]
			isLong := false

			if opt[0] == '-' { // lang option: --opt
				opt = strings.TrimLeft(opt, "-=")
				if opt == "" { // invalid. eg "--="
					continue
				}

				isLong = true
				if strings.IndexByte(opt, '=') > -1 { // has val: --opt=VAL
					ss := strings.SplitN(opt, "=", 2)
					opt, val = ss[0], ss[1]
					noVal = false
				}
			} else { // short option
				opt = strings.TrimLeft(opt, "-=")
				if opt == "" { // invalid. eg "-="
					continue
				}

				// has val: -s=VAL
				if len(opt) > 1 && strings.IndexByte(opt, '=') > -1 {
					ss := strings.SplitN(opt, "=", 2)
					opt, val = ss[0], ss[1]
					noVal = false
				}
			}

			// get next elem value
			nxt, valid := p.next()

			// current opt no value and next is value.
			if valid && noVal && !p.isBoolOpt(opt) && p.isValue(nxt) {
				val = nxt
				p.index++
			} else if !isLong && noVal { // short bool opts. like -e -abc
				for _, n := range []rune(opt) {
					p.result.shortOpts[string(n)] = val
				}
				continue
			}

			// collect option and value
			p.collectOption(opt, val, isLong)
		} else { // args
			p.result.args = append(p.result.args, cur)
		}
	}

	return p.result
}

func (p *ArgsParser) next() (val string, valid bool) {
	if p.index >= p.length { // end
		return
	}

	return p.args[p.index], true
}

func (p *ArgsParser) collectOption(opt, val string, isLong bool) {
	isArray := p.isArrayOpt(opt)
	if isLong {
		if isArray {
			vs, ok := p.result.longOpts[opt]
			if !ok {
				vs = []string{val}
			} else {
				vs = append(vs.([]string), val)
			}

			p.result.longOpts[opt] = vs
		} else {
			bl, err := ParseBool(val)
			if err != nil {
				p.result.longOpts[opt] = val
			} else {
				p.result.longOpts[opt] = bl
			}
		}

		return
	}

	// short
	if isArray {
		vs, ok := p.result.shortOpts[opt]
		if !ok {
			vs = []string{val}
		} else {
			vs = append(vs.([]string), val)
		}

		p.result.shortOpts[opt] = vs
	} else {
		bl, err := ParseBool(val)
		if err != nil {
			p.result.shortOpts[opt] = val
		} else {
			p.result.shortOpts[opt] = bl
		}
	}
}

func (p *ArgsParser) isValue(str string) bool {
	if str == "" {
		return true
	}

	return str[0] != '-'
}

func (p *ArgsParser) prepare() {
	if len(p.BoolOpts) > 0 {
		p.boolOpts = p.flipSlice(p.BoolOpts)
	}

	if len(p.ArrayOpts) > 0 {
		p.arrayOpts = p.flipSlice(p.ArrayOpts)
	}

	p.result = &ResultSet{
		longOpts:  make(map[string]interface{}),
		shortOpts: make(map[string]interface{}),
	}

}

func (p *ArgsParser) flipSlice(ss []string) map[string]bool {
	m := make(map[string]bool)
	for _, v := range ss {
		m[v] = false
	}

	return m
}

func (p *ArgsParser) isBoolOpt(n string) bool {
	if p.boolOpts != nil {
		_, ok := p.boolOpts[n]
		return ok
	}

	return false
}

func (p *ArgsParser) isArrayOpt(n string) bool {
	if p.arrayOpts != nil {
		_, ok := p.arrayOpts[n]
		return ok
	}

	return false
}
