package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
)

// Config struct
type Config struct {
	// will exec command. eg: "go run main.go"
	Cmd string `json:"cmd"`
	// serve name for will exec command
	Name string `json:"name"`
	// run in the background
	Daemon bool
	// the pid file. eg "/var/run/serve.pid"
	PidFile string `json:"pidFile"`
	// the command run dir.
	WorkDir string `json:"workDir"`
}

var (
	config = new(Config)
	// config file
	confFile string
)

// ServerStart command
var ServerStart = &gcli.Command{
	Name: "start",
	Desc: "start server",
	Func: func(c *gcli.Command, args []string) error {
		return startServer(c.BinName())
	},
	Config: func(c *gcli.Command) {
		// c.StrOpt(&config.Pid, "pid", "", "", "the running server PID file")
		c.StrOpt(&confFile, "config", "c", "serve-config.json", "the running json config file path")
		c.BoolOpt(&config.Daemon, "daemon", "d", false, "the running server PID file")
	},
}

func startServer(binFile string) (err error) {
	if config.Daemon {
		cmd := exec.Command(binFile, "start")
		if err = cmd.Start(); err != nil {
			return
		}

		pid := cmd.Process.Pid
		color.Greenf("Server start, [PID] %d running...\n", pid)
		err = ioutil.WriteFile(config.PidFile, []byte(fmt.Sprintf("%d", pid)), 0666)
		config.Daemon = false
		return
	}

	color.Infoln("Server started")
	// front run
	// startHttp()
	return
}

var ServerStop = &gcli.Command{
	Name: "stop",
	Desc: "stop the running server(by PID file)",
	Func: func(_ *gcli.Command, _ []string) error {
		return stopServer()
	},
}

// ServerRestart Server restart
var ServerRestart = &gcli.Command{
	Name: "restart",
	Desc: "restart the running server by PID file",
	Func: func(c *gcli.Command, _ []string) (err error) {
		// c.App().SubRun("stop", []string{"-c", confFile})
		if err = stopServer(); err != nil {
			return
		}

		return startServer(c.BinName())
	},
}

func stopServer() error {
	bs, _ := ioutil.ReadFile(config.PidFile)
	cmd := exec.Command("kill", string(bs))
	err := cmd.Start()

	color.Success.Println("server stopped")
	return err
}
