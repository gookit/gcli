package cmd

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show/emoji"
)

func EmojiDemoCmd() *gcli.Command {
	return &gcli.Command{
		Name:    "emoji",
		Desc:    "this is a emoji usage example command",
		Aliases: []string{"emoj"},
		Config: func(c *gcli.Command) {
			c.AddArg("subcmd", "The name of the subcommand you want to run. allow: render, search", true)
			c.AddArg("param", "Used in the previous subcommand. It's message string OR keywords for search", true)
		},
		Func: func(c *gcli.Command, _ []string) error {
			subCmd := c.Arg("subcmd").String()
			param := c.Arg("param").String()
			switch subCmd {
			case "render":
				return renderEmoji(param)
			case "search":
				return searchEmoji(param)
			default:
				return c.Errorf("invalid sub-command name for %s, only allow: render, search", c.Name)
			}
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
