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

// Formatter interface
type Formatter interface {
	Format()
}

type FormatFunc func()

// Format implement FormatterFace
func (fn FormatFunc) Format() {
	fn()
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

// Base formatter. NOTE: must config the FormatFn before use
type Base struct {
	// out  comdef.ByteStringWriter
	// TODO lock sync.Mutex
	Out io.Writer
	// Buf store formatted string
	Buf *bytes.Buffer
	Err error
	// FormatFn function
	FormatFn  FormatFunc
	formatted bool
}

// SetOutput for print message
func (b *Base) SetOutput(out io.Writer) { b.Out = out }

// SetBuffer field
func (b *Base) SetBuffer(buf *bytes.Buffer) { b.Buf = buf }

// InitBuffer instance
func (b *Base) InitBuffer() {
	if b.Buf == nil {
		b.Buf = new(bytes.Buffer)
	} else {
		b.Buf.Reset()
	}
}

// Buffer get buffer instance
func (b *Base) Buffer() *bytes.Buffer {
	if b.Buf == nil {
		b.Buf = new(bytes.Buffer)
	}
	return b.Buf
}

// String format given data to string
func (b *Base) String() string {
	b.format()
	return b.Buf.String()
}

// Format given data to string
func (b *Base) format() {
	if b.formatted {
		return
	}
	b.formatted = true

	if b.FormatFn == nil {
		panic("gcli/show: please set the FormatFn")
	}
	b.FormatFn()
}

// WriteTo format to string and write to w.
func (b *Base) WriteTo(w io.Writer) (int64, error) {
	b.format()
	return b.Buf.WriteTo(w)
}

// Print formatted message
func (b *Base) Print() {
	if b.Out == nil {
		b.Out = gclicom.Output
	}

	// call format
	b.format()

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

// SetErr set error
func (b *Base) SetErr(err error) { b.Err = err }
