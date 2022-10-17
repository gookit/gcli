package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
)

var bgrOpts = struct {
	deamon bool
}{}

var DaemonRun = &gcli.Command{
	Name:    "bgrun",
	Desc:    "an example for background run program",
	Func:    handleDaemonRun,
	Aliases: []string{"bgr"},
	Config: func(c *gcli.Command) {
		c.BoolOpt(&bgrOpts.deamon, "daemon", "d", false, "want background run")

	},
}

func handleDaemonRun(c *gcli.Command, _ []string) (err error) {
	if bgrOpts.deamon {
		newArgs := clearDaemonOpt("--daemon", "-d", c.Ctx.OsArgs()[1:])
		newCmd := exec.Command(c.BinName(), newArgs...)

		if err = newCmd.Start(); err != nil {
			return
		}

		pid := newCmd.Process.Pid
		color.Magenta.Printf("server start, process [PID:%d] running...\n", pid)

		err = ioutil.WriteFile("./server.pid", []byte(fmt.Sprintf("%d", pid)), 0666)
		if err != nil {
			return
		}

		bgrOpts.deamon = false
		os.Exit(0)
	}

	// block process
	for {
		fmt.Println(time.Now())
		time.Sleep(time.Second * 2)
	}
}

// newArgs := clearDaemonOpt("--daemon", "-d", c.OsArgs()[1:])
func clearDaemonOpt(name, short string, args []string) []string {
	var newArgs []string
	for _, val := range args {
		if val == name || val == short {
			continue
		}

		newArgs = append(newArgs, val)
	}

	return newArgs
}
