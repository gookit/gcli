package show

// show shown
type IShow interface {
	Print()
	String() string
}

// table a cli table show
type table struct {
	Name       string
	Cols       []string
	Rows       []string
	Border     bool
	HeadBorder bool
	RowBorder  bool
}

// Table
func Table(name string) *table {
	return &table{Name: name}
}

func (t *table) Print() {

}

func (t *table) String() string {
	return ""
}
