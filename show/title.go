package show

import "github.com/gookit/cliapp/show/symbols"

// some constants
const (
	// position
	Left = iota
	Middle
	Right
	// other
)

// Title definition
type Title struct {
	Title     string
	Style     string
	Formatter func(t *Title) string
	// Formatter IFormatter
	Char       rune
	Width      int
	Indent     int
	Position   int
	ShowBorder bool
}

// NewTitle instance
func NewTitle(title string) *Title {
	return &Title{
		Title:    title,
		Width:    80,
		Char:     symbols.Equal,
		Indent:   2,
		Position: Left,
		Style:    "comment",
	}
}
