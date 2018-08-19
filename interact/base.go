// Package interact collect some interactive methods for CLI
package interact

import (
	"fmt"
	"github.com/gookit/color"
	"os"
)

const (
	// OK success exit code
	OK = 0
	// ERR error exit code
	ERR = 2
)

type Interactive struct {
	Name string
}

func New(name string) *Interactive {
	return &Interactive{Name: name}
}

func exitWithErr(format string, v ...interface{}) {
	fmt.Println(color.Red.Render("ERROR:"), fmt.Sprintf(format, v...))
	os.Exit(2)
}

func exitWithMsg(exitCode int, messages ...interface{}) {
	fmt.Println(messages...)
	os.Exit(exitCode)
}

func intsToMap(is []int) map[string]string {
	ms := make(map[string]string, len(is))
	for i, val := range is {
		k := fmt.Sprint(i)
		ms[k] = fmt.Sprint(val)
	}

	return ms
}

func stringsToMap(ss []string) map[string]string {
	ms := make(map[string]string, len(ss))
	for i, val := range ss {
		k := fmt.Sprint(i)
		ms[k] = val
	}

	return ms
}
