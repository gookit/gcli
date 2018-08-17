package interact

import (
	"bufio"
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
	return string(answer[0]), err
}
