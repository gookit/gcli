package builtin

import "github.com/gookit/cliapp"

type genEmojiMap struct {
	// muan - https://github.com/muan/emoji/blob/gh-pages/javascripts/emojilib/emojis.json
	// gemoji - https://raw.githubusercontent.com/github/gemoji/master/db/emoji.json
	// unicode - https://unicode.org/Public/emoji/11.0/emoji-test.txt
	source  string // allow: gemoji
	saveDir string
}

type Gemoji struct {
	Aliases     []string `json:"aliases"`
	Description string   `json:"description"`
	Emoji       string   `json:"emoji"`
	Tags        []string `json:"tags"`
}

// GenEmojiMapCommand create
func GenEmojiMapCommand() *cliapp.Command {
	gem := &genEmojiMap{}

	return &cliapp.Command{
		Name:    "gen:emojis",
		Aliases: []string{"gen:emoji"},
		// handler func
		Func: gem.run,
		// des
		Description: "fetch all emoji codes form https://www.unicode.org, generate a go file.",
		Init: func(c *cliapp.Command) {
			c.StrOpt(&gem.source, "source", "s", "gemoji",
				"the emoji data source, allow: muan, gemoji, unicode")
			c.StrOpt(&gem.saveDir, "dir", "d", "./", "the generated go file save `DIR` path")
		},
	}
}

func (g *genEmojiMap) run(c *cliapp.Command, _ []string) int {

	return 0
}

const templateString = `
package {{.PkgName}}
// NOTE: THIS FILE WAS PRODUCED BY THE
// EMOJI COD EMAP CODE GENERATION TOOL (https://github.com/gookit/cliapp)
// DO NOT EDIT
// Mapping from character to concrete escape code.
var emojiMap = map[string]string{
	{{range $key, $val := .CodeMap}}":{{$key}}:": {{$val}},
{{end}}
}
`
