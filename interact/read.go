package interact

import (
	"bufio"
	"os"
	"github.com/gookit/cliapp/color"
	"strings"
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
