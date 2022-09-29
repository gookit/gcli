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
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/strutil"
)

// Flags struct definition
type Flags struct {
	// Desc message
	Desc string
	// AfterParse options hook
	AfterParse func(fs *Flags) error

	// cfg option for the flags
	cfg *Config
	// the options flag set
	fSet *flag.FlagSet
	// buf for build help message
	buf *bytes.Buffer
	// output for print help message
	out io.Writer

	// all option names of the command. {name: length} // TODO delete, move len to meta.
	names map[string]int
	// metadata for all options
	metas map[string]*Option // TODO support option category
	// short names map for options. format: {short: name}
	// eg. {"n": "name", "o": "opt"}
	shorts map[string]string
	// support option category
	categories []OptCategory
	// flag name max length. useful for render help
	// eg: "-V, --version" length is 13
	flagMaxLen int
	// exist short names. useful for render help
	existShort bool

	// --- arguments
	Arguments
}

func newDefaultFlagConfig() *Config {
	return &Config{
		Alignment: AlignLeft,
		TagName:   FlagTagName,
	}
}

// New create a new Flags
func New(nameWithDesc ...string) *Flags {
	fs := &Flags{
		out: os.Stdout,
		cfg: newDefaultFlagConfig(),
	}
	// fs.ExitFunc = os.Exit

	fName := "gflag"
	if num := len(nameWithDesc); num > 0 {
		fName = nameWithDesc[0]
		if num > 1 {
			fs.Desc = nameWithDesc[1]
		}
	}

	fs.InitFlagSet(fName)
	return fs
}

// InitFlagSet create and init flag.FlagSet
func (fs *Flags) InitFlagSet(name string) {
	if fs.fSet != nil {
		return
	}

	if fs.cfg == nil {
		fs.cfg = newDefaultFlagConfig()
	}

	fs.fSet = flag.NewFlagSet(name, flag.ContinueOnError)
	// disable output internal error message on parse flags
	fs.fSet.SetOutput(io.Discard)
	// nothing to do ... render usage on after parsed
	fs.fSet.Usage = func() {}
}

// SetConfig for the object.
func (fs *Flags) SetConfig(opt *Config) { fs.cfg = opt }

// UseSimpleRule for the parse tag value rule string. see TagRuleSimple
func (fs *Flags) UseSimpleRule() *Flags {
	fs.cfg.TagRuleType = TagRuleSimple
	return fs
}

