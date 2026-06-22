package gflag

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"

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
	// reorderStop optional predicate for args reorder. when it returns true for a
	// token, reordering stops at that token (used to not cross a sub-command).
	reorderStop func(name string) bool
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

// SetReorderStop set the stop predicate for args reorder.
//
// When it returns true for a token, the auto-reorder stops at that token and
// keeps the rest verbatim. gcli uses it to avoid reordering across a sub-command
// boundary, so only the final executed command's args are reordered.
func (p *Parser) SetReorderStop(fn func(name string) bool) { p.reorderStop = fn }

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
		// NOTE: 必须赋值给具名返回值 err，否则 panic 会被静默吞掉、对外仍返回 nil。
		if re := recover(); re != nil {
			err = fmt.Errorf("gflag: parse options panic: %v", re)
		}
	}()

	// prepare
	if err := p.prepare(); err != nil {
		return err
	}

	// POSIX 短选项规范化预处理(EnhanceShort>0 时生效，默认 0 不改变行为)
	if p.cfg.EnhanceShort > 0 {
		args = expandShortArgs(args, p.fSet.shorts, p.fSet.isBoolShort, p.cfg.EnhanceShort)
	}

	// 自动重排 args 为标准 "options... arguments" 形态(默认开启)。
	// 让写在 arguments 之后的 options 仍能被正确解析。
	if !p.cfg.DisableReorderArgs {
		args = rearrangeArgs(args, p.fSet, p.reorderStop)
	}

	// do parsing options
	if err = p.fSet.Parse(args); err != nil {
		return err
	}

	// after options parse hook
	if p.AfterParse != nil {
		if err := p.AfterParse(p); err != nil {
			return err
		}
	}

	// call options validations(shared with ParseOpts)
	// parsing named arguments. TODO err = p.ParseArgs(p.fSet.Args())
	return p.validateAll()
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

	if v.Type().Kind() != reflect.Struct {
		return errNotAnStruct
	}

	if len(ruleType) > 0 {
		p.SetRuleType(ruleType[0])
	}

	tagName := p.cfg.GetTagName()
	return p.fromStructValue(v, tagName)
}

// fromStructValue parse the struct value fields and bind them as flag options.
// it is split from FromStruct so that anonymous nested fields can recurse.
func (p *Parser) fromStructValue(v reflect.Value, tagName string) error {
	t := v.Type()

	var mp maputil.SMap
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		name := sf.Name
		// skip cannot export field
		if name[0] >= 'a' && name[0] <= 'z' {
			continue
		}

		// support anonymous struct field: recursively expand its inner fields
		if sf.Anonymous {
			aft := sf.Type
			afv := v.Field(i)
			if aft.Kind() == reflect.Ptr {
				if afv.IsNil() {
					continue // anonymous nil pointer, skip expand
				}
				aft = aft.Elem()
				afv = afv.Elem()
			}
			if aft.Kind() == reflect.Struct {
				if err := p.fromStructValue(afv, tagName); err != nil {
					return err
				}
			}
			continue // anonymous field itself is not a single option
		}

		// field rule: use field name as option name, read meta from independent tag keys.
		// only treat as an option field when one of flag/desc/default/required tag exists.
		var str string
		if p.cfg.TagRuleType == TagRuleField {
			flagTag, hasFlag := sf.Tag.Lookup(tagName)
			if hasFlag && flagTag == "-" {
				continue // explicit skip by `flag:"-"`
			}

			_, hasDesc := sf.Tag.Lookup("desc")
			_, hasDef := sf.Tag.Lookup("default")
			_, hasReq := sf.Tag.Lookup("required")
			if !hasFlag && !hasDesc && !hasDef && !hasReq {
				continue // not an option field
			}
		} else {
			// eg: "name=int0;shorts=i;required=true;desc=int option message"
			str = sf.Tag.Get(tagName)
			if str == "" {
				continue
			}
		}

		fv := v.Field(i)
		ft := sf.Type
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
			var err error
			mp, err = structs.ParseTagValueNamed(name, str, namedTagKeys...)
			if err != nil {
				return err
			}
		} else if p.cfg.TagRuleType == TagRuleSimple {
			mp = parseSimpleRule(str)
		} else if p.cfg.TagRuleType == TagRuleField {
			mp = maputil.SMap{
				"desc":     sf.Tag.Get("desc"),
				"default":  sf.Tag.Get("default"),
				"required": sf.Tag.Get("required"),
				"shorts":   sf.Tag.Get(tagName), // flag tag as shorts
			}
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

		// get the field address safely. struct is passed by pointer, so its
		// fields (and deref'd pointer fields) are always addressable.
		if !fv.CanAddr() {
			return fmt.Errorf("field: %s - is not addressable for binding flag", name)
		}
		addr := fv.Addr().Interface()
		switch ft.Kind() {
		case reflect.Bool:
			p.BoolVar(addr.(*bool), opt)
		case reflect.Int:
			p.IntVar(addr.(*int), opt)
		case reflect.Int64:
			p.Int64Var(addr.(*int64), opt)
		case reflect.Uint:
			p.UintVar(addr.(*uint), opt)
		case reflect.Uint64:
			p.Uint64Var(addr.(*uint64), opt)
		case reflect.Float64:
			p.Float64Var(addr.(*float64), opt)
		case reflect.String:
			p.StrVar(addr.(*string), opt)
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
