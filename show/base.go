package show

// show shown
type IShow interface {
	// print current message
	Print()
	// trans to string
	String() string

}

type Title struct {
	Title string
	Formatter func(t *Title) string
	// Formatter IFormatter
}

func NewTitle(title string) *Title {
	return &Title{Title: title}
}
