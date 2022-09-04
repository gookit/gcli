package cmd

import (
	"os"
	"runtime"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show"
)

// options for the command
var eiOpts = struct {
	id    int
	c     string
	dir   string
	opt   string
	names Names
}{}

// EnvInfo command
var EnvInfo = &gcli.Command{
	Name:    "env",
	Desc:    "collect project info by git info",
	Aliases: []string{"env-info", "ei"},
	Config: func(c *gcli.Command) {
		c.IntOpt(&eiOpts.id, "id", "", 0, "the id option")
		c.StrOpt(&eiOpts.c, "c", "", "", "the config option")
		c.StrOpt(&eiOpts.dir, "dir", "d", "", "the dir option")

	},

	Func: func(c *gcli.Command, _ []string) error {
		eAble, _ := os.Executable()

		data := map[string]any{
			"os":       runtime.GOOS,
			"binName":  c.BinName(),
			"workDir":  c.WorkDir(),
			"rawArgs":  os.Args,
			"execAble": eAble,
			"env":      os.Environ(),
		}

		show.JSON(&data)
		return nil
	},
}
