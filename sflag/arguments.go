package sflag

import "github.com/gookit/gcli/v2"

// Arguments definition
type Arguments struct {
	// args definition for a command.
	// eg. {
	// 	{"arg0", "this is first argument", false, false},
	// 	{"arg1", "this is second argument", false, false},
	// }
	args []*gcli.Argument
	// record min length for args
	// argsMinLen int
	// record argument names and defined positional relationships
	// {
	// 	// name: position
	// 	"arg0": 0,
	// 	"arg1": 1,
	// }
	argsIndexes    map[string]int
	hasArrayArg    bool
	hasOptionalArg bool
}

// Add a new argument
func (ags *Arguments) Add(name, description string) {
	// todo ...
}
