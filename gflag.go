package gcli

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/strutil"
)

// The options alignment type
// - Align right, padding left
// - Align left, padding right
const (
	AlignLeft  = strutil.PosRight
	AlignRight = strutil.PosLeft
)

// GFlagOption for render help information
type GFlagOption struct {
	// WithoutType dont display flag data type on print help
	WithoutType bool
	// NameDescOL flag and desc at one line on print help
	NameDescOL bool
	// Alignment flag align left or right. default is: right
	Alignment uint8
}

// GFlags definition
type GFlags struct {
	// GFlagOption option for render help message
	GFlagOption
	// raw flag set
	fs *flag.FlagSet
	// buf for build help message
	buf *bytes.Buffer
	// output for print help message
	out io.Writer
	// all option names of the command. {name: length} // TODO delete, move len to meta.
	names map[string]int
	// metadata for all options
	metas map[string]*FlagMeta
	// shortcuts for command options. {short:name}
	// eg. {"n": "name", "o": "opt"}
	shortcuts map[string]string
	// mapping for name to shortcut {"name": {"n", "m"}}
	name2shorts map[string][]string
	// flag name max length
	// eg: "-V, --version" length is 13
	flagMaxLen int
}

// NewGFlags create an GFlags
func NewGFlags(name string) *GFlags {
	gf := &GFlags{
		out: os.Stdout,
		fs:  flag.NewFlagSet(name, flag.ContinueOnError),
	}

	// disable output internal error message on parse flags
	gf.fs.SetOutput(ioutil.Discard)
	// nothing to do ... render usage on after parsed
	gf.fs.Usage = func() {}

	return gf
}

// Parse given arguments
func (gf *GFlags) Parse(args []string) error {
	return gf.fs.Parse(args)
}

// FromStruct binding options
func (gf *GFlags) FromStruct(ptr interface{}) error {
	// TODO WIP
	return nil
}

// WithOption for render help panel message
func (gf *GFlags) WithOption(cfg GFlagOption) *GFlags {
	gf.GFlagOption = cfg
	return gf
}

/***********************************************************************
 * GFlag:
 * - binding option var
 ***********************************************************************/

// --- bool option

// BoolVar binding an bool option flag
func (gf *GFlags) BoolVar(p *bool, meta FlagMeta) {
	gf.boolOpt(p, &meta)
}

// BoolOpt binding an bool option
func (gf *GFlags) BoolOpt(p *bool, name string, defValue bool, desc string, shortcuts ...string) {
	gf.boolOpt(p, &FlagMeta{
		Name:   name,
		Desc:   desc,
		DefVal: defValue,
		Shorts: shortcuts,
	})
}

// binding option and shortcuts
func (gf *GFlags) boolOpt(p *bool, meta *FlagMeta) {
	defValue := meta.DValue().Bool()
	fmtName := gf.checkName(meta.Name, meta)

	// binding option to flag.FlagSet
	gf.fs.BoolVar(p, fmtName, defValue, meta.Desc)

	// binding all short name options to flag.FlagSet
	for _, s := range meta.Shorts {
		gf.fs.BoolVar(p, s, defValue, "") // dont add description for short name
	}
}

// --- float option

// Float64Var binding an float64 option flag
func (gf *GFlags) Float64Var(p *float64, meta FlagMeta) {
	gf.float64Opt(p, &meta)
}

// Float64Opt binding an float64 option
func (gf *GFlags) Float64Opt(p *float64, name string, defValue float64, desc string, shortcuts ...string) {
	gf.float64Opt(p, &FlagMeta{
		Name:   name,
		Desc:   desc,
		DefVal: defValue,
		Shorts: shortcuts,
	})
}

func (gf *GFlags) float64Opt(p *float64, meta *FlagMeta) {
	defValue := meta.DValue().Float64()
	fmtName := gf.checkName(meta.Name, meta)

	// binding option to flag.FlagSet
	gf.fs.Float64Var(p, fmtName, defValue, meta.Desc)

	// binding all short name options to flag.FlagSet
	for _, s := range meta.Shorts {
		gf.fs.Float64Var(p, s, defValue, "") // dont add description for short name
	}
}

// --- string option

// StrVar binding an string option flag
func (gf *GFlags) StrVar(p *string, meta FlagMeta) {
	gf.strOpt(p, &meta)
}

