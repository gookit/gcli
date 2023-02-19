package gflag

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"unsafe"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/helper"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/structs"
	"github.com/gookit/goutil/strutil"
)

// Flags type
type Flags = Parser

// HandleFunc type
type HandleFunc func(p *Parser) error

// Parser cli flag options and arguments binding management and parsing.
type Parser struct {
	// --- cli options ---
	CliOpts

	// --- cli arguments ---
	CliArgs

	name string
	// Desc message
	Desc string
	// AfterParse options hook
	AfterParse func(fs *Parser) error

	// cfg option for the flags parser
	cfg *Config
	// buf for build help message
	buf *bytes.Buffer
	// output for print help message
	out io.Writer
	// handle func
	handle HandleFunc
}

func newDefaultFlagConfig() *Config {
	return &Config{
		Alignment: AlignLeft,
		TagName:   FlagTagName,
	}
}

// New create a new Parser
func New(nameWithDesc ...string) *Parser {
	p := &Parser{
		out: os.Stdout,
		cfg: newDefaultFlagConfig(),
	}
	// fs.ExitFunc = os.Exit

	fName := "gflag"
	if num := len(nameWithDesc); num > 0 {
		fName = nameWithDesc[0]
		if num > 1 {
			p.Desc = nameWithDesc[1]
		}
	}

	p.InitFlagSet(fName)
	return p
}

// Init for parser
func (p *Parser) Init(name string) {
	if p.out != nil {
		return
	}

	p.out = os.Stdout
	if p.cfg == nil {
		p.cfg = newDefaultFlagConfig()
	}

	p.SetName(name)
	p.InitFlagSet(name)
}

// SetName for parser
func (p *Parser) SetName(name string) {
	p.name = name
	p.CliOpts.SetName(name)
	p.CliArgs.SetName(name)
}

// SetConfig for the object.
func (p *Parser) SetConfig(opt *Config) { p.cfg = opt }

// UseSimpleRule for the parse tag value rule string. see TagRuleSimple
func (p *Parser) UseSimpleRule() *Parser {
	p.cfg.TagRuleType = TagRuleSimple
	return p
}

// WithConfigFn for the object.
func (p *Parser) WithConfigFn(fns ...func(cfg *Config)) *Parser {
	for _, fn := range fns {
		fn(p.cfg)
	}
	return p
}

// SetHandle func
func (p *Parser) SetHandle(fn HandleFunc) *Parser {
	p.handle = fn
	return p
}

/***********************************************************************
 * Flags:
 * - parse input flags
 ***********************************************************************/

// Run flags parse and handle help render
//
// Usage:
//
//		gf := gflag.New()
//	 ...
//		// OR: gf.Run(nil)
//		gf.Run(os.Args)
func (p *Parser) Run(args []string) {
	if args == nil {
		args = os.Args
	}

	// split binFile and args
	binFile, waitArgs := args[0], args[1:]

	// register help render
	p.SetHelpRender(func() {
		if p.Desc != "" {
			color.Infoln(p.Desc)
		}

		color.Comment.Println("Usage:")
		color.Cyan.Println(" ", binFile, "[--Options...] [CliArgs...]\n")

		p.PrintHelpPanel()
	})

	// do parsing
	if err := p.Parse(waitArgs); err != nil {
		if err == flag.ErrHelp {
			return // ignore help error
		}

		color.Errorf("Parse error: %s\n", err.Error())
	}

	if p.handle != nil {
		if err := p.handle(p); err != nil {
			color.Errorln(err)
		}
	}
}

// Parse given arguments
//
// Usage:
//
//	gf := gflag.New()
//	gf.BoolOpt(&debug, "debug", "", defDebug, "open debug mode")
//	gf.UintOpt(&port, "port", "p", 18081, "the http server port")
//
//	err := gf.Parse(os.Args[1:])
func (p *Parser) Parse(args []string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			color.Errorln("Flags.Parse Error:", err)
		}
	}()

	// prepare
	if err := p.prepare(); err != nil {
		return err
	}

	if len(p.shorts) > 0 && len(args) > 0 {
		args = cflag.ReplaceShorts(args, p.shorts)
		// TODO gcli.Debugf("replace shortcuts. now, args: %v", args)
	}

	// do parsing
	if err = p.fSet.Parse(args); err != nil {
		return err
	}

	// after hook
	if p.AfterParse != nil {
		if err := p.AfterParse(p); err != nil {
			return err
		}
	}

	// call flags validate
	for name, opt := range p.opts {
		fItem := p.fSet.Lookup(name)
		err = opt.Validate(fItem.Value.String())
		if err != nil {
			return err
		}
	}
	return
}

