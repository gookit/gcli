package gcli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gookit/goutil/strutil"
)

/*************************************************************************
 * options: some special flag vars
 * - implemented flag.Value interface
 *************************************************************************/

// Ints The int flag list, implemented flag.Value interface
type Ints []int

// String to string
func (s *Ints) String() string {
	return fmt.Sprintf("%v", *s)
}

// Set new value
func (s *Ints) Set(value string) error {
	intVal, err := strconv.Atoi(value)
	if err == nil {
		*s = append(*s, intVal)
	}

	return err
}

// Strings The string flag list, implemented flag.Value interface
type Strings []string

// String to string
func (s *Strings) String() string {
	return fmt.Sprintf("%v", *s)
}

// Set new value
func (s *Strings) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// Booleans The bool flag list, implemented flag.Value interface
type Booleans []bool

// String to string
func (s *Booleans) String() string {
	return fmt.Sprintf("%v", *s)
}

// Set new value
func (s *Booleans) Set(value string) error {
	boolVal, err := strconv.ParseBool(value)
	if err == nil {
		*s = append(*s, boolVal)
	}

	return err
}

// EnumString The string flag list, implemented flag.Value interface
type EnumString struct {
	val  string
	enum []string
}

// String to string
func (s *EnumString) String() string {
	return s.val
}

// SetEnum values
func (s *EnumString) SetEnum(enum []string) {
	s.enum = enum
}

// Set new value, will check value is right
func (s *EnumString) Set(value string) error {
	var ok bool
	for _, item := range s.enum {
		if value == item {
			ok = true
			break
		}
	}

	if !ok {
		return fmt.Errorf("value must one of the: %v", s.enum)
	}
	return nil
}

// String type, a special string
//
// Usage:
//
//	// case 1:
//	var names gcli.String
//	c.VarOpt(&names, "names", "", "multi name by comma split")
//
//	--names "tom,john,joy"
//	 names.Split(",") -> []string{"tom","john","joy"}
//
//	// case 2:
//	var ids gcli.String
//	c.VarOpt(&ids, "ids", "", "multi id by comma split")
//
//	--names "23,34,56"
//	 names.Ints(",") -> []int{23,34,56}
type String string

// Set value
func (s *String) Set(val string) error {
	*s = String(val)
	return nil
}

// String to string
func (s *String) String() string {
	return string(*s)
}

// Split value to []string
func (s *String) Split(sep string) []string {
	return strutil.ToStrings(string(*s), sep)
}

// Ints value to []int
func (s *String) Ints(sep string) []int {
	return strutil.Ints(string(*s), sep)
}

/*************************************************************************
 * verbose level
 *************************************************************************/

// VerbLevel type.
type VerbLevel uint

// Int verbose level to int.
func (vl *VerbLevel) Int() int {
	return int(*vl)
}

// String verbose level to string.
func (vl *VerbLevel) String() string {
	return fmt.Sprintf("%d=%s", *vl, vl.Name())
}

// Upper verbose level to string.
func (vl *VerbLevel) Upper() string {
	return strings.ToUpper(vl.Name())
}

// Name verbose level to string.
func (vl *VerbLevel) Name() string {
	switch *vl {
	case VerbQuiet:
		return "quiet"
	case VerbError:
		return "error"
	case VerbWarn:
		return "warn"
	case VerbInfo:
		return "info"
	case VerbDebug:
		return "debug"
	case VerbCrazy:
		return "crazy"
	}
	return "unknown"
}

// Set value from option binding.
func (vl *VerbLevel) Set(value string) error {
	// int: level value.
	if iv, err := strconv.Atoi(value); err == nil {
		if iv > int(VerbCrazy) {
			*vl = VerbCrazy
		} else if iv < 0 { // fallback to default level.
			*vl = DefaultVerb
		} else { // 0 - 5
			*vl = VerbLevel(iv)
		}

		return nil
	}

	// string: level name.
	*vl = name2verbLevel(value)
	return nil
}
