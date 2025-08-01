package gflag

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"unsafe"

	"github.com/gookit/color"
	"github.com/gookit/color/colorp"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/structs"
	"github.com/gookit/goutil/strutil"
)

// Parser type, alias of Flags type
type Parser = Flags

// HandleFunc type
type HandleFunc func(fs *Flags) error

// Flags cli flag options and arguments management, parsing, and binding.
type Flags struct {
	// --- cli options ---
	CliOpts
	// --- cli arguments ---
	CliArgs

	name string
	// Desc message
	Desc string
	// AfterParse options hook
	AfterParse func(fs *Flags) error

	// cfg option for the flags parser
	cfg *Config
	// buf for build help message
	buf *bytes.Buffer
	// output for print help message
	out io.Writer
	// HandleFunc handle func, will call on Run()
	HandleFunc HandleFunc
}

func newDefaultFlagConfig() *Config {
	return &Config{
		Alignment: AlignLeft,
		TagName:   FlagTagName,
	}
}

// New create a new Parser and init it.
func New(nameWithDesc ...string) *Flags {
	p := &Flags{
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

	p.SetName(fName)
	p.InitFlagSet()
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
	p.InitFlagSet()
}

// SetName for parser
func (p *Parser) SetName(name string) {
	p.name = name
	p.CliOpts.SetName(name)
	p.CliArgs.SetName(name)
}

// ParserCfg for the parser.
func (p *Parser) ParserCfg() *Config {
	return p.cfg
}

// SetConfig for the object.
func (p *Parser) SetConfig(opt *Config) { p.cfg = opt }

// UseSimpleRule for the parse tag value rule string. see TagRuleSimple
func (p *Parser) UseSimpleRule() *Parser {
	p.cfg.TagRuleType = TagRuleSimple
	return p
}

// SetRuleType for the parse tag value rule string.
func (p *Parser) SetRuleType(rt uint8) *Parser {
	p.cfg.TagRuleType = rt
	return p
}

// WithConfigFn for the object.
func (p *Parser) WithConfigFn(fns ...ConfigFunc) *Parser {
	for _, fn := range fns {
		fn(p.cfg)
	}
	return p
}

// SetHandle func on flags parsed
func (p *Parser) SetHandle(fn HandleFunc) *Flags {
	p.HandleFunc = fn
	return p
}

/***********************************************************************
 * Flags:
 * - parse input flags
 ***********************************************************************/

// Run parse options and arguments, and HandleFunc help render
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
		colorp.Cyanln(" ", binFile, "[--Options ...] [Arguments ...]\n")
		p.PrintHelpPanel()
	})

	// do parsing options
	if err := p.Parse(waitArgs); err != nil {
		if err != flag.ErrHelp {
			color.Errorf("Parse options error: %s\n", err.Error())
		}
		return // ignore help error
	}

	// parsing named arguments.
	if err := p.ParseArgs(p.fSet.Args()); err != nil {
		color.Errorf("Parse arguments error: %s\n", err.Error())
		return
	}

	if p.HandleFunc != nil {
		if err := p.HandleFunc(p); err != nil {
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
			color.Errorln("gflag.Parse Error:", err)
		}
	}()

	// prepare
	if err := p.prepare(); err != nil {
		return err
	}

	// do parsing options. TODO p.ParseOpts(args)
	if err = p.fSet.Parse(args); err != nil {
		return err
	}

	// after options parse hook
	if p.AfterParse != nil {
		if err := p.AfterParse(p); err != nil {
			return err
		}
	}

	// call options validations
	for name, opt := range p.opts {
		fItem := p.fSet.Lookup(name)
		if err = opt.Validate(fItem.Value.String()); err != nil {
			return err
		}
	}

	// parsing named arguments. TODO
	// err = p.ParseArgs(p.fSet.Args())
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

// MustFromStruct from struct tag binding options, panic if error
//
// more see FromStruct()
func (p *Parser) MustFromStruct(ptr any, ruleType ...uint8) {
	goutil.MustOK(p.FromStruct(ptr, ruleType...))
}

// FromStruct from struct tag binding options
//
// ## Named rule(default)
//
//	// tag format: name=val0;shorts=i;required=true;desc=a message
//	type UserCmdOpts struct {
//		Name string `flag:"name=name;shorts=n;required=true;desc=input username"`
//		Age int `flag:"name=age;shorts=a;required=true;desc=input user age"`
//	}
//	opt := &UserCmdOpts{}
//	p.FromStruct(opt)
//
// ## Simple rule
//
//	// tag format1: name;desc;required;default;shorts
//	// tag format2: desc;required;default;shorts
//	type UserCmdOpts struct {
//		Name string `flag:"input username;true;;n"`
//		Age int `flag:"age;input user age;true;;o"`
//	}
//	opt := &UserCmdOpts{}
//	p.FromStruct(opt, gflag.TagRuleSimple)
func (p *Parser) FromStruct(ptr any, ruleType ...uint8) (err error) {
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

	tagName := p.cfg.GetTagName()
	if len(ruleType) > 0 {
		p.SetRuleType(ruleType[0])
	}

	var mp maputil.SMap
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		name := sf.Name
		// skip cannot export field
		if name[0] >= 'a' && name[0] <= 'z' {
			continue
		}

		// TODO support anonymous field by sf.Anonymous
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
			mp, err = structs.ParseTagValueNamed(name, str, namedTagKeys...)
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

		opt := newOpt(optName, mp["desc"], mp["default"], mp.StrOne("shorts", "short"))
		if must, has := mp["required"]; has {
			opt.Required = strutil.QuietBool(must)
		}

		// field is implements flag.Value
		if ft.Implements(flagValueType) {
			p.Var(fv.Interface().(flag.Value), opt)
			continue
		}

		// field is addressable and implements flag.Value
		if fv.CanAddr() {
			if addrV := fv.Addr(); addrV.Type().Implements(flagValueType) {
				p.Var(addrV.Interface().(flag.Value), opt)
				continue
			}
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
			return fmt.Errorf("field: %s - unsupport type(%s) for binding flag", name, ft.String())
		}
	}
	return nil
}

// Required flag option name(s)
func (p *Parser) Required(names ...string) {
	for _, name := range names {
		opt, ok := p.opts[name]
		if !ok {
			panicf("config undefined option '%s'", cflag.AddPrefix(name))
		}
		opt.Required = true
	}
}

/***********************************************************************
 * Flags:
 * - helper methods
 ***********************************************************************/

// Name of the Flags
func (p *Parser) Name() string { return p.name }

// Len of the Flags
func (p *Parser) Len() int { return len(p.names) }

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
