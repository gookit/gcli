package interact

import (
	"strings"
	"fmt"
)

// check user inputted answer is right
// fmt.Print("are you OK? [yes|no]")
// ok := EnsureUserAnswer()
func EnsureAnswerIsOk() bool {
	var answer string
	_, err := fmt.Scanln(&answer)
	if err != nil {
		panic(err)
	}

	answer = strings.TrimSpace(answer)

	if len(answer) > 0 {
		fChar := strings.ToLower(string(answer[0]))

		if fChar == "y" {
			return true
		} else if fChar == "n" {
			return false
		}
	}

	fmt.Println("Please type yes or no and then press enter:")
	return EnsureAnswerIsOk()
}
