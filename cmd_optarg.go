package gcli

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

/*************************************************************
 * options: some special vars
 *************************************************************/

// Ints The int flag list, implemented flag.Value interface
type Ints []int

// String to string
func (s *Ints) String() string {
	return fmt.Sprintf("%v", *s)
}

// Set new value
func (s *Ints) Set(value string) error {
	intVal, err := strconv.Atoi(value)
	if err == nil {
		*s = append(*s, intVal)
	}

	return err
}

// Strings The string flag list, implemented flag.Value interface
type Strings []string

// String to string
func (s *Strings) String() string {
	return fmt.Sprintf("%v", *s)
}

// Set new value
func (s *Strings) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// Booleans The bool flag list, implemented flag.Value interface
type Booleans []bool

// String to string
func (s *Booleans) String() string {
	return fmt.Sprintf("%v", *s)
}

// Set new value
func (s *Booleans) Set(value string) error {
	boolVal, err := strconv.ParseBool(value)
	if err == nil {
		*s = append(*s, boolVal)
	}

	return err
}

/*************************************************************
 * command options
 *************************************************************/

// IntOpt binding a int option
func (c *Command) IntOpt(p *int, name, short string, defValue int, description string) *Command {
	c.Flags.IntVar(p, name, defValue, description)

	if s, ok := c.addShortcut(name, short); ok {
		c.Flags.IntVar(p, s, defValue, "")
	}

	return c
}

// UintOpt binding a uint option
func (c *Command) UintOpt(p *uint, name, short string, defValue uint, description string) *Command {
	c.Flags.UintVar(p, name, defValue, description)

	if s, ok := c.addShortcut(name, short); ok {
		c.Flags.UintVar(p, s, defValue, "")
	}

	return c
}

// StrOpt binding a string option
func (c *Command) StrOpt(p *string, name, short string, defValue, description string) *Command {
	c.Flags.StringVar(p, name, defValue, description)

	if s, ok := c.addShortcut(name, short); ok {
		c.Flags.StringVar(p, s, defValue, "")
	}

	return c
}

// BoolOpt binding a bool option
func (c *Command) BoolOpt(p *bool, name, short string, defValue bool, description string) *Command {
	c.Flags.BoolVar(p, name, defValue, description)

	if s, ok := c.addShortcut(name, short); ok {
		c.Flags.BoolVar(p, s, defValue, "")
	}

	return c
}

// VarOpt binding a custom var option
// Usage:
// 		cmd.VarOpt(&opts.Strings, "tables", "t", "description ...")
func (c *Command) VarOpt(p flag.Value, name string, short string, description string) *Command {
	c.Flags.Var(p, name, description)

	if s, ok := c.addShortcut(name, short); ok {
		c.Flags.Var(p, s, "")
	}

	return c
}

// addShortcut add a shortcut name for a option name
func (c *Command) addShortcut(name, short string) (string, bool) {
	// record all option names
	if c.optNames == nil {
		c.optNames = map[string]string{}
	}
	c.optNames[name] = ""

	// empty string
	if len(short) == 0 {
		return "", false
	}

	// first add short
	if c.shortcuts == nil {
		c.shortcuts = map[string]string{}
	}

	// ensure it is one char
	short = string(short[0])
	if n, ok := c.shortcuts[short]; ok {
		panicf("The shortcut name '%s' has been used by option '%s'", short, n)
	}

	c.shortcuts[short] = name
	return short, true
}

// isShortOpt alias of the `isShortcut`
func (c *Command) isShortOpt(short string) bool {
	return c.isShortcut(short)
}

// isShortcut check it is a shortcut name
func (c *Command) isShortcut(short string) bool {
	if len(short) != 1 {
		return false
	}

	_, ok := c.shortcuts[short]
	return ok
}

// ShortName get a shortcut name by option name
func (c *Command) ShortName(name string) string {
	for s, n := range c.shortcuts {
		if n == name {
			return s
		}
	}

	return ""
}

// OptFlag get option Flag by option name
func (c *Command) OptFlag(name string) *flag.Flag {
	if _, ok := c.optNames[name]; ok {
		return c.Flags.Lookup(name)
	}

	return nil
}

// OptDes get option description by option name
func (c *Command) OptDes(name string) string {
	if _, ok := c.optNames[name]; ok {
		return c.Flags.Lookup(name).Usage
	}

	return ""
}

// OptNames return all option names
func (c *Command) OptNames() map[string]string {
	return c.optNames
}

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

// RawArgs get all raw arguments
func (c *Command) RawArgs() []string {
	return c.Flags.Args()
}

// RawArg get an argument value by index
func (c *Command) RawArg(i int) string {
	return c.Flags.Arg(i)
}

/*************************************************************
 * Argument definition
 *************************************************************/

// Argument a command argument definition
type Argument struct {
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
