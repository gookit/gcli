package cmd

import (
	"github.com/gookit/cliapp"
	"github.com/gookit/cliapp/progress"
	"time"
)

type progressDemo struct {
	maxSteps int
	overwrite bool
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
			c.BoolOpt(&pd.overwrite, "overwrite", "o", true, "setting overwrite progress bar line")
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
	case "txt", "text":
		txtProgressBar(max)
	case "load", "loading":
		runLoadingBar(max)
	case "rt", "roundTrip":
		runRoundTripBar(max)
	default:
		return c.Errorf("the progress bar type name only allow: bar,txt,loading,roundTrip. input is: %s", name)
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
	cs := progress.BarCharsStyle4

	p := progress.CustomBar(60, cs)
	p.MaxSteps = uint(maxStep)
	p.Format = progress.FullBarFormat
	// p.Overwrite = false

	// p.AddMessage("message", " handling ...")
	// use dynamic message
	p.AddWidget("message", progress.DynamicTextWidget(map[int]string{
		20: " Prepare ...",
		40: " Request ...",
		65: " Transport ...",
		95: " Saving ...",
		100: " Handle Complete.",
	}))

	// running
	runProgressBar(p, maxStep, 100)
	p.Finish()
}

func runLoadingBar(maxStep int) {
	p := progress.LoadingBar(progress.LoadingTheme7)
	p.MaxSteps = uint(maxStep)
	p.AddMessage("message", "data loading ... ...")

	// running
	runProgressBar(p, maxStep, 70)

	// p.Finish()
	p.Finish("data load complete")
}

// running
func runProgressBar(p *progress.Progress, maxSteps int, speed int) {
	p.Start()
	for i := 0; i < maxSteps; i++ {
		time.Sleep(time.Duration(speed) * time.Millisecond)
		p.Advance()
	}
}
