package show

import (
	"fmt"
	"github.com/gookit/color"
)

const (
	// OK success exit code
	OK = 0
	// ERR error exit code
	ERR = 2
)

// show shown
type IShow interface {
	// print current message
	Print()
	// trans to string
	String() string
}

type Title struct {
	Title     string
	Formatter func(t *Title) string
	// Formatter IFormatter
}

func NewTitle(title string) *Title {
	return &Title{Title: title}
}

// Error message print
func Error(format string, v ...interface{}) int {
	fmt.Println(color.FgRed.Render("ERROR:"), fmt.Sprintf(format, v...))
	return ERR
}
