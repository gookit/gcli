package cparam

import (
	"github.com/gookit/gcli/v3/interact"
	"github.com/gookit/goutil/errorx"
)

// ChoicesParam definition
type ChoicesParam struct {
	InputParam
	// Choices for select
	Choices  []string
	selected []string
}

// NewChoicesParam instance
func NewChoicesParam(name, desc string) *ChoicesParam {
	return &ChoicesParam{
		InputParam: InputParam{
			typ:  TypeChoicesParam,
			name: name,
			desc: desc,
		},
	}
}

// WithChoices to param definition
func (p *ChoicesParam) WithChoices(Choices []string) *ChoicesParam {
	p.Choices = Choices
	return p
}

// Selected values get
func (p *ChoicesParam) Selected() []string {
	return p.val.Strings()
}

// Set value
func (p *ChoicesParam) Set(v string) error {
	if err := p.Valid(v); err != nil {
		return err
	}

	p.selected = append(p.selected, v)
	p.val.Set(p.selected)
	return nil
}

// Run param and get user input
func (p *ChoicesParam) Run() (err error) {
	if len(p.Choices) == 0 {
		return errorx.Raw("must provide items for choices")
	}

	s := interact.NewSelect(p.Desc(), p.Choices)
	s.EnableMulti()

	sr := s.Run()
	p.val.Set(sr.Val())
	return
}
