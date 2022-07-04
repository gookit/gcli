package gcli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gookit/goutil/strutil"
)

/*************************************************************
 * Arguments definition
 *************************************************************/

// an empty argument
var emptyArg = &Argument{}

// Arguments definition
type Arguments struct {
	// Inherited from Command
	name string
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
	argsIndexes map[string]int
	// validate the args number is right
	validateNum bool
	// mark exists array argument
	hasArrayArg bool
	// mark exists optional argument
	hasOptionalArg bool
}

// SetName for Arguments
func (ags *Arguments) SetName(name string) {
	ags.name = name
}

// SetValidateNum check
func (ags *Arguments) SetValidateNum(validateNum bool) {
	ags.validateNum = validateNum
}

// ParseArgs for Arguments
func (ags *Arguments) ParseArgs(args []string) (err error) {
	var num int
	inNum := len(args)

	for i, arg := range ags.args {
		// num is equals to "index + 1"
		num = i + 1
		if num > inNum { // not enough args
			if arg.Required {
				return fmt.Errorf("must set value for the argument: %s(position#%d)", arg.ShowName, arg.index)
			}
			break
		}

		if arg.Arrayed {
			err = arg.bindValue(args[i:])
			inNum = num // must reset inNum
		} else {
			err = arg.bindValue(args[i])
		}

		// has error on binding arg value
		if err != nil {
			return
		}
	}

	if ags.validateNum && inNum > num {
		return fmt.Errorf("entered too many arguments: %v", args[num:])
	}
	return
}

/*************************************************************
 * command arguments
 *************************************************************/

// AddArg binding an named argument for the command.
// Notice:
// 	- Required argument cannot be defined after optional argument
//  - Only one array parameter is allowed
// 	- The (array) argument of multiple values can only be defined at the end
//
// Usage:
// 	cmd.AddArg("name", "description")
// 	cmd.AddArg("name", "description", true) // required
// 	cmd.AddArg("names", "description", true, true) // required and is array
func (ags *Arguments) AddArg(name, desc string, requiredAndArrayed ...bool) *Argument {
	// create new argument
	newArg := NewArgument(name, desc, requiredAndArrayed...)

	return ags.AddArgument(newArg)
}

// AddArgByRule add an arg by simple string rule
func (ags *Arguments) AddArgByRule(name, rule string) *Argument {
	mp := parseSimpleRule(name, rule)

	required := strutil.MustBool(mp["required"])
	newArg := NewArgument(name, mp["desc"], required)

	if defVal := mp["default"]; defVal != "" {
		newArg.Value = defVal
	}

	return ags.AddArgument(newArg)
}

// BindArg alias of the AddArgument()
func (ags *Arguments) BindArg(arg *Argument) *Argument {
	return ags.AddArgument(arg)
}

// AddArgument binding a named argument for the command.
//
// Notice:
// 	- Required argument cannot be defined after optional argument
//  - Only one array parameter is allowed
// 	- The (array) argument of multiple values can only be defined at the end
//
func (ags *Arguments) AddArgument(arg *Argument) *Argument {
	if ags.argsIndexes == nil {
		ags.argsIndexes = make(map[string]int)
	}

	// validate argument name
	name := arg.goodArgument()
	if _, has := ags.argsIndexes[name]; has {
		panicf("the argument name '%s' already exists in command '%s'", name, ags.name)
	}

	if ags.hasArrayArg {
		panicf("have defined an array argument, you cannot add argument '%s'", name)
	}

	if arg.Required && ags.hasOptionalArg {
		panicf("required argument '%s' cannot be defined after optional argument", name)
	}

	// add argument index record
	arg.index = len(ags.args)
	ags.argsIndexes[name] = arg.index

	// add argument
	ags.args = append(ags.args, arg)
	if !arg.Required {
		ags.hasOptionalArg = true
	}

	if arg.Arrayed {
		ags.hasArrayArg = true
	}

	return arg
}

