// Package gclicom provides common types, definitions for gcli
package gclicom

import (
	"io"
	"os"
)

// Output the global input out stream
var Output io.Writer = os.Stdout

// SetOutput stream
func SetOutput(out io.Writer) { Output = out }

// ResetOutput stream
func ResetOutput() { Output = os.Stdout }
