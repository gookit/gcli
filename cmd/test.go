package cmd

import (
    "feedscenter/console/cli"
    "fmt"
)

var gitCmd = cli.Command{
    Name: "git",
    Description: "this is a description message",
    Aliases: []string{"git-info"},
}

func GitCommand() *cli.Command {
    gitCmd.Execute = gitExecute

    return &gitCmd
}

func gitExecute(cmd *cli.Command, args []string) int  {
    fmt.Print("hello, in test command\n")
    return 0
}
