package gcli

import "flag"

// GFlags definition
type GFlags struct {
	fs *flag.FlagSet
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
func (gf *GFlags) StrOpt() {

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

// Info for an flag(option/argument)
type Info struct {
	Name, UseFor string
	// short names
	Shorts []string
	// default value for the option
	DefValue interface{}
	Required bool
}



