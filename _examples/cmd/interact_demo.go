package cmd

import (
	"fmt"
	"os/exec"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/interact"
	"github.com/gookit/gcli/v3/show/emoji"
)

var InteractDemo = &gcli.Command{
	Name:    "interact",
	Func:    interactDemo,
	Aliases: []string{"itt"},
	Config: func(c *gcli.Command) {
		c.AddArg("name", "want running interact method name", true)
	},
	Desc: "the command will show some interactive methods",
	Examples: `{$fullCmd} confirm
  {$fullCmd} select
`,
	Help: `
Supported interactive methods:
  read           read user input text
  answerIsYes    check user answer is Yes
  confirm        confirm message
  select         select one from given options
  password       read user hidden input
  multiSelect    select multi from given options
`,
}

var funcMap = map[string]func(c *gcli.Command){
	"read":     demoReadInput,
	"select":   demoSelect,
	"confirm":  demoConfirm,
	"password": demoPassword,

	"ms": demoMultiSelect,

	"multiSelect": demoMultiSelect,
	"answerIsYes": demoAnswerIsYes,
}

func demoReadInput(c *gcli.Command) {
	ans, _ := interact.ReadLine("Your name?")

	if ans != "" {
		color.Println("Your input: ", ans)
	} else {
		color.Cyan.Println("No input!")
	}
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
	color.Green.Println("This's An Select Demo")
	fmt.Println("----------------------------------------------------------")

	ans := interact.SelectOne(
		"Your city name(use string slice/array)?",
		[]string{"chengdu", "beijing", "shanghai"},
		"",
	)
	color.Info.Println("your select is:", ans)
	fmt.Println("----------------------------------------------------------")

	ans1 := interact.Choice(
		"Your age(use int slice/array)?",
		[]int{23, 34, 45},
		"",
	)
	color.Info.Println("your select is:", ans1)

	fmt.Println("----------------------------------------------------------")

	ans2 := interact.SingleSelect(
		"Your city name(use map)?",
		map[string]string{"a": "chengdu", "b": "beijing", "c": "shanghai"},
		"a",
	)
	color.Info.Println("your select is:", ans2)

	s := interact.NewSelect("Your city", []string{"chengdu", "beijing", "shanghai"})
	s.DefOpt = "2"
	r := s.Run()
	color.Info.Println("your select key:", r.K.String())
	color.Info.Println("your select val:", r.String())
}

func demoMultiSelect(_ *gcli.Command) {
	color.Green.Println("This's An MultiSelect Demo")

	ans := interact.MultiSelect(
		"Your city name(use array)?",
		[]string{"chengdu", "beijing", "shanghai"},
		nil,
	)
	color.Comment.Println("your select is: ", ans)
	fmt.Println("----------------------------------------------------------")

	ans2 := interact.Checkbox(
		"Your city name(use map)?",
		map[string]string{"a": "chengdu", "b": "beijing", "c": "shanghai"},
		[]string{"a"},
	)
	color.Comment.Println("your select is:", ans2)
}

func demoConfirm(_ *gcli.Command) {
	color.Green.Println("This's An Confirm Demo")

	if interact.Confirm("Ensure continue") {
		fmt.Println(emoji.Render(":smile: Confirmed"))
	} else {
		color.Warn.Println("Unconfirmed")
	}
}

func demoPassword(_ *gcli.Command) {
	color.Green.Println("This's An ReadPassword Demo")
	// hiddenInputTest()
	// return
	// pwd := interact.GetHiddenInput("Enter Password:", true)
	// color.Comment.Println("you input password is: ", pwd)

	pwd := interact.ReadPassword()
	color.Comment.Println("Your input password is:", pwd)
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

func demoQuestion(_ *gcli.Command) {
	ans := interact.Ask("Your name? ", "", nil, 3)
	color.Comment.Println("Your answer is:", ans)
}
