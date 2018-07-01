package cliapp

import (
	"flag"
	"github.com/gookit/color"
)

// IntOpt set a int option
func (c *Command) IntOpt(p *int, name string, short string, defaultValue int, description string) *Command {
	c.Flags.IntVar(p, name, defaultValue, description)

	if len(short) == 1 {
		c.addShortcut(name, short)
		c.Flags.IntVar(p, short, defaultValue, "")
	}

	return c
}

// UintOpt set a int option
func (c *Command) UintOpt(p *uint, name string, short string, defaultValue uint, description string) *Command {
	c.Flags.UintVar(p, name, defaultValue, description)

	if len(short) == 1 {
		c.addShortcut(name, short)
		c.Flags.UintVar(p, short, defaultValue, "")
	}

	return c
}

// StrOpt set a str option
func (c *Command) StrOpt(p *string, name string, short string, defaultValue string, description string) *Command {
	c.Flags.StringVar(p, name, defaultValue, description)

	if len(short) == 1 {
		c.addShortcut(name, short)
		c.Flags.StringVar(p, short, defaultValue, "")
	}

	return c
}

// BoolOpt set a bool option
func (c *Command) BoolOpt(p *bool, name string, short string, defaultValue bool, description string) *Command {
	c.Flags.BoolVar(p, name, defaultValue, description)

	if len(short) == 1 {
		c.addShortcut(name, short)
		c.Flags.BoolVar(p, short, defaultValue, "")
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

	if len(short) == 1 {
		c.addShortcut(name, short)
		c.Flags.Var(p, short, "")
	}

	return c
}

// addShortcut add a shortcut name for a option name
func (c *Command) addShortcut(name string, short string) {
	if n, ok := c.shortcuts[short]; ok {
		color.Tips("error").Printf("The shortcut name '%s' has been used by option '%s'", short, n)
		Exit(-2)
	}

	// first add
	if c.shortcuts == nil {
		c.shortcuts = map[string]string{short: name}
		return
	}

	c.shortcuts[short] = name
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
