package show

// List list
type List struct {
	Name  string
	Items interface{}
}

// NewList
func NewList(name string, items interface{}) *List {
	return &List{name, items}
}

func (l *List) Format() string {
	panic("implement me")
}

func (l *List) String() string {
	panic("implement me")
}

func (l *List) Print() {

}

// Lists lists
type Lists struct {
	Title string
	Rows  []List
}

// NewLists
func NewLists(title string, lists []List) *Lists {
	return &Lists{title, lists}
}

func (ls *Lists) Format() {
	panic("implement me")
}

func (ls *Lists) Print() {
	panic("implement me")
}

func (ls *Lists) String() string {
	panic("implement me")
}
