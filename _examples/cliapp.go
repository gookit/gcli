package main

import (
	"github.com/gookit/cliapp"
	"github.com/gookit/cliapp/builtin"
	"github.com/gookit/cliapp/_examples/cmd"
	"runtime"
	"github.com/gookit/cliapp/builtin/filewatcher"
)

// for test run: go run ./_examples/cliapp.go
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app := cliapp.NewApp()
	app.Version = "1.0.3"
	app.Description = "this is my cli application"

	app.SetVerbose(cliapp.VerbDebug)
	// app.DefaultCmd("exampl")

	app.Add(cmd.ExampleCommand())
	app.Add(cmd.EnvInfoCommand())
	app.Add(cmd.GitCommand())
	app.Add(cmd.ColorCommand())
	app.Add(filewatcher.FileWatcher())
	app.Add(&cliapp.Command{
		Name:        "test",
		Aliases:     []string{"ts"},
		Description: "this is a description <info>message</> for command {$cmd}",
		Fn: func(cmd *cliapp.Command, args []string) int {
			cliapp.Print("hello, in the test command\n")
			return 0
		},
	})

	app.Add(builtin.GenShAutoComplete())
	// fmt.Printf("%+v\n", cliapp.CommandNames())
	app.Run()
}
