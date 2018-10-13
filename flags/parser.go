package flags

import (
	"fmt"
	"strings"
)

// ArgsParser definition. a simple command line args parser
type ArgsParser struct {
	// BoolOpts bool option names. list all bool value options.
	// eg. "--debug -h" -> []string{"debug", "h"}
	BoolOpts []string
	// bool option name map. it's flip from BoolOpts.
	// eg. {"name": false}
	boolOpts map[string]bool
	// ArrayOpts array option names. list all array value options.
	// eg. "--name tom --name john" should add []string{"name"}
	ArrayOpts []string
	// array option name map. it's flip from ArrayOpts.
	// eg. {"name": false}
	arrayOpts map[string]bool
	// ValidOpts list all valid option names.
	// ValidOpts []string
	// current index for loop rawArgs
	index int
	// raw args
	rawArgs []string
	// raw args length
	length int
	// parsed longs options. value allow: bool, string, array
	longOpts map[string]interface{}
	// parsed shorts options. value allow: bool, string, array
	shortOpts map[string]interface{}
	// parsed arguments
	args []string
}

// ParseArgs parse os.Args to options.
func ParseArgs(args []string, boolOpts []string, arrayOpts []string) *ArgsParser {
	p := &ArgsParser{
		BoolOpts:  boolOpts,
		ArrayOpts: arrayOpts,
	}

	p.Parse(args)
	return p
}

// Opts get parsed opts
func (p *ArgsParser) Opts() map[string]interface{} {
	return map[string]interface{}{
		"longs":  p.longOpts,
		"shorts": p.shortOpts,
	}
}

// Args get parsed args
func (p *ArgsParser) Args() []string {
	return p.args
}

// OptsString convert all options to string
func (p *ArgsParser) OptsString() string {
	return fmt.Sprintf("long opts: %#v\nshort opts: %#v\n", p.longOpts, p.shortOpts)
}

func (p *ArgsParser) prepare() {
	if len(p.BoolOpts) > 0 {
		p.boolOpts = p.flipSlice(p.BoolOpts)
	}

	if len(p.ArrayOpts) > 0 {
		p.arrayOpts = p.flipSlice(p.ArrayOpts)
	}

	p.longOpts = make(map[string]interface{})
	p.shortOpts = make(map[string]interface{})
}

/*************************************************************
 * command options and arguments parse
 *************************************************************/

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
func (p *ArgsParser) Parse(args []string) {
	p.rawArgs = args
	p.prepare()
	p.length = len(args)

	for p.index < p.length {
		current := args[p.index]
		p.index++

		if current != "" && current[0] == '-' { // is option
			p.parseOne(current)
		} else { // is argument
			p.args = append(p.args, current)
		}
	}
}

func (p *ArgsParser) parseOne(current string) {
	var val string
	val = "true"
	noVal := true // mark current option is no value assigned
	isLong := false
	opt := current[1:]

	if opt[0] == '-' { // lang option: --opt
		opt = strings.TrimLeft(opt, "-=")
		if opt == "" { // invalid. eg "--="
			return
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
			return
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
			p.shortOpts[string(n)] = noVal
		}
		return
	}

	// collect option and value
	p.collectOption(opt, val, isLong)
}

func (p *ArgsParser) next() (val string, valid bool) {
	if p.index >= p.length { // end
		return
	}

	return p.rawArgs[p.index], true
}

func (p *ArgsParser) collectOption(opt, val string, isLong bool) {
	isArray := p.isArrayOpt(opt)
	if isLong {
		if isArray {
			vs, ok := p.longOpts[opt]
			if !ok {
				vs = []string{val}
			} else {
				vs = append(vs.([]string), val)
			}

			p.longOpts[opt] = vs
		} else {
			bl, err := parseBool(val)
			if err != nil {
				p.longOpts[opt] = val
			} else {
				p.longOpts[opt] = bl
			}
		}
		return
	}

	// short
	if isArray {
		vs, ok := p.shortOpts[opt]
		if !ok {
			vs = []string{val}
		} else {
			vs = append(vs.([]string), val)
		}

		p.shortOpts[opt] = vs
	} else {
		bl, err := parseBool(val)
		if err != nil {
			p.shortOpts[opt] = val
		} else {
			p.shortOpts[opt] = bl
		}
	}
}

func (p *ArgsParser) isValue(str string) bool {
	if str == "" {
		return true
	}

	return str[0] != '-'
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

// parseBool parse string to bool
func parseBool(str string) (bool, error) {
	lower := strings.ToLower(str)
	switch lower {
	case "1", "on", "yes", "true":
		return true, nil
	case "0", "off", "no", "false":
		return false, nil
	}

	return false, fmt.Errorf("'%s' cannot convert to bool", str)
}
