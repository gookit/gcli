package show

import (
	"bytes"
	"github.com/gookit/color"
	"github.com/gookit/filter"
	"github.com/gookit/goutil/str"
	"os"
)

// ListOption definition
type ListOption struct {
	// IgnoreEmpty ignore empty value item
	IgnoreEmpty bool
	// UpperFirst upper first char for the item.value
	UpperFirst  bool
	SepChar     string // split key value
	LeftIndent  string
	KeyWidth 	int // if not set, will auto detect.
	KeyMinWidth int
	KeyStyle    string
	ValueStyle  string
	TitleStyle  string
}


/*************************************************************
 * List
 *************************************************************/

// List definition
//
// String len:
// 	len("你好"), len("hello"), len("hello你好") -> 6 5 11
type List struct {
	Base // use for internal
	// Title list title name
	title string
	// list data. allow type: struct, slice, array, map
	data interface{}
	buffer *bytes.Buffer
	// IgnoreEmpty ignore empty value item
	IgnoreEmpty bool
	// UpperFirst upper first char for the item.value
	UpperFirst  bool
	SepChar     string // split key value
	LeftIndent  string
	KeyWidth 	int // if not set, will auto detect.
	KeyMinWidth int
	KeyStyle    string
	ValueStyle  string
	TitleStyle  string
}

// SetBuffer field
func (l *List) SetBuffer(buffer *bytes.Buffer) {
	l.buffer = buffer
}

// NewList instance
func NewList(title string, data interface{}) *List {
	return &List{
		title: title,
		data:  data,
		// base
		Base: Base{output: os.Stdout},
		// options
		SepChar:    " ",
		LeftIndent: "  ",
		KeyStyle:   "info",
		//
		KeyMinWidth: 8,
		IgnoreEmpty: true,
		TitleStyle:  "comment",
	}
}

// Format as string
func (l *List) Format() string {
	if l.data == nil || l.formatted != "" {
		return l.formatted
	}

	if l.buffer == nil {
		l.buffer = new(bytes.Buffer)
	}

	if l.title != "" { // has title
		title := filter.UpperWord(l.title)
		l.buffer.WriteString(color.WrapTag(title, l.TitleStyle) + "\n")
	}

	items := NewItems(l.data) // build items
	keyWidth := l.KeyWidth
	if keyWidth <= 0 {
		keyWidth = items.KeyMaxWidth()
	}

	if keyWidth < l.KeyMinWidth {
		keyWidth = l.KeyMinWidth
	}

	for _, item := range items.List {
		if l.IgnoreEmpty && item.Val == "" {
			continue
		}

		if l.LeftIndent != "" {
			l.buffer.WriteString(l.LeftIndent)
		}

		// parsed from map, struct
		if items.itemType == ItemMap {
			key := str.PadRight(item.Key, " ", keyWidth)
			key = color.WrapTag(key, l.KeyStyle)
			l.buffer.WriteString(key + l.SepChar)
		}

		val := item.Val
		if l.UpperFirst {
			val = str.UpperFirst(val)
		}

		l.buffer.WriteString(val + "\n")
	}

	l.formatted = l.buffer.String()
	return l.formatted
}

// String returns formatted string
func (l *List) String() string {
	return l.Format()
}

// Println formatted message
func (l *List) Print() {
	l.Format()
	l.Base.Print()
}

// Println formatted message with newline
func (l *List) Println() {
	l.Format()
	l.Base.Println()
}

// Flush formatted message to console
func (l *List) Flush() {
	l.Println()
	l.buffer.Reset()
	l.formatted = ""
}

/*************************************************************
 * Lists
 *************************************************************/

// Lists definition
type Lists struct {
	Base // use for internal
	rows  []*List
	// NewLine print "\n" at last
	NewLine bool
	buffer *bytes.Buffer
}

// NewLists create lists
func NewLists(listMap map[string]interface{}) *Lists {
	ls := &Lists{}

	for title, data := range listMap {
		ls.rows = append(ls.rows, NewList(title, data))
	}

	return ls
}

// Format as string
func (ls *Lists) Format() string {
	if len(ls.rows) == 0 || ls.formatted != "" {
		return ls.formatted
	}

	ls.buffer = new(bytes.Buffer)

	for _, list := range ls.rows {
		list.SetBuffer(ls.buffer)
		list.Format()
	}

	ls.formatted = ls.buffer.String()
	return ls.formatted
}

// String returns formatted string
func (ls *Lists) String() string {
	return ls.Format()
}

// Println formatted message
func (ls *Lists) Print() {
	ls.Format()
	ls.Base.Print()
}

// Println formatted message with newline
func (ls *Lists) Println() {
	ls.Format()
	ls.Base.Println()
}

// Flush formatted message to console
func (ls *Lists) Flush() {
	ls.Println()
	ls.buffer.Reset()
	ls.formatted = ""
}
