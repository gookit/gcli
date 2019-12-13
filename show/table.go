package show

/*
━━━┯━━━━━━━┯━━━━━━━━━━━━━━━━━┯━━━━━━━━━━┯━━━━━━━━━━
 # │ pid   │ name            │ status   │ cpu
───┼───────┼─────────────────┼──────────┼──────────
 0 │   992 │ chrome          │ Sleeping │ 6.988768
 2 │ 13973 │ qemu-system-x86 │ Sleeping │ 4.996551
━━━┷━━━━━━━┷━━━━━━━━━━━━━━━━━┷━━━━━━━━━━┷━━━━━━━━━━
*/

// Table a cli Table show
type Table struct {
	Base // use for internal
	// Title for the table
	Title string
	// Cols the table head col data
	Cols []string
	// Rows table data rows
	Rows []interface{}
	// options ...
	// HasBorder show border line
	HasBorder bool
	// RowBorder show row border
	RowBorder bool
	// HeadBorder show head border
	HeadBorder bool
	// WrapBorder wrap border for table
	WrapBorder bool
}

// NewTable create table
func NewTable(title string) *Table {
	return &Table{Title: title}
}

// Format as string
func (t *Table) Format() string {
	panic("implement me")
}
