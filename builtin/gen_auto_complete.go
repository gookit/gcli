package builtin

import (
	"os"
	"strings"

	"github.com/gookit/cliui/interact"
	"github.com/gookit/cliui/show"
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/cliutil"
)

// current supported shell
const (
	ZshShell  = gcli.ZshShell
	BashShell = gcli.BashShell
)

// generate options
var genOpts = &struct {
	shell   string
	binName string
	output  string
	static  bool
}{}

// GenAutoComplete create command
func GenAutoComplete(fns ...func(c *gcli.Command)) *gcli.Command {
	c := &gcli.Command{
		Func:    doGen,
		Name:    "genac",
		Aliases: []string{"gen-ac"},
		Desc:    "generate auto complete scripts for current application",
	}

	shell := cliutil.CurrentShell(true)
	if shell == "" {
		shell = "bash"
	}

	c.StrOpt(
		&genOpts.shell,
		"shell",
		"s",
		shell,
		"the shell env name for want generated, allow: zsh,bash,pwsh(pwsh only for dynamic)",
	)
	c.StrOpt(
		&genOpts.binName,
		"bin-name",
		"b",
		"",
		"your packaged application bin file name.",
	)
	c.StrOpt(
		&genOpts.output,
		"output",
		"o",
		"auto-completion."+shell,
		"output shell auto completion script file name.",
	)
	c.BoolOpt(
		&genOpts.static,
		"static",
		"S",
		false,
		"generate static(embedded) completion script instead of the default dynamic(thin) one.",
	)

	for _, fn := range fns {
		fn(c)
	}
	return c
}

func doGen(c *gcli.Command, _ []string) (err error) {
	if len(genOpts.binName) == 0 {
		genOpts.binName = c.Ctx.BinName()
	}

	genOpts.binName = strings.TrimSuffix(strings.Trim(genOpts.binName, "./"), ".exe")

	if len(genOpts.output) == 0 {
		genOpts.output = genOpts.binName + "." + genOpts.shell
	} else {
		genOpts.output = strings.Replace(genOpts.output, "{shell}", genOpts.shell, 1)

		// check suffix
		if !strings.Contains(genOpts.output, ".") {
			genOpts.output += "." + genOpts.shell
		}
	}

	// color.Info.Tips("\n  %+v\n", genOpts)
	data := map[string]any{
		"Shell":    genOpts.shell,
		"BinName":  genOpts.binName,
		"FileName": genOpts.output,
	}

	show.AList("Information", data)

	if interact.Unconfirmed("Please confirm the above information", true) {
		color.Info.Print("\nBye :)\n")
		return
	}

	// 默认生成动态(瘦)脚本; --static/-S 时生成静态(嵌入式)脚本。
	// 传入 genOpts.binName 以保留 --bin-name 定制能力。
	var str string
	if genOpts.static {
		str, err = c.App().GenStaticCompletionScript(genOpts.shell, genOpts.binName)
	} else {
		str, err = c.App().GenCompletionScript(genOpts.shell, genOpts.binName)
	}
	if err != nil {
		return c.NewErrf("%s", err.Error())
	}

	color.Infoln("Now, will write content to file ", genOpts.output)
	color.Normal.Print("Continue?")

	if !interact.AnswerIsYes(true) {
		color.Info.Print("\nBye :)\n")
		return
	}

	// Open the file for reading and writing, if it does not exist, create it
	err = os.WriteFile(genOpts.output, []byte(str), 0664)
	if err != nil {
		return c.NewErrf("Write file error: %s", err.Error())
	}

	color.Success.Println("\nOK, auto-complete file generate successful")
	return
}
