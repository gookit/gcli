package cmd

import (
	"feedscenter/console/cli"
	"fmt"
)

var testOpts TestOpts
type TestOpts struct {
	id  int
	c   string
	dir string
}

// TestCommand command definition
func TestCommand() *cli.Command {
	cmd := cli.Command{
		Name:        "test",
		Description: "this is a description message",
		Aliases:     []string{"ts"},
		Execute:     testExecute,
	}

	testOpts = TestOpts{}

	f := &cmd.Flags
	f.IntVar(&testOpts.id, "id", 2, "the id option")
	f.StringVar(&testOpts.c, "c", "value", "the short option")
	f.StringVar(&testOpts.dir, "dir", "", "the dir option")

	return &cmd
}

// command running
// test run:
// 	go build console/cliapp.go && ./cliapp test --id 12 -c val ag0 ag1
func testExecute(cmd *cli.Command, args []string) int {
	fmt.Print("hello, in test command\n")

	// fmt.Printf("%+v\n", cmd.Flags)
	fmt.Printf("opts %+v\n", testOpts)
	fmt.Printf("args is %v\n", args)

	return 0
}