// WithConfigFn for the object.
func (fs *Flags) WithConfigFn(fns ...func(cfg *Config)) *Flags {
	for _, fn := range fns {
		fn(fs.cfg)
	}
	return fs
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
//		gf.Run(os.Args)
func (fs *Flags) Run(args []string) {
	if args == nil {
		args = os.Args
	}

	// split binFile and args
	binFile, waitArgs := args[0], args[1:]

	// register help render
	fs.SetHelpRender(func() {
		if fs.Desc != "" {
			color.Infoln(fs.Desc)
		}

		color.Comment.Println("Usage:")
		color.Cyan.Println(" ", binFile, "[--OPTIONS...]\n")
		color.Comment.Println("Options:")

		fs.PrintHelpPanel()
	})

	// do parsing
	if err := fs.Parse(waitArgs); err != nil {
		if err == flag.ErrHelp {
			return // ignore help error
		}

		color.Errorf("Parse error: %s\n", err.Error())
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
func (fs *Flags) Parse(args []string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			color.Errorln("Flags.Parse Error:", err)
		}
	}()

	// prepare
	if err := fs.prepare(); err != nil {
		return err
	}

	if len(fs.shorts) > 0 && len(args) > 0 {
		args = cflag.ReplaceShorts(args, fs.shorts)
		// TODO gcli.Debugf("replace shortcuts. now, args: %v", args)
	}

	// do parsing
	if err = fs.fSet.Parse(args); err != nil {
		return err
	}

	// after hook
	if fs.AfterParse != nil {
		if err := fs.AfterParse(fs); err != nil {
			return err
		}
	}

	// call flags validate
	for name, meta := range fs.metas {
		fItem := fs.fSet.Lookup(name)
		err = meta.Validate(fItem.Value.String())
		if err != nil {
			return err
		}
	}
	return
}

func (fs *Flags) prepare() error {
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
func (fs *Flags) FromStruct(ptr any) error {
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

	tagName := fs.cfg.TagName
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

		if fs.cfg.TagRuleType == TagRuleNamed {
			mp = parseNamedRule(name, str)
		} else if fs.cfg.TagRuleType == TagRuleSimple {
			mp = ParseSimpleRule(name, str)
		} else {
			return errTagRuleType
		}

		// for create flag meta
		optName, has := mp["name"]
		if !has { // use field as option name.
			optName = strutil.SnakeCase(name, "-")
		}

		meta := newFlagOpt(optName, mp["desc"], mp["default"], mp["shorts"])
		if must, has := mp["required"]; has {
			meta.Required = strutil.MustBool(must)
		}

		// field is implements flag.Value
		if ft.Implements(flagValueType) {
			fs.Var(fv.Interface().(flag.Value), meta)
			continue
		}

		// get field ptr addr
		ptr := unsafe.Pointer(fv.UnsafeAddr())
		switch ft.Kind() {
		case reflect.Bool:
			fs.BoolVar((*bool)(ptr), meta)
		case reflect.Int:
			fs.IntVar((*int)(ptr), meta)
			// if isNilPtr {
			// 	fv.SetInt(0)
			// 	newPtr := unsafe.Pointer(fv.UnsafeAddr())
			// 	fs.IntVar((*int)(newPtr), meta)
			// } else {
			// 	fs.IntVar((*int)(ptr), meta)
			// }
		case reflect.Int64:
			fs.Int64Var((*int64)(ptr), meta)
		case reflect.Uint:
			fs.UintVar((*uint)(ptr), meta)
		case reflect.Uint64:
			fs.Uint64Var((*uint64)(ptr), meta)
		case reflect.Float64:
			fs.Float64Var((*float64)(ptr), meta)
		case reflect.String:
			fs.StrVar((*string)(ptr), meta)
		default:
			return fmt.Errorf("field: %s - invalid type for binding flag", name)
		}
	}
	return nil
}

/***********************************************************************
 * Flags:
 * - binding option var
 ***********************************************************************/

// --- bool option

// Bool binding a bool option flag, return pointer
func (fs *Flags) Bool(name, shorts string, defVal bool, desc string) *bool {
	meta := newFlagOpt(name, desc, defVal, shorts)
	name = fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	p := fs.fSet.Bool(name, defVal, meta.Desc)
	meta.flag = fs.fSet.Lookup(name)

	return p
}

// BoolVar binding a bool option flag
func (fs *Flags) BoolVar(p *bool, meta *FlagMeta) { fs.boolOpt(p, meta) }

// BoolOpt binding a bool option
func (fs *Flags) BoolOpt(p *bool, name, shorts string, defVal bool, desc string) {
	fs.boolOpt(p, newFlagOpt(name, desc, defVal, shorts))
}

// binding option and shorts
func (fs *Flags) boolOpt(p *bool, meta *FlagMeta) {
	defVal := meta.DValue().Bool()
	name := fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	fs.fSet.BoolVar(p, name, defVal, meta.Desc)
	meta.flag = fs.fSet.Lookup(name)
}

// --- float option

// Float64Var binding an float64 option flag
func (fs *Flags) Float64Var(p *float64, meta *FlagMeta) { fs.float64Opt(p, meta) }

// Float64Opt binding a float64 option
func (fs *Flags) Float64Opt(p *float64, name, shorts string, defVal float64, desc string) {
	fs.float64Opt(p, newFlagOpt(name, desc, defVal, shorts))
}

func (fs *Flags) float64Opt(p *float64, meta *FlagMeta) {
	defVal := meta.DValue().Float64()
	name := fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	fs.fSet.Float64Var(p, name, defVal, meta.Desc)
	meta.flag = fs.fSet.Lookup(name)
}

// --- string option

// Str binding an string option flag, return pointer
func (fs *Flags) Str(name, shorts string, defValue, desc string) *string {
	meta := newFlagOpt(name, desc, defValue, shorts)
	name = fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	p := fs.fSet.String(name, defValue, meta.Desc)
	meta.flag = fs.fSet.Lookup(name)

	return p
}

// StrVar binding an string option flag
func (fs *Flags) StrVar(p *string, meta *FlagMeta) { fs.strOpt(p, meta) }

// StrOpt binding an string option
func (fs *Flags) StrOpt(p *string, name, shorts, defValue, desc string) {
	fs.strOpt(p, newFlagOpt(name, desc, defValue, shorts))
}

// binding option and shorts
func (fs *Flags) strOpt(p *string, meta *FlagMeta) {
	defVal := meta.DValue().String()
	name := fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	fs.fSet.StringVar(p, meta.Name, defVal, meta.Desc)
	meta.flag = fs.fSet.Lookup(name)
}

// --- intX option

// Int binding an int option flag, return pointer
func (fs *Flags) Int(name, shorts string, defValue int, desc string) *int {
	meta := newFlagOpt(name, desc, defValue, shorts)
	name = fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	p := fs.fSet.Int(name, defValue, meta.Desc)
	meta.flag = fs.fSet.Lookup(name)

	return p
}

// IntVar binding an int option flag
func (fs *Flags) IntVar(p *int, meta *FlagMeta) { fs.intOpt(p, meta) }

// IntOpt binding an int option
func (fs *Flags) IntOpt(p *int, name, shorts string, defValue int, desc string) {
	fs.intOpt(p, newFlagOpt(name, desc, defValue, shorts))
}

func (fs *Flags) intOpt(p *int, meta *FlagMeta) {
	defValue := meta.DValue().Int()
	name := fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	fs.fSet.IntVar(p, name, defValue, meta.Desc)
	meta.flag = fs.fSet.Lookup(name)
}

// Int64 binding an int64 option flag, return pointer
func (fs *Flags) Int64(name, shorts string, defValue int64, desc string) *int64 {
	meta := newFlagOpt(name, desc, defValue, shorts)
	name = fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	p := fs.fSet.Int64(name, defValue, meta.Desc)
	meta.flag = fs.fSet.Lookup(name)

	return p
}

// Int64Var binding an int64 option flag
func (fs *Flags) Int64Var(p *int64, meta *FlagMeta) { fs.int64Opt(p, meta) }

// Int64Opt binding an int64 option
func (fs *Flags) Int64Opt(p *int64, name, shorts string, defValue int64, desc string) {
	fs.int64Opt(p, newFlagOpt(name, desc, defValue, shorts))
}

func (fs *Flags) int64Opt(p *int64, meta *FlagMeta) {
	defVal := meta.DValue().Int64()
	name := fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	fs.fSet.Int64Var(p, name, defVal, meta.Desc)
	meta.flag = fs.fSet.Lookup(name)
}

// --- uintX option

// Uint binding an int option flag, return pointer
func (fs *Flags) Uint(name, shorts string, defVal uint, desc string) *uint {
	meta := newFlagOpt(name, desc, defVal, shorts)
	name = fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	p := fs.fSet.Uint(name, defVal, meta.Desc)
	meta.flag = fs.fSet.Lookup(name)

	return p
}

// UintVar binding an uint option flag
func (fs *Flags) UintVar(p *uint, meta *FlagMeta) { fs.uintOpt(p, meta) }

// UintOpt binding an uint option
func (fs *Flags) UintOpt(p *uint, name, shorts string, defValue uint, desc string) {
	fs.uintOpt(p, newFlagOpt(name, desc, defValue, shorts))
}

func (fs *Flags) uintOpt(p *uint, meta *FlagMeta) {
	defVal := meta.DValue().Int()
	name := fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	fs.fSet.UintVar(p, name, uint(defVal), meta.Desc)
	meta.flag = fs.fSet.Lookup(name)
}

// Uint64 binding an int option flag, return pointer
func (fs *Flags) Uint64(name, shorts string, defVal uint64, desc string) *uint64 {
	meta := newFlagOpt(name, desc, defVal, shorts)
	name = fs.checkFlagInfo(meta)

	p := fs.fSet.Uint64(name, defVal, meta.Desc)
	meta.flag = fs.fSet.Lookup(name)

	return p
}

// Uint64Var binding an uint option flag
func (fs *Flags) Uint64Var(p *uint64, meta *FlagMeta) { fs.uint64Opt(p, meta) }

// Uint64Opt binding an uint64 option
func (fs *Flags) Uint64Opt(p *uint64, name, shorts string, defVal uint64, desc string) {
	fs.uint64Opt(p, newFlagOpt(name, desc, defVal, shorts))
}

// binding option and shorts
func (fs *Flags) uint64Opt(p *uint64, meta *FlagMeta) {
	defVal := meta.DValue().Int64()
	name := fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	fs.fSet.Uint64Var(p, name, uint64(defVal), meta.Desc)
	meta.flag = fs.fSet.Lookup(name)
}

// Var binding an custom var option flag
func (fs *Flags) Var(p flag.Value, meta *FlagMeta) { fs.varOpt(p, meta) }

// VarOpt binding a custom var option
//
// Usage:
//
//	var names gcli.Strings
//	cmd.VarOpt(&names, "tables", "t", "description ...")
func (fs *Flags) VarOpt(p flag.Value, name, shorts, desc string) {
	fs.varOpt(p, newFlagOpt(name, desc, nil, shorts))
}

// binding option and shorts
func (fs *Flags) varOpt(p flag.Value, meta *Option) {
	name := fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	fs.fSet.Var(p, name, meta.Desc)
	meta.flag = fs.fSet.Lookup(name)
}

// Required flag option name(s)
func (fs *Flags) Required(names ...string) {
	for _, name := range names {
		meta, ok := fs.metas[name]
		if !ok {
			helper.Panicf("undefined option flag '%s'", name)
		}
		meta.Required = true
	}
}

// check flag option name and short-names
func (fs *Flags) checkFlagInfo(meta *Option) string {
	// NOTICE: must init some required fields
	if fs.names == nil {
		fs.names = map[string]int{}
		fs.metas = map[string]*Option{}
		fs.InitFlagSet("flags-" + meta.Name)
	}

	// check flag name
	name := meta.initCheck()
	if _, ok := fs.metas[name]; ok {
		helper.Panicf("redefined option flag '%s'", name)
	}

	// is a short name
	helpLen := meta.helpNameLen()
	// fix: must exclude Hidden option
	if !meta.Hidden {
		fs.flagMaxLen = mathutil.MaxInt(fs.flagMaxLen, helpLen)
	}

	// check short names
	fs.checkShortNames(name, meta.Shorts)

	// update name length
	fs.names[name] = helpLen
	// storage meta and name
	fs.metas[name] = meta
	return name
}

// check short names
func (fs *Flags) checkShortNames(name string, shorts []string) {
	if len(shorts) == 0 {
		return
	}

	fs.existShort = true
	if fs.shorts == nil {
		fs.shorts = map[string]string{}
	}

	for _, short := range shorts {
		if name == short {
			helper.Panicf("short name '%s' has been used as the current option name", short)
		}

		if _, ok := fs.names[short]; ok {
			helper.Panicf("short name '%s' has been used as an option name", short)
		}

		if n, ok := fs.shorts[short]; ok {
			helper.Panicf("short name '%s' has been used by option '%s'", short, n)
		}

		// storage short name
		fs.shorts[short] = name
	}

}

/***********************************************************************
 * Flags:
 * - render help message
 ***********************************************************************/

// SetHelpRender set the raw *flag.FlagSet.Usage
func (fs *Flags) SetHelpRender(fn func()) {
	fs.fSet.Usage = fn
}

// PrintHelpPanel for all options to the gf.out
func (fs *Flags) PrintHelpPanel() {
	color.Fprint(fs.out, fs.String())
}

// String for all flag options
func (fs *Flags) String() string {
	return fs.BuildHelp()
}

// BuildHelp string for all flag options
func (fs *Flags) BuildHelp() string {
	if fs.buf == nil {
		fs.buf = new(bytes.Buffer)
	}

	// repeat call the method
	if fs.buf.Len() < 1 {
		fs.buf.WriteString("Options:\n")
		fs.buf.WriteString(fs.BuildOptsHelp())
		fs.buf.WriteByte('\n')

		if fs.HasArgs() {
			fs.buf.WriteString("Arguments:\n")
			fs.buf.WriteString(fs.BuildArgsHelp())
			fs.buf.WriteByte('\n')
		}
	}

	return fs.buf.String()
}

// BuildOptsHelp string.
func (fs *Flags) BuildOptsHelp() string {
	var sb strings.Builder

	fs.FSet().VisitAll(func(f *flag.Flag) {
		sb.WriteString(fs.formatOneFlag(f))
		sb.WriteByte('\n')
	})

	return sb.String()
}

func (fs *Flags) formatOneFlag(f *flag.Flag) (s string) {
	// Skip render:
	// - meta is not exists(Has ensured that it is not a short name)
	// - it is hidden flag option
	// - flag desc is empty
	meta, has := fs.metas[f.Name]
	if !has || meta.Hidden {
		return
	}

	var fullName string
	name := f.Name
	// eg: "-V, --version" length is: 13
	nameLen := fs.names[name]
	// display description on new line
	descNl := fs.cfg.DescNewline

	var nlIndent string
	if descNl {
		nlIndent = "\n        "
	} else {
		nlIndent = "\n      " + strings.Repeat(" ", fs.flagMaxLen)
	}

	// add prefix '-' to option
	fullName = cflag.AddPrefixes2(name, meta.Shorts, true)
	s = fmt.Sprintf("  <info>%s</>", fullName)

	// - build flag type info
	typeName, desc := flag.UnquoteUsage(f)
	// typeName: option value data type: int, string, ..., bool value will return ""
	if !fs.cfg.WithoutType && len(typeName) > 0 {
		typeLen := len(typeName) + 1
		if !descNl && nameLen+typeLen > fs.flagMaxLen {
			descNl = true
		} else {
			nameLen += typeLen
		}

		s += fmt.Sprintf(" <magenta>%s</>", typeName)
	}

	if descNl {
		s += nlIndent
	} else {
		// padding space to flagMaxLen width.
		if padLen := fs.flagMaxLen - nameLen; padLen > 0 {
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

	s += getRequiredMark(meta.Required) + desc

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

// IterAll Iteration all flag options with metadata
func (fs *Flags) IterAll(fn func(f *flag.Flag, meta *FlagMeta)) {
	fs.FSet().VisitAll(func(f *flag.Flag) {
		if _, ok := fs.metas[f.Name]; ok {
			fn(f, fs.metas[f.Name])
		}
	})
}

// ShortNames get all short-names of the option
func (fs *Flags) ShortNames(name string) (ss []string) {
	if opt, ok := fs.metas[name]; ok {
		ss = opt.Shorts
	}
	return
}

// IsShortOpt alias of the IsShortcut()
func (fs *Flags) IsShortOpt(short string) bool { return fs.IsShortName(short) }

// IsShortName check it is a shortcut name
func (fs *Flags) IsShortName(short string) bool {
	if len(short) != 1 {
		return false
	}

	_, ok := fs.shorts[short]
	return ok
}

// IsOption check it is a option name
func (fs *Flags) IsOption(name string) bool { return fs.HasOption(name) }

// HasOption check it is a option name
func (fs *Flags) HasOption(name string) bool {
	_, ok := fs.names[name]
	return ok
}

// HasFlag check it is a option name. alias of HasOption()
func (fs *Flags) HasFlag(name string) bool {
	_, ok := fs.names[name]
	return ok
}

// HasFlagMeta check it is has FlagMeta
func (fs *Flags) HasFlagMeta(name string) bool {
	_, ok := fs.metas[name]
	return ok
}

// LookupFlag get flag.Flag by name
func (fs *Flags) LookupFlag(name string) *flag.Flag { return fs.fSet.Lookup(name) }

// FlagMeta get FlagMeta by name
func (fs *Flags) FlagMeta(name string) *FlagMeta { return fs.metas[name] }

// Metas get all flag metas
func (fs *Flags) Metas() map[string]*FlagMeta { return fs.metas }

// Opts get all flag options
func (fs *Flags) Opts() map[string]*Option { return fs.metas }

// Hidden there are given option names
// func (gf *Flags) Hidden(names ...string) {
// 	for _, name := range names {
// 		if !gf.HasOption(name) { // not registered
// 			continue
// 		}
//
// 		gf.metas[name].Hidden = true
// 	}
// }

// Name of the Flags
func (fs *Flags) Name() string { return fs.fSet.Name() }

// Len of the Flags
func (fs *Flags) Len() int { return len(fs.names) }

// FSet get the raw *flag.FlagSet
func (fs *Flags) FSet() *flag.FlagSet { return fs.fSet }

// SetFlagSet set the raw *flag.FlagSet
func (fs *Flags) SetFlagSet(fSet *flag.FlagSet) { fs.fSet = fSet }

// SetOutput for the Flags
func (fs *Flags) SetOutput(out io.Writer) { fs.out = out }

// FlagNames return all option names
func (fs *Flags) FlagNames() map[string]int { return fs.names }

// RawArg get an argument value by index
func (fs *Flags) RawArg(i int) string { return fs.fSet.Arg(i) }

// RawArgs get all raw arguments.
// if have been called parse, the return is remaining args.
func (fs *Flags) RawArgs() []string { return fs.fSet.Args() }

// FSetArgs get all raw arguments. alias of the RawArgs()
// if have been called parse, the return is remaining args.
func (fs *Flags) FSetArgs() []string { return fs.fSet.Args() }
