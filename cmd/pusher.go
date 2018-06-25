package cmd

import "feedscenter/consumer/cli"

// moment pusher

var pusher = cli.Command{
	Name:        "pusher",
	Description: "push moment to user repo",
}

func NewPusher() *cli.Command {
	pusher.Execute = executePush

	return &pusher
}

func executePush(cmd *cli.Command, args []string) int {

	return 0
}
