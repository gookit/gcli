package cmd

import (
	"github.com/gookit/cliapp"
	"github.com/gookit/cliapp/interact"
	"github.com/gookit/color"
)

// InteractDemoCommand create
func InteractDemoCommand() *cliapp.Command {
	c := &cliapp.Command{
		Name: "interact",
		Func: interactDemo,
		// Aliases: []string{"im"}
		Description: "the command will show some interactive methods",
		Examples: `{$fullCmd} confirm
  {$fullCmd} select
`,
		Help: `Supported interactive methods:
answerIsYes 	check user answer is Yes
confirm 		confirm message
select			select one from multi options
`,
	}

	c.AddArg("name", "want running interact method name", true)
	return c
}

var funcMap = map[string]func(c *cliapp.Command){
	"select":      demoSelect,
	"confirm":     demoConfirm,
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
		map[string]string{"a":"chengdu", "b":"beijing", "c":"shanghai"},
		"a",
	)
	color.Comment.Println("your select is: ", ans)
}

func demoConfirm(_ *cliapp.Command) {

}

func demoAnswerIsYes(_ *cliapp.Command) {

}
