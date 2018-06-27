package interact

import (
	"bufio"
	"os"
	"strings"
	"fmt"
	"github.com/golangkit/cliapp/color"
)

// ReadLine
// in := ReadLine("")
// ans := ReadLine("your name?")
func ReadLine(question string) string {
	reader := bufio.NewReader(os.Stdin)

	if len(question) > 0 {
		color.Print(question)
	}

	answer, _, _ := reader.ReadLine()

	return strings.TrimSpace(string(answer))
}

// ReadFirst read first char
func ReadFirst(question string) string {
	answer := ReadLine(question)

	return string(answer[0])
}

// check user inputted answer is right
// fmt.Print("are you OK? ")
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