func (p *Parser) prepare() error {
	return nil
}

/***********************************************************************
 * Flags:
 * - binding option from struct
 ***********************************************************************/

var (
	flagValueType  = reflect.TypeOf(new(flag.Value)).Elem()
	errNotPtrValue = errors.New("must provide an ptr value")
	errNotAnStruct = errors.New("must provide an struct ptr")
	errTagRuleType = errors.New("invalid tag rule type on struct")
)

// FromStruct from struct tag binding options
func (p *Parser) FromStruct(ptr any) (err error) {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		return errNotPtrValue
	}

	if !v.IsNil() {
		v = v.Elem()
	}

	t := v.Type()
	if t.Kind() != reflect.Struct {
		return errNotAnStruct
	}

	tagName := p.cfg.TagName
	if tagName == "" {
		tagName = FlagTagName
	}

	var mp map[string]string
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		name := sf.Name

		// skip cannot export field
		if name[0] >= 'a' && name[0] <= 'z' {
			continue
		}

		// eg: "name=int0;shorts=i;required=true;desc=int option message"
		str := sf.Tag.Get(tagName)
		if str == "" {
			continue
		}

		fv := v.Field(i)
		ft := t.Field(i).Type
		if !fv.CanInterface() {
			continue
		}

		// is pointer
		// var isPtr bool
		// var isNilPtr bool
		if ft.Kind() == reflect.Ptr {
			// isPtr = true
			if fv.IsNil() {
				return fmt.Errorf("field: %s - nil pointer dereference", name)
			}

			ft = ft.Elem()
			fv = fv.Elem()
		}

		// eg: "name=int0;shorts=i;required=true;desc=int option message"
		if p.cfg.TagRuleType == TagRuleNamed {
			// mp = parseNamedRule(name, str)
			mp, err = structs.ParseTagValueNamed(name, str, flagTagKeys...)
			if err != nil {
				return err
			}
		} else if p.cfg.TagRuleType == TagRuleSimple {
			mp = parseSimpleRule(str)
		} else {
			return errTagRuleType
		}

		// for create flag opt
		optName, has := mp["name"]
		if !has { // use field as option name.
			optName = strutil.SnakeCase(name, "-")
		}

		opt := newOpt(optName, mp["desc"], mp["default"], mp["shorts"])
		if must, has := mp["required"]; has {
			opt.Required = strutil.QuietBool(must)
		}

		// field is implements flag.Value
		if ft.Implements(flagValueType) {
			p.Var(fv.Interface().(flag.Value), opt)
			continue
		}

		// get field ptr addr
		ptr := unsafe.Pointer(fv.UnsafeAddr())
		switch ft.Kind() {
		case reflect.Bool:
			p.BoolVar((*bool)(ptr), opt)
		case reflect.Int:
			p.IntVar((*int)(ptr), opt)
			// if isNilPtr {
			// 	fv.SetInt(0)
			// 	newPtr := unsafe.Pointer(fv.UnsafeAddr())
			// 	p.IntVar((*int)(newPtr), opt)
			// } else {
			// 	p.IntVar((*int)(ptr), opt)
			// }
		case reflect.Int64:
			p.Int64Var((*int64)(ptr), opt)
		case reflect.Uint:
			p.UintVar((*uint)(ptr), opt)
		case reflect.Uint64:
			p.Uint64Var((*uint64)(ptr), opt)
		case reflect.Float64:
			p.Float64Var((*float64)(ptr), opt)
		case reflect.String:
			p.StrVar((*string)(ptr), opt)
		default:
			return fmt.Errorf("field: %s - invalid type for binding flag", name)
		}
	}
	return nil
}

// Required flag option name(s)
func (p *Parser) Required(names ...string) {
	for _, name := range names {
		opt, ok := p.opts[name]
		if !ok {
			helper.Panicf("config undefined option '%s'", cflag.AddPrefix(name))
		}
		opt.Required = true
	}
}

/***********************************************************************
 * Flags:
 * - render help message
 ***********************************************************************/

// SetHelpRender set the raw *flag.FlagSet.Usage
func (p *Parser) SetHelpRender(fn func()) {
	p.fSet.Usage = fn
}

