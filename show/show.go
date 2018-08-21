package show

import (
	"encoding/json"
	"fmt"
	"github.com/gookit/color"
)

// Error tips message print
func Error(format string, v ...interface{}) int {
	color.Red.Print("ERROR: ")
	fmt.Printf(format+"\n", v...)
	return ERR
}

// Success tips message print
func Success(format string, v ...interface{}) int {
	color.Green.Print("SUCCESS: ")
	fmt.Printf(format+"\n", v...)
	return OK
}

// JSON print pretty JSON data
func JSON(v interface{}, settings ...string) int {
	prefix := ""
	indent := "    "

	l := len(settings)
	if l > 0 {
		prefix = settings[0]
		if l > 1 {
			indent = settings[1]
		}
	}

	bs, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bs))
	return OK
}

// AList show a list data
func AList() {

}

// AList show multi list data
func MList() {

}
