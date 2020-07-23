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
	"github.com/gookit/goutil/strutil"
)

// GFlags definition
type GFlags struct {
	fs *flag.FlagSet
	// output for print help message
	out io.Writer
	// buf for build help message
	buf *bytes.Buffer
	// all option names of the command. {name: length}
	names map[string]int
	// shortcuts for command options. {short:name}
	// eg. {"n": "name", "o": "opt"}
	shortcuts map[string]string
	// mapping for name to shortcut {"name": {"n", "m"}}
	name2shorts map[string][]string

	// dont display flag data type
	flagNoType bool
	// flag and desc at one line
	flagDescOL bool
	flagMaxLen int
}

// NewGFlags create an GFlags
func NewGFlags(name string) *GFlags {
	gf := &GFlags{
		out: os.Stdout,
		buf: new(bytes.Buffer),
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

/***********************************************************************
 * GFlag:
 * - binding option var
 ***********************************************************************/

// --- bool opt

// BoolOpt binding an bool option flag
func (gf *GFlags) BoolOpt(p *bool, info *Meta) {
	info.Name = gf.checkName(info.Name)
	defValue := info.DValue().Bool()

	// binding option and shortcuts
	gf.boolOpt(p, info.Name, defValue, info.Description(), info.Shortcuts)
}

// BoolVar binding an bool option
func (gf *GFlags) BoolVar(p *bool, name string, defValue bool, description string, shortcuts ...string) {
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

// --- string opt

// StrOpt binding an string option flag
func (gf *GFlags) StrOpt(p *string, info *Meta) {
	info.Name = gf.checkName(info.Name)
	defValue := info.DValue().String()

	// binding option and shortcuts
	gf.strOpt(p, info.Name, defValue, info.Description(), info.Shortcuts)
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

// --- uintX opt

// UintOpt binding an uint option flag
func (gf *GFlags) UintOpt(p *uint, info *Meta) {
	info.Name = gf.checkName(info.Name)
	defValue := info.DValue().Int()

	// binding option and shortcuts
	gf.uintOpt(p, info.Name, uint(defValue), info.Description(), info.Shortcuts)
}

// UintVar binding an uint option
func (gf *GFlags) UintVar(p *uint, name string, defValue uint, description string, shortcuts ...string) {
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

// UintOpt binding an uint option flag
func (gf *GFlags) Uint64Opt(p *uint64, info *Meta) {
	info.Name = gf.checkName(info.Name)
	defValue := info.DValue().Int64()

	// binding option and shortcuts
	gf.uint64Opt(p, info.Name, uint64(defValue), info.Description(), info.Shortcuts)
}

// Uint64Var binding an uint64 option
func (gf *GFlags) Uint64Var(p *uint64, name string, defValue uint64, description string, shortcuts ...string) {
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
	// init gf.names
	if gf.names == nil {
		gf.names = map[string]int{}
	}

	name = strings.Trim(name, "- ")
	if name == "" {
		panicf("option flag name cannot be empty")
	}

	if _, ok := gf.names[name]; ok {
		panicf("redefined option flag: %s", name)
	}

	nameLength := len(name)
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
		if !isAlphabet(char) {
			panicf("shortcut name only allow: A-Za-z(given: '%s')", short)
		}

		if n, ok := gf.shortcuts[short]; ok {
			panicf("shortcut name '%s' has been used by option '%s'", short, n)
		}

		fmtShorts = append(fmtShorts, short)
		// storage short name
		gf.shortcuts[short] = name
	}

	gf.names[name] += len(fmtShorts)
	gf.name2shorts[name] = fmtShorts
	return fmtShorts
}

/***********************************************************************
 * GFlag:
 * - helper methods
 ***********************************************************************/

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

// RenderUsage for all options
func (gf *GFlags) RenderUsage(w io.Writer) {

}

// PrintHelp for all options
func (gf *GFlags) PrintHelpPanel() {
	buf := new(bytes.Buffer)

	gf.Fs().PrintDefaults()

	gf.Fs().VisitAll(gf.formatOneFlag)

	color.Fprint(gf.out, buf.String())
}

func (gf *GFlags) formatOneFlag(f *flag.Flag)  {
	// if desc is empty(hidden flag), skip it
	if f.Usage == "" {
		return
	}

	var s string
	name := f.Name
	fLen := len(f.Name)

	// - build flag name info
	// is long option
	if fLen > 1 {
		// find shortcuts
		shortcuts := gf.ShortNames(name)
		shortLen := len(shortcuts)
		if shortLen == 0 {
			s = fmt.Sprintf("      <info>--%s</>", name)
		} else {
			s = fmt.Sprintf("  <info>%s, --%s</>", shortcuts2str(shortcuts), name)
		}
	} else {
		// is short option name, skip it
		if gf.IsShortcut(name) {
			return
		}

		// only short option
		s = fmt.Sprintf("  <info>-%s</>", name)
	}

	// - build flag type info
	typeName, usage := flag.UnquoteUsage(f)
	// option value data type: int, string, ...
	if len(typeName) > 0 {
		s += fmt.Sprintf(" <magenta>%s</>", typeName)
	}

	// Boolean flags of one ASCII letter are so common we
	// treat them specially, putting their usage on the same line.
	if len(s) <= 4 { // space, space, '-', 'x'.
		s += "\t"
	} else {
		// Four spaces before the tab triggers good alignment
		// for both 4- and 8-space tab stops.
		s += "\n    \t"
	}

	// - build description
	s += strings.Replace(strutil.UpperFirst(usage), "\n", "\n    \t", -1)

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
 * GFlag:
 * - flag metadata
 ***********************************************************************/

// Meta for an flag(option/argument)
type Meta struct {
	varPtr interface{}
	// defVal *Value
	// name and description
	Name, UseFor string
	// short names
	Shortcuts []string
	// default value for the option
	DefValue interface{}
	Required bool
}

// DValue wrap the default value
func (m *Meta) DValue() *Value {
	return &Value{V: m.DefValue}
}

// Description of the flag
func (m *Meta) Description() string {
	if len(m.UseFor) > 0 {
		return m.UseFor
	}

	return "no description"
}
