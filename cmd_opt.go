package cliapp

import (
	"flag"
	"github.com/gookit/color"
)

// IntOpt set a int option
func (c *Command) IntOpt(p *int, name string, short string, defValue int, description string) *Command {
	c.Flags.IntVar(p, name, defValue, description)

	if s, ok := c.addShortcut(name, short); ok {
		c.Flags.IntVar(p, s, defValue, "")
	}

	return c
}

// UintOpt set a int option
func (c *Command) UintOpt(p *uint, name string, short string, defValue uint, description string) *Command {
	c.Flags.UintVar(p, name, defValue, description)

	if s, ok := c.addShortcut(name, short); ok {
		c.Flags.UintVar(p, s, defValue, "")
	}

	return c
}

// StrOpt set a str option
func (c *Command) StrOpt(p *string, name string, short string, defValue string, description string) *Command {
	c.Flags.StringVar(p, name, defValue, description)

	if s, ok := c.addShortcut(name, short); ok {
		c.Flags.StringVar(p, s, defValue, "")
	}

	return c
}

// BoolOpt set a bool option
func (c *Command) BoolOpt(p *bool, name string, short string, defValue bool, description string) *Command {
	c.Flags.BoolVar(p, name, defValue, description)

	if s, ok := c.addShortcut(name, short); ok {
		c.Flags.BoolVar(p, s, defValue, "")
	}

	return c
}

// VarOpt set a custom option
// raw usage:
// cmd.Flags.Var(&opts.Strings, "tables", "Description ...")
// in here:
// cmd.VarOpt(&opts.Strings, "tables", "t", "Description ...")
func (c *Command) VarOpt(p flag.Value, name string, short string, description string) *Command {
	c.Flags.Var(p, name, description)

	if s, ok := c.addShortcut(name, short); ok {
		c.Flags.Var(p, s, "")
	}

	return c
}

// addShortcut add a shortcut name for a option name
func (c *Command) addShortcut(name string, short string) (string, bool) {
	// first add
	if c.optNames == nil {
		c.optNames = map[string]string{}
	}

	// empty string
	if len(short) == 0 {
		c.optNames[name] = ""
		return "", false
	}

	// first add
	if c.shortcuts == nil {
		c.shortcuts = map[string]string{}
	}

	// ensure it is one char
	short = string(short[0])

	if n, ok := c.shortcuts[short]; ok {
		color.Tips("error").Printf("The shortcut name '%s' has been used by option '%s'", short, n)
		Exit(-2)
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

// getShortName get a shortcut name by option name
func (c *Command) getShortName(name string) string {
	for s, n := range c.shortcuts {
		if n == name {
			return s
		}
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
	Name string
	Short string
	DType string // int, string, bool, value

	Required bool
	DefValue string
	Description string
}
