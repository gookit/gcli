package cmd

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/dump"
)

var gitRmtOPts = struct {
	v bool
}{}

// GitRemote remote command of the git.
var GitRemote = &gcli.Command{
	Name: "remote",
	Desc: "remote command of the git",
	Aliases: []string{"rmt"},
	Config: func(c *gcli.Command) {
		c.BoolOpt(&gitRmtOPts.v, "v", "", false, "option for git remote")
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
			Config: func(c *gcli.Command) {
				c.AddArg("name", "the remote name", true)
				c.AddArg("address", "the remote address", true)
			},
			Func: func(c *gcli.Command, args []string) error {
				dump.P(c.Path())
				return nil
			},
		},
	},
}
