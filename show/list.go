package show

import (
	"bytes"
	"reflect"

	"github.com/gookit/color"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/maputil"
	"github.com/gookit/goutil/reflects"
	"github.com/gookit/goutil/structs"
	"github.com/gookit/goutil/strutil"
)

// ListOption definition
type ListOption struct {
	// IgnoreEmpty ignore empty value item. default: true
	IgnoreEmpty bool
	// UpperFirst upper first char for the item.value. default: false
	UpperFirst  bool
	SepChar     string // split key value
	LeftIndent  string
	KeyWidth    int // if not set, will be auto-detected.
	KeyMinWidth int
	KeyStyle    string
	ValueStyle  string
	TitleStyle  string
	// FilterFunc filter item.
	//  - return true to show item, otherwise hide item.
	FilterFunc func(item *Item) bool
}

// ListOpFunc define
type ListOpFunc func(opts *ListOption)

// NewListOption instance
func NewListOption() *ListOption {
	return &ListOption{
		SepChar:  " ",
		KeyStyle: "info",
		// more
		LeftIndent:  "  ",
		KeyMinWidth: 8,
		IgnoreEmpty: true,
		TitleStyle:  "comment",
	}
}

/*************************************************************
 * region List
 *************************************************************/

// List definition. data allow type: struct, slice, array, map
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
	data any
}

// NewList instance.
//
// data allow type:
//
//	struct, slice, array, map
func NewList(title string, data any, fns ...ListOpFunc) *List {
	l := &List{
		title: title,
		data:  data,
		// base
		Base: Base{out: Output},
		// options
		Opts: NewListOption(),
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
func (l *List) Format() {
	if l.data == nil {
		return
	}

	if l.buf == nil {
		l.buf = new(bytes.Buffer)
	}

	if l.title != "" { // has title
		title := strutil.UpperWord(l.title)
		l.buf.WriteString(color.WrapTag(title, l.Opts.TitleStyle) + "\n")
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
		// call filter func
		if l.Opts.FilterFunc != nil && !l.Opts.FilterFunc(item) {
			continue
		}

		if l.Opts.LeftIndent != "" {
			l.buf.WriteString(l.Opts.LeftIndent)
		}

		// format key - parsed from map, struct
		if items.itemType == ItemMap {
			key := strutil.PadRight(item.Key, " ", keyWidth)
			key = color.WrapTag(key, l.Opts.KeyStyle)
			l.buf.WriteString(key + l.Opts.SepChar)
		}

		// format value
		if item.IsArray() {
			arrutil.NewFormatter(item.rftVal).WithFn(func(f *arrutil.ArrFormatter) {
				f.Indent = mlIndent
				f.ClosePrefix = "  "
				// f.AfterReset = true
				f.SetOutput(l.buf)
			}).Format()
			l.buf.WriteByte('\n')
		} else if item.IsKind(reflect.Map) {
			maputil.NewFormatter(item.rftVal).WithFn(func(f *maputil.MapFormatter) {
				f.Indent = mlIndent
				f.ClosePrefix = "  "
				// f.AfterReset = true
				f.SetOutput(l.buf)
			}).Format()
			l.buf.WriteByte('\n')
		} else {
			val := item.ValString()
			if l.Opts.UpperFirst {
				val = strutil.UpperFirst(val)
			}
			l.buf.WriteString(val + "\n")
		}

	}
}

// String returns formatted string
func (l *List) String() string {
	l.Format()
	return l.buf.String()
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
func (l *List) Flush() { l.Println() }

/*************************************************************
 * region Lists
 *************************************************************/

// Lists use for formatting and printing multi list data
type Lists struct {
	Base // use for internal
	err error
	// options
	Opts *ListOption
	rows []*List
	// data buffer
	buffer *bytes.Buffer
}

// NewEmptyLists create empty lists
func NewEmptyLists(fns ...ListOpFunc) *Lists {
	ls := &Lists{
		Base: Base{out: Output},
		Opts: NewListOption(),
	}
	return ls.WithOptionFns(fns)
}

// NewLists create lists. allow: map[string]any, struct-ptr
func NewLists(mlist any, fns ...ListOpFunc) *Lists {
	ls := NewEmptyLists()
	rv := reflect.Indirect(reflect.ValueOf(mlist))

	if rv.Kind() == reflect.Map {
		ls.err = reflects.EachStrAnyMap(rv, func(key string, val any) {
			ls.AddSublist(key, val)
		})
	} else if rv.Kind() == reflect.Struct {
		for title, data := range structs.ToMap(mlist) {
			ls.rows = append(ls.rows, NewList(title, data))
		}
	} else {
		panic("Lists: not support type: " + rv.Kind().String())
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

// WithOptions with an options func list
func (ls *Lists) WithOptions(fns ...ListOpFunc) *Lists {
	return ls.WithOptionFns(fns)
}

// AddSublist with options func list
func (ls *Lists) AddSublist(title string, data any) *Lists {
	ls.rows = append(ls.rows, NewList(title, data))
	return ls
}

// Format as string
func (ls *Lists) Format() {
	if len(ls.rows) == 0 {
		return
	}

	ls.buffer = new(bytes.Buffer)

	for _, list := range ls.rows {
		list.Opts = ls.Opts
		list.SetBuffer(ls.buffer)
		list.Format()
	}
}

// String returns formatted string
func (ls *Lists) String() string {
	ls.Format()
	return ls.buf.String()
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
func (ls *Lists) Flush() { ls.Println() }
