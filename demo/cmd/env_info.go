package cmd


import (
	cli "github.com/gookit/cliapp"
	"fmt"
	"runtime"
	"os"
	"github.com/gookit/cliapp/utils"
)

// options for the command
var eiOpts = defEiOpts{}

type defEiOpts struct {
	id    int
	c     string
	dir   string
	opt   string
	names Names
}

// EnvInfoCommand
func EnvInfoCommand() *cli.Command {
	cmd := cli.Command{
		Name:        "env",
		Aliases:     []string{"env-info", "ei"},
		Description: "collect project info by git info",

		Fn: envInfoRun,
	}

	cmd.IntOpt(&eiOpts.id, "id", "", 0, "the id option")
	cmd.StrOpt(&eiOpts.c, "c", "", "", "the config option")
	cmd.StrOpt(&eiOpts.dir, "dir", "d", "", "the dir option")

	return &cmd
}

// do run
func envInfoRun(cmd *cli.Command, args []string) int   {
	eAble,_ := os.Executable()

	data := map[string]interface{}{
		"os": runtime.GOOS,
		"binName": cli.BinName(),
		"workDir": cli.WorkDir(),
		"rawArgs": os.Args,
		"execAble": eAble,
		"env": os.Environ(),
	}

	str, _ := utils.PrettyJson(&data)

	fmt.Println(str)
	return 0
}
