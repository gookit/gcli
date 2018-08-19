package interact

import (
	"fmt"
	"github.com/gookit/color"
	"strings"
)

// AnswerIsYes check user inputted answer is right
// fmt.Print("are you OK?")
// ok := AnswerIsYes()
// ok := AnswerIsYes(true)
func AnswerIsYes(def ...bool) bool {
	var answer string
	mark := " [yes|no]: "

	if len(def) > 0 {
		var defShow string
		if def[0] {
			defShow = "yes"
		} else {
			defShow = "no"
		}

		mark = fmt.Sprintf(" [yes|no](default <cyan>%s</>): ", defShow)
	}

	// _, err := fmt.Scanln(&answer)
	// _, err := fmt.Scan(&answer)
	answer, err := ReadLine(mark)
	if err != nil {
		panic(err)
	}

	if len(answer) > 0 {
		fChar := strings.ToLower(string(answer[0]))
		if fChar == "y" {
			return true
		} else if fChar == "n" {
			return false
		}
	} else if len(def) > 0 {
		return def[0]
	}

	fmt.Print("Please try again")
	return AnswerIsYes()
}

// Confirm a question, returns bool
func Confirm(message string, defVal ...bool) bool {
	color.Print(message)
	return AnswerIsYes(defVal...)
}

// QuickSelect select one of the options, returns selected option value
// map options:
// 	{
//    // option value => option name
//    'a' => 'chengdu',
//    'b' => 'beijing'
// 	}
// array options:
// 	{
//    // only name, value will use index
//    'chengdu',
//    'beijing'
// 	}
func QuickSelect(title string, options interface{}, defOpt string, allowQuit ...bool) string {
	s := NewSelect(title, options)
	s.DefOpt = defOpt

	if len(allowQuit) > 0 {
		s.NoQuit = !allowQuit[0]
	}

	return s.Run().String()
}