// StrOpt binding an string option
func (gf *GFlags) StrOpt(p *string, name, defValue, desc string, shortcuts ...string) {
	gf.strOpt(p, &FlagMeta{
		Name:   name,
		Desc:   desc,
		DefVal: defValue,
		Shorts: shortcuts,
	})
}

// binding option and shortcuts
func (gf *GFlags) strOpt(p *string, meta *FlagMeta) {
	defValue := meta.DValue().String()
	fmtName := gf.checkName(meta.Name, meta)

	// binding option to flag.FlagSet
	gf.fs.StringVar(p, fmtName, defValue, meta.Desc)

	// binding all short name options to flag.FlagSet
	for _, s := range meta.Shorts {
		gf.fs.StringVar(p, s, defValue, "") // dont add description for short name
	}
}

// --- intX option

// IntVar binding an int option flag
func (gf *GFlags) IntVar(p *int, meta FlagMeta) {
	gf.intOpt(p, &meta)
}

// IntOpt binding an int option
func (gf *GFlags) IntOpt(p *int, name string, defValue int, desc string, shortcuts ...string) {
	gf.intOpt(p, &FlagMeta{
		Name:   name,
		Desc:   desc,
		DefVal: defValue,
		Shorts: shortcuts,
	})
}

func (gf *GFlags) intOpt(p *int, meta *FlagMeta) {
	defValue := meta.DValue().Int()
	fmtName := gf.checkName(meta.Name, meta)

	// binding option to flag.FlagSet
	gf.fs.IntVar(p, fmtName, defValue, meta.Desc)

	// binding all short name options to flag.FlagSet
	for _, s := range meta.Shorts {
		gf.fs.IntVar(p, s, defValue, "") // dont add description for short name
	}
}

// Int64Var binding an uint option flag
func (gf *GFlags) Int64Var(p *int64, meta FlagMeta) {
	gf.int64Opt(p, &meta)
}

// Int64Opt binding an int64 option
func (gf *GFlags) Int64Opt(p *int64, name string, defValue int64, desc string, shortcuts ...string) {
	gf.int64Opt(p, &FlagMeta{
		Name:   name,
		Desc:   desc,
		DefVal: defValue,
		Shorts: shortcuts,
	})
}

func (gf *GFlags) int64Opt(p *int64, meta *FlagMeta) {
	defValue := meta.DValue().Int64()
	fmtName := gf.checkName(meta.Name, meta)

	// binding option to flag.FlagSet
	gf.fs.Int64Var(p, fmtName, defValue, meta.Desc)

	// binding all short options to flag.FlagSet
	for _, s := range meta.Shorts {
		gf.fs.Int64Var(p, s, defValue, "") // dont add description for short name
	}
}

// --- uintX option

// UintVar binding an uint option flag
func (gf *GFlags) UintVar(p *uint, meta FlagMeta) {
	gf.uintOpt(p, &meta)
}

// UintOpt binding an uint option
func (gf *GFlags) UintOpt(p *uint, name string, defValue uint, desc string, shortcuts ...string) {
	gf.uintOpt(p, &FlagMeta{
		Name:   name,
		Desc:   desc,
		DefVal: defValue,
		Shorts: shortcuts,
	})
}

func (gf *GFlags) uintOpt(p *uint, meta *FlagMeta) {
	defValue := meta.DValue().Int()
	fmtName := gf.checkName(meta.Name, meta)

	// binding option to flag.FlagSet
	gf.fs.UintVar(p, fmtName, uint(defValue), meta.Desc)

	// binding all short options to flag.FlagSet
	for _, s := range meta.Shorts {
		gf.fs.UintVar(p, s, uint(defValue), "") // dont add description for short name
	}
}

// Uint64Var binding an uint option flag
func (gf *GFlags) Uint64Var(p *uint64, meta FlagMeta) {
	// binding option and shortcuts
	gf.uint64Opt(p, &meta)
}

// Uint64Opt binding an uint64 option
func (gf *GFlags) Uint64Opt(p *uint64, name string, defValue uint64, desc string, shortcuts ...string) {
	// binding option and shortcuts
	gf.uint64Opt(p, &FlagMeta{
		Name:   name,
		Desc:   desc,
		DefVal: defValue,
		Shorts: shortcuts,
	})
}

