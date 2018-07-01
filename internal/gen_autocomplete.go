package internal

import (
	"github.com/gookit/cliapp"
	"fmt"
)

//
var genOpts = struct {
	shell   string
	binName string
	output string
}{}

func GenShAutoComplete() *cliapp.Command {
	cmd := cliapp.Command{
		Name:    "gen",
		Aliases: []string{"gen-ac"},

		Description: "generate script file for command auto complete",

		Fn: doGen,
	}

	cmd.StrOpt(
		&genOpts.shell,
		"shell",
		"",
		"bash",
		"the shell env name for want generated, allow: sh,zsh,bash",
	).StrOpt(
		&genOpts.binName,
		"bin-name",
		"b",
		"",
		"your packaged application bin file name.",
	).StrOpt(
		&genOpts.output,
		"output",
		"o",
		"auto-completion.{shell}",
		"output shell auto completion script file name.",
	)

	return &cmd
}

func doGen(cmd *cliapp.Command, args []string) int {
	fmt.Println(genOpts)

	return 0
}
