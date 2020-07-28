package gcli

import (
	"regexp"
	"strconv"
	"strings"
)

/*************************************************************
 * command arguments
 *************************************************************/

// AddArg binding an named argument for the command.
// Notice:
// 	- Required argument cannot be defined after optional argument
//  - Only one array parameter is allowed
// 	- The (array) argument of multiple values ​​can only be defined at the end
//
// Usage:
// 	cmd.AddArg("name", "description")
// 	cmd.AddArg("name", "description", true) // required
// 	cmd.AddArg("names", "description", true, true) // required and is array
func (c *Command) AddArg(name, description string, requiredAndIsArray ...bool) *Argument {
	// create new argument
	newArg := NewArgument(name, description, requiredAndIsArray...)

	return c.AddArgument(newArg)
}

// BindArg alias of the AddArgument()
func (c *Command) BindArg(arg *Argument) *Argument {
	return c.AddArgument(arg)
}

// AddArgument binding an named argument for the command.
//
// Notice:
// 	- Required argument cannot be defined after optional argument
//  - Only one array parameter is allowed
// 	- The (array) argument of multiple values ​​can only be defined at the end
//
func (c *Command) AddArgument(arg *Argument) *Argument {
	if c.argsIndexes == nil {
		c.argsIndexes = make(map[string]int)
	}

	// validate argument
	arg.goodArgument()

	name := arg.Name
	if _, has := c.argsIndexes[name]; has {
		panicf("the argument name '%s' already exists in command '%s'", name, c.Name)
	}

	if c.hasArrayArg {
		panicf("have defined an array argument, you cannot add argument '%s'", name)
	}

	if arg.Required && c.hasOptionalArg {
		panicf("required argument '%s' cannot be defined after optional argument", name)
	}

	// add argument index record
	arg.index = len(c.args)
	c.argsIndexes[name] = arg.index

	// add argument
	c.args = append(c.args, arg)
	if !arg.Required {
		c.hasOptionalArg = true
	}

	if arg.IsArray {
		c.hasArrayArg = true
	}

	return arg
}

// Args get all defined argument
func (c *Command) Args() []*Argument {
	return c.args
}

// Arg get arg by defined name.
// Usage:
// 	intVal := c.Arg("name").Int()
// 	strVal := c.Arg("name").String()
// 	arrVal := c.Arg("names").Array()
func (c *Command) Arg(name string) *Argument {
	i, ok := c.argsIndexes[name]
	if !ok {
		return emptyArg
	}
	return c.args[i]
}

// ArgByIndex get named arg by index
func (c *Command) ArgByIndex(i int) *Argument {
	if i < len(c.args) {
		return c.args[i]
	}
	return emptyArg
}

/*************************************************************
 * Argument definition
 *************************************************************/

// Arguments definition
type Arguments struct {
	// args definition for a command.
	// eg. {
	// 	{"arg0", "this is first argument", false, false},
	// 	{"arg1", "this is second argument", false, false},
	// }
	args []*Argument
	// record min length for args
	// argsMinLen int
	// record argument names and defined positional relationships
	// {
	// 	// name: position
	// 	"arg0": 0,
	// 	"arg1": 1,
	// }
	argsIndexes  map[string]int
	hasArrayable bool
	hasOptional  bool
}

// Add a new argument
func (ags *Arguments) Add(name, description string) {
	// todo ...
}

// Args get all defined argument
func (ags *Arguments) Args() []*Argument {
	return ags.args
}

// Arg get arg by defined name.
// Usage:
// 	intVal := c.Arg("name").Int()
// 	strVal := c.Arg("name").String()
// 	arrVal := c.Arg("names").Array()
func (ags *Arguments) Arg(name string) *Argument {
	i, ok := ags.argsIndexes[name]
	if !ok {
		return emptyArg
	}
	return ags.args[i]
}

// ArgByIndex get named arg by index
func (ags *Arguments) ArgByIndex(i int) *Argument {
	if i < len(ags.args) {
		return ags.args[i]
	}
	return emptyArg
}

/*************************************************************
 * Argument definition
 *************************************************************/

