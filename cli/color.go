package cli

import "fmt"

// Color represents a text color.
type Color uint8

// Foreground colors.
const (
	FgBlack Color = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Background colors.
const (
	BgBlack Color = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

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

