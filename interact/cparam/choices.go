package cparam

import (
	"github.com/gookit/gcli/v3/interact"
	"github.com/gookit/goutil/errorx"
)

// ChoiceParam definition
type ChoiceParam struct {
	InputParam
	// Choices for select
	Choices  []string
	selected []string
}

// NewChoiceParam instance
func NewChoiceParam(name, desc string) *ChoiceParam {
	return &ChoiceParam{
		InputParam: InputParam{
			typ:  TypeChoicesParam,
			name: name,
			desc: desc,
		},
	}
}

// WithChoices to param definition
func (p *ChoiceParam) WithChoices(Choices []string) *ChoiceParam {
	p.Choices = Choices
	return p
}

// Selected values get
func (p *ChoiceParam) Selected() []string {
	return p.val.Strings()
}

// Set value
func (p *ChoiceParam) Set(v string) error {
	if err := p.Valid(v); err != nil {
		return err
	}

	p.selected = append(p.selected, v)
	p.val.Set(p.selected)
	return nil
}

// Run param and get user input
func (p *ChoiceParam) Run() (err error) {
	if len(p.Choices) == 0 {
		return errorx.Raw("must provide items for choices")
	}

	s := interact.NewSelect(p.Desc(), p.Choices)
	s.EnableMulti()

	sr := s.Run()
	p.val.Set(sr.Val())
	return
}