// Argument a command argument definition
type Argument struct {
	// valWrapper Value TODO ...
	// Name argument name. it's required
	Name string
	// ShowName is a name for display help. default is equals to Name.
	ShowName string
	// Description argument description message
	Description string
	// Required arg is required
	Required bool
	// IsArray if is array, can allow accept multi values, and must in last.
	IsArray bool
	// value store parsed argument data. (type: string, []string)
	Value interface{}
	// Handler custom argument value parse handler
	Handler func(val interface{}) interface{}
	// Validator you can add an validator, will call it on binding argument value
	Validator func(val interface{}) (interface{}, error)
	// the argument position index in all arguments(cmd.args[index])
	index int
}

var (
	emptyArg    = &Argument{}
	goodArgName = regexp.MustCompile(`^[\w-]+$`)
)

// NewArgument quick create an new command argument
func NewArgument(name, description string, requiredAndIsArray ...bool) *Argument {
	var isArray, required bool

	length := len(requiredAndIsArray)
	if length > 0 {
		required = requiredAndIsArray[0]

		if length > 1 {
			isArray = requiredAndIsArray[1]
		}
	}

	// create new argument
	return &Argument{
		Name: name, ShowName: name, Description: description, Required: required, IsArray: isArray,
	}
}

func (a *Argument) goodArgument() {
	a.Name = strings.TrimSpace(a.Name)
	if a.Name == "" {
		panicf("the command argument name cannot be empty")
	}

	if !goodArgName.MatchString(a.Name) {
		panicf("the command argument name '%s' is invalid, only allow: a-Z 0-9 _ -", a.Name)
	}
}

// Config the argument
func (a *Argument) Config(fn func(arg *Argument)) {
	if fn != nil {
		fn(a)
	}
}

// WithValidator set an value validator of the argument
func (a *Argument) WithValidator(fn func(interface{}) (interface{}, error)) *Argument {
	a.Validator = fn
	return a
}

// WithValue set an value of the argument
func (a *Argument) WithValue(val interface{}) *Argument {
	a.Value = val
	return a
}

// Int argument value to int
func (a *Argument) Int(defVal ...int) int {
	def := 0
	if len(defVal) == 1 {
		def = defVal[0]
	}

	if a.Value == nil || a.IsArray {
		return def
	}

	if intVal, ok := a.Value.(int); ok {
		return intVal
	}

	if str, ok := a.Value.(string); ok {
		val, err := strconv.Atoi(str)
		if err == nil {
			return val
		}
	}
	return def
}

// String argument value to string
func (a *Argument) String(defVal ...string) string {
	def := ""
	if len(defVal) == 1 {
		def = defVal[0]
	}

	if a.Value == nil || a.IsArray {
		return def
	}

	if str, ok := a.Value.(string); ok {
		return str
	}
	return def
}

// StringSplit quick split a string argument to string slice
func (a *Argument) StringSplit(sep ...string) (ss []string) {
	str := a.String()
	if str == "" {
		return
	}

	char := ","
	if len(sep) > 0 {
		char = sep[0]
	}

	return strings.Split(str, char)
}

// Array alias of the Strings()
func (a *Argument) Array() (ss []string) {
	return a.Strings()
}

// Strings argument value to string array, if argument isArray = true.
func (a *Argument) Strings() (ss []string) {
	if a.Value != nil && a.IsArray {
		ss = a.Value.([]string)
	}

	return
}

// GetValue get value by custom handler func
func (a *Argument) GetValue() interface{} {
	val := a.Value
	if a.Handler != nil {
		return a.Handler(val)
	}

	return val
}

// HasValue value is empty
func (a *Argument) HasValue() bool {
	return a.Value != nil
}

// IsEmpty argument is empty
func (a *Argument) IsEmpty() bool {
	return a.Name == ""
}

// Index get argument index in the command
func (a *Argument) Index() int {
	return a.index
}

// bind an value of the argument
func (a *Argument) bindValue(val interface{}) (err error) {
	// has validator
	if a.Validator != nil {
		val, err = a.Validator(val)
	}

	a.Value = val
	return
}
