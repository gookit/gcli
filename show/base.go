package show

import (
	"io"
	"os"
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

// var errInvalidType = errors.New("invalid input data type")

// FormatterFace interface
type FormatterFace interface {
	Format() string
}

// ShownFace shown interface
type ShownFace interface {
	// String data to string
	String() string
	// Print print current message
	Print()
	// Println print current message
	Println()
}

// Base formatter
type Base struct {
	output io.Writer
	// formatted string
	formatted string
}

// SetOutput for print message
func (b *Base) SetOutput(output io.Writer) {
	b.output = output
}

// Format given data to string
func (b *Base) Format() string {
	panic("please implement the method")
}

// Print formatted message
func (b *Base) Print() {
	if b.output == nil {
		b.output = os.Stdout
	}

	if b.formatted != "" {
		color.Fprint(b.output, b.formatted)
		// clear data
		b.formatted = ""
	}
}

// Println formatted message and print newline
func (b *Base) Println() {
	if b.output == nil {
		b.output = os.Stdout
	}

	if b.formatted != "" {
		color.Fprintln(b.output, b.formatted)
		// clear data
		b.formatted = ""
	}
}

/*************************************************************
 * Data item(s)
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
	data interface{}
	// inner context
	itemType    string
	rowNumber   int
	keyMaxWidth int
}

// NewItems create a Items for data.
func NewItems(data interface{}) *Items {
	items := &Items{
		data:     data,
		itemType: ItemMap,
	}
	rv := reflect.ValueOf(data)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	var keyWidth int
	switch rv.Kind() {
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
			// skip don't exported field
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
	default:
		panic("GCLI: invalid data type, only allow: array, map, slice, struct")
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

// Item definition
type Item struct {
	// Val string
	Key string
	// info
	index int
	// valLen int
	keyLen int
	// rawVal interface{}
	rftVal reflect.Value
}

func newItem(key interface{}, rv reflect.Value, index int) *Item {
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
func (item *Item) Kind() reflect.Kind {
	return item.rftVal.Kind()
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
	}
	return false
}

// ValString get
func (item *Item) ValString() string {
	return strutil.QuietString(item.rftVal.Interface())
}

// RftVal get
func (item *Item) RftVal() reflect.Value {
	return item.rftVal
}

func (item *Item) maxLen(ln int) int {
	if item.keyLen > ln {
		return item.keyLen
	}
	return ln
}
