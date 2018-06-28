package cmd

import (
	"github.com/golangkit/cliapp"
	"github.com/golangkit/cliapp/color"
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
		Examples: "{$script} {$cmd} --id 12 -c val ag0 ag1",
	}

	colorOpts = ColorOpts{}

	f := &cmd.Flags
	f.IntVar(&colorOpts.id, "id", 2, "the id option")
	f.StringVar(&colorOpts.c, "c", "value", "the short option")
	f.StringVar(&colorOpts.dir, "dir", "", "the dir option")

	return &cmd
}

func colorUsage(cmd *cliapp.Command, args []string) int {
	// simple usage
	color.FgCyan.Printf("Simple to use %s\n", "color")

	// custom color
	color.New(color.FgWhite, color.BgBlack).Println("custom color style")

	// use defined color tag
	color.Print("<suc>he</><comment>llo</>, <cyan>wel</><red>come</>\n")

	// use custom color tag
	color.Print("<fg=yellow;bg=black;op=underscore;>hello, welcome</>\n")

	// set a color tag
	color.Tag("info").Print("info style\n")

	// use info style tips
	color.Tips("info").Print("tips style")

	// use info style blocked tips
	color.BlockTips("info").Print("blocked tips style")

	return 0
}
