package cmd

import (
    "feedscenter/console/cli"
    "fmt"
)

var testCmd = cli.Command{
    Name: "test",
    Description: "this is a description message",
    Aliases: []string{"ts","test1"},
}

func TestCommand() *cli.Command {
    testCmd.Execute = testExecute

    return &testCmd
}

func testExecute(cmd *cli.Command, args []string) int  {
	// latest commit id by: git log --pretty="%H" -n1 HEAD

	// latest commit date by: git log -n1 --pretty=%ci HEAD


	// get tag: git describe --tags --exact-match HEAD
	// get branch: git branch -a | grep "*"

    fmt.Print("hello, in test command\n")
    return 0
}
