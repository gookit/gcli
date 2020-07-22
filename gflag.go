package gcli

import (
	"flag"
	"strings"

	"github.com/gookit/goutil"
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
	return &GFlags{
		fs: flag.NewFlagSet(name, flag.ContinueOnError),
	}
}

// FromStruct binding options
func (gf *GFlags) FromStruct(ptr interface{}) error {
	// TODO WIP
	return nil
}

// StrVar definition
func (gf *GFlags) StrOpt(p *string, info *Meta) {
	info.Name = gf.checkName(info.Name)
	defValue := info.DValue().String()

	// binding option to flag.FlagSet
	gf.fs.StringVar(p, info.Name, defValue, info.Description())

	shortNames := gf.checkShortNames(info.Name, info.Shortcuts)
	if len(shortNames) > 0 {
		for _, s := range shortNames {
			gf.fs.StringVar(p, s, defValue, "")
		}
	}
}

// AddFlag definition
func (gf *GFlags) strOpt(name string) string {
}

// AddFlag definition
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

// addShortcuts definition
func (gf *GFlags) checkShortNames(name string, shorts []string) []string {
	for _, short := range shorts {
		// ensure it is one char
		short = string(short[0])
		if n, ok := gf.shortcuts[short]; ok {
			panicf("shortcut name '%s' has been used by option '%s'", short, n)
		}
	}


}

// AddFlag definition
func (gf *GFlags) AddFlag(f Flag) {

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

// Flag definition
type Flag struct {
	Name, UseFor string
	// short names
	Shorts []string
	// default value for the option
	DefValue interface{}
	Required bool
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



