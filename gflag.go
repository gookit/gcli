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
	// all option names of the command. {name: length}
	names map[string]int
	// metadata for all options
	metas map[string]*Meta
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
		fs: flag.NewFlagSet(name, flag.ContinueOnError),
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
func (gf *GFlags) BoolVar(p *bool, info Meta) {
	info.Name = gf.checkName(info.Name)
	defValue := info.DValue().Bool()

	// binding option and shortcuts
	gf.boolOpt(p, info.Name, defValue, info.Description(), info.Shortcuts)
}

// BoolOpt binding an bool option
func (gf *GFlags) BoolOpt(p *bool, name string, defValue bool, description string, shortcuts ...string) {
	name = gf.checkName(name)

	// binding option and shortcuts
	gf.boolOpt(p, name, defValue, description, shortcuts)
}

func (gf *GFlags) boolOpt(p *bool, name string, defValue bool, description string, shortcuts []string) {
	// binding option to flag.FlagSet
	gf.fs.BoolVar(p, name, defValue, description)

	// check and format
	fmtNames := gf.checkShortNames(name, shortcuts)
	if len(fmtNames) > 0 {
		for _, s := range fmtNames {
			gf.fs.BoolVar(p, s, defValue, "") // dont add description for short name
		}
	}
}

// --- float option

// Float64Var binding an float64 option flag
func (gf *GFlags) Float64Var(p *float64, info Meta) {
	info.Name = gf.checkName(info.Name)
	defValue := info.DValue().Float64()

	// binding option and shortcuts
	gf.float64Opt(p, info.Name, defValue, info.Description(), info.Shortcuts)
}

// Float64Opt binding an float64 option
func (gf *GFlags) Float64Opt(p *float64, name string, defValue float64, description string, shortcuts ...string) {
	name = gf.checkName(name)

	// binding option and shortcuts
	gf.float64Opt(p, name, defValue, description, shortcuts)
}

func (gf *GFlags) float64Opt(p *float64, name string, defValue float64, description string, shortNames []string) {
	// binding option to flag.FlagSet
	gf.fs.Float64Var(p, name, defValue, description)

	// check and format
	fmtNames := gf.checkShortNames(name, shortNames)
	if len(fmtNames) > 0 {
		for _, s := range fmtNames {
			gf.fs.Float64Var(p, s, defValue, "") // dont add description for short name
		}
	}
}

// --- string option

// StrVar binding an string option flag
func (gf *GFlags) StrVar(p *string, info Meta) {
	info.Name = gf.checkName(info.Name)
	defValue := info.DValue().String()

	// binding option and shortcuts
	gf.strOpt(p, info.Name, defValue, info.Description(), info.Shortcuts)
}

// StrOpt binding an string option
func (gf *GFlags) StrOpt(p *string, name, defValue, description string, shortcuts ...string) {
	name = gf.checkName(name)

	// binding option and shortcuts
	gf.strOpt(p, name, defValue, description, shortcuts)
}

func (gf *GFlags) strOpt(p *string, name, defValue, description string, shortNames []string) {
	// binding option to flag.FlagSet
	gf.fs.StringVar(p, name, defValue, description)

	// check and format
	fmtNames := gf.checkShortNames(name, shortNames)
	if len(fmtNames) > 0 {
		for _, s := range fmtNames {
			gf.fs.StringVar(p, s, defValue, "") // dont add description for short name
		}
	}
}

// --- intX option

// IntVar binding an int option flag
func (gf *GFlags) IntVar(p *int, info Meta) {
	info.Name = gf.checkName(info.Name)
	defValue := info.DValue().Int()

	// binding option and shortcuts
	gf.intOpt(p, info.Name, defValue, info.Description(), info.Shortcuts)
}

// IntOpt binding an int option
func (gf *GFlags) IntOpt(p *int, name string, defValue int, description string, shortcuts ...string) {
	name = gf.checkName(name)

	// binding option and shortcuts
	gf.intOpt(p, name, defValue, description, shortcuts)
}

func (gf *GFlags) intOpt(p *int, name string, defValue int, description string, shortNames []string) {
	// binding option to flag.FlagSet
	gf.fs.IntVar(p, name, defValue, description)

	// check and format
	fmtNames := gf.checkShortNames(name, shortNames)
	if len(fmtNames) > 0 {
		for _, s := range fmtNames {
			gf.fs.IntVar(p, s, defValue, "") // dont add description for short name
		}
	}
}

// Int64Var binding an uint option flag
func (gf *GFlags) Int64Var(p *int64, info Meta) {
	info.Name = gf.checkName(info.Name)
	defValue := info.DValue().Int64()

	// binding option and shortcuts
	gf.int64Opt(p, info.Name, defValue, info.Description(), info.Shortcuts)
}

// Int64Opt binding an int64 option
func (gf *GFlags) Int64Opt(p *int64, name string, defValue int64, description string, shortcuts ...string) {
	name = gf.checkName(name)

	// binding option and shortcuts
	gf.int64Opt(p, name, defValue, description, shortcuts)
}

