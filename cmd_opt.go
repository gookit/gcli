package cliapp

import "flag"

// IntOpt set a int option
func (c *Command) IntOpt(p *int, name string, short string, defaultValue int, description string) *Command {
	c.Flags.IntVar(p, name, defaultValue, description)

	if len(short) == 1 {
		c.Flags.IntVar(p, short, defaultValue, "")
	}

	return c
}

// UintOpt set a int option
func (c *Command) UintOpt(p *uint, name string, short string, defaultValue uint, description string) *Command {
	c.Flags.UintVar(p, name, defaultValue, description)

	if len(short) == 1 {
		c.Flags.UintVar(p, short, defaultValue, "")
	}

	return c
}

// StrOpt set a str option
func (c *Command) StrOpt(p *string, name string, short string, defaultValue string, description string) *Command {
	c.Flags.StringVar(p, name, defaultValue, description)

	if len(short) == 1 {
		c.Flags.StringVar(p, short, defaultValue, "")
	}

	return c
}

// BoolOpt set a bool option
func (c *Command) BoolOpt(p *bool, name string, short string, defaultValue bool, description string) *Command {
	c.Flags.BoolVar(p, name, defaultValue, description)

	if len(short) == 1 {
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
		c.Flags.Var(p, short, "")
	}

	return c
}
