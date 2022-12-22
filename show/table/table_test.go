package table_test

import (
	"testing"

	"github.com/gookit/gcli/v3/show/table"
)

func TestNewTable(t *testing.T) {
	tb := table.New("Table example1")
	tb.SetRows([]any{
		// TODO ...
	})

	// tb.Println()
}
