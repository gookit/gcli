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
    fmt.Print("hello, in test command\n")
    return 0
}
