package cmd

import (
	"os"
	"runtime"

	"github.com/gookit/gcli/v2"
	"github.com/gookit/gcli/v2/show"
)

// options for the command
var eiOpts = struct {
	id    int
	c     string
	dir   string
	opt   string
	names Names
}{}

// EnvInfoCommand
func EnvInfoCommand() *gcli.Command {
	cmd := gcli.Command{
		Name:    "module:env",
		Aliases: []string{"env-info", "ei"},
		Desc:    "collect project info by git info",

		Func: envInfoRun,
	}

	cmd.IntOpt(&eiOpts.id, "id", "", 0, "the id option")
	cmd.StrOpt(&eiOpts.c, "c", "", "", "the config option")
	cmd.StrOpt(&eiOpts.dir, "dir", "d", "", "the dir option")

	return &cmd
}

// do run
func envInfoRun(c *gcli.Command, _ []string) error {
	eAble, _ := os.Executable()

	data := map[string]interface{}{
		"os":       runtime.GOOS,
		"binName":  c.BinName(),
		"workDir":  c.WorkDir(),
		"rawArgs":  os.Args,
		"execAble": eAble,
		"env":      os.Environ(),
	}

	show.JSON(&data)
	return nil
}
