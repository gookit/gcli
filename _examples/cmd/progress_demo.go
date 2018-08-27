package cmd

import "github.com/gookit/cliapp"

type ProgressDemo struct {
}

func ProgressDemoCommand() *cliapp.Command {
	return &cliapp.Command{
		Name: "prog:demo",
		UseFor: "there are some progress bar run demos",
		Aliases: []string{"prg:demo", "progress"},
	}
}

func (d *ProgressDemo) Run(c *cliapp.Command, _ []string) int {
	return 0
}