func (gf *GFlags) int64Opt(p *int64, name string, defValue int64, description string, shortNames []string) {
	// binding option to flag.FlagSet
	gf.fs.Int64Var(p, name, defValue, description)

	// check and format
	fmtNames := gf.checkShortNames(name, shortNames)
	if len(fmtNames) > 0 {
		for _, s := range fmtNames {
			gf.fs.Int64Var(p, s, defValue, "") // dont add description for short name
		}
	}
}

// --- uintX option

// UintVar binding an uint option flag
func (gf *GFlags) UintVar(p *uint, info Meta) {
	info.Name = gf.checkName(info.Name)
	defValue := info.DValue().Int()

	// binding option and shortcuts
	gf.uintOpt(p, info.Name, uint(defValue), info.Description(), info.Shortcuts)
}

// UintOpt binding an uint option
func (gf *GFlags) UintOpt(p *uint, name string, defValue uint, description string, shortcuts ...string) {
	name = gf.checkName(name)

	// binding option and shortcuts
	gf.uintOpt(p, name, defValue, description, shortcuts)
}

func (gf *GFlags) uintOpt(p *uint, name string, defValue uint, description string, shortNames []string) {
	// binding option to flag.FlagSet
	gf.fs.UintVar(p, name, defValue, description)

	// check and format
	fmtNames := gf.checkShortNames(name, shortNames)
	if len(fmtNames) > 0 {
		for _, s := range fmtNames {
			gf.fs.UintVar(p, s, defValue, "") // dont add description for short name
		}
	}
}

// Uint64Var binding an uint option flag
func (gf *GFlags) Uint64Var(p *uint64, info Meta) {
	info.Name = gf.checkName(info.Name)
	defValue := info.DValue().Int64()

	gf.metas[info.Name] = &info

	// binding option and shortcuts
	gf.uint64Opt(p, info.Name, uint64(defValue), info.Description(), info.Shortcuts)
}

// Uint64Opt binding an uint64 option
func (gf *GFlags) Uint64Opt(p *uint64, name string, defValue uint64, description string, shortcuts ...string) {
	name = gf.checkName(name)

	// binding option and shortcuts
	gf.uint64Opt(p, name, defValue, description, shortcuts)
}

func (gf *GFlags) uint64Opt(p *uint64, name string, defValue uint64, description string, shortNames []string) {
	// binding option to flag.FlagSet
	gf.fs.Uint64Var(p, name, defValue, description)

	// check and format
	fmtNames := gf.checkShortNames(name, shortNames)
	if len(fmtNames) > 0 {
		for _, s := range fmtNames {
			gf.fs.Uint64Var(p, s, defValue, "") // dont add description for short name
		}
	}
}

// check option name and return clean name
func (gf *GFlags) checkName(name string) string {
	// init gf.names, gf.metas
	if gf.names == nil {
		gf.names = map[string]int{}
		gf.metas = map[string]*Meta{}
	}

	name = strings.Trim(name, "- ")
	if name == "" {
		panicf("option flag name cannot be empty")
	}

	if _, ok := gf.names[name]; ok {
		panicf("redefined option flag: %s", name)
	}

	nameLength := len(name)
	// is an short name
	if nameLength == 1 {
		nameLength += 1 // prefix: "-"
	} else {
		nameLength += 2 // prefix: "--"
	}

	if gf.flagMaxLen < nameLength {
		gf.flagMaxLen = nameLength
	}

	// storage name
	gf.names[name] = nameLength
	return name
}

// check short names
func (gf *GFlags) checkShortNames(name string, shorts []string) []string {
	if len(shorts) == 0 {
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
	// eg: "-o"
	// eg: "-o, -a"
	shortsLen := 4 * len(fmtShorts)
	nameLength := gf.names[name] + shortsLen

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

func (gf *GFlags) formatOneFlag(f *flag.Flag)  {
	meta := gf.metas[f.Name]

	// skip render:
	// - it is hidden flag option
	// - flag desc is empty
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
		if gf.IsShortcut(name) {
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
func (gf *GFlags) IterAll(fn func(f *flag.Flag, meta *Meta)) {
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
	return gf.IsShortcut(short)
}

// IsShortcut check it is a shortcut name
func (gf *GFlags) IsShortcut(short string) bool {
	if len(short) != 1 {
		return false
	}

	_, ok := gf.shortcuts[short]
	return ok
}

// Exists check it is a option name
func (gf *GFlags) Exists(name string) bool {
	_, ok := gf.names[name]
	return ok
}

// Hidden there are given option names
func (gf *GFlags) Hidden(names ...string) {
	for _, name := range names {
		if !gf.Exists(name) { // not registered
			continue
		}

		gf.metas[name].Hidden = true
	}
}

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

// Meta for an flag(option/argument)
type Meta struct {
	// varPtr interface{}
	// defVal *Value
	// name and description
	Name, UseFor string
	Hidden, Required bool
	// short names
	Shortcuts []string
	// default value for the option
	DefValue interface{}
}

// DValue wrap the default value
func (m *Meta) DValue() *goutil.Value {
	return &goutil.Value{V: m.DefValue}
}

// Description of the flag
func (m *Meta) Description() string {
	if len(m.UseFor) > 0 {
		return m.UseFor
	}

	return "no description"
}
