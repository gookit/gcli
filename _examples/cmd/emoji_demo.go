package cmd

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/show/emoji"
)

var EmojiDemo = &gcli.Command{
	Name:    "emoji",
	Desc:    "this is a emoji usage example command",
	Aliases: []string{"emoj"},
	// Func: ,
	Examples: `
An render example
  {$fullCmd} render ":car: a message text, contains emoji :smile:"
An search example
  {$fullCmd} search smi`,
  	Subs: []*gcli.Command{
		{
			Name: "render",
			Desc: "render given string, will replace special char to emoji",
			Aliases: []string{"r"},
			Config: func(c *gcli.Command) {
				c.AddArg("msg", "The message string for render", true)
			},
			Func: func(c *gcli.Command, args []string) error {
				fmt.Println(emoji.Render(c.Arg("msg").String()))
				return nil
			},
		},
		{
			Name: "search",
			Desc: "search emojis by given keywords",
			Aliases: []string{"s"},
			Config: func(c *gcli.Command) {
				c.AddArg("keyword", "The keyword string for search", true)
			},
			Func: func(c *gcli.Command, args []string) error {
				kw := c.Arg("keyword").String()

				return searchEmoji(kw)
			},
		},
	},
}

func searchEmoji(kw string) (err error) {
	ret := emoji.Search(kw, 15)
	if len(ret) == 0 {
		color.Note.Tips(":( no matched emoji found! keyword: %s", kw)
		return
	}

	color.Success.Println("OK, successfully found some emojis:")
	for name, code := range ret {
		fmt.Println(code, name)
	}
	return
}
