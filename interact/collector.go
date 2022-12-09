package interact

import (
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/structs"
	"github.com/gookit/goutil/strutil"
)

// InputParameter interface
type InputParameter interface {
	Type() string
	Name() string
	Desc() string
	Value() structs.Value
	Set(v string) error
	Run() error
}

// Collector information collector 信息收集者
// cli input values collector
type Collector struct {
	// input parameters
	ps  map[string]InputParameter
	ret maputil.Data
	err error

	ns []string
}

// NewCollector instance
func NewCollector() *Collector {
	return &Collector{
		ps:  make(map[string]InputParameter),
		ret: make(maputil.Data),
	}
}

// AddParams definitions at once.
func (c *Collector) AddParams(ps ...InputParameter) error {
	for _, p := range ps {
		if err := c.AddParam(p); err != nil {
			return err
		}
	}
	return nil
}

// Param get from collector
func (c *Collector) Param(name string) (InputParameter, bool) {
	p, ok := c.ps[name]
	return p, ok
}

// MustParam get from collector
func (c *Collector) MustParam(name string) InputParameter {
	p, ok := c.ps[name]
	if !ok {
		panic("not found the param: " + name)
	}

	return p
}

// AddParam to collector
func (c *Collector) AddParam(p InputParameter) error {
	if strutil.IsBlank(p.Name()) {
		return errorx.Raw("input parameter name cannot be empty")
	}

	name := p.Name()
	if arrutil.Contains(c.ns, name) {
		return errorx.Rawf("input parameter name %s has been exists", name)
	}

	c.ns = append(c.ns, name)
	c.ps[name] = p

	return nil
}

// Results for collector
func (c *Collector) Results() maputil.Data {
	return c.ret
}

// Run collector
func (c *Collector) Run() error {
	if len(c.ns) == 0 {
		return errorx.Raw("empty params definitions")
	}

	for _, name := range c.ns {
		p := c.ps[name]

		// has input value
		if c.ret.Has(name) {
			err := p.Set(c.ret.Str(name))
			if err != nil {
				return err
			}
			continue
		}

		// require input value
		if err := p.Run(); err != nil {
			return err
		}

		c.ret.Set(name, p.Value().V)
	}

	return nil
}
