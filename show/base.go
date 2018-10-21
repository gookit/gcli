package show

import (
	"github.com/gookit/color"
)

const (
	// OK success exit code
	OK = 0
	// ERR error exit code
	ERR = 2
)

// FormatterFace interface
type FormatterFace interface {
	Format() string
}

// ShownFace shown interface
type ShownFace interface {
	// data to string
	String() string
	// print current message
	Print()
	// print current message
	Println()
}

// Base formatter
type Base struct {
	// formatted string
	formatted string
}

// Format given data to string
func (f *Base) Format() string {
	panic("please implement the method")
}

// String returns formatted string
func (f *Base) String() string {
	return f.Format()
}

// Print formatted message
func (f *Base) Print() {
	color.Print(f.Format())
}

// Println formatted message and print newline
func (f *Base) Println() {
	color.Println(f.Format())
}
