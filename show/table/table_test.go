package table_test

import (
	"testing"

	"github.com/gookit/gcli/v3/show"
)

func TestNewTable(t *testing.T) {
	tb := show.NewTable("Table example1")
	tb.SetRows()

	tb.Println()
}
