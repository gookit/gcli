package lists

import (
	"bytes"
	"reflect"

	"github.com/gookit/gcli/v3/gclicom"
	"github.com/gookit/gcli/v3/show/showcom"
	"github.com/gookit/goutil/reflects"
	"github.com/gookit/goutil/structs"
)

/*************************************************************
 * region Lists
 *************************************************************/

// Lists use for formatting and printing multi list data
type Lists struct {
	showcom.Base // use for internal
	// options
	Opts *ListOption
	rows []*List
}

// NewEmptyLists create empty lists
func NewEmptyLists(fns ...ListOpFunc) *Lists {
	ls := &Lists{
		Base: showcom.Base{Out: gclicom.Output},
		Opts: NewListOption(),
	}
	return ls.WithOptionFns(fns)
}

// NewLists create lists. allow: map[string]any, struct-ptr
func NewLists(mlist any, fns ...ListOpFunc) *Lists {
	ls := NewEmptyLists()
	rv := reflect.Indirect(reflect.ValueOf(mlist))

	if rv.Kind() == reflect.Map {
		ls.Err = reflects.EachStrAnyMap(rv, func(key string, val any) {
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

	ls.Buf = new(bytes.Buffer)
	for _, list := range ls.rows {
		list.Opts = ls.Opts
		list.SetBuffer(ls.Buf)
		list.Format()
	}
}

// String returns formatted string
func (ls *Lists) String() string {
	ls.Format()
	return ls.Buf.String()
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
