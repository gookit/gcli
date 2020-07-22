package gcli

import (
	"flag"
	"io"
	"io/ioutil"
	"strings"
)

// GFlags definition
type GFlags struct {
	fs *flag.FlagSet
	// all option names of the command. {name: length}
	names map[string]int
	// shortcuts for command options. {short:name}
	// eg. {"n": "name", "o": "opt"}
	shortcuts map[string]string
}

// NewGFlags create an GFlags
func NewGFlags(name string) *GFlags {
	gf := &GFlags{
		fs: flag.NewFlagSet(name, flag.ContinueOnError),
	}

	// disable output internal error message on parse flags
	gf.fs.SetOutput(ioutil.Discard)
	gf.fs.Usage = func() {
		// nothing to do ... render usage on after
	}

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

// RenderUsage for all options
func (gf *GFlags) RenderUsage(w io.Writer) {
}

// StrOpt binding an bool option flag
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

// check option name
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

	// storage name
	gf.names[name] = len(name)
	return name
}

// check short names
func (gf *GFlags) checkShortNames(name string, shorts []string) []string {
	var fmtShorts []string
	for _, short := range shorts {
		short = strings.Trim(short, "- ")
		if short == "" {
			continue
		}

		// ensure it is one char
		char := short[0]
		short = string(char)
		if char < 'a' || char > 'Z'{
			panicf("shortcut name only allow: a-zA-Z(given: '%s')", short)
		}

		if n, ok := gf.shortcuts[short]; ok {
			panicf("shortcut name '%s' has been used by option '%s'", short, n)
		}

		fmtShorts = append(fmtShorts, short)
		// storage short name
		gf.shortcuts[short] = name
	}

	return fmtShorts
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

// Meta for an flag(option/argument)
type Meta struct {
	varPtr interface{}
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



