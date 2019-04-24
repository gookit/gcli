package cmd

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"time"

	"github.com/gookit/color"
	"github.com/gookit/gcli"
)

var bgrOpts = struct {
	deamon bool
}{}

func DaemonRunCommand() *gcli.Command {
	c := &gcli.Command{
		Name:   "bg:run",
		UseFor: "an example for background run program",
		Func:   handleDaemonRun,
	}

	c.BoolOpt(&bgrOpts.deamon, "daemon", "d", false, "want background run")

	return c
}

func handleDaemonRun(c *gcli.Command, _ []string) (err error) {
	if bgrOpts.deamon {
		newArgs := clearDaemonOpt("--daemon", "-d", c.OsArgs()[1:])
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
		gcli.Exit(0)
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
