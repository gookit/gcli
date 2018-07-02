package interact

import (
	"strings"
	"fmt"
)

// check user inputted answer is right
// fmt.Print("are you OK?")
// ok := EnsureAnswerIsOk()
// ok := EnsureAnswerIsOk(true)
func EnsureAnswerIsOk(def ...bool) bool {
	var answer string

	// _, err := fmt.Scanln(&answer)
	// _, err := fmt.Scan(&answer)
	answer, err := ReadLine(" [yes|no]: ")

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
	return EnsureAnswerIsOk()
}
