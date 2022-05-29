package gcli

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"unsafe"

	"github.com/gookit/color"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/stdutil"
	"github.com/gookit/goutil/strutil"
)

// The options text alignment type
// - Align right, padding left
// - Align left, padding right
const (
	AlignLeft  = strutil.PosRight
	AlignRight = strutil.PosLeft

	// default desc
	defaultDesc = "No description"

	// TagRuleNamed struct tag use named k-v rule.
	// eg: `flag:"name=int0;shorts=i;required=true;desc=int option message"`
	TagRuleNamed = 0
	// TagRuleSimple struct tag use simple rule.
	// format: "desc;required;default;shorts"
	// eg: `flag:"int option message;required;;i"`
	TagRuleSimple = 1
)

var (
	// FlagTagName default tag name on struct
	FlagTagName = "flag"
	// allowed keys on struct tag.
	flagTagKeys = arrutil.Strings{"name", "shorts", "desc", "default", "required"}
)

// FlagsOption for render help information
type FlagsOption struct {
	// WithoutType don't display flag data type on print help
	WithoutType bool
	// NameDescOL flag and desc at one line on print help
	NameDescOL bool
	// Alignment flag name align left or right. default is: right
	Alignment uint8
	// TagName on struct
	TagName string
	// TagRuleType for struct tag value. default is TagRuleNamed
	TagRuleType uint8
}

// Flags struct definition
type Flags struct {
	// Desc message
	Desc string
	// ExitFunc for handle exit
	ExitFunc func(code int)
	// FlagsOption option for render help message
	opt *FlagsOption
	// raw flag set
	fSet *flag.FlagSet
	// buf for build help message
	buf *bytes.Buffer
	// output for print help message
	out io.Writer
	// all option names of the command. {name: length} // TODO delete, move len to meta.
	names map[string]int
	// metadata for all options
	metas map[string]*FlagMeta // TODO support option category
	// short names for options. format: {short:name}
	// eg. {"n": "name", "o": "opt"}
	shorts map[string]string
	// mapping for name to shortcut {"name": {"n", "m"}}
	name2shorts map[string][]string
	// flag name max length. useful for render help
	// eg: "-V, --version" length is 13
	flagMaxLen int
	// exist short names. useful for render help
	existShort bool
}

func newDefaultFlagOption() *FlagsOption {
	return &FlagsOption{
		Alignment: AlignRight,
		TagName:   FlagTagName,
	}
}

// NewFlags create an new Flags
func NewFlags(nameWithDesc ...string) *Flags {
	fs := &Flags{
		out: os.Stdout,
		opt: newDefaultFlagOption(),
	}
	fs.ExitFunc = os.Exit

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

var flagValueType = reflect.TypeOf(new(flag.Value)).Elem()
var errNotPtrValue = errors.New("must provide an ptr value")
var errNotAnStruct = errors.New("must provide an struct ptr")
var errTagRuleType = errors.New("invalid tag rule type on struct")

// FromStruct from struct tag binding options
func (fs *Flags) FromStruct(ptr interface{}) error {
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

	tagName := fs.opt.TagName
	if tagName == "" {
		tagName = FlagTagName
	}

	var mp map[string]string
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		name := sf.Name

		// skip cannot exported field
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

		if fs.opt.TagRuleType == TagRuleNamed {
			mp = parseNamedRule(name, str)
		} else if fs.opt.TagRuleType == TagRuleSimple {
			mp = parseSimpleRule(name, str)
		} else {
			return errTagRuleType
		}

		// for create flag meta
		defVal := mp["default"]
		shorts := splitShortcut(mp["shorts"])
		optName, has := mp["name"]
		if !has { // use field as option name.
			optName = strutil.SnakeCase(name, "-")
		}

		meta := newFlagMeta(optName, mp["desc"], defVal, shorts)
		if must, has := mp["required"]; has {
			meta.Required = strutil.MustBool(must)
		}

		// field is implements flag.Value
		if ft.Implements(flagValueType) {
			// use assert
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

// FromText from text desc binding options
// func (fs *Flags) FromText(s string) error {
// 	return nil
// }

// SetOptions for the object.
func (fs *Flags) SetOptions(opt *FlagsOption) {
	fs.opt = opt
}

// UseSimpleRule for the parse tag value rule string. see TagRuleSimple
func (fs *Flags) UseSimpleRule() *Flags {
	fs.opt.TagRuleType = TagRuleSimple
	return fs
}

// WithOptions for the object.
func (fs *Flags) WithOptions(fns func(opt *FlagsOption)) *Flags {
	fns(fs.opt)
	return fs
}

// InitFlagSet create and init flag.FlagSet
func (fs *Flags) InitFlagSet(name string) {
	if fs.fSet != nil {
		return
	}

	fs.fSet = flag.NewFlagSet(name, flag.ContinueOnError)
	// disable output internal error message on parse flags
	fs.fSet.SetOutput(ioutil.Discard)
	// nothing to do ... render usage on after parsed
	fs.fSet.Usage = func() {}
}

// Run flags parse and handle help render
// Usage:
// 	gf := gcli.NewFlags()
//  ...
// 	gf.Run(os.Args)
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
			fs.ExitFunc(0)
		} else {
			color.Errorf("flag parse error - %s", err.Error())
			fs.ExitFunc(2)
		}
	}
}

