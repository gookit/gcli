package cmd

import (
	"github.com/gookit/cliapp"
	"github.com/gookit/cliapp/interact"
	"github.com/gookit/color"
	"fmt"
	"os/exec"
)

// InteractDemoCommand create
func InteractDemoCommand() *cliapp.Command {
	c := &cliapp.Command{
		Name: "interact",
		Func: interactDemo,
		Aliases: []string{"itt"},
		Description: "the command will show some interactive methods",
		Examples: `{$fullCmd} confirm
  {$fullCmd} select
`,
		Help: `Supported interactive methods:
answerIsYes 	check user answer is Yes
confirm 		confirm message
select			select one from multi options
password		read user hidden input
`,
	}

	c.AddArg("name", "want running interact method name", true)
	return c
}

var funcMap = map[string]func(c *cliapp.Command){
	"select":   demoSelect,
	"confirm":  demoConfirm,
	"password": demoPassword,

	"answerIsYes": demoAnswerIsYes,
}

func interactDemo(c *cliapp.Command, _ []string) int {
	name := c.Arg("name").String()
	if handler, ok := funcMap[name]; ok {
		handler(c)
	} else {
		return c.Errorf("want run unknown demo method: %s", name)
	}

	return 0
}

func demoSelect(_ *cliapp.Command) {
	// s := interact.NewSelect("Your city", []string{"chengdu", "beijing", "shanghai"})
	// s.DefOpt = "2"
	// val := s.Run()
	// color.Comment.Println("your select is: ", val.String())

	ans := interact.QuickSelect(
		"Your city name(use array)?",
		[]string{"chengdu", "beijing", "shanghai"},
		"",
	)
	color.Comment.Println("your select is: ", ans)

	ans = interact.QuickSelect(
		"Your city name(use map)?",
		map[string]string{"a": "chengdu", "b": "beijing", "c": "shanghai"},
		"a",
	)
	color.Comment.Println("your select is: ", ans)
}

func demoConfirm(_ *cliapp.Command) {

}

func demoPassword(_ *cliapp.Command) {
	// hiddenInputTest()
	// return
	// pwd := interact.GetHiddenInput("Enter Password:", true)
	// color.Comment.Println("you input password is: ", pwd)

	pwd := interact.ReadPassword()
	color.Comment.Println("you input password is: ", pwd)
}

func hiddenInputTest()  {
	// COMMAND: sh -c 'read -p "Enter Password:" -s user_input && echo $user_input'
	// str := fmt.Sprintf(`'read -p "%s" -s user_input && echo $user_input'`, "Enter Password:")
	// cmd := exec.CommandContext()
	cmd := exec.Command("sh", "-c", `read -p "Enter Password:" -s user_input && echo $user_input`)
	err := cmd.Start()
	fmt.Println("start", err)
	err = cmd.Wait()
	fmt.Println("wait", err, cmd.Process.Pid, cmd.ProcessState.Pid())

	cmd = exec.Command("sh", "./read-pwd.sh")
	bs, err := cmd.Output()
	fmt.Println(string(bs), err)
}

func demoAnswerIsYes(_ *cliapp.Command) {

}
