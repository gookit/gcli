package interact

import (
	"fmt"
	"strconv"
)

type Option struct {
	Quit bool
	// default value
	DefVal string
}

// Value data store
type Value struct {
	val interface{}
}

// Val get
func (v Value) Val() interface{} {
	return v.val
}

// Int convert
func (v Value) Int() (val int) {
	if v.val == nil {
		return
	}
	switch tpVal := v.val.(type) {
	case int:
		return tpVal
	case string:
		val, err := strconv.Atoi(tpVal)
		if err == nil {
			return val
		}
	}
	return
}

// String convert
func (v Value) String() string {
	if v.val == nil {
		return ""
	}

	return fmt.Sprintf("%v", v.val)
}

// IsEmpty value
func (v Value) IsEmpty() bool {
	return v.val == nil
}

func Checkbox() {

}

func MultiSelect(message string) {

}
