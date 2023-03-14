package gflag

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/gookit/gcli/v3/helper"
	"github.com/gookit/goutil/cflag"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/stdutil"
	"github.com/gookit/goutil/structs"
	"github.com/gookit/goutil/strutil"
)

const (
	shortSepRune = ','
	shortSepChar = ","
)

// DefaultOptWidth for render help
var DefaultOptWidth = 20

// CliOpts cli options management
type CliOpts struct {
	// name inherited from gcli.Command
	name string

	// the options flag set TODO remove flag.FlagSet, custom implement parse
	fSet *flag.FlagSet
	// all cli option names.
	// format: {name: length} // TODO delete, move len to opts.
	names map[string]int
	// metadata for all options
	opts map[string]*CliOpt // TODO support option category
	// short names map for options. format: {short: name}
	// eg. {"n": "name", "o": "opt"}
	shorts map[string]string
	// support option category
	categories []OptCategory
	// flag name max length. useful for render help
	// eg: "-V, --version" length is 13
	optMaxLen int
	// exist short names. useful for render help
	hasShort bool
}

// InitFlagSet create and init flag.FlagSet
func (ops *CliOpts) InitFlagSet(name string) {
	if ops.fSet != nil {
		return
	}

	ops.name = name
	ops.fSet = flag.NewFlagSet(name, flag.ContinueOnError)
	// disable output internal error message on parse flags
	ops.fSet.SetOutput(io.Discard)
	// nothing to do ... render usage on after parsed
	ops.fSet.Usage = func() {}
	ops.optMaxLen = DefaultOptWidth
}

// SetName for CliArgs
func (ops *CliOpts) SetName(name string) {
	ops.name = name
}

/***********************************************************************
 * Options:
 * - binding option var
 ***********************************************************************/

// --- bool option

// Bool binding a bool option flag, return pointer
func (ops *CliOpts) Bool(name, shorts string, defVal bool, desc string) *bool {
	opt := newOpt(name, desc, defVal, shorts)
	name = ops.checkFlagInfo(opt)

	// binding option to flag.FlagSet
	ptr := ops.fSet.Bool(name, defVal, opt.Desc)
	opt.flag = ops.fSet.Lookup(name)

	return ptr
}

// BoolVar binding a bool option flag
func (ops *CliOpts) BoolVar(ptr *bool, opt *CliOpt) { ops.boolOpt(ptr, opt) }

// BoolOpt binding a bool option
func (ops *CliOpts) BoolOpt(ptr *bool, name, shorts string, defVal bool, desc string) {
	ops.boolOpt(ptr, newOpt(name, desc, defVal, shorts))
}

// BoolOpt2 binding a bool option, and allow with CliOptFn for config option.
func (ops *CliOpts) BoolOpt2(p *bool, nameAndShorts, desc string, setFns ...CliOptFn) {
	ops.boolOpt(p, NewOpt(nameAndShorts, desc, false, setFns...))
}

// binding option and shorts
func (ops *CliOpts) boolOpt(ptr *bool, opt *CliOpt) {
	defVal := opt.DValue().Bool()
	name := ops.checkFlagInfo(opt)

	// binding option to flag.FlagSet
	ops.fSet.BoolVar(ptr, name, defVal, opt.Desc)
	opt.flag = ops.fSet.Lookup(name)
}

// --- float option

// Float64Var binding an float64 option flag
func (ops *CliOpts) Float64Var(ptr *float64, opt *CliOpt) { ops.float64Opt(ptr, opt) }

// Float64Opt binding a float64 option
func (ops *CliOpts) Float64Opt(p *float64, name, shorts string, defVal float64, desc string) {
	ops.float64Opt(p, newOpt(name, desc, defVal, shorts))
}

func (ops *CliOpts) float64Opt(p *float64, opt *CliOpt) {
	defVal := opt.DValue().Float64()
	name := ops.checkFlagInfo(opt)

	// binding option to flag.FlagSet
	ops.fSet.Float64Var(p, name, defVal, opt.Desc)
	opt.flag = ops.fSet.Lookup(name)
}

// --- string option

