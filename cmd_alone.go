package cliapp

import (
	"log"
	"os"
)

// AloneRun current command
func (c *Command) AloneRun() int {
	// don't display date on print log
	log.SetFlags(0)
	// mark is alone
	c.alone = true
	// args := parseGlobalOpts()
	// init
	c.Init()
	// parse args and opts
	c.Flags.Parse(os.Args[1:])

	return c.Execute(c.Flags.Args())
}

// IsAlone running
func (c *Command) IsAlone() bool {
	return c.alone
}

// NotAlone running
func (c *Command) NotAlone() bool {
	return !c.alone
}
