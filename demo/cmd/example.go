package cmd

import (
    cli "github.com/golangkit/cliapp"
	"fmt"
)

var exampleOpts ExampleOpts
type ExampleOpts struct {
	id  int
	c   string
	dir string
}

// ExampleCommand command definition
func ExampleCommand() *cli.Command {
	cmd := cli.Command{
		Name:        "example",
		Description: "this is a description message",
		Aliases:     []string{"exp", "ex"},
		Fn:          exampleExecute,
		ArgList: map[string]string{
			"arg0": "the first argument",
			"arg1": "the second argument",
		},
		Examples: "{$script} {$cmd} --id 12 -c val ag0 ag1",
	}

	exampleOpts = ExampleOpts{}

	f := &cmd.Flags
	f.IntVar(&exampleOpts.id, "id", 2, "the id option")
	f.StringVar(&exampleOpts.c, "c", "value", "the short option")
	f.StringVar(&exampleOpts.dir, "dir", "", "the dir option")

	return &cmd
}

// command running
// example run:
// 	go build cliapp.go && ./cliapp example --id 12 -c val ag0 ag1
func exampleExecute(cmd *cli.Command, args []string) int {
	fmt.Print("hello, in example command\n")

	// fmt.Printf("%+v\n", cmd.Flags)
	fmt.Printf("opts %+v\n", exampleOpts)
	fmt.Printf("args is %v\n", args)

	return 0
}
