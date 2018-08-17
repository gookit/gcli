package interact

import (
	"fmt"
	"strings"
	"github.com/gookit/color"
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
