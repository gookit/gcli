package show

// Title definition
type Title struct {
	Title     string
	Formatter func(t *Title) string
	// Formatter IFormatter
}

// NewTitle instance
func NewTitle(title string) *Title {
	return &Title{Title: title}
}
