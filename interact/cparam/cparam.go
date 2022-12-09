package cparam

import (
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/structs"
)

// params types
const (
	TypeStringParam  = "string"
	TypeChoicesParam = "choices"
)

// RunFn func type
type RunFn func() (val string, err error)

// InputParam struct
type InputParam struct {
	typ  string
	name string
	desc string
	// Default value
	Default any
	ValidFn func(val string) error
	runFn   func() (val string, err error)
	// Value for input
	val structs.Value
	err error
}

// NewInputParam instance
func NewInputParam(typ, name, desc string) *InputParam {
	return &InputParam{
		typ:  typ,
		name: name,
		desc: desc,
	}
}

// Type name get
func (p *InputParam) Type() string {
	return p.typ
}

// Name get
func (p *InputParam) Name() string {
	return p.name
}

// Desc message
func (p *InputParam) Desc() string {
	return p.desc
}

// Valid value validate
func (p *InputParam) Valid(v string) error {
	if p.ValidFn != nil {
		return p.ValidFn(v)
	}
	return nil
}

// Set value and with validate
func (p *InputParam) Set(v string) error {
	if err := p.Valid(v); err != nil {
		return err
	}

	p.val.Set(v)
	return nil
}

// Value data get
func (p *InputParam) Value() structs.Value {
	return p.val
}

// Value get
func (p *InputParam) String() string {
	return p.val.String()
}

// SetFunc for run
func (p *InputParam) SetFunc(fn RunFn) {
	p.runFn = fn
}

// SetValidFn for run
func (p *InputParam) SetValidFn(fn func(val string) error) {
	p.ValidFn = fn
}

// Run param and get user input
func (p *InputParam) Run() (err error) {
	if p.runFn != nil {
		val, err := p.runFn()
		if err != nil {
			return err
		}

		err = p.Set(val)
	} else {
		err = errorx.Raw("please implement me")
	}

	return err
}
