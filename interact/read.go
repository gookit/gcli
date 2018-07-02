package interact

import (
	"bufio"
	"os"
	"github.com/gookit/color"
	"strings"
)

// ReadLine read line for user input
// in := ReadLine("")
// ans := ReadLine("your name?")
func ReadLine(question string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	if len(question) > 0 {
		color.Print(question)
	}

	answer, _, err := reader.ReadLine()

	return strings.TrimSpace(string(answer)), err
}

// ReadFirst read first char
func ReadFirst(question string) (string, error) {
	answer, err := ReadLine(question)

	return string(answer[0]), err
}
