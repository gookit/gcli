package interact

import (
	"bufio"
	"fmt"
	"github.com/gookit/color"
	"os"
	"strings"
)

// ReadLine read line for user input
// in := ReadLine("")
// ans := ReadLine("your name?")
func ReadLine(question string) (string, error) {
	if len(question) > 0 {
		color.Print(question)
	}

	reader := bufio.NewReader(os.Stdin)
	answer, _, err := reader.ReadLine()
	return strings.TrimSpace(string(answer)), err
}

// ReadFirst read first char
func ReadFirst(question string) (string, error) {
	answer, err := ReadLine(question)
	if len(answer) == 0 {
		return "", err
	}

	return string(answer[0]), err
}

// AnswerIsYes check user inputted answer is right
// fmt.Print("are you OK?")
// ok := AnswerIsYes()
// ok := AnswerIsYes(true)
func AnswerIsYes(def ...bool) bool {
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
	fChar, err := ReadFirst(mark)
	if err != nil {
		panic(err)
	}

	if len(fChar) > 0 {
		fChar := strings.ToLower(fChar)
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
