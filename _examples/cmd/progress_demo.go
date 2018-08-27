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
		Name:    "prog",
		UseFor:  "there are some progress bar run demos",
		Aliases: []string{"prg:demo", "progress"},
		Func:    pd.Run,
		Config: func(c *cliapp.Command) {
			c.IntOpt(&pd.maxSteps, "max-step", "", 110, "setting the max step value")
			c.AddArg("name",
				"progress bar type name. allow: bar,txt,loading,roundTrip",
				true,
			)
		},
		Examples: `Text progress bar:
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
	case "loading":
		runLoadingBar(max)
	case "rt", "roundTrip":
		runRoundTripBar(max)
	default:
		return c.Errorf("the progress bar type name only allow: bar,txt. input is: %s", name)
	}
	return 0
}

func runRoundTripBar(max int) {
	p := progress.RoundTrip(0).WithMaxSteps(max)

	// running
	runProgressBar(p, max, 120)

	p.Finish()
}

func txtProgressBar(maxStep int) {
	txt := progress.Txt(maxStep)
	txt.AddMessage("message", "handling ... ")
	// txt.Overwrite = false

	// running
	runProgressBar(txt, maxStep, 80)

	txt.Finish()
}

func imgProgressBar(maxStep int) {
	p := progress.Bar(maxStep)
	// p.Overwrite = false
	p.Chars.Completed = progress.CharWell
	p.Chars.Processing = '>'
	p.Chars.Remaining = '-'
	// p.AddMessage("message", " handling ...")
	// use dynamic message
	p.AddWidget("message", func(p *progress.Progress) string {
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

	// running
	runProgressBar(p, maxStep, 90)

	p.Finish()
}

func runLoadingBar(maxStep int) {
	// chars := []rune(`∷∵∴∶`)
	p := progress.LoadingBar(progress.LoadingTheme1)
	p.MaxSteps = uint(maxStep)
	p.AddMessage("message", "data loading ... ...")

	// running
	runProgressBar(p, maxStep, 70)

	// p.Finish()
	p.Finish("data load complete")
}

// running
func runProgressBar(p progress.ProgressFace, maxStep int, speed int) {
	p.Start()
	for i := 0; i < maxStep; i++ {
		time.Sleep(time.Duration(speed) * time.Millisecond)
		p.Advance()
	}
}
