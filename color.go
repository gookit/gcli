package cliapp

import "fmt"

// Color represents a text color.
type Color uint8

// Foreground colors.
const (
	// basic Foreground colors 30 - 37
	FgBlack   Color = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite

	FgDefault Color = 39

	// extra Foreground color 90 - 97
	FgDarkGray     Color = iota + 90
	FgLightRed
	FgLightGreen
	FgLightYellow
	FgLightBlue
	FgLightMagenta
	FgLightCyan
	FgWhiteEx
)

// Background colors.
const (
	// basic Background colors 40 - 47
	BgBlack   Color = iota + 40
	BgRed
	BgGreen
	BgYellow   // BgBrown like yellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
	BgDefault Color = 49

	// extra Background color 100 - 107
	BgDarkGray     Color = iota + 100
	BgLightRed
	BgLightGreen
	BgLightYellow
	BgLightBlue
	BgLightMagenta
	BgLightCyan
	BgWhiteEx
)

// color options
const (
	OpBold       = 1 // 加粗
	OpFuzzy      = 2 // 模糊(不是所有的终端仿真器都支持)
	OpItalic     = 3 // 斜体(不是所有的终端仿真器都支持)
	OpUnderscore = 4 // 下划线
	OpBlink      = 5 // 闪烁
	OpReverse    = 7 // 颠倒的 交换背景色与前景色
	OpConcealed  = 8 // 隐匿的
)

// CLI color template
// const ColorTpl = "\033[%sm%s\033[0m"
// const ColorTpl = "\x1b[%sm%s\x1b[0m"
const ColorTpl = "\x1b[%dm%s\x1b[0m"

// S adds the coloring to the given string.
// usage: cli.Color(cli.FgCyan).S("string")
func (c Color) S(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}

// F adds the coloring to the given string.
// usage: cli.Color(cli.FgCyan).F("string %s", "arg0")
func (c Color) F(s string, args ...interface{}) string {
	s = fmt.Sprintf(s, args...)

	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}
