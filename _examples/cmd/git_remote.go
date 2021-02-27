package cmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/dump"
)

// GitRemote remote command of the git.
var GitRemote = &gcli.Command{
	Name: "remote",
	Desc: "remote command of the git",
	Aliases: []string{"rmt"},
	Config: func(c *gcli.Command) {
		c.BoolOpt()
	},
	Func: func(c *gcli.Command, args []string) error {
		dump.P(c.Path())
		return nil
	},
	Subs: []*gcli.Command{
		{
			Name: "set-url",
			Desc: "set-url command of git remote",
			Aliases: []string{"su"},
			Func: func(c *gcli.Command, args []string) error {
				dump.P(c.Path())
				return nil
			},
		},
	},
}
