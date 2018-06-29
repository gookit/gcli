package cliapp

import (
	"log"
	"os"
)

// AloneRun
func (c *Command) AloneRun() int {
	c.alone = true

	// init some tpl vars
	c.Vars = map[string]string{
		"workDir": workDir,
		"binName": binName,
	}

	c.Flags.Usage = func() {
		c.ShowHelp(true)
	}

	// don't display date on print log
	log.SetFlags(0)

	c.Flags.Parse(os.Args[1:])

	return c.Fn(c, c.Flags.Args())
}

// IsAlone
func (c *Command) IsAlone() bool {
	return c.alone
}

// NotAlone
func (c *Command) NotAlone() bool {
	return !c.alone
}
