package cmd

import (
	"github.com/gookit/cliapp"
	"github.com/gookit/cliapp/color"
	"fmt"
)

var colorOpts = struct {
	id  int
	c   string
	dir string
}{}

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
		Examples: "{$binName} {$cmd} --id 12 -c val ag0 ag1",
	}

	cmd.IntOpt(&colorOpts.id, "id", "", 2, "the id option")
	cmd.StrOpt(&colorOpts.c, "c", "", "value", "the config option")
	cmd.StrOpt(&colorOpts.dir, "dir", "", "", "the dir option")

	return &cmd
}

func colorUsage(cmd *cliapp.Command, args []string) int {
	// simple usage
	color.FgCyan.Printf("Simple to use %s\n", "color")

	// custom color
	color.New(color.FgWhite, color.BgBlack).Println("custom color style")
	// can also:
	color.Style{color.FgCyan, color.OpBold}.Println("custom color style")

	// use defined color tag
	color.Print("<suc>he</><comment>llo</>, <cyan>wel</><red>come</>\n")

	// use custom color tag
	color.Print("<fg=yellow;bg=black;op=underscore;>hello, welcome</>\n")

	// set a color tag
	color.Tag("info").Println("info style message")

	// tips
	color.Tips("info").Print("tips style message")
	color.Tips("warn").Print("tips style message")

	// lite tips
	color.LiteTips("info").Print("lite tips style message")
	color.LiteTips("warn").Print("lite tips style message")

	i := 0

	fmt.Print("\n- All Available Tags: \n\n")

	for tag, _ := range color.TagColors {
		i++
		color.Tag(tag).Print(tag)

		if i%5 == 0 {
			fmt.Print("\n")
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Print("\n")

	return 0
}

func byte8color()  {

}
