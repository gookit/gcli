package cmd

import (
	"fmt"
	"github.com/gookit/cliapp"
	"github.com/gookit/cliapp/show/emoji"
	"github.com/gookit/color"
)

func EmojiDemoCmd() *cliapp.Command {
	return &cliapp.Command{
		Name: "emoji",
		UseFor: "this is a emoji usage example command",
		Aliases: []string{"emoj"},
		Config: func(c *cliapp.Command) {
			c.AddArg("subcmd", "The name of the subcommand you want to run. allow: render, search", true)
			c.AddArg("param", "Used in the previous subcommand. It's message string OR keywords for search", true)
		},
		Func: func(c *cliapp.Command, _ []string) int {
			subCmd := c.Arg("subcmd").String()
			param := c.Arg("param").String()
			switch subCmd {
			case "render":
				renderEmoji(param)
			case "search":
				searchEmoji(param)
			default:
				return c.Errorf("invalid sub-command name for %s, only allow: render, search", c.Name)
			}

			return 0
		},
		Examples: `An render example
  {$fullCmd} render ":car: a message text, contains emoji :smile:"
An search example
  {$fullCmd} search smi`,
	}
}

func renderEmoji(msg string) (err error) {
	fmt.Println(emoji.Render(msg))
	return
}

func searchEmoji(kw string) (err error) {
	ret := emoji.Search(kw, 15)
	if len(ret) == 0 {
		color.Info.Tips(":( no matched emoji found! keyword: %s", kw)
	}

	color.Success.Println("OK, successfully found some emoji:")
	for name, code := range ret {
		fmt.Println(code, name)
	}

	return
}
