package show

import "io"

// Writer definition
type Writer struct {
	// buf bytes.Buffer
	out io.Writer
}

// NewWriter create a new writer
func NewWriter(output io.Writer) *Writer {
	if output == nil {
		output = Output
	}

	return &Writer{
		out: output,
	}
}

// Write bytes message
func (w *Writer) Write(p []byte) (n int, err error) {
	return w.out.Write(p)
}

// Print data to io.Writer
func (w *Writer) Print() {

}

// Flush data to io.Writer
func (w *Writer) Flush() error {
	return nil
}
