package builtin

import (
	"fmt"
	"github.com/gookit/gcli"
	"io"
	"net/http"
	"os"
	"time"
)

type genEmojiMap struct {
	c *gcli.Command
	// muan - https://raw.githubusercontent.com/muan/emoji/gh-pages/javascripts/emojilib/emojis.json
	// gemoji - https://raw.githubusercontent.com/github/gemoji/master/db/emoji.json
	// unicode - https://unicode.org/Public/emoji/11.0/emoji-test.txt
	source  string // allow: gemoji
	saveDir string
	onlyGen bool
}

// Gemoji definition
type Gemoji struct {
	Aliases     []string `json:"aliases"`
	Description string   `json:"description"`
	Emoji       string   `json:"emoji"`
	Tags        []string `json:"tags"`
}

// GenEmojiMapCommand create
func GenEmojiMapCommand() *gcli.Command {
	gem := &genEmojiMap{}

	return &gcli.Command{
		Name:    "gen:emojis",
		Aliases: []string{"gen:emoji", "gen:emj"},
		// handler func
		Func: gem.run,
		// des
		UseFor: "fetch emoji codes form data source url, then generate a go file.",
		// config options
		Config: func(c *gcli.Command) {
			gem.c = c
			c.StrOpt(
				&gem.source, "source", "s", "gemoji",
				"the emoji data source, allow: muan, gemoji, unicode",
			)
			c.StrOpt(&gem.saveDir, "dir", "d", "./", "the generated go file save `DIR` path")
			c.BoolOpt(&gem.onlyGen, "onlyGen", "", false, "whether only generate go file from exists emoji data file")
		},
		Help: `source allow:
 muan - https://raw.githubusercontent.com/muan/emoji/gh-pages/javascripts/emojilib/emojis.json
 gemoji - https://raw.githubusercontent.com/github/gemoji/master/db/emoji.json
 unicode - https://unicode.org/Public/emoji/11.0/emoji-test.txt
`,
	}
}

func (g *genEmojiMap) run(c *gcli.Command, _ []string) int {

	return 0
}

// Download 实现单个文件的下载
func (g *genEmojiMap) Download(remoteFile string, saveAs string) error {
	nt := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s]To download %s\n", nt, remoteFile)

	newFile, err := os.Create(saveAs)
	if err != nil {
		return err
	}
	defer newFile.Close()

	client := http.Client{Timeout: 900 * time.Second}
	resp, err := client.Get(remoteFile)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(newFile, resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

const templateString = `
package {{.PkgName}}
// NOTE: THIS FILE WAS PRODUCED BY THE
// EMOJI COD EMAP CODE GENERATION TOOL (https://github.com/gookit/gcli)
// DO NOT EDIT
// Mapping from character to concrete escape code.
var emojiMap = map[string]string{
	{{range $key, $val := .CodeMap}}":{{$key}}:": {{$val}},
{{end}}
}
`
