package show

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gookit/color"
	"io"
	"os"
	"reflect"
	"unicode/utf8"
)

const (
	// OK success exit code
	OK = 0
	// ERR error exit code
	ERR = 2
)

var errInvalidType = errors.New("invalid input data type")

// FormatterFace interface
type FormatterFace interface {
	Format() string
}

// ShownFace shown interface
type ShownFace interface {
	// data to string
	String() string
	// print current message
	Print()
	// print current message
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
	}
}

// Println formatted message and print newline
func (b *Base) Println() {
	if b.output == nil {
		b.output = os.Stdout
	}

	if b.formatted != "" {
		color.Fprintln(b.output, b.formatted)
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
	data interface{}
	//
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
			item := newItem(key.Interface(), rv.MapIndex(key).Interface(), i)
			items.List = append(items.List, item)
			// max len
			keyWidth = item.maxLen(keyWidth)
		}
	case reflect.Slice, reflect.Array:
		items.itemType = ItemList
		for i := 0; i < rv.Len(); i++ {
			item := newItem("", rv.Index(i).Interface(), i)
			items.List = append(items.List, item)
			// max len
			keyWidth = item.maxLen(keyWidth)
		}
	case reflect.Struct:
		bs, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}

		mp := make(map[string]interface{})
		if err = json.Unmarshal(bs, &mp); err != nil {
			panic(err)
		}

		for key, val := range mp {
			item := newItem(key, val, 0)
			items.List = append(items.List, item)
			// max len
			keyWidth = item.maxLen(keyWidth)
		}
	default:
		panic("invalid data type, only allow: array, map, slice, struct")
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
	Key string
	Val string
	// info
	index  int
	keyLen int
	valLen int
}

func newItem(key, value interface{}, index int) *Item {
	item := &Item{
		Key:   fmt.Sprint(key),
		Val:   fmt.Sprint(value),
		index: index,
	}

	if item.Key != "" {
		item.keyLen = utf8.RuneCountInString(item.Key)
	}

	item.valLen = utf8.RuneCountInString(item.Val)
	return item
}

func (item *Item) maxLen(ln int) int {
	if item.keyLen > ln {
		return item.keyLen
	}

	return ln
}
