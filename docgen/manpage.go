package docgen

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/gcli/v3"
)

// escapeRoff 转义 man(roff) 文本中的特殊字符。
//   - `\` -> `\\`(避免被当作转义引导)
//   - `-`  -> `\-`(roff 中连字符需转义, 否则被渲染为短横/连字)
func escapeRoff(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "-", "\\-")
	return s
}

// roffLine 转义并输出一行文本。
// roff 中行首的 `.` 与 `'` 是控制字符, 需用 `\&` 零宽前缀保护。
func roffLine(s string) string {
	s = escapeRoff(cleanLine(s))
	if s != "" && (s[0] == '.' || s[0] == '\'') {
		s = "\\&" + s
	}
	return s
}

// CmdMan 渲染单个命令为 man page(roff 格式)。
func CmdMan(c *gcli.Command) string {
	var buf strings.Builder

	// .TH 头: 标题(命令名大写) section date source manual
	buf.WriteString(".TH \"" + strings.ToUpper(escapeRoff(c.Name)) + "\" \"1\" \"\" \"\" \"\"\n")

	// NAME
	buf.WriteString(".SH NAME\n")
	buf.WriteString(roffLine(c.Path()) + " \\- " + roffLine(renderText(c, c.Desc)) + "\n")

	// SYNOPSIS
	buf.WriteString(".SH SYNOPSIS\n")
	buf.WriteString(roffLine(c.Path()) + " [\\-\\-options ...] [arguments ...]\n")

	// DESCRIPTION(长帮助)
	if c.Help != "" {
		buf.WriteString(".SH DESCRIPTION\n")
		buf.WriteString(roffLine(renderText(c, c.Help)) + "\n")
	}

	// OPTIONS(跳过 Hidden)
	opts := c.Opts()
	if hasVisibleOpts(opts) {
		buf.WriteString(".SH OPTIONS\n")
		for _, name := range sortedOptNames(opts) {
			opt := opts[name]
			if opt.Hidden {
				continue
			}
			buf.WriteString(".TP\n")
			// 选项名转义, eg: --name -> \-\-name
			buf.WriteString("\\fB" + escapeRoff(optHelpName(opt)) + "\\fR\n")
			buf.WriteString(roffLine(renderText(c, opt.Desc)) + "\n")
		}
	}

	// ARGUMENTS
	args := c.Args()
	if len(args) > 0 {
		buf.WriteString(".SH ARGUMENTS\n")
		for _, arg := range args {
			argName := arg.Name
			if arg.Arrayed {
				argName += "..."
			}
			buf.WriteString(".TP\n")
			buf.WriteString("\\fB" + escapeRoff(argName) + "\\fR\n")
			buf.WriteString(roffLine(renderText(c, arg.Desc)) + "\n")
		}
	}

	// EXAMPLES
	if c.Examples != "" {
		buf.WriteString(".SH EXAMPLES\n")
		buf.WriteString(roffLine(renderText(c, c.Examples)) + "\n")
	}

	return buf.String()
}

// ManTree 在 dir 下为每个命令(含子命令递归)写一个 `.1` 文件(命名同 markdown, 扩展名 .1)。
func ManTree(app *gcli.App, dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	for _, c := range app.Commands() {
		if err := writeCmdMan(c, dir); err != nil {
			return err
		}
	}
	return nil
}

// writeCmdMan 递归写单个命令及其子命令的 man 文件。
func writeCmdMan(c *gcli.Command, dir string) error {
	file := filepath.Join(dir, cmdFileName(c)+".1")
	if err := os.WriteFile(file, []byte(CmdMan(c)), 0644); err != nil {
		return err
	}

	for _, sub := range c.Commands() {
		if err := writeCmdMan(sub, dir); err != nil {
			return err
		}
	}
	return nil
}