// Str binding an string option flag, return pointer
func (ops *CliOpts) Str(name, shorts string, defVal, desc string) *string {
	opt := newOpt(name, desc, defVal, shorts)
	name = ops.checkFlagInfo(opt)

	// binding option to flag.FlagSet
	p := ops.fSet.String(name, defVal, opt.Desc)
	opt.flag = ops.fSet.Lookup(name)

	return p
}

// StrVar binding an string option flag
func (ops *CliOpts) StrVar(p *string, opt *CliOpt) { ops.strOpt(p, opt) }

// StrOpt binding a string option.
//
// If defValAndDesc only one elem, will as desc message.
func (ops *CliOpts) StrOpt(p *string, name, shorts string, defValAndDesc ...string) {
	var defVal, desc string
	if ln := len(defValAndDesc); ln > 0 {
		if ln >= 2 {
			defVal, desc = defValAndDesc[0], defValAndDesc[1]
		} else { // only one as desc
			desc = defValAndDesc[0]
		}
	}

	ops.StrOpt2(p, name, desc, func(opt *CliOpt) {
		opt.DefVal = defVal
		opt.Shorts = strutil.Split(shorts, shortSepChar)
	})
}

// StrOpt2 binding a string option, and allow with CliOptFn for config option.
func (ops *CliOpts) StrOpt2(p *string, nameAndShorts, desc string, setFns ...CliOptFn) {
	ops.strOpt(p, NewOpt(nameAndShorts, desc, "", setFns...))
}

// binding option and shorts
func (ops *CliOpts) strOpt(p *string, opt *CliOpt) {
	defVal := opt.DValue().String()
	name := ops.checkFlagInfo(opt)

	// use *p as default value
	if defVal == "" && *p != "" {
		defVal = *p
	}

	// binding option to flag.FlagSet
	ops.fSet.StringVar(p, opt.Name, defVal, opt.Desc)
	opt.flag = ops.fSet.Lookup(name)
}

// --- intX option

// Int binding an int option flag, return pointer
func (ops *CliOpts) Int(name, shorts string, defVal int, desc string) *int {
	opt := newOpt(name, desc, defVal, shorts)
	name = ops.checkFlagInfo(opt)

	// binding option to flag.FlagSet
	ptr := ops.fSet.Int(name, defVal, opt.Desc)
	opt.flag = ops.fSet.Lookup(name)

	return ptr
}

// IntVar binding an int option flag
func (ops *CliOpts) IntVar(p *int, opt *CliOpt) { ops.intOpt(p, opt) }

// IntOpt binding an int option
func (ops *CliOpts) IntOpt(p *int, name, shorts string, defVal int, desc string) {
	ops.intOpt(p, newOpt(name, desc, defVal, shorts))
}

// IntOpt2 binding an int option and with config func.
func (ops *CliOpts) IntOpt2(p *int, nameAndShorts, desc string, setFns ...CliOptFn) {
	opt := newOpt(nameAndShorts, desc, 0, "")

	ops.intOpt(p, opt.WithOptFns(setFns...))
}

func (ops *CliOpts) intOpt(ptr *int, opt *CliOpt) {
	defVal := opt.DValue().Int()
	name := ops.checkFlagInfo(opt)

	// use *p as default value
	if defVal == 0 && *ptr != 0 {
		defVal = *ptr
	}

	// binding option to flag.FlagSet
	ops.fSet.IntVar(ptr, name, defVal, opt.Desc)
	opt.flag = ops.fSet.Lookup(name)
}

// Int64 binding an int64 option flag, return pointer
func (ops *CliOpts) Int64(name, shorts string, defVal int64, desc string) *int64 {
	opt := newOpt(name, desc, defVal, shorts)
	name = ops.checkFlagInfo(opt)

	// binding option to flag.FlagSet
	p := ops.fSet.Int64(name, defVal, opt.Desc)
	opt.flag = ops.fSet.Lookup(name)
	return p
}

// Int64Var binding an int64 option flag
func (ops *CliOpts) Int64Var(ptr *int64, opt *CliOpt) { ops.int64Opt(ptr, opt) }

