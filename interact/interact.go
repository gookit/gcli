package interact

import (
	"io"
	"os"

	"github.com/gookit/color"
)

// the global input output stream
var (
	Input  io.Reader = os.Stdin
	Output io.Writer = os.Stdout
)

// SetInput stream
func SetInput(in io.Reader) { Input = in }

// SetOutput stream
func SetOutput(out io.Writer) { Output = out }

// ResetIO stream
func ResetIO() {
	Input = os.Stdin
	Output = os.Stdout
}

// Interactive definition
type Interactive struct {
	Name string
}

// New Interactive instance
func New(name string) *Interactive {
	return &Interactive{Name: name}
}

// Options definition
type Options struct {
	Quit bool
	// default value
	DefVal string
}

// Confirm a question, returns bool
func Confirm(message string, defVal ...bool) bool {
	color.Print(message)
	return AnswerIsYes(defVal...)
}

// Unconfirmed a question, returns bool
func Unconfirmed(message string, defVal ...bool) bool {
	return !Confirm(message, defVal...)
}

// Ask a question and return the result of the input.
//
// Usage:
//
//	answer := Ask("Your name?", "", nil)
//	answer := Ask("Your name?", "tom", nil)
//	answer := Ask("Your name?", "", nil, 3)
func Ask(question, defVal string, fn func(ans string) error, maxTimes ...int) string {
	q := &Question{Q: question, Func: fn, DefVal: defVal}
	if len(maxTimes) > 0 {
		q.MaxTimes = maxTimes[0]
	}

	return q.Run().String()
}

// Query is alias of method Ask()
func Query(question, defVal string, fn func(ans string) error, maxTimes ...int) string {
	return Ask(question, defVal, fn, maxTimes...)
}

// Choice is alias of method SelectOne()
func Choice(title string, options any, defOpt string, allowQuit ...bool) string {
	return SelectOne(title, options, defOpt, allowQuit...)
}

// SingleSelect is alias of method SelectOne()
func SingleSelect(title string, options any, defOpt string, allowQuit ...bool) string {
	return SelectOne(title, options, defOpt, allowQuit...)
}

// SelectOne select one of the options, returns selected option value
//
// Map options:
//
//	{
//	   	// option value => option name
//	   	'a' => 'chengdu',
//	   	'b' => 'beijing'
//	}
//
// Array options:
//
//	{
//	   // only name, value will use index
//	   'chengdu',
//	   'beijing'
//	}
func SelectOne(title string, options any, defOpt string, allowQuit ...bool) string {
	s := &Select{Title: title, Options: options, DefOpt: defOpt}

	if len(allowQuit) > 0 {
		s.DisableQuit = !allowQuit[0]
	}

	return s.Run().String()
}

// Checkbox is alias of method MultiSelect()
func Checkbox(title string, options any, defOpts []string, allowQuit ...bool) []string {
	return MultiSelect(title, options, defOpts, allowQuit...)
}

// MultiSelect select multi of the options, returns selected option values.
// like SingleSelect(), but allow select multi option
func MultiSelect(title string, options any, defOpts []string, allowQuit ...bool) []string {
	s := &Select{Title: title, Options: options, DefOpts: defOpts, MultiSelect: true}

	if len(allowQuit) > 0 {
		s.DisableQuit = !allowQuit[0]
	}

	return s.Run().Strings()
}
