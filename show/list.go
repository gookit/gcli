package show

// List definition
type List struct {
	Base // use for internal
	// Title list title name
	Title string
	// Items list items.
	// allow:
	// 	struct, []int, []string, []interface{}, map[string]string, map[string]interface{}
	Items interface{}
	// NewLine print "\n" at last
	NewLine bool
	// formatted string
	formatted string
}

// NewList instance
func NewList(title string, items interface{}) *List {
	return &List{Title: title, Items: items, NewLine: true}
}

func (l *List) Format() string {
	if l.Items == nil {
		return ""
	}

	return ""
}

func (l *List) Print() {

}

// Lists definition
type Lists struct {
	Base  // use for internal
	Title string
	Rows  []List
}

// NewLists
func NewLists(title string, lists []List) *Lists {
	return &Lists{Title: title, Rows: lists}
}

func (ls *Lists) Format() string {
	panic("implement me")
}

func (ls *Lists) Print() {
	panic("implement me")
}
