package show

import (
	"bytes"
	"os"
	"reflect"

	"github.com/gookit/color"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/strutil"
)

// ListOption definition
type ListOption struct {
	// IgnoreEmpty ignore empty value item
	IgnoreEmpty bool
	// UpperFirst upper first char for the item.value
	UpperFirst  bool
	SepChar     string // split key value
	LeftIndent  string
	KeyWidth    int // if not set, will be auto-detected.
	KeyMinWidth int
	KeyStyle    string
	ValueStyle  string
	TitleStyle  string
}

// ListOpFunc define
type ListOpFunc func(opts *ListOption)

/*************************************************************
 * List
 *************************************************************/

// List definition
//
// String len:
//
//	len("你好"), len("hello"), len("hello你好") -> 6 5 11
type List struct {
	Base // use for internal
	// options
	Opts *ListOption
	// Title list title name
	title string
	// list data. allow type: struct, slice, array, map
	data interface{}
	// formatted data buffer
	buffer *bytes.Buffer
}

// SetBuffer field
func (l *List) SetBuffer(buffer *bytes.Buffer) {
	l.buffer = buffer
}

// NewList instance.
//
// data allow type:
//
//	struct, slice, array, map
func NewList(title string, data interface{}, fns ...ListOpFunc) *List {
	l := &List{
		title: title,
		data:  data,
		// base
		Base: Base{output: os.Stdout},
		// options
		Opts: &ListOption{
			SepChar:    " ",
			KeyStyle:   "info",
			LeftIndent: "  ",
			// more settings
			KeyMinWidth: 8,
			IgnoreEmpty: true,
			TitleStyle:  "comment",
		},
	}

	return l.WithOptionFns(fns)
}

// WithOptionFns with options func
func (l *List) WithOptionFns(fns []ListOpFunc) *List {
	for _, fn := range fns {
		if fn != nil {
			fn(l.Opts)
		}
	}
	return l
}

// WithOptions with options func
func (l *List) WithOptions(fns ...ListOpFunc) *List {
	return l.WithOptionFns(fns)
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
		title := strutil.UpperWord(l.title)
		l.buffer.WriteString(color.WrapTag(title, l.Opts.TitleStyle) + "\n")
	}

	items := NewItems(l.data) // build items
	keyWidth := items.KeyMaxWidth(l.Opts.KeyWidth)

	if keyWidth < l.Opts.KeyMinWidth {
		keyWidth = l.Opts.KeyMinWidth
	}

	// multi line indent
	mlIndent := l.Opts.LeftIndent + strutil.Repeat(" ", len(l.Opts.LeftIndent))

	for _, item := range items.List {
		if l.Opts.IgnoreEmpty && item.IsEmpty() {
			continue
		}

		if l.Opts.LeftIndent != "" {
			l.buffer.WriteString(l.Opts.LeftIndent)
		}

		// format key - parsed from map, struct
		if items.itemType == ItemMap {
			key := strutil.PadRight(item.Key, " ", keyWidth)
			key = color.WrapTag(key, l.Opts.KeyStyle)
			l.buffer.WriteString(key + l.Opts.SepChar)
		}

		// format value
		if item.IsArray() {
			arrutil.NewFormatter(item.rftVal).WithFn(func(f *arrutil.ArrFormatter) {
				f.Indent = mlIndent
				f.ClosePrefix = "  "
				// f.AfterReset = true
				f.SetOutput(l.buffer)
			}).Format()
			l.buffer.WriteByte('\n')
		} else if item.Kind() == reflect.Map {
			maputil.NewFormatter(item.rftVal).WithFn(func(f *maputil.MapFormatter) {
				f.Indent = mlIndent
				f.ClosePrefix = "  "
				// f.AfterReset = true
				f.SetOutput(l.buffer)
			}).Format()
			l.buffer.WriteByte('\n')
		} else {
			val := item.ValString()
			if l.Opts.UpperFirst {
				val = strutil.UpperFirst(val)
			}
			l.buffer.WriteString(val + "\n")
		}

	}

	l.formatted = l.buffer.String()
	return l.formatted
}

// String returns formatted string
func (l *List) String() string {
	return l.Format()
}

// Print formatted message
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
	// options
	Opts *ListOption
	rows []*List
	// data buffer
	buffer *bytes.Buffer
}

// NewLists create lists
func NewLists(listMap map[string]interface{}, fns ...ListOpFunc) *Lists {
	ls := &Lists{
		Opts: &ListOption{
			SepChar:  " ",
			KeyStyle: "info",
			// more
			LeftIndent:  "  ",
			KeyMinWidth: 8,
			IgnoreEmpty: true,
			TitleStyle:  "comment",
		},
	}

	for title, data := range listMap {
		ls.rows = append(ls.rows, NewList(title, data))
	}
	return ls.WithOptionFns(fns)
}

// WithOptionFns with options func
func (ls *Lists) WithOptionFns(fns []ListOpFunc) *Lists {
	for _, fn := range fns {
		fn(ls.Opts)
	}
	return ls
}

// WithOptions with options func list
func (ls *Lists) WithOptions(fns ...ListOpFunc) *Lists {
	return ls.WithOptionFns(fns)
}

// Format as string
func (ls *Lists) Format() string {
	if len(ls.rows) == 0 || ls.formatted != "" {
		return ls.formatted
	}

	ls.buffer = new(bytes.Buffer)

	for _, list := range ls.rows {
		list.Opts = ls.Opts
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

// Print formatted message
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
