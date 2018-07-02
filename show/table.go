package show

// Table a cli Table show
type Table struct {
	Name       string
	Cols       []string
	Rows       []interface{}
	Border     bool
	RowBorder  bool
	HeadBorder bool
	WrapBorder  bool
}

// NewTable
func NewTable(name string) *Table {
	return &Table{Name: name}
}

func (t *Table) Format() string {
	panic("implement me")
}

func (t *Table) Print() {

}

func (t *Table) String() string {
	return ""
}

