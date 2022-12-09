package cparam

import (
	"github.com/gookit/color"
	"github.com/gookit/goutil/cliutil"
)

// StringParam definition
type StringParam struct {
	InputParam
}

// NewStringParam instance
func NewStringParam(name, desc string) *StringParam {
	return &StringParam{
		InputParam: InputParam{
			typ:  TypeStringParam,
			name: name,
			desc: desc,
		},
	}
}

// Config param
func (p *StringParam) Config(fn func(p *StringParam)) *StringParam {
	fn(p)
	return p
}

// Run param and get user input
func (p *StringParam) Run() (err error) {
	var val string
	if p.runFn != nil {
		val, err = p.runFn()
		if err != nil {
			return err
		}

		return p.Set(val)
	}

	val, err = cliutil.ReadLine(color.WrapTag(p.desc+"? ", "yellow"))
	if err != nil {
		return err
	}
	return p.Set(val)
}
