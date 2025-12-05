package lists

import (
	"bytes"
	"reflect"

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
	Opts *Options
	rows []*List
}

// NewEmptyLists create empty lists
func NewEmptyLists(fns ...ListOpFunc) *Lists {
	ls := &Lists{
		Opts: NewOptions(),
	}

	ls.FormatFn = ls.Format
	return ls.WithOptionFns(fns)
}

// NewLists create lists. allow: map[string]any, struct-ptr
func NewLists(mList any, fns ...ListOpFunc) *Lists {
	ls := NewEmptyLists()
	rv := reflect.Indirect(reflect.ValueOf(mList))

	if rv.Kind() == reflect.Map {
		ls.Err = reflects.EachStrAnyMap(rv, func(key string, val any) {
			ls.AddSublist(key, val)
		})
	} else if rv.Kind() == reflect.Struct {
		for title, data := range structs.ToMap(mList) {
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

// Flush formatted message to console
func (ls *Lists) Flush() { ls.Println() }
