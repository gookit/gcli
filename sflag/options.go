package sflag

// Options the values for parsed arguments and options
type Options struct {
	// all option map
	optMap map[string]*Option
	// record relationship
	s2l map[string]string
	l2s map[string]string
}

// VisitAll options
func (r *Options) VisitAll() {

}

// Bindings all parsed long and short option values to Option
func (r *Options) Bindings(longs, shorts map[string]interface{}) {
	for _, opt := range r.optMap {
		if val, ok := shorts[opt.Short]; ok {
			opt.Value = val
		} else if val, ok := longs[opt.Name]; ok {
			opt.Value = val
		}
	}
}

// Opt get
func (r *Options) Opt(name string) *Option {
	return nil
}

// String the options
func (r *Options) String() string {
	return ""
}

// Option is config info for a option
// Usage:
// cmd.AddOpt(Option{
// 	Name: "name"
// 	Short: "n"
// 	DType: "string"
// })
// cmd.Flags.String()
type Option struct {
	Name  string
	Short string
	// Type value type. allow: int, string, bool, ints, strings, bools, custom
	Type string
	// Value of the option. allow: bool, string, slice(ints, strings)
	Value interface{}

	// refer flag.Value interface
	Setter func(string) error

	Required bool
	DefValue interface{}
	// Description
	Description string
}

// NewOpt create new option
func NewOpt(name, description string) *Option {
	return &Option{}
}

// StrVar binding
func (opt *Option) StrVar(s *string) *Option {
	return opt
}
