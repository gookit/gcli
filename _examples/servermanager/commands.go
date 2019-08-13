package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/gookit/color"
	"github.com/gookit/gcli"
)

type Config struct {
	// will exec command. eg. "go run main.go"
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

var config *Config

// config file
var confFile string

// eg: cliapp serve:start
func ServerStart() *gcli.Command {
	c := &gcli.Command{
		Name: "start",
	}

	// c.StrOpt(&config.Pid, "pid", "", "", "the running server PID file")
	c.StrOpt(&confFile, "config", "c", "serve-config.json", "the running json config file path")
	c.BoolOpt(&config.Daemon, "daemon", "d", false, "the running server PID file")

	return c
}

func startServer() int {
	if config.Daemon {
		command := exec.Command("gonne", "start")
		command.Start()
		fmt.Printf("server start, [PID] %d running...\n", command.Process.Pid)
		ioutil.WriteFile(config.PidFile, []byte(fmt.Sprintf("%d", command.Process.Pid)), 0666)
		config.Daemon = false
		return 0
	} else {
		fmt.Println("gonne start")
	}

	// front run
	// startHttp()

	return 0
}

func ServerStop() *gcli.Command {
	cmd := &gcli.Command{
		Name:   "stop",
		UseFor: "stop the running server by PID file",
	}

	cmd.Func = func(_ *gcli.Command, _ []string) int {
		return stopServer()
	}

	return cmd
}

func stopServer() int {
	bs, _ := ioutil.ReadFile(config.PidFile)
	command := exec.Command("kill", string(bs))
	command.Start()

	color.Success.Println("server stopped")
	return 0
}

func ServerRestart() *gcli.Command {
	cmd := &gcli.Command{
		Name:   "restart",
		UseFor: "restart the running server by PID file",
	}

	cmd.Func = func(c *gcli.Command, _ []string) int {
		// c.App().SubRun("stop", []string{"-c", confFile})
		stopServer()
		startServer()

		return 0
	}

	return cmd
}
