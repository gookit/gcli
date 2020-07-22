package gflag

import "flag"

// GFlags definition
type GFlags struct {
	fs *flag.FlagSet
}

// StrVar definition
func (gf *GFlags) StrVar() {

}

// GFlag definition
type GFlag struct {
	Name, UseFor string
	// short names
	Shorts []string
}