// Int64Opt binding an int64 option
func (ops *CliOpts) Int64Opt(ptr *int64, name, shorts string, defValue int64, desc string) {
	ops.int64Opt(ptr, newOpt(name, desc, defValue, shorts))
}

func (ops *CliOpts) int64Opt(ptr *int64, opt *CliOpt) {
	defVal := opt.DValue().Int64()
	name := ops.checkFlagInfo(opt)

	// use *p as default value
	if defVal == 0 && *ptr != 0 {
		defVal = *ptr
	}

	// binding option to flag.FlagSet
	ops.fSet.Int64Var(ptr, name, defVal, opt.Desc)
	opt.flag = ops.fSet.Lookup(name)
}

// --- uintX option

// Uint binding an int option flag, return pointer
func (ops *CliOpts) Uint(name, shorts string, defVal uint, desc string) *uint {
	opt := newOpt(name, desc, defVal, shorts)
	name = ops.checkFlagInfo(opt)

	// binding option to flag.FlagSet
	ptr := ops.fSet.Uint(name, defVal, opt.Desc)
	opt.flag = ops.fSet.Lookup(name)

	return ptr
}

// UintVar binding an uint option flag
func (ops *CliOpts) UintVar(ptr *uint, opt *CliOpt) { ops.uintOpt(ptr, opt) }

// UintOpt binding an uint option
func (ops *CliOpts) UintOpt(ptr *uint, name, shorts string, defValue uint, desc string) {
	ops.uintOpt(ptr, newOpt(name, desc, defValue, shorts))
}

func (ops *CliOpts) uintOpt(ptr *uint, opt *CliOpt) {
	defVal := opt.DValue().Int()
	name := ops.checkFlagInfo(opt)

	// binding option to flag.FlagSet
	ops.fSet.UintVar(ptr, name, uint(defVal), opt.Desc)
	opt.flag = ops.fSet.Lookup(name)
}

// Uint64 binding an int option flag, return pointer
func (ops *CliOpts) Uint64(name, shorts string, defVal uint64, desc string) *uint64 {
	opt := newOpt(name, desc, defVal, shorts)
	name = ops.checkFlagInfo(opt)

	ptr := ops.fSet.Uint64(name, defVal, opt.Desc)
	opt.flag = ops.fSet.Lookup(name)

	return ptr
}

// Uint64Var binding an uint option flag
func (ops *CliOpts) Uint64Var(ptr *uint64, opt *CliOpt) { ops.uint64Opt(ptr, opt) }

// Uint64Opt binding an uint64 option
func (ops *CliOpts) Uint64Opt(ptr *uint64, name, shorts string, defVal uint64, desc string) {
	ops.uint64Opt(ptr, newOpt(name, desc, defVal, shorts))
}

// binding option and shorts
func (ops *CliOpts) uint64Opt(ptr *uint64, opt *CliOpt) {
	defVal := opt.DValue().Int64()
	name := ops.checkFlagInfo(opt)

	// binding option to flag.FlagSet
	ops.fSet.Uint64Var(ptr, name, uint64(defVal), opt.Desc)
	opt.flag = ops.fSet.Lookup(name)
}

// Var binding an custom var option flag
func (ops *CliOpts) Var(ptr flag.Value, opt *CliOpt) { ops.varOpt(ptr, opt) }

// VarOpt binding a custom var option
//
// Usage:
//
//	var names gcli.Strings
//	cmd.VarOpt(&names, "tables", "t", "description ...")
func (ops *CliOpts) VarOpt(v flag.Value, name, shorts, desc string) {
	ops.varOpt(v, newOpt(name, desc, nil, shorts))
}

// VarOpt2 binding an int option and with config func.
func (ops *CliOpts) VarOpt2(v flag.Value, nameAndShorts, desc string, setFns ...CliOptFn) {
	ops.varOpt(v, NewOpt(nameAndShorts, desc, nil, setFns...))
}

// binding option and shorts
func (ops *CliOpts) varOpt(v flag.Value, opt *CliOpt) {
	name := ops.checkFlagInfo(opt)

	// binding option to flag.FlagSet
	ops.fSet.Var(v, name, opt.Desc)
	opt.flag = ops.fSet.Lookup(name)
}

