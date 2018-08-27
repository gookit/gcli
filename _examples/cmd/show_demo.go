package cmd

import "github.com/gookit/cliapp"

// ShowDemoCommand create
func ShowDemoCommand() *cliapp.Command {
	c := &cliapp.Command{
		Name: "show",
		Func: runShow,
		//
		UseFor: "the command will show some data format methods",
	}

	return c
}

func runShow(c *cliapp.Command, _ []string) int {

	return 0
}
