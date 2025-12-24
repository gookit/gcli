package main

import (
	"fmt"

	"github.com/gookit/gcli/v3"
)

var run = &gcli.Command{
	Name: "run",
	Desc: "run app client",
	Func: func(c *gcli.Command, args []string) error {
		// server.Run()
		fmt.Println("test default command")
		return nil
	},
	// Hidden: true,
}

func Run() {
	app := gcli.NewApp()
	app.Version = "1.0.0"
	app.Name = "Test Client"
	app.Desc = ""
	app.Add(run)
	app.SetDefaultCommand(run.Name)
	app.Run(nil)
}

// RUN: go run ./_examples/issues/issue130.go
func main() {
	Run()
}
