package show

import "io"

// Writer definition
type Writer struct {
	output io.Writer
}

// NewWriter create a new writer
func NewWriter(output io.Writer) *Writer {
	return &Writer{}
}

// Write
func (w *Writer) Write(buf []byte) (n int, err error) {
	return
}

// Print data to io.Writer
func (w *Writer) Print() {

}

// Flush data to io.Writer
func (w *Writer) Flush() error {
	return nil
}
