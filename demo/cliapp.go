package main

import (
	"runtime"
	"github.com/golangkit/cliapp"
	"github.com/golangkit/cliapp/demo/cmd"
	"fmt"
)

// for test run: go build ./demo/cliapp.go && ./cliapp
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	str := `abc <err>err-text</> 
def <info>info text
</>`

	s := cliapp.ReplaceTag(str)

	fmt.Printf("%s\n", s)
	return

	app := cliapp.NewApp()
	app.Version = "1.0.3"
	app.Verbose = cliapp.VerbDebug
	app.Description = "this is my cli application"

	app.Add(cmd.ExampleCommand())
	app.Add(cmd.GitCommand())
	app.Add(&cliapp.Command{
		Name:        "demo",
		Aliases:     []string{"dm"},
		Description: "this is a description message for demo",
		Execute: func(cmd *cliapp.Command, args []string) int {
			cliapp.Stdout("hello, in the demo command\n")
			return 0
		},
	})

	app.Run()
}
