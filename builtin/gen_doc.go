package builtin

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/docgen"
)

// gen doc options
var docOpts = &struct {
	format string
	output string
}{}

// GenDoc create a command to export application commands documentation(markdown/man).
//
// 参照 GenAutoComplete 范式: app 添加后即可 `./cliapp gendoc -f md -o ./docs` 导出文档。
func GenDoc(fns ...func(c *gcli.Command)) *gcli.Command {
	c := &gcli.Command{
		Func:    doGenDoc,
		Name:    "gendoc",
		Aliases: []string{"gen-doc"},
		Desc:    "generate documentation(markdown/man) for current application commands",
	}

	c.StrOpt(
		&docOpts.format,
		"format",
		"f",
		"md",
		"the documentation format for generated, allow: md, man",
	)
	c.StrOpt(
		&docOpts.output,
		"output",
		"o",
		"./docs",
		"the output directory for generated documentation files.",
	)

	for _, fn := range fns {
		fn(c)
	}
	return c
}

func doGenDoc(c *gcli.Command, _ []string) error {
	app := c.App()
	dir := docOpts.output

	// 按 format 选择渲染器; 非法格式直接返回 error。
	switch docOpts.format {
	case "md", "markdown":
		if err := docgen.MarkdownTree(app, dir); err != nil {
			return c.NewErrf("generate markdown docs error: %s", err.Error())
		}
	case "man":
		if err := docgen.ManTree(app, dir); err != nil {
			return c.NewErrf("generate man docs error: %s", err.Error())
		}
	default:
		return c.NewErrf("invalid format %q, allow: md, man", docOpts.format)
	}

	color.Success.Printf("\nOK, documentation(%s) generated to: %s\n", docOpts.format, dir)
	return nil
}
