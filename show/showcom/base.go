package showcom

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/gclicom"
	"github.com/gookit/goutil/strutil"
)

const (
	// OK success exit code
	OK = 0
	// ERR error exit code
	ERR = 2
)

// PosFlag type
type PosFlag = strutil.PosFlag

// some position constants
const (
	PosLeft PosFlag = iota
	PosMiddle
	PosRight
)

const PosCenter = PosMiddle

// var errInvalidType = errors.New("invalid input data type")

// FormatterFace interface
type FormatterFace interface {
	Format()
}

// ShownFace shown interface
type ShownFace interface {
	// io.WriterTo TODO
	// Format()
	// Buffer()

	// String data to string
	String() string
	// Print print current message
	Print()
	// Println print current message
	Println()
}

// Base formatter
type Base struct {
	// out  comdef.ByteStringWriter
	// TODO lock sync.Mutex
	Out io.Writer
	// Buf store formatted string
	Buf *bytes.Buffer
	Err error
}

// SetOutput for print message
func (b *Base) SetOutput(out io.Writer) { b.Out = out }

// SetBuffer field
func (b *Base) SetBuffer(buf *bytes.Buffer) { b.Buf = buf }

// Buffer get
func (b *Base) Buffer() *bytes.Buffer {
	if b.Buf == nil {
		b.Buf = new(bytes.Buffer)
	}
	return b.Buf
}

// String format given data to string
func (b *Base) String() string {
	panic("please implement the method")
}

// Format given data to string
func (b *Base) Format() { panic("please implement the method") }

// SetErr set error
func (b *Base) SetErr(err error) { b.Err = err }

// Print formatted message
func (b *Base) Print() {
	if b.Out == nil {
		b.Out = gclicom.Output
	}

	if b.Buf != nil && b.Buf.Len() > 0 {
		color.Fprint(b.Out, b.Buf.String())
		b.Buf.Reset()
	}
}

// Println formatted message and print newline
func (b *Base) Println() {
	b.Print()
	fmt.Println()
}
