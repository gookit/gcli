package show

// List definition
//
// String len:
// 	len("你好"), len("hello"), len("hello你好") -> 6 5 11
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

// Format as string
func (l *List) Format() string {
	if l.Items == nil {
		return ""
	}

	return ""
}

// Print to console
func (l *List) Print() {

}

// Lists definition
type Lists struct {
	Base  // use for internal
	Title string
	Rows  []List
}

// NewLists create lists
func NewLists(title string, lists []List) *Lists {
	return &Lists{Title: title, Rows: lists}
}

// Format as string
func (ls *Lists) Format() string {
	panic("implement me")
}

// Print to console
func (ls *Lists) Print() {
	panic("implement me")
}
