package gflag

import (
	"flag"
	"time"
)

// BindVar binds a typed pointer as an option in a type-safe, generic way.
//
// It dispatches on the concrete type of ptr to the matching typed binder, so a
// single call replaces the per-type BoolVar/IntVar/StrVar/... methods. The opt
// carries the option metadata (name, shorts, desc, default, validator, ...).
//
// Supported T: bool, int, int64, uint, uint64, float64, string, time.Duration,
// []string, []int, []bool, map[string]string, and any type whose pointer
// implements flag.Value. Other types panic with a clear message.
func BindVar[T any](fs *Parser, ptr *T, opt *CliOpt) {
	switch pv := any(ptr).(type) {
	case *bool:
		fs.BoolVar(pv, opt)
	case *int:
		fs.IntVar(pv, opt)
	case *int64:
		fs.Int64Var(pv, opt)
	case *uint:
		fs.UintVar(pv, opt)
	case *uint64:
		fs.Uint64Var(pv, opt)
	case *float64:
		fs.Float64Var(pv, opt)
	case *string:
		fs.StrVar(pv, opt)
	case *time.Duration:
		fs.DurationVar(pv, opt)
	case *[]string:
		fs.Var((*Strings)(pv), opt)
	case *[]int:
		fs.Var((*Ints)(pv), opt)
	case *[]bool:
		fs.Var((*Booleans)(pv), opt)
	case *map[string]string:
		fs.Var(&mapStrValue{ref: pv, sep: "="}, opt)
	case flag.Value:
		// T's pointer already implements flag.Value (e.g. a custom value type)
		fs.Var(pv, opt)
	default:
		panicf("BindVar: unsupported type %T for option %q", ptr, opt.Name)
	}
}

// Opt binds a typed option with name/shorts/default/desc in one generic call.
//
//	var name string
//	gflag.Opt(fs, &name, "name", "n", "tom", "the user name")
//
//	var tags []string
//	gflag.Opt(fs, &tags, "tag", "t", nil, "the tags, repeatable")
//
// NOTE: the default value is applied for scalar and time.Duration types; for
// slice / map types the default is the zero value (use the field directly to
// pre-fill). For richer per-option config use the setFns (WithValidator, ...).
func Opt[T any](fs *Parser, ptr *T, name, shorts string, defVal T, desc string, setFns ...CliOptFn) {
	BindVar(fs, ptr, newOpt(name, desc, defVal, shorts, setFns...))
}