func (gf *GFlags) uint64Opt(p *uint64, meta *FlagMeta) {
	defValue := meta.DValue().Int64()
	fmtName := gf.checkName(meta.Name, meta)

	// binding option to flag.FlagSet
	gf.fs.Uint64Var(p, fmtName, uint64(defValue), meta.Desc)

	// binding all short options to flag.FlagSet
	for _, s := range meta.Shorts {
		gf.fs.Uint64Var(p, s, uint64(defValue), "") // dont add description for short name
	}
}

// check option name and return clean name
func (gf *GFlags) checkName(name string, meta *FlagMeta) string {
	// init gf.names, gf.metas
	if gf.names == nil {
		gf.names = map[string]int{}
		gf.metas = map[string]*FlagMeta{}
	}

	// check name
	name = strings.Trim(name, "- ")
	if name == "" {
		panicf("option flag name cannot be empty")
	}

	if _, ok := gf.metas[name]; ok {
		panicf("redefined option flag: %s", name)
	}

	nameLength := len(name)
	// is an short name
	if nameLength == 1 {
		nameLength += 1 // prefix: "-"
	} else {
		nameLength += 2 // prefix: "--"
	}

	// fix: must exclude Hidden option
	if !meta.Hidden && gf.flagMaxLen < nameLength {
		gf.flagMaxLen = nameLength
	}

	// update name
	meta.Name = name
	// check and format short names
	meta.Shorts = gf.checkShortNames(name, nameLength, meta.Shorts)

	// storage meta and name
	gf.metas[name] = meta
	return name
}

// check short names
func (gf *GFlags) checkShortNames(name string, nameLength int, shorts []string) []string {
	if len(shorts) == 0 {
		// record name without shorts
		gf.names[name] = nameLength
		return shorts
	}

	// init gf.shortcuts and gf.name2shorts
	if gf.shortcuts == nil {
		gf.shortcuts = map[string]string{}
		gf.name2shorts = map[string][]string{}
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
			panicf("shortcut name only allow: A-Za-z(given: '%s')", short)
		}

		if n, ok := gf.shortcuts[short]; ok {
			panicf("shortcut name '%s' has been used by option '%s'", short, n)
		}

		fmtShorts = append(fmtShorts, short)
		// storage short name
		gf.shortcuts[short] = name
	}

	// one short = '-' + 'x' + ',' + ' '
	// eg: "-o, " len=4
	// eg: "-o, -a, " len=8
	nameLength += 4 * len(fmtShorts)

	// update name length
	gf.names[name] = nameLength
	if gf.flagMaxLen < nameLength {
		gf.flagMaxLen = nameLength
	}

	gf.name2shorts[name] = fmtShorts
	return fmtShorts
}

/***********************************************************************
 * GFlag:
 * - render help message
 ***********************************************************************/

// PrintHelpPanel for all options to the gf.out
func (gf *GFlags) PrintHelpPanel() {
	dump.P(gf.flagMaxLen)
	color.Fprint(gf.out, gf.String())
}

// String for all flag options
func (gf *GFlags) String() string {
	if gf.buf == nil {
		gf.buf = new(bytes.Buffer)
	}

	// repeat call
	if gf.buf.Len() < 1 {
		// refer gf.Fs().PrintDefaults()
		gf.Fs().VisitAll(gf.formatOneFlag)
	}

	return gf.buf.String()
}

