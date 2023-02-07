package gflag

import (
	"fmt"
	"strings"

	"github.com/gookit/gcli/v3/helper"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/mathutil"
	"github.com/gookit/goutil/structs"
	"github.com/gookit/goutil/strutil"
)

/*************************************************************
 * Cli Arguments definition
 *************************************************************/

// Arguments alias of CliArgs
type Arguments = CliArgs

// CliArgs definition
type CliArgs struct {
	// name inherited from gcli.Command
	name string
	// args definition for a command.
	//
	// eg. [
	// 	{"arg0", "this is first argument", false, false},
	// 	{"arg1", "this is second argument", false, false},
	// ]
	args []*CliArg
	// arg name max width
	argWidth int
	// record min length for args
	// argsMinLen int
	// record argument names and defined positional relationships
	//
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
	// remain extra args after parse named arguments
	remainArgs []string
}

// SetName for CliArgs
func (ags *CliArgs) SetName(name string) {
	ags.name = name
}

// SetValidateNum check
func (ags *CliArgs) SetValidateNum(validateNum bool) {
	ags.validateNum = validateNum
}

// ParseArgs for CliArgs
func (ags *CliArgs) ParseArgs(args []string) (err error) {
	var num int // parsed num
	inNum := len(args)

	for i, arg := range ags.args {
		// num is equals to "index + 1"
		num = i + 1
		if num > inNum { // not enough args
			if arg.Required {
				return errorx.Rawf("must set value for the argument: %s(position#%d)", arg.ShowName, arg.index)
			}
			num = i
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

	if inNum > num {
		if ags.validateNum {
			return errorx.Rawf("entered too many arguments: %v", args[num:])
		}
		ags.remainArgs = args[num:]
	}
	return
}

/*************************************************************
 * command arguments
 *************************************************************/

// AddArg binding a named argument for the command.
//
// Notice:
//   - Required argument cannot be defined after optional argument
//   - Only one array parameter is allowed
//   - The (array) argument of multiple values can only be defined at the end
//
// Usage:
//
//	cmd.AddArg("name", "description")
//	cmd.AddArg("name", "description", true) // required
//	cmd.AddArg("names", "description", true, true) // required and is arrayed
func (ags *CliArgs) AddArg(name, desc string, requiredAndArrayed ...bool) *CliArg {
	newArg := NewArgument(name, desc, requiredAndArrayed...)
	return ags.AddArgument(newArg)
}

// AddArgByRule add an arg by simple string rule
func (ags *CliArgs) AddArgByRule(name, rule string) *CliArg {
	mp := ParseSimpleRule(name, rule)

	required := strutil.QuietBool(mp["required"])
	newArg := NewArgument(name, mp["desc"], required)

	if defVal := mp["default"]; defVal != "" {
		newArg.Set(defVal)
	}

	return ags.AddArgument(newArg)
}

// BindArg alias of the AddArgument()
func (ags *CliArgs) BindArg(arg *CliArg) *CliArg {
	return ags.AddArgument(arg)
}

// AddArgument binding a named argument for the command.
//
// Notice:
//   - Required argument cannot be defined after optional argument
//   - Only one array parameter is allowed
//   - The (array) argument of multiple values can only be defined at the end
func (ags *CliArgs) AddArgument(arg *CliArg) *CliArg {
	if ags.argsIndexes == nil {
		ags.argWidth = 12 // default width
		ags.argsIndexes = make(map[string]int)
	}

	// validate argument name
	name := arg.goodArgument()
	if _, has := ags.argsIndexes[name]; has {
		helper.Panicf("the argument name '%s' already exists in command '%s'", name, ags.name)
	}

	if ags.hasArrayArg {
		helper.Panicf("have defined an array argument, you cannot add argument '%s'", name)
	}

	if arg.Required && ags.hasOptionalArg {
		helper.Panicf("required argument '%s' cannot be defined after optional argument", name)
	}

	// add argument index record
	arg.index = len(ags.args)
	ags.argsIndexes[name] = arg.index
	ags.argWidth = mathutil.MaxInt(ags.argWidth, len(name))

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
func (ags *CliArgs) Args() []*CliArg {
	return ags.args
}

// HasArg check named argument is defined
func (ags *CliArgs) HasArg(name string) bool {
	_, ok := ags.argsIndexes[name]
	return ok
}

// HasArgs defined. alias of the HasArguments()
func (ags *CliArgs) HasArgs() bool {
	return len(ags.argsIndexes) > 0
}

// HasArguments defined
func (ags *CliArgs) HasArguments() bool {
	return len(ags.argsIndexes) > 0
}

// Arg get arg by defined name.
//
// Usage:
//
//	intVal := ags.Arg("name").Int()
//	strVal := ags.Arg("name").String()
//	arrVal := ags.Arg("names").Array()
func (ags *CliArgs) Arg(name string) *CliArg {
	i, ok := ags.argsIndexes[name]
	if !ok {
		helper.Panicf("get not exists argument '%s'", name)
	}
	return ags.args[i]
}

// ArgByIndex get named arg by index
func (ags *CliArgs) ArgByIndex(i int) *CliArg {
	if i >= len(ags.args) {
		helper.Panicf("get not exists argument #%d", i)
	}
	return ags.args[i]
}

// String build args help string
func (ags *CliArgs) String() string {
	return ags.BuildArgsHelp()
}

// BuildArgsHelp string
func (ags *CliArgs) BuildArgsHelp() string {
	if len(ags.args) < 1 {
		return ""
	}

	var sb strings.Builder
	for _, arg := range ags.args {
		sb.WriteString(fmt.Sprintf(
			"<info>%s</> %s%s\n",
			strutil.PadRight(arg.HelpName(), " ", ags.argWidth),
			getRequiredMark(arg.Required),
			strutil.UpperFirst(arg.Desc),
		))
	}

	return sb.String()
}

// ExtraArgs remain extra args after collect parse.
func (ags *CliArgs) ExtraArgs() []string {
	return ags.remainArgs
}

/*************************************************************
 * Cli Argument definition
 *************************************************************/

// Argument alias of CliArg
type Argument = CliArg

// CliArg a command argument definition
type CliArg struct {
	*structs.Value
	// the argument position index in all arguments(cmd.args[index])
	index int
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
	// Arrayed if is array, can allow to accept multi values, and must in last.
	Arrayed bool

	// Handler custom argument value handler on call GetValue()
	Handler func(val any) any
	// Validator you can add a validator, will call it on binding argument value
	Validator func(val any) (any, error)
}

// NewArg quick create a new command argument
func NewArg(name, desc string, val any, requiredAndArrayed ...bool) *CliArg {
	var arrayed, required bool
	if ln := len(requiredAndArrayed); ln > 0 {
		required = requiredAndArrayed[0]
		if ln > 1 {
			arrayed = requiredAndArrayed[1]
		}
	}

	return &CliArg{
		Name:  name,
		Desc:  desc,
		Value: structs.NewValue(val),
		// other settings
		// ShowName: name,
		Required: required,
		Arrayed:  arrayed,
	}
}

// NewArgument quick create a new command argument
func NewArgument(name, desc string, requiredAndArrayed ...bool) *CliArg {
	return NewArg(name, desc, nil, requiredAndArrayed...)
}

// SetArrayed the argument
func (a *CliArg) SetArrayed() *CliArg {
	a.Arrayed = true
	return a
}

// WithValue to the argument
func (a *CliArg) WithValue(val any) *CliArg {
	a.Value.Set(val)
	return a
}

// WithFn a func for config the argument
func (a *CliArg) WithFn(fn func(arg *CliArg)) *CliArg {
	if fn != nil {
		fn(a)
	}
	return a
}

// WithValidator set a value validator of the argument
func (a *CliArg) WithValidator(fn func(any) (any, error)) *CliArg {
	a.Validator = fn
	return a
}

// SetValue set an validated value
func (a *CliArg) SetValue(val any) error {
	return a.bindValue(val)
}

// Init the argument
func (a *CliArg) Init() *CliArg {
	a.goodArgument()
	return a
}

func (a *CliArg) goodArgument() string {
	name := strings.TrimSpace(a.Name)
	if name == "" {
		helper.Panicf("the command argument name cannot be empty")
	}

	if !helper.IsGoodName(name) {
		helper.Panicf("the argument name '%s' is invalid, must match: %s", name, helper.RegGoodName)
	}

	a.Name = name
	if a.ShowName == "" {
		a.ShowName = name
	}

	if a.Value == nil {
		a.Value = structs.NewValue(nil)
	}
	return name
}

// GetValue get value by custom handler func
func (a *CliArg) GetValue() any {
	val := a.Value.Val()
	if a.Handler != nil {
		return a.Handler(val)
	}
	return val
}

// Array alias of the Strings()
func (a *CliArg) Array() (ss []string) {
	return a.Strings()
}

// HasValue value is empty
func (a *CliArg) HasValue() bool {
	return a.V != nil
}

// Index get argument index in the command
func (a *CliArg) Index() int {
	return a.index
}

// HelpName for render help message
func (a *CliArg) HelpName() string {
	if a.Arrayed {
		return a.ShowName + "..."
	}
	return a.ShowName
}

// bind a value to the argument
func (a *CliArg) bindValue(val any) (err error) {
	if a.Validator != nil {
		val, err = a.Validator(val)
		if err != nil {
			return
		}
	}

	if a.Handler != nil {
		val = a.Handler(val)
	}

	a.Value.V = val
	return
}
