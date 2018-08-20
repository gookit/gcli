package interact

import (
	"fmt"
	"github.com/gookit/cliapp/utils"
	"github.com/gookit/color"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
)

// Confirm a question, returns bool
func Confirm(message string, defVal ...bool) bool {
	color.Print(message)
	return AnswerIsYes(defVal...)
}

// Ask a question and return the result of the input.
// Usage:
// 	answer := Ask("Your name?", "", nil)
// 	answer := Ask("Your name?", "tom", nil)
// 	answer := Ask("Your name?", "", nil, 3)
func Ask(question, defVal string, fn func(ans string) (errMsg string), maxTimes ...int) string {
	q := &Question{Q: question, Func: fn, DefVal: defVal}
	if len(maxTimes) > 0 {
		q.MaxTimes = maxTimes[0]
	}

	return q.Run().String()
}

// AskQuestion is alias of method Ask()
func AskQuestion(question, defVal string, fn func(ans string) (errMsg string), maxTimes ...int) string {
	return Ask(question, defVal, fn, maxTimes...)
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

// Choice is alias of method QuickSelect()
func Choice(title string, options interface{}, defOpt string, allowQuit ...bool) string {
	return QuickSelect(title, options, defOpt, allowQuit...)
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
