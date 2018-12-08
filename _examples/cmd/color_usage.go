package cmd

import (
	"fmt"
	"github.com/gookit/gcli"
	"github.com/gookit/color"
)

var colorOpts = struct {
	id  int
	c   string
	dir string
}{}

// ColorCommand command definition
func ColorCommand() *gcli.Command {
	cmd := gcli.Command{
		Name:     "color",
		UseFor:   "this is a example for cli color usage",
		Aliases:  []string{"clr", "colors"},
		Func:     colorUsage,
		Examples: "{$binName} {$cmd} --id 12 -c val ag0 ag1",
	}

	cmd.IntOpt(&colorOpts.id, "id", "", 2, "the id option")
	cmd.StrOpt(&colorOpts.c, "c", "", "value", "the config option")
	cmd.StrOpt(&colorOpts.dir, "dir", "", "", "the dir option")

	return &cmd
}

func colorUsage(cmd *gcli.Command, args []string) int {
	// simple usage
	color.FgCyan.Printf("Simple to use %s\n", "color")

	// custom color
	color.New(color.FgMagenta, color.BgBlack).Println("custom color style")
	// can also:
	color.Style{color.FgCyan, color.OpBold}.Println("custom color style")

	// use defined color tag
	color.Print("<suc>he</><comment>llo</>, <cyan>wel</><red>come</>\n")

	// use custom color tag
	color.Print("<fg=yellow;bg=black;op=underscore;>hello, welcome</>\n")

	// set a color tag
	color.Tag("info").Println("info style message")

	// prompt message
	color.Info.Prompt("prompt style message")
	color.Warn.Prompt("prompt style message")

	// tips message
	color.Info.Tips("tips style message")
	color.Warn.Tips("tips style message")

	i := 0
	fmt.Print("\n- All Available color Tags: \n\n")

	for tag, _ := range color.GetColorTags() {
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

func byte8color() {

}
