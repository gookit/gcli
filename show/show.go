package show

import (
	"encoding/json"
	"fmt"
)

// New a
func New() {
	//
}

// JSON print pretty JSON data
func JSON(v interface{}, settings ...string) {
	l := len(settings)
	prefix := ""
	indent := "    "

	if l > 1 {
		prefix = settings[0]
		indent = settings[1]
	} else if l > 0 {
		prefix = settings[0]
	}

	bs, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bs))
}
