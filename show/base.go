package show

type IShow interface {
	Show()
	String() string
}

// Table a cli table show
type Table struct {
	Name       string
	Cols       []string
	Rows       []string
	Border     bool
	HeadBorder bool
	RowBorder  bool
}

// Table
func NewTable(name string) *Table {
	return &Table{Name: name}
}

func (t *Table) Show() {

}

func (t *Table) String() string {
	return ""
}
