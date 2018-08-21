package show

// Table a cli Table show
type Table struct {
	Base // use for internal
	// Title for the table
	Title string
	// Cols the table head col data
	Cols []string
	// Rows table data rows
	Rows []interface{}
	// Border show border
	Border bool
	// RowBorder show row border
	RowBorder bool
	// HeadBorder show head border
	HeadBorder bool
	// WrapBorder wrap border for table
	WrapBorder bool
}

// NewTable
func NewTable(title string) *Table {
	return &Table{Title: title}
}

func (t *Table) Format() string {
	panic("implement me")
}