// check flag option name and short-names
func (ops *CliOpts) checkFlagInfo(opt *CliOpt) string {
	// check flag name
	name := opt.initCheck()
	if _, ok := ops.opts[name]; ok {
		helper.Panicf("redefined option flag '%s'", name)
	}

	// NOTICE: must init some required fields
	if ops.names == nil {
		ops.names = map[string]int{}
		ops.opts = map[string]*CliOpt{}
		ops.InitFlagSet("flags-" + opt.Name)
	}

	// is a short name
	helpLen := opt.helpNameLen()
	// fix: must exclude Hidden option
	if !opt.Hidden {
		// +7: type placeholder width
		ops.optMaxLen = mathutil.MaxInt(ops.optMaxLen, helpLen+6)
	}

	// check short names
	ops.checkShortNames(name, opt.Shorts)

	// update name length
	ops.names[name] = helpLen
	// storage opt and name
	ops.opts[name] = opt
	return name
}

// check short names
func (ops *CliOpts) checkShortNames(name string, shorts []string) {
	if len(shorts) == 0 {
		return
	}

	ops.hasShort = true
	if ops.shorts == nil {
		ops.shorts = map[string]string{}
	}

	for _, short := range shorts {
		if name == short {
			helper.Panicf("short name '%s' has been used as the current option name", short)
		}

		if _, ok := ops.names[short]; ok {
			helper.Panicf("short name '%s' has been used as an option name", short)
		}

		if n, ok := ops.shorts[short]; ok {
			helper.Panicf("short name '%s' has been used by option '%s'", short, n)
		}

		// storage short name
		ops.shorts[short] = name
	}
}

/***********************************************************************
 * Options:
 * - helper methods
 ***********************************************************************/

// IterAll Iteration all flag options with metadata
func (ops *CliOpts) IterAll(fn func(f *flag.Flag, opt *CliOpt)) {
	ops.fSet.VisitAll(func(f *flag.Flag) {
		if _, ok := ops.opts[f.Name]; ok {
			fn(f, ops.opts[f.Name])
		}
	})
}

// ShortNames get all short-names of the option
func (ops *CliOpts) ShortNames(name string) (ss []string) {
	if opt, ok := ops.opts[name]; ok {
		ss = opt.Shorts
	}
	return
}

// IsShortOpt alias of the IsShortcut()
func (ops *CliOpts) IsShortOpt(short string) bool { return ops.IsShortName(short) }

// IsShortName check it is a shortcut name
func (ops *CliOpts) IsShortName(short string) bool {
	if len(short) != 1 {
		return false
	}

	_, ok := ops.shorts[short]
	return ok
}

// IsOption check it is an option name
func (ops *CliOpts) IsOption(name string) bool { return ops.HasOption(name) }

// HasOption check it is an option name
func (ops *CliOpts) HasOption(name string) bool {
	_, ok := ops.names[name]
	return ok
}

// LookupFlag get flag.Flag by name
func (ops *CliOpts) LookupFlag(name string) *flag.Flag { return ops.fSet.Lookup(name) }

// Opt get CliOpt by name
func (ops *CliOpts) Opt(name string) *CliOpt { return ops.opts[name] }

// Opts get all flag options
func (ops *CliOpts) Opts() map[string]*CliOpt { return ops.opts }

/***********************************************************************
 * flag options metadata
 ***********************************************************************/

// CliOptFn opt config func type
type CliOptFn func(opt *CliOpt)

// WithRequired setting for option
func WithRequired() CliOptFn {
	return func(opt *CliOpt) { opt.Required = true }
}

// WithDefault value setting for option
func WithDefault(defVal any) CliOptFn {
	return func(opt *CliOpt) { opt.DefVal = defVal }
}

// WithShorts setting for option
func WithShorts(shorts ...string) CliOptFn {
	return func(opt *CliOpt) { opt.Shorts = shorts }
}

// WithShortcut setting for option
func WithShortcut(shortcut string) CliOptFn {
	return func(opt *CliOpt) { opt.Shorts = strutil.Split(shortcut, shortSepChar) }
}

