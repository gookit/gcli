package cliapp

import (
	"flag"
	"fmt"
	"strconv"
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
func (s *Ints) Set(value int) error {
	*s = append(*s, value)
	return nil
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

/*************************************************************
 * command options
 *************************************************************/

// IntOpt binding a int option
func (c *Command) IntOpt(p *int, name string, short string, defValue int, description string) *Command {
	c.Flags.IntVar(p, name, defValue, description)

	if s, ok := c.addShortcut(name, short); ok {
		c.Flags.IntVar(p, s, defValue, "")
	}

	return c
}

// UintOpt binding a uint option
func (c *Command) UintOpt(p *uint, name string, short string, defValue uint, description string) *Command {
	c.Flags.UintVar(p, name, defValue, description)

	if s, ok := c.addShortcut(name, short); ok {
		c.Flags.UintVar(p, s, defValue, "")
	}

	return c
}

// StrOpt binding a string option
func (c *Command) StrOpt(p *string, name string, short string, defValue string, description string) *Command {
	c.Flags.StringVar(p, name, defValue, description)

	if s, ok := c.addShortcut(name, short); ok {
		c.Flags.StringVar(p, s, defValue, "")
	}

	return c
}

// BoolOpt binding a bool option
func (c *Command) BoolOpt(p *bool, name string, short string, defValue bool, description string) *Command {
	c.Flags.BoolVar(p, name, defValue, description)

	if s, ok := c.addShortcut(name, short); ok {
		c.Flags.BoolVar(p, s, defValue, "")
	}

	return c
}

// VarOpt binding a custom var option
// usage:
// 		cmd.VarOpt(&opts.Strings, "tables", "t", "Description ...")
func (c *Command) VarOpt(p flag.Value, name string, short string, description string) *Command {
	c.Flags.Var(p, name, description)

	if s, ok := c.addShortcut(name, short); ok {
		c.Flags.Var(p, s, "")
	}

	return c
}

// addShortcut add a shortcut name for a option name
func (c *Command) addShortcut(name string, short string) (string, bool) {
	if c.optNames == nil {
		c.optNames = map[string]string{}
	}

	// empty string
	if len(short) == 0 {
		c.optNames[name] = "" // record all option names
		return "", false
	}

	// first add
	if c.shortcuts == nil {
		c.shortcuts = map[string]string{}
	}

	// ensure it is one char
	short = string(short[0])
	if n, ok := c.shortcuts[short]; ok {
		exitWithErr("The shortcut name '%s' has been used by option '%s'", short, n)
	}

	c.optNames[name] = short
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

// Option is config info for a option
// usage:
// cmd.AddOpt(Option{
// 	Name: "name"
//	Short: "n"
// 	DType: "string"
// })
// cmd.Flags.String()
type Option struct {
	Name   string
	Short  string
	DType  string // int, string, bool, value
	VarPtr interface{}

	Required    bool
	DefValue    interface{}
	Description string
}

/*************************************************************
 * command arguments
 *************************************************************/

// Argument a command argument definition
type Argument struct {
	// Name argument name
	Name string
	// Description argument description message
	Description string
	// IsArray if is array, can allow accept multi values, and must in last.
	IsArray bool
	// Required arg is required
	Required bool
	// value store parsed argument data. (string, []string)
	Value interface{}
	// the argument position index in all arguments(cmd.args[index])
	index int
}

var emptyArg = &Argument{}

// Int argument value to int
func (a *Argument) Int(defVal ...int) int {
	def := 0
	if len(defVal) == 1 {
		def = defVal[0]
	}

	if a.Value == nil {
		return def
	}

	str := a.Value.(string)
	if str != "" {
		val, err := strconv.Atoi(str)
		if err != nil {
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

	if a.Value == nil {
		return def
	}

	return a.Value.(string)
}

// Strings argument value to string array, if argument isArray = true.
func (a *Argument) Strings() (ss []string) {
	if a.Value != nil {
		ss = a.Value.([]string)
	}

	return
}

// AddArg add a command argument.
// Notice:
// 	- Required argument cannot be defined after optional argument
//  - Only one array parameter is allowed
// 	- The (array) argument of multiple values ​​can only be defined at the end
//
// usage:
//	cmd.AddArg("name", "description")
//	cmd.AddArg("name", "description", true) // required
//	cmd.AddArg("name", "description", true, true) // required and is array
func (c *Command) AddArg(name, description string, requiredAndIsArray ...bool) *Argument {
	if c.argsIndexes == nil {
		c.argsIndexes = make(map[string]int)
	}

	if _, has := c.argsIndexes[name]; has {
		exitWithErr("the argument name '%s' already exists", name)
	}

	if c.hasArrayArg {
		exitWithErr("have defined an array argument, you can not add argument '%s'", name)
	}

	var isArray, required bool
	length := len(requiredAndIsArray)
	if length > 0 {
		required = requiredAndIsArray[0]

		if length > 1 {
			isArray = requiredAndIsArray[1]
		}
	}

	if required && c.hasOptionalArg {
		exitWithErr("required argument '%s' cannot be defined after optional argument", name)
	}

	// add argument index record
	argIndex := len(c.args)
	c.argsIndexes[name] = argIndex

	// add argument
	newArg := &Argument{
		Name: name, Description: description, Required: required, IsArray: isArray, index: argIndex,
	}
	c.args = append(c.args, newArg)

	if !required {
		c.hasOptionalArg = true
	}

	if isArray {
		c.hasArrayArg = true
	}

	return newArg
}

// Args get all defined argument
func (c *Command) Args() []*Argument {
	return c.args
}

// Arg get arg by defined name.
// usage:
// 	intVal := c.Arg("name").Int()
// 	strVal := c.Arg("name").String()
// 	arrVal := c.Arg("name").Array()
func (c *Command) Arg(name string) *Argument {
	i, ok := c.argsIndexes[name]
	if !ok {
		return emptyArg
	}

	return c.args[i]
}

// ArgByIndex get named arg by index
func (c *Command) ArgByIndex(i int) *Argument {
	if i <= len(c.args) {
		return c.args[i]
	}

	return emptyArg
}

// RawArgs get Flags args
func (c *Command) RawArgs() []string {
	return c.Flags.Args()
}

// RawArg get Flags arg value
func (c *Command) RawArg(i int) string {
	return c.Flags.Arg(i)
}
