package cmd

import (
	"github.com/gookit/cliapp"
	"github.com/gookit/cliapp/progress"
	"time"
)

type progressDemo struct {
	maxSteps int
}

func ProgressDemoCmd() *cliapp.Command {
	pd := &progressDemo{}

	return &cliapp.Command{
		Name: "prog",
		UseFor: "there are some progress bar run demos",
		Aliases: []string{"prg:demo", "progress"},
		Func: pd.Run,
		Config: func(c *cliapp.Command) {
			c.IntOpt(&pd.maxSteps, "max-step", "", 111, "setting the max step value")
			c.AddArg("name", "the progress bar type name. allow: bar,txt", true)
		},
		Examples:`Text progress bar:
  {$fullCmd} txt
Image progress bar:
  {$fullCmd} bar`,
	}
}

// Run command
func (d *progressDemo) Run(c *cliapp.Command, _ []string) int {
	name := c.Arg("name").String()
	max := d.maxSteps

	switch name {
	case "bar":
		imgProgressBar(max)
	case "txt":
		txtProgressBar(max)
	default:
		return c.Errorf("the progress bar type name only allow: bar,txt. input is: %s", name)
	}
	return 0
}

func txtProgressBar(max int)  {
	maxStep := 110
	txt := progress.Txt(maxStep)
	txt.AddMessage("message", "handling ... ")
	// txt.Overwrite = false
	txt.Start()

	// running
	for i := 0; i < maxStep; i++ {
		time.Sleep(200 * time.Millisecond)
		txt.Advance()
	}

	txt.Finish()
}

func imgProgressBar(max int)  {
	maxStep := 110
	img := progress.Bar(maxStep)
	img.AddMessage("message", " handling ...")
	// txt.Overwrite = false
	img.Start()

	// running
	for i := 0; i < maxStep; i++ {
		time.Sleep(200 * time.Millisecond)
		img.Advance()
	}

	img.Finish()
}
