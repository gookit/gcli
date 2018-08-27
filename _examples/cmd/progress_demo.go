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
			c.IntOpt(&pd.maxSteps, "max-step", "", 110, "setting the max step value")
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

func txtProgressBar(maxStep int)  {
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

func imgProgressBar(maxStep int)  {
	img := progress.Bar(maxStep)
	// img.Overwrite = false
	img.Chars.Processing = ' '
	// img.AddMessage("message", " handling ...")
	// use dynamic message
	img.AddWidget("message", func(p *progress.Progress) string {
		var message string
		percent := int(p.Percent() * 100)
		if percent < 20 {
			message = " Prepare ..."
		} else if percent < 40 {
			message = " Request ..."
		} else if percent < 60 {
			message = " Transport ..."
		} else if percent < 95 {
			message = " Saving ..."
		} else {
			message = " Complete."
		}

		return message
	})
	img.Start()

	// running
	for i := 0; i < maxStep; i++ {
		time.Sleep(150 * time.Millisecond)
		img.Advance()
	}

	img.Finish()
}