func (gf *GFlags) formatOneFlag(f *flag.Flag) {
	// Skip render:
	// - meta is not exists
	// - it is hidden flag option
	// - flag desc is empty
	meta, ok := gf.metas[f.Name]
	if !ok {
		return
	}
	if meta.Hidden || f.Usage == "" {
		return
	}

	var s, fullName string
	name := f.Name
	// eg: "-V, --version" length is: 13
	fLen := gf.names[name]

	// - build flag name info
	// is long option
	if len(name) > 1 {
		// find shortcuts
		shortcuts := gf.ShortNames(name)
		if len(shortcuts) == 0 {
			fullName = "--" + name
			// s = fmt.Sprintf("      <info>--%s</>", name)
		} else {
			fullName = fmt.Sprintf("%s, --%s", shortcuts2str(shortcuts), name)
			// s = fmt.Sprintf("  <info>%s, --%s</>", shortcuts2str(shortcuts), name)
		}
	} else {
		// is short option name, skip it
		if gf.IsShortName(name) {
			return
		}

		// only short option
		// s = fmt.Sprintf("  <info>-%s</>", name)
		fullName = "-" + name
	}

	// padding space to same width.
	fullName = strutil.Padding(fullName, " ", gf.flagMaxLen, gf.Alignment)

	s = fmt.Sprintf("  <info>%s</>", fullName)

	// - build flag type info
	typeName, usage := flag.UnquoteUsage(f)
	// typeName: option value data type: int, string, ..., bool value will return ""
	if gf.WithoutType == false && len(typeName) > 0 {
		s += fmt.Sprintf(" <magenta>%s</>", typeName)
	}

	// - flag and description at one line
	// - Boolean flags of one ASCII letter are so common we
	// treat them specially, putting their usage on the same line.
	if gf.NameDescOL || fLen <= 4 { // space, space, '-', 'x'.
		s += "    "
	} else {
		// display description on new line
		s += "\n        "
	}

	// - build description
	s += strings.Replace(strutil.UpperFirst(usage), "\n", "\n        ", -1)

	if !isZeroValue(f, f.DefValue) {
		if _, ok := f.Value.(*stringValue); ok {
			// put quotes on the value
			s += fmt.Sprintf(" (default <magentaB>%q</>)", f.DefValue)
		} else {
			s += fmt.Sprintf(" (default <magentaB>%v</>)", f.DefValue)
		}
	}

	// save to buffer
	gf.buf.WriteString(s)
	gf.buf.WriteByte('\n')
}

/***********************************************************************
 * GFlags:
 * - helper methods
 ***********************************************************************/

// IterAll Iteration all flag options with metadata
func (gf *GFlags) IterAll(fn func(f *flag.Flag, meta *FlagMeta)) {
	gf.Fs().VisitAll(func(f *flag.Flag) {
		if _, ok := gf.metas[f.Name]; ok {
			fn(f, gf.metas[f.Name])
		}
	})
}

// ShortNames get all short-names of the option
func (gf *GFlags) ShortNames(name string) (ss []string) {
	if len(gf.name2shorts) == 0 {
		return
	}

	return gf.name2shorts[name]
}

// IsShortOpt alias of the IsShortcut()
func (gf *GFlags) IsShortOpt(short string) bool {
	return gf.IsShortName(short)
}

// IsShortcut check it is a shortcut name
func (gf *GFlags) IsShortName(short string) bool {
	if len(short) != 1 {
		return false
	}

	_, ok := gf.shortcuts[short]
	return ok
}

// HasOption check it is a option name
func (gf *GFlags) HasOption(name string) bool {
	_, ok := gf.names[name]
	return ok
}

// HasFlagMeta check it is has FlagMeta
func (gf *GFlags) HasFlagMeta(name string) bool {
	_, ok := gf.metas[name]
	return ok
}

// FlagMeta get FlagMeta by name
func (gf *GFlags) FlagMeta(name string) *FlagMeta {
	return gf.metas[name]
}

// Metas get all flag metas
func (gf *GFlags) Metas() map[string]*FlagMeta {
	return gf.metas
}

// Hidden there are given option names
// func (gf *GFlags) Hidden(names ...string) {
// 	for _, name := range names {
// 		if !gf.HasOption(name) { // not registered
// 			continue
// 		}
//
// 		gf.metas[name].Hidden = true
// 	}
// }

// Name of the Flags
func (gf *GFlags) Name() string {
	return gf.fs.Name()
}

// Fs get the raw *flag.FlagSet
func (gf *GFlags) Fs() *flag.FlagSet {
	return gf.fs
}

// Fs set the raw *flag.FlagSet
func (gf *GFlags) SetFs(fs *flag.FlagSet) {
	gf.fs = fs
}

// SetOutput for the GFlags
func (gf *GFlags) SetOutput(out io.Writer) {
	gf.out = out
}

/***********************************************************************
 * GFlag:
 * - flag metadata
 ***********************************************************************/

// FlagMeta for an flag(option/argument)
type FlagMeta struct {
	// varPtr interface{}
	// defVal *goutil.Value
	// name and description
	Name, Desc string
	// default value for the option
	DefVal interface{}
	// short names. eg: ["o", "a"]
	Shorts []string
	// special setting
	Hidden, Required bool
}

// DValue wrap the default value
func (m *FlagMeta) DValue() *goutil.Value {
	return &goutil.Value{V: m.DefVal}
}

// Description of the flag
func (m *FlagMeta) Description() string {
	if len(m.Desc) > 0 {
		return m.Desc
	}

	return "no description"
}
