// Package interact collect some interactive methods for CLI
package interact

import (
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/goutil/structs"
)

const (
	// OK success exit code
	OK = 0
	// ERR error exit code
	ERR = 2
)

// ComOptions struct
type ComOptions struct {
	// ValidFn check input value
	ValidFn func(val any) (any, error)
}

// Value alias of structs.Value
type Value = structs.Value

// RunFace for interact methods
type RunFace interface {
	Run() *Value
}

/*************************************************************
 * value for select
 *************************************************************/

// SelectResult data store
type SelectResult struct {
	Value // V the select value(s)
	// K the select key(s)
	K Value
}

// create SelectResult create
func newSelectResult(key, val any) *SelectResult {
	return &SelectResult{
		K:     Value{V: key},
		Value: Value{V: val},
	}
}

// KeyString get
func (sv *SelectResult) KeyString() string {
	return sv.K.String()
}

// KeyStrings get
func (sv *SelectResult) KeyStrings() []string {
	return sv.K.Strings()
}

// Key value get
func (sv *SelectResult) Key() any {
	return sv.K.Val()
}

// WithKey value
func (sv *SelectResult) WithKey(key any) *SelectResult {
	sv.K.Set(key)
	return sv
}

/*************************************************************
 * helper methods
 *************************************************************/

func exitWithErr(format string, v ...any) {
	color.Error.Tips(format, v...)
	os.Exit(ERR)
}

func exitWithMsg(exitCode int, messages ...any) {
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

func stringToArr(str, sep string) (arr []string) {
	str = strings.TrimSpace(str)
	ss := strings.Split(str, sep)
	for _, val := range ss {
		if val = strings.TrimSpace(val); val != "" {
			arr = append(arr, val)
		}
	}

	return arr
}

func stringsToMap(ss []string) map[string]string {
	ms := make(map[string]string, len(ss))
	for i, val := range ss {
		k := fmt.Sprint(i)
		ms[k] = val
	}

	return ms
}