// PrintHelpPanel for all options to the gf.out
func (p *Parser) PrintHelpPanel() {
	color.Fprint(p.out, p.String())
}

// String for all flag options
func (p *Parser) String() string {
	return p.BuildHelp()
}

// BuildHelp string for all flag options
func (p *Parser) BuildHelp() string {
	if p.buf == nil {
		p.buf = new(bytes.Buffer)
	}

	// repeat call the method
	if p.buf.Len() < 1 {
		p.buf.WriteString("Options:\n")
		p.buf.WriteString(p.BuildOptsHelp())
		p.buf.WriteByte('\n')

		if p.HasArgs() {
			p.buf.WriteString("Arguments:\n")
			p.buf.WriteString(p.BuildArgsHelp())
			p.buf.WriteByte('\n')
		}
	}

	return p.buf.String()
}

// BuildOptsHelp string.
func (p *Parser) BuildOptsHelp() string {
	var sb strings.Builder

	p.fSet.VisitAll(func(f *flag.Flag) {
		line := p.formatOneFlag(f)
		if line != "" {
			sb.WriteString(line)
			sb.WriteByte('\n')
		}
	})

	return sb.String()
}

func (p *Parser) formatOneFlag(f *flag.Flag) (s string) {
	// Skip render:
	// - opt is not exists(Has ensured that it is not a short name)
	// - it is hidden flag option
	// - flag desc is empty
	opt, has := p.opts[f.Name]
	if !has || opt.Hidden {
		return
	}

	var fullName string
	name := f.Name
	// eg: "-V, --version" length is: 13
	nameLen := p.names[name]
	// display description on new line
	descNl := p.cfg.DescNewline

	var nlIndent string
	if descNl {
		nlIndent = "\n        "
	} else {
		nlIndent = "\n      " + strings.Repeat(" ", p.optMaxLen)
	}

	// add prefix '-' to option
	fullName = cflag.AddPrefixes2(name, opt.Shorts, true)
	s = fmt.Sprintf("  <info>%s</>", fullName)

	// - build flag type info
	typeName, desc := flag.UnquoteUsage(f)
	// typeName: option value data type: int, string, ..., bool value will return ""
	if !p.cfg.WithoutType && len(typeName) > 0 {
		typeLen := len(typeName) + 1
		if !descNl && nameLen+typeLen > p.optMaxLen {
			descNl = true
		} else {
			nameLen += typeLen
		}

		s += fmt.Sprintf(" <magenta>%s</>", typeName)
	}

	if descNl {
		s += nlIndent
	} else {
		// padding space to optMaxLen width.
		if padLen := p.optMaxLen - nameLen; padLen > 0 {
			s += strings.Repeat(" ", padLen)
		}
		s += "    "
	}

	// --- build description
	if desc == "" {
		desc = defaultDesc
	} else {
		desc = strings.Replace(strutil.UpperFirst(desc), "\n", nlIndent, -1)
	}

	s += getRequiredMark(opt.Required) + desc

	// ---- append default value
	if isZero, isStr := cflag.IsZeroValue(f, f.DefValue); !isZero {
		if isStr {
			s += fmt.Sprintf(" (default <magentaB>%q</>)", f.DefValue)
		} else {
			s += fmt.Sprintf(" (default <magentaB>%v</>)", f.DefValue)
		}
	}

	return s
}

/***********************************************************************
 * Flags:
 * - helper methods
 ***********************************************************************/

// Name of the Flags
func (p *Parser) Name() string { return p.fSet.Name() }

// Len of the Flags
func (p *Parser) Len() int { return len(p.names) }

// FSet get the raw *flag.FlagSet
func (p *Parser) FSet() *flag.FlagSet { return p.fSet }

// SetFlagSet set the raw *flag.FlagSet
func (p *Parser) SetFlagSet(fSet *flag.FlagSet) { p.fSet = fSet }

// SetOutput for the Flags
func (p *Parser) SetOutput(out io.Writer) { p.out = out }

// FlagNames return all option names
func (p *Parser) FlagNames() map[string]int { return p.names }

// RawArg get an argument value by index
func (p *Parser) RawArg(i int) string { return p.fSet.Arg(i) }

// RawArgs get all raw arguments.
// if have been called parse, the return is remaining args.
func (p *Parser) RawArgs() []string { return p.fSet.Args() }

// FSetArgs get all raw arguments. alias of the RawArgs()
// if have been called parse, the return is remaining args.
func (p *Parser) FSetArgs() []string { return p.fSet.Args() }
