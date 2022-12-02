package show

import "github.com/gookit/gcli/v3/show/symbols"

// Title definition
type Title struct {
	Title     string
	Style     string
	Formatter func(t *Title) string
	// Formatter IFormatter
	Char       rune
	Width      int
	Indent     int
	Align      PosFlag
	ShowBorder bool
}

// NewTitle instance
func NewTitle(title string) *Title {
	return &Title{
		Title:  title,
		Width:  80,
		Char:   symbols.Equal,
		Indent: 2,
		Align:  PosLeft,
		Style:  "comment",
	}
}
