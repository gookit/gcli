package cmd

import (
	"fmt"
	"os/exec"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v2"
	"github.com/gookit/gcli/v2/interact"
)

// InteractDemoCommand create
func InteractDemoCommand() *gcli.Command {
	c := &gcli.Command{
		Name:    "interact",
		Func:    interactDemo,
		Aliases: []string{"itt"},
		UseFor:  "the command will show some interactive methods",
		Examples: `{$fullCmd} confirm
  {$fullCmd} select
`,
		Help: `
Supported interactive methods:
  answerIsYes    check user answer is Yes
  confirm        confirm message
  select         select one from given options
  password       read user hidden input
  multiSelect    select multi from given options
`,
	}

	c.AddArg("name", "want running interact method name", true)
	return c
}

var funcMap = map[string]func(c *gcli.Command){
	"select":   demoSelect,
	"confirm":  demoConfirm,
	"password": demoPassword,

	"ms": demoMultiSelect,

	"multiSelect": demoMultiSelect,
	"answerIsYes": demoAnswerIsYes,
}

func interactDemo(c *gcli.Command, _ []string) error {
	name := c.Arg("name").String()
	if handler, ok := funcMap[name]; ok {
		handler(c)
	} else {
		return c.Errorf("want run unknown demo method: %s", name)
	}

	return nil
}

func demoSelect(_ *gcli.Command) {
	color.Green.Println("Run Select Demo")
	fmt.Println("----------------------------------------------------------")

	// s := interact.NewSelect("Your city", []string{"chengdu", "beijing", "shanghai"})
	// s.DefOpt = "2"
	// val := s.Run()
	// color.Comment.Println("your select is: ", val.String())

	ans := interact.SingleSelect(
		"Your city name(use array)?",
		[]string{"chengdu", "beijing", "shanghai"},
		"",
	)
	color.Comment.Println("your select is: ", ans)
	fmt.Println("----------------------------------------------------------")

	ans1 := interact.Choice(
		"Your age(use int array)?",
		[]int{23, 34, 45},
		"",
	)
	color.Comment.Println("your select is: ", ans1)

	fmt.Println("----------------------------------------------------------")

	ans2 := interact.SingleSelect(
		"Your city name(use map)?",
		map[string]string{"a": "chengdu", "b": "beijing", "c": "shanghai"},
		"a",
	)
	color.Comment.Println("your select is: ", ans2)
}

func demoMultiSelect(_ *gcli.Command) {
	color.Green.Println("run MultiSelect demo")

	ans := interact.MultiSelect(
		"Your city name(use array)?",
		[]string{"chengdu", "beijing", "shanghai"},
		nil,
	)
	color.Comment.Println("your select is: ", ans)

	ans2 := interact.Checkbox(
		"Your city name(use map)?",
		map[string]string{"a": "chengdu", "b": "beijing", "c": "shanghai"},
		[]string{"a"},
	)
	color.Comment.Println("your select is: ", ans2)
}

func demoConfirm(_ *gcli.Command) {
	color.Green.Println("run Confirm demo")

}

func demoPassword(_ *gcli.Command) {
	color.Green.Println("run ReadPassword demo")
	// hiddenInputTest()
	// return
	// pwd := interact.GetHiddenInput("Enter Password:", true)
	// color.Comment.Println("you input password is: ", pwd)

	pwd := interact.ReadPassword()
	color.Comment.Println("you input password is: ", pwd)
}

func hiddenInputTest() {
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

func demoAnswerIsYes(_ *gcli.Command) {

}
