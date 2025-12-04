package lists

import (
	"reflect"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/gclicom"
	"github.com/gookit/gcli/v3/show/showcom"
	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/maputil"
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
	showcom.Base // use for internal
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
		Base: showcom.Base{Out: gclicom.Output},
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

	buf := l.Buffer()
	if l.title != "" { // has title
		title := strutil.UpperWord(l.title)
		buf.WriteString(color.WrapTag(title, l.Opts.TitleStyle) + "\n")
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
			buf.WriteString(l.Opts.LeftIndent)
		}

		// format key - parsed from map, struct
		if items.itemType == ItemMap {
			key := strutil.PadRight(item.Key, " ", keyWidth)
			key = color.WrapTag(key, l.Opts.KeyStyle)
			buf.WriteString(key + l.Opts.SepChar)
		}

		// format value
		if item.IsArray() {
			arrutil.NewFormatter(item.rftVal).WithFn(func(f *arrutil.ArrFormatter) {
				f.Indent = mlIndent
				f.ClosePrefix = "  "
				// f.AfterReset = true
				f.SetOutput(buf)
			}).Format()
			buf.WriteByte('\n')
		} else if item.IsKind(reflect.Map) {
			maputil.NewFormatter(item.rftVal).WithFn(func(f *maputil.MapFormatter) {
				f.Indent = mlIndent
				f.ClosePrefix = "  "
				// f.AfterReset = true
				f.SetOutput(buf)
			}).Format()
			buf.WriteByte('\n')
		} else {
			val := item.ValString()
			if l.Opts.UpperFirst {
				val = strutil.UpperFirst(val)
			}
			buf.WriteString(val + "\n")
		}

	}
}

// String returns formatted string
func (l *List) String() string {
	l.Format()
	return l.Buf.String()
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
