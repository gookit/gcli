package cliapp

import (
	"log"
	"os"
)

// AloneRun
func (c *Command) AloneRun() int {
	// mark is alone
	c.alone = true

	// init
	c.Init()

	// set help handler
	c.Flags.Usage = func() {
		// init some tpl vars
		c.Vars = map[string]string{
			"workDir": workDir,
			"binName": binName,
		}

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
