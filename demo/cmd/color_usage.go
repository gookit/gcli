package cmd

import (
	"github.com/golangkit/cliapp"
	"github.com/golangkit/cliapp/color"
	"fmt"
)

var colorOpts ColorOpts
type ColorOpts struct {
	id  int
	c   string
	dir string
}

// ColorCommand command definition
func ColorCommand() *cliapp.Command {
	cmd := cliapp.Command{
		Name:        "color",
		Description: "this is a example for cli color usage",
		Aliases:     []string{"clr", "colors"},
		Fn:          colorUsage,
		ArgList: map[string]string{
			"arg0": "the first argument",
			"arg1": "the second argument",
		},
		Examples: "{{script}} {{cmd}} --id 12 -c val ag0 ag1",
	}

	colorOpts = ColorOpts{}

	f := &cmd.Flags
	f.IntVar(&colorOpts.id, "id", 2, "the id option")
	f.StringVar(&colorOpts.c, "c", "value", "the short option")
	f.StringVar(&colorOpts.dir, "dir", "", "the dir option")

	return &cmd
}

func colorUsage(cmd *cliapp.Command, args []string) int  {
	str := color.UseStyle("red", "text message")

	fmt.Println(str)

	return 0
}
