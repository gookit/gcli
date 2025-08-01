package show

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/gookit/color"
	"github.com/gookit/goutil/reflects"
	"github.com/gookit/goutil/strutil"
)

const (
	// OK success exit code
	OK = 0
	// ERR error exit code
	ERR = 2
)

// PosFlag type
type PosFlag = strutil.PosFlag

// some position constants
const (
	PosLeft PosFlag = iota
	PosRight
	PosMiddle
)

// var errInvalidType = errors.New("invalid input data type")

// FormatterFace interface
type FormatterFace interface {
	Format()
}

// ShownFace shown interface
type ShownFace interface {
	// io.WriterTo TODO
	// Format()
	// Buffer()

	// String data to string
	String() string
	// Print print current message
	Print()
	// Println print current message
	Println()
}

// Base formatter
type Base struct {
	// TODO lock sync.Mutex
	out io.Writer
	// formatted string
	buf *bytes.Buffer
	err error
}

// SetOutput for print message
func (b *Base) SetOutput(out io.Writer) {
	b.out = out
}

// SetBuffer field
func (b *Base) SetBuffer(buf *bytes.Buffer) {
	b.buf = buf
}

// Buffer get
func (b *Base) Buffer() *bytes.Buffer {
	if b.buf == nil {
		b.buf = new(bytes.Buffer)
	}
	return b.buf
}

// String format given data to string
func (b *Base) String() string {
	panic("please implement the method")
}

// Format given data to string
func (b *Base) Format() {
	panic("please implement the method")
}

// Err get
func (b *Base) Err() error {
	return b.err
}

// Print formatted message
func (b *Base) Print() {
	if b.out == nil {
		b.out = Output
	}

	if b.buf != nil && b.buf.Len() > 0 {
		color.Fprint(b.out, b.buf.String())
		b.buf.Reset()
	}
}

// Println formatted message and print newline
func (b *Base) Println() {
	b.Print()
	fmt.Println()
}

/*************************************************************
 * region Data item(s)
 *************************************************************/

const (
	// ItemMap parsed from map, struct
	ItemMap = "map"
	// ItemList parsed from array, slice
	ItemList = "list"
)

// Items definition
type Items struct {
	List []*Item
	// raw data
	data any
	// inner context
	itemType    string
	rowNumber   int
	keyMaxWidth int
}

// NewItems create an Items for data.
func NewItems(data any) *Items {
	items := &Items{
		data:     data,
		itemType: ItemMap,
	}
	rv := reflect.ValueOf(data)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	var keyWidth int
	kind := rv.Kind()

	switch kind {
	case reflect.Map:
		mapKeys := rv.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			key := mapKeys[i]
			rftV := rv.MapIndex(key)
			item := newItem(key.Interface(), rftV, i)

			items.List = append(items.List, item)
			keyWidth = item.maxLen(keyWidth)
		}
	case reflect.Slice, reflect.Array:
		items.itemType = ItemList
		for i := 0; i < rv.Len(); i++ {
			rftV := rv.Index(i)
			item := newItem("", rftV, i)

			items.List = append(items.List, item)
			keyWidth = item.maxLen(keyWidth)
		}
	case reflect.Struct:
		// structs.ToMap()
		rt := rv.Type()
		for i := 0; i < rt.NumField(); i++ {
			ft := rt.Field(i)
			fv := rv.Field(i)
			// skip unexported field
			name := ft.Name
			if name[0] >= 'a' && name[0] <= 'z' {
				continue
			}

			tagName := ft.Tag.Get("json")
			if tagName == "" {
				tagName = ft.Name
			} else if pos := strings.IndexByte(tagName, ','); pos > 0 {
				tagName = tagName[:pos]
			}

			item := newItem(tagName, fv, i)
			items.List = append(items.List, item)
			keyWidth = item.maxLen(keyWidth)
		}
	case reflect.String:
		items.List = append(items.List, newItem("", rv, 0))
	default:
		if reflects.IsAnyInt(kind) {
			items.List = append(items.List, newItem("", rv, 0))
		} else {
			panic("GCLI.show: unsupported data type: " + rv.Kind().String())
		}
	}

	// settings
	items.rowNumber = len(items.List)
	items.keyMaxWidth = keyWidth
	return items
}

// KeyMaxWidth get
func (its *Items) KeyMaxWidth(userSetting int) int {
	if userSetting > 0 {
		return userSetting
	}
	return its.keyMaxWidth
}

// ItemType get
func (its *Items) ItemType() string {
	return its.itemType
}

// Each handle item in the items.List
func (its *Items) Each(fn func(item *Item)) {
	for _, item := range its.List {
		fn(item)
	}
}

/*************************************************************
 * region Data item
 *************************************************************/

// Item definition
type Item struct {
	// Val string
	Key string
	// info
	index int
	// valLen int
	keyLen int
	// rawVal any
	rftVal reflect.Value
}

func newItem(key any, rv reflect.Value, index int) *Item {
	item := &Item{
		Key: strutil.QuietString(key),
		// Val:    fmt.Sprint(value),
		index: index,
		// rawVal: value,
		rftVal: reflects.Elem(rv),
	}

	if item.Key != "" {
		item.keyLen = utf8.RuneCountInString(item.Key)
	}

	// item.valLen = utf8.RuneCountInString(item.Val)
	return item
}

// Kind get
func (item *Item) Kind() reflect.Kind { return item.rftVal.Kind() }

// IsKind check
func (item *Item) IsKind(kind reflect.Kind) bool {
	return item.rftVal.Kind() == kind
}

// IsArray get is array, slice
func (item *Item) IsArray() bool {
	return item.rftVal.Kind() == reflect.Array || item.rftVal.Kind() == reflect.Slice
}

// IsEmpty get value is empty: nil, empty string.
func (item *Item) IsEmpty() bool {
	switch item.rftVal.Kind() {
	case reflect.String:
		return item.rftVal.Len() == 0
	case reflect.Interface, reflect.Slice, reflect.Ptr:
		return item.rftVal.IsNil()
	default:
		return !item.rftVal.IsValid()
	}
}

// ValString get
func (item *Item) ValString() string {
	if item.IsEmpty() {
		return ""
	}
	return strutil.QuietString(item.rftVal.Interface())
}

// RftVal get
func (item *Item) RftVal() reflect.Value { return item.rftVal }

func (item *Item) maxLen(ln int) int {
	if item.keyLen > ln {
		return item.keyLen
	}
	return ln
}
