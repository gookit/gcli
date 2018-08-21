package interact

import (
	"bufio"
	"fmt"
	"github.com/gookit/cliapp/utils"
	"github.com/gookit/color"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
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

// ReadPassword from terminal
func ReadPassword(message ...string) string {
	if len(message) > 0 {
		print(message[0])
	} else {
		print("Enter Password: ")
	}

	bs, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return ""
	}

	println() // new line
	return string(bs)
}

// GetHiddenInput interactively prompts for input without echoing to the terminal.
// usage:
// 	// askPassword
// 	pwd := GetHiddenInput("Enter Password:")
func GetHiddenInput(message string, trimmed bool) string {
	var err error
	var input string
	var hasResult bool

	// like *nix, git-bash ...
	if utils.HasShellEnv("sh") {
		// COMMAND: sh -c 'read -p "Enter Password:" -s user_input && echo $user_input'
		cmd := fmt.Sprintf(`'read -p "%s" -s user_input && echo $user_input'`, message)
		input, err = utils.ShellExec(cmd)
		if err != nil {
			fmt.Println("error:", err)
			return ""
		}

		println() // new line
		hasResult = true
	} else if utils.IsWin() { // at windows cmd.exe
		// create a temp VB script file
		vbFile, err := ioutil.TempFile("", "cliapp")
		if err != nil {
			return ""
		}
		defer func() {
			// delete file
			vbFile.Close()
			os.Remove(vbFile.Name())
		}()

		script := fmt.Sprintf(`wscript.echo(InputBox("%s", "", "password here"))`, message)
		vbFile.WriteString(script)
		hasResult = true

		// exec VB script
		// COMMAND: cscript //nologo vbFile.Name()
		input, err = utils.ExecCmd("cscript", []string{"//nologo", vbFile.Name()})
		if err != nil {
			return ""
		}
	}

	if hasResult {
		if trimmed {
			return strings.TrimSpace(input)
		}
		return input
	}

	panic("current env is not support the method")
}
