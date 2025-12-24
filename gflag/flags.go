package gflag

import (
	"encoding"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/gookit/goutil/cflag"
)

// Flag type. alias of flag.Flag
type Flag = flag.Flag

// sortFlags returns the flags as a slice in lexicographical sorted order.
//
// from the go flag package.
func sortFlags(flags map[string]*Flag) []*Flag {
	result := make([]*Flag, len(flags))
	i := 0
	for _, f := range flags {
		result[i] = f
		i++
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

// FlagSet custom implement flag set, like flag.FlagSet.
// But support short names for options.
type FlagSet struct {
	name string
	args []string // remain arguments after parse

	parsed bool
	actual map[string]*Flag
	formal map[string]*Flag

	// error handling strategy
	errorHandling flag.ErrorHandling
	// output  io.Writer // nil means stderr; use Output() accessor

	Usage func()

	// short names map for options. format: {short: name}
	//
	// eg. {"n": "name", "o": "opt"}
	shorts map[string]string
}

// NewFlagSet create a new FlagSet
func NewFlagSet(name string, errorHandling flag.ErrorHandling) *FlagSet {
	return &FlagSet{
		name:          name,
		errorHandling: errorHandling,
	}
}

// NFlag returns the number of flags that have been set.
func (f *FlagSet) NFlag() int { return len(f.actual) }

// Args returns the non-flag arguments.
func (f *FlagSet) Args() []string { return f.args }

// Arg returns the i'th argument. Arg(0) is the first remaining argument
// after flags have been processed. Arg returns an empty string if the
// requested element does not exist.
func (f *FlagSet) Arg(i int) string {
	if i < 0 || i >= len(f.args) {
		return ""
	}
	return f.args[i]
}

// Lookup returns the Flag structure of the named flag, returning nil if none exists.
func (f *FlagSet) Lookup(name string) *Flag {
	return f.formal[name]
}

// VisitAll visits the flags in lexicographical order, calling fn for each.
// It visits all flags, even those not set.
func (f *FlagSet) VisitAll(fn func(*Flag)) {
	for _, flg := range sortFlags(f.formal) {
		fn(flg)
	}
}

// Set sets the value of the named flag.
func (f *FlagSet) Set(name, value string) error {
	flg, ok := f.formal[name]
	if !ok {
		return fmt.Errorf("no such option flag %q", name)
	}
	err := flg.Value.Set(value)
	if err != nil {
		return err
	}

	if f.actual == nil {
		f.actual = make(map[string]*Flag)
	}
	f.actual[name] = flg
	return nil
}

// Var defines a flag with the specified name and usage string.
// from the go flag.Var()
func (f *FlagSet) Var(value Value, name string, usage string) *Flag {
	if !cflag.IsGoodName(name) {
		panicf("option flag name '%s' is not a good name", name)
	}

	// Remember the default value as a string; it won't change.
	flg := &Flag{Name: name, Usage: usage, Value: value, DefValue: value.String()}
	_, exists := f.formal[name]
	if exists {
		var msg string
		if f.name == "" {
			msg = fmt.Sprintf("flag redefined: %s", name)
		} else {
			msg = fmt.Sprintf("%s flag redefined: %s", f.name, name)
		}
		panic(msg) // Happens only if flags are declared with identical names
	}

	if f.formal == nil {
		f.formal = make(map[string]*Flag)
	}
	f.formal[name] = flg
	return flg
}

// BoolVar defines a bool flag with specified name, default value, and usage string.
// from the go flag.BoolVar()
func (f *FlagSet) BoolVar(p *bool, name string, value bool, usage string) *Flag {
	return f.Var(newBoolValue(value, p), name, usage)
}

// Bool defines a bool flag with specified name, default value, and usage string.
// from the go flag.Bool()
func (f *FlagSet) Bool(name string, value bool, usage string) *bool {
	p := new(bool)
	f.BoolVar(p, name, value, usage)
	return p
}

// IntVar defines an int flag with specified name, default value, and usage string.
// from the go flag.IntVar()
func (f *FlagSet) IntVar(p *int, name string, value int, usage string) *Flag {
	return f.Var(newIntValue(value, p), name, usage)
}

// Int defines an int flag with specified name, default value, and usage string.
// from the go flag.Int()
func (f *FlagSet) Int(name string, value int, usage string) *int {
	p := new(int)
	f.IntVar(p, name, value, usage)
	return p
}

// Int64Var defines an int64 flag with specified name, default value, and usage string.
// from the go flag.Int64Var()
func (f *FlagSet) Int64Var(p *int64, name string, value int64, usage string) *Flag {
	return f.Var(newInt64Value(value, p), name, usage)
}

// Int64 defines an int64 flag with specified name, default value, and usage string.
// from the go flag.Int64()
func (f *FlagSet) Int64(name string, value int64, usage string) *int64 {
	p := new(int64)
	f.Int64Var(p, name, value, usage)
	return p
}

// UintVar defines a uint flag with specified name, default value, and usage string.
// from the go flag.UintVar()
func (f *FlagSet) UintVar(p *uint, name string, value uint, usage string) *Flag {
	return f.Var(newUintValue(value, p), name, usage)
}

// Uint defines a uint flag with specified name, default value, and usage string.
// from the go flag.Uint()
func (f *FlagSet) Uint(name string, value uint, usage string) *uint {
	p := new(uint)
	f.UintVar(p, name, value, usage)
	return p
}

// Uint64Var defines a uint64 flag with specified name, default value, and usage string.
// from the go flag.Uint64Var()
func (f *FlagSet) Uint64Var(p *uint64, name string, value uint64, usage string) *Flag {
	return f.Var(newUint64Value(value, p), name, usage)
}

// Uint64 defines a uint64 flag with specified name, default value, and usage string.
// from the go flag.Uint64()
func (f *FlagSet) Uint64(name string, value uint64, usage string) *uint64 {
	p := new(uint64)
	f.Uint64Var(p, name, value, usage)
	return p
}

// StringVar defines a string flag with specified name, default value, and usage string.
// from the go flag.StringVar()
func (f *FlagSet) StringVar(p *string, name, value, usage string) *Flag {
	return f.Var(newStringValue(value, p), name, usage)
}

// String defines a string flag with specified name, default value, and usage string.
// from the go flag.String()
func (f *FlagSet) String(name, value, usage string) *string {
	p := new(string)
	f.StringVar(p, name, value, usage)
	return p
}

// Float64Var defines a float64 flag with specified name, default value, and usage string.
// from the go flag.Float64Var()
func (f *FlagSet) Float64Var(p *float64, name string, value float64, usage string) *Flag {
	return f.Var(newFloat64Value(value, p), name, usage)
}

// Float64 defines a float64 flag with specified name, default value, and usage string.
// from the go flag.Float64()
func (f *FlagSet) Float64(name string, value float64, usage string) *float64 {
	p := new(float64)
	f.Float64Var(p, name, value, usage)
	return p
}

// DurationVar defines a time.Duration flag with specified name, default value, and usage string.
// from the go flag.DurationVar()
func (f *FlagSet) DurationVar(p *time.Duration, name string, value time.Duration, usage string) *Flag {
	return f.Var(newDurationValue(value, p), name, usage)
}

// Duration defines a time.Duration flag with specified name, default value, and usage string.
// from the go flag.Duration()
func (f *FlagSet) Duration(name string, value time.Duration, usage string) *time.Duration {
	p := new(time.Duration)
	f.DurationVar(p, name, value, usage)
	return p
}

// TextVar defines a text flag with specified name, default value, and usage string.
// from the go flag.TextVar()
func (f *FlagSet) TextVar(p encoding.TextUnmarshaler, name string, value encoding.TextMarshaler, usage string) *Flag {
	return f.Var(newTextValue(value, p), name, usage)
}

// Func defines a flag with the specified name and usage string.
// from the go flag.Func()
func (f *FlagSet) Func(name, usage string, fn func(string) error) *Flag {
	return f.Var(funcValue(fn), name, usage)
}

// Parsed reports whether f.Parse has been called.
func (f *FlagSet) Parsed() bool {
	return f.parsed
}

// Parse parses flag definitions from the argument list, which should not
// include the command name. Must be called after all flags in the FlagSet
// are defined and before flags are accessed by the program.
//
// The return value will be ErrHelp if -help or -h were set but not defined.
//
// NOTE: refer from flag.FlagSet#Parse()
func (f *FlagSet) Parse(arguments []string) error {
	f.parsed = true
	f.args = arguments
	for {
		seen, err := f.parseOne()
		if seen {
			continue
		}
		if err == nil {
			break
		}

		switch f.errorHandling {
		case flag.ContinueOnError:
			return err
		case flag.ExitOnError:
			if errors.Is(err, flag.ErrHelp) {
				os.Exit(0)
			}
			os.Exit(2)
		case flag.PanicOnError:
			panic(err)
		}
	}
	return nil
}

// parseOne parses one flag. It reports whether a flag was seen.
//
// NOTE: refer from flag.FlagSet#parseOne()
func (f *FlagSet) parseOne() (bool, error) {
	if len(f.args) == 0 {
		return false, nil
	}

	s := f.args[0]
	if len(s) < 2 || s[0] != '-' {
		return false, nil
	}

	numMinuses := 1
	if s[1] == '-' {
		numMinuses++
		if len(s) == 2 { // "--" terminates the flags
			f.args = f.args[1:]
			return false, nil
		}
	}
	name := s[numMinuses:]
	if len(name) == 0 || name[0] == '-' || name[0] == '=' {
		return false, fmt.Errorf("bad flag syntax: %s", s)
	}

	// it's a flag. does it have an argument?
	f.args = f.args[1:]
	hasValue := false
	value := ""
	for i := 1; i < len(name); i++ { // equals cannot be first
		if name[i] == '=' {
			value = name[i+1:]
			hasValue = true
			name = name[0:i]
			break
		}
	}

	// resolve shortcut name
	if rName, ok := f.shorts[name]; ok {
		name = rName
	}

	flg, ok := f.formal[name]
	if !ok {
		if name == "help" || name == "h" { // special case for nice help message.
			// f.usage()
			return false, flag.ErrHelp
		}
		return false, fmt.Errorf("option provided but not defined: %s", cflag.AddPrefix(name))
	}

	if fv, ok := flg.Value.(boolFlag); ok && fv.IsBoolFlag() { // special case: doesn't need an arg
		if hasValue {
			if err := fv.Set(value); err != nil {
				return false, fmt.Errorf("invalid boolean value %q for %s: %v", value, cflag.AddPrefix(name), err)
			}
		} else {
			if err := fv.Set("true"); err != nil {
				return false, fmt.Errorf("invalid boolean flag %s: %v", cflag.AddPrefix(name), err)
			}
		}
	} else {
		// It must have a value, which might be the next argument.
		if !hasValue && len(f.args) > 0 {
			// value is the next arg
			hasValue = true
			value, f.args = f.args[0], f.args[1:]
		}
		if !hasValue {
			return false, fmt.Errorf("flag option needs an argument: %s", cflag.AddPrefix(name))
		}
		if err := flg.Value.Set(value); err != nil {
			return false, fmt.Errorf("invalid value %q for option %s: %v", value, cflag.AddPrefix(name), err)
		}
	}

	if f.actual == nil {
		f.actual = make(map[string]*Flag)
	}
	f.actual[name] = flg
	return true, nil
}