// Parse given arguments
//
// Usage:
// 	gf := gcli.NewFlags()
// 	gf.BoolOpt(&debug, "debug", "", defDebug, "open debug mode")
// 	gf.UintOpt(&port, "port", "p", 18081, "the http server port")
//
// 	err := gf.Parse(os.Args[1:])
func (fs *Flags) Parse(args []string) (err error) {
	err = fs.fSet.Parse(args)
	if err != nil {
		return
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

// RawArg get an argument value by index
func (fs *Flags) RawArg(i int) string {
	return fs.fSet.Arg(i)
}

// RawArgs get all raw arguments.
// if have been called parse, the return is remaining args.
func (fs *Flags) RawArgs() []string {
	return fs.fSet.Args()
}

// FSetArgs get all raw arguments. alias of the RawArgs()
// if have been called parse, the return is remaining args.
func (fs *Flags) FSetArgs() []string {
	return fs.fSet.Args()
}

/***********************************************************************
 * Flags:
 * - binding option var
 ***********************************************************************/

// --- bool option

// Bool binding an bool option flag, return pointer
func (fs *Flags) Bool(name, shorts string, defValue bool, desc string) *bool {
	meta := newFlagMeta(name, desc, defValue, splitShortcut(shorts))
	name = fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	p := fs.fSet.Bool(name, defValue, meta.Desc)

	// binding all short name options to flag.FlagSet
	for _, s := range meta.Shorts {
		fs.fSet.BoolVar(p, s, defValue, "") // dont add description for short name
	}
	return p
}

// BoolVar binding an bool option flag
func (fs *Flags) BoolVar(p *bool, meta *FlagMeta) {
	fs.boolOpt(p, meta)
}

// BoolOpt binding an bool option
func (fs *Flags) BoolOpt(p *bool, name, shorts string, defValue bool, desc string) {
	fs.boolOpt(p, newFlagMeta(name, desc, defValue, splitShortcut(shorts)))
}

// binding option and shorts
func (fs *Flags) boolOpt(p *bool, meta *FlagMeta) {
	defVal := meta.DValue().Bool()
	fmtName := fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	fs.fSet.BoolVar(p, fmtName, defVal, meta.Desc)

	// binding all short name options to flag.FlagSet
	for _, s := range meta.Shorts {
		fs.fSet.BoolVar(p, s, defVal, "") // dont add description for short name
	}
}

// --- float option

// Float64Var binding an float64 option flag
func (fs *Flags) Float64Var(p *float64, meta *FlagMeta) {
	fs.float64Opt(p, meta)
}

// Float64Opt binding an float64 option
func (fs *Flags) Float64Opt(p *float64, name, shorts string, defValue float64, desc string) {
	fs.float64Opt(p, newFlagMeta(name, desc, defValue, splitShortcut(shorts)))
}

func (fs *Flags) float64Opt(p *float64, meta *FlagMeta) {
	defValue := meta.DValue().Float64()
	fmtName := fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	fs.fSet.Float64Var(p, fmtName, defValue, meta.Desc)

	// binding all short name options to flag.FlagSet
	for _, s := range meta.Shorts {
		fs.fSet.Float64Var(p, s, defValue, "") // dont add description for short name
	}
}

// --- string option

// Str binding an string option flag, return pointer
func (fs *Flags) Str(name, shorts string, defValue, desc string) *string {
	meta := newFlagMeta(name, desc, defValue, splitShortcut(shorts))
	name = fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	p := fs.fSet.String(name, defValue, meta.Desc)

	// binding all short name options to flag.FlagSet
	for _, s := range meta.Shorts {
		fs.fSet.StringVar(p, s, defValue, "") // dont add description for short name
	}
	return p
}

// StrVar binding an string option flag
func (fs *Flags) StrVar(p *string, meta *FlagMeta) {
	fs.strOpt(p, meta)
}

// StrOpt binding an string option
func (fs *Flags) StrOpt(p *string, name, shorts, defValue, desc string) {
	fs.strOpt(p, newFlagMeta(name, desc, defValue, splitShortcut(shorts)))
}

// binding option and shorts
func (fs *Flags) strOpt(p *string, meta *FlagMeta) {
	defValue := meta.DValue().String()
	fmtName := fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	fs.fSet.StringVar(p, fmtName, defValue, meta.Desc)

	// binding all short name options to flag.FlagSet
	for _, s := range meta.Shorts {
		fs.fSet.StringVar(p, s, defValue, "") // dont add description for short name
	}
}

// --- intX option

// Int binding an int option flag, return pointer
func (fs *Flags) Int(name, shorts string, defValue int, desc string) *int {
	meta := newFlagMeta(name, desc, defValue, splitShortcut(shorts))
	name = fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	p := fs.fSet.Int(name, defValue, meta.Desc)

	// binding all short name options to flag.FlagSet
	for _, s := range meta.Shorts {
		fs.fSet.IntVar(p, s, defValue, "") // dont add description for short name
	}
	return p
}

// IntVar binding an int option flag
func (fs *Flags) IntVar(p *int, meta *FlagMeta) {
	fs.intOpt(p, meta)
}

// IntOpt binding an int option
func (fs *Flags) IntOpt(p *int, name, shorts string, defValue int, desc string) {
	fs.intOpt(p, newFlagMeta(name, desc, defValue, splitShortcut(shorts)))
}

func (fs *Flags) intOpt(p *int, meta *FlagMeta) {
	defValue := meta.DValue().Int()
	fmtName := fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	fs.fSet.IntVar(p, fmtName, defValue, meta.Desc)

	// binding all short name options to flag.FlagSet
	for _, s := range meta.Shorts {
		fs.fSet.IntVar(p, s, defValue, "") // dont add description for short name
	}
}

// Str binding an int64 option flag, return pointer
func (fs *Flags) Int64(name, shorts string, defValue int64, desc string) *int64 {
	meta := newFlagMeta(name, desc, defValue, splitShortcut(shorts))
	name = fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	p := fs.fSet.Int64(name, defValue, meta.Desc)

	// binding all short name options to flag.FlagSet
	for _, s := range meta.Shorts {
		fs.fSet.Int64Var(p, s, defValue, "") // dont add description for short name
	}
	return p
}

// Int64Var binding an int64 option flag
func (fs *Flags) Int64Var(p *int64, meta *FlagMeta) {
	fs.int64Opt(p, meta)
}

// Int64Opt binding an int64 option
func (fs *Flags) Int64Opt(p *int64, name, shorts string, defValue int64, desc string) {
	fs.int64Opt(p, newFlagMeta(name, desc, defValue, splitShortcut(shorts)))
}

func (fs *Flags) int64Opt(p *int64, meta *FlagMeta) {
	defValue := meta.DValue().Int64()
	fmtName := fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	fs.fSet.Int64Var(p, fmtName, defValue, meta.Desc)

	// binding all short options to flag.FlagSet
	for _, s := range meta.Shorts {
		fs.fSet.Int64Var(p, s, defValue, "") // dont add description for short name
	}
}

// --- uintX option

// Uint binding an int option flag, return pointer
func (fs *Flags) Uint(name, shorts string, defValue uint, desc string) *uint {
	meta := newFlagMeta(name, desc, defValue, splitShortcut(shorts))
	name = fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	p := fs.fSet.Uint(name, defValue, meta.Desc)

	// binding all short name options to flag.FlagSet
	for _, s := range meta.Shorts {
		fs.fSet.UintVar(p, s, defValue, "") // dont add description for short name
	}
	return p
}

// UintVar binding an uint option flag
func (fs *Flags) UintVar(p *uint, meta *FlagMeta) {
	fs.uintOpt(p, meta)
}

// UintOpt binding an uint option
func (fs *Flags) UintOpt(p *uint, name, shorts string, defValue uint, desc string) {
	fs.uintOpt(p, newFlagMeta(name, desc, defValue, splitShortcut(shorts)))
}

func (fs *Flags) uintOpt(p *uint, meta *FlagMeta) {
	defValue := meta.DValue().Int()
	fmtName := fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	fs.fSet.UintVar(p, fmtName, uint(defValue), meta.Desc)

	// binding all short options to flag.FlagSet
	for _, s := range meta.Shorts {
		fs.fSet.UintVar(p, s, uint(defValue), "") // dont add description for short name
	}
}

// Uint binding an int option flag, return pointer
func (fs *Flags) Uint64(name, shorts string, defValue uint64, desc string) *uint64 {
	meta := newFlagMeta(name, desc, defValue, splitShortcut(shorts))
	name = fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	p := fs.fSet.Uint64(name, defValue, meta.Desc)

	// binding all short name options to flag.FlagSet
	for _, s := range meta.Shorts {
		fs.fSet.Uint64Var(p, s, defValue, "") // dont add description for short name
	}
	return p
}

// Uint64Var binding an uint option flag
func (fs *Flags) Uint64Var(p *uint64, meta *FlagMeta) {
	fs.uint64Opt(p, meta)
}

// Uint64Opt binding an uint64 option
func (fs *Flags) Uint64Opt(p *uint64, name, shorts string, defValue uint64, desc string) {
	fs.uint64Opt(p, newFlagMeta(name, desc, defValue, splitShortcut(shorts)))
}

// binding option and shorts
func (fs *Flags) uint64Opt(p *uint64, meta *FlagMeta) {
	defValue := meta.DValue().Int64()
	fmtName := fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	fs.fSet.Uint64Var(p, fmtName, uint64(defValue), meta.Desc)

	// binding all short options to flag.FlagSet
	for _, s := range meta.Shorts {
		fs.fSet.Uint64Var(p, s, uint64(defValue), "") // dont add description for short name
	}
}

// Var binding an custom var option flag
func (fs *Flags) Var(p flag.Value, meta *FlagMeta) {
	fs.varOpt(p, meta)
}

// VarOpt binding a custom var option
// Usage:
//		var names gcli.Strings
// 		cmd.VarOpt(&names, "tables", "t", "description ...")
func (fs *Flags) VarOpt(p flag.Value, name, shorts, desc string) {
	fs.varOpt(p, newFlagMeta(name, desc, nil, splitShortcut(shorts)))
}

// binding option and shorts
func (fs *Flags) varOpt(p flag.Value, meta *FlagMeta) {
	fmtName := fs.checkFlagInfo(meta)

	// binding option to flag.FlagSet
	fs.fSet.Var(p, fmtName, meta.Desc)

	// binding all short options to flag.FlagSet
	for _, s := range meta.Shorts {
		fs.fSet.Var(p, s, "") // dont add description for short name
	}
}

// Required flag option name(s)
func (fs *Flags) Required(names ...string) {
	for _, name := range names {
		meta, ok := fs.metas[name]
		if !ok {
			panicf("undefined option flag '%s'", name)
		}

		meta.Required = true
	}
}

// check flag option name and short-names
func (fs *Flags) checkFlagInfo(meta *FlagMeta) string {
	// NOTICE: must init some required fields
	if fs.names == nil {
		fs.names = map[string]int{}
		fs.metas = map[string]*FlagMeta{}
		fs.InitFlagSet("flags")
	}

	// check flag name
	name := meta.goodName()
	if _, ok := fs.metas[name]; ok {
		panicf("redefined option flag '%s'", name)
	}

	nameLength := len(name)
	// is an short name
	if nameLength == 1 {
		nameLength += 1 // prefix: "-"
		fs.existShort = true
	} else {
		nameLength += 2 // prefix: "--"
	}

	// fix: must exclude Hidden option
	if !meta.Hidden && fs.flagMaxLen < nameLength {
		fs.flagMaxLen = nameLength
	}

	// check and format short names
	meta.Shorts = fs.checkShortNames(name, nameLength, meta.Shorts)

	// storage meta and name
	fs.metas[name] = meta
	return name
}

// check short names
func (fs *Flags) checkShortNames(name string, nameLength int, shorts []string) []string {
	if len(shorts) == 0 {
		// record name without shorts
		fs.names[name] = nameLength
		return shorts
	}

	// init fs.shorts and fs.name2shorts
	if fs.shorts == nil {
		fs.shorts = map[string]string{}
		fs.name2shorts = map[string][]string{}
	}

	var fmtShorts []string
	for _, short := range shorts {
		short = strings.Trim(short, "- ")
		if short == "" {
			continue
		}

		// ensure it is one char
		char := short[0]
		short = string(char)
		// A 65 -> Z 90, a 97 -> z 122
		if !strutil.IsAlphabet(char) {
			panicf("short name only allow: A-Za-z given: '%s'", short)
		}

		if name == short {
			panicf("short name '%s' has been used as the current option name", short)
		}

		if _, ok := fs.names[short]; ok {
			panicf("short name '%s' has been used as an option name", short)
		}

		if n, ok := fs.shorts[short]; ok {
			panicf("short name '%s' has been used by option '%s'", short, n)
		}

		fmtShorts = append(fmtShorts, short)
		// storage short name
		fs.shorts[short] = name
	}

	// one short = '-' + 'x' + ',' + ' '
	// eg: "-o, " len=4
	// eg: "-o, -a, " len=8
	nameLength += 4 * len(fmtShorts)
	fs.existShort = true

	// update name length
	fs.names[name] = nameLength
	if fs.flagMaxLen < nameLength {
		fs.flagMaxLen = nameLength
	}

	fs.name2shorts[name] = fmtShorts
	return fmtShorts
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
	if fs.buf == nil {
		fs.buf = new(bytes.Buffer)
	}

	// repeat call the method
	if fs.buf.Len() < 1 {
		if fs.existShort { // add 4 space prefix for flag
			fs.flagMaxLen += 4
		}

		// refer fs.Fs().PrintDefaults()
		fs.FSet().VisitAll(fs.formatOneFlag)
	}

	return fs.buf.String()
}

func (fs *Flags) formatOneFlag(f *flag.Flag) {
	// Skip render:
	// - meta is not exists(Has ensured that it is not a short name)
	// - it is hidden flag option
	// - flag desc is empty
	meta, has := fs.metas[f.Name]
	if !has || meta.Hidden {
		return
	}

	var s, fullName string
	name := f.Name
	// eg: "-V, --version" length is: 13
	fLen := fs.names[name]

	// - build flag name info
	// is long option
	if len(name) > 1 {
		// find shorts
		if shorts := fs.ShortNames(name); len(shorts) > 0 {
			fullName = fmt.Sprintf("%s, --%s", shorts2string(shorts), name)
		} else {
			fullName = "--" + name
			// if has short names. add 4 space. eg: "-s, "
			if fs.existShort {
				fullName = "    " + fullName
			}
		}
	} else {
		// only short option
		// s = fmt.Sprintf("  <info>-%s</>", name)
		fullName = "-" + name
	}

	// fs.NameDescOL = true: padding space to same width.
	if fs.opt.NameDescOL {
		fullName = strutil.Padding(fullName, " ", fs.flagMaxLen, fs.opt.Alignment)
	}

	s = fmt.Sprintf("  <info>%s</>", fullName)

	// - build flag type info
	typeName, desc := flag.UnquoteUsage(f)
	// typeName: option value data type: int, string, ..., bool value will return ""
	if fs.opt.WithoutType == false && len(typeName) > 0 {
		s += fmt.Sprintf(" <magenta>%s</>", typeName)
	}

	// - flag and description at one line
	// - Boolean flags of one ASCII letter are so common we
	// treat them specially, putting their usage on the same line.
	if fs.opt.NameDescOL || (typeName == "" && fLen <= 4) { // space, space, '-', 'x'.
		s += "    "
	} else {
		// display description on new line
		s += "\n        "
	}

	// --- build description
	if desc == "" {
		desc = defaultDesc
	}

	// flag is required
	if meta.Required {
		s += "<red>*</>"
	}

	s += strings.Replace(strutil.UpperFirst(desc), "\n", "\n        ", -1)

	// ---- append default value
	if !isZeroValue(f, f.DefValue) {
		if _, ok := f.Value.(*stringValue); ok {
			// put quotes on the value
			s += fmt.Sprintf(" (default <magentaB>%q</>)", f.DefValue)
		} else {
			s += fmt.Sprintf(" (default <magentaB>%v</>)", f.DefValue)
		}
	}

	// save to buffer
	fs.buf.WriteString(s)
	fs.buf.WriteByte('\n')
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
	if len(fs.name2shorts) == 0 {
		return
	}

	return fs.name2shorts[name]
}

// IsShortOpt alias of the IsShortcut()
func (fs *Flags) IsShortOpt(short string) bool {
	return fs.IsShortName(short)
}

// IsShortcut check it is a shortcut name
func (fs *Flags) IsShortName(short string) bool {
	if len(short) != 1 {
		return false
	}

	_, ok := fs.shorts[short]
	return ok
}

// IsOption check it is a option name
func (fs *Flags) IsOption(name string) bool {
	return fs.HasOption(name)
}

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
func (fs *Flags) LookupFlag(name string) *flag.Flag {
	return fs.fSet.Lookup(name)
}

// FlagMeta get FlagMeta by name
func (fs *Flags) FlagMeta(name string) *FlagMeta {
	return fs.metas[name]
}

// Metas get all flag metas
func (fs *Flags) Metas() map[string]*FlagMeta {
	return fs.metas
}

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
func (fs *Flags) Name() string {
	return fs.fSet.Name()
}

// Len of the Flags
func (fs *Flags) Len() int {
	return len(fs.names)
}

// FSet get the raw *flag.FlagSet
func (fs *Flags) FSet() *flag.FlagSet {
	return fs.fSet
}

// SetFlagSet set the raw *flag.FlagSet
func (fs *Flags) SetFlagSet(fSet *flag.FlagSet) {
	fs.fSet = fSet
}

// SetOutput for the Flags
func (fs *Flags) SetOutput(out io.Writer) {
	fs.out = out
}

// FlagNames return all option names
func (fs *Flags) FlagNames() map[string]int {
	return fs.names
}

/***********************************************************************
 * Flags:
 * - flag metadata
 ***********************************************************************/

// FlagMeta for an flag(option/argument)
type FlagMeta struct {
	// varPtr interface{}
	// name and description
	Name, Desc string
	// Alias of the name. isn't shorts. eg: name='dry-run' alias='dr' TODO
	Alias string
	// default value for the flag option
	DefVal interface{}
	// wrapped the default value
	defVal *goutil.Value
	// short names. eg: ["o", "a"]
	Shorts []string
	// advanced settings
	// hidden the option on help
	Hidden bool
	// the option is required
	Required bool
	// Validator support validate the option flag value
	Validator func(val string) error
}

// newFlagMeta quick create an FlagMeta
func newFlagMeta(name, desc string, defVal interface{}, shorts []string) *FlagMeta {
	return &FlagMeta{
		Name: name,
		Desc: desc,
		// other info
		DefVal: defVal,
		Shorts: shorts,
	}
}

// Shorts2String join shorts to an string
func (m *FlagMeta) Shorts2String(sep ...string) string {
	if len(m.Shorts) == 0 {
		return ""
	}

	char := ","
	if len(sep) > 0 {
		char = sep[0]
	}

	return strings.Join(m.Shorts, char)
}

// Validate the binding value
func (m *FlagMeta) Validate(val string) error {
	// check required
	if m.Required && val == "" {
		return fmt.Errorf("flag '%s' is required", m.Name)
	}

	// call user custom validator
	if m.Validator != nil {
		return m.Validator(val)
	}

	return nil
}

// DValue wrap the default value
func (m *FlagMeta) DValue() *goutil.Value {
	if m.defVal == nil {
		m.defVal = &stdutil.Value{V: m.DefVal}
	}

	return m.defVal
}

// good name of the flag
func (m *FlagMeta) goodName() string {
	name := strings.Trim(m.Name, "- ")
	if name == "" {
		panicf("option flag name cannot be empty")
	}

	if !goodName.MatchString(name) {
		panicf("option flag name '%s' is invalid, must match: %s", name, regGoodName)
	}

	// update self name
	m.Name = name
	return name
}

/***********************************************************************
 * Flags:
 * - flag value check methods
 ***********************************************************************/

// isZeroValue guesses whether the string represents the zero
// value for a flag. It is not accurate but in practice works OK.
//
// NOTICE: the func is copied from package 'flag', func 'isZeroValue'
func isZeroValue(fg *flag.Flag, value string) bool {
	// Build a zero value of the flag's Value type, and see if the
	// result of calling its String method equals the value passed in.
	// This works unless the Value type is itself an interface type.
	typ := reflect.TypeOf(fg.Value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	if value == z.Interface().(flag.Value).String() {
		return true
	}

	switch value {
	case "false", "", "0":
		return true
	}
	return false
}

// -- string Value
// NOTICE: the var is copied from package 'flag'
type stringValue string

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}
func (s *stringValue) Get() interface{} { return string(*s) }
func (s *stringValue) String() string   { return string(*s) }