// WithValidator setting for option
func WithValidator(fn func(val string) error) CliOptFn {
	return func(opt *CliOpt) { opt.Validator = fn }
}

// CliOpt define for a flag option
type CliOpt struct {
	// go flag value
	flag *flag.Flag
	// Name of flag and description
	Name, Desc string
	// default value for the flag option
	DefVal any
	// wrapped the default value
	defVal *structs.Value
	// short names. eg: ["o", "a"]
	Shorts []string
	// EnvVar allow set flag value from ENV var
	EnvVar string

	// --- advanced settings

	// Hidden the option on help
	Hidden bool
	// Required the option is required
	Required bool
	// Validator support validate the option flag value
	Validator func(val string) error
	// TODO interactive question for collect value
	Question string
}

// NewOpt quick create an CliOpt instance
func NewOpt(nameAndShorts, desc string, defVal any, setFns ...CliOptFn) *CliOpt {
	return newOpt(nameAndShorts, desc, defVal, "").WithOptFns(setFns...)
}

// newOpt quick create an CliOpt instance
func newOpt(nameAndShorts, desc string, defVal any, shortcut string) *CliOpt {
	return &CliOpt{
		Name: nameAndShorts,
		Desc: desc,
		// other info
		DefVal: defVal,
		Shorts: strutil.Split(shortcut, shortSepChar),
	}
}

// WithOptFns set for current option
func (m *CliOpt) WithOptFns(fns ...CliOptFn) *CliOpt {
	for _, fn := range fns {
		fn(m)
	}
	return m
}

func (m *CliOpt) initCheck() string {
	// feat: support add shorts by option name. eg: "name,n"
	if strings.ContainsRune(m.Name, shortSepRune) {
		ss := strings.Split(m.Name, shortSepChar)
		m.Name = ss[0]
		m.Shorts = append(m.Shorts, ss[1:]...)
	}

	if m.Desc != "" {
		desc := strings.Trim(m.Desc, "; ")
		if strings.ContainsRune(desc, ';') {
			// format: desc;required
			// format: desc;required;env TODO parse ENV var
			parts := strutil.SplitNTrimmed(desc, ";", 2)
			if ln := len(parts); ln > 1 {
				bl, err := strutil.Bool(parts[1])
				if err == nil && bl {
					desc = parts[0]
					m.Required = true
				}
			}
		}

		m.Desc = desc
	}

	// filter shorts
	if len(m.Shorts) > 0 {
		m.Shorts = cflag.FilterNames(m.Shorts)
	}
	return m.goodName()
}

// good name of the flag
func (m *CliOpt) goodName() string {
	name := strings.Trim(m.Name, "- ")
	if name == "" {
		helper.Panicf("option flag name cannot be empty")
	}

	if !helper.IsGoodName(name) {
		helper.Panicf("option flag name '%s' is invalid, must match: %s", name, helper.RegGoodName)
	}

	// update self name
	m.Name = name
	return name
}

// Shorts2String join shorts to a string
func (m *CliOpt) Shorts2String(sep ...string) string { return m.ShortsString(sep...) }

// ShortsString join shorts to a string
func (m *CliOpt) ShortsString(sep ...string) string {
	if len(m.Shorts) == 0 {
		return ""
	}
	return strings.Join(m.Shorts, sepStr(sep))
}

// HelpName for show help
func (m *CliOpt) HelpName() string {
	return cflag.AddPrefixes(m.Name, m.Shorts)
}

func (m *CliOpt) helpNameLen() int {
	return len(m.HelpName())
}

// Validate the binding value
func (m *CliOpt) Validate(val string) error {
	if m.Required && val == "" {
		return fmt.Errorf("flag '%s' is required", m.Name)
	}

	// call user custom validator
	if m.Validator != nil {
		return m.Validator(val)
	}
	return nil
}

// Flag value
func (m *CliOpt) Flag() *flag.Flag {
	return m.flag
}

// DValue wrap the default value
func (m *CliOpt) DValue() *stdutil.Value {
	if m.defVal == nil {
		m.defVal = &stdutil.Value{V: m.DefVal}
	}
	return m.defVal
}