// Args get all defined argument
func (ags *Arguments) Args() []*Argument {
	return ags.args
}

// HasArg check named argument is defined
func (ags *Arguments) HasArg(name string) bool {
	_, ok := ags.argsIndexes[name]
	return ok
}

// HasArgs defined. alias of the HasArguments()
func (ags *Arguments) HasArgs() bool {
	return len(ags.argsIndexes) > 0
}

// HasArguments defined
func (ags *Arguments) HasArguments() bool {
	return len(ags.argsIndexes) > 0
}

// Arg get arg by defined name.
//
// Usage:
// 	intVal := ags.Arg("name").Int()
// 	strVal := ags.Arg("name").String()
// 	arrVal := ags.Arg("names").Array()
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
	// Name argument name. it's required
	Name string
	// Desc argument description message
	Desc string
	// Type name. eg: string, int, array
	// Type string
	// ShowName is a name for display help. default is equals to Name.
	ShowName string
	// Required arg is required
	Required bool
	// Arrayed if is array, can allow accept multi values, and must in last.
	Arrayed bool
	// value *goutil.Value TODO ...
	// value store parsed argument data. (type: string, []string)
	Value interface{}
	// Handler custom argument value handler on call GetValue()
	Handler func(val interface{}) interface{}
	// Validator you can add an validator, will call it on binding argument value
	Validator func(val interface{}) (interface{}, error)
	// the argument position index in all arguments(cmd.args[index])
	index int
}

// NewArgument quick create an new command argument
func NewArgument(name, desc string, requiredAndArrayed ...bool) *Argument {
	var arrayed, required bool
	if ln := len(requiredAndArrayed); ln > 0 {
		required = requiredAndArrayed[0]
		if ln > 1 {
			arrayed = requiredAndArrayed[1]
		}
	}

	return &Argument{
		Name: name,
		Desc: desc,
		// other options
		// ShowName: name,
		Required: required,
		Arrayed:  arrayed,
	}
}

// SetArrayed the argument
func (a *Argument) SetArrayed() *Argument {
	a.Arrayed = true
	return a
}

// Init the argument
func (a *Argument) Init() *Argument {
	a.goodArgument()
	return a
}

func (a *Argument) goodArgument() string {
	name := strings.TrimSpace(a.Name)
	if name == "" {
		panicf("the command argument name cannot be empty")
	}

	if !goodName.MatchString(name) {
		panicf("the command argument name '%s' is invalid, must match: %s", name, regGoodName)
	}

	a.Name = name
	if a.ShowName == "" {
		a.ShowName = name
	}
	return name
}

// HelpName for render help message
func (a *Argument) HelpName() string {
	if a.Arrayed {
		return a.ShowName + "..."
	}

	return a.ShowName
}

// With an func for config the argument
func (a *Argument) With(fn func(arg *Argument)) *Argument {
	if fn != nil {
		fn(a)
	}

	return a
}

// WithValidator set an value validator of the argument
func (a *Argument) WithValidator(fn func(interface{}) (interface{}, error)) *Argument {
	a.Validator = fn
	return a
}

// SetValue set an validated value
func (a *Argument) SetValue(val interface{}) error {
	return a.bindValue(val)
}

// GetValue get value by custom handler func
func (a *Argument) GetValue() interface{} {
	val := a.Value
	if a.Handler != nil {
		return a.Handler(val)
	}

	return val
}

// Int argument value to int
func (a *Argument) Int(defVal ...int) int {
	def := 0
	if len(defVal) == 1 {
		def = defVal[0]
	}

	if a.Value == nil || a.Arrayed {
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

	if a.Value == nil || a.Arrayed {
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
	if a.Value != nil && a.Arrayed {
		ss = a.Value.([]string)
	}

	return
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

// bind an value to the argument
func (a *Argument) bindValue(val interface{}) (err error) {
	// has validator
	if a.Validator != nil {
		val, err = a.Validator(val)
		if err == nil {
			a.Value = val
		}
	} else {
		a.Value = val
	}
	return
}
